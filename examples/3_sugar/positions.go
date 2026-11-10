package sugar

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// Positions covers MT4Sugar groups 8-9: watchers & waiters, and order info
// (positions, exposure, profit). All read-only.
func Positions(ctx context.Context, sugar *mt4.MT4Sugar, symbol string) {
	fmt.Println("\n=== 🍬 MT4Sugar · Positions & Watchers ===")

	// --- Order info ---
	// In MT4 an open MARKET order IS the position; pendings live in the same list.
	// Orders returns everything; Positions filters to the market ones.
	orders, _ := sugar.Orders(ctx)
	positions, _ := sugar.Positions(ctx)
	fmt.Printf("  Orders: %d  Positions: %d\n", len(orders), len(positions))

	// Filtered views — by symbol or by the EA magic number you stamped on entry.
	if ofs, err := sugar.OrdersForSymbol(ctx, symbol); err == nil {
		fmt.Printf("  OrdersForSymbol %s: %d\n", symbol, len(ofs))
	}
	if obm, err := sugar.OrdersByMagic(ctx, 777); err == nil {
		fmt.Printf("  OrdersByMagic(777): %d\n", len(obm))
	}
	// Look one order up by ticket (nil if it isn't open).
	if len(positions) > 0 {
		if o, err := sugar.FindTicket(ctx, positions[0].GetTicket()); err == nil && o != nil {
			fmt.Printf("  FindTicket #%d: %s %.2f lot\n", o.GetTicket(), o.GetSymbol(), o.GetLots())
		}
	}
	// Money views. TotalProfit/PositionsValue sum profit+swap+commission; pass ""
	// / -1 to not filter on symbol / magic.
	if tp, err := sugar.TotalProfit(ctx, ""); err == nil {
		fmt.Printf("  TotalProfit: %+.2f\n", tp)
	}
	if pv, err := sugar.PositionsValue(ctx, "", -1); err == nil {
		fmt.Printf("  PositionsValue: %+.2f\n", pv)
	}
	// Net lots per symbol (buy positive, sell negative) — your directional risk.
	if exp, err := sugar.ExposureSummary(ctx); err == nil {
		fmt.Printf("  ExposureSummary: %v\n", exp)
	}
	if a, err := sugar.GetAccountInfo(ctx); err == nil {
		fmt.Printf("  GetAccountInfo: #%d equity %.2f\n", a.GetAccountLogin(), a.GetAccountEquity())
	}

	// --- Watchers (channel-based, bounded) ---
	// A watcher pushes live frames onto a channel in the background. We bound each
	// one with a timeout context so the demo returns; in real code you'd keep the
	// context alive and range over the channel. readStream reads N frames or until
	// the timeout, then stops (cancel() drops the subscription).
	wctx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	qd, qe := sugar.SubscribeQuotes(wctx, []string{symbol}) // live ticks
	common.ReadStream("SubscribeQuotes", qd, qe, 2, 5*time.Second, func(d *pb.OnSymbolTickData) string {
		t := d.GetSymbolTick()
		return fmt.Sprintf("%s %.5f/%.5f", t.GetSymbol(), t.GetBid(), t.GetAsk())
	})
	pd, pe := sugar.WatchProfit(wctx, 1000) // per-order floating profit every 1s
	common.ReadStream("WatchProfit", pd, pe, 1, 4*time.Second, func(d *pb.OnOpenedOrdersProfitData) string {
		return fmt.Sprintf("%d orders", len(d.GetOpenedOrdersWithProfitUpdated()))
	})
	td, te := sugar.WatchTickets(wctx, 1000) // the open-ticket list every 1s
	common.ReadStream("WatchTickets", td, te, 1, 4*time.Second, func(d *pb.OnOpenedOrdersTicketsData) string {
		return fmt.Sprintf("%d position tickets", len(d.GetPositionTickets()))
	})

	// --- Waiters (poll until a condition holds, or time out) ---
	// Unlike watchers, a waiter blocks and returns once. Here we pick a target
	// just above the current ask so "ask <= target" is already true and it
	// resolves immediately — just to show the shape of the call.
	if q, err := sugar.LastQuote(ctx, symbol); err == nil {
		if hit, err := sugar.WaitPrice(ctx, symbol, q.GetAsk()+1.0, false, 2*time.Second); err == nil {
			fmt.Printf("  WaitPrice: hit at ask %.5f\n", hit.GetAsk())
		}
	}
	fmt.Println("  (WaitFilled / WaitClosed are demonstrated in the Strategies demo)")
}
