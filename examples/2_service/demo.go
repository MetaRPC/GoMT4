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

	// Account info.
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

	// Orders — facade unwraps to slices.
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

	// Symbols & quotes.
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

	// Streaming (clean names, bounded).
	sctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	d, e := svc.StreamTicks(sctx, []string{symbol})
	common.ReadStream("StreamTicks", d, e, 2, 5*time.Second, func(x *pb.OnSymbolTickData) string {
		t := x.GetSymbolTick()
		return fmt.Sprintf("%s %.5f/%.5f", t.GetSymbol(), t.GetBid(), t.GetAsk())
	})

	if !common.TradingEnabled {
		fmt.Println("  (trading demos skipped — set MT4_ENABLE_TRADING=1)")
		return
	}
	// At the facade layer trading is the same shape as the raw account (the facade
	// passes it straight through) — side-based one-liners with defaults live one
	// level up in MT4Sugar (BuyMarket/…). Here: market BUY 0.01 → modify → close.
	fmt.Println("  ⚠️  Trading: OrderSend 0.01 BUY → modify → close")
	if r, err := svc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil); err == nil {
		t := r.GetTicket() // ticket of the new order
		fmt.Printf("  OrderSend #%d\n", t)
		_, _ = svc.OrderModify(ctx, t, nil, nil, nil, nil) // nil args = no-op, shown for shape
		if _, err := svc.OrderClose(ctx, t, nil, nil, nil); err == nil {
			fmt.Printf("  OrderClose #%d ok\n", t) // lots nil = full close
		}
	}
}
