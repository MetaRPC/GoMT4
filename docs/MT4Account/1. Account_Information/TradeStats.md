# TradeStats

> Aggregated trading performance — profit factor, win rates, drawdowns, expectancy and more, in one call.

**API Information:**

* **Low-level API:** `MT4Account.TradeStats(ctx, from, fromUtc)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **RPC:** `TradeStats(...) → ...Reply`
* **Payload:** `StatsData`

## 💬 Just the essentials

* **What it is.** A performance report computed from the account's closed trades.
* **Why you need it.** KPIs (profit factor, expectancy, drawdown, win %) without crunching history yourself.
* **When to use.** Dashboards, periodic reports, strategy health checks.

---

## 🎯 Purpose

* Profitability KPIs: profit factor, GHPR, expectancy, Sharpe
* Win rates for longs and shorts
* Balance / equity drawdown (raw and relative)
* Totals: trades, lots, pips, deposits/withdrawals
* Best / worst trade

---

## Method Signature

```go
func (a *MT4Account) TradeStats(ctx context.Context, from, fromUtc *string) (*pb.StatsData, error)
```

**Request message:**

```protobuf
message TradeStatsRequest {
  optional string from = 1;       // server-time "yyyy-MM-ddTHH:mm:ss"; nil = since account start
  optional string from_utc = 2;   // same moment in UTC — send one of from / from_utc
}
```

**Payload message (excerpt):**

```protobuf
message StatsData {
  int64  trades = ...;
  double profit_factor = ...;
  double sharpe_ratio = ...;
  double growth_percent = ...;
  double max_balance_drawdown_relative = ...;
  double max_equity_drawdown_relative = ...;
  ExpectancyData expectancy = ...;
  WonData longs_won = ...;
  WonData shorts_won = ...;
  ProfitData best_trade = ...;
  ProfitData worst_trade = ...;
  repeated MarketTradeCountData markets = ...;
  // … and more (see getters below)
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `from` | `*string` | no | Server-time start `"yyyy-MM-ddTHH:mm:ss"`; `nil` = since account start. |
| `fromUtc` | `*string` | no | Same start in UTC. Send one of `from` / `fromUtc`, not both. |

---

## ⬆️ Output — `*pb.StatsData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `trades` | int64 | `GetTrades() int64` | Number of closed trades |
| `profit_factor` | double | `GetProfitFactor() float64` | Gross profit / gross loss |
| `sharpe_ratio` | double | `GetSharpeRatio() float64` | Risk-adjusted return |
| `growth_percent` | double | `GetGrowthPercent() float64` | Account growth % |
| `ghpr` | double | `GetGhpr() float64` | Geometric holding-period return |
| `pips` / `lots` | double | `GetPips()` / `GetLots()` | Totals |
| `max_balance_drawdown_relative` | double | `GetMaxBalanceDrawdownRelative() float64` | Worst balance drawdown (%) |
| `max_equity_drawdown_relative` | double | `GetMaxEquityDrawdownRelative() float64` | Worst equity drawdown (%) |
| `expectancy` | ExpectancyData | `GetExpectancy()` | → `GetDollar()`, `GetPips()` |
| `longs_won` / `shorts_won` | WonData | `GetLongsWon()` / `GetShortsWon()` | → `GetAll()`, `GetWonCount()`, `GetWonPersent()` |
| `best_trade` / `worst_trade` | ProfitData | `GetBestTrade()` / `GetWorstTrade()` | → `GetProfit()`, `GetDate()` |
| `markets` | repeated MarketTradeCountData | `GetMarkets()` | Per-symbol trade counts |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **`from` is server time**, formatted `"yyyy-MM-ddTHH:mm:ss"`; use `fromUtc` to pass UTC instead.
* **Nested getters are nillable.** `GetExpectancy()` / `GetBestTrade()` can be `nil` on an empty history — guard before dereferencing.
* **Note the spelling** — `WonData.GetWonPersent()` and `ProfitData.GetTiket()` are as generated.

---

## 🔗 Usage Examples

### 1) Performance summary

```go
func perfSummary(ctx context.Context, acc *mt4.MT4Account) error {
    st, err := acc.TradeStats(ctx, nil, nil)
    if err != nil { return err }
    fmt.Printf("trades=%d  PF=%.2f  growth=%.1f%%  maxDD=%.1f%%\n",
        st.GetTrades(), st.GetProfitFactor(), st.GetGrowthPercent(),
        st.GetMaxEquityDrawdownRelative())
    if e := st.GetExpectancy(); e != nil {
        fmt.Printf("expectancy: %.2f$ / %.1f pips per trade\n", e.GetDollar(), e.GetPips())
    }
    return nil
}
```

### 2) Which markets get traded most?

```go
func topMarkets(ctx context.Context, acc *mt4.MT4Account) error {
    st, err := acc.TradeStats(ctx, nil, nil)
    if err != nil { return err }
    for _, m := range st.GetMarkets() {
        fmt.Printf("%s: %d trades  win %.0f%%  profit %+.2f\n",
            m.GetMarketName(), m.GetCount(), m.GetWinRate(), m.GetProfit())
    }
    return nil
}
```

### 3) Month-to-date health check

```go
func monthToDate(ctx context.Context, acc *mt4.MT4Account) (bool, error) {
    from := time.Now().Format("2006-01") + "-01T00:00:00"
    st, err := acc.TradeStats(ctx, &from, nil)
    if err != nil { return false, err }
    healthy := st.GetProfitFactor() >= 1.0 && st.GetMaxEquityDrawdownRelative() < 20.0
    return healthy, nil
}
```

---

## ⚠️ Pitfalls

* **Empty history.** With no closed trades the ratios are 0 and nested messages may be `nil` — guard.
* **Two time axes.** Don't send both `from` and `fromUtc`.
* **Server vs UTC.** A naive `from` is interpreted as **server** time; mind the broker offset.

---

## 🧪 Testing Suggestions

* **Happy path (demo with history).** `GetTrades() > 0`; `GetProfitFactor()` finite; nested getters non-nil.
* **Windowed.** Pass a recent `from` and confirm the counts shrink vs the all-time call.
* **Empty window.** A future `from` yields `GetTrades() == 0` and no panic on nested getters.

---

## 📚 See also

* [AccountSummary](./AccountSummary.md) · [TradeStatsEquityHistory](./TradeStatsEquityHistory.md) · [EquityHistory](./EquityHistory.md) · [TradeUnrealizedPL](./TradeUnrealizedPL.md)
* **Back to:** [1. Account Information — Overview](./Account_Information.Overview.md) · [MT4Account — Overview](../Overview.md)
