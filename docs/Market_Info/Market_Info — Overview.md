# Market\_Info — Overview

This section contains **market data and instrument metadata** methods for MT4. Use it to discover symbols, read quotes, fetch historical candles, and get contract/tick parameters.

---

## 📂 Methods in this folder

* [ShowAllSymbols.md](ShowAllSymbols.md)
  Full catalogue of all instruments available in the terminal (names + indices).

* [ShowSymbols.md](ShowSymbols.md)
  Lightweight list of symbol names and indices (basic variant).

* [ShowSymbolParams.md](ShowSymbolParams.md)
  Extended parameters for a symbol: digits, volume rules, currencies, trade mode.

* [ShowQuote.md](ShowQuote.md)
  Latest **bid/ask** snapshot for a single symbol.

* [ShowQuotesMany.md](ShowQuotesMany.md)
  Snapshot quotes for **multiple** symbols at once; useful before subscribing to ticks.

* [ShowQuoteHistory.md](ShowQuoteHistory.md)
  Historical **OHLC** data (candles) for a symbol over a timeframe.

* [ShowTickValues.md](ShowTickValues.md)
  **TickValue / TickSize / ContractSize** for one or more symbols (P/L & sizing math).

---

## ⚡ Typical Workflows

### 1) Build a watchlist and show live prices

```go
// Discover and pick symbols
syms, _ := svc.ShowAllSymbols(ctx)

// Get initial snapshot for the shortlist
_ = svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD"})

// (optional) Subscribe to ticks in the Streaming section
// svc.StreamQuotes(ctx)
```

### 2) Validate trading inputs and display instrument info

```go
// Fetch parameters for validation and formatting
_ = svc.ShowSymbolParams(ctx, "EURUSD")

// Compute monetary values
_ = svc.ShowTickValues(ctx, []string{"EURUSD"})
```

### 3) Chart a symbol

```go
// Pull historical candles and render
_ = svc.ShowQuoteHistory(ctx, "EURUSD")
```

---

## ✅ Best Practices

1. **Key by name, not index.** `SymbolIndex` can change across sessions; persist `SymbolName`.
2. **Format by Digits.** Use `Digits` from `ShowSymbolParams` for UI; keep raw doubles for math.
3. **Batch when possible.** Prefer `ShowQuotesMany` and `ShowTickValues` to reduce round‑trips.
4. **Time is UTC.** Quotes and candles come with UTC timestamps—convert only at presentation.
5. **Broker suffixes are real.** Treat `EURUSD.m` / `XAUUSD-RAW` as distinct symbols.

---

## 🎯 Purpose

The `Market_Info` block helps you:

* Discover and organize instruments.
* Read **current** and **historical** market data.
* Validate and format symbol‑specific values (digits, steps, currencies).
* Compute monetary effects (tick value/size, contract size) for risk & P/L.

---

👉 Use this page as a **map**. Jump into each `.md` file for full method details, parameters, and pitfalls.
