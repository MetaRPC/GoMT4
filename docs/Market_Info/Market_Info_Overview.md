# Market Info — Overview

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
// Discover and pick symbols (svc.ShowAllSymbols(ctx) prints the same list)
syms, err := account.ShowAllSymbols(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("symbols:", len(syms.GetSymbolNameInfos()))

// Print an initial snapshot for the shortlist
svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD"})

// (optional) Subscribe to ticks in the Streaming section
// svc.StreamQuotes(ctx)
```

### 2) Validate trading inputs and display instrument info

```go
// Fetch parameters for validation and formatting (prints them; returns an error on failure)
if err := svc.ShowSymbolParams(ctx, "EURUSD"); err != nil {
    log.Printf("symbol params: %v", err)
}

// Print monetary values (tick value / tick size / contract size)
svc.ShowTickValues(ctx, []string{"EURUSD"})
```

### 3) Chart a symbol

```go
// Print the last 5 days of H1 candles
svc.ShowQuoteHistory(ctx, "EURUSD")

// ...or get the candles to render them yourself
to := time.Now()
bars, err := account.QuoteHistory(ctx, "EURUSD", pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, to.AddDate(0, 0, -5), to)
if err == nil {
    fmt.Println("bars:", len(bars.GetHistoricalQuotes()))
}
```

The `Show*` methods of `MT4Service` print their results and return nothing (`ShowSymbolParams` returns only an `error`); to work with the data, call the `MT4Account` methods (`ShowAllSymbols`, `Symbols`, `SymbolParams`, `SymbolParamsMany`, `Quote`, `QuoteMany`, `QuoteHistory`, `TickValueWithSize`).

---

## ✅ Best Practices

1. **Key by name, not index.** `SymbolIndex` can change across sessions; persist `SymbolName`.
2. **Format by Digits.** Use `Digits` from `SymbolParams` for UI; keep raw doubles for math.
3. **Batch when possible.** Prefer `QuoteMany` and `TickValueWithSize` to reduce round‑trips.
4. **Know the time zone.** Times like `DateTime` / `Time` are **broker server time**; each has a UTC sibling (`DateTimeUtc`, `TimeUtc`) for the same moment. History ranges you pass (`from`, `to`) are read as broker server time.
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
