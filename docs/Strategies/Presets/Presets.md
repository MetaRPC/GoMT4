# Presets

> Two small structs that decouple **what** to trade from **how much** to risk — plus a ready-made catalog.

* **Source:** [`mt4/strategies/presets.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/strategies/presets.go)

The orchestrators take one of each: a `StrategyPreset` (the *what*) and a `RiskPreset` (the *how much / how to protect*).

---

## StrategyPreset

**WHAT to trade.**

```go
type StrategyPreset struct {
    Symbol     string
    UseMarket  bool
    EntryPrice *float64 // limit price when !UseMarket
    Lots       *float64 // fixed volume; nil = size by risk
    Magic      int32
    Comment    string
}
```

| Field | Type | Description |
|-------|------|-------------|
| `Symbol` | `string` | Instrument, e.g. `"EURUSD"`. |
| `UseMarket` | `bool` | `true` = market entry; `false` = pending (uses `EntryPrice`). |
| `EntryPrice` | `*float64` | Limit price when `!UseMarket`; `nil` lets the orchestrator pick one. |
| `Lots` | `*float64` | Fixed volume; `nil` = size by risk from the `RiskPreset`. |
| `Magic` | `int32` | Magic number tag. |
| `Comment` | `string` | Order comment. |

```go
// build one by hand
lot := 0.10
sp := strategies.StrategyPreset{Symbol: "GBPUSD", UseMarket: true, Magic: 101}
sp.Lots = &lot // fixed volume instead of risk-based sizing
```

---

## RiskPreset

**HOW MUCH to risk and how to protect.** Fields left nil/zero are simply not applied (no TP, no trailing, no breakeven).

```go
type RiskPreset struct {
    RiskPercent   float64
    SLPips        float64
    TPPips        *float64
    TrailingPips  *float64
    BETriggerPips *float64
    BEPlusPips    float64
}
```

| Field | Type | Description |
|-------|------|-------------|
| `RiskPercent` | `float64` | % of balance to risk (used when `Lots` is nil). |
| `SLPips` | `float64` | Stop-loss distance in pips (also drives risk sizing). |
| `TPPips` | `*float64` | Take-profit distance in pips; `nil` = no TP. |
| `TrailingPips` | `*float64` | Trailing distance in pips; `nil` = no trailing. |
| `BETriggerPips` | `*float64` | Profit at which breakeven arms; `nil` = no breakeven. |
| `BEPlusPips` | `float64` | Pips of profit to lock when breakeven fires. |

```go
tp := 40.0
rp := strategies.RiskPreset{RiskPercent: 1.0, SLPips: 20, TPPips: &tp}
```

---

## Catalog

Ready-made presets so demos read as intent:

| Preset | Kind | Values |
|--------|------|--------|
| `MarketEURUSD()` | Strategy | Market EURUSD, magic 777 |
| `LimitEURUSD(price)` | Strategy | BUY LIMIT EURUSD at `price`, magic 778 |
| `Conservative()` | Risk | 0.5% risk, SL 25, TP 50 |
| `Balanced()` | Risk | 1.0% risk, SL 20, TP 40, trail 15, breakeven 10 / +2 |
| `Aggressive()` | Risk | 2.0% risk, SL 15, TP 45, trail 10, breakeven 8 / +1 |

```go
ticket, _ := strategies.RunMarketOneShot(ctx, sugar,
    strategies.MarketEURUSD(), strategies.Balanced())
```

---

## 📚 See also

* [Strategies — Master Overview](../Strategies.Master.Overview.md)
* [Market One-Shot](../Orchestrators/MarketOneShot.md) · [Bracket + Trailing](../Orchestrators/BracketTrailing.md) · [Pending Bracket](../Orchestrators/PendingBracket.md)
* [`CalcLotByRisk`](../../MT4Sugar/4.%20Risk_Sizing.md)
