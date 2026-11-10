# 3. Symbols & Quotes — Overview

> Everything about instruments and prices: the symbol list, Market Watch select, bulk contract params, latest ticks, bar history, and the server timezone.

Mostly read-only (except `SymbolSelect`). On `mt4_term_api.MarketInfo`, `AccountHelper` (params / tick value) and `ServerTimezone`.

---

## Methods in this group

| Method | Returns | One-liner |
|--------|---------|-----------|
| [`Quote(ctx, symbol)`](./Quote.md) | `*pb.QuoteData` | Latest tick (bid/ask + session high/low). |
| [`QuoteMany(ctx, symbols)`](./QuoteMany.md) | `*pb.QuoteManyData` | Latest ticks for several symbols. |
| [`QuoteHistory(ctx, symbol, tf, from, to)`](./QuoteHistory.md) | `*pb.QuoteHistoryData` | OHLC bars for a symbol/timeframe range. |
| [`Symbols(ctx)`](./Symbols.md) | `*pb.SymbolsData` | Every symbol name the terminal knows. |
| [`SymbolSelect(ctx, symbol, selected)`](./SymbolSelect.md) | `bool` | Add / remove a symbol in Market Watch. |
| [`SymbolParams(ctx, symbol)`](./SymbolParams.md) | `*pb.SymbolParamsManyInfo` | Full contract params for one symbol. |
| [`SymbolParamsMany(ctx, symbols)`](./SymbolParamsMany.md) | `*pb.SymbolParamsManyData` | Contract params for several symbols. |
| [`TickValueWithSize(ctx, symbols)`](./TickValueWithSize.md) | `*pb.TickValueWithSizeData` | Tick value / size / contract size. |
| [`ServerTimezone(ctx)`](./ServerTimezone.md) | `*pb.ServerTimezoneData` | The broker server's UTC offset. |

---

**Back to:** [MT4Account — Overview](../Overview.md)
