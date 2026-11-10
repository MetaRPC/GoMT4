# SymbolSelect

> Add a symbol to Market Watch (or remove it) so it publishes quotes.

**API Information:**

* **Low-level API:** `MT4Account.SymbolSelect(ctx, symbol, selected)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.MarketInfo`
* **RPC:** `SymbolSelect(SymbolSelectRequest) → SymbolSelectReply`
* **Returns:** `bool`

## 💬 Just the essentials

* **What it is.** MQL `SymbolSelect()` — show/hide a symbol in Market Watch.
* **Why you need it.** A symbol must be selected to stream ticks / return fresh quotes.
* **When to use.** Before quoting or trading a symbol you haven't touched this session.

---

## 🎯 Purpose

* Enable a symbol so it publishes quotes
* Prepare a symbol before `Quote` / `OnSymbolTick`
* Optionally hide a symbol you no longer need

---

## Method Signature

```go
func (a *MT4Account) SymbolSelect(ctx context.Context, symbol string, selected bool) (bool, error)
```

**Request message:**

```protobuf
message SymbolSelectRequest {
  string symbol = 1;
  bool   selected = 2;   // true = add to Market Watch, false = remove
}
```

**Returns:** `bool` — whether the server accepted the change.

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbol` | `string` | yes | Instrument to select, e.g. `"EURUSD"`. |
| `selected` | `bool` | yes | `true` adds to Market Watch; `false` removes. |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| ok | `bool` | `true` if the server accepted the select/deselect |

> On failure the call returns `false` and a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Select before you quote.** A symbol that isn't in Market Watch may return stale/empty quotes.
* **Fresh feeds need a moment.** After selecting, the first quote can be `bid == ask`; poll until it's two-sided (the sugar `EnsureSymbol` does this).
* **Don't over-remove.** Deselecting a symbol you later need just forces a re-select.

---

## 🔗 Usage Examples

### 1) Select then quote

```go
func selectAndQuote(ctx context.Context, acc *mt4.MT4Account, symbol string) error {
    if _, err := acc.SymbolSelect(ctx, symbol, true); err != nil { return err }
    q, err := acc.Quote(ctx, symbol)
    if err != nil { return err }
    fmt.Printf("%s %.5f / %.5f\n", symbol, q.GetBid(), q.GetAsk())
    return nil
}
```

### 2) Ensure a usable quote after selecting

```go
func ensure(ctx context.Context, acc *mt4.MT4Account, symbol string, wait time.Duration) (*pb.QuoteData, error) {
    if _, err := acc.SymbolSelect(ctx, symbol, true); err != nil { return nil, err }
    deadline := time.Now().Add(wait)
    for {
        if q, err := acc.Quote(ctx, symbol); err == nil && q.GetBid() > 0 && q.GetAsk() > q.GetBid() {
            return q, nil
        }
        if time.Now().After(deadline) {
            return nil, fmt.Errorf("%s: no usable quote in %s", symbol, wait)
        }
        time.Sleep(500 * time.Millisecond)
    }
}
```

### 3) Hide a symbol

```go
func hide(ctx context.Context, acc *mt4.MT4Account, symbol string) (bool, error) {
    return acc.SymbolSelect(ctx, symbol, false)
}
```

---

## ⚠️ Pitfalls

* **First quote may be degenerate** (`bid == ask`/zeros) — poll until two-sided before trading.
* **Unknown symbol** returns an error — validate with [`Symbols`](./Symbols.md) first.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Selecting a known symbol returns `true`; a subsequent `Quote` is two-sided.
* **Unknown.** Selecting a bogus symbol returns an error / `false`.
* **Round-trip.** Select then deselect without error.

---

## 📚 See also

* [Symbols](./Symbols.md) · [Quote](./Quote.md) · [OnSymbolTick](../5.%20Streaming_Methods/OnSymbolTick.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
