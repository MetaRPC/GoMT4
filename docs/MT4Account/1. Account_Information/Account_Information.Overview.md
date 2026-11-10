# 1. Account Information — Overview

> Account snapshot, balance/equity getters, and aggregated trade statistics. All read-only.

These live on the `mt4_term_api.AccountHelper` and `mt4_term_api.TradeStats` gRPC services. The one call you'll use most is [`AccountSummary`](./AccountSummary.md); the statistics calls back a performance report.

---

## Methods in this group

| Method | Returns | One-liner |
|--------|---------|-----------|
| [`AccountSummary(ctx)`](./AccountSummary.md) | `*pb.AccountSummaryData` | Full account snapshot (balance, equity, margin, leverage, permissions). |
| [`TradeStats(ctx, from, fromUtc)`](./TradeStats.md) | `*pb.StatsData` | Aggregated performance statistics. |
| [`TradeStatsEquityHistory(ctx, from, fromUtc, excludeSameBars)`](./TradeStatsEquityHistory.md) | `*pb.StatsWithChartsData` | Stats + equity-curve points + trades. |
| [`EquityHistory(ctx, from, fromUtc, excludeSameBars)`](./EquityHistory.md) | `*pb.EquityHistoryData` | Just the equity-curve points. |
| [`TradeUnrealizedPL(ctx, from, fromUtc)`](./TradeUnrealizedPL.md) | `*pb.TradeUnrealizedPLData` | Per-open-order unrealized-P/L extrema. |

---

## Convenience getters (wrap `AccountSummary`)

These read `AccountSummary` under the hood and return a single field — handy one-liners, one round-trip each:

| Method | Returns | Field |
|--------|---------|-------|
| `AccountLogin(ctx)` | `int64` | `GetAccountLogin()` |
| `AccountBalance(ctx)` | `float64` | `GetAccountBalance()` |
| `AccountCredit(ctx)` | `float64` | `GetAccountCredit()` |
| `AccountEquity(ctx)` | `float64` | `GetAccountEquity()` |
| `AccountLeverage(ctx)` | `int64` | `GetAccountLeverage()` |
| `AccountName(ctx)` | `string` | `GetAccountUserName()` |
| `AccountCompany(ctx)` | `string` | `GetAccountCompanyName()` |
| `AccountCurrency(ctx)` | `string` | `GetAccountCurrency()` |

> If you need several of these, call `AccountSummary` once and read the getters off the result instead of making one round-trip per field.

---

**Back to:** [MT4Account — Overview](../Overview.md)
