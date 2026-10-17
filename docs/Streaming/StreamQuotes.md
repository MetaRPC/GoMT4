# 🧊 Streaming Real-Time Quotes

> **Request:** subscribe to live tick updates for predefined symbols
> Continuously streams the latest bid/ask prices and timestamps for each symbol.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Streams live ticks for EURUSD and GBPUSD; the demo stops after 3 ticks or 5s without a tick.
svc.StreamQuotes(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

symbols := []string{"EURUSD", "GBPUSD"}

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

tickCh, errCh := account.OnSymbolTick(ctx, symbols)

fmt.Println("🔄 Streaming ticks...")
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

    case <-time.After(30 * time.Second): // demo timeout
        fmt.Println("⏱️ Timeout reached.")
        return
    }
}

```

---

### Method Signature

```go
func (s *MT4Service) StreamQuotes(ctx context.Context)
func (a *MT4Account) OnSymbolTick(ctx context.Context, symbols []string) (<-chan *pb.OnSymbolTickData, <-chan error)
```

---

## 🔽 Input

| Field | Type              | Description                                |
| ----- | ----------------- | ------------------------------------------ |
| `ctx` | `context.Context` | Controls stream lifetime and cancellation. |

> The wrapper uses a fixed symbol list (`EURUSD`, `GBPUSD`). For your own symbols call `account.OnSymbolTick(ctx, symbols)`.

---

## ⬆️ Output

Stream of `*pb.OnSymbolTickData` packets. Each packet may contain `SymbolTick` with:

| Field     | Type                     | Description                                   |
| --------- | ------------------------ | --------------------------------------------- |
| `Symbol`  | `string`                 | Trading symbol name.                          |
| `Bid`     | `float64`                | Current bid price.                            |
| `Ask`     | `float64`                | Current ask price.                            |
| `Last`    | `float64`                | Price of the last deal.                       |
| `Volume`  | `uint64`                 | Volume for the last price.                    |
| `Time`    | `*timestamppb.Timestamp` | Time of the tick (broker server time).        |
| `TimeUtc` | `*timestamppb.Timestamp` | The same moment in UTC.                       |
| `TimeMsc` / `TimeMscUtc` | `int64`   | Tick time in milliseconds (broker time / UTC). |

The packet also carries `TerminalInstanceGuidId` (the terminal that sent it).

---

## 🎯 Purpose

Receive continuous **real-time market data** for selected symbols — ideal for live dashboards, widgets, and spread tracking.

---

## 🧩 Notes & Tips

* **Per-symbol cache:** Keep a `map[string]*pb.OnSymbolMqlTickInfo` of last values and update only on change to reduce UI churn.
* **Both channels:** Always consume **data and error** channels to avoid leaks.
* **Display precision:** Use symbol `Digits` (from `SymbolParams`) for formatting; keep raw doubles for math.

---

## ⚠️ Pitfalls

* **Nil checks:** `SymbolTick` can be nil in a packet — guard before reading fields (as in example).
* **Ordering:** Do not assume packets are strictly chronological across symbols.
* **Bursts:** Rapid bursts can overwhelm rendering; debounce or batch prints.

