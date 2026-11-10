package mt4

// MT4Service is the mid-level facade over MT4Account (Layer 2).
//
// It mirrors the three-layer SDK design: friendly arguments and clean return
// values (unwrapped slices/scalars, simplified parameters, clean stream names)
// over the low-level MT4Account. It never prints — demonstrations live in the
// examples. Order construction helpers, pip math and risk sizing live one layer
// up in MT4Sugar.

import (
	"context"
	"time"

	pb "github.com/MetaRPC/GoMT4/package"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// MT4Service is the mid-level facade. Where MT4Account returns the raw protobuf
// wrappers, MT4Service hands back the useful part directly — a slice, a scalar,
// a bool — and takes simpler arguments.
//
// Example — the same read, one layer apart:
//
//	raw, _ := acc.OpenedOrders(ctx)       // *pb.OpenedOrdersData (wrapper)
//	orders := raw.GetOrderInfos()         // ...then unwrap yourself
//
//	orders, _ := svc.OpenedOrders(ctx)    // []*pb.OpenedOrderInfo (already unwrapped)
type MT4Service struct {
	account *MT4Account
}

// NewMT4Service wraps a low-level account:
//
//	svc := mt4.NewMT4Service(acc)
func NewMT4Service(acc *MT4Account) *MT4Service {
	return &MT4Service{account: acc}
}

// Account returns the underlying MT4Account for callers that need the raw layer.
func (s *MT4Service) Account() *MT4Account { return s.account }

// ===== 1. Account information =====

// AccountSummary returns the full account snapshot.
func (s *MT4Service) AccountSummary(ctx context.Context) (*pb.AccountSummaryData, error) {
	return s.account.AccountSummary(ctx)
}

// TradeStats returns aggregated trade statistics (optionally from a server-time
// "yyyy-MM-ddTHH:mm:ss" string).
func (s *MT4Service) TradeStats(ctx context.Context, from *string) (*pb.StatsData, error) {
	return s.account.TradeStats(ctx, from, nil)
}

// TradeStatsEquityHistory returns stats plus equity points and trades.
func (s *MT4Service) TradeStatsEquityHistory(ctx context.Context, from *string, excludeSameBars *bool) (*pb.StatsWithChartsData, error) {
	return s.account.TradeStatsEquityHistory(ctx, from, nil, excludeSameBars)
}

// EquityHistory returns the equity-curve points.
func (s *MT4Service) EquityHistory(ctx context.Context, from *string, excludeSameBars *bool) (*pb.EquityHistoryData, error) {
	return s.account.EquityHistory(ctx, from, nil, excludeSameBars)
}

// TradeUnrealizedPL returns the unrealized-P/L extrema per open order.
func (s *MT4Service) TradeUnrealizedPL(ctx context.Context, from *string) (*pb.TradeUnrealizedPLData, error) {
	return s.account.TradeUnrealizedPL(ctx, from, nil)
}

// ===== 2. Orders & history =====

// OpenedOrders returns the open market + pending orders, unwrapped to a slice
// (no wrapper struct to unpack).
//
// Example:
//
//	for _, o := range must(svc.OpenedOrders(ctx)) {
//	    fmt.Printf("#%d %s %.2f P/L %+.2f\n", o.GetTicket(), o.GetSymbol(), o.GetLots(), o.GetProfit())
//	}
func (s *MT4Service) OpenedOrders(ctx context.Context) ([]*pb.OpenedOrderInfo, error) {
	data, err := s.account.OpenedOrders(ctx)
	if err != nil {
		return nil, err
	}
	return data.GetOrderInfos(), nil
}

// OpenedOrdersTickets returns just the ticket numbers of open/pending orders.
func (s *MT4Service) OpenedOrdersTickets(ctx context.Context) ([]int32, error) {
	data, err := s.account.OpenedOrdersTickets(ctx)
	if err != nil {
		return nil, err
	}
	return data.GetTickets(), nil
}

// OrdersHistory returns closed/cancelled order history over a window (paging
// defaults applied).
func (s *MT4Service) OrdersHistory(ctx context.Context, sortType pb.EnumOrderHistorySortType, from, to *time.Time) (*pb.OrdersHistoryData, error) {
	return s.account.OrdersHistory(ctx, sortType, from, to, nil, nil)
}

// ===== 3. Symbols & quotes =====

// Symbols returns every symbol name the terminal knows, unwrapped to a slice.
func (s *MT4Service) Symbols(ctx context.Context) ([]*pb.SymbolNameInfo, error) {
	data, err := s.account.Symbols(ctx)
	if err != nil {
		return nil, err
	}
	return data.GetSymbolNameInfos(), nil
}

// SymbolSelect adds/removes a symbol in Market Watch; returns whether it succeeded.
func (s *MT4Service) SymbolSelect(ctx context.Context, symbol string, selected bool) (bool, error) {
	return s.account.SymbolSelect(ctx, symbol, selected)
}

// SymbolParamsMany returns full parameter sets for the given symbols.
func (s *MT4Service) SymbolParamsMany(ctx context.Context, symbols []string) (*pb.SymbolParamsManyData, error) {
	return s.account.SymbolParamsMany(ctx, symbols)
}

// TickValueWithSize returns tick value / size / contract size for the given symbols.
func (s *MT4Service) TickValueWithSize(ctx context.Context, symbols []string) (*pb.TickValueWithSizeData, error) {
	return s.account.TickValueWithSize(ctx, symbols)
}

// Quote returns the latest tick for one symbol.
func (s *MT4Service) Quote(ctx context.Context, symbol string) (*pb.QuoteData, error) {
	return s.account.Quote(ctx, symbol)
}

// QuoteMany returns the latest ticks for several symbols.
func (s *MT4Service) QuoteMany(ctx context.Context, symbols []string) (*pb.QuoteManyData, error) {
	return s.account.QuoteMany(ctx, symbols)
}

// QuoteHistory returns OHLC bars for a symbol/timeframe over a time range.
func (s *MT4Service) QuoteHistory(ctx context.Context, symbol string, timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME, from, to time.Time) (*pb.QuoteHistoryData, error) {
	return s.account.QuoteHistory(ctx, symbol, timeframe, from, to)
}

// ServerTimezone returns the broker server's UTC offset.
func (s *MT4Service) ServerTimezone(ctx context.Context) (*pb.ServerTimezoneData, error) {
	return s.account.ServerTimezone(ctx)
}

// ===== 4. Trading operations =====

// OrderSend sends a market or pending order. The optional parameters are
// pointers so you can omit them (nil) — pass price only for pending orders.
// For side-based one-liners with defaults, use MT4Sugar's BuyMarket/BuyLimit/….
//
// Example — market BUY 0.10 lot with a stop and target:
//
//	r, err := svc.OrderSend(ctx, "EURUSD",
//	    pb.OrderSendOperationType_OC_OP_BUY, 0.10,
//	    nil,            // price (market)
//	    ptr(int32(20)), // slippage
//	    ptr(1.0800),    // stoploss
//	    ptr(1.0900),    // takeprofit
//	    nil, nil, nil)  // comment, magic, expiration
func (s *MT4Service) OrderSend(
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
	return s.account.OrderSend(ctx, symbol, operationType, volume, price, slippage, stoploss, takeprofit, comment, magicNumber, expiration)
}

// OrderModify changes an order's price/SL/TP/expiry. Returns whether it was modified.
func (s *MT4Service) OrderModify(ctx context.Context, ticket int32, price, stoploss, takeprofit *float64, expiration *timestamppb.Timestamp) (bool, error) {
	return s.account.OrderModify(ctx, ticket, price, stoploss, takeprofit, expiration)
}

// OrderClose closes a market order (full/partial) or deletes a pending.
func (s *MT4Service) OrderClose(ctx context.Context, ticket int32, lots, price *float64, slippage *int32) (*pb.OrderCloseDeleteData, error) {
	return s.account.OrderClose(ctx, ticket, lots, price, slippage)
}

// OrderCloseBy closes an order by an opposite one.
func (s *MT4Service) OrderCloseBy(ctx context.Context, ticket, oppositeTicket int32) (*pb.OrderCloseByData, error) {
	return s.account.OrderCloseBy(ctx, ticket, oppositeTicket)
}

// ===== 5. Streaming =====

// StreamTicks streams live ticks for the given symbols as a (data, err) channel
// pair. Cancel the context to stop.
//
// Example:
//
//	data, errs := svc.StreamTicks(ctx, []string{"EURUSD"})
//	for {
//	    select {
//	    case d := <-data:
//	        t := d.GetSymbolTick(); fmt.Println(t.GetSymbol(), t.GetBid(), t.GetAsk())
//	    case err := <-errs:
//	        if err != nil { return err }
//	    case <-ctx.Done():
//	        return nil
//	    }
//	}
func (s *MT4Service) StreamTicks(ctx context.Context, symbols []string) (<-chan *pb.OnSymbolTickData, <-chan error) {
	return s.account.OnSymbolTick(ctx, symbols)
}

// StreamTrades streams live trade events.
func (s *MT4Service) StreamTrades(ctx context.Context) (<-chan *pb.OnTradeData, <-chan error) {
	return s.account.OnTrade(ctx)
}

// StreamOpenedTickets streams the open-order ticket list every intervalMs.
func (s *MT4Service) StreamOpenedTickets(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersTicketsData, <-chan error) {
	return s.account.OnOpenedOrdersTickets(ctx, intervalMs)
}

// StreamOpenedProfit streams per-order floating profit every intervalMs.
func (s *MT4Service) StreamOpenedProfit(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersProfitData, <-chan error) {
	return s.account.OnOpenedOrdersProfit(ctx, intervalMs)
}

// StreamConnectState streams terminal connection-state events.
func (s *MT4Service) StreamConnectState(ctx context.Context, terminalIds []string) (<-chan *pb.ConnectStateData, <-chan error) {
	return s.account.OnConnectState(ctx, terminalIds)
}
