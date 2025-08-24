# 🚀 Performance Notes (GoMT4)

Practical tips to keep GoMT4 **fast, stable, and resource‑efficient**. Focus is on *your* current codebase (`examples/`, pb API, MT4 terminal on Windows).

---

## 🎯 Goals

* **Low latency** for quotes/updates.
* **Throughput** for bulk reads (history, many symbols).
* **Stability** under loss/reconnects.
* **Predictable CPU/RAM** usage.

---

## 🔥 Hot paths to watch

* **Streaming**: [`StreamQuotes`](../Streaming/StreamQuotes.md), [`StreamOpenedOrderProfits`](../Streaming/StreamAccountProfits.md), [`StreamTradeUpdates`](../Streaming/StreamTradeUpdates.md).
* **Batch RPCs**: [`ShowQuotesMany`](../Market_Info/ShowQuotesMany.md), paged/streamed **Orders History**.
* **Symbol metadata**: [`ShowSymbolParams`](../Market_Info/ShowSymbolParams.md), [`ShowTickValues`](../Market_Info/ShowTickValues.md).

➡️ In your code: see `examples/mt4/MT4_service.go` (methods `ShowQuotesMany`, `StreamQuotes`, etc.) — these are the hot loops.

---

## 📦 Batch, don’t loop RPCs

* Prefer **batch** APIs over per‑symbol calls:

  * ✅ `ShowQuotesMany([]string)` instead of `ShowQuote()` in a loop.
* History: use **paged**/streaming endpoints rather than huge single responses.

  * ✅ `OrdersHistoryPagedRequest{page_size: N}` or `OrdersHistoryStream{...}`.

➡️ Example: `examples/mt4/MT4_service.go: ShowQuotesMany` demonstrates a batched call.

---

## 📡 Streaming tuning

* **Backpressure**: tune `OpenedOrdersProfitStreamRequest{ buffer_size }` (see `MT4_service.go: StreamOpenedOrderProfits`).
* **Consumer loop**: drain fast, do work in a goroutine pool if processing is heavy.
* **Tick throttling**: if you don’t need *every* tick, aggregate before updating UI/logs.
* **Timeouts**: long‑lived streams → use `WithCancel`. Add **idle watchdog** for reconnect.

➡️ Example: `examples/mt4/MT4_service.go: StreamOpenedOrderProfits` — already has timeout (`time.After`) and error handling.

---

## 🧵 Concurrency patterns (Go)

* **One gRPC connection** per MT4 instance; reuse `*grpc.ClientConn` (created in `examples/mt4/MT4_service.go: NewMT4Service`).
* **Worker pool** for CPU‑bound post‑processing (statistics/analytics), *not* for parallel trade RPCs to the same terminal.
* **Channel fan‑out**: if multiple consumers need the same stream, broadcast from one reader; don’t duplicate streams.

---

## 🧰 Symbol metadata cache

* Cache `SymbolParams` and `TickValues` in memory (map by symbol).
* Fetch once at startup (see `examples/mt4/MT4_service.go: ShowSymbolParams`).
* Use cached `digits/point/lot_step` for rounding before any order (Cookbook → `RoundVolumePrice.md`).

---

## 🧮 History I/O

* Request in **chunks** (days/weeks) and merge client‑side.
* Prefer higher TF (`H1/H4/D1`) for long ranges.
* For analytics: process streamed chunks on the fly.

➡️ Example: `examples/mt4/MT4_service.go: ShowOrdersHistory` streams history with `from/to` ranges.

---

## 🗜️ Compression & payloads

* gRPC **compression** helps on large history, but avoid on ticks.
* Avoid logging full payloads in hot paths (quotes). Log counts/latency only.

➡️ In your code: `main.go` prints full quotes for demo. In production → reduce log volume.

---

## 🧠 Allocations & GC

* Reuse slices for quotes aggregation.
* Convert `time ↔ Timestamp` carefully in loops.
* Use `strings.Builder` for big text exports.

➡️ Example: see `examples/mt4/MT4_service.go: StreamQuotes` — every tick prints. For production, consider slice reuse.

---

## ⏱️ Timeouts & retries

* **Unary** calls: `context.WithTimeout` (100ms–3s typical). Add exponential backoff.
* **Streams**: no short deadline. Reconnect on error/idle with jitter.

➡️ Cookbook refs: `HandleReconnect.md`, `UnaryRetries.md`.

---

## 🖥️ Windows & MT4 specifics

* Keep MT4 responsive: exclude data dir from heavy antivirus scans.
* Power plan: High performance → avoid throttling.
* Ensure symbols visible before requests (`examples/mt4/MT4_service.go: EnsureSymbolVisible`).

---

## 📊 Observability

* Add lightweight metrics: latency, stream msg rate, reconnect counters.
* Sampled logs (1/N) in hot paths.

---

## ✅ Checklist (TL;DR)

* [ ] Use batch requests + paged/streamed history (`ShowQuotesMany`, `ShowOrdersHistory`).
* [ ] Keep one gRPC conn (`NewMT4Service`).
* [ ] Cache symbol params (`ShowSymbolParams`).
* [ ] Fast stream consumers (see `StreamQuotes`, `StreamOpenedOrderProfits`).
* [ ] Unary: timeouts. Streams: cancel + backoff.
* [ ] Reduce logging in hot paths (quotes).
* [ ] Ensure symbols visible early (`EnsureSymbolVisible`).
* [ ] History: chunk requests, prefer higher TF.

With these, your GoMT4 will stay lean and responsive 🚀

---

## 🗺️ Code map → where each tip lives in your repo

### 1) Batch, don’t loop RPCs

* **Batch quotes** → `examples/mt4/MT4_service.go` → `func (s *MT4Service) ShowQuotesMany(ctx context.Context, symbols []string)`
  *Search:* `ShowQuotesMany(`
* **Single quote (avoid loops)** → same file → `ShowQuote(ctx, symbol string)`
  *Search:* `ShowQuote(ctx`
* **History (paged/stream)** → `examples/mt4/MT4_service.go` → `StreamOrdersHistoryExample`, `StreamQuoteHistoryExample`
  *Search:* `StreamOrdersHistoryExample` / `StreamQuoteHistoryExample`

### 2) Streaming tuning & backpressure

* **Quotes stream** → `examples/mt4/MT4_service.go` → `StreamQuotes(ctx context.Context)`
  *Search:* `StreamQuotes(`
* **Opened order profits** → `examples/mt4/MT4_service.go` → `StreamOpenedOrderProfits(ctx context.Context)`
  *Search:* `StreamOpenedOrderProfits(`
  *(this is the method you pasted earlier with `OnOpenedOrdersProfit(ctx, 1000)` and the 30s timeout)*
* **Trade updates** → same file → `StreamTradeUpdates(ctx context.Context)`
  *Search:* `StreamTradeUpdates(`
* **Opened tickets** → same file → `StreamOpenedOrderTickets(ctx context.Context)`
  *Search:* `StreamOpenedOrderTickets(`

**Tip → buffer\_size**: if you expose `buffer_size` for profits, it is wired at account layer:

* `examples/mt4/MT4Account.go` → look for `OnOpenedOrdersProfit(ctx,`
  *Search:* `OnOpenedOrdersProfit(`

### 3) One gRPC connection, reuse clients

* **Connection bootstrap** → `examples/main.go` → creation of service & clients
  *Search:* `NewMT4Service(` / `grpc.Dial(`
* **Account/session holder** → `examples/mt4/MT4Account.go`
  *Search:* `type MT4Account struct` / `connect` / `login`

### 4) Symbol metadata cache (params/tick-values)

* **Fetch params** → `examples/mt4/MT4_service.go` → `ShowSymbolParams(ctx, symbol string)`
  *Search:* `ShowSymbolParams(`
* **Fetch tick values** → same file → `ShowTickValues(ctx context.Context, symbols []string)`
  *Search:* `ShowTickValues(`
* **Used for rounding** → `examples/mt4/MT4_service.go` → see order helpers (e.g., `ShowOrderSendExample`) where `Digits`, `Point`, `LotStep` are applied before send.
  *Search:* `ShowOrderSendExample(` / `Round` / `lot`

### 5) Rounding & validation before trades

* **Volume/price rounding** (from Cookbook) is backed by these call sites:
  `examples/mt4/MT4_service.go` → `ShowOrderSendExample`, `ShowOrderModifyExample`
  *Search:* `ShowOrderModifyExample(` / `Normalize` / `math.Round`
* **Ensure stops distance** → check usage of `stops_level` / `freeze_level` from `SymbolParams`.
  *Search:* `stops_level` / `freeze_level`

### 6) Ensure symbol visible / pre‑subscribe

* **Ensure visible** → `examples/mt4/MT4Account.go`
  *Search:* `EnsureSymbolVisible` / `SymbolSelect` / `SymbolsGet`
* **Pre‑subscribe quotes** (lower first‑tick latency) → `examples/mt4/MT4_service.go` → `StreamQuotes` setup
  *Search:* `QuoteStreamRequest` / `symbols:`

### 7) History I/O in chunks

* **Quote history (time windows)** → `examples/mt4/MT4_service.go` → `ShowQuoteHistory(ctx, symbol string)`
  *Search:* `ShowQuoteHistory(`
* **Orders history (paged)** → same file → `StreamOrdersHistoryExample`
  *Search:* `OrdersHistoryPaged`

### 8) Timeouts & retries

* **Unary with timeout** → `examples/mt4/MT4_service.go` → places that use `context.WithTimeout(...)` around quote/symbol calls.
  *Search:* `context.WithTimeout(`
* **Reconnect with backoff** → reliability helpers you documented (Cookbook):
  `Cookbook/Reliability/HandleReconnect.md` and `Cookbook/Reliability/UnaryRetries.md` correspond to the code in `examples/mt4/MT4_service.go` (stream loops) and in connection handling in `MT4Account.go`.

### 9) Observability (logs/metrics)

* **Light logs around hot paths** → `examples/mt4/MT4_service.go` stream handlers print `Tick/Profit/Trade` lines — replace with counters/rate meters in production.
  *Search:* `fmt.Println("[Tick]"` / `"[Profit]"` / `"[Trade]"`

