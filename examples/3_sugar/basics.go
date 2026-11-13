/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/3_sugar/basics.go - HIGH-LEVEL (MT4Sugar) BASICS DEMO
 PURPOSE:
   MT4Sugar groups 1-4 - connection/defaults, prices & quotes, symbol math and
   risk sizing. All read-only and safe to run.

 🎯 WHO SHOULD USE THIS:
   • Developers starting with the high-level MT4Sugar API
   • Anyone needing price math / risk sizing helpers (pips, lots, cash risk)
   • Learning the read side before trading

 📚 WHAT THIS DEMO COVERS (4 sections):

   1. CONNECTION & DEFAULTS
      • Ping()                               - cheap liveness check

   2. PRICES & QUOTES
      • EnsureSymbol() / LastQuote() / Quotes() / MidPrice() / SpreadPips() / Bars()

   3. SYMBOL MATH
      • Point() / Digits() / PipSize() / PipsToPrice() / PriceToPips()
      • NormalizeLot() / NormalizePrice() / TickValue()

   4. RISK SIZING
      • CalcCashRisk()                       - stop distance -> cash at risk
      • CalcLotByRisk()                      - cash budget -> normalized lot

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - raw gRPC calls
   MID-LEVEL  (MT4Service) - typed wrappers
   HIGH-LEVEL (MT4Sugar)   - THIS FILE: one-liner helpers with smart defaults

 ℹ️  SAFE: read-only - no orders are sent.

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 31 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

// Package sugar holds the MT4Sugar (high-level) examples, split by group.
package sugar

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// Basics covers MT4Sugar groups 1-4: connection/defaults, prices & quotes,
// symbol math, and risk sizing. These are all read-only and safe to run.
func Basics(ctx context.Context, sugar *mt4.MT4Sugar, symbol string) {
	fmt.Println("\n=== 🍬 MT4Sugar · Basics (prices · math · risk) ===")

	//#region 1. CONNECTION & DEFAULTS
	// ══════════════════════════════════════════════════════════════
	// 1. CONNECTION & DEFAULTS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// Ping()
	//   Cheap liveness check — true if an account snapshot succeeds.
	// ──────────────────────────────────────────────────────────────
	fmt.Printf("  Ping: %v\n", sugar.Ping(ctx))
	//#endregion 1. CONNECTION & DEFAULTS

	//#region 2. PRICES & QUOTES
	// ══════════════════════════════════════════════════════════════
	// 2. PRICES & QUOTES
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// EnsureSymbol()
	//   Add the symbol to Market Watch and wait for a usable two-sided
	//   quote. Always do this before quoting/trading a fresh symbol —
	//   demo feeds often send bid==ask or zeros for the first moment.
	// ──────────────────────────────────────────────────────────────
	if _, err := sugar.EnsureSymbol(ctx, symbol, 8*time.Second); err != nil {
		fmt.Println("  EnsureSymbol error:", err)
	}

	// ──────────────────────────────────────────────────────────────
	// LastQuote() / Quotes() / MidPrice() / SpreadPips() / Bars()
	//   Latest Bid/Ask (one and many symbols), mid price, spread in
	//   pips, and OHLC candles.
	// ──────────────────────────────────────────────────────────────
	if q, err := sugar.LastQuote(ctx, symbol); err == nil {
		fmt.Printf("  LastQuote %s: %.5f / %.5f\n", symbol, q.GetBid(), q.GetAsk())
	}
	if qs, err := sugar.Quotes(ctx, []string{"EURUSD", "GBPUSD"}); err == nil {
		fmt.Printf("  Quotes: %d\n", len(qs)) // several symbols in one round-trip
	}
	if mid, err := sugar.MidPrice(ctx, symbol); err == nil {
		fmt.Printf("  MidPrice: %.5f\n", mid) // (bid+ask)/2
	}
	if sp, err := sugar.SpreadPips(ctx, symbol); err == nil {
		fmt.Printf("  SpreadPips: %.1f\n", sp) // ask-bid expressed in pips
	}
	if b, err := sugar.Bars(ctx, symbol, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, 50); err == nil {
		fmt.Printf("  Bars (H1 x50): %d bars\n", len(b.GetHistoricalQuotes())) // OHLC candles
	}
	//#endregion 2. PRICES & QUOTES

	//#region 3. SYMBOL MATH
	// ══════════════════════════════════════════════════════════════
	// 3. SYMBOL MATH
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// Point() / Digits() / PipSize() / PipsToPrice() / PriceToPips()
	//   Read the symbol's contract spec. A "pip" is point×10 on 3/5-digit
	//   symbols, else point — PipSize hides that rule.
	// ──────────────────────────────────────────────────────────────
	point, _ := sugar.Point(ctx, symbol)
	digits, _ := sugar.Digits(ctx, symbol)
	pip, _ := sugar.PipSize(ctx, symbol)
	fmt.Printf("  Point %.5f  Digits %d  PipSize %.5f\n", point, digits, pip)
	if px, err := sugar.PipsToPrice(ctx, symbol, 20); err == nil {
		fmt.Printf("  PipsToPrice(20): %.5f\n", px) // 20 pips → a price distance
	}
	if pips, err := sugar.PriceToPips(ctx, symbol, 0.0035); err == nil {
		fmt.Printf("  PriceToPips(0.0035): %.1f\n", pips) // and back the other way
	}

	// ──────────────────────────────────────────────────────────────
	// NormalizeLot() / NormalizePrice() / TickValue()
	//   Round a lot to the broker's volume step (clamped to [min,max];
	//   skip it and a computed lot may be rejected with error 131),
	//   round a price to the symbol's digits, and the money per tick.
	// ──────────────────────────────────────────────────────────────
	if lot, err := sugar.NormalizeLot(ctx, symbol, 0.137); err == nil {
		// Round to the broker's volume step and clamp to [min,max] — always do
		// this before sending a computed lot, or it may be rejected (error 131).
		fmt.Printf("  NormalizeLot(0.137): %.2f\n", lot)
	}
	if np, err := sugar.NormalizePrice(ctx, symbol, 1.234567); err == nil {
		fmt.Printf("  NormalizePrice(1.234567): %.5f\n", np) // round to the symbol's digits
	}
	if tv, err := sugar.TickValue(ctx, symbol); err == nil {
		fmt.Printf("  TickValue: %.4f\n", tv) // money per tick per lot (account currency)
	}
	//#endregion 3. SYMBOL MATH

	//#region 4. RISK SIZING
	// ══════════════════════════════════════════════════════════════
	// 4. RISK SIZING
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// CalcCashRisk() / CalcLotByRisk()
	//   Turn a stop distance into cash at risk, and a cash budget into
	//   a normalized lot size (the two directions of position sizing).
	// ──────────────────────────────────────────────────────────────
	if risk, err := sugar.CalcCashRisk(ctx, symbol, 0.10, 20); err == nil {
		// "If a 0.10-lot trade is stopped out 20 pips away, I lose this much."
		fmt.Printf("  CalcCashRisk(0.10, 20p): %.2f\n", risk)
	}
	if lots, err := sugar.CalcLotByRisk(ctx, symbol, 50, 20); err == nil {
		// The inverse: "size a trade that risks 50 over a 20-pip stop" (normalized).
		fmt.Printf("  CalcLotByRisk(50 cash, 20p): %.2f lot\n", lots)
	}
	//#endregion 4. RISK SIZING
}
