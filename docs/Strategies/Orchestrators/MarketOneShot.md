# Market One-Shot

> Open a market order with an automatic pip-based SL/TP bracket — the simplest orchestrator.

* **Source:** [`mt4/strategies/market_one_shot.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/strategies/market_one_shot.go)
* **Entry:** `RunMarketOneShot(ctx, sugar, sp, rp) (int32, error)`
* **Demo:** examples app → menu `41` (`strategydemo.Demo`)

---

## What it does

1. Sizes the lot — fixed `sp.Lots`, or risk-based from `rp.RiskPercent` over `rp.SLPips` via [`CalcLotByRisk`](../../MT4Sugar/4.%20Risk_Sizing.md).
2. Opens a **market BUY** via [`MT4Sugar.BuyMarket`](../../MT4Sugar/5.%20Order_Entry.md).
3. Attaches a bracket — SL/TP as pip distances with [`ModifySLTPByPips`](../../MT4Sugar/6.%20Order_Management.md).
4. Manages for ~5s (auto-breakeven, then gated trailing — both from the preset).
5. Closes the position and returns the ticket.

Use it when you simply want a protected market entry in one call.

---

## Signature

```go
func RunMarketOneShot(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset) (int32, error)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Cancellation / timeout. |
| `sugar` | `*mt4.MT4Sugar` | A connected high-level client (drives entry, bracket, manage, close). |
| `sp` | [`StrategyPreset`](../Presets/Presets.md#strategypreset) | `Symbol`, `Lots` (optional), `Magic`, `Comment`. |
| `rp` | [`RiskPreset`](../Presets/Presets.md#riskpreset) | `RiskPercent`, `SLPips`, `TPPips`, optional trailing/breakeven. |

**Returns:** the order `ticket` (`int32`) and an `error`.

---

## Example

```go
sugar := mt4.NewMT4Sugar(mt4.NewMT4Service(acc)).SetDefaults(777, 20, "demo")

ticket, err := strategies.RunMarketOneShot(ctx, sugar,
    strategies.MarketEURUSD(), // WHAT: market EURUSD, magic 777
    strategies.Balanced())     // HOW: 1% risk, SL 20 / TP 40, trail 15, BE 10/+2
if err != nil { return err }
fmt.Println("one-shot ticket:", ticket)
```

> If `sp.Lots` is nil, sizing comes from `rp.RiskPercent` + `rp.SLPips` via `CalcLotByRisk`.

---

## 📚 See also

* [Strategies — Master Overview](../Strategies.Master.Overview.md)
* [Presets](../Presets/Presets.md)
* [Bracket + Trailing](./BracketTrailing.md) · [Pending Bracket](./PendingBracket.md)
