# Bracket + Trailing

> Market entry with a bracket, where the trailing stop only switches on once the trade is far enough in profit.

* **Source:** [`mt4/strategies/bracket_trailing.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/strategies/bracket_trailing.go)
* **Entry:** `RunBracketTrailing(ctx, sugar, sp, rp, activateAtPips) (int32, error)`
* **Demo:** examples app → menu `41` (`strategydemo.Demo`)

---

## What it does

1. Sizes the lot and opens a **market BUY** (like [Market One-Shot](./MarketOneShot.md)).
2. Attaches the pip SL/TP bracket.
3. Manages for ~10s: auto-breakeven applies early; **trailing only arms once the trade is `activateAtPips` in profit** — letting the position breathe before the stop starts following it.
4. Closes and returns the ticket.

Use it when you want a trade to run before the stop begins trailing, rather than trailing from the first pip.

---

## Signature

```go
func RunBracketTrailing(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset, activateAtPips float64) (int32, error)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Cancellation / timeout. |
| `sugar` | `*mt4.MT4Sugar` | Connected high-level client. |
| `sp` | [`StrategyPreset`](../Presets/Presets.md#strategypreset) | What to trade. |
| `rp` | [`RiskPreset`](../Presets/Presets.md#riskpreset) | Risk & protection; `TrailingPips` must be set for trailing to do anything. |
| `activateAtPips` | `float64` | Profit (in pips) at which trailing switches on. |

**Returns:** the order `ticket` (`int32`) and an `error`.

---

## Example

```go
// trailing arms once the trade is +12 pips in profit
ticket, err := strategies.RunBracketTrailing(ctx, sugar,
    strategies.MarketEURUSD(), strategies.Balanced(), 12)
if err != nil { return err }
fmt.Println("bracket+trailing ticket:", ticket)
```

> `Balanced()` sets `TrailingPips = 15` and `BETriggerPips = 10`, so breakeven locks at +10 while the trail follows 15 pips behind once +12 is reached.

---

## 📚 See also

* [Strategies — Master Overview](../Strategies.Master.Overview.md)
* [Presets](../Presets/Presets.md)
* [Market One-Shot](./MarketOneShot.md) · [Pending Bracket](./PendingBracket.md)
* [`TrailStop` / `AutoBreakeven`](../../MT4Sugar/7.%20Trailing_Breakeven.md)
