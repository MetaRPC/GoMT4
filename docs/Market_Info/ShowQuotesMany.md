# Getting Quotes for Multiple Symbols with Ticks

> **Request:** fetch quotes for multiple symbols and stream real-time price ticks
> Combines a one-time quote snapshot (`QuoteMany`) with a live tick stream (`OnSymbolTick`) for each symbol.

---

### Code Example

```go
// Using service wrapper
symbols := []string{"EURUSD", "GBPUSD"}
service.ShowQuotesMany(context.Background(), symbols)

// Or directly from MT4Account (QuoteMany)
data, err := mt4.QuoteMany(context.Background(), symbols)
if err != nil {
    log.Fatalf("Error fetching multiple quotes: %v", err)
}

for _, q := range data.GetQuotes() {
    fmt.Printf("Symbol: %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
        q.GetSymbol(),
        q.GetBid(),
        q.GetAsk(),
        q.GetDateTime().AsTime().Format("2006-01-02 15:04:05"),
    )
}

// Stream real-time tick updates
streamSymbols := []string{"EURUSD"}
tickCh, errCh := mt4.OnSymbolTick(context.Background(), streamSymbols)
for {
    select {
    case tick := <-tickCh:
        if tick != nil && tick.GetSymbolTick() != nil {
            q := tick.GetSymbolTick()
            fmt.Printf("[Tick] %s | Bid: %.5f | Ask: %.5f | Time: %s\n",
                q.GetSymbol(),
                q.GetBid(),
                q.GetAsk(),
                q.GetTime().AsTime().Format("2006-01-02 15:04:05"))
            break // For test/demo
        }
    case err := <-errCh:
        log.Fatalf("Tick stream error: %v", err)
    case <-time.After(5 * time.Second):
        fmt.Println("Timeout reached")
        return
    }
}
```

---

## 🔽 Input

### For `QuoteMany`

| Field     | Type              | Description                              |
| --------- | ----------------- | ---------------------------------------- |
| `symbols` | `[]string`        | List of trading symbols to fetch quotes. |
| `ctx`     | `context.Context` | Cancellation and timeout control.        |

### For `OnSymbolTick`

| Field     | Type              | Description                               |
| --------- | ----------------- | ----------------------------------------- |
| `symbols` | `[]string`        | Symbols to subscribe for real-time ticks. |
| `ctx`     | `context.Context` | Required for stream control.              |

---

### Method Signatures

```go
func (s *MT4Service) ShowQuotesMany(ctx context.Context, symbols []string)

func (a *MT4Account) QuoteMany(ctx context.Context, symbols []string) (*pb.QuoteManyReply, error)

func (a *MT4Account) OnSymbolTick(ctx context.Context, symbols []string) (<-chan *pb.SymbolTickDataPacket, <-chan error)
```

---

## ⬆️ Output

### From `QuoteMany`

Returns slice of `*pb.QuoteData`:

| Field    | Type        | Description           |
| -------- | ----------- | --------------------- |
| `Symbol` | `string`    | Trading symbol name   |
| `Bid`    | `float64`   | Current bid price     |
| `Ask`    | `float64`   | Current ask price     |
| `Time`   | `timestamp` | UTC time of the quote |

### From `OnSymbolTick`

Returns real-time stream of `*pb.SymbolTickDataPacket`:

| Field        | Type        | Description                    |
| ------------ | ----------- | ------------------------------ |
| `SymbolTick` | `QuoteData` | Real-time tick data for symbol |

---

## 🎯 Purpose

Use this method when working with **multiple symbols**:

1. `QuoteMany` gives instant snapshot of bid/ask prices — good for validation or display.
2. `OnSymbolTick` streams live updates — ideal for dashboards or pricing alerts.

Perfect for trading UIs, price monitors, or auto-trading logic with symbol watchlists.

---

## 🧩 Notes & Tips

* **Stream lifecycle:** Always cancel context to stop `OnSymbolTick`. Otherwise the channel stays open.
* **Initial state:** Combine `QuoteMany` first, then subscribe to ticks — avoids waiting for first tick.
* **Channel handling:** Both data and error channels must be consumed to avoid goroutine leaks.

---

## ⚠️ Pitfalls

* **High-frequency symbols:** Subscribing to many active symbols can flood your app with ticks. Consider rate-limiting.
* **Connection drops:** Streams end silently if connection dies — always monitor `errCh`.
* **Duplicate updates:** Rapid ticks may repeat Bid/Ask without visible change. Filter if necessary.

---

## 🧪 Testing Suggestions

* **Snapshot check:** Call `QuoteMany` on 2–3 symbols, verify consistent prices.
* **Stream test:** Subscribe to one symbol, ensure ticks arrive within expected intervals.
* **Timeout/Cancel:** Cancel context mid-stream → channels should close cleanly.
