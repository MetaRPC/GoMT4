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

	// 1) ENTRY. BuyMarket opens a market buy at the current ask. We pass nil for
	//    sl/tp here and attach them by pips in the next step. The default magic /
	//    slippage / comment (from SetDefaults) are stamped automatically.
	r, err := sugar.BuyMarket(ctx, symbol, 0.01, nil, nil)
	if err != nil {
		// A rejected trade comes back as an error (e.g. market closed, not enough
		// money) — there is no "status code" to inspect separately.
		fmt.Println("  BuyMarket error:", err)
		return
	}
	t := r.GetTicket() // the ticket identifies this order for every later call
	fmt.Printf("  BuyMarket #%d @ %.5f\n", t, r.GetPrice())

	// 2) BRACKET. Set the protective stop/target as PIP DISTANCES from the entry.
	//    ModifySLTPByPips is direction-aware: for a BUY the stop goes below and the
	//    target above the open price (mirrored for a SELL). ptr(x) just takes the
	//    address of the literal, because the args are optional *float64.
	if err := sugar.ModifySLTPByPips(ctx, t, common.Ptr(20.0), common.Ptr(40.0)); err != nil {
		fmt.Println("  ModifySLTPByPips error:", err)
	}

	// 3) MANAGE. These two are one-shot primitives — in real code you'd call them
	//    on a timer. AutoBreakeven moves the stop to entry+2p once +10p in profit;
	//    TrailStop then follows price 15p behind. Both only ever tighten the stop,
	//    so calling them when nothing should move simply returns moved=false.
	if moved, err := sugar.AutoBreakeven(ctx, t, 10, 2); err == nil {
		fmt.Printf("  AutoBreakeven moved: %v\n", moved)
	}
	if moved, err := sugar.TrailStop(ctx, t, 15); err == nil {
		fmt.Printf("  TrailStop moved: %v\n", moved)
	}

	// 4) EXIT. ClosePartial trims part of the position; Close(…, nil) closes the
	//    rest in full. (On a 0.01 lot the partial may be a no-op below min volume.)
	if _, err := sugar.ClosePartial(ctx, t, 0.01); err == nil {
		fmt.Printf("  ClosePartial #%d ok\n", t)
	}
	if _, err := sugar.Close(ctx, t, nil); err == nil {
		fmt.Printf("  Close #%d ok\n", t)
	}

	// 5) BULK HELPERS. Handy for a clean slate — they iterate the open book for
	//    you. Both are no-ops (return 0) when nothing matches.
	if n, err := sugar.CloseAll(ctx, symbol); err == nil {
		fmt.Printf("  CloseAll(%s): %d closed\n", symbol, n) // close every market order on the symbol
	}
	if n, err := sugar.CancelPendings(ctx, ""); err == nil {
		fmt.Printf("  CancelPendings: %d deleted\n", n) // delete every pending (all symbols)
	}
}
