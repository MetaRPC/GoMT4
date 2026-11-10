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

	lots, err := sizeLot(ctx, sugar, sp, rp)
	if err != nil {
		return 0, err
	}
	r, err := sugar.BuyMarket(ctx, sp.Symbol, lots, nil, nil)
	if err != nil {
		return 0, err
	}
	ticket := r.GetTicket()
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
