package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"GoMT4/internal/account"
	pb "GoMT4/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MT4Service provides higher-level methods wrapping MT4Account operations.
type MT4Service struct {
	account *account.MT4Account
}

// NewMT4Service constructs a new MT4Service.
func NewMT4Service(acc *account.MT4Account) *MT4Service {
	return &MT4Service{account: acc}
}

// ShowAccountSummary prints basic account information.
func (s *MT4Service) ShowAccountSummary(ctx context.Context) {
	fmt.Println("=== Account Summary ===")
	summary, err := s.account.AccountSummary(ctx)
	if err != nil {
		log.Printf("❌ AccountSummary error: %v", err)
		return
	}
	fmt.Printf("Balance: %.2f, Equity: %.2f, Currency: %s\n",
		summary.GetAccountBalance(),
		summary.GetAccountEquity(),
		summary.GetAccountCurrency(),
	)
}

// ShowOpenedOrders lists currently opened orders.
func (s *MT4Service) ShowOpenedOrders(ctx context.Context) {
	fmt.Println("=== Opened Orders ===")
	ordersData, err := s.account.OpenedOrders(ctx)
	if err != nil {
		log.Printf("❌ OpenedOrders error: %v", err)
		return
	}
	infos := ordersData.GetOrderInfos()
	if len(infos) == 0 {
		fmt.Println("No opened orders.")
		return
	}
	for _, o := range infos {
		fmt.Printf("[%s] Ticket: %d | Symbol: %s | Lots: %.2f | OpenPrice: %.5f | Profit: %.2f\n",
			o.GetOrderType(),
			o.GetTicket(),
			o.GetSymbol(),
			o.GetLots(),
			o.GetOpenPrice(),
			o.GetProfit(),
		)
	}
}

// ShowOpenedOrderTickets prints the ticket IDs of opened orders.
func (s *MT4Service) ShowOpenedOrderTickets(ctx context.Context) {
	fmt.Println("=== Opened Order Tickets ===")
	ticketsData, err := s.account.OpenedOrdersTickets(ctx)
	if err != nil {
		log.Printf("❌ OpenedOrdersTickets error: %v", err)
		return
	}
	tickets := ticketsData.GetTickets()
	if len(tickets) == 0 {
		fmt.Println("📭 No open order tickets found.")
		return
	}
	fmt.Println("Open Order Tickets:")
	for _, t := range tickets {
		fmt.Printf(" - Ticket: %d\n", t)
	}
}

// ShowOrdersHistory prints the order history for the last 7 days.
func (s *MT4Service) ShowOrdersHistory(ctx context.Context) {
	fmt.Println("=== Orders History (last 7 days) ===")
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now()
	history, err := s.account.OrdersHistory(
		ctx,
		pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC,
		&from, &to, nil, nil,
	)
	if err != nil {
		log.Printf("❌ OrdersHistory error: %v", err)
		return
	}
	orders := history.GetOrdersInfo()
	if len(orders) == 0 {
		fmt.Println("📭 No historical orders found.")
		return
	}
	for _, o := range orders {
		fmt.Printf("[%s] Ticket: %d | Symbol: %s | Lots: %.2f | Open: %.5f | Close: %.5f | Profit: %.2f | Closed: %s\n",
			o.GetOrderType(),
			o.GetTicket(),
			o.GetSymbol(),
			o.GetLots(),
			o.GetOpenPrice(),
			o.GetClosePrice(),
			o.GetProfit(),
			o.GetCloseTime().AsTime().Format("2006-01-02 15:04:05"),
		)
	}
}

// ShowOrderSendExample sends a simple market Buy order.
func (s *MT4Service) ShowOrderSendExample(ctx context.Context, symbol string) {
	fmt.Println("📤 Sending market order...")
	data, err := s.account.OrderSend(
		ctx,
		symbol,
		pb.OrderSendOperationType_OC_OP_BUY,
		0.1,
		nil,
		ptrInt32(5),
		nil,
		nil,
		ptrString("Go order test"),
		ptrInt32(123456),
		nil,
	)
	if err != nil {
		log.Printf("❌ OrderSend error: %v", err)
		return
	}
	fmt.Printf("✅ Order opened! Ticket: %d, Price: %.5f, Time: %s\n",
		data.GetTicket(), data.GetPrice(), data.GetOpenTime().AsTime().Format("2006-01-02 15:04:05"),
	)
}

// ShowOrderModifyExample demonstrates modifying SL/TP.
func (s *MT4Service) ShowOrderModifyExample(ctx context.Context, ticket int32) {
	fmt.Println("=== Modify Order ===")
	newSL := 1.0500
	newTP := 1.0900
	modified, err := s.account.OrderModify(ctx, ticket, nil, &newSL, &newTP, nil)
	if err != nil {
		log.Printf("❌ OrderModify error: %v", err)
		return
	}
	if modified {
		fmt.Println("✅ Order successfully modified.")
	} else {
		fmt.Println("⚠️ Order was NOT modified.")
	}
}

// ShowOrderCloseExample shows closing a market order.
func (s *MT4Service) ShowOrderCloseExample(ctx context.Context, ticket int32) {
	fmt.Println("=== Close Order ===")
	result, err := s.account.OrderClose(ctx, ticket, nil, nil, nil)
	if err != nil {
		log.Printf("❌ OrderClose error: %v", err)
		return
	}
	fmt.Printf("✅ Order closed. Mode: %s", result.GetMode())
	if c := result.GetHistoryOrderComment(); c != "" {
		fmt.Printf(" | Comment: %s", c)
	}
	fmt.Println()
}

// ShowOrderCloseByExample closes one order using an opposite order.
func (s *MT4Service) ShowOrderCloseByExample(ctx context.Context, ticket int32, oppositeTicket int32) {
	fmt.Println("=== Close Order By Opposite ===")
	data, err := s.account.OrderCloseBy(ctx, ticket, oppositeTicket)
	if err != nil {
		log.Printf("❌ OrderCloseBy error: %v", err)
		return
	}
	fmt.Printf("✅ Closed by opposite: Profit=%.2f, Price=%.5f, Time: %s\n",
		data.GetProfit(), data.GetClosePrice(), data.GetCloseTime().AsTime().Format("2006-01-02 15:04:05"),
	)
}

// ShowOrderDeleteExample deletes a pending order by ticket.
func (s *MT4Service) ShowOrderDeleteExample(ctx context.Context, ticket int32) {
	fmt.Println("=== Delete Pending Order ===")
	data, err := s.account.OrderDelete(ctx, ticket)
	if err != nil {
		log.Printf("❌ OrderDelete error: %v", err)
		return
	}
	fmt.Printf("✅ Order deleted. Mode: %s, Comment: %s\n", data.GetMode(), data.GetHistoryOrderComment())
}

// ShowQuote retrieves and prints the current Bid/Ask quote for the given symbol.
func (s *MT4Service) ShowQuote(ctx context.Context, symbol string) {
	fmt.Printf("=== Current Quote for %s ===\n", symbol)
	data, err := s.account.Quote(ctx, symbol)
	if err != nil {
		log.Printf("❌ Quote error: %v", err)
		return
	}
	fmt.Printf("✅ Symbol: %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
		symbol, data.GetBid(), data.GetAsk(),
		data.GetDateTime().AsTime().Format("2006-01-02 15:04:05"),
	)
}

// ShowQuotesMany retrieves and displays quotes for multiple symbols.
func (s *MT4Service) ShowQuotesMany(ctx context.Context, symbols []string) {
	fmt.Println("=== Quotes for Multiple Symbols ===")
	data, err := s.account.QuoteMany(ctx, symbols)
	if err != nil {
		log.Printf("❌ QuoteMany error: %v", err)
		return
	}
	for _, q := range data.GetQuotes() {
		fmt.Printf("📈 Symbol: %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
			q.GetSymbol(), q.GetBid(), q.GetAsk(),
			q.GetDateTime().AsTime().Format("2006-01-02 15:04:05"),
		)
	}
}

// ShowQuoteHistory displays historical OHLC data for a given symbol.
func (s *MT4Service) ShowQuoteHistory(ctx context.Context, symbol string) {
	fmt.Printf("=== Historical Quote History for %s ===\n", symbol)
	from := time.Now().AddDate(0, 0, -5)
	to := time.Now()
	timeframe := pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QhPeriodH1
	data, err := s.account.QuoteHistory(ctx, symbol, timeframe, from, to)
	if err != nil {
		log.Printf("❌ QuoteHistory error: %v", err)
		return
	}
	for _, c := range data.GetHistoricalQuotes() {
		fmt.Printf("[%s] O: %.5f H: %.5f L: %.5f C: %.5f\n",
			c.GetTime().AsTime().Format("2006-01-02 15:04:05"),
			c.GetOpen(), c.GetHigh(), c.GetLow(), c.GetClose(),
		)
	}
}

// ShowAllSymbols retrieves and displays all available trading symbols.
func (s *MT4Service) ShowAllSymbols(ctx context.Context) {
	fmt.Println("=== All Available Symbols ===")
	data, err := s.account.Symbols(ctx)
	if err != nil {
		log.Printf("❌ Symbols error: %v", err)
		return
	}
	for _, e := range data.Get
}


// ShowTickValues fetches and displays tick value, tick size, and contract size for each symbol.
func (s *MT4Service) ShowTickValues(ctx context.Context, symbols []string) {
    fmt.Println("=== Tick Value, Size, and Contract Size ===")

    data, err := s.account.TickValueWithSize(ctx, symbols)
    if err != nil {
        log.Printf("❌ TickValueWithSize error: %v", err)
        return
    }
    for _, info := range data.GetInfos() {
        fmt.Printf("💹 Symbol: %s\n", info.GetSymbolName())
        fmt.Printf("  TickValue: %.5f\n", info.GetTradeTickValue())
        fmt.Printf("  TickSize: %.5f\n", info.GetTradeTickSize())
        fmt.Printf("  ContractSize: %.2f\n\n", info.GetTradeContractSize())
    }
}

// === Streaming (Streaming data) ===

// StreamQuotes subscribes to live tick updates.
func (s *MT4Service) StreamQuotes(ctx context.Context) {
    symbols := []string{"EURUSD", "GBPUSD"}
    tickCh, errCh := s.account.OnSymbolTick(ctx, symbols)

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    fmt.Println("🔄 Streaming ticks...")

    for {
        select {
        case tick, ok := <-tickCh:
            if !ok {
                fmt.Println("✅ Tick stream ended.")
                return
            }
            if sym := tick.GetSymbolTick(); sym != nil {
                fmt.Printf("[Tick] %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
                    sym.GetSymbol(), sym.GetBid(), sym.GetAsk(),
                    sym.GetTime().AsTime().Format("2006-01-02 15:04:05"))
            }
        case err := <-errCh:
            log.Printf("❌ Stream error: %v", err)
            return
        case <-time.After(30 * time.Second):
            fmt.Println("⏱️ Timeout reached.")
            return
        }
    }
}

// StreamOpenedOrderProfits streams live profit updates for open orders.
func (s *MT4Service) StreamOpenedOrderProfits(ctx context.Context) {
    profitCh, errCh := s.account.OnOpenedOrdersProfit(ctx, 1000)

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    fmt.Println("🔄 Streaming order profits...")
    for {
        select {
        case profit, ok := <-profitCh:
            if !ok {
                fmt.Println("✅ Profit stream ended.")
                return
            }
            info := profit.GetProfitInfo()
            fmt.Printf("[Profit] Ticket: %d | Symbol: %s | Profit: %.2f\n",
                info.GetTicket(), info.GetSymbol(), info.GetProfit())
        case err := <-errCh:
            log.Printf("❌ Stream error: %v", err)
            return
        case <-time.After(30 * time.Second):
            fmt.Println("⏱️ Timeout reached.")
            return
        }
    }
}

// StreamOpenedOrderTickets streams real‑time updates of open order tickets.
func (s *MT4Service) StreamOpenedOrderTickets(ctx context.Context) {
    ticketCh, errCh := s.account.OnOpenedOrdersTickets(ctx, 1000)

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    fmt.Println("🔄 Streaming opened order tickets...")
    for {
        select {
        case pkt, ok := <-ticketCh:
            if !ok {
                fmt.Println("✅ Ticket stream ended.")
                return
            }
            tix := pkt.GetTickets()
            fmt.Printf("[Tickets] %d open tickets: %v\n", len(tix), tix)
        case err := <-errCh:
            log.Printf("❌ Stream error: %v", err)
            return
        case <-time.After(30 * time.Second):
            fmt.Println("⏱️ Timeout reached.")
            return
        }
    }
}

// StreamTradeUpdates streams real‑time trade events.
func (s *MT4Service) StreamTradeUpdates(ctx context.Context) {
    tradeCh, errCh := s.account.OnTrade(ctx)

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    fmt.Println("🔄 Streaming trade updates...")
    for {
        select {
        case trade, ok := <-tradeCh:
            if !ok {
                fmt.Println("✅ Trade stream ended.")
                return
            }
            info := trade.GetTradeInfo()
            fmt.Printf("[Trade] Ticket: %d | Symbol: %s | Type: %v | Volume: %.2f | Profit: %.2f\n",
                info.GetTicket(), info.GetSymbol(), info.GetOrderType(),
                info.GetLots(), info.GetProfit())
        case err := <-errCh:
            log.Printf("❌ Stream error: %v", err)
            return
        case <-time.After(30 * time.Second):
            fmt.Println("⏱️ Timeout reached.")
            return
        }
    }
}
