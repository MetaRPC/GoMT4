# TradeStatsEquityHistory

> Stats **and** the equity curve **and** the trade list — the combined analytics payload in one call.

**API Information:**

* **Low-level API:** `MT4Account.TradeStatsEquityHistory(ctx, from, fromUtc, excludeSameBars)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **Payload:** `StatsWithChartsData`

## 💬 Just the essentials

* **What it is.** Everything [`TradeStats`](./TradeStats.md) and [`EquityHistory`](./EquityHistory.md) return, plus the trade list — in a single round-trip.
* **Why you need it.** Build a full report page without three separate calls.
* **When to use.** Report/export screens that need KPIs, the curve and the trades together.

---

## 🎯 Purpose

* The full `StatsData` KPI block
* The equity-curve points (`EquityPointData`)
* The closed-trade list and the currently-open trades
* One round-trip for a complete analytics view

---

## Method Signature

```go
func (a *MT4Account) TradeStatsEquityHistory(ctx context.Context, from, fromUtc *string, excludeSameBars *bool) (*pb.StatsWithChartsData, error)
```

**Request message:**

```protobuf
message TradeStatsEquityHistoryRequest {
  optional string from = 1;              // server-time "yyyy-MM-ddTHH:mm:ss"; nil = account start
  optional string from_utc = 2;
  optional bool   exclude_same_bars = 3; // collapse consecutive equal equity points
}
```

**Payload message:**

```protobuf
message StatsWithChartsData {
  StatsData stats = 1;
  repeated EquityPointData charts = 2;
  repeated TradeHistoryItemData trades = 3;
  repeated TradeHistoryItemData open_trades = 4;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `from` | `*string` | no | Server-time start `"yyyy-MM-ddTHH:mm:ss"`; `nil` = account start. |
| `fromUtc` | `*string` | no | Same start in UTC (send one of `from` / `fromUtc`). |
| `excludeSameBars` | `*bool` | no | `true` collapses consecutive identical equity points. |

---

## ⬆️ Output — `*pb.StatsWithChartsData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `stats` | StatsData | `GetStats()` | Full KPI block — see [`TradeStats`](./TradeStats.md) |
| `charts` | repeated EquityPointData | `GetCharts()` | Equity-curve points — see [`EquityHistory`](./EquityHistory.md) |
| `trades` | repeated TradeHistoryItemData | `GetTrades()` | Closed trades in the window |
| `open_trades` | repeated TradeHistoryItemData | `GetOpenTrades()` | Currently open trades |

---

## 🧩 Notes & Tips

* **One call, three datasets.** Prefer this over calling `TradeStats` + `EquityHistory` separately for a report.
* **Nested `stats` can be nil** on an empty history — guard `GetStats()` before reading KPIs.
* **`excludeSameBars`** applies to `charts` (the equity points), trimming flat stretches.

---

## 🔗 Usage Examples

### 1) One-call report header

```go
func reportHeader(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.TradeStatsEquityHistory(ctx, nil, nil, nil)
    if err != nil { return err }
    if s := d.GetStats(); s != nil {
        fmt.Printf("trades=%d  PF=%.2f  maxDD=%.1f%%\n",
            s.GetTrades(), s.GetProfitFactor(), s.GetMaxEquityDrawdownRelative())
    }
    fmt.Printf("curve points=%d  closed=%d  open=%d\n",
        len(d.GetCharts()), len(d.GetTrades()), len(d.GetOpenTrades()))
    return nil
}
```

### 2) Export the equity curve

```go
func exportCurve(ctx context.Context, acc *mt4.MT4Account) error {
    collapse := true
    d, err := acc.TradeStatsEquityHistory(ctx, nil, nil, &collapse)
    if err != nil { return err }
    for _, p := range d.GetCharts() {
        fmt.Printf("%s,%.2f,%.2f\n", p.GetTime().AsTime().Format(time.RFC3339),
            p.GetBalance(), p.GetEquity())
    }
    return nil
}
```

### 3) Month-to-date bundle

```go
func monthBundle(ctx context.Context, acc *mt4.MT4Account) (*pb.StatsWithChartsData, error) {
    from := time.Now().Format("2006-01") + "-01T00:00:00"
    return acc.TradeStatsEquityHistory(ctx, &from, nil, nil)
}
```

---

## ⚠️ Pitfalls

* **Heavier payload.** It carries stats + curve + trades — window and `excludeSameBars` it for big histories.
* **Nil sub-messages.** `GetStats()` may be `nil`; the slices may be empty — guard both.
* **Two time axes.** Don't send both `from` and `fromUtc`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `GetStats() != nil`; `len(GetCharts()) > 0`; `len(GetTrades())` matches `GetStats().GetTrades()` for the all-time call.
* **Collapse.** `excludeSameBars=true` reduces `len(GetCharts())`.
* **Window.** A recent `from` shrinks all three datasets.

---

## 📚 See also

* [TradeStats](./TradeStats.md) · [EquityHistory](./EquityHistory.md) · [TradeUnrealizedPL](./TradeUnrealizedPL.md)
* **Back to:** [1. Account Information — Overview](./Account_Information.Overview.md) · [MT4Account — Overview](../Overview.md)
