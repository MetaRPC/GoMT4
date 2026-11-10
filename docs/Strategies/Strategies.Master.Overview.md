# Strategies — Master Overview

> One page to orient fast: the ready-made **orchestrators** and **presets** that ship with GoMT4, how to pick one, and links to each spec.

---

## ⚠️ Educational & demonstration only

The orchestrators and presets here are **demonstration examples** — templates for building your own strategies, not production trading systems.

- Simplified logic, no portfolio-level risk controls, no economic-calendar integration.
- **Demo accounts only.** Study the code, test on demo, add your own risk management before considering anything live.
- Trading involves substantial risk of loss; you are responsible for your decisions.

They live in the repo under [`mt4/strategies/`](https://github.com/MetaRPC/GoMT4/tree/main/mt4/strategies) and are driven by [`MT4Sugar`](../MT4Sugar/Overview.md) (which composes [`MT4Service`](../MT4Service/Overview.md) → [`MT4Account`](../MT4Account/Overview.md)). The examples app runs the demo via menu option **`41`**.

---

## 🚦 Orchestrators

Each orchestrator is a single `func Run*(ctx, sugar, sp, rp, …) (int32, error)` that takes a [`StrategyPreset`](./Presets/Presets.md#strategypreset) (what / where to enter) and a [`RiskPreset`](./Presets/Presets.md#riskpreset) (how much to risk / how to protect), runs the trade lifecycle, and returns the ticket.

| Orchestrator | Entry func | What it does |
|--------------|-----------|--------------|
| **[Market One-Shot](./Orchestrators/MarketOneShot.md)** | `RunMarketOneShot(ctx, sugar, sp, rp)` | Size → market BUY → pip SL/TP bracket → manage briefly → close. |
| **[Bracket + Trailing](./Orchestrators/BracketTrailing.md)** | `RunBracketTrailing(ctx, sugar, sp, rp, activateAtPips)` | Market entry + bracket; trailing stop switches on once price is far enough in profit. |
| **[Pending Bracket](./Orchestrators/PendingBracket.md)** | `RunPendingBracket(ctx, sugar, sp, rp, timeoutSeconds)` | Place a BUY LIMIT, wait for the fill (cancel on timeout), then bracket + manage. |

All three BUY for the demo, auto-size by risk (unless the preset pins `Lots`), and flatten at the end.

---

## 🎛️ Presets

Two small structs decouple **what** from **how much**:

* **[`StrategyPreset`](./Presets/Presets.md#strategypreset)** — symbol, market/limit, entry price, fixed lots, magic, comment.
* **[`RiskPreset`](./Presets/Presets.md#riskpreset)** — risk %, SL/TP pips, trailing, breakeven.

Ready-made catalog: `MarketEURUSD()`, `LimitEURUSD(price)` · `Conservative()`, `Balanced()`, `Aggressive()`. See **[Presets](./Presets/Presets.md)**.

---

## 🏁 Running the demo

```go
sugar := mt4.NewMT4Sugar(mt4.NewMT4Service(acc)).SetDefaults(777, 20, "demo")

// one protected market trade
ticket, err := strategies.RunMarketOneShot(ctx, sugar,
    strategies.MarketEURUSD(), // WHAT
    strategies.Balanced())     // HOW
```

Or run all three from the examples app: launch it and choose **`41`** (`strategydemo.Demo`).

---

## 📚 See also

* [Market One-Shot](./Orchestrators/MarketOneShot.md) · [Bracket + Trailing](./Orchestrators/BracketTrailing.md) · [Pending Bracket](./Orchestrators/PendingBracket.md) · [Presets](./Presets/Presets.md)
* **Layers:** [MT4Sugar](../MT4Sugar/Overview.md) · [MT4Service](../MT4Service/Overview.md) · [MT4Account](../MT4Account/Overview.md)
