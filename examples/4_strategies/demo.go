/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/4_strategies/demo.go - STRATEGIES (ORCHESTRATORS + PRESETS) DEMO
 PURPOSE:
   Runs the three ready-made orchestrators from mt4/strategies with presets:
   a full size -> entry -> bracket -> manage -> close lifecycle each.

 🎯 WHO SHOULD USE THIS:
   • Developers wanting a complete, automated trade lifecycle to copy
   • Seeing how StrategyPreset (what) + RiskPreset (how much) compose
   • Learning the orchestrator entry points before writing your own

 📚 WHAT THIS DEMO COVERS:

   PRESETS
      • StrategyPreset: MarketEURUSD() / LimitEURUSD()
      • RiskPreset:     Conservative() / Balanced() / Aggressive()

   ORCHESTRATORS
      • RunMarketOneShot()   - size by risk -> market BUY -> bracket -> manage -> close
      • RunBracketTrailing() - same, trailing stop arms once +N pips in profit
      • RunPendingBracket()  - BUY LIMIT, wait for fill, cancel if unfilled, else bracket

 🔄 API LEVELS:
   Orchestrators sit ABOVE MT4Sugar (which wraps MT4Service / MT4Account).

 ⚠️  WARNING: these SEND REAL ORDERS on a demo account - gated behind
     MT4_ENABLE_TRADING=1 (skipped otherwise).

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 41 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

// Package strategydemo runs the ready-made orchestrators from mt4/strategies.
package strategydemo

import (
	"context"
	"fmt"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
	"github.com/MetaRPC/GoMT4/mt4/strategies"
)

// Demo runs the three orchestrators with presets. These SEND REAL ORDERS, so
// they are gated behind MT4_ENABLE_TRADING=1 (demo accounts only).
func Demo(ctx context.Context, sugar *mt4.MT4Sugar, symbol string) {
	fmt.Println("\n=== 🧭 Strategies · orchestrators + presets ===")
	fmt.Println("  Presets: MarketEURUSD / LimitEURUSD · Conservative / Balanced / Aggressive")
	fmt.Println("  Orchestrators: RunMarketOneShot · RunBracketTrailing · RunPendingBracket")

	//#region GUARD & PRESETS
	// ══════════════════════════════════════════════════════════════
	// GUARD & PRESETS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// Trading guard
	//   The orchestrators place real orders, so bail out unless
	//   MT4_ENABLE_TRADING=1 is set.
	// ──────────────────────────────────────────────────────────────
	if !common.TradingEnabled {
		fmt.Println("  (orchestrators skipped — they send real orders; set MT4_ENABLE_TRADING=1)")
		return
	}

	// ──────────────────────────────────────────────────────────────
	// Presets
	//   An orchestrator takes a StrategyPreset (WHAT to trade) and a
	//   RiskPreset (HOW MUCH to risk + how to protect). MarketEURUSD()
	//   gives a market strategy with magic 777; point it at the symbol.
	// ──────────────────────────────────────────────────────────────
	sp := strategies.MarketEURUSD()
	sp.Symbol = symbol
	//#endregion GUARD & PRESETS

	//#region ORCHESTRATORS
	// ══════════════════════════════════════════════════════════════
	// ORCHESTRATORS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// RunMarketOneShot()
	//   Size by risk → market BUY → pip bracket → manage → close.
	//   Balanced() = 1% risk, SL 20 / TP 40, trail 15, breakeven 10/+2.
	// ──────────────────────────────────────────────────────────────
	fmt.Println("\n  ▶ Market One-Shot (balanced)")
	if t, err := strategies.RunMarketOneShot(ctx, sugar, sp, strategies.Balanced()); err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Printf("  done, ticket %d\n", t)
	}

	// ──────────────────────────────────────────────────────────────
	// RunBracketTrailing()
	//   Same entry, but the trailing stop only arms once the trade is
	//   +12 pips in profit (so noise doesn't stop it out first).
	// ──────────────────────────────────────────────────────────────
	fmt.Println("\n  ▶ Bracket + Trailing (balanced, trail arms at +12p)")
	if t, err := strategies.RunBracketTrailing(ctx, sugar, sp, strategies.Balanced(), 12); err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Printf("  done, ticket %d\n", t)
	}

	// ──────────────────────────────────────────────────────────────
	// RunPendingBracket()
	//   Place a BUY LIMIT, wait up to 30s for the fill (via WaitFilled),
	//   cancel if unfilled, otherwise bracket and manage.
	// ──────────────────────────────────────────────────────────────
	fmt.Println("\n  ▶ Pending Bracket (conservative, 30s fill window)")
	if t, err := strategies.RunPendingBracket(ctx, sugar, sp, strategies.Conservative(), 30); err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Printf("  done, ticket %d\n", t)
	}
	//#endregion ORCHESTRATORS
}
