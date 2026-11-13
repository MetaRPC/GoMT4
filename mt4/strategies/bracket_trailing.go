/*══════════════════════════════════════════════════════════════════════════════
 FILE: mt4/strategies/bracket_trailing.go - ORCHESTRATOR: Bracket + Trailing
 PURPOSE:
   Enters with a pip bracket, then starts trailing only once the trade is
   activateAtPips in profit - letting the position breathe before the stop
   follows it. Auto-breakeven (from the preset) still applies earlier.
   Returns the ticket.

 📚 FLOW:
   1. sizeLot()          - fixed or risk-based volume
   2. BuyMarket()        - open the position at market
   3. ModifySLTPByPips() - attach the pip SL/TP bracket
   4. manage()           - breakeven early; trailing arms at +activateAtPips
   5. Close()            - flatten

 ⚠️  Places a REAL market order - DEMO accounts only.
══════════════════════════════════════════════════════════════════════════════*/

package strategies

import (
	"context"
	"fmt"

	"github.com/MetaRPC/GoMT4/mt4"
)

// RunBracketTrailing enters with a bracket, then starts trailing only once the
// trade is activateAtPips in profit — this lets the position breathe before the
// stop starts following it. Auto-breakeven (from the preset) still applies
// earlier. Returns the ticket.
//
// Example — trailing arms once +12 pips in profit:
//
//	ticket, err := strategies.RunBracketTrailing(ctx, sugar,
//	    strategies.MarketEURUSD(), strategies.Balanced(), 12)
func RunBracketTrailing(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset, activateAtPips float64) (int32, error) {
	fmt.Printf("── Bracket + Trailing · %s ──\n", sp.Symbol)

	// ── sizeLot() + BuyMarket(): risk-size the lot, then open at market ──
	lots, err := sizeLot(ctx, sugar, sp, rp)
	if err != nil {
		return 0, err
	}
	r, err := sugar.BuyMarket(ctx, sp.Symbol, lots, nil, nil)
	if err != nil {
		return 0, err
	}
	ticket := r.GetTicket()
	// ── ModifySLTPByPips(): attach the pip SL/TP bracket ──
	if err := sugar.ModifySLTPByPips(ctx, ticket, &rp.SLPips, rp.TPPips); err != nil {
		return ticket, err
	}
	fmt.Printf("entry : BUY #%d @ %.5f; trailing starts at +%.0fp\n", ticket, r.GetPrice(), activateAtPips)

	if err := manage(ctx, sugar, ticket, sp.Symbol, rp, 10, activateAtPips); err != nil {
		return ticket, err
	}

	fmt.Printf("close : #%d\n", ticket)
	_, err = sugar.Close(ctx, ticket, nil)
	return ticket, err
}
