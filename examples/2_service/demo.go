/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/2_service/demo.go - MID-LEVEL MT4 SERVICE (MT4Service) DEMO
 PURPOSE:
   Exercises the MT4Service facade - friendly arguments and unwrapped
   (slice/scalar) return values across every group: account info, orders,
   symbols & quotes, streaming, and (optional) trading.

 🎯 WHO SHOULD USE THIS:
   • Developers who want the typed facade over the raw MT4Account
   • Users who prefer slices/scalars to protobuf wrapper messages
   • Learning the mid-level API before the one-liner MT4Sugar

 📚 WHAT THIS DEMO COVERS (5 sections):

   1. ACCOUNT INFO
      • AccountSummary() / TradeStats() / TradeStatsEquityHistory()
      • EquityHistory() / TradeUnrealizedPL()

   2. ORDERS
      • OpenedOrders() / OpenedOrdersTickets() / OrdersHistory()

   3. SYMBOLS & QUOTES
      • Symbols() / SymbolSelect() / Quote() / QuoteMany()
      • SymbolParamsMany() / TickValueWithSize() / QuoteHistory() / ServerTimezone()

   4. STREAMING
      • StreamTicks()                        - bounded live Bid/Ask stream

   5. TRADING (optional)
      • OrderSend() / OrderModify() / OrderClose() - market BUY 0.01, self-closed

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - raw gRPC calls
   MID-LEVEL  (MT4Service) - THIS FILE: typed facade, friendly args, clean returns
   HIGH-LEVEL (MT4Sugar)   - one-liner helpers with smart defaults

 ⚠️  WARNING: section 5 sends a real 0.01 market order (BUY -> modify -> close,
     self-cleaning) and runs only when MT4_ENABLE_TRADING=1.

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 21 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

// Package service holds the MT4Service (mid-level facade) example.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// Demo exercises the MT4Service facade — one file showing the friendly arguments
// and unwrapped (slice/scalar) return values across every group.
func Demo(ctx context.Context, svc *mt4.MT4Service, symbol string) {
	fmt.Println("\n=== 🧩 MT4Service · facade (friendly args, clean returns) ===")

	//#region 1. ACCOUNT INFO
	// ══════════════════════════════════════════════════════════════
	// 1. ACCOUNT INFO
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// AccountSummary() / TradeStats() / TradeStatsEquityHistory() / EquityHistory() / TradeUnrealizedPL()
	//   Money snapshot plus the aggregated analytics, via the facade.
	// ──────────────────────────────────────────────────────────────
	if a, err := svc.AccountSummary(ctx); err == nil {
		fmt.Printf("  AccountSummary: #%d  %.2f %s\n", a.GetAccountLogin(), a.GetAccountBalance(), a.GetAccountCurrency())
	}
	if _, err := svc.TradeStats(ctx, nil); err == nil {
		fmt.Println("  TradeStats: ok")
	}
	if _, err := svc.TradeStatsEquityHistory(ctx, nil, nil); err == nil {
		fmt.Println("  TradeStatsEquityHistory: ok")
	}
	if _, err := svc.EquityHistory(ctx, nil, nil); err == nil {
		fmt.Println("  EquityHistory: ok")
	}
	if _, err := svc.TradeUnrealizedPL(ctx, nil); err == nil {
		fmt.Println("  TradeUnrealizedPL: ok")
	}
	//#endregion 1. ACCOUNT INFO

	//#region 2. ORDERS
	// ══════════════════════════════════════════════════════════════
	// 2. ORDERS (facade unwraps to slices)
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// OpenedOrders() / OpenedOrdersTickets() / OrdersHistory()
	//   Open orders as a ready slice, their tickets, and 30 days of
	//   closed-order history.
	// ──────────────────────────────────────────────────────────────
	if orders, err := svc.OpenedOrders(ctx); err == nil {
		fmt.Printf("  OpenedOrders: %d (already []*OpenedOrderInfo)\n", len(orders))
	}
	if tickets, err := svc.OpenedOrdersTickets(ctx); err == nil {
		fmt.Printf("  OpenedOrdersTickets: %v\n", tickets)
	}
	now := time.Now()
	from := now.Add(-30 * 24 * time.Hour)
	if _, err := svc.OrdersHistory(ctx, pb.EnumOrderHistorySortType(0), &from, &now); err == nil {
		fmt.Println("  OrdersHistory (30d): ok")
	}
	//#endregion 2. ORDERS

	//#region 3. SYMBOLS & QUOTES
	// ══════════════════════════════════════════════════════════════
	// 3. SYMBOLS & QUOTES
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// Symbols() / SymbolSelect() / Quote() / QuoteMany()
	//   The symbol catalog, Market Watch select, and Bid/Ask for one
	//   or several symbols.
	// ──────────────────────────────────────────────────────────────
	if syms, err := svc.Symbols(ctx); err == nil {
		fmt.Printf("  Symbols: %d (slice)\n", len(syms))
	}
	if ok, err := svc.SymbolSelect(ctx, symbol, true); err == nil {
		fmt.Printf("  SymbolSelect %s: %v\n", symbol, ok)
	}
	if q, err := svc.Quote(ctx, symbol); err == nil {
		fmt.Printf("  Quote %s: %.5f / %.5f\n", symbol, q.GetBid(), q.GetAsk())
	}
	if _, err := svc.QuoteMany(ctx, []string{"EURUSD", "GBPUSD"}); err == nil {
		fmt.Println("  QuoteMany: ok")
	}

	// ──────────────────────────────────────────────────────────────
	// SymbolParamsMany() / TickValueWithSize() / QuoteHistory() / ServerTimezone()
	//   Trading params, tick value, H1 history and the server timezone.
	// ──────────────────────────────────────────────────────────────
	if _, err := svc.SymbolParamsMany(ctx, []string{symbol}); err == nil {
		fmt.Println("  SymbolParamsMany: ok")
	}
	if _, err := svc.TickValueWithSize(ctx, []string{symbol}); err == nil {
		fmt.Println("  TickValueWithSize: ok")
	}
	if _, err := svc.QuoteHistory(ctx, symbol, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, now.Add(-24*time.Hour), now); err == nil {
		fmt.Println("  QuoteHistory: ok")
	}
	if tz, err := svc.ServerTimezone(ctx); err == nil {
		fmt.Printf("  ServerTimezone: %s\n", tz.GetServer())
	}
	//#endregion 3. SYMBOLS & QUOTES

	//#region 4. STREAMING
	// ══════════════════════════════════════════════════════════════
	// 4. STREAMING (clean names, bounded)
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// StreamTicks()
	//   Live Bid/Ask stream, read a couple of ticks with a 5s cap.
	// ──────────────────────────────────────────────────────────────
	sctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	d, e := svc.StreamTicks(sctx, []string{symbol})
	common.ReadStream("StreamTicks", d, e, 2, 5*time.Second, func(x *pb.OnSymbolTickData) string {
		t := x.GetSymbolTick()
		return fmt.Sprintf("%s %.5f/%.5f", t.GetSymbol(), t.GetBid(), t.GetAsk())
	})
	//#endregion 4. STREAMING

	//#region 5. TRADING
	// ══════════════════════════════════════════════════════════════
	// 5. TRADING (⚠️ real order; MT4_ENABLE_TRADING=1 only)
	// ══════════════════════════════════════════════════════════════

	if !common.TradingEnabled {
		fmt.Println("  (trading demos skipped — set MT4_ENABLE_TRADING=1)")
		return
	}
	// At the facade layer trading is the same shape as the raw account (the facade
	// passes it straight through) — side-based one-liners with defaults live one
	// level up in MT4Sugar (BuyMarket/…). Here: market BUY 0.01 → modify → close.

	// ──────────────────────────────────────────────────────────────
	// OrderSend() / OrderModify() / OrderClose()
	//   Open a 0.01 market BUY, show a no-op modify, then close it.
	// ──────────────────────────────────────────────────────────────
	fmt.Println("  ⚠️  Trading: OrderSend 0.01 BUY → modify → close")
	if r, err := svc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil); err == nil {
		t := r.GetTicket() // ticket of the new order
		fmt.Printf("  OrderSend #%d\n", t)
		_, _ = svc.OrderModify(ctx, t, nil, nil, nil, nil) // nil args = no-op, shown for shape
		if _, err := svc.OrderClose(ctx, t, nil, nil, nil); err == nil {
			fmt.Printf("  OrderClose #%d ok\n", t) // lots nil = full close
		}
	}
	//#endregion 5. TRADING
}
