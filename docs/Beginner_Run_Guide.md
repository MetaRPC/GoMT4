# 🚦 Beginner Run Guide for GoMT4

This guide is based on `examples/main.go`. Code is already wired; just enter your credentials and API key in `examples/config/config.json` and keep (or comment out) the calls you want to run.

---

## ⚡ How to Run

```bash
git clone https://github.com/MetaRPC/GoMT4.git
cd GoMT4/examples
go run .
```

The entry point is `examples/main.go`. It connects to the MetaRPC cloud API (`mt4.mrpc.pro:443`), which runs the MT4 terminal for your account; nothing else has to be installed or started. You can also pass the API key on the command line: `go run . --api-key=YOUR_KEY`.

---

## 🧪 Safe First Steps (read‑only)

Start with these — these **do not change** account state:

```go
// Account snapshot
svc.ShowAccountSummary(ctx)

// Discover all instruments
svc.ShowAllSymbols(ctx)

// One‑shot quote for your default symbol
svc.ShowQuote(ctx, cfg.DefaultSymbol)

// Live quotes for predefined symbols (stops by timeout in example)
svc.StreamQuotes(ctx)
```

---

## 📊 Getting Data (account & market)

Handy readers to inspect the environment:

```go
// Recent closed orders (default window inside helper)
svc.ShowOrdersHistory(ctx)

// All active orders (incl. pending)
svc.ShowOpenedOrders(ctx)

// Full instrument profile
svc.ShowSymbolParams(ctx, cfg.DefaultSymbol)

// Monetary metrics for sizing
svc.ShowTickValues(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})

// Multi‑symbol snapshot quotes
svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})

// Historical OHLC for charting
svc.ShowQuoteHistory(ctx, cfg.DefaultSymbol)
```

---

## ⚠️ Trading Operations (danger zone)

These **modify state** (even on demo). Use real tickets from output of previous steps.

```go
// Place an order (market BUY, 0.1 lot, on the given symbol)
svc.ShowOrderSendExample(ctx, cfg.DefaultSymbol)

// Modify SL/TP — requires a valid ticket
svc.ShowOrderModifyExample(ctx, 12345678)

// Close order by ticket
svc.ShowOrderCloseExample(ctx, 12345678)

// Close by opposite order (two tickets)
svc.ShowOrderCloseByExample(ctx, 12345678, 12345679)

// Delete a pending order by ticket
svc.ShowOrderDeleteExample(ctx, 111222)
```

> Replace placeholders with actual tickets printed by `ShowOpenedOrders` or `ShowOpenedOrderTickets`.

---

## 📡 Streaming

Real‑time subscriptions with graceful cancellation inside examples:

```go
// Live price updates (ticks)
svc.StreamQuotes(ctx)

// Floating P/L per opened order (interval inside helper)
svc.StreamOpenedOrderProfits(ctx)

// Ticket IDs of open orders only
svc.StreamOpenedOrderTickets(ctx)

// Trade activity events
svc.StreamTradeUpdates(ctx)
```

---

## 🎬 Combo Scenarios

### A) Read‑only dashboard (safe)

```go
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, cfg.DefaultSymbol)
svc.ShowOpenedOrders(ctx)
svc.StreamQuotes(ctx)
svc.StreamOpenedOrderProfits(ctx)
```

### B) Adjust risk → then close

```go
// 1) Inspect active orders and pick a ticket
svc.ShowOpenedOrders(ctx)

// 2) Update SL/TP for that ticket
svc.ShowOrderModifyExample(ctx, /*ticket=*/ 12345678)

// 3) Close when ready
svc.ShowOrderCloseExample(ctx, /*ticket=*/ 12345678)
```

### C) Tickets stream + lazy details

```go
// Low‑overhead ticket stream
svc.StreamOpenedOrderTickets(ctx)

// Fetch full details on change
svc.ShowOpenedOrders(ctx)
```

### D) History snapshot (safe)

```go
svc.ShowOrdersHistory(ctx)
```

---

## 🧠 Tips

* Begin with **safe readers** before trading helpers.
* Never hardcode tickets; copy from console output.
* Protect `examples/config/config.json` — it contains your MT4 password and API key.
* Prefer demo until you are confident.
* For price formatting, use `Digits` from `ShowSymbolParams`; keep raw values for math.
