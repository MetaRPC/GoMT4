# 📚 GetMultipleQuotes (GoMT4)

**Goal:** fetch quotes for **many symbols** efficiently.

> Real code refs:
>
> * Single quote: `examples/mt4/MT4Account.go` (`Quote`)
> * Demo: `examples/mt4/MT4_service.go` (see `ShowQuotesMany` flow)

---

## ✅ 1) Preconditions

* All symbols exist and are visible in MT4 (*Market Watch* → Show All).
* `config.json` has at least one known symbol to start with.

---

## 🚀 2) Fast path A — Batch RPC (if available in your build)

Some builds expose a **batch** RPC (e.g., `QuotesMany(ctx, symbols)`), mirroring the `ShowQuotesMany` example. If you have it, use it:

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
qs, err := account.QuotesMany(ctx, symbols)
if err != nil { return err }
for _, q := range qs.GetQuotes() {
    fmt.Printf("%s Bid=%.5f Ask=%.5f Spread=%.1f pips\n",
        q.GetSymbol(), q.GetBid(), q.GetAsk(), (q.GetAsk()-q.GetBid())/q.GetPoint())
}
```

> If your current API does not expose this convenience call, use **fan‑out** concurrency (next section).

---

## 🧵 3) Fast path B — Fan‑out with goroutines (always works)

Run `Quote` requests in parallel and collect results.

```go
func GetQuotesParallel(ctx context.Context, a *MT4Account, syms []string) ([]*pb.Quote, error) {
    type item struct { q *pb.Quote; err error }
    out := make([]*pb.Quote, 0, len(syms))
    ch  := make(chan item, len(syms))
    var wg sync.WaitGroup

    // limit concurrency to avoid hammering the terminal
    sem := make(chan struct{}, 8) // 8 concurrent calls

    for _, s := range syms {
        s := s
        wg.Add(1)
        go func() {
            defer wg.Done()
            sem <- struct{}{}
            defer func(){ <-sem }()

            q, err := a.Quote(ctx, s)
            ch <- item{q, err}
        }()
    }
    wg.Wait()
    close(ch)

    for it := range ch {
        if it.err != nil { return nil, it.err }
        out = append(out, it.q)
    }
    return out, nil
}
```

**Usage:**

```go
qs, err := GetQuotesParallel(ctx, account, []string{"EURUSD","GBPUSD","USDJPY"})
if err != nil { log.Fatal(err) }
for _, q := range qs {
    fmt.Printf("%s %.5f/%.5f @ %s\n", q.GetSymbol(), q.GetBid(), q.GetAsk(), q.GetTime().AsTime().Format(time.RFC3339))
}
```

---

## 🧪 4) Validation & fallbacks

* If a symbol fails (hidden or has a suffix), log it and continue.
* Optionally, map missing symbols to their broker‑specific names (`EURUSD` → `EURUSD.m`).

```go
for _, s := range syms {
    q, err := a.Quote(ctx, s)
    if err != nil {
        log.Printf("skip %s: %v", s, err)
        continue
    }
    // ...
}
```

---

## ⚠️ Pitfalls

* **Storming MT4**: keep concurrency reasonable (4–8) to avoid timeouts.
* **Suffixes**: `EURUSD` vs `EURUSD.m` — normalize names or take from *Market Watch*.
* **Timeouts**: wrap each call with a short timeout (2–3s) and retry only **transport** errors.

---

## 🔗 See also

* `GetQuote.md` — single symbol quote.
* `StreamQuotes.md` — live stream for continuous updates.
* `SymbolParams.md` — get `Digits`, `Point`, `LotStep` for rounding & display.
