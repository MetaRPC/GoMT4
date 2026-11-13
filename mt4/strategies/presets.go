/*══════════════════════════════════════════════════════════════════════════════
 FILE: mt4/strategies/presets.go - STRATEGY & RISK PRESETS + SHARED HELPERS
 PURPOSE:
   The building blocks the orchestrators are made of: a StrategyPreset (WHAT to
   trade), a RiskPreset (HOW MUCH to risk and how to protect), a small catalog of
   ready-made presets, and the shared sizeLot / manage helpers.

 📚 WHAT THIS FILE PROVIDES (3 sections):

   1. TYPES
      • StrategyPreset - symbol, market/limit, entry price, lots, magic, comment
      • RiskPreset     - risk %, SL/TP pips, trailing, breakeven

   2. PRESET CATALOG
      • MarketEURUSD() / LimitEURUSD(price)          - WHAT presets
      • Conservative() / Balanced() / Aggressive()   - HOW (risk) presets

   3. SHARED HELPERS (used by the orchestrators)
      • sizeLot() - fixed lot, or a risk-based lot over the stop distance
      • manage()  - auto-breakeven + gated trailing for a few seconds

 ⚠️  DEMO accounts only - the orchestrators built on these place real orders.
══════════════════════════════════════════════════════════════════════════════*/

// Package strategies provides ready-made orchestrators and presets (demonstration
// templates — DEMO accounts only), driven by MT4Sugar.
//
// A StrategyPreset says WHAT to trade; a RiskPreset says HOW MUCH and how to protect.
package strategies

import (
	"context"
	"time"

	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

//#region 1. TYPES
// ══════════════════════════════════════════════════════════════
// 1. TYPES
// ══════════════════════════════════════════════════════════════

// StrategyPreset — WHAT to trade. Use the catalog (MarketEURUSD, LimitEURUSD) or
// build one by hand:
//
//	sp := strategies.StrategyPreset{Symbol: "GBPUSD", UseMarket: true, Magic: 101}
//	lot := 0.10
//	sp.Lots = &lot // fixed volume instead of risk-based sizing
type StrategyPreset struct {
	Symbol     string
	UseMarket  bool
	EntryPrice *float64 // limit price when !UseMarket
	Lots       *float64 // fixed volume; nil = size by risk
	Magic      int32
	Comment    string
}

// RiskPreset — HOW MUCH to risk and how to protect. Fields left nil/zero are
// simply not applied (no TP, no trailing, no breakeven).
//
//	tp := 40.0
//	rp := strategies.RiskPreset{RiskPercent: 1.0, SLPips: 20, TPPips: &tp}
type RiskPreset struct {
	RiskPercent   float64
	SLPips        float64
	TPPips        *float64
	TrailingPips  *float64
	BETriggerPips *float64
	BEPlusPips    float64
}

//#endregion 1. TYPES

//#region 2. PRESET CATALOG
// ══════════════════════════════════════════════════════════════
// 2. PRESET CATALOG
// ══════════════════════════════════════════════════════════════

// MarketEURUSD: market EURUSD (magic 777).
func MarketEURUSD() StrategyPreset {
	return StrategyPreset{Symbol: "EURUSD", UseMarket: true, Magic: 777, Comment: "market EURUSD"}
}

// LimitEURUSD: BUY LIMIT EURUSD at price (magic 778).
func LimitEURUSD(price float64) StrategyPreset {
	return StrategyPreset{Symbol: "EURUSD", UseMarket: false, EntryPrice: &price, Magic: 778, Comment: "limit EURUSD"}
}

// Conservative: 0.5% risk, SL 25, TP 50.
func Conservative() RiskPreset {
	tp := 50.0
	return RiskPreset{RiskPercent: 0.5, SLPips: 25, TPPips: &tp}
}

// Balanced: 1.0% risk, SL 20, TP 40, trail 15, breakeven 10/+2.
func Balanced() RiskPreset {
	tp, tr, be := 40.0, 15.0, 10.0
	return RiskPreset{RiskPercent: 1.0, SLPips: 20, TPPips: &tp, TrailingPips: &tr, BETriggerPips: &be, BEPlusPips: 2}
}

// Aggressive: 2.0% risk, SL 15, TP 45, trail 10, breakeven 8/+1.
func Aggressive() RiskPreset {
	tp, tr, be := 45.0, 10.0, 8.0
	return RiskPreset{RiskPercent: 2.0, SLPips: 15, TPPips: &tp, TrailingPips: &tr, BETriggerPips: &be, BEPlusPips: 1}
}

//#endregion 2. PRESET CATALOG

//#region 3. SHARED HELPERS
// ══════════════════════════════════════════════════════════════
// 3. SHARED HELPERS (used by the orchestrators)
// ══════════════════════════════════════════════════════════════

// sizeLot returns the fixed lot from the strategy, or a risk-based lot over the stop.
func sizeLot(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset) (float64, error) {
	if sp.Lots != nil {
		return *sp.Lots, nil
	}
	a, err := sugar.GetAccountInfo(ctx)
	if err != nil {
		return 0, err
	}
	return sugar.CalcLotByRisk(ctx, sp.Symbol, a.GetAccountBalance()*rp.RiskPercent/100.0, rp.SLPips)
}

// manage runs auto-breakeven and (gated) trailing on an open BUY order for
// manageSecs. Returns early if the order is no longer open (hit SL/TP).
func manage(ctx context.Context, sugar *mt4.MT4Sugar, ticket int32, symbol string, rp RiskPreset, manageSecs int, trailActivatePips float64) error {
	var pip float64
	if rp.TrailingPips != nil {
		pip, _ = sugar.PipSize(ctx, symbol)
	}
	deadline := time.Now().Add(time.Duration(manageSecs) * time.Second)
	for time.Now().Before(deadline) {
		o, err := sugar.FindTicket(ctx, ticket)
		if err != nil {
			return err
		}
		if o == nil {
			break // closed by SL/TP
		}
		if rp.BETriggerPips != nil {
			_, _ = sugar.AutoBreakeven(ctx, ticket, *rp.BETriggerPips, rp.BEPlusPips)
		}
		if rp.TrailingPips != nil && pip > 0 {
			if o2, _ := sugar.FindTicket(ctx, ticket); o2 != nil {
				if q, err := sugar.LastQuote(ctx, symbol); err == nil {
					var profitPips float64
					if o2.GetOrderType() == pb.OpenedOrderType_OO_OP_BUY {
						profitPips = (q.GetBid() - o2.GetOpenPrice()) / pip
					} else {
						profitPips = (o2.GetOpenPrice() - q.GetAsk()) / pip
					}
					if profitPips >= trailActivatePips {
						_, _ = sugar.TrailStop(ctx, ticket, *rp.TrailingPips)
					}
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

//#endregion 3. SHARED HELPERS
