# TradeUnrealizedPL

> Per-trade floating-P/L extrema — the best (MFE) and worst (MAE) each trade reached while open.

**API Information:**

* **Low-level API:** `MT4Account.TradeUnrealizedPL(ctx, from, fromUtc)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **Payload:** `TradeUnrealizedPLData` → `GetExtrema()` → `TradeUnrealizedExtremaData`

## 💬 Just the essentials

* **What it is.** For each trade, the maximum favorable and maximum adverse floating P/L it ever showed.
* **Why you need it.** MFE/MAE analysis — were stops too tight? targets too greedy?
* **When to use.** Post-trade analytics, stop/target tuning.

---

## 🎯 Purpose

* Max favorable excursion (MFE) per trade — peak unrealized profit
* Max adverse excursion (MAE) per trade — worst unrealized loss
* The timestamps at which each extreme occurred
* Ticket + symbol to join back to the trade

---

## Method Signature

```go
func (a *MT4Account) TradeUnrealizedPL(ctx context.Context, from, fromUtc *string) (*pb.TradeUnrealizedPLData, error)
```

**Request message:**

```protobuf
message TradeUnrealizedPLRequest {
  optional string from = 1;      // server-time "yyyy-MM-ddTHH:mm:ss"; nil = account start
  optional string from_utc = 2;  // same moment in UTC
}
```

**Payload message:**

```protobuf
message TradeUnrealizedPLData { repeated TradeUnrealizedExtremaData extrema = 1; }

message TradeUnrealizedExtremaData {
  int64  ticket = ...;
  string symbol = ...;
  double max_unrealized_pl = ...;   // MFE
  Timestamp max_unrealized_time = ...;
  double min_unrealized_pl = ...;   // MAE
  Timestamp min_unrealized_time = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `from` | `*string` | no | Server-time start `"yyyy-MM-ddTHH:mm:ss"`; `nil` = account start. |
| `fromUtc` | `*string` | no | Same start in UTC (send one of `from` / `fromUtc`). |

---

## ⬆️ Output — `*pb.TradeUnrealizedPLData`

`GetExtrema()` returns `[]*pb.TradeUnrealizedExtremaData`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `ticket` | int64 | `GetTicket() int64` | Trade ticket |
| `symbol` | string | `GetSymbol() string` | Instrument |
| `max_unrealized_pl` | double | `GetMaxUnrealizedPl() float64` | Peak floating profit (MFE) |
| `max_unrealized_time` | Timestamp | `GetMaxUnrealizedTime()` / `…Utc()` | When MFE occurred |
| `min_unrealized_pl` | double | `GetMinUnrealizedPl() float64` | Worst floating loss (MAE) |
| `min_unrealized_time` | Timestamp | `GetMinUnrealizedTime()` / `…Utc()` | When MAE occurred |

---

## 🧩 Notes & Tips

* **MFE/MAE** = how far in profit / loss a trade went *while open*, independent of where it closed.
* **Tune stops with MAE.** If winners rarely dip below −X, a tighter stop may be safe.
* **Join on `GetTicket()`** to combine with [`OrdersHistory`](../2.%20Orders_History/OrdersHistory.md) outcomes.

---

## 🔗 Usage Examples

### 1) List each trade's MFE/MAE

```go
func excursions(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.TradeUnrealizedPL(ctx, nil, nil)
    if err != nil { return err }
    for _, e := range d.GetExtrema() {
        fmt.Printf("#%d %s  MFE=%+.2f  MAE=%+.2f\n",
            e.GetTicket(), e.GetSymbol(), e.GetMaxUnrealizedPl(), e.GetMinUnrealizedPl())
    }
    return nil
}
```

### 2) Worst adverse excursion across all trades

```go
func worstMAE(ctx context.Context, acc *mt4.MT4Account) (int64, float64, error) {
    d, err := acc.TradeUnrealizedPL(ctx, nil, nil)
    if err != nil { return 0, 0, err }
    var ticket int64
    worst := 0.0
    for _, e := range d.GetExtrema() {
        if e.GetMinUnrealizedPl() < worst {
            worst, ticket = e.GetMinUnrealizedPl(), e.GetTicket()
        }
    }
    return ticket, worst, nil
}
```

### 3) Trades that gave back a big profit

```go
func gaveBackProfit(ctx context.Context, acc *mt4.MT4Account, minPeak float64) error {
    d, err := acc.TradeUnrealizedPL(ctx, nil, nil)
    if err != nil { return err }
    for _, e := range d.GetExtrema() {
        if e.GetMaxUnrealizedPl() >= minPeak {
            fmt.Printf("#%d peaked at %+.2f\n", e.GetTicket(), e.GetMaxUnrealizedPl())
        }
    }
    return nil
}
```

---

## ⚠️ Pitfalls

* **Empty result.** No qualifying trades → zero extrema; guard `len(GetExtrema())`.
* **Two time axes.** Don't send both `from` and `fromUtc`.
* **Sampling.** Extrema are computed from recorded floating-P/L samples — treat them as close, not tick-exact.

---

## 🧪 Testing Suggestions

* **Happy path (demo with history).** `len(GetExtrema()) > 0`; for each, `MaxUnrealizedPl ≥ MinUnrealizedPl`.
* **Join.** Every `GetTicket()` also appears in [`OrdersHistory`](../2.%20Orders_History/OrdersHistory.md).
* **Window.** A recent `from` returns fewer extrema.

---

## 📚 See also

* [TradeStats](./TradeStats.md) · [OrdersHistory](../2.%20Orders_History/OrdersHistory.md) · [AccountSummary](./AccountSummary.md)
* **Back to:** [1. Account Information — Overview](./Account_Information.Overview.md) · [MT4Account — Overview](../Overview.md)
