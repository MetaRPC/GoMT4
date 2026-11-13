/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/3_sugar/trading.go - HIGH-LEVEL (MT4Sugar) TRADING DEMO
 PURPOSE:
   MT4Sugar groups 5-7 - order entry by side, order management, and
   trailing/breakeven. Walks the full lifecycle of one trade so you can see how
   the sugar helpers compose.

 🎯 WHO SHOULD USE THIS:
   • Developers wanting the shortest path to a managed trade
   • Learning the open → bracket → manage → close lifecycle
   • Quick trading scripts and bots

 📚 WHAT THIS DEMO COVERS (5 sections + a safety guard):

   1. ENTRY        • BuyMarket()                     - market buy at the ask
   2. BRACKET      • ModifySLTPByPips()              - SL/TP as pip distances
   3. MANAGE       • AutoBreakeven() / TrailStop()   - tighten the stop
   4. EXIT         • ClosePartial() / Close()        - trim then flatten
   5. BULK HELPERS • CloseAll() / CancelPendings()   - clean-slate helpers

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - raw gRPC calls
   MID-LEVEL  (MT4Service) - typed wrappers
   HIGH-LEVEL (MT4Sugar)   - THIS FILE: one-liner helpers with smart defaults

 ⚠️  WARNING: sends REAL orders on a demo account. Gated by MT4_ENABLE_TRADING=1
     (skipped otherwise); self-cleaning — it closes everything it opens.

 🚀 HOW TO RUN:
   cd examples && MT4_ENABLE_TRADING=1 go run .        (then choose 32 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

package sugar

import (
	"context"
	"fmt"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
)

// Trading covers MT4Sugar groups 5-7: order entry by side, order management, and
// trailing/breakeven. It walks the full lifecycle of one trade so you can see how
// the sugar helpers compose. Gated behind MT4_ENABLE_TRADING=1.
func Trading(ctx context.Context, sugar *mt4.MT4Sugar, symbol string) {
	fmt.Println("\n=== 🍬 MT4Sugar · Trading (entry · manage · trail) ===")

	//#region SAFETY GUARD
	// ══════════════════════════════════════════════════════════════
	// SAFETY GUARD (opt in with MT4_ENABLE_TRADING=1)
	// ══════════════════════════════════════════════════════════════

	// Reads are safe to run anywhere; sending orders is not. The whole block
	// below is skipped unless you explicitly opt in with MT4_ENABLE_TRADING=1.
	if !common.TradingEnabled {
		fmt.Println("  (trading demos skipped — set MT4_ENABLE_TRADING=1)")
		fmt.Println("  Available: BuyMarket/SellMarket, BuyLimit/SellLimit, BuyStop/SellStop,")
		fmt.Println("             Close/ClosePartial/CloseBy/CloseAll/CancelPendings,")
		fmt.Println("             ModifySLTPByPrice/ByPips, AutoBreakeven, TrailStop")
		return
	}

	fmt.Println("  BuyMarket 0.01 → pip SL/TP → breakeven → trail → close")
	//#endregion SAFETY GUARD

	//#region 1. ENTRY
	// ══════════════════════════════════════════════════════════════
	// 1. ENTRY
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// BuyMarket()
	//   Open a market buy at the current ask. We pass nil for sl/tp
	//   here and attach them by pips in the next step. The default
	//   magic / slippage / comment (from SetDefaults) are stamped
	//   automatically. A rejected trade comes back as an error.
	// ──────────────────────────────────────────────────────────────
	r, err := sugar.BuyMarket(ctx, symbol, 0.01, nil, nil)
	if err != nil {
		// A rejected trade comes back as an error (e.g. market closed, not enough
		// money) — there is no "status code" to inspect separately.
		fmt.Println("  BuyMarket error:", err)
		return
	}
	t := r.GetTicket() // the ticket identifies this order for every later call
	fmt.Printf("  BuyMarket #%d @ %.5f\n", t, r.GetPrice())
	//#endregion 1. ENTRY

	//#region 2. BRACKET
	// ══════════════════════════════════════════════════════════════
	// 2. BRACKET
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// ModifySLTPByPips()
	//   Set the protective stop/target as PIP DISTANCES from the entry.
	//   Direction-aware: for a BUY the stop goes below and the target
	//   above the open price (mirrored for a SELL). common.Ptr(x) takes
	//   the address of the literal (the args are optional *float64).
	// ──────────────────────────────────────────────────────────────
	if err := sugar.ModifySLTPByPips(ctx, t, common.Ptr(20.0), common.Ptr(40.0)); err != nil {
		fmt.Println("  ModifySLTPByPips error:", err)
	}
	//#endregion 2. BRACKET

	//#region 3. MANAGE
	// ══════════════════════════════════════════════════════════════
	// 3. MANAGE
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// AutoBreakeven() / TrailStop()
	//   One-shot primitives — in real code you'd call them on a timer.
	//   AutoBreakeven moves the stop to entry+2p once +10p in profit;
	//   TrailStop then follows price 15p behind. Both only ever tighten
	//   the stop, so calling them early simply returns moved=false.
	// ──────────────────────────────────────────────────────────────
	if moved, err := sugar.AutoBreakeven(ctx, t, 10, 2); err == nil {
		fmt.Printf("  AutoBreakeven moved: %v\n", moved)
	}
	if moved, err := sugar.TrailStop(ctx, t, 15); err == nil {
		fmt.Printf("  TrailStop moved: %v\n", moved)
	}
	//#endregion 3. MANAGE

	//#region 4. EXIT
	// ══════════════════════════════════════════════════════════════
	// 4. EXIT
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// ClosePartial() / Close()
	//   ClosePartial trims part of the position; Close(…, nil) closes
	//   the rest in full. (On a 0.01 lot the partial may be a no-op
	//   below min volume.)
	// ──────────────────────────────────────────────────────────────
	if _, err := sugar.ClosePartial(ctx, t, 0.01); err == nil {
		fmt.Printf("  ClosePartial #%d ok\n", t)
	}
	if _, err := sugar.Close(ctx, t, nil); err == nil {
		fmt.Printf("  Close #%d ok\n", t)
	}
	//#endregion 4. EXIT

	//#region 5. BULK HELPERS
	// ══════════════════════════════════════════════════════════════
	// 5. BULK HELPERS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// CloseAll() / CancelPendings()
	//   Handy for a clean slate — they iterate the open book for you.
	//   Both are no-ops (return 0) when nothing matches.
	// ──────────────────────────────────────────────────────────────
	if n, err := sugar.CloseAll(ctx, symbol); err == nil {
		fmt.Printf("  CloseAll(%s): %d closed\n", symbol, n) // close every market order on the symbol
	}
	if n, err := sugar.CancelPendings(ctx, ""); err == nil {
		fmt.Printf("  CancelPendings: %d deleted\n", n) // delete every pending (all symbols)
	}
	//#endregion 5. BULK HELPERS
}
