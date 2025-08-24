# 🚀 Performance Notes (GoMT4)

Practical tips to keep GoMT4 **fast, stable, and resource‑efficient**. Focus is on *your* current codebase (examples/, pb API, MT4 terminal on Windows).

---

## 🎯 Goals

* **Low latency** for quotes/updates.
* **Throughput** for bulk reads (history, many symbols).
* **Stability** under loss/reconnects.
* **Predictable CPU/RAM** usage.

---

## 🔥 Hot paths to watch

* **Streaming**: [`StreamQuotes`](Streaming/StreamQuotes.md), [`StreamOpenedOrderProfits`](Streaming/StreamOpenedOrderProfits.md), [`StreamTradeUpdates`](Streaming/StreamTradeUpdates.md).
* **Batch RPCs**: [`ShowQuotesMany`](Market_Info/ShowQuotesMany.md), paged/streamed **Orders History** ([`ShowOrdersHistory`](Order_Operations/ShowOrdersHistory.md), [`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md)).
* **Symbol metadata**: [`ShowSymbolParams`](Market_Info/ShowSymbolParams.md), [`ShowTickValues`](Market_Info/ShowTickValues.md) used in loops.

Keep these tight: avoid heavy logging, allocations, and repeated RPCs.

---

## 📦 Batch, don’t loop RPCs

* Prefer **batch** APIs over per‑symbol calls:

  * ✅ [`ShowQuotesMany`](Market_Info/ShowQuotesMany.md) instead of [`ShowQuote`](Market_Info/ShowQuote.md) in a loop.
* History: use **paged**/streaming endpoints rather than huge single responses.

  * ✅ [`OrdersHistory`](Order_Operations/ShowOrdersHistory.md) or [`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md).

**Why**: fewer round‑trips → lower latency and CPU per message.

---

## 📡 Streaming tuning

* **Backpressure**: use the server’s knobs when available.

  * Example: `OpenedOrdersProfitStreamRequest{ buffer_size }` — set buffer high enough to absorb short spikes, but not so high that consumer lags unnoticed.
* **Consumer loop**: drain fast, do work in a goroutine pool if processing is heavy.
* **Tick throttling**: if you don’t need *every* tick, aggregate with a `time.Ticker` (e.g., 100–200ms) before updating UI/logs.
* **Timeouts**: long‑lived streams should use `WithCancel`, not short timeouts; add an **idle watchdog** to reconnect when no data for X seconds.

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()
updates, errCh := svc.StreamQuotes(ctx, symbols)
for {
    select {
    case u := <-updates: /* cheap handling */
    case err := <-errCh: /* reconnect + backoff */
    case <-idle.C:      /* no data → ping/refresh */
    }
}
```

---

## 🧵 Concurrency patterns (Go)

* **One gRPC connection** per MT4 instance; reuse `*grpc.ClientConn` (creating many is expensive).
* **Worker pool** for CPU‑bound post‑processing (e.g., statistics on history), *not* for issuing parallel trade RPCs to the same MT4 terminal.
* **Channel fan‑out**: if several consumers need the same stream, broadcast from **one** reader goroutine; don’t open multiple identical streams.

---

## 🧰 Symbol metadata cache

* Cache [`ShowSymbolParams`](Market_Info/ShowSymbolParams.md) and [`ShowTickValues`](Market_Info/ShowTickValues.md) in memory (map by symbol). These values change rarely; re‑fetch only on **symbol list** changes or once per session.
* Use cached `digits/point/lot_step` for **rounding** before any trade call (see Cookbook → `RoundVolumePrice.md`).

---

## 🧮 History I/O

* Request time windows in **chunks** (days/weeks) and merge client‑side.
* Prefer **H1/H4/D1** where possible; M1 over long ranges is expensive.
* For analytics, stream and **process as you go** instead of storing everything first.

---

## 🗜️ Compression & payloads

* gRPC **compression** helps on large history responses; avoid compressing ultra‑short ticks.
* Avoid logging full payloads in hot paths; log only *counts* and *latency* metrics.

---

## 🧠 Allocations & GC

* Reuse buffers/slices for quote aggregation; pre‑allocate with `make(..., cap)` when size is known.
* Convert `time` ↔ `Timestamp` carefully in tight loops; avoid repeated parsing/formatting.
* Use `strings.Builder` for building CSV/TSV outputs.

---

## ⏱️ Timeouts & retries

* **Unary** calls: `context.WithTimeout` (100ms–3s typical for local terminal). Add **idempotency** where safe and **exponential backoff**.
* **Streams**: no short deadline; reconnect on error/idle with jittered backoff (see Cookbook → `HandleReconnect.md`, `UnaryRetries.md`).

---

## 🖥️ Windows & MT4 specifics

* **Keep MT4 responsive**: exclude MT4 data directory from real‑time antivirus scanning if safe in your environment.
* **Power plan**: use *High performance* to avoid CPU sleep throttling for low‑latency streams.
* **Symbol visibility**: ensure symbols are visible before quote/history calls (Cookbook → `EnsureSymbolVisible.md`). Pre‑subscribe to reduce first‑tick latency.

---

## 📊 Observability

* Add lightweight metrics: per‑RPC latency, stream message rate, reconnect counters.
* Sampled logs (1/N) for hot paths; full debug logs only on demand.

---

## ✅ Checklist (TL;DR)

* [ ] Use **batch** requests and **paged/streamed** history.
* [ ] Keep **one** gRPC conn; reuse clients.
* [ ] Cache **SymbolParams** / **TickValues**.
* [ ] Fast stream consumers; apply throttling/aggregation if UI‑facing.
* [ ] Timeouts for unary; **cancel** for streams; exponential **backoff** with jitter.
* [ ] Minimal logging in hot paths; metrics for rates/latency.
* [ ] Pre‑make symbols visible or subscribe early.
* [ ] Prefer higher‑TF history when possible; chunk long ranges.

---

## 🗺️ Code map → where each tip lives in your repo

### 1) Batch, don’t loop RPCs

* **Batch quotes** → `examples/mt4/MT4_service.go` → `func (s *MT4Service) ShowQuotesMany(ctx context.Context, symbols []string)`
  *Search:* `ShowQuotesMany(`
* **Single quote (avoid loops)** → same file → `ShowQuote(ctx, symbol string)`
  *Search:* `ShowQuote(ctx`
* **History (paged/stream)** → `examples/mt4/MT4_service.go` → [`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md), [`StreamQuoteHistoryExample`](Streaming/StreamQuoteHistoryExample.md)

### 2) Streaming tuning & backpressure

* **Quotes stream** → `examples/mt4/MT4_service.go` → [`StreamQuotes`](Streaming/StreamQuotes.md)
* **Opened order profits** → same file → [`StreamOpenedOrderProfits`](Streaming/StreamOpenedOrderProfits.md)
  *(this is the method you pasted earlier with `OnOpenedOrdersProfit(ctx, 1000)` and the 30s timeout)*
* **Trade updates** → same file → [`StreamTradeUpdates`](Streaming/StreamTradeUpdates.md)
* **Opened tickets** → same file → [`StreamOpenedOrderTickets`](Streaming/StreamOpenedOrderTickets.md)

**Tip → buffer\_size**: if you expose `buffer_size` for profits, it is wired at account layer:

* `examples/mt4/MT4Account.go` → look for `OnOpenedOrdersProfit(ctx,`

### 3) One gRPC connection, reuse clients

* **Connection bootstrap** → `examples/main.go` → creation of service & clients
  *Search:* `NewMT4Service(` / `grpc.Dial(`
* **Account/session holder** → `examples/mt4/MT4Account.go`
  *Search:* `type MT4Account struct` / `connect` / `login`

### 4) Symbol metadata cache (params/tick-values)

* **Fetch params** → `examples/mt4/MT4_service.go` → [`ShowSymbolParams`](Market_Info/ShowSymbolParams.md)
* **Fetch tick values** → same file → [`ShowTickValues`](Market_Info/ShowTickValues.md)
* **Used for rounding** → see order helpers (e.g., [`ShowOrderSendExample`](Order_Operations/ShowOrderSendExample.md)) where `Digits`, `Point`, `LotStep` are applied before send.

### 5) Rounding & validation before trades

* **Volume/price rounding** (from Cookbook) is backed by these call sites:
  `examples/mt4/MT4_service.go` → [`ShowOrderSendExample`](Order_Operations/ShowOrderSendExample.md), [`ShowOrderModifyExample`](Order_Operations/ShowOrderModifyExample.md)
* **Ensure stops distance** → check usage of `stops_level` / `freeze_level` from [`ShowSymbolParams`](Market_Info/ShowSymbolParams.md).

### 6) Ensure symbol visible / pre‑subscribe

* **Ensure visible** → `examples/mt4/MT4Account.go`
  *Search:* `EnsureSymbolVisible` / `SymbolSelect` / `SymbolsGet`
* **Pre‑subscribe quotes** (lower first‑tick latency) → `examples/mt4/MT4_service.go` → [`StreamQuotes`](Streaming/StreamQuotes.md)

### 7) History I/O in chunks

* **Quote history (time windows)** → `examples/mt4/MT4_service.go` → [`ShowQuoteHistory`](Market_Info/ShowQuoteHistory.md)
* **Orders history (paged)** → same file → [`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md)

### 8) Timeouts & retries

* **Unary with timeout** → `examples/mt4/MT4_service.go` → places that use `context.WithTimeout(...)` around quote/symbol calls.
* **Reconnect with backoff** → reliability helpers you documented (Cookbook):
  `Cookbook/Reliability/HandleReconnect.md` and `Cookbook/Reliability/UnaryRetries.md` correspond to the code in `examples/mt4/MT4_service.go` (stream loops) and in connection handling in `MT4Account.go`.

### 9) Observability (logs/metrics)

* **Light logs around hot paths** → `examples/mt4/MT4_service.go` stream handlers print `Tick/Profit/Trade` lines — replace with counters/rate meters in production.

---

If you want, я добавлю **точные вытяжки кода** под каждый пункт (короткие блоки `▶ snippet`) прямо в этот файл. Скажи, какие пункты показать первыми — например: `StreamOpenedOrderProfits`, `ShowQuotesMany`, `ShowOrderSendExample`.
