# 🚀 Performance Notes (GoMT4)

Practical tips to keep GoMT4 **fast, stable, and resource‑efficient**. Focus is on *your* current codebase (`examples/`, the `package/` pb API, and the MT4 terminal that MetaRPC runs for you behind `mt4.mrpc.pro:443`).

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

Keep these tight: avoid heavy logging, allocations, and repeated RPCs. Every call is a network round‑trip to the cloud API, so the number of calls matters more than local CPU.

---

## 📦 Batch, don’t loop RPCs

* Prefer **batch** APIs over per‑symbol calls:

  * ✅ `account.QuoteMany(ctx, symbols)` ([`ShowQuotesMany`](Market_Info/ShowQuotesMany.md)) instead of `account.Quote` ([`ShowQuote`](Market_Info/ShowQuote.md)) in a loop.
  * ✅ `account.TickValueWithSize(ctx, symbols)` takes a list.
* History: use **paged**/chunked requests rather than huge single responses.

  * ✅ `account.OrdersHistory` with `page` / `itemsPerPage`, or `account.OrdersHistoryStream` ([`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md)).

**Why**: fewer round‑trips → lower latency and CPU per message.

---

## 📡 Streaming tuning

* **Server‑side interval**: the opened‑orders streams take an interval in milliseconds — `account.OnOpenedOrdersProfit(ctx, intervalMs)` and `account.OnOpenedOrdersTickets(ctx, intervalMs)`. Use the largest interval your logic tolerates.
* **Consumer loop**: drain fast, do work in a goroutine pool if processing is heavy. The data channels are unbuffered: a slow reader slows the stream.
* **Tick throttling**: if you don’t need *every* tick, aggregate with a `time.Ticker` (e.g., 100–200ms) before updating UI/logs.
* **Timeouts**: long‑lived streams should use `WithCancel`, not short timeouts; add an **idle watchdog** if you want to act when no data arrives for X seconds.

```go
ctx, cancel := context.WithCancel(parent)
defer cancel()
updates, errCh := account.OnSymbolTick(ctx, symbols)
idle := time.NewTimer(30 * time.Second)
for {
    select {
    case u, ok := <-updates:
        if !ok {
            return
        }
        _ = u // cheap handling
        idle.Reset(30 * time.Second)
    case err := <-errCh:
        log.Printf("stream ended: %v", err) // the helper already retried; reopen if needed
        return
    case <-idle.C:
        // no data for 30s (closed market?) → log, check with account.Quote, or cancel and reopen
    }
}
```

---

## 🧵 Concurrency patterns (Go)

* **One `MT4Account`** (one `*grpc.ClientConn`) per MT4 account; reuse it from all goroutines (creating many connections is expensive).
* **Worker pool** for CPU‑bound post‑processing (e.g., statistics on history), *not* for issuing parallel trade RPCs for the same account.
* **Channel fan‑out**: if several consumers need the same stream, broadcast from **one** reader goroutine; don’t open multiple identical streams.

---

## 🧰 Symbol metadata cache

* Cache `account.SymbolParams` ([`ShowSymbolParams`](Market_Info/ShowSymbolParams.md)) and `account.TickValueWithSize` ([`ShowTickValues`](Market_Info/ShowTickValues.md)) in memory (map by symbol). These values change rarely; re‑fetch only on **symbol list** changes or once per session.
* Use cached `Digits` / `Point` / `VolumeStep` for **rounding** before any trade call (see Cookbook → [Round Volume/Price](Cookbook/Utils_Helpers/RoundVolumePrice.md)).

---

## 🧮 History I/O

* Request time windows in **chunks** (days/weeks) and merge client‑side (`account.QuoteHistoryStream` does the chunking for you).
* Prefer **H1/H4/D1** where possible; M1 over long ranges is expensive.
* For analytics, stream and **process as you go** instead of storing everything first.

---

## 🗜️ Payloads

* Avoid logging full payloads in hot paths; log only *counts* and *latency* metrics.

---

## 🧠 Allocations & GC

* Reuse buffers/slices for quote aggregation; pre‑allocate with `make(..., cap)` when size is known.
* Convert `time` ↔ `Timestamp` carefully in tight loops; avoid repeated parsing/formatting.
* Use `strings.Builder` for building CSV/TSV outputs.

---

## ⏱️ Timeouts & retries

* **Unary** calls: `context.WithTimeout` (1–5s typical for a cloud call; the methods default to 3s/5s/8s when your `ctx` has no deadline). Retries with **exponential backoff** are built in (`ExecuteWithReconnect`); add your own only for idempotent operations.
* **Streams**: no short deadline; the client reconnects on transient errors with jittered backoff (see Cookbook → [Handle Reconnect](Cookbook/Reliability_Connection/HandleReconnect.md), [Unary Retries](Cookbook/Reliability_Connection/UnaryRetries.md)).

---

## 🖥️ MT4 specifics

* **Symbol visibility**: ensure symbols are visible in Market Watch before quote/history calls (Cookbook → [Ensure Symbol Visible](Cookbook/Utils_Helpers/EnsureSymbolVisible.md)). Subscribing early with `OnSymbolTick` reduces first‑tick latency.

---

## 📊 Observability

* Add lightweight metrics: per‑RPC latency, stream message rate, reconnect counters.
* Sampled logs (1/N) for hot paths; full debug logs only on demand.

---

## ✅ Checklist (TL;DR)

* [ ] Use **batch** requests and **paged/chunked** history.
* [ ] Keep **one** `MT4Account` (gRPC conn) per account; reuse it.
* [ ] Cache **SymbolParams** / **TickValues**.
* [ ] Fast stream consumers; apply throttling/aggregation if UI‑facing.
* [ ] Timeouts for unary; **cancel** for streams; exponential **backoff** with jitter (built in).
* [ ] Minimal logging in hot paths; metrics for rates/latency.
* [ ] Pre‑make symbols visible or subscribe early.
* [ ] Prefer higher‑TF history when possible; chunk long ranges.

---

## 🗺️ Code map → where each tip lives in your repo

### 1) Batch, don’t loop RPCs

* **Batch quotes** → `examples/mt4/MT4Account.go` → `func (a *MT4Account) QuoteMany(ctx context.Context, symbols []string)`; printed by `MT4Service.ShowQuotesMany` in `examples/mt4/MT4_service.go`
* **Single quote (avoid loops)** → `MT4Account.Quote(ctx, symbol)`; `MT4Service.ShowQuote`
* **History (paged/chunked)** → `MT4Account.OrdersHistoryStream`, `MT4Account.QuoteHistoryStream`; demos [`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md), [`StreamQuoteHistoryExample`](Streaming/StreamQuoteHistoryExample.md)

### 2) Streaming tuning

* **Quotes stream** → `examples/mt4/MT4_service.go` → [`StreamQuotes`](Streaming/StreamQuotes.md) (uses `OnSymbolTick`)
* **Opened order profits** → same file → [`StreamOpenedOrderProfits`](Streaming/StreamOpenedOrderProfits.md) (`OnOpenedOrdersProfit(ctx, 1000)`, stops after 3s without data)
* **Trade updates** → same file → [`StreamTradeUpdates`](Streaming/StreamTradeUpdates.md)
* **Opened tickets** → same file → [`StreamOpenedOrderTickets`](Streaming/StreamOpenedOrderTickets.md)

### 3) One gRPC connection, reuse clients

* **Connection bootstrap** → `examples/main.go` → `mt4.NewMT4Account(...)`, `ConnectByServerName(...)`, `mt4.NewMT4Service(account)`
* **Account/session holder** → `examples/mt4/MT4Account.go`
  *Search:* `type MT4Account struct` / `grpc.Dial(` / `ConnectByServerName`

### 4) Symbol metadata cache (params/tick-values)

* **Fetch params** → `MT4Account.SymbolParams` / `SymbolParamsMany`; printed by [`ShowSymbolParams`](Market_Info/ShowSymbolParams.md)
* **Fetch tick values** → `MT4Account.TickValueWithSize`; printed by [`ShowTickValues`](Market_Info/ShowTickValues.md)

### 5) Rounding & validation before trades

* The examples (`ShowOrderSendExample`, `ShowOrderModifyExample`) use fixed demo values and do **not** round; add the helpers from [Round Volume/Price](Cookbook/Utils_Helpers/RoundVolumePrice.md) in your own code.
* **Stops distance** → `TradeStopsLevel` / `TradeFreezeLevel` from `SymbolParams`.

### 6) Ensure symbol visible / pre‑subscribe

* **Ensure visible** → the generated `MarketInfoClient.SymbolSelect` RPC (no wrapper in `MT4Account`), see [Ensure Symbol Visible](Cookbook/Utils_Helpers/EnsureSymbolVisible.md)
* **Pre‑subscribe quotes** (lower first‑tick latency) → `MT4Account.OnSymbolTick`

### 7) History I/O in chunks

* **Quote history (time windows)** → `MT4Account.QuoteHistory` / `QuoteHistoryStream`; [`ShowQuoteHistory`](Market_Info/ShowQuoteHistory.md)
* **Orders history (paged)** → `MT4Account.OrdersHistoryStream`; [`StreamOrdersHistoryExample`](Streaming/StreamOrdersHistoryExample.md)

### 8) Timeouts & retries

* **Default per‑call timeouts** → every unary method in `examples/mt4/MT4Account.go` (`context.WithTimeout` when `ctx` has no deadline).
* **Reconnect with backoff** → `ExecuteWithReconnect` / `ExecuteStreamWithReconnect` in `MT4Account.go`; Cookbook pages [Handle Reconnect](Cookbook/Reliability_Connection/HandleReconnect.md) and [Unary Retries](Cookbook/Reliability_Connection/UnaryRetries.md).

### 9) Observability (logs/metrics)

* **Light logs around hot paths** → `examples/mt4/MT4_service.go` stream handlers print `Tick/Profit/Trade` lines — replace with counters/rate meters in production.
