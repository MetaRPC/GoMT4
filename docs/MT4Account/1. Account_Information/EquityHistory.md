# EquityHistory

> The equity curve — a time series of balance, equity, realized/unrealized P/L and drawdown points.

**API Information:**

* **Low-level API:** `MT4Account.EquityHistory(ctx, from, fromUtc, excludeSameBars)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **Payload:** `EquityHistoryData` → `GetPoints()` → `EquityPointData`

## 💬 Just the essentials

* **What it is.** A series of equity-curve points over time.
* **Why you need it.** Plot the curve, compute drawdown, feed analytics.
* **When to use.** Reporting, equity charts, risk review.

---

## 🎯 Purpose

* Per-point balance and equity
* Realized and unrealized P/L at each point
* Balance / equity drawdown (raw and relative)
* Cash flow (deposits/withdrawals) markers
* Server and UTC timestamps

---

## Method Signature

```go
func (a *MT4Account) EquityHistory(ctx context.Context, from, fromUtc *string, excludeSameBars *bool) (*pb.EquityHistoryData, error)
```

**Request message:**

```protobuf
message EquityHistoryRequest {
  optional string from = 1;              // server-time "yyyy-MM-ddTHH:mm:ss"; nil = account start
  optional string from_utc = 2;          // same moment in UTC
  optional bool   exclude_same_bars = 3; // collapse consecutive equal points
}
```

**Payload message:**

```protobuf
message EquityHistoryData { repeated EquityPointData points = 1; }

message EquityPointData {
  double balance = ...;
  double equity = ...;
  double realized_pl = ...;
  double unrealized_pl = ...;
  double cash_flow = ...;
  double balance_drawdown_relative = ...;
  double equity_drawdown_relative = ...;
  Timestamp time = ...;
  Timestamp time_utc = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `from` | `*string` | no | Server-time start `"yyyy-MM-ddTHH:mm:ss"`; `nil` = account start. |
| `fromUtc` | `*string` | no | Same start in UTC (send one of `from` / `fromUtc`). |
| `excludeSameBars` | `*bool` | no | `true` collapses consecutive identical points (smaller payload). |

---

## ⬆️ Output — `*pb.EquityHistoryData`

`GetPoints()` returns `[]*pb.EquityPointData`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `balance` | double | `GetBalance() float64` | Balance at this point |
| `equity` | double | `GetEquity() float64` | Equity at this point |
| `realized_pl` | double | `GetRealizedPl() float64` | Realized P/L so far |
| `unrealized_pl` | double | `GetUnrealizedPl() float64` | Floating P/L at this point |
| `cash_flow` | double | `GetCashFlow() float64` | Deposit/withdrawal at this point |
| `balance_drawdown_relative` | double | `GetBalanceDrawdownRelative() float64` | Balance drawdown (%) |
| `equity_drawdown_relative` | double | `GetEquityDrawdownRelative() float64` | Equity drawdown (%) |
| `time` / `time_utc` | Timestamp | `GetTime()` / `GetTimeUtc()` | Point time (server / UTC) |

---

## 🧩 Notes & Tips

* **`excludeSameBars`** trims flat stretches — handy for long histories / lighter charts.
* **Points are ordered in time** — iterate directly to build a curve.
* **Drawdown is provided** per point (`*DrawdownRelative`) — no need to recompute the running peak.

---

## 🔗 Usage Examples

### 1) Print the curve length and last equity

```go
func equityTail(ctx context.Context, acc *mt4.MT4Account) error {
    h, err := acc.EquityHistory(ctx, nil, nil, nil)
    if err != nil { return err }
    pts := h.GetPoints()
    fmt.Printf("%d points\n", len(pts))
    if n := len(pts); n > 0 {
        last := pts[n-1]
        fmt.Printf("last: equity=%.2f balance=%.2f dd=%.2f%%\n",
            last.GetEquity(), last.GetBalance(), last.GetEquityDrawdownRelative())
    }
    return nil
}
```

### 2) Worst drawdown over the curve

```go
func worstDrawdown(ctx context.Context, acc *mt4.MT4Account) (float64, error) {
    h, err := acc.EquityHistory(ctx, nil, nil, nil)
    if err != nil { return 0, err }
    var worst float64
    for _, p := range h.GetPoints() {
        if dd := p.GetEquityDrawdownRelative(); dd > worst { worst = dd }
    }
    return worst, nil
}
```

### 3) Windowed + collapsed (lighter payload)

```go
func monthCurve(ctx context.Context, acc *mt4.MT4Account) (int, error) {
    from := time.Now().AddDate(0, -1, 0).Format("2006-01-02T15:04:05")
    collapse := true
    h, err := acc.EquityHistory(ctx, &from, nil, &collapse)
    if err != nil { return 0, err }
    return len(h.GetPoints()), nil
}
```

---

## ⚠️ Pitfalls

* **Large histories.** Without `excludeSameBars` the payload can be big — collapse or window it.
* **Two time axes.** Don't send both `from` and `fromUtc`.
* **Empty result.** A future `from` returns zero points — guard `len(GetPoints())`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `len(GetPoints()) > 0`; timestamps increase; equity is finite.
* **Collapse.** With `excludeSameBars=true` the count is ≤ the uncollapsed count.
* **Window.** A recent `from` yields fewer points than the all-time call.

---

## 📚 See also

* [TradeStats](./TradeStats.md) · [TradeStatsEquityHistory](./TradeStatsEquityHistory.md) · [AccountSummary](./AccountSummary.md)
* **Back to:** [1. Account Information — Overview](./Account_Information.Overview.md) · [MT4Account — Overview](../Overview.md)
