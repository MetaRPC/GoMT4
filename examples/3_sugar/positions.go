/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/3_sugar/positions.go - HIGH-LEVEL (MT4Sugar) POSITIONS & WATCHERS
 PURPOSE:
   MT4Sugar groups 8-9 - order info (positions, exposure, profit) plus watchers
   and waiters. All read-only.

 🎯 WHO SHOULD USE THIS:
   • Developers building position / exposure dashboards
   • Reacting to live ticks / profit / ticket updates
   • Learning the channel-based watchers and the poll-until waiters

 📚 WHAT THIS DEMO COVERS (3 sections):

   1. ORDER INFO
      • Orders() / Positions() / OrdersForSymbol() / OrdersByMagic() / FindTicket()
      • TotalProfit() / PositionsValue() / ExposureSummary() / GetAccountInfo()

   2. WATCHERS (channel-based, bounded)
      • SubscribeQuotes() / WatchProfit() / WatchTickets()

   3. WAITERS (poll until a condition holds, or time out)
      • WaitPrice()                          (WaitFilled/WaitClosed are in Strategies)

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - raw gRPC calls
   MID-LEVEL  (MT4Service) - typed wrappers
   HIGH-LEVEL (MT4Sugar)   - THIS FILE: one-liner helpers with smart defaults

 ℹ️  SAFE: read-only - only reads positions and opens short-lived watch streams.

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 33 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

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

	//#region 1. ORDER INFO
	// ══════════════════════════════════════════════════════════════
	// 1. ORDER INFO
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// Orders() / Positions()
	//   In MT4 an open MARKET order IS the position; pendings live in
	//   the same list. Orders returns everything; Positions filters to
	//   the market ones.
	// ──────────────────────────────────────────────────────────────
	orders, _ := sugar.Orders(ctx)
	positions, _ := sugar.Positions(ctx)
	fmt.Printf("  Orders: %d  Positions: %d\n", len(orders), len(positions))

	// ──────────────────────────────────────────────────────────────
	// OrdersForSymbol() / OrdersByMagic() / FindTicket()
	//   Filtered views — by symbol or by the EA magic number you
	//   stamped on entry — and a single ticket lookup.
	// ──────────────────────────────────────────────────────────────
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

	// ──────────────────────────────────────────────────────────────
	// TotalProfit() / PositionsValue() / ExposureSummary() / GetAccountInfo()
	//   Money views (sum profit+swap+commission; "" / -1 = don't filter),
	//   net lots per symbol (buy +, sell -), and an account snapshot.
	// ──────────────────────────────────────────────────────────────
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
	//#endregion 1. ORDER INFO

	//#region 2. WATCHERS
	// ══════════════════════════════════════════════════════════════
	// 2. WATCHERS (channel-based, bounded)
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// SubscribeQuotes() / WatchProfit() / WatchTickets()
	//   A watcher pushes live frames onto a channel in the background.
	//   We bound each with a timeout context so the demo returns; in
	//   real code you'd keep the context alive and range the channel.
	//   readStream reads N frames or until the timeout, then stops.
	// ──────────────────────────────────────────────────────────────
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
	//#endregion 2. WATCHERS

	//#region 3. WAITERS
	// ══════════════════════════════════════════════════════════════
	// 3. WAITERS (poll until a condition holds, or time out)
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// WaitPrice()
	//   Unlike watchers, a waiter blocks and returns once. Here we pick
	//   a target just above the current ask so "ask <= target" is already
	//   true and it resolves immediately — just to show the call shape.
	// ──────────────────────────────────────────────────────────────
	if q, err := sugar.LastQuote(ctx, symbol); err == nil {
		if hit, err := sugar.WaitPrice(ctx, symbol, q.GetAsk()+1.0, false, 2*time.Second); err == nil {
			fmt.Printf("  WaitPrice: hit at ask %.5f\n", hit.GetAsk())
		}
	}
	fmt.Println("  (WaitFilled / WaitClosed are demonstrated in the Strategies demo)")
	//#endregion 3. WAITERS
}
