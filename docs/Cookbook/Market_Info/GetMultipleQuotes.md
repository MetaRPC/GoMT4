# 📚 GetMultipleQuotes (GoMT4)

**Goal:** fetch quotes for **many symbols** efficiently.

> Real code refs:
>
> * Batch quote: `examples/mt4/MT4Account.go` (`QuoteMany`)
> * Single quote: `examples/mt4/MT4Account.go` (`Quote`)
> * Demo: `examples/mt4/MT4_service.go` (`ShowQuotesMany`)

---

## ✅ 1) Preconditions

* `account` is connected; all symbols exist on the broker (exact names).
* `config.json` has at least one known symbol to start with.

---

## 🚀 2) Fast path A — Batch RPC

`QuoteMany` gets all quotes in one round‑trip:

```go
symbols := []string{"EURUSD", "GBPUSD", "USDJPY"}
qs, err := account.QuoteMany(ctx, symbols)
if err != nil { return err }
for _, q := range qs.GetQuotes() {
    fmt.Printf("%s Bid=%.5f Ask=%.5f Time=%s\n",
        q.GetSymbol(), q.GetBid(), q.GetAsk(), q.GetDateTime().AsTime().Format(time.RFC3339))
}
```

> An empty `symbols` slice returns `no symbols provided` without calling the server.

---

## 🧵 3) Fast path B — Fan‑out with goroutines

If you need per‑symbol error handling (one bad symbol should not fail the whole batch), run `Quote` requests in parallel and collect results:

```go
func GetQuotesParallel(ctx context.Context, a *mt4.MT4Account, syms []string) ([]*pb.QuoteData, error) {
    type item struct { q *pb.QuoteData; err error }
    out := make([]*pb.QuoteData, 0, len(syms))
    ch  := make(chan item, len(syms))
    var wg sync.WaitGroup

    // limit concurrency to avoid hammering the API
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
    fmt.Printf("%s %.5f/%.5f @ %s\n", q.GetSymbol(), q.GetBid(), q.GetAsk(), q.GetDateTime().AsTime().Format(time.RFC3339))
}
```

---

## 🧪 4) Validation & fallbacks

* If a symbol fails (wrong name or suffix), log it and continue.
* Optionally, map missing symbols to their broker‑specific names (`EURUSD` → `EURUSD.m`).

```go
for _, s := range syms {
    q, err := account.Quote(ctx, s)
    if err != nil {
        log.Printf("skip %s: %v", s, err)
        continue
    }
    _ = q // ...
}
```

---

## ⚠️ Pitfalls

* **Storming the API**: keep concurrency reasonable (4–8) to avoid timeouts; prefer `QuoteMany`.
* **Suffixes**: `EURUSD` vs `EURUSD.m` — take the exact names from `account.Symbols(ctx)`.
* **Timeouts**: each call has a short default timeout (3s); the client retries transient errors itself.

---

## 🔗 See also

* [GetQuote](GetQuote.md) — single symbol quote.
* [StreamQuotes](StreamQuotes.md) — live stream for continuous updates.
* [SymbolParams](SymbolParams.md) — get `Digits`, `Point`, `VolumeStep` for rounding & display.
