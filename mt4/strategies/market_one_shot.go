/*══════════════════════════════════════════════════════════════════════════════
 FILE: mt4/strategies/market_one_shot.go - ORCHESTRATOR: Market One-Shot
 PURPOSE:
   The simplest orchestrator: size → market BUY → pip bracket → manage briefly →
   close. Returns the ticket.

 📚 FLOW:
   1. sizeLot()          - fixed or risk-based volume
   2. BuyMarket()        - open the position at market
   3. ModifySLTPByPips() - attach the pip SL/TP bracket
   4. manage()           - breakeven + (gated) trailing for a few seconds
   5. Close()            - flatten

 ⚠️  Places a REAL market order - DEMO accounts only.
══════════════════════════════════════════════════════════════════════════════*/

package strategies

import (
	"context"
	"fmt"

	"github.com/MetaRPC/GoMT4/mt4"
)

// RunMarketOneShot is the simplest orchestrator: size → market BUY → pip bracket →
// manage briefly → close. Returns the ticket.
//
// Example — run it with ready-made presets (DEMO accounts only):
//
//	sugar := mt4.NewMT4Sugar(mt4.NewMT4Service(acc)).SetDefaults(777, 20, "demo")
//	ticket, err := strategies.RunMarketOneShot(ctx, sugar,
//	    strategies.MarketEURUSD(), // WHAT: market EURUSD, magic 777
//	    strategies.Balanced())     // HOW: 1% risk, SL 20 / TP 40, trail 15, BE 10/+2
func RunMarketOneShot(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset) (int32, error) {
	fmt.Printf("── Market One-Shot · %s ──\n", sp.Symbol)

	// ── sizeLot(): fixed lot from the preset, or a risk-based lot over the stop ──
	lots, err := sizeLot(ctx, sugar, sp, rp)
	if err != nil {
		return 0, err
	}
	// ── BuyMarket(): open the position at market ──
	r, err := sugar.BuyMarket(ctx, sp.Symbol, lots, nil, nil)
	if err != nil {
		return 0, err
	}
	ticket := r.GetTicket()
	fmt.Printf("entry : BUY #%d %.2f lot @ %.5f\n", ticket, lots, r.GetPrice())

	// pip-based SL/TP bracket
	if err := sugar.ModifySLTPByPips(ctx, ticket, &rp.SLPips, rp.TPPips); err != nil {
		return ticket, err
	}

	// manage a few seconds: breakeven, then trail (both from the preset)
	trailFrom := 0.0
	if rp.BETriggerPips != nil {
		trailFrom = *rp.BETriggerPips
	}
	if err := manage(ctx, sugar, ticket, sp.Symbol, rp, 5, trailFrom); err != nil {
		return ticket, err
	}

	fmt.Printf("close : #%d\n", ticket)
	_, err = sugar.Close(ctx, ticket, nil)
	return ticket, err
}
