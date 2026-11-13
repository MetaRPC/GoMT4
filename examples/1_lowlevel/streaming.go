/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/1_lowlevel/streaming.go - LOW-LEVEL MT4 STREAMING DEMO
 PURPOSE:
   Exercises the MT4Account server streams, each bounded by a short timeout /
   frame cap so the demo returns. Read-only - it only subscribes, no orders.

 🎯 WHO SHOULD USE THIS:
   • Developers learning the MT4 gRPC streaming API
   • Building live tick / order / P&L monitors
   • Seeing how bounded stream reads are structured

 📚 WHAT THIS DEMO COVERS (2 sections):

   1. MARKET STREAM
      • OnSymbolTick()            - live Bid/Ask ticks for a symbol

   2. ACCOUNT & TRADE STREAMS
      • OnOpenedOrdersProfit()    - periodic P/L of open orders
      • OnOpenedOrdersTickets()   - periodic open-position ticket list
      • OnTrade()                 - fires on trade events
      • OnConnectState()          - terminal connection-state changes

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - THIS FILE: raw gRPC streams
   MID-LEVEL  (MT4Service) - typed wrappers, less boilerplate
   HIGH-LEVEL (MT4Sugar)   - one-liner helpers with smart defaults

 ℹ️  SAFE: read-only streaming, opens no orders.

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 13 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

package lowlevel

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// Streaming exercises the MT4Account server streams, each bounded by a short
// timeout / frame cap so the demo returns.
func Streaming(ctx context.Context, acc *mt4.MT4Account, symbol string) {
	fmt.Println("\n=== 📂 MT4Account · Streaming (bounded) ===")

	sctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	//#region 1. MARKET STREAM
	// ══════════════════════════════════════════════════════════════
	// 1. MARKET STREAM
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// OnSymbolTick()
	//   Live Bid/Ask ticks for the symbol (first 3 frames / 6s cap).
	// ──────────────────────────────────────────────────────────────
	d1, e1 := acc.OnSymbolTick(sctx, []string{symbol})
	common.ReadStream("OnSymbolTick", d1, e1, 3, 6*time.Second, func(d *pb.OnSymbolTickData) string {
		t := d.GetSymbolTick()
		return fmt.Sprintf("%s bid %.5f ask %.5f", t.GetSymbol(), t.GetBid(), t.GetAsk())
	})
	//#endregion 1. MARKET STREAM

	//#region 2. ACCOUNT & TRADE STREAMS
	// ══════════════════════════════════════════════════════════════
	// 2. ACCOUNT & TRADE STREAMS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// OnOpenedOrdersProfit() / OnOpenedOrdersTickets()
	//   Periodic P/L of open orders and the open-position ticket list.
	// ──────────────────────────────────────────────────────────────
	d2, e2 := acc.OnOpenedOrdersProfit(sctx, 1000)
	common.ReadStream("OnOpenedOrdersProfit", d2, e2, 2, 5*time.Second, func(d *pb.OnOpenedOrdersProfitData) string {
		return fmt.Sprintf("%d orders", len(d.GetOpenedOrdersWithProfitUpdated()))
	})

	d3, e3 := acc.OnOpenedOrdersTickets(sctx, 1000)
	common.ReadStream("OnOpenedOrdersTickets", d3, e3, 2, 5*time.Second, func(d *pb.OnOpenedOrdersTicketsData) string {
		return fmt.Sprintf("%d position tickets", len(d.GetPositionTickets()))
	})

	// ──────────────────────────────────────────────────────────────
	// OnTrade() / OnConnectState()
	//   Trade events, and terminal connection-state changes.
	// ──────────────────────────────────────────────────────────────
	d4, e4 := acc.OnTrade(sctx)
	common.ReadStream("OnTrade", d4, e4, 1, 3*time.Second, func(d *pb.OnTradeData) string { return "trade event" })

	d5, e5 := acc.OnConnectState(sctx, nil)
	common.ReadStream("OnConnectState", d5, e5, 1, 3*time.Second, func(d *pb.ConnectStateData) string { return d.GetState() })
	//#endregion 2. ACCOUNT & TRADE STREAMS
}
