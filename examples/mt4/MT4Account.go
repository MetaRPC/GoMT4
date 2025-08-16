package mt4


import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"time"


	pb "git.mtapi.io/root/mrpc-proto.git/mt4/libraries/go"


	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MT4Account represents a client session for interacting with the MT4 terminal API over gRPC.
type MT4Account struct {
	
	// User is the MT4 account login number.
	User uint64

	// Password for the user account.
	Password string

	// Host is the IP/domain of the MT4 server.
	Host string

	// Port is the MT4 server port (typically 443).
	Port int

	// ServerName is the MT4 server name (used for cluster connections).
	ServerName string

	// BaseChartSymbol is the default chart symbol (e.g., "EURUSD").
	BaseChartSymbol string

	// ConnectTimeout is the timeout for connection readiness, in seconds.
	ConnectTimeout int

	// GrpcServer is the address (host:port) of the gRPC API endpoint.
	GrpcServer string

	// GrpcConn is the underlying gRPC client connection (must be closed when done).
	GrpcConn *grpc.ClientConn

	// Per-service gRPC API clients.
	ConnectionClient   pb.ConnectionClient
	SubscriptionClient pb.SubscriptionServiceClient
	AccountClient      pb.AccountHelperClient
	TradeClient        pb.TradingHelperClient
	MarketInfoClient   pb.MarketInfoClient
    AccountHelper      pb.AccountHelperClient
	// Id is a unique identifier (UUID) for this account session/instance.
	Id uuid.UUID
}

// NewMT4Account initializes a new MT4Account and establishes the underlying gRPC connection.
// Returns a pointer to the account object and any error encountered while connecting.
func NewMT4Account(user uint64, password string, grpcServer string, id uuid.UUID) (*MT4Account, error) {
	// If no endpoint specified, use production default
	if grpcServer == "" {
		grpcServer = "mt4.mrpc.pro:443"
	}

	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return nil, err
	}

	// Instantiate API service clients using the shared gRPC connection
	return &MT4Account{
		User:               user,
		Password:           password,
		GrpcServer:         grpcServer,
		GrpcConn:           conn,
		ConnectionClient:   pb.NewConnectionClient(conn),
		SubscriptionClient: pb.NewSubscriptionServiceClient(conn),
		AccountClient:      pb.NewAccountHelperClient(conn),
		TradeClient:        pb.NewTradingHelperClient(conn),
		MarketInfoClient:   pb.NewMarketInfoClient(conn),
		Id:                 id,
		Port:               443,
		ConnectTimeout:     30,
	}, nil
}

// isConnected returns true if this account is associated with any host or server name.
func (a *MT4Account) isConnected() bool {
	return a.Host != "" || a.ServerName != ""
}

// getHeaders builds the gRPC metadata headers (adds "id" if present).
func (a *MT4Account) getHeaders() metadata.MD {
	if a.Id == uuid.Nil {
		return nil
	}
	return metadata.Pairs("id", a.Id.String())
}

// ConnectByHostPort connects to the MT4 terminal using a host/port pair.
// Updates the session state fields (Host, Port, etc.) upon success.
func (a *MT4Account) ConnectByHostPort(
	ctx context.Context,
	host string,
	port int,
	baseChartSymbol string,
	waitForTerminalIsAlive bool,
	timeoutSeconds int,
) error {
	// Ensure ctx is non-nil (metadata.NewOutgoingContext would panic on nil)
	if ctx == nil {
		ctx = context.Background()
	}

	// Build the protobuf request struct
	req := &pb.ConnectRequest{
		User:                                   a.User,
		Password:                               a.Password,
		Host:                                   host,
		Port:                                   int32(port),
		BaseChartSymbol:                        proto.String(baseChartSymbol),
		WaitForTerminalIsAlive:                 proto.Bool(waitForTerminalIsAlive),
		TerminalReadinessWaitingTimeoutSeconds: proto.Int32(int32(timeoutSeconds)),
	}

	// Set metadata if available
	md := a.getHeaders()
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Make the actual gRPC call
	res, err := a.ConnectionClient.Connect(ctx, req)
	if err != nil {
		return err
	}
	// API errors are delivered via .GetError()
	if res.GetError() != nil {
		return fmt.Errorf("API error: %v", res.GetError())
	}

	// Store session properties if connection is established
	a.Host = host
	a.Port = port
	a.BaseChartSymbol = baseChartSymbol
	a.ConnectTimeout = timeoutSeconds

	// Set the session UUID if present in response
	if data := res.GetData(); data != nil && data.GetTerminalInstanceGuid() != "" {
		id, _ := uuid.Parse(data.GetTerminalInstanceGuid())
		a.Id = id
	}
	return nil
}


// ConnectByServerName connects to the MT4 terminal using the cluster/server name.
func (a *MT4Account) ConnectByServerName(
	ctx context.Context,
	serverName string,
	baseChartSymbol string,
	waitForTerminalIsAlive bool,
	timeoutSeconds int,
) error {
	// Ensure ctx is non-nil (metadata.NewOutgoingContext would panic on nil)
	if ctx == nil {
		ctx = context.Background()
	}

	req := &pb.ConnectExRequest{
		User:                                   a.User,
		Password:                               a.Password,
		MtClusterName:                          serverName,
		BaseChartSymbol:                        proto.String(baseChartSymbol),
		TerminalReadinessWaitingTimeoutSeconds: proto.Int32(int32(timeoutSeconds)),
	}

	md := a.getHeaders()
	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := a.ConnectionClient.ConnectEx(ctx, req)
	if err != nil {
		return err
	}
	if res.GetError() != nil {
		return fmt.Errorf("API error: %v", res.GetError())
	}

	a.ServerName = serverName
	a.BaseChartSymbol = baseChartSymbol
	a.ConnectTimeout = timeoutSeconds

	if data := res.GetData(); data != nil && data.GetTerminalInstanceGuid() != "" {
		id, _ := uuid.Parse(data.GetTerminalInstanceGuid())
		a.Id = id
	}
	return nil
}


// ExecuteWithReconnect retries a gRPC call on recoverable errors (network/instance-not-found).
//
// T:          Type of the response object (e.g., *pb.AccountSummaryReply)
// a:          Pointer to MT4Account (used for headers, etc.)
// ctx:        Context for cancellation/deadline
// grpcCall:   Function taking metadata and returning (response, error)
// errorSelector: Function taking response and extracting *pb.Error (returns nil if no API error)
//
// Returns the response or error.
func ExecuteWithReconnect[T any](
	a *MT4Account,
	ctx context.Context,
	grpcCall func(metadata.MD) (T, error),
	errorSelector func(T) *pb.Error,
) (T, error) {
	// Ensure ctx is non-nil for <-ctx.Done() usage
	if ctx == nil {
		ctx = context.Background()
	}

	var zeroT T // Zero value of T (used for error returns)

	for {
		// Prepare gRPC headers for session (may be nil)
		headers := a.getHeaders()

		// Call the gRPC method (returns (reply, error))
		res, err := grpcCall(headers)
		if err != nil {
			// If it's a gRPC Unavailable (connection/server issue), wait and retry
			if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
				select {
				case <-time.After(500 * time.Millisecond):
					continue // Try again after delay
				case <-ctx.Done():
					return zeroT, ctx.Err() // Cancelled by caller
				}
			}
			// Other errors: return immediately
			return zeroT, err
		}

		// Check for API (business logic) error in the response
		apiErr := errorSelector(res)
		if apiErr != nil {
			// If terminal instance not found (e.g., dropped session), wait and retry
			if apiErr.GetErrorCode() == "TERMINAL_INSTANCE_NOT_FOUND" {
				select {
				case <-time.After(500 * time.Millisecond):
					continue // Try again after delay
				case <-ctx.Done():
					return zeroT, ctx.Err()
				}
			}
			// All other API errors: return as Go errors
			return zeroT, fmt.Errorf("API error: %v", apiErr)
		}

		// Success! Return the response.
		return res, nil
	}
}


//=== 📂 Account Info ===

// AccountSummary retrieves summary information about the connected MT4 trading account.
//
// Parameters:
//   - ctx: Context for timeout or cancellation (e.g., context.Background(), context.WithTimeout).
//
// Returns:
//   - Pointer to AccountSummaryData (or nil if error).
//   - Error if not connected or if the gRPC/API call fails.
//
// This method handles automatic retries on network or "terminal instance not found" errors
// using ExecuteWithReconnect. It accesses protobuf fields via generated Get...() methods.
func (a *MT4Account) AccountSummary(ctx context.Context) (*pb.AccountSummaryData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	// If caller didn't set a deadline, add a short per-call timeout to avoid hanging RPCs.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	// 2) Ensure the account is connected to a server before making a request.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// 3) Construct the empty request message (no parameters for account summary).
	req := &pb.AccountSummaryRequest{}

	// 4) gRPC call closure: attach per-call metadata to the outgoing context.
	grpcCall := func(headers metadata.MD) (*pb.AccountSummaryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.AccountClient.AccountSummary(c, req)
	}

	// 5) Application-level error selector.
	errorSelector := func(reply *pb.AccountSummaryReply) *pb.Error {
		return reply.GetError()
	}

	// 6) Execute with reconnect/retry wrapper.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// 7) Return payload.
	return reply.GetData(), nil
}


// NOTE: All values are retrieved via AccountSummary,
// which already uses ExecuteWithReconnect internally.

// AccountLogin returns the account login number.
func (a *MT4Account) AccountLogin(ctx context.Context) (int64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountLogin(), nil
}

// AccountBalance returns the account balance.
func (a *MT4Account) AccountBalance(ctx context.Context) (float64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountBalance(), nil
}

// AccountCredit returns the account credit.
func (a *MT4Account) AccountCredit(ctx context.Context) (float64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountCredit(), nil
}

// AccountEquity returns the account equity.
func (a *MT4Account) AccountEquity(ctx context.Context) (float64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountEquity(), nil
}

// AccountLeverage returns the account leverage.
func (a *MT4Account) AccountLeverage(ctx context.Context) (int64, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return 0, err
	}
	return summary.GetAccountLeverage(), nil
}

// AccountName returns the account user name.
func (a *MT4Account) AccountName(ctx context.Context) (string, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return "", err
	}
	return summary.GetAccountUserName(), nil
}

// AccountCompany returns the account company name.
func (a *MT4Account) AccountCompany(ctx context.Context) (string, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return "", err
	}
	return summary.GetAccountCompanyName(), nil
}

// AccountCurrency returns the account currency code.
func (a *MT4Account) AccountCurrency(ctx context.Context) (string, error) {
	summary, err := a.AccountSummary(ctx)
	if err != nil {
		return "", err
	}
	return summary.GetAccountCurrency(), nil
}

// ExecuteStreamWithReconnect wraps a gRPC server-streaming call with automatic reconnection
// on network and recoverable API errors, sending extracted data to a channel.
// - ctx: Context for cancellation and deadline
// - a:   Your session/account struct (for session headers etc.)
// - request: The protobuf request message
// - streamInvoker: function to open the gRPC stream (returns grpc.ClientStream, error)
// - getError: function to extract API error from a reply message (returns *pb.Error or nil)
// - getData: function to extract (data, ok) from a reply (ok=false means skip this message)
// - newReply: function that creates a new TReply instance (needed because Go generics can't new(T))
// Returns:
//   - dataCh: channel of received TData messages (e.g. *pb.OnSymbolTickData)
//   - errCh:  channel for any errors or end-of-stream events.
func ExecuteStreamWithReconnect[TRequest any, TReply any, TData any](
	ctx context.Context,
	a *MT4Account,
	request TRequest,
	streamInvoker func(TRequest, metadata.MD, context.Context) (grpc.ClientStream, error),
	getError func(TReply) *pb.Error,
	getData func(TReply) (TData, bool),
	newReply func() TReply,
) (<-chan TData, <-chan error) {
	// Ensure ctx is non-nil for stream lifecycle and selects
	if ctx == nil {
		ctx = context.Background()
	}

	// Context-aware wait helper (respects cancellation/deadlines)
waitWithCtx := func(d time.Duration) error {
    tctx, cancel := context.WithTimeout(ctx, d)
    defer cancel()
    select {
    case <-tctx.Done():
        return tctx.Err()
		}
	}

	// Small symmetric jitter around a base duration to avoid thundering herd
	jitter := func(base, plusMinus time.Duration) time.Duration {
		n := time.Now().UnixNano()
		off := time.Duration(int64(n)%int64(plusMinus*2)) - plusMinus
		return base + off
	}

	dataCh := make(chan TData)
	errCh := make(chan error, 1)

	go func() {
		defer close(dataCh)
		defer close(errCh)

		for {
			reconnectRequired := false
			headers := a.getHeaders()

			// Open the gRPC streaming call with headers and context
			stream, err := streamInvoker(request, headers, ctx)
			if err != nil {
				// If network/server unavailable, retry unless cancelled
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					if err := waitWithCtx(jitter(500*time.Millisecond, 250*time.Millisecond)); err != nil {
						errCh <- err
						return
					}
					continue // retry connection
				}
				errCh <- err
				return
			}

			for {
				// Create a new empty reply message of the correct type
				reply := newReply() // pointer to proto message

				// Receive a message from the stream
				recvErr := stream.RecvMsg(reply)
				if recvErr != nil {
					// Network/server error: attempt reconnect if recoverable
					if s, ok := status.FromError(recvErr); ok && s.Code() == codes.Unavailable {
						reconnectRequired = true
						break // retry stream
					}
					// Treat EOF as transient: server closed stream; reconnect
					if errors.Is(recvErr, io.EOF) {
						reconnectRequired = true
						break
					}
					// User cancelled or deadline exceeded
					if errors.Is(recvErr, context.Canceled) || errors.Is(recvErr, context.DeadlineExceeded) {
						errCh <- recvErr
						return
					}
					// All other errors: fail and close
					errCh <- recvErr
					return
				}

				// Check for logical/API errors inside the proto reply
				apiErr := getError(reply)
				if apiErr != nil {
					code := apiErr.GetErrorCode()
					// Certain terminal errors are recoverable (reconnect)
					if code == "TERMINAL_INSTANCE_NOT_FOUND" || code == "TERMINAL_REGISTRY_TERMINAL_NOT_FOUND" {
						reconnectRequired = true
						break
					}
					// All other API errors: report and end
					errCh <- fmt.Errorf("API error: %v", apiErr)
					return
				}

				// Extract the real data from reply (skip if not present)
				if data, ok := getData(reply); ok {
					select {
					case dataCh <- data:
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
				}
			}

			// Handle reconnect logic
			if reconnectRequired {
				if err := waitWithCtx(jitter(500*time.Millisecond, 250*time.Millisecond)); err != nil {
					errCh <- err
					return
				}
				continue // retry outer loop (reconnect)
			}
			break // Exit outer loop if not reconnecting
		}
	}()
	return dataCh, errCh
}


//=== 📂 Order Operations ===

// OrderSend places a new trade order (market or pending) on the connected MT4 terminal.
//
// Parameters:
//   - ctx: Context used for request cancellation or timeout control.
//   - symbol: Trading symbol (e.g., "EURUSD") for the order.
//   - operationType: Type of order operation (e.g., Buy, Sell, BuyLimit, etc.), defined by OrderSendOperationType enum.
//   - volume: Number of lots to trade (e.g., 0.1, 1.0).
//   - price: Price for the order (used for pending orders). Ignored for market orders if nil.
//   - slippage: Maximum allowed slippage in points for market orders (optional).
//   - stoploss: Price level for the Stop Loss (optional).
//   - takeprofit: Price level for the Take Profit (optional).
//   - comment: Optional text comment associated with the order (visible in the terminal).
//   - magicNumber: Identifier used by expert advisors to track orders (optional).
//   - expiration: Expiration time for pending orders (optional; ignored for market orders).
//
// Returns:
//   - Pointer to OrderSendData structure containing result details (e.g., ticket number, status).
//   - Error object if the operation fails due to network issues, validation, or terminal errors.
//
// The method wraps the gRPC call to the TradeClient.OrderSend method, preparing all optional fields if provided.
// It uses ExecuteWithReconnect to automatically retry on transient network or session errors,
// and checks for application-level errors returned by the terminal via the response.Error field.
func (a *MT4Account) OrderSend(
	ctx context.Context,
	symbol string,
	operationType pb.OrderSendOperationType,
	volume float64,
	price *float64,
	slippage *int32,
	stoploss *float64,
	takeprofit *float64,
	comment *string,
	magicNumber *int32,
	expiration *timestamppb.Timestamp,
) (*pb.OrderSendData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second) // trades usually tolerate a slightly longer timeout
		defer cancel()
	}

	// Ensure connection before making the call.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Build request with optional fields preserved as pointers.
	req := &pb.OrderSendRequest{
		Symbol:        symbol,
		OperationType: operationType,
		Volume:        volume,
	}
	if price != nil {
		req.Price = price
	}
	if slippage != nil {
		req.Slippage = slippage
	}
	if stoploss != nil {
		req.Stoploss = stoploss
	}
	if takeprofit != nil {
		req.Takeprofit = takeprofit
	}
	if comment != nil {
		req.Comment = comment
	}
	if magicNumber != nil {
		req.MagicNumber = magicNumber
	}
	if expiration != nil {
		req.Expiration = expiration
	}

	// gRPC call wrapper: attach metadata to the normalized context.
	grpcCall := func(headers metadata.MD) (*pb.OrderSendReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderSend(c, req)
	}

	// Extract API-level error if present.
	errorSelector := func(reply *pb.OrderSendReply) *pb.Error {
		return reply.GetError()
	}

	// Execute with reconnect/retry semantics.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}


// OrderClose closes or deletes an existing market or pending order on the MT4 terminal.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Unique identifier (ticket number) of the order to be closed.
//   - lots: Optional parameter to partially close the order (e.g., 0.5 lots). If nil, the full order is closed.
//   - price: Optional closing price. Used for market orders; if nil, market price is used.
//   - slippage: Optional maximum allowable slippage in points. Ignored for pending orders.
//
// Returns:
//   - Pointer to OrderCloseDeleteData structure with details of the operation (e.g., confirmation).
//   - Error if not connected, or if the gRPC/API call fails or is rejected by the terminal.
//
// This method sends a request to close or delete an order via the TradeClient.OrderCloseDelete RPC call.
// Optional fields (lots, price, slippage) are only included if explicitly provided.
// It uses ExecuteWithReconnect to automatically retry on connection or session-related errors,
// and checks for application-level errors in the gRPC response via the Error field.
func (a *MT4Account) OrderClose(
	ctx context.Context,
	ticket int32,
	lots, price *float64,
	slippage *int32,
) (*pb.OrderCloseDeleteData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second) // closing trades may take a bit longer
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	req := &pb.OrderCloseDeleteRequest{
		OrderTicket: ticket,
	}
	if lots != nil {
		req.Lots = lots
	}
	if price != nil {
		req.ClosingPrice = price
	}
	if slippage != nil {
		req.Slippage = slippage
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderCloseDeleteReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.TradeClient.OrderCloseDelete(c, req)
	}

	errorSelector := func(reply *pb.OrderCloseDeleteReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}


// OrderDelete removes a pending order from the MT4 terminal using its ticket number.
//
// Parameters:
//   - ctx: Context for cancellation or timeout.
//   - ticket: Ticket number of the pending order to delete.
//
// Returns:
//   - Pointer to OrderCloseDeleteData containing result information (e.g., confirmation).
//   - Error if the order could not be deleted or if the operation fails.
//
// This is a convenience wrapper around OrderClose with all optional parameters set to nil,
// meaning the full pending order is deleted using default behavior.
// Only applicable to pending orders; attempting to delete an open market order may result in an error.

func (a *MT4Account) OrderDelete(ctx context.Context, ticket int32) (*pb.OrderCloseDeleteData, error) {
	return a.OrderClose(ctx, ticket, nil, nil, nil)
}

// OrderSelect retrieves detailed information about a specific opened order by its ticket number.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Ticket number (unique ID) of the order to search for.
//
// Returns:
//   - Pointer to OpenedOrderInfo struct containing full order details.
//   - Error if the terminal is not connected, or if the order is not found.
//
// This method fetches the list of currently opened orders by calling OpenedOrders(),
// then iterates through them to find a match with the specified ticket.
// Returns an error if the order does not exist or if the API call fails.
//
// Note:
//   - This method does not query closed or pending orders — only currently opened ones.
//   - Matching is done locally by comparing tickets from the retrieved list.

func (a *MT4Account) OrderSelect(ctx context.Context, ticket int32) (*pb.OpenedOrderInfo, error) {
	if !a.isConnected() {
		return nil, errors.New("not connected to terminal")
	}

	orders, err := a.OpenedOrders(ctx)
	if err != nil {
		return nil, err
	}

	for _, order := range orders.GetOrderInfos() {
		if order.GetTicket() == ticket {
			return order, nil
		}
	}

	return nil, fmt.Errorf("order with ticket %d not found", ticket)
}


// OrderCloseBy closes a market order by pairing it with an opposite order (i.e., a hedge).
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//   - ticketToClose: Ticket number of the order to be closed.
//   - oppositeTicket: Ticket number of the opposite-direction order used to close the first one.
//
// Returns:
//   - Pointer to OrderCloseByData structure with details of the closing operation.
//   - Error if not connected or if the API call fails or is rejected by the terminal.
//
// This method is used when closing one position using another of the opposite direction,
// typically in hedge-mode accounts. The operation is atomic — both orders are closed simultaneously.
//
// It sends a gRPC request to the terminal using TradeClient.OrderCloseBy,
// and applies automatic retry logic through ExecuteWithReconnect.
// Any application-level errors returned by the terminal are extracted from the reply.
//
// Note:
//   - Both orders must be opened and of equal volume.
//   - This operation is only available in hedge-enabled accounts.
func (a *MT4Account) OrderCloseBy(
	ctx context.Context,
	ticketToClose int32,
	oppositeTicket int32,
) (*pb.OrderCloseByData, error) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second) // close-by may take a bit longer
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	req := &pb.OrderCloseByRequest{
		TicketToClose:           ticketToClose,
		OppositeTicketClosingBy: oppositeTicket,
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderCloseByReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers)
		return a.TradeClient.OrderCloseBy(c, req)
	}

	errorSelector := func(reply *pb.OrderCloseByReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}


// OrdersTotal returns the total number of currently opened orders on the connected MT4 trading account.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//
// Returns:
//   - Integer value representing the number of open orders (int32).
//   - Error if the terminal is not connected or the API call fails.
//
// This method calls OpenedOrders() to retrieve the list of currently active market and pending orders,
// and returns the total count. It does not include closed or historical orders.
//
// Note:
//   - This is a runtime query and reflects the current state of the trading account.
//   - If the connection to the terminal is not established, an error is returned immediately.

// OrdersTotal returns the total number of currently opened orders.
func (a *MT4Account) OrdersTotal(ctx context.Context) (int32, error) {
	if !a.isConnected() {
		return 0, errors.New("not connected")
	}

	orders, err := a.OpenedOrders(ctx)
	if err != nil {
		return 0, err
	}
	if orders == nil {
		return 0, nil // no opened orders
	}

	return int32(len(orders.GetOrderInfos())), nil
}


// OrderModify updates the parameters of an existing order on the MT4 terminal,
// such as entry price, stop loss, take profit, or expiration time.
//
// Parameters:
//   - ctx: Context for request timeout or cancellation.
//   - ticket: Ticket number (unique ID) of the order to be modified.
//   - price: New entry price (optional). Used primarily for pending orders.
//   - stoploss: New Stop Loss level (optional).
//   - takeprofit: New Take Profit level (optional).
//   - expiration: New expiration time for pending orders (optional).
//
// Returns:
//   - Boolean indicating whether the order was successfully modified.
//   - Error if the modification failed due to terminal state, invalid parameters, or connectivity issues.
//
// This method constructs an OrderModifyRequest message with any provided fields,
// and sends it to the terminal via the TradeClient.OrderModify RPC call.
// Fields that are nil are omitted, so only specified parameters are updated.
//
// The call is wrapped in ExecuteWithReconnect to automatically retry on transient connection issues,
// and inspects the gRPC response for any terminal-level errors.
//
// Note:
//   - Only pending orders can have price or expiration modified.
//   - The order must still be valid (not closed or already filled).
func (a *MT4Account) OrderModify(
	ctx context.Context,
	ticket int32,
	price, stoploss, takeprofit *float64,
	expiration *timestamppb.Timestamp,
) (bool, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second) // order modify can take slightly longer
		defer cancel()
	}

	if !a.isConnected() {
		return false, errors.New("not connected")
	}

	req := &pb.OrderModifyRequest{
		OrderTicket: ticket,
	}
	if price != nil {
		req.NewPrice = price
	}
	if stoploss != nil {
		req.NewStopLoss = stoploss
	}
	if takeprofit != nil {
		req.NewTakeProfit = takeprofit
	}
	if expiration != nil {
		req.NewExpiration = expiration
		// NOTE: if your proto supports expiration_time_type, set it here accordingly.
	}

	grpcCall := func(headers metadata.MD) (*pb.OrderModifyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.TradeClient.OrderModify(c, req)
	}

	errorSelector := func(reply *pb.OrderModifyReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return false, err
	}
	if reply.GetData() == nil {
		return false, fmt.Errorf("empty reply data")
	}
	return reply.GetData().GetOrderWasModified(), nil
}


// OpenedOrders retrieves a list of currently opened orders for the connected MT4 account.
//
// This method performs a gRPC call to the AccountHelper.OpenedOrders service.
// It automatically handles retries if the connection is lost or if the terminal session is dropped.
//
// Parameters:
//   - ctx: Context for timeout or cancellation (e.g. context.WithTimeout).
//
// Returns:
//   - A pointer to OpenedOrdersData struct containing all open orders.
//   - An error if the request fails or account is not connected.
func (a *MT4Account) OpenedOrders(ctx context.Context) (*pb.OpenedOrdersData, error) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call; 2–3s достаточно
		defer cancel()
	}

	// Check if the account is connected (either by host/port or server name).
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Prepare the request message — this doesn't require any fields.
	req := &pb.OpenedOrdersRequest{}

	// Define the actual gRPC call as a closure that uses metadata headers.
	grpcCall := func(headers metadata.MD) (*pb.OpenedOrdersReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.OpenedOrders(c, req)
	}

	// Define how to extract the API-level error from the response (if any).
	errorSelector := func(reply *pb.OpenedOrdersReply) *pb.Error {
		return reply.GetError()
	}

	// Execute the request using the reconnect helper.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Extract and return the useful data from the reply.
	return reply.GetData(), nil
}


// OpenedOrdersTickets retrieves the ticket IDs of all currently opened orders.
//
// This method sends a simple gRPC request to the AccountHelper.OpenedOrdersTickets method.
// It returns the list of tickets (with symbols and order type) for all active orders.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//
// Returns:
//   - A pointer to OpenedOrdersTicketsData containing the order tickets and details.
//   - An error if the request fails or the session is not connected.
func (a *MT4Account) OpenedOrdersTickets(ctx context.Context) (*pb.OpenedOrdersTicketsData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call; short timeout
		defer cancel()
	}

	// Ensure the account is connected before making any requests.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Create the request object. This message doesn't require any input fields.
	req := &pb.OpenedOrdersTicketsRequest{}

	// Define the gRPC call logic, including session metadata headers.
	grpcCall := func(headers metadata.MD) (*pb.OpenedOrdersTicketsReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.OpenedOrdersTickets(c, req)
	}

	// Define how to extract an API error from the server response.
	errorSelector := func(reply *pb.OpenedOrdersTicketsReply) *pb.Error {
		return reply.GetError()
	}

	// Perform the call using retry/reconnect logic in case of temporary failures.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return the data field, which contains the actual ticket information.
	return reply.GetData(), nil
}


// OrdersHistory retrieves a paginated list of historical orders for the connected MT4 account.
//
// This method allows filtering by time range, sorting method, and pagination.
//
// Parameters:
//   - ctx: Context for cancellation or timeout.
//   - sortType: Sorting option (e.g. newest to oldest, oldest to newest).
//   - from: Optional start time (nil = no lower bound).
//   - to: Optional end time (nil = no upper bound).
//   - page: Optional page number (nil = no pagination).
//   - itemsPerPage: Optional number of items per page (nil = default set by server).
//
// Returns:
//   - Pointer to OrdersHistoryData containing order history records.
//   - An error if the gRPC request fails or connection is invalid.
func (a *MT4Account) OrdersHistory(
	ctx context.Context,
	sortType pb.EnumOrderHistorySortType,
	from, to *time.Time,
	page, itemsPerPage *int32,
) (*pb.OrdersHistoryData, error) {
	// Normalize context: ensure non-nil and add a per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	// History calls may scan larger ranges; give them a bit more time.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}

	// Check that the account is connected before calling the API.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Build the request message with provided filters and options.
	req := &pb.OrdersHistoryRequest{
		InputSortMode: sortType,
	}

	// Convert optional time filters to protobuf timestamps.
	if from != nil {
		req.InputFrom = timestamppb.New(*from)
	}
	if to != nil {
		req.InputTo = timestamppb.New(*to)
	}

	// Set pagination if provided.
	if page != nil {
		req.PageNumber = page
	}
	if itemsPerPage != nil {
		req.ItemsPerPage = itemsPerPage
	}

	// Define the gRPC call with headers and context.
	grpcCall := func(headers metadata.MD) (*pb.OrdersHistoryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.OrdersHistory(c, req)
	}

	// Extract error from the response if it exists.
	errorSelector := func(reply *pb.OrdersHistoryReply) *pb.Error {
		return reply.GetError()
	}

	// Call the method using retry wrapper for reconnects and transient errors.
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	// Return the resulting historical order data.
	return reply.GetData(), nil
}


// === 📂 Market Info / Symbol Info ===

// Return the actual quote result from the server.
// Quote retrieves the latest market quote (bid/ask) for a given symbol.
//
// This method uses MarketInfo.Quote and returns bid, ask, and timestamp data.
//
// Parameters:
//   - ctx: Context for cancellation or timeout.
//   - symbol: Symbol name (e.g., "EURUSD").
//
// Returns:
//   - Pointer to QuoteData (bid, ask, high, low, etc.)
//   - Error if connection or API fails.
func (a *MT4Account) Quote(ctx context.Context, symbol string) (*pb.QuoteData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // quote is a fast read
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	req := &pb.QuoteRequest{
		Symbol: symbol,
	}

	grpcCall := func(headers metadata.MD) (*pb.QuoteReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.MarketInfoClient.Quote(c, req)
	}

	errorSelector := func(reply *pb.QuoteReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}
	return reply.GetData(), nil
}


// QuoteMany retrieves latest market quotes for multiple trading symbols.
//
// It makes a gRPC call to MarketInfo.QuoteMany and returns the data as a list of QuoteData entries.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//   - symbols: Slice of symbol names (e.g., []string{"EURUSD", "GBPUSD"}).
//
// Returns:
//   - A pointer to QuoteManyData containing all quotes.
//   - An error if the request fails or if the account is not connected.
func (a *MT4Account) QuoteMany(ctx context.Context, symbols []string) (*pb.QuoteManyData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // fast read
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if len(symbols) == 0 {
		return nil, errors.New("no symbols provided")
	}

	req := &pb.QuoteManyRequest{
		Symbols: symbols,
	}

	grpcCall := func(headers metadata.MD) (*pb.QuoteManyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.MarketInfoClient.QuoteMany(c, req)
	}

	errorSelector := func(reply *pb.QuoteManyReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}


// ShowAllSymbols retrieves all available trading symbols from the server.
//
// This method sends a request to the MetaTrader terminal (via gRPC) asking for
// a list of all known trading instruments (symbols), such as "EURUSD", "GBPJPY", etc.
// These symbols include both visible (in Market Watch) and hidden ones.
func (a *MT4Account) ShowAllSymbols(ctx context.Context) (*pb.SymbolsData, error) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call
		defer cancel()
	}

	// Ensure the account is connected to a server.
	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	req := &pb.SymbolsRequest{}

	grpcCall := func(headers metadata.MD) (*pb.SymbolsReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.MarketInfoClient.Symbols(c, req)
	}

	errorSelector := func(reply *pb.SymbolsReply) *pb.Error {
		return reply.GetError()
	}

	// Execute the request with reconnection logic
	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}


// SymbolParams retrieves detailed trading parameters for a single symbol.
//
// Parameters:
//   - ctx: Context for cancellation or timeout.
//   - symbol: Name of the trading symbol (e.g., "EURUSD").
//
// Returns:
//   - *pb.SymbolParamsManyInfo: Detailed symbol info (digits, volume limits, trade mode, etc.).
//   - error: If request fails or symbol not found.
//
// Notes:
//   - Internally calls SymbolParamsMany with one symbol.
//   - Returns the first result in the symbol info list.
//   - Performs automatic reconnect if the terminal connection is lost.
func (a *MT4Account) SymbolParams(ctx context.Context, symbol string) (*pb.SymbolParamsManyInfo, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call
		defer cancel()
	}

	if !a.isConnected() {
		return nil, fmt.Errorf("not connected")
	}

	req := &pb.SymbolParamsManyRequest{
		SymbolName: proto.String(symbol), // optional string; keep pointer form
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolParamsManyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.SymbolParamsMany(c, req)
	}

	errorSelector := func(reply *pb.SymbolParamsManyReply) *pb.Error {
		return reply.GetError()
	}

	resp, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	if resp.GetData() == nil {
		return nil, fmt.Errorf("empty reply data")
	}
	symbols := resp.GetData().GetSymbolInfos()
	if len(symbols) == 0 {
		return nil, fmt.Errorf("no parameters returned for symbol: %s", symbol)
	}

	return symbols[0], nil
}


// QuoteHistory retrieves historical candlestick data (OHLC) for a symbol within a time range.
//
// This method uses the MarketInfo.QuoteHistory endpoint. Timeframe options include M1, H1, D1, etc.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//   - symbol: Trading symbol (e.g., "EURUSD").
//   - timeframe: Chart timeframe (e.g., QH_PERIOD_M1, QH_PERIOD_H1).
//   - from: Start time (UTC).
//   - to: End time (UTC).
//
// Returns:
//   - A pointer to QuoteHistoryData containing historical OHLC candles.
//   - An error if the request fails or session is not connected.
func (a *MT4Account) QuoteHistory(
	ctx context.Context,
	symbol string,
	timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME,
	from, to time.Time,
) (*pb.QuoteHistoryData, error) {
	// Normalize context: ensure non-nil and add a per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	// History calls may scan larger ranges; give them a bit more time.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if to.Before(from) {
		return nil, fmt.Errorf("invalid time range: to < from")
	}

	req := &pb.QuoteHistoryRequest{
		Symbol:    symbol,
		Timeframe: timeframe,
		FromTime:  timestamppb.New(from),
		ToTime:    timestamppb.New(to),
	}

	grpcCall := func(headers metadata.MD) (*pb.QuoteHistoryReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.MarketInfoClient.QuoteHistory(c, req)
	}

	errorSelector := func(reply *pb.QuoteHistoryReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}


// Symbols retrieves all available trading symbols from the server.
//
// Returns a list of SymbolNameInfo entries (symbol name and index).
func (a *MT4Account) Symbols(ctx context.Context) (*pb.SymbolsData, error) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	req := &pb.SymbolsRequest{}

	grpcCall := func(headers metadata.MD) (*pb.SymbolsReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.MarketInfoClient.Symbols(c, req)
	}

	errorSelector := func(reply *pb.SymbolsReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}


// SymbolParamsMany retrieves trading symbol parameters for one or more symbols.
//
// Parameters:
//   - ctx: Context for timeout/cancel.
//   - symbols: List of symbol names.
//
// Returns:
//   - SymbolParamsManyData containing all symbol param info.
func (a *MT4Account) SymbolParamsMany(ctx context.Context, symbols []string) (*pb.SymbolParamsManyData, error) {
	
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}

	// Build request: empty => all symbols; otherwise filter by first symbol.
	req := &pb.SymbolParamsManyRequest{}
	if len(symbols) > 0 {
		req.SymbolName = proto.String(symbols[0])
	}

	grpcCall := func(headers metadata.MD) (*pb.SymbolParamsManyReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.SymbolParamsMany(c, req)
	}

	errorSelector := func(reply *pb.SymbolParamsManyReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	if reply.GetData() == nil {
		return nil, fmt.Errorf("empty reply data")
	}
	return reply.GetData(), nil
}


// TickValueWithSize calculates tick value, lot size, and contract size for specified symbols.
//
// Parameters:
//   - ctx: Context for cancel/timeout.
//   - symbols: List of symbol+lot pairs.
//
// Returns:
//   - TickValueWithSizeData containing value calculations.
func (a *MT4Account) TickValueWithSize(ctx context.Context, symbolNames []string) (*pb.TickValueWithSizeData, error) {
	// Normalize context: ensure non-nil and add a short per-call timeout if none is set.
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second) // read-only call
		defer cancel()
	}

	if !a.isConnected() {
		return nil, errors.New("not connected")
	}
	if len(symbolNames) == 0 {
		return nil, errors.New("no symbols provided")
	}

	req := &pb.TickValueWithSizeRequest{
		SymbolNames: symbolNames,
	}

	grpcCall := func(headers metadata.MD) (*pb.TickValueWithSizeReply, error) {
		c := metadata.NewOutgoingContext(ctx, headers) // attach metadata to normalized ctx
		return a.AccountClient.TickValueWithSize(c, req)
	}

	errorSelector := func(reply *pb.TickValueWithSizeReply) *pb.Error {
		return reply.GetError()
	}

	reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
	if err != nil {
		return nil, err
	}

	if reply.GetData() == nil {
		return nil, fmt.Errorf("empty reply data")
	}
	return reply.GetData(), nil
}


// === 📂 Streaming ===

// OnTrade subscribes to real-time trade updates (open, close, modify).
//
// Returns:
//   - A receive-only channel of TradeData messages.
//   - A receive-only error channel (closes when the stream ends).
func (a *MT4Account) OnTrade(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error) {
	// Normalize context for stream lifecycle.
	if ctx == nil {
		ctx = context.Background()
	}

	if a.Id == uuid.Nil {
		dataCh := make(chan *pb.OnTradeData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- errors.New("not connected")
		}()
		return dataCh, errCh
	}

	req := &pb.OnTradeRequest{}

	getError := func(reply *pb.OnTradeReply) *pb.Error {
		return reply.GetError()
	}
	getData := func(reply *pb.OnTradeReply) (*pb.OnTradeData, bool) {
		data := reply.GetData()
		return data, data != nil
	}
	newReply := func() *pb.OnTradeReply { return new(pb.OnTradeReply) }

	return ExecuteStreamWithReconnect(
		ctx, a, req,
		func(r *pb.OnTradeRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			return a.SubscriptionClient.OnTrade(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError, getData, newReply,
	)
}


// OnOpenedOrdersProfit subscribes to periodic updates of profits for all open orders.
//
// Parameters:
//   - ctx: Context for cancel/timeout.
//   - intervalMs: Update interval in milliseconds (e.g., 1000 = 1s).
//
// Returns:
//   - A receive-only channel of OnOpenedOrdersProfitData updates.
//   - A receive-only error channel.
func (a *MT4Account) OnOpenedOrdersProfit(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersProfitData, <-chan error) {
	// Ensure ctx is non-nil (stream invoker and selects rely on it)
	if ctx == nil {
		ctx = context.Background()
	}

	if a.Id == uuid.Nil {
		dataCh := make(chan *pb.OnOpenedOrdersProfitData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- errors.New("not connected")
		}()
		return dataCh, errCh
	}

	req := &pb.OnOpenedOrdersProfitRequest{
		TimerPeriodMilliseconds: intervalMs,
	}

	getError := func(reply *pb.OnOpenedOrdersProfitReply) *pb.Error {
		return reply.GetError()
	}
	getData := func(reply *pb.OnOpenedOrdersProfitReply) (*pb.OnOpenedOrdersProfitData, bool) {
		data := reply.GetData()
		return data, data != nil
	}
	newReply := func() *pb.OnOpenedOrdersProfitReply { return new(pb.OnOpenedOrdersProfitReply) }

	return ExecuteStreamWithReconnect(
		ctx, a, req,
		func(r *pb.OnOpenedOrdersProfitRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			return a.SubscriptionClient.OnOpenedOrdersProfit(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError, getData, newReply,
	)
}


// OnOpenedOrdersTickets subscribes to periodic updates of opened order ticket IDs.
//
// Parameters:
//   - ctx: Context for cancel/timeout.
//   - intervalMs: Update interval in milliseconds.
//
// Returns:
//   - A receive-only channel of OnOpenedOrdersTicketsData updates.
//   - A receive-only error channel.
func (a *MT4Account) OnOpenedOrdersTickets(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersTicketsData, <-chan error) {
	
	if ctx == nil {
		ctx = context.Background()
	}

	if a.Id == uuid.Nil {
		dataCh := make(chan *pb.OnOpenedOrdersTicketsData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- errors.New("not connected")
		}()
		return dataCh, errCh
	}

	req := &pb.OnOpenedOrdersTicketsRequest{
		PullIntervalMilliseconds: intervalMs,
	}

	getError := func(reply *pb.OnOpenedOrdersTicketsReply) *pb.Error {
		return reply.GetError()
	}
	getData := func(reply *pb.OnOpenedOrdersTicketsReply) (*pb.OnOpenedOrdersTicketsData, bool) {
		data := reply.GetData()
		return data, data != nil
	}
	newReply := func() *pb.OnOpenedOrdersTicketsReply { return new(pb.OnOpenedOrdersTicketsReply) }

	return ExecuteStreamWithReconnect(
		ctx, a, req,
		func(r *pb.OnOpenedOrdersTicketsRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			return a.SubscriptionClient.OnOpenedOrdersTickets(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError, getData, newReply,
	)
}


// OnSymbolTick subscribes to real-time tick data for specified symbols, with reconnection logic.
//
// Parameters:
//   - ctx: Context for cancellation/timeouts
//   - symbols: Slice of symbol names (e.g., []string{"EURUSD", "USDJPY"})
//
// Returns:
//   - Receive-only channel of *pb.OnSymbolTickData (each tick)
//   - Receive-only error channel
func (a *MT4Account) OnSymbolTick(
	ctx context.Context,
	symbols []string,
) (<-chan *pb.OnSymbolTickData, <-chan error) {
	// Ensure ctx is non-nil for stream lifecycle
	if ctx == nil {
		ctx = context.Background()
	}

	// Check that the account is connected (has a valid UUID)
	if a.Id == uuid.Nil {
		dataCh := make(chan *pb.OnSymbolTickData)
		errCh := make(chan error, 1)
		go func() {
			defer close(dataCh)
			defer close(errCh)
			errCh <- errors.New("not connected")
		}()
		return dataCh, errCh
	}

	// Build the request message for the stream
	req := &pb.OnSymbolTickRequest{SymbolNames: symbols}

	// Function to extract API error from the proto reply
	getError := func(reply *pb.OnSymbolTickReply) *pb.Error {
		return reply.GetError()
	}
	// Function to extract the tick data (returns (data, ok))
	getData := func(reply *pb.OnSymbolTickReply) (*pb.OnSymbolTickData, bool) {
		data := reply.GetData()
		return data, data != nil
	}

	// The "newReply" function returns a new pointer to your proto reply type.
	newReply := func() *pb.OnSymbolTickReply { return new(pb.OnSymbolTickReply) }

	// Call the generic streaming helper.
	dataCh, errCh := ExecuteStreamWithReconnect(
		ctx, a, req,
		func(r *pb.OnSymbolTickRequest, md metadata.MD, ctx context.Context) (grpc.ClientStream, error) {
			return a.SubscriptionClient.OnSymbolTick(metadata.NewOutgoingContext(ctx, md), r)
		},
		getError,
		getData,
		newReply,
	)
	return dataCh, errCh
}







