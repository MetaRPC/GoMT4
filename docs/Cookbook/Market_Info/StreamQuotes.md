# 📡 StreamQuotes (GoMT4)

**Goal:** subscribe to **live quotes** for one or many symbols with auto‑reconnect and clean cancellation.

> Real code refs:
>
> * Streaming helpers: `examples/mt4/MT4Account.go` (e.g., `OnSymbolTick` / stream wrapper)
> * Demo usage: `examples/mt4/MT4_service.go` (see `StreamQuotes` flow)

---

## ✅ 1) Preconditions

* MT4 terminal is connected and symbols are visible in *Market Watch*.
* You have a parent `ctx` to control lifetime.

---

## ▶️ 2) Subscribe to symbols

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

symbols := []string{"EURUSD", "GBPUSD"}

// Helper starts a stream and returns two channels:
//   dataCh: ticks
//   errCh : terminal errors (non‑recoverable)
dataCh, errCh := account.OnSymbolTick(ctx, symbols)

for {
    select {
    case <-ctx.Done():
        return // stop gracefully
    case err := <-errCh:
        if err != nil {
            log.Printf("stream stopped: %v", err)
            return
        }
    case t := <-dataCh:
        fmt.Printf("%s Bid=%.5f Ask=%.5f @ %s\n",
            t.GetSymbol(), t.GetBid(), t.GetAsk(), t.GetTime().AsTime().Format(time.RFC3339))
    }
}
```

> Under the hood the helper auto‑reconnects on `io.EOF` / `codes.Unavailable` using backoff+jitter from `MT4Account.go`.

---

## 🧭 3) Backoff & health checks

* Reconnects use exponential backoff with jitter (see constants in `MT4Account.go`).
* A short health‑check (e.g., `AccountSummary` with 3s timeout) can be run after initial connect to ensure MT4 is ready.

---

## 🧰 4) Processing tips

* Keep the loop **non‑blocking**: if you do heavy work (DB, strategy), hand off ticks into a buffered channel/worker.
* Log spread in pips: `(Ask-Bid)/Point` to catch broker anomalies.
* If you need OHLC bars, aggregate incoming ticks by timeframe on your side.

---

## ⚠️ Pitfalls

* **No reads** → if you stop reading `dataCh`, producer back‑pressure can stall the stream.
* **Hidden symbol** → ensure symbol is visible; suffixes like `EURUSD.m` are different instruments.
* **Context canceled** → stream ends; always watch `<-ctx.Done()>`.
* **Network flaps** → tune `backoffBase`/`backoffMax` (Reliability chapter).

---

## 🔄 Variations

* **Single symbol**: pass `[]string{"EURUSD"}`.
* **Dynamic subscribe**: keep a registry and restart the stream with a new symbol set when needed.
* **Parallel consumers**: fan‑out ticks to multiple goroutines via a broadcast channel.

---

## 📎 See also

* `GetQuote.md` — one‑shot quote.
* `GetMultipleQuotes.md` — batch fetch.
* `Reliability (en)` — timeouts, reconnects, backoff.
