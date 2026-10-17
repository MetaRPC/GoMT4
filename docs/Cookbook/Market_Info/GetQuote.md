# 💱 GetQuote (GoMT4)

**Goal:** fetch the **latest quote** (Bid/Ask/Spread) for a symbol.

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (`Quote`)
> * Example: `examples/mt4/MT4_service.go` (`ShowQuote`)

---

## ✅ 1) Preconditions

* `account` is connected; the symbol exists on the broker (exact name, including any suffix).
* `config.json` has a valid `DefaultSymbol` or you specify a symbol manually.

---

## 📝 2) Request one quote

`QuoteData` has no point/digits; take `Point` from `SymbolParams` to express the spread in points:

```go
q, err := account.Quote(ctx, symbol)
if err != nil { return err }

p, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

fmt.Printf("%s: Bid=%.5f Ask=%.5f Spread=%.1f points Time=%s\n",
    symbol,
    q.GetBid(),
    q.GetAsk(),
    (q.GetAsk()-q.GetBid())/p.GetPoint(),
    q.GetDateTime().AsTime().Format(time.RFC3339))
```

---

## 🔍 3) Inspect fields

`q` (`*pb.QuoteData`) contains:

* `Symbol` — symbol name.
* `Bid` — price you sell at.
* `Ask` — price you buy at.
* `High` / `Low` — high / low price.
* `DateTime` — quote time (broker server time); `DateTimeUtc` — the same moment in UTC.

`Point` and `Digits` come from `account.SymbolParams(ctx, symbol)`.

---

## ⚠️ Pitfalls

* **symbol not found** → check suffix (`EURUSD.m`, etc.) with `account.Symbols(ctx)`.
* **timeout** → the call has a 3s default timeout; the client retries transient errors itself.
* **spread units** → 5‑digit brokers quote fractional pips; `(Ask-Bid)/Point` is in **points** (10 points = 1 pip on 5‑digit FX).

---

## 🔄 Variations

* **Default symbol from config**:

```go
sym := cfg.DefaultSymbol
q, err := account.Quote(ctx, sym)
```

* **Many symbols at once**: see [GetMultipleQuotes](GetMultipleQuotes.md).
* **Subscribe stream**: see [StreamQuotes](StreamQuotes.md).

---

## 📎 See also

* [GetMultipleQuotes](GetMultipleQuotes.md) — batch query for multiple symbols.
* [StreamQuotes](StreamQuotes.md) — live stream of quotes.
* [SymbolParams](SymbolParams.md) — `Digits`, `Point`, `VolumeStep`, `VolumeMin`.
