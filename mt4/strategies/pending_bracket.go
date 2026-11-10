package strategies

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/mt4"
)

// RunPendingBracket places a BUY LIMIT (at the preset's EntryPrice, or 5 pips
// below market by default), waits up to timeoutSeconds for the fill, and cancels
// the pending if it doesn't fill. Once filled it brackets and manages. Returns
// the ticket (filled-and-managed, or the cancelled pending).
//
// Example — explicit limit price, 60s fill window:
//
//	ticket, err := strategies.RunPendingBracket(ctx, sugar,
//	    strategies.LimitEURUSD(1.0750), strategies.Conservative(), 60)
func RunPendingBracket(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset, timeoutSeconds int) (int32, error) {
	fmt.Printf("── Pending Bracket · %s ──\n", sp.Symbol)
	pip, err := sugar.PipSize(ctx, sp.Symbol)
	if err != nil {
		return 0, err
	}

	// entry price: explicit, or 5 pips below market for the demo
	var entry float64
	if sp.EntryPrice != nil {
		entry = *sp.EntryPrice
	} else {
		q, err := sugar.LastQuote(ctx, sp.Symbol)
		if err != nil {
			return 0, err
		}
		entry = q.GetBid() - 5.0*pip
	}

	lots, err := sizeLot(ctx, sugar, sp, rp)
	if err != nil {
		return 0, err
	}
	p, err := sugar.BuyLimit(ctx, sp.Symbol, lots, entry, nil, nil)
	if err != nil {
		return 0, err
	}
	ticket := p.GetTicket()
	fmt.Printf("placed : BUY LIMIT #%d %.2f lot @ %.5f\n", ticket, lots, entry)

	// wait for the fill, or cancel on timeout
	filled, err := sugar.WaitFilled(ctx, ticket, time.Duration(timeoutSeconds)*time.Second)
	if err != nil {
		return ticket, err
	}
	if !filled {
		fmt.Printf("timeout: not filled in %ds — cancelling\n", timeoutSeconds)
		_, _ = sugar.Close(ctx, ticket, nil) // delete the resting pending order
		return ticket, nil
	}
	fmt.Printf("filled : #%d\n", ticket)

	if err := sugar.ModifySLTPByPips(ctx, ticket, &rp.SLPips, rp.TPPips); err != nil {
		return ticket, err
	}
	trailFrom := 0.0
	if rp.BETriggerPips != nil {
		trailFrom = *rp.BETriggerPips
	}
	if err := manage(ctx, sugar, ticket, sp.Symbol, rp, 5, trailFrom); err != nil {
		return ticket, err
	}

	_, err = sugar.Close(ctx, ticket, nil)
	fmt.Println("close  : flat")
	return ticket, err
}
