# QuoteMany

> Latest quotes for several symbols in one round-trip.

**API Information:**

* **Low-level API:** `MT4Account.QuoteMany(ctx, symbols)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.MarketInfo`
* **RPC:** `QuoteMany(QuoteManyRequest) → QuoteManyReply`
* **Payload:** `QuoteManyData` → `GetQuotes()` → `QuoteData`

## 💬 Just the essentials

* **What it is.** The batch version of [`Quote`](./Quote.md) — many symbols, one call.
* **Why you need it.** Price a basket/watchlist without N separate round-trips.
* **When to use.** Dashboards, multi-symbol scans, correlation checks.

---

## 🎯 Purpose

* Fetch current bid/ask for a list of symbols at once
* Build a watchlist / quote board
* Reduce round-trips vs calling `Quote` in a loop

---

## Method Signature

```go
func (a *MT4Account) QuoteMany(ctx context.Context, symbols []string) (*pb.QuoteManyData, error)
```

**Request message:**

```protobuf
message QuoteManyRequest { repeated string symbols = 1; }
```

**Payload message:**

```protobuf
message QuoteManyData { repeated QuoteData quotes = 1; }
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbols` | `[]string` | yes | Symbols to price, e.g. `[]string{"EURUSD","USDJPY"}`. |

---

## ⬆️ Output — `*pb.QuoteManyData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `quotes` | repeated QuoteData | `GetQuotes()` | One `QuoteData` per symbol (see [`Quote`](./Quote.md) for its getters) |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Order / presence.** Match results by `GetSymbol()`, not by slice position — unknown symbols may be omitted.
* **One call beats a loop.** Prefer this over calling `Quote` per symbol for a watchlist.
* **Each element** is a full `QuoteData` (bid/ask/high/low/time).

---

## 🔗 Usage Examples

### 1) Price a watchlist

```go
func board(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.QuoteMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})
    if err != nil { return err }
    for _, q := range d.GetQuotes() {
        fmt.Printf("%-7s %.5f / %.5f\n", q.GetSymbol(), q.GetBid(), q.GetAsk())
    }
    return nil
}
```

### 2) Build a symbol→mid map

```go
func mids(ctx context.Context, acc *mt4.MT4Account, symbols []string) (map[string]float64, error) {
    d, err := acc.QuoteMany(ctx, symbols)
    if err != nil { return nil, err }
    out := make(map[string]float64, len(d.GetQuotes()))
    for _, q := range d.GetQuotes() {
        out[q.GetSymbol()] = (q.GetBid() + q.GetAsk()) / 2
    }
    return out, nil
}
```

### 3) Widest-spread symbol in the basket

```go
func widest(ctx context.Context, acc *mt4.MT4Account, symbols []string) (string, float64, error) {
    d, err := acc.QuoteMany(ctx, symbols)
    if err != nil { return "", 0, err }
    var sym string
    var wide float64
    for _, q := range d.GetQuotes() {
        if s := q.GetAsk() - q.GetBid(); s > wide { wide, sym = s, q.GetSymbol() }
    }
    return sym, wide, nil
}
```

---

## ⚠️ Pitfalls

* **Missing symbols.** An unknown/unselected symbol may simply not appear — don't index by position.
* **Snapshot.** Prices are point-in-time.
* **Large baskets.** Very long lists are heavier; split if needed.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `len(GetQuotes())` matches the number of valid symbols; each has `GetBid() > 0`.
* **Mixed list.** Include one bogus symbol; the valid ones still return.
* **Map.** The symbol→mid map has an entry per returned quote.

---

## 📚 See also

* [Quote](./Quote.md) · [Symbols](./Symbols.md) · [SymbolParamsMany](./SymbolParamsMany.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
