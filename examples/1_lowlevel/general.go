/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/1_lowlevel/general.go - LOW-LEVEL MT4 GENERAL (READ-ONLY) DEMO
 PURPOSE:
   Exercises the read-only MT4Account methods - account info (summary, getters,
   trade stats) and market/symbol info. No trades here.

 🎯 WHO SHOULD USE THIS:
   • Developers learning the low-level MT4 gRPC API (read side)
   • Building monitoring / dashboard / reporting tools
   • Seeing exactly what each MT4Account method returns

 📚 WHAT THIS DEMO COVERS (2 sections):

   1. ACCOUNT INFO
      • AccountSummary()                     - login / balance / equity / leverage
      • AccountBalance/Equity/Credit/Name/Company/Currency() - individual getters
      • TradeStats() / TradeStatsEquityHistory() / EquityHistory() / TradeUnrealizedPL()

   2. MARKET & SYMBOLS
      • SymbolSelect() / Quote() / QuoteMany()
      • Symbols() / SymbolParams() / SymbolParamsMany() / TickValueWithSize()
      • QuoteHistory() / ServerTimezone()

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - THIS FILE: raw gRPC calls
   MID-LEVEL  (MT4Service) - typed wrappers, less boilerplate
   HIGH-LEVEL (MT4Sugar)   - one-liner helpers with smart defaults

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 11 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

// Package lowlevel holds the MT4Account (low-level gRPC) examples, split by group.
package lowlevel

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// General exercises the read-only MT4Account methods: account info (summary +
// getters + trade stats) and market/symbol info.
func General(ctx context.Context, acc *mt4.MT4Account, symbol string) {
	fmt.Println("\n=== 📂 MT4Account · General (Account Info + Market) ===")

	//#region 1. ACCOUNT INFO
	// ══════════════════════════════════════════════════════════════
	// 1. ACCOUNT INFO
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// AccountSummary()
	//   Login / balance / equity / currency / leverage in one call.
	// ──────────────────────────────────────────────────────────────
	if a, err := acc.AccountSummary(ctx); err != nil {
		fmt.Println("  AccountSummary error:", err)
	} else {
		fmt.Printf("  Summary: #%d  balance %.2f %s  equity %.2f  leverage 1:%d\n",
			a.GetAccountLogin(), a.GetAccountBalance(), a.GetAccountCurrency(), a.GetAccountEquity(), a.GetAccountLeverage())
	}

	// ──────────────────────────────────────────────────────────────
	// AccountBalance/Equity/Credit/Name/Company/Currency()
	//   Individual property getters (each a single RPC).
	// ──────────────────────────────────────────────────────────────
	if v, err := acc.AccountBalance(ctx); err == nil {
		fmt.Printf("  AccountBalance: %.2f\n", v)
	}
	if v, err := acc.AccountEquity(ctx); err == nil {
		fmt.Printf("  AccountEquity:  %.2f\n", v)
	}
	if v, err := acc.AccountCredit(ctx); err == nil {
		fmt.Printf("  AccountCredit:  %.2f\n", v)
	}
	if v, err := acc.AccountName(ctx); err == nil {
		fmt.Printf("  AccountName:    %s\n", v)
	}
	if v, err := acc.AccountCompany(ctx); err == nil {
		fmt.Printf("  AccountCompany: %s\n", v)
	}
	if v, err := acc.AccountCurrency(ctx); err == nil {
		fmt.Printf("  AccountCurrency:%s\n", v)
	}

	// ──────────────────────────────────────────────────────────────
	// TradeStats() / TradeStatsEquityHistory() / EquityHistory() / TradeUnrealizedPL()
	//   Aggregated analytics: trade stats, the equity curve and
	//   per-position unrealized P/L.
	// ──────────────────────────────────────────────────────────────
	if st, err := acc.TradeStats(ctx, nil, nil); err != nil {
		fmt.Println("  TradeStats error:", err)
	} else {
		fmt.Printf("  TradeStats: %d trades, profit factor %.2f\n", st.GetTrades(), st.GetProfitFactor())
	}
	if _, err := acc.TradeStatsEquityHistory(ctx, nil, nil, nil); err == nil {
		fmt.Println("  TradeStatsEquityHistory: ok")
	}
	if eq, err := acc.EquityHistory(ctx, nil, nil, nil); err == nil {
		fmt.Printf("  EquityHistory: %d points\n", len(eq.GetPoints()))
	}
	if _, err := acc.TradeUnrealizedPL(ctx, nil, nil); err == nil {
		fmt.Println("  TradeUnrealizedPL: ok")
	}
	//#endregion 1. ACCOUNT INFO

	//#region 2. MARKET & SYMBOLS
	// ══════════════════════════════════════════════════════════════
	// 2. MARKET & SYMBOLS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// SymbolSelect() / Quote() / QuoteMany()
	//   Add a symbol to Market Watch, then read its Bid/Ask - one
	//   symbol and several at once.
	// ──────────────────────────────────────────────────────────────
	if ok, err := acc.SymbolSelect(ctx, symbol, true); err == nil {
		fmt.Printf("  SymbolSelect %s: %v\n", symbol, ok)
	}
	if q, err := acc.Quote(ctx, symbol); err == nil {
		fmt.Printf("  Quote %s: bid %.5f ask %.5f\n", symbol, q.GetBid(), q.GetAsk())
	}
	if d, err := acc.QuoteMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"}); err == nil {
		fmt.Printf("  QuoteMany: %d quotes\n", len(d.GetQuotes()))
	}

	// ──────────────────────────────────────────────────────────────
	// Symbols() / SymbolParams() / SymbolParamsMany() / TickValueWithSize()
	//   The symbol catalog and per-symbol trading params / tick value.
	// ──────────────────────────────────────────────────────────────
	if d, err := acc.Symbols(ctx); err == nil {
		fmt.Printf("  Symbols: %d\n", len(d.GetSymbolNameInfos()))
	}
	if si, err := acc.SymbolParams(ctx, symbol); err == nil {
		fmt.Printf("  SymbolParams %s: digits %d, point %.5f, step %.2f\n", symbol, si.GetDigits(), si.GetPoint(), si.GetVolumeStep())
	}
	if d, err := acc.SymbolParamsMany(ctx, []string{symbol}); err == nil {
		fmt.Printf("  SymbolParamsMany: %d infos\n", len(d.GetSymbolInfos()))
	}
	if d, err := acc.TickValueWithSize(ctx, []string{"EURUSD", "GBPUSD"}); err == nil {
		fmt.Printf("  TickValueWithSize: %d infos\n", len(d.GetInfos()))
	}

	// ──────────────────────────────────────────────────────────────
	// QuoteHistory() / ServerTimezone()
	//   H1 history bars over the last week and the trade-server timezone.
	// ──────────────────────────────────────────────────────────────
	now := time.Now()
	if h, err := acc.QuoteHistory(ctx, symbol, pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, now.Add(-7*24*time.Hour), now); err == nil {
		fmt.Printf("  QuoteHistory (H1, 7d): %d bars\n", len(h.GetHistoricalQuotes()))
	}
	if tz, err := acc.ServerTimezone(ctx); err == nil {
		fmt.Printf("  ServerTimezone: %s UTC%+d\n", tz.GetServer(), tz.GetUtcOffsetSeconds()/3600)
	}
	//#endregion 2. MARKET & SYMBOLS
}
