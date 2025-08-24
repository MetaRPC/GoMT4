# 💱 GetQuote (GoMT4)

**Goal:** fetch the **latest quote** (Bid/Ask/Spread) for a symbol.

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (`Quote`)
> * Example: `examples/mt4/MT4_service.go` (`ShowQuote`)

---

## ✅ 1) Preconditions

* Symbol exists and is visible in MT4 (*Market Watch* → Show All).
* `config.json` has a valid `DefaultSymbol` or you specify a symbol manually.

---

## 📝 2) Request one quote

```go
q, err := account.Quote(ctx, symbol)
if err != nil { return err }

fmt.Printf("%s: Bid=%.5f Ask=%.5f Spread=%.1f pips Time=%s\n",
    symbol,
    q.GetBid(),
    q.GetAsk(),
    (q.GetAsk()-q.GetBid())/q.GetPoint(),
    q.GetTime().AsTime().Format(time.RFC3339))
```

---

## 🔍 3) Inspect fields

`q` (`*pb.Quote`) contains:

* `Bid` — broker’s buy price.
* `Ask` — broker’s sell price.
* `Point` — 10^-Digits (used to compute spread in points/pips).
* `Digits` — precision for price rounding.
* `Time` — server timestamp.

---

## ⚠️ Pitfalls

* **symbol not found** → check suffix (`EURUSD.m`, etc.) or ensure visible in MT4.
* **timeout** → wrap in `context.WithTimeout` (2–3s) and retry on transient errors only.
* **spread mismatch** → some brokers quote in fractional pips; always divide `(Ask-Bid)` by `Point`.

---

## 🔄 Variations

* **Default symbol from config**:

```go
sym := cfg.DefaultSymbol
q, _ := account.Quote(ctx, sym)
```

* **Loop for multiple quotes**: see `GetMultipleQuotes.md` recipe.
* **Subscribe stream**: see `StreamQuotes.md` recipe.

---

## 📎 See also

* `GetMultipleQuotes.md` — batch query for multiple symbols.
* `StreamQuotes.md` — live stream of quotes.
* `SymbolParams.md` — detailed info about Digits/LotStep/MinLot.
