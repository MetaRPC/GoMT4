# Quote (RECOMMENDED)

> The latest bid/ask (and day high/low) for one symbol.

**API Information:**

* **Low-level API:** `MT4Account.Quote(ctx, symbol)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.MarketInfo`
* **RPC:** `Quote(QuoteRequest) → QuoteReply`
* **Payload:** `QuoteData`

## 💬 Just the essentials

* **What it is.** A one-shot current price for a symbol.
* **Why you need it.** Pricing, spread checks, converting pips↔price, pre-trade reads.
* **When to use.** Any time you need "the price now" without a live subscription.

---

## 🎯 Purpose

* Read current bid / ask
* Derive mid price and spread
* Fetch the day's high / low
* Feed price into order sizing and SL/TP math

---

## Method Signature

```go
func (a *MT4Account) Quote(ctx context.Context, symbol string) (*pb.QuoteData, error)
```

**Request message:**

```protobuf
message QuoteRequest { string symbol = 1; }
```

**Payload message:**

```protobuf
message QuoteData {
  string symbol = ...;
  double bid = ...;
  double ask = ...;
  double high = ...;
  double low = ...;
  Timestamp date_time = ...;
  Timestamp date_time_utc = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbol` | `string` | yes | Instrument, e.g. `"EURUSD"`. |

---

## ⬆️ Output — `*pb.QuoteData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `symbol` | string | `GetSymbol() string` | Instrument |
| `bid` | double | `GetBid() float64` | Best bid |
| `ask` | double | `GetAsk() float64` | Best ask |
| `high` / `low` | double | `GetHigh()` / `GetLow()` | Session high / low |
| `date_time` / `date_time_utc` | Timestamp | `GetDateTime()` / `GetDateTimeUtc()` | Quote time (server / UTC) |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Spread** = `GetAsk() - GetBid()`; divide by the pip size for pips.
* **Fresh symbols.** A just-added symbol may send `bid == ask` or zeros at first — see [`SymbolSelect`](./SymbolSelect.md) / the sugar `EnsureSymbol`.
* **Time getter** is `GetDateTime()` (not `GetTime()`).
* **For a live feed** use [`OnSymbolTick`](../5.%20Streaming_Methods/OnSymbolTick.md).

---

## 🔗 Usage Examples

### 1) Print bid/ask

```go
func printQuote(ctx context.Context, acc *mt4.MT4Account, symbol string) error {
    q, err := acc.Quote(ctx, symbol)
    if err != nil { return err }
    fmt.Printf("%s  %.5f / %.5f  (H %.5f L %.5f)\n",
        q.GetSymbol(), q.GetBid(), q.GetAsk(), q.GetHigh(), q.GetLow())
    return nil
}
```

### 2) Mid price and spread in pips

```go
func midAndSpread(ctx context.Context, acc *mt4.MT4Account, symbol string) (mid, spreadPips float64, err error) {
    q, err := acc.Quote(ctx, symbol)
    if err != nil { return 0, 0, err }
    sp, err := acc.SymbolParams(ctx, symbol)
    if err != nil { return 0, 0, err }
    pip := sp.GetPoint()
    if sp.GetDigits() == 3 || sp.GetDigits() == 5 { pip *= 10 }
    mid = (q.GetBid() + q.GetAsk()) / 2
    spreadPips = (q.GetAsk() - q.GetBid()) / pip
    return mid, spreadPips, nil
}
```

### 3) Skip trading on a wide spread

```go
func spreadOK(ctx context.Context, acc *mt4.MT4Account, symbol string, maxPips float64) (bool, error) {
    _, sp, err := midAndSpread(ctx, acc, symbol)
    if err != nil { return false, err }
    return sp <= maxPips, nil
}
```

---

## ⚠️ Pitfalls

* **Snapshot, not a stream.** Re-query for fresh prices, or subscribe via [`OnSymbolTick`](../5.%20Streaming_Methods/OnSymbolTick.md).
* **`bid == ask` on fresh feeds.** Normalize/ensure the symbol before trusting the first quote.
* **Market closed.** Prices may be stale on weekends/holidays.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `GetBid() > 0` and `GetAsk() >= GetBid()` for a liquid symbol.
* **Unknown symbol.** A bogus name returns an error.
* **Spread.** The derived pip spread is a small positive number in normal conditions.

---

## 📚 See also

* [QuoteMany](./QuoteMany.md) · [QuoteHistory](./QuoteHistory.md) · [SymbolParams](./SymbolParams.md) · [OnSymbolTick](../5.%20Streaming_Methods/OnSymbolTick.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
