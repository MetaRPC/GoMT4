# OnSymbolTick (RECOMMENDED)

> Real-time bid/ask/last ticks for one or more symbols, pushed as they change.

**API Information:**

* **Low-level API:** `MT4Account.OnSymbolTick(ctx, symbols)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.Subscriptions`
* **RPC:** `OnSymbolTick(OnSymbolTickRequest) → stream OnSymbolTickReply`
* **Payload:** `OnSymbolTickData` → `GetSymbolTick()` → `OnSymbolMqlTickInfo`

## 💬 Just the essentials

* **What it is.** A live tick feed — one frame per price change.
* **Why you need it.** Drive pricing, trailing, dashboards without polling.
* **When to use.** Any time you need the current market as it moves.

---

## 🎯 Purpose

* Stream live bid/ask/last for a set of symbols
* React to price changes in real time
* Feed trailing/breakeven and live charts
* Auto-reconnects on transient transport errors

---

## Method Signature

```go
func (a *MT4Account) OnSymbolTick(ctx context.Context, symbols []string) (<-chan *pb.OnSymbolTickData, <-chan error)
```

**Request message:**

```protobuf
message OnSymbolTickRequest { repeated string symbol_names = 1; }
```

**Payload message:**

```protobuf
message OnSymbolTickData { OnSymbolMqlTickInfo symbol_tick = 1; }

message OnSymbolMqlTickInfo {
  string symbol = ...;
  double bid = ...;  double ask = ...;  double last = ...;
  uint64 volume = ...;
  Timestamp time = ...;  Timestamp time_utc = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancel it to unsubscribe and close the channels. |
| `symbols` | `[]string` | yes | Symbols to stream, e.g. `[]string{"EURUSD","USDJPY"}`. |

---

## ⬆️ Frame — `*pb.OnSymbolTickData`

`GetSymbolTick()` returns `*pb.OnSymbolMqlTickInfo`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `symbol` | string | `GetSymbol() string` | Symbol |
| `bid` / `ask` | double | `GetBid()` / `GetAsk()` | Best bid / ask |
| `last` | double | `GetLast() float64` | Last trade price |
| `volume` | uint64 | `GetVolume() uint64` | Tick volume |
| `time` / `time_utc` | Timestamp | `GetTime()` / `GetTimeUtc()` | Tick time (server / UTC) |

> A non-nil value on the **error channel** ends the stream. Both channels close when `ctx` is cancelled.

---

## 🧩 Notes & Tips

* **Subscribe once per symbol set** — pass every symbol in the one slice.
* **`cancel()` unsubscribes** — the returned channels close.
* **Auto-reconnects** on `codes.Unavailable`; you keep reading the same channel.
* **Guard the tick** — `GetSymbolTick()` can be `nil` on a keep-alive frame.

---

## 🔗 Usage Examples

### 1) Stream ticks for 10 seconds

```go
func streamTicks(ctx context.Context, acc *mt4.MT4Account, symbols []string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    data, errs := acc.OnSymbolTick(ctx, symbols)
    for {
        select {
        case d := <-data:
            t := d.GetSymbolTick()
            if t == nil { continue } // keep-alive
            fmt.Printf("%s  %.5f / %.5f\n", t.GetSymbol(), t.GetBid(), t.GetAsk())
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 2) Take the first N ticks, then stop

```go
func firstN(ctx context.Context, acc *mt4.MT4Account, symbol string, n int) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    data, errs := acc.OnSymbolTick(ctx, []string{symbol})
    for got := 0; got < n; {
        select {
        case d := <-data:
            if t := d.GetSymbolTick(); t != nil {
                fmt.Printf("tick %d: %.5f\n", got+1, t.GetBid())
                got++
            }
        case err := <-errs:
            if err != nil { return err }
        }
    }
    return nil // cancel() unsubscribes
}
```

### 3) Keep a latest-price cache

```go
func priceCache(ctx context.Context, acc *mt4.MT4Account, symbols []string, cache map[string]float64) error {
    data, errs := acc.OnSymbolTick(ctx, symbols)
    for {
        select {
        case d := <-data:
            if t := d.GetSymbolTick(); t != nil {
                cache[t.GetSymbol()] = (t.GetBid() + t.GetAsk()) / 2
            }
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

---

## ⚠️ Pitfalls

* **Nil tick** on keep-alive frames — guard before use.
* **High frequency.** Don't block the `select`; offload heavy work to another goroutine.
* **Leaking the context** leaks the stream goroutine — always `defer cancel()`.

---

## 🧪 Testing Suggestions

* **Happy path (demo, market open).** Frames arrive within seconds; `GetBid()`/`GetAsk()` are positive.
* **Cancel.** Cancelling the context ends the loop and closes the channels.
* **Idle market.** On a closed market, expect few/no frames (not an error).

---

## 📚 See also

* [Quote](../3.%20Symbols_Quotes/Quote.md) · [OnTrade](./OnTrade.md) · [Streaming — Overview](./Streaming_Methods.Overview.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
