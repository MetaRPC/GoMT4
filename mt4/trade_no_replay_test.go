package mt4

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/MetaRPC/GoMT4/package"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// fakeTrading answers the trade RPCs with a scripted result and counts calls.
type fakeTrading struct {
	pb.UnimplementedTradingHelperServer
	calls  atomic.Int32
	send   func(call int32) (*pb.OrderSendReply, error)
	failFn func() error // OrderModify / OrderCloseDelete / OrderCloseBy transport error
}

func (f *fakeTrading) OrderSend(context.Context, *pb.OrderSendRequest) (*pb.OrderSendReply, error) {
	return f.send(f.calls.Add(1))
}

func (f *fakeTrading) OrderModify(context.Context, *pb.OrderModifyRequest) (*pb.OrderModifyReply, error) {
	f.calls.Add(1)
	return nil, f.failFn()
}

func (f *fakeTrading) OrderCloseDelete(context.Context, *pb.OrderCloseDeleteRequest) (*pb.OrderCloseDeleteReply, error) {
	f.calls.Add(1)
	return nil, f.failFn()
}

func (f *fakeTrading) OrderCloseBy(context.Context, *pb.OrderCloseByRequest) (*pb.OrderCloseByReply, error) {
	f.calls.Add(1)
	return nil, f.failFn()
}

// fakeAccount fails AccountSummary once with UNAVAILABLE, then answers.
type fakeAccount struct {
	pb.UnimplementedAccountHelperServer
	calls atomic.Int32
}

func (f *fakeAccount) AccountSummary(context.Context, *pb.AccountSummaryRequest) (*pb.AccountSummaryReply, error) {
	if f.calls.Add(1) == 1 {
		return nil, status.Error(codes.Unavailable, "connection reset")
	}
	return &pb.AccountSummaryReply{Response: &pb.AccountSummaryReply_Data{Data: &pb.AccountSummaryData{AccountLogin: 42}}}, nil
}

func newFakeMT4Account(t *testing.T, trading *fakeTrading, account *fakeAccount) *MT4Account {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	if trading != nil {
		pb.RegisterTradingHelperServer(srv, trading)
	}
	if account != nil {
		pb.RegisterAccountHelperServer(srv, account)
	}
	go func() { _ = srv.Serve(lis) }()
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close(); srv.Stop() })
	return &MT4Account{
		GrpcConn:      conn,
		TradeClient:   pb.NewTradingHelperClient(conn),
		AccountClient: pb.NewAccountHelperClient(conn),
		Id:            uuid.New(),
		ApiKey:        "test-key",
	}
}

func assertTradeOutcomeUnknown(t *testing.T, err error, op string, code codes.Code) {
	t.Helper()
	var tou *TradeOutcomeUnknownError
	if !errors.As(err, &tou) {
		t.Fatalf("expected *TradeOutcomeUnknownError, got %T: %v", err, err)
	}
	if tou.Code != "TRADE_OUTCOME_UNKNOWN" || tou.ErrorCode() != "TRADE_OUTCOME_UNKNOWN" || tou.Operation != op {
		t.Fatalf("bad fields: code=%q op=%q", tou.Code, tou.Operation)
	}
	if !strings.HasPrefix(err.Error(), "TRADE_OUTCOME_UNKNOWN: "+op+" failed with gRPC ") {
		t.Fatalf("bad message: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "NOT re-sent") {
		t.Fatalf("message does not say it was not re-sent: %s", err.Error())
	}
	if got := status.Code(err); got != code {
		t.Fatalf("status.Code = %s, want %s", got, code)
	}
	if s, ok := status.FromError(err); !ok || s.Code() != code {
		t.Fatalf("status.FromError = %v,%v", s, ok)
	}
	if errors.Unwrap(tou) == nil {
		t.Fatal("Unwrap returned nil")
	}
}

func TestOrderSendUnavailableIsNotReplayed(t *testing.T) {
	f := &fakeTrading{send: func(int32) (*pb.OrderSendReply, error) {
		return nil, status.Error(codes.Unavailable, "socket closed")
	}}
	a := newFakeMT4Account(t, f, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil)
	assertTradeOutcomeUnknown(t, err, "OrderSend", codes.Unavailable)
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("OrderSend called %d times, want 1", n)
	}

	// Through the service layer too.
	_, err = NewMT4Service(a).OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil)
	assertTradeOutcomeUnknown(t, err, "OrderSend", codes.Unavailable)
	if n := f.calls.Load(); n != 2 {
		t.Fatalf("OrderSend called %d times after service call, want 2", n)
	}
}

func TestOtherTradesOutcomeUnknownCodes(t *testing.T) {
	for _, c := range []codes.Code{codes.Unavailable, codes.DeadlineExceeded, codes.Unknown, codes.Internal} {
		f := &fakeTrading{failFn: func() error { return status.Error(c, "boom") }}
		a := newFakeMT4Account(t, f, nil)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := a.OrderModify(ctx, 1, nil, nil, nil, nil)
		assertTradeOutcomeUnknown(t, err, "OrderModify", c)
		lots := 0.01
		_, err = a.OrderClose(ctx, 1, &lots, nil, nil)
		assertTradeOutcomeUnknown(t, err, "OrderCloseDelete", c)
		_, err = a.OrderDelete(ctx, 2)
		assertTradeOutcomeUnknown(t, err, "OrderCloseDelete", c)
		_, err = a.OrderCloseBy(ctx, 1, 2)
		assertTradeOutcomeUnknown(t, err, "OrderCloseBy", c)
		cancel()
		if n := f.calls.Load(); n != 4 {
			t.Fatalf("%s: trade calls = %d, want 4 (one each)", c, n)
		}
	}
}

func TestOrderSendTerminalNotFoundIsResent(t *testing.T) {
	f := &fakeTrading{send: func(call int32) (*pb.OrderSendReply, error) {
		if call == 1 {
			return &pb.OrderSendReply{Response: &pb.OrderSendReply_Error{Error: &pb.Error{ErrorCode: "TERMINAL_INSTANCE_NOT_FOUND"}}}, nil
		}
		return &pb.OrderSendReply{Response: &pb.OrderSendReply_Data{Data: &pb.OrderSendData{Ticket: 7}}}, nil
	}}
	a := newFakeMT4Account(t, f, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := a.OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.GetTicket() != 7 {
		t.Fatalf("unexpected data: %v", data)
	}
	if n := f.calls.Load(); n != 2 {
		t.Fatalf("OrderSend called %d times, want 2", n)
	}
}

func TestReadUnavailableIsRetried(t *testing.T) {
	acc := &fakeAccount{}
	a := newFakeMT4Account(t, nil, acc)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := a.AccountSummary(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.GetAccountLogin() != 42 {
		t.Fatalf("unexpected data: %v", data)
	}
	if n := acc.calls.Load(); n != 2 {
		t.Fatalf("AccountSummary called %d times, want 2", n)
	}
}

func TestApiKeyInvalidFailsFast(t *testing.T) {
	f := &fakeTrading{send: func(int32) (*pb.OrderSendReply, error) {
		return &pb.OrderSendReply{Response: &pb.OrderSendReply_Error{Error: &pb.Error{ErrorCode: "API_KEY_INVALID", ErrorMessage: "unknown api key"}}}, nil
	}}
	a := newFakeMT4Account(t, f, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := a.OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "API_KEY_INVALID") {
		t.Fatalf("expected API_KEY_INVALID error, got %v", err)
	}
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("OrderSend called %d times, want 1", n)
	}

	// Refused key carried in a transport status message: no retry for reads either.
	g := &fakeTrading{send: func(int32) (*pb.OrderSendReply, error) {
		return nil, status.Error(codes.Unavailable, "API_KEY_REQUIRED: no api key")
	}}
	b := newFakeMT4Account(t, g, nil)
	_, err = ExecuteWithReconnect(b, ctx, func(md metadata.MD) (*pb.OrderSendReply, error) {
		return b.TradeClient.OrderSend(metadata.NewOutgoingContext(ctx, md), &pb.OrderSendRequest{})
	}, func(r *pb.OrderSendReply) *pb.Error { return r.GetError() })
	if err == nil || !strings.Contains(err.Error(), "API_KEY_REQUIRED") {
		t.Fatalf("expected API_KEY_REQUIRED error, got %v", err)
	}
	if n := g.calls.Load(); n != 1 {
		t.Fatalf("call made %d times, want 1", n)
	}
}
