# 🧊 Streaming Real-Time Quotes

> **Request:** subscribe to live tick updates for predefined symbols
> Continuously streams the latest bid/ask prices and timestamps for each symbol.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Streams live ticks for default symbols (EURUSD, GBPUSD in demo); stops after ~30s.
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
```

---

## 🔽 Input

| Field | Type              | Description                                |
| ----- | ----------------- | ------------------------------------------ |
| `ctx` | `context.Context` | Controls stream lifetime and cancellation. |

> The wrapper uses a predefined symbol list (e.g., `EURUSD`, `GBPUSD`). Adjust as needed.

---

## ⬆️ Output

Stream of `*pb.OnSymbolTickData` packets. Each packet may contain `SymbolTick` with:

| Field    | Type        | Description            |
| -------- | ----------- | ---------------------- |
| `Symbol` | `string`    | Trading symbol name.   |
| `Bid`    | `float64`   | Current bid price.     |
| `Ask`    | `float64`   | Current ask price.     |
| `Time`   | `timestamp` | UTC time of the quote. |

---

## 🎯 Purpose

Receive continuous **real-time market data** for selected symbols — ideal for live dashboards, widgets, and spread tracking.

---

## 🧩 Notes & Tips

* **Per-symbol cache:** Keep a `map[string]Quote]` of last values and update only on change to reduce UI churn.
* **Both channels:** Always consume **data and error** channels to avoid leaks.
* **Display precision:** Use symbol `Digits` (from `SymbolParams`) for formatting; keep raw doubles for math.

---

## ⚠️ Pitfalls

* **Nil checks:** `SymbolTick` can be nil in a packet — guard before reading fields (as in example).
* **Ordering:** Do not assume packets are strictly chronological across symbols.
* **Bursts:** Rapid bursts can overwhelm rendering; debounce or batch prints.

