# Streaming Real-Time Quotes

> **Request:** subscribe to live tick updates for predefined symbols
> Continuously streams the latest bid/ask prices and timestamps for each symbol.

---

### Code Example

```go
// Using service wrapper
service.StreamQuotes(context.Background())

// Or directly from MT4Account
symbols := []string{"EURUSD", "GBPUSD"}
tickCh, errCh := mt4.OnSymbolTick(context.Background(), symbols)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

fmt.Println("\uD83D\uDD04 Streaming ticks...")
for {
    select {
    case tick, ok := <-tickCh:
        if !ok {
            fmt.Println("✅ Tick stream ended.")
            return
        }
        if sym := tick.GetSymbolTick(); sym != nil {
            fmt.Printf("[Tick] %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
                sym.GetSymbol(),
                sym.GetBid(),
                sym.GetAsk(),
                sym.GetTime().AsTime().Format("2006-01-02 15:04:05"))
        }
    case err := <-errCh:
        log.Printf("❌ Stream error: %v", err)
        return
    case <-time.After(30 * time.Second):
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}
```

---

### Method Signature

```go
func (s *MT4Service) StreamQuotes(ctx context.Context)
```

---

## 🔽 Input

| Field | Type              | Description                                |
| ----- | ----------------- | ------------------------------------------ |
| `ctx` | `context.Context` | Controls stream lifetime and cancellation. |

> The symbol list is hardcoded in the method (`EURUSD`, `GBPUSD`). Modify as needed.

---

## ⬆️ Output

Returns a stream of **tick updates** for each subscribed symbol:

| Field    | Type        | Description            |
| -------- | ----------- | ---------------------- |
| `Symbol` | `string`    | Trading symbol name.   |
| `Bid`    | `float64`   | Current bid price.     |
| `Ask`    | `float64`   | Current ask price.     |
| `Time`   | `timestamp` | UTC time of the quote. |

---

## 🎯 Purpose

Use this method to receive real-time market data continuously for selected symbols.
Ideal for:

* Building live dashboards
* Updating price widgets in trading UIs
* Tracking bid/ask spreads and reacting to ticks

Streams can be combined with context timeouts or cancellation logic for graceful control.
