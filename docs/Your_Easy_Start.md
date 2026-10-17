# 🛡️ Safe-by-Default Examples & Feature Toggles

> Run the examples **safely by default**. Turn on the powerful (but potentially dangerous) trading calls only when you explicitly decide to.

---

## 🪪 Requirements

* A clone of the repository; the example app is `examples/main.go` and is run from the `examples` folder.
* `examples/config/config.json` filled with **Login**, **Password**, **Server**, **DefaultSymbol** and **ApiKey** (with `"Login": 0` the example opens a fresh demo account for you; with `"ApiKey": "TRIAL"` or no key the calls run as a trial).
* Standard build/run works: `go build ./...` or `go run .` (inside `examples`).
* `main.go` creates the `MT4Account`, connects it with `ConnectByServerName`, and then uses only the `MT4Service` methods (`svc.Show...`, `svc.Stream...`).

---

## Quick Start (Safe Mode)

1. Fill `examples/config/config.json`.
2. Run the example app (`main.go`) as-is — only **read-only** calls and safe streams will run.
3. When you’re ready to try trading (open/modify/close orders), enable it via **Feature Toggles** below.

---

## ✨ Feature Toggles (Safe by Default)

`main.go` has these three flags right after connecting the account:

```go
// ===== Feature toggles (safe by default) =====
enableTradingExamples := false   // ⚠️ Real trading calls (OrderSend/Modify/Close...)
enableStreams         := true    // Read-only streaming: quotes/profits/tickets/trades
enableHistoryStreams  := true    // New: history streams (pagination/time-chunks)
```

### 🔧 What each toggle does

* **`enableTradingExamples`**

  * `false` (recommended default): trading methods are **not** executed.
  * `true`: allows examples that **open/modify/close** orders.
* **`enableStreams`**

  * Enables read‑only streaming examples (quotes, opened‑orders profit/tickets, trade updates).
* **`enableHistoryStreams`**

  * Enables **new** read‑only history streaming examples:

    * paginated order history (last 30 days),
    * chunked quote history (last 90 days, weekly chunks).

---

## Example Layout in `main.go`

Keep `main` readable and declarative. The toggles control the “dangerous” parts.

```go
// --- 📂 Account Info ---
svc.ShowAccountSummary(ctx)

// --- 📂 Order Operations (read-only) ---
svc.ShowOpenedOrders(ctx)
svc.ShowOpenedOrderTickets(ctx)
svc.ShowOrdersHistory(ctx)

// --- ⚠️ Trading (dangerous): turn on consciously ---
if enableTradingExamples {
    // Will open a real order (demo or real — depends on your login)
    svc.ShowOrderSendExample(ctx, cfg.DefaultSymbol)

    // Require valid tickets:
    svc.ShowOrderModifyExample(ctx, 12345678)
    svc.ShowOrderCloseExample(ctx, 12345678)
    svc.ShowOrderCloseByExample(ctx, 12345678, 12345679)

    // Deletes ONLY pending orders
    svc.ShowOrderDeleteExample(ctx, 12345678)
}

// --- 📂 Market Info / Symbol Info ---
svc.ShowQuote(ctx, cfg.DefaultSymbol)
svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})
svc.ShowQuoteHistory(ctx, cfg.DefaultSymbol)
svc.ShowAllSymbols(ctx)
svc.ShowSymbolParams(ctx, cfg.DefaultSymbol)
svc.ShowTickValues(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})

// --- 📂 Streaming / Subscriptions (read-only) ---
if enableStreams {
    svc.StreamQuotes(ctx)
    svc.StreamOpenedOrderProfits(ctx)
    svc.StreamOpenedOrderTickets(ctx)
    svc.StreamTradeUpdates(ctx)
}

// --- 🧾 New history streams (read-only) ---
if enableHistoryStreams {
    svc.StreamOrdersHistoryExample(ctx)                      // last 30 days, paginated
    svc.StreamQuoteHistoryExample(ctx, cfg.DefaultSymbol)   // last 90 days, chunked
}
```

> Tip: If you still have an old commented "demo" block at the bottom of `main.go`, remove it. The list above is the single source of truth now.

---

## What’s “dangerous” vs “safe”

**Dangerous (disabled by default):**

* `ShowOrderSendExample` — opens a market BUY of 0.1 lot on the given symbol.
* `ShowOrderModifyExample` — changes SL/TP.
* `ShowOrderCloseExample` — closes an order by ticket.
* `ShowOrderCloseByExample` — closes using an opposite order.
* `ShowOrderDeleteExample` — deletes **only** pending orders.

**Safe (read‑only):**

* All other `Show*` info calls and `Stream*` subscriptions (quotes, profits, tickets, trades, history).

---

## How to Use Feature Toggles (Step by Step)

1. Open `main.go`.
2. Find the **Feature toggles** block:

   ```go
   enableTradingExamples := false
   enableStreams := true
   enableHistoryStreams := true
   ```
3. To run **only safe** examples → keep `enableTradingExamples = false`.
4. When you’re ready to try trading:

   * Set `enableTradingExamples = true`.
   * Put **real, valid tickets** into modify/close/close‑by/delete calls.
   * Prefer **DEMO** first.
5. Run normally from the `examples` folder (`go run .` or `go build && ./examples`).

---

## 🏗️ Safety Checklist (before enabling trading)

* ✅ You’re on a **DEMO** account (strongly recommended for first runs).
* ✅ `ShowAccountSummary` looks sane (it prints balance, equity and currency).
* ✅ Symbol, volume, SL/TP comply with your broker’s rules.
* ✅ Tickets are **real**, e.g., from `ShowOpenedOrders`.

---

## 🤖 FAQ

**Do I need to call `MT4Account` directly?**
Only to create and connect it (`main.go` already does that). After that you can use `MT4Service` methods (`svc.Show...`, `svc.Stream...`), which call the `MT4Account` methods for you; call `MT4Account` directly when you need the returned data instead of printed output.

**Where can I change date ranges for history?**

* `StreamOrdersHistoryExample` — last 30 days by default.
* `StreamQuoteHistoryExample` — last 90 days (weekly chunks).
  For custom ranges call `svc.StreamOrdersHistory(ctx, sortType, &from, &to, pageSize)` or `svc.StreamQuoteHistory(ctx, symbol, timeframe, from, to, chunk)` directly; they take the range as parameters.

**How do I prevent accidental trades in a team repo?**
Keep `enableTradingExamples = false` in committed code, add a big comment banner above the toggle. Flip it locally for tests.

---

## Troubleshooting

* *“Nothing happens in trading examples”* → You likely still have `enableTradingExamples = false`.
* *“Order modify/close fails”* → The ticket is invalid or not applicable (e.g., trying to delete a market order).
* *“Streams stop too soon”* → Example wrappers include a safety timeout. Adjust/remove `time.After(...)` in wrapper methods if you need long-running streams.

---

**That’s it — clean `main`, safe defaults, and powerful features behind a single switch. When you’re ready, flip the toggle and proceed step by step.**
