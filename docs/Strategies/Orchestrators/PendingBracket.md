# Pending Bracket

> Place a BUY LIMIT, wait for the fill (cancel on timeout), then bracket and manage.

* **Source:** [`mt4/strategies/pending_bracket.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/strategies/pending_bracket.go)
* **Entry:** `RunPendingBracket(ctx, sugar, sp, rp, timeoutSeconds) (int32, error)`
* **Demo:** examples app → menu `41` (`strategydemo.Demo`)

---

## What it does

1. Picks the entry price — `sp.EntryPrice` if set, otherwise **5 pips below the current bid** (demo default).
2. Sizes the lot and places a **BUY LIMIT** via [`MT4Sugar.BuyLimit`](../../MT4Sugar/5.%20Order_Entry.md).
3. Waits up to `timeoutSeconds` for the fill with [`WaitFilled`](../../MT4Sugar/8.%20Watchers_Waiters.md). **If it doesn't fill, the pending is cancelled** and the ticket is returned.
4. Once filled: attaches the pip bracket, manages ~5s (breakeven + gated trailing), then closes.

Use it for limit entries where you want a bounded wait and automatic cleanup if price never comes back.

---

## Signature

```go
func RunPendingBracket(ctx context.Context, sugar *mt4.MT4Sugar, sp StrategyPreset, rp RiskPreset, timeoutSeconds int) (int32, error)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `ctx` | `context.Context` | Cancellation / timeout. |
| `sugar` | `*mt4.MT4Sugar` | Connected high-level client. |
| `sp` | [`StrategyPreset`](../Presets/Presets.md#strategypreset) | `Symbol`, `EntryPrice` (optional — else 5 pips below bid), `Lots`, `Magic`. |
| `rp` | [`RiskPreset`](../Presets/Presets.md#riskpreset) | Risk & protection. |
| `timeoutSeconds` | `int` | How long to wait for the pending to fill before cancelling. |

**Returns:** the `ticket` (`int32`) — filled-and-managed, or the cancelled pending — and an `error`.

---

## Example

```go
// explicit limit price, 60s fill window
ticket, err := strategies.RunPendingBracket(ctx, sugar,
    strategies.LimitEURUSD(1.0750), strategies.Conservative(), 60)
if err != nil { return err }
fmt.Println("pending-bracket ticket:", ticket)
```

> With `strategies.MarketEURUSD()` (no `EntryPrice`) the limit is placed 5 pips below the current bid — handy for a quick demo fill.

---

## 📚 See also

* [Strategies — Master Overview](../Strategies.Master.Overview.md)
* [Presets](../Presets/Presets.md)
* [Market One-Shot](./MarketOneShot.md) · [Bracket + Trailing](./BracketTrailing.md)
* [`WaitFilled`](../../MT4Sugar/8.%20Watchers_Waiters.md) · [`BuyLimit`](../../MT4Sugar/5.%20Order_Entry.md)
