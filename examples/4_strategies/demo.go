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

	if !common.TradingEnabled {
		fmt.Println("  (orchestrators skipped — they send real orders; set MT4_ENABLE_TRADING=1)")
		return
	}

	// An orchestrator takes two presets: a StrategyPreset (WHAT to trade) and a
	// RiskPreset (HOW MUCH to risk + how to protect). MarketEURUSD() gives a market
	// strategy with magic 777; we point it at the configured symbol.
	sp := strategies.MarketEURUSD()
	sp.Symbol = symbol

	// Market One-Shot: size by risk → market BUY → pip bracket → manage → close.
	// Balanced() = 1% risk, SL 20 / TP 40, trail 15, breakeven 10/+2.
	fmt.Println("\n  ▶ Market One-Shot (balanced)")
	if t, err := strategies.RunMarketOneShot(ctx, sugar, sp, strategies.Balanced()); err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Printf("  done, ticket %d\n", t)
	}

	// Bracket + Trailing: same entry, but the trailing stop only arms once the
	// trade is +12 pips in profit (so it isn't stopped out by normal noise first).
	fmt.Println("\n  ▶ Bracket + Trailing (balanced, trail arms at +12p)")
	if t, err := strategies.RunBracketTrailing(ctx, sugar, sp, strategies.Balanced(), 12); err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Printf("  done, ticket %d\n", t)
	}

	// Pending Bracket: place a BUY LIMIT, wait up to 30s for the fill (internally
	// via WaitFilled), cancel if unfilled, otherwise bracket and manage.
	fmt.Println("\n  ▶ Pending Bracket (conservative, 30s fill window)")
	if t, err := strategies.RunPendingBracket(ctx, sugar, sp, strategies.Conservative(), 30); err != nil {
		fmt.Println("  error:", err)
	} else {
		fmt.Printf("  done, ticket %d\n", t)
	}
}
