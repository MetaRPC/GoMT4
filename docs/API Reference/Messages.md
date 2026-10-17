# 📘 API Reference — Messages (GoMT4)

This page documents the **message types** the client (`MT4Account`) sends and receives, as generated in `github.com/MetaRPC/GoMT4/package` (`pb`).
Field names and Go types match the generated code. Enums are on a separate page.

---

## 🔖 Conventions

* **`*timestamppb.Timestamp`** → ⏰ fields without suffix are broker server time; the `...Utc` sibling is the same moment in UTC.
* **Pointer scalars (`*float64`, `*int32`, `*string`, `*bool`)** → 🎛 optional fields with presence (`nil` = not set).
* **Prices & volumes** → 💹 respect `SymbolParamsManyInfo` (`Digits`, `Point`, `VolumeStep`, `VolumeMin`/`VolumeMax`).
* **Currency** → 💵 monetary fields are in **account currency** (see `AccountSummaryData.AccountCurrency`).
* **Replies** → every `...Reply` holds either `Data` (the `...Data` message below) or `Error` (`*pb.Error`); `MT4Account` methods return the `Data` part and turn `Error` into a Go error.

---

## 🧾 Account & Orders

### 📊 AccountSummaryData — `account.AccountSummary(ctx)`

| Field                       | Type                      |
| --------------------------- | ------------------------- |
| `AccountLogin`              | `int64`                   |
| `AccountBalance`            | `float64`                 |
| `AccountEquity`             | `float64`                 |
| `AccountCredit`             | `float64`                 |
| `AccountUserName`           | `string`                  |
| `AccountLeverage`           | `int64`                   |
| `AccountTradeMode`          | `pb.EnumAccountTradeMode` |
| `AccountCompanyName`        | `string`                  |
| `AccountCurrency`           | `string`                  |
| `IsInvestor`                | `bool`                    |
| `ServerTime` / `ServerTimeUtc` | `*timestamppb.Timestamp` |
| `UtcServerTimeShiftMinutes` | `int32`                   |
| `ServerUtcOffsetSeconds`    | `*int32` (nil while unknown) |

**Notes:** snapshot for health/risk checks. Margin / free margin / margin level come with the `OnOpenedOrdersProfit` and `OnTrade` streams (`AccountInfo`).

### 📌 OpenedOrderInfo — `account.OpenedOrders(ctx)` → `OpenedOrdersData.OrderInfos`

| Field            | Type                     |
| ---------------- | ------------------------ |
| `Ticket`         | `int32`                  |
| `Symbol`         | `string`                 |
| `OrderType`      | `pb.OpenedOrderType`     |
| `Lots`           | `float64`                |
| `OpenPrice`      | `float64`                |
| `StopLoss`       | `float64`                |
| `TakeProfit`     | `float64`                |
| `Profit`         | `float64`                |
| `Swap`           | `float64`                |
| `Commision`      | `float64` (spelled so in the proto) |
| `MagicNumber`    | `int32`                  |
| `Comment`        | `string`                 |
| `OpenTime` / `OpenTimeUtc` | `*timestamppb.Timestamp` |
| `ExpirationTime` / `ExpirationTimeUtc` | `*timestamppb.Timestamp` |
| `PositionIndex`, `SortIndex` | `int32`      |
| `AccountLogin`   | `int64`                  |

**Notes:** live positions and pending orders. `OpenedOrdersTicketsData` (from `account.OpenedOrdersTickets`) has only `Tickets []int32`.

### 📜 HistoryOrderInfo — `account.OrdersHistory(...)` → `OrdersHistoryData.OrdersInfo`

Same fields as `OpenedOrderInfo`, plus:

| Field                        | Type                     |
| ---------------------------- | ------------------------ |
| `ClosePrice`                 | `float64`                |
| `CloseTime` / `CloseTimeUtc` | `*timestamppb.Timestamp` |

`OrdersHistoryData` also has `TotalCount int32` and `NextPageToken string`.

### 🔎 OrdersHistoryRequest (built by `account.OrdersHistory`)

| Field                        | Type                          | Set by the client from |
| ---------------------------- | ----------------------------- | ---------------------- |
| `InputSortMode`              | `pb.EnumOrderHistorySortType` | `sortType`             |
| `InputFrom` / `InputTo`      | `*timestamppb.Timestamp`      | `from` / `to` (order **open** time window, broker server time; nil = no bound) |
| `PageNumber` / `ItemsPerPage`| `*int32`                      | `page` / `itemsPerPage` |
| `InputFromUtc` / `InputToUtc`, `PageToken` / `PageSize` | | not used by the client (UTC range, cursor paging) |

### 📨 OrderSendRequest (built by `account.OrderSend`)

| Field           | Type                        | Optional |
| --------------- | --------------------------- | :------: |
| `Symbol`        | `string`                    |          |
| `OperationType` | `pb.OrderSendOperationType` |          |
| `Volume`        | `float64`                   |          |
| `Price`         | `*float64`                  |     ✓    |
| `Slippage`      | `*int32`                    |     ✓    |
| `Stoploss`      | `*float64`                  |     ✓    |
| `Takeprofit`    | `*float64`                  |     ✓    |
| `Comment`       | `*string`                   |     ✓    |
| `MagicNumber`   | `*int32`                    |     ✓    |
| `Expiration`    | `*timestamppb.Timestamp`    |     ✓    |

* Market orders: `Price` may be omitted (current market price, with `Slippage`).
* Pending orders: set `Price` and (optionally) `Expiration` (broker server time).

### 📬 OrderSendData — result of `account.OrderSend`

| Field                        | Type                     |
| ---------------------------- | ------------------------ |
| `Ticket`                     | `int32`                  |
| `Volume`                     | `float64`                |
| `Price`                      | `float64`                |
| `OpenTime` / `OpenTimeUtc`   | `*timestamppb.Timestamp` |

**Notes:** a broker rejection arrives as `Error` in the reply; the client returns it as a Go error.

### ✏️ OrderModifyRequest (built by `account.OrderModify`)

| Field           | Type                     | Optional |
| --------------- | ------------------------ | :------: |
| `OrderTicket`   | `int32`                  |          |
| `NewPrice`      | `*float64`               |     ✓    |
| `NewStopLoss`   | `*float64`               |     ✓    |
| `NewTakeProfit` | `*float64`               |     ✓    |
| `NewExpiration` | `*timestamppb.Timestamp` |     ✓    |

Result: `OrderModifyData.OrderWasModified bool` (returned as the method's `bool`).

### ❌ OrderCloseDeleteRequest (built by `account.OrderClose` / `account.OrderDelete`)

| Field          | Type       | Optional |
| -------------- | ---------- | :------: |
| `OrderTicket`  | `int32`    |          |
| `Lots`         | `*float64` |     ✓    |
| `ClosingPrice` | `*float64` |     ✓    |
| `Slippage`     | `*int32`   |     ✓    |

Result: `OrderCloseDeleteData` — `Mode pb.OrderCloseDeleteMode` (market order closed / pending order deleted), `HistoryOrderComment *string`.

### 🔄 OrderCloseByRequest (built by `account.OrderCloseBy`)

| Field                     | Type    |
| ------------------------- | ------- |
| `TicketToClose`           | `int32` |
| `OppositeTicketClosingBy` | `int32` |

Result: `OrderCloseByData` — `ClosePrice float64`, `Profit float64`, `CloseTime` / `CloseTimeUtc *timestamppb.Timestamp`.

---

## 💹 Quotes & Market Info

### 💱 QuoteData — `account.Quote(ctx, symbol)`, `account.QuoteMany(...)` → `QuoteManyData.Quotes`

| Field                        | Type                     |
| ---------------------------- | ------------------------ |
| `Symbol`                     | `string`                 |
| `Bid`                        | `float64`                |
| `Ask`                        | `float64`                |
| `High` / `Low`               | `float64`                |
| `DateTime` / `DateTimeUtc`   | `*timestamppb.Timestamp` |

**Notes:** compute spread in points: `(Ask - Bid) / Point` (`Point` from `SymbolParams`).

### 🗂️ SymbolsData — `account.Symbols(ctx)` / `account.ShowAllSymbols(ctx)`

`SymbolNameInfos []*pb.SymbolNameInfo`, each with `SymbolName string`, `SymbolIndex int32`.

### 📊 SymbolParamsManyInfo — `account.SymbolParams(ctx, symbol)`, `account.SymbolParamsMany(...)` → `SymbolParamsManyData.SymbolInfos`

Key fields (the message has ~70; see `package/mt4-term-api-account-helper.pb.go`):

| Field               | Type                                |
| ------------------- | ----------------------------------- |
| `SymbolName`        | `string`                            |
| `SymDescription`    | `string`                            |
| `Digits`            | `int32`                             |
| `Point`             | `float64`                           |
| `VolumeMin` / `VolumeMax` / `VolumeStep` | `float64`      |
| `TradeStopsLevel`   | `int32`                             |
| `TradeFreezeLevel`  | `int32`                             |
| `TradeContractSize` | `float64`                           |
| `TradeTickValue` / `TradeTickSize` | `float64`            |
| `TradeMode`         | `pb.SP_ENUM_SYMBOL_TRADE_MODE`      |
| `TradeExeMode`      | `pb.SP_ENUM_SYMBOL_TRADE_EXECUTION` |
| `TradeCalcMode`     | `pb.SP_ENUM_TRADE_CALC_MODE`        |
| `CurrencyBase` / `CurrencyProfit` / `CurrencyMargin` | `string` |
| `Bid` / `Ask`       | `float64`                           |
| `SwapLong` / `SwapShort` | `float64`                      |
| `IsSelected` / `Visible` | `bool`                         |

**Notes:** use for rounding, SL/TP distance checks and lot validation.

### 💹 SymbolTickValueWithSizeInfo — `account.TickValueWithSize(ctx, symbols)` → `TickValueWithSizeData.Infos`

`Index int32`, `SymbolName string`, `TradeTickValue float64`, `TradeTickSize float64`, `TradeContractSize float64`.

### 🕯️ HistoryQuote — `account.QuoteHistory(...)` → `QuoteHistoryData.HistoricalQuotes`

| Field                  | Type                     |
| ---------------------- | ------------------------ |
| `Symbol`               | `string`                 |
| `Index`                | `int32`                  |
| `Time` / `TimeUtc`     | `*timestamppb.Timestamp` (bar open time) |
| `Open` / `High` / `Low` / `Close` | `float64`     |
| `TickVolume` / `RealVolume` | `int64`             |
| `Spread`               | `int32`                  |

Request (`QuoteHistoryRequest`, built by the client): `Symbol`, `Timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME`, `FromTime`, `ToTime` (broker server time).

**Notes:** the server marks `MarketInfo.QuoteHistory` deprecated; its paged successor is `QuoteHistory.Bars` (`pb.QuoteHistoryBarsRequest` / `pb.QuoteHistoryBarsData`, client `pb.NewQuoteHistoryClient`), which `MT4Account` does not wrap.

---

## 🔌 Connection

### 🔐 ConnectExRequest — `account.ConnectByServerName(ctx, serverName, baseChartSymbol, waitForTerminalIsAlive, timeoutSeconds)`

| Field            | Type      | Set by the client from |
| ---------------- | --------- | ---------------------- |
| `User`           | `uint64`  | account login          |
| `Password`       | `string`  | account password       |
| `MtClusterName`  | `string`  | `serverName` (broker server name as in MetaTrader) |
| `TimeoutSeconds` | `*uint32` | `timeoutSeconds`       |

### 🔐 ConnectRequest — `account.ConnectByHostPort(ctx, host, port, baseChartSymbol, waitForTerminalIsAlive, timeoutSeconds)`

`User`, `Password`, `Host`, `Port int32`, `TimeoutSeconds *uint32`.

Both replies carry `TerminalInstanceGuid` in their data; the client stores it in `account.TerminalInstanceGuid` and sends it as the `id` metadata on every later call. `baseChartSymbol` is only stored on the account, and `waitForTerminalIsAlive` is not used by the client; readiness is checked with an `AccountSummary` call right after connecting.

### 🆔 GetIdRequest / GetIdData — `account.GetId(ctx...)`

Request: `User string`, `Password string`. Data: `Id string` — a deterministic terminal id for these credentials. If the call fails, the client computes the same id locally (`ComputeDeterministicTerminalId`).

### 📴 DisconnectRequest — `account.Disconnect(args...)`

`Reason *string`, `Delete *bool` (the client sets `Delete` when you pass `true`).

---

## ⚠️ Errors

### Error (`*pb.Error`)

| Field                 | Type               | Meaning                                  |
| --------------------- | ------------------ | ---------------------------------------- |
| `Type`                | `pb.ErrorType`     | Where the error comes from (MRPC, validation, timeout, MQL execution, …). |
| `ErrorCode`           | `string`           | Code, e.g. `TERMINAL_INSTANCE_NOT_FOUND`. |
| `ErrorMessage`        | `string`           | Human-readable message.                  |
| `MqlErrorCode`        | `pb.MqlErrorCode`  | MT4 error code (when MT4 refused).        |
| `MqlErrorIntCode`     | `int32`            | Same as a number.                        |
| `MqlErrorDescription` | `string`           | MT4 error text.                          |
| `StackTrace`, `CommandTypeName`, `CommandId`, `Properties` | | Diagnostics. |

The client returns it as `API error (code=<ErrorCode>): <ErrorMessage>`.

**See also:** the **Enums** page for enum value tables.
