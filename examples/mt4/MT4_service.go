package mt4

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "git.mtapi.io/root/mrpc-proto.git/mt4/libraries/go"
)

type MT4Service struct {
	account *MT4Account
}

func NewMT4Service(acc *MT4Account) *MT4Service {
	return &MT4Service{account: acc}
}

// === 📂 Account Info ===

// ShowAccountSummary fetches and prints the account's balance, equity, and currency.
//
// Parameters:
//   - ctx: Context for timeout or cancellation.
//
// If the request fails, logs the error. Otherwise, prints key account metrics to stdout.
// Intended for CLI or diagnostic output; does not return any data.

func (s *MT4Service) ShowAccountSummary(ctx context.Context) {
	summary, err := s.account.AccountSummary(ctx)
	if err != nil {
		log.Printf("❌ AccountSummary error: %v", err)
		return
	}
	fmt.Printf("Balance: %.2f, Equity: %.2f, Currency: %s\n",
		summary.GetAccountBalance(),
		summary.GetAccountEquity(),
		summary.GetAccountCurrency())
}

// === 📂 Order Operations ===

// ShowOpenedOrders prints all currently opened orders with basic trade info.
//
// Logs an error if the request fails. Prints "No opened orders" if the list is empty.

func (s *MT4Service) ShowOpenedOrders(ctx context.Context) {
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

// ShowOpenedOrderTickets prints the ticket numbers of all open orders.
//
// Logs an error if the request fails. Shows a message if no tickets are found.

func (s *MT4Service) ShowOpenedOrderTickets(ctx context.Context) {
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

// ShowOrdersHistory prints closed orders from the last 7 days, sorted by close time.
//
// Logs an error if the request fails. Prints a message if no historical orders are found.

func (s *MT4Service) ShowOrdersHistory(ctx context.Context) {
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

// ShowOrderSendExample sends a sample BUY order with test parameters.
//
// Logs the result or prints the ticket and price if successful.

func (s *MT4Service) ShowOrderSendExample(ctx context.Context, symbol string) {
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
		data.GetTicket(), data.GetPrice(), data.GetOpenTime().AsTime().Format("2006-01-02 15:04:05"))
}

// ShowOrderModifyExample attempts to modify SL/TP for the given order ticket.
//
// Logs the result and indicates whether the modification was successful.

func (s *MT4Service) ShowOrderModifyExample(ctx context.Context, ticket int32) {
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

// ShowOrderCloseExample closes an order by ticket and prints the result.
//
// Displays close mode and optional order comment if available.

func (s *MT4Service) ShowOrderCloseExample(ctx context.Context, ticket int32) {
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

// ShowOrderCloseByExample closes an order using an opposite one and prints the result.
//
// Shows profit, close price, and time of closure.

func (s *MT4Service) ShowOrderCloseByExample(ctx context.Context, ticket int32, oppositeTicket int32) {
	data, err := s.account.OrderCloseBy(ctx, ticket, oppositeTicket)
	if err != nil {
		log.Printf("❌ OrderCloseBy error: %v", err)
		return
	}
	fmt.Printf("✅ Closed by opposite: Profit=%.2f, Price=%.5f, Time: %s\n",
		data.GetProfit(), data.GetClosePrice(), data.GetCloseTime().AsTime().Format("2006-01-02 15:04:05"))
}

// ShowOrderDeleteExample deletes a pending order by ticket and prints the result.
//
// Displays the close mode and any associated comment.

func (s *MT4Service) ShowOrderDeleteExample(ctx context.Context, ticket int32) {
	data, err := s.account.OrderDelete(ctx, ticket)
	if err != nil {
		log.Printf("❌ OrderDelete error: %v", err)
		return
	}
	fmt.Printf("✅ Order deleted. Mode: %s, Comment: %s\n", data.GetMode(), data.GetHistoryOrderComment())
}

// === 📂 Market Info / Symbol Info ===

// ShowQuote fetches and prints the latest quote (bid/ask) for a given symbol.
//
// Also shows the timestamp of the quote.

func (s *MT4Service) ShowQuote(ctx context.Context, symbol string) {
	data, err := s.account.Quote(ctx, symbol)
	if err != nil {
		log.Printf("❌ Quote error: %v", err)
		return
	}
	fmt.Printf("✅ Symbol: %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
		symbol, data.GetBid(), data.GetAsk(),
		data.GetDateTime().AsTime().Format("2006-01-02 15:04:05"))
}

// ShowQuotesMany prints the latest quotes (bid/ask) for multiple symbols.
//
// Logs an error if the request fails.

func (s *MT4Service) ShowQuotesMany(ctx context.Context, symbols []string) {
	data, err := s.account.QuoteMany(ctx, symbols)
	if err != nil {
		log.Printf("❌ QuoteMany error: %v", err)
		return
	}
	for _, q := range data.GetQuotes() {
		fmt.Printf("📈 Symbol: %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
			q.GetSymbol(), q.GetBid(), q.GetAsk(),
			q.GetDateTime().AsTime().Format("2006-01-02 15:04:05"))
	}
}

// ShowQuoteHistory prints historical OHLC data for a symbol (last 5 days, H1).
//
// Displays time, open, high, low, and close prices for each candle.

func (s *MT4Service) ShowQuoteHistory(ctx context.Context, symbol string) {
	from := time.Now().AddDate(0, 0, -5)
	to := time.Now()
	timeframe := pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1
	data, err := s.account.QuoteHistory(ctx, symbol, timeframe, from, to)
	if err != nil {
		log.Printf("❌ QuoteHistory error: %v", err)
		return
	}
	for _, c := range data.GetHistoricalQuotes() {
		fmt.Printf("[%s] O: %.5f H: %.5f L: %.5f C: %.5f\n",
			c.GetTime().AsTime().Format("2006-01-02 15:04:05"),
			c.GetOpen(), c.GetHigh(), c.GetLow(), c.GetClose())
	}
}

// ShowAllSymbols fetches and prints all symbols available in the MT4 terminal.
//
// Calls ShowAllSymbols from the MT4Account wrapper and displays name and description
// for each symbol, if available.
func (s *MT4Service) ShowAllSymbols(ctx context.Context) {
	data, err := s.account.ShowAllSymbols(ctx)
	if err != nil {
		log.Printf("❌ ShowAllSymbols error: %v", err)
		return
	}

	symbols := data.SymbolNameInfos
	if len(symbols) == 0 {
		fmt.Println("📭 No symbols found.")
		return
	}

	fmt.Println("=== 🧾 All Available Symbols ===")
	for _, sym := range symbols {
		fmt.Printf("• %s (Index: %d)\n", sym.GetSymbolName(), sym.GetSymbolIndex())
	}
}

// ShowSymbolParams prints detailed parameters for a specific trading symbol.
//
// Retrieves and displays advanced symbol information such as description, digits,
// volume limits, trade mode, and related currencies using the SymbolParamsMany gRPC method.
func (s *MT4Service) ShowSymbolParams(ctx context.Context, symbol string) error {
	info, err := s.account.SymbolParams(ctx, symbol)
	if err != nil {
		return err
	}

	fmt.Println("📊 Symbol Parameters:")
	fmt.Printf("• Symbol: %s\n", info.GetSymbolName())
	fmt.Printf("• Description: %s\n", info.GetSymDescription())
	fmt.Printf("• Digits: %d\n", info.GetDigits())
	fmt.Printf("• Volume Min: %.2f\n", info.GetVolumeMin())
	fmt.Printf("• Volume Max: %.2f\n", info.GetVolumeMax())
	fmt.Printf("• Volume Step: %.2f\n", info.GetVolumeStep())
	fmt.Printf("• Trade Mode: %s\n", tradeModeToString(info.GetTradeMode()))
	fmt.Printf("• Currency Base: %s\n", info.GetCurrencyBase())
	fmt.Printf("• Currency Profit: %s\n", info.GetCurrencyProfit())
	fmt.Printf("• Currency Margin: %s\n", info.GetCurrencyMargin())

	return nil
}

// ShowSymbols prints the available symbols along with their indices.
//
// Displays the symbol name and its corresponding index from the SymbolsData response.

func (s *MT4Service) ShowSymbols(ctx context.Context) {
	data, err := s.account.Symbols(ctx)
	if err != nil {
		log.Printf("❌ Symbols error: %v", err)
		return
	}

	fmt.Println("=== Available Symbols ===")
	for _, symbolInfo := range data.GetSymbolNameInfos() {
		fmt.Printf("Symbol: %s, Index: %d\n", symbolInfo.GetSymbolName(), symbolInfo.GetSymbolIndex())
	}
}

// ShowTickValues prints tick value, tick size, and contract size for each symbol.
//
// Useful for risk and volume calculations.

func (s *MT4Service) ShowTickValues(ctx context.Context, symbols []string) {
	data, err := s.account.TickValueWithSize(ctx, symbols)
	if err != nil {
		log.Printf("❌ TickValueWithSize error: %v", err)
		return
	}
	for _, info := range data.Infos {
		fmt.Printf("💹 Symbol: %s\n  TickValue: %.5f\n  TickSize: %.5f\n  ContractSize: %.2f\n\n",
			info.GetSymbolName(),
			info.GetTradeTickValue(),
			info.GetTradeTickSize(),
			info.GetTradeContractSize(),
		)
	}
}

// === 📂 Streaming / Subscriptions ===

// StreamQuotes subscribes to live tick updates for predefined symbols.
//
// Prints bid/ask updates until stream ends or times out.

func (s *MT4Service) StreamQuotes(ctx context.Context) {
	symbols := []string{"EURUSD", "GBPUSD"}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tickCh, errCh := s.account.OnSymbolTick(ctx, symbols)

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
					sym.GetSymbol(), sym.GetBid(), sym.GetAsk(), sym.GetTime().AsTime().Format("2006-01-02 15:04:05"))
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

// StreamOpenedOrderProfits subscribes to live profit updates for opened orders.
//
// Prints real-time profit data for each updated order, including ticket number,
// symbol, and current profit. Automatically stops on stream timeout or error.

func (s *MT4Service) StreamOpenedOrderProfits(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	profitCh, errCh := s.account.OnOpenedOrdersProfit(ctx, 1000)

	fmt.Println("🔄 Streaming order profits...")

	for {
		select {
		case profit, ok := <-profitCh:
			if !ok {
				fmt.Println("✅ Profit stream ended.")
				return
			}

			// profit.OpenedOrdersWithProfitUpdated — это массив []*OnOpenedOrdersProfitOrderInfo
			for _, info := range profit.OpenedOrdersWithProfitUpdated {
				fmt.Printf("[Profit] Ticket: %d | Symbol: %s | Profit: %.2f\n",
					info.Ticket, info.Symbol, info.OrderProfit)
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

// StreamOpenedOrderTickets subscribes to live updates of open order tickets.
//
// Prints current ticket list on each update.

func (s *MT4Service) StreamOpenedOrderTickets(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ticketCh, errCh := s.account.OnOpenedOrdersTickets(ctx, 1000)

	fmt.Println("🔄 Streaming opened order tickets...")
	for {
		select {
		case pkt, ok := <-ticketCh:
			if !ok {
				fmt.Println("✅ Ticket stream ended.")
				return
			}
			tix := append(pkt.PositionTickets, pkt.PendingOrderTickets...)
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

// StreamTradeUpdates listens for real-time trade events (opens, closes, modifies).
//
// Prints order info on each trade update.

func (s *MT4Service) StreamTradeUpdates(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tradeCh, errCh := s.account.OnTrade(ctx)

	fmt.Println("🔄 Streaming trade updates...")

	for {
		select {
		case trade, ok := <-tradeCh:
			if !ok {
				fmt.Println("✅ Trade stream ended.")
				return
			}

			info := trade.EventData
			if info == nil {
				log.Println("❌ No trade info available.")
				continue
			}

			if len(info.NewOrders) > 0 {
				order := info.NewOrders[0]
				fmt.Printf("[Trade] Ticket: %d | Symbol: %s | Type: %v | Volume: %.2f | Profit: %.2f\n",
					order.Ticket, order.Symbol, order.Type, order.Lots, order.OrderProfit)
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

// streamorderhistory — a page-by-page stream of the order history.
// Proxies the OrdersHistoryStream from MT4Account.
func (s *MT4Service) StreamOrdersHistory(
	ctx context.Context,
	sortType pb.EnumOrderHistorySortType,
	from, to *time.Time,
	pageSize int32,
) (<-chan *pb.OrdersHistoryData, <-chan error) {
	return s.account.OrdersHistoryStream(ctx, sortType, from, to, pageSize)
}

// StreamQuoteHistory — gives the history of quotes in chunks of time.
// Proxies the QuoteHistoryStream from MT4Account.
func (s *MT4Service) StreamQuoteHistory(
	ctx context.Context,
	symbol string,
	timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME,
	from, to time.Time,
	chunk time.Duration,
) (<-chan *pb.QuoteHistoryData, <-chan error) {
	return s.account.QuoteHistoryStream(ctx, symbol, timeframe, from, to, chunk)
}

// StreamOrdersHistoryExample streams account order history as paginated pages
// for a fixed demo window (last 30 days). It uses the high-level service
// wrapper s.StreamOrdersHistory which returns two channels:
//   - pagesCh: emits *pb.OrdersHistoryData (one page at a time)
//   - errCh:   emits terminal/stream errors
//
// Demo behavior:
//   - Sorts by close time (DESC) so the most recent trades come first.
//   - Uses a page size of 200 orders per page.
//   - Stops after ~30 seconds via a demo timeout (time.After).
//   - Prints each order with type, ticket, symbol and realized PnL.
//
// Production tips:
//   - Replace the fixed date range with user input or config.
//   - Remove the demo timeout and control lifetime via ctx (cancel).
//   - Batch/process pages off the main goroutine if heavy work is needed.
//   - Always read errCh to avoid goroutine leaks.
func (s *MT4Service) StreamOrdersHistoryExample(ctx context.Context) {
	// Define the historical window for the demo: 30 days back from now.
	from := time.Now().AddDate(0, 0, -30)
	to := time.Now()

	// Open the paginated history stream via the service wrapper.
	// Sort: by close time descending; Page size: 200.
	pagesCh, errCh := s.StreamOrdersHistory(
		ctx,
		pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC,
		&from, &to,
		200, // page size
	)

	fmt.Println("📥 Streaming orders history (last 30 days)...")

	for {
		select {
		case page, ok := <-pagesCh:
			if !ok {
				// pagesCh closed: the stream finished gracefully.
				fmt.Println("✅ Orders history stream finished.")
				return
			}

			// Iterate over a single page of orders and print key fields.
			for _, order := range page.GetOrdersInfo() {
				fmt.Printf("[HIST] %s | Ticket: %d | %s | PnL: %.2f\n",
					order.GetOrderType(),
					order.GetTicket(),
					order.GetSymbol(),
					order.GetProfit(),
				)
			}

		case err := <-errCh:
			// Any transport/API error during streaming comes here.
			// In production, consider retry/backoff or user notification.
			log.Printf("❌ Orders history stream error: %v", err)
			return

		case <-time.After(30 * time.Second):
			// Demo timeout to keep examples short and safe.
			// In production, remove this and rely on ctx cancellation.
			fmt.Println("⏱ Timeout reached.")
			return
		}
	}
}

// StreamQuoteHistoryExample streams historical OHLC bars in fixed-size time chunks
// for a demo window (last 90 days). It uses the service wrapper s.StreamQuoteHistory
// which returns two channels:
//   - barsCh: emits *pb.QuoteHistoryData (one chunk/batch at a time)
//   - errCh:  emits terminal/stream errors
//
// Demo behavior:
//   - Timeframe: H1 (hourly bars).
//   - Range: last 90 days from now.
//   - Chunking: 7 days per batch (weekly chunks).
//   - Stops after ~30 seconds via a demo timeout.
//   - Prints time, Open and Close for each bar.
//
// Production tips:
//   - Replace the fixed date range and chunk size with user input/config.
//   - Remove the demo timeout; control lifetime via ctx (cancel).
//   - If you need strict ordering, sort batches by time before processing.
//   - Offload heavy processing to worker goroutines; never block the reader loop.
func (s *MT4Service) StreamQuoteHistoryExample(ctx context.Context, symbol string) {
	// Define the historical window for the demo: 90 days back from now.
	from := time.Now().AddDate(0, 0, -90)
	to := time.Now()

	// Start the time-chunked history stream:
	// - timeframe: H1
	// - chunk: 1 week per emitted batch
	barsCh, errCh := s.StreamQuoteHistory(
		ctx,
		symbol,
		pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1,
		from, to,
		7*24*time.Hour, // chunks per week
	)

	fmt.Printf("📊 Streaming quote history for %s...\n", symbol)

	for {
		select {
		case batch, ok := <-barsCh:
			if !ok {
				// barsCh closed: the stream finished gracefully.
				fmt.Println("✅ Quote history stream finished.")
				return
			}

			// Iterate all bars in the current batch and print key fields.
			for _, c := range batch.GetHistoricalQuotes() {
				fmt.Printf("[%s] O: %.5f C: %.5f\n",
					c.GetTime().AsTime().Format("2006-01-02 15:04:05"),
					c.GetOpen(), c.GetClose(),
				)
			}

		case err := <-errCh:
			// Any transport/API error during streaming comes here.
			// Consider retry/backoff or surfacing to the user in production.
			log.Printf("❌ Quote history stream error: %v", err)
			return

		case <-time.After(30 * time.Second):
			// Demo timeout to keep examples short and safe.
			// In production, remove this and rely on ctx cancellation.
			fmt.Println("⏱ Timeout reached.")
			return
		}
	}
}

func tradeModeToString(mode pb.SP_ENUM_SYMBOL_TRADE_MODE) string {
	switch mode {
	case pb.SP_ENUM_SYMBOL_TRADE_MODE_SYMBOL_TRADE_MODE_DISABLED:
		return "Disabled"
	case pb.SP_ENUM_SYMBOL_TRADE_MODE_SYMBOL_TRADE_MODE_LONGONLY:
		return "Long Only"
	case pb.SP_ENUM_SYMBOL_TRADE_MODE_SYMBOL_TRADE_MODE_SHORTONLY:
		return "Short Only"
	case pb.SP_ENUM_SYMBOL_TRADE_MODE_SYMBOL_TRADE_MODE_CLOSEONLY:
		return "Close Only"
	case pb.SP_ENUM_SYMBOL_TRADE_MODE_SYMBOL_TRADE_MODE_FULL:
		return "Full Access"
	default:
		return "Unknown"
	}
}

// === 🔧Utilities ===

func ptrInt32(v int32) *int32    { return &v }
func ptrString(v string) *string { return &v }
