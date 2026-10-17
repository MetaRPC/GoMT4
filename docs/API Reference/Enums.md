# 🎛️ API Reference — Enums (GoMT4)

This page documents the **enumerations** used by the messages the client works with. Names and values match the generated Go code in `github.com/MetaRPC/GoMT4/package` (`pb`). The Go constant is `<Type>_<Name>`, e.g. `pb.OrderSendOperationType_OC_OP_BUY`; `.String()` returns the name.

---

## 📊 Orders

### `pb.OrderSendOperationType` — `OrderSend(..., operationType, ...)`

| Name           | Value | Meaning            |
| -------------- | ----: | ------------------ |
| `OC_OP_BUY`       |     0 | Market Buy         |
| `OC_OP_SELL`      |     1 | Market Sell        |
| `OC_OP_BUYLIMIT`  |     2 | Pending Buy Limit  |
| `OC_OP_SELLLIMIT` |     3 | Pending Sell Limit |
| `OC_OP_BUYSTOP`   |     4 | Pending Buy Stop   |
| `OC_OP_SELLSTOP`  |     5 | Pending Sell Stop  |

### `pb.OpenedOrderType` — `OpenedOrderInfo.OrderType`, `HistoryOrderInfo.OrderType`

| Name              | Value | Meaning            |
| ----------------- | ----: | ------------------ |
| `OO_OP_BUY`       |     0 | Buy                |
| `OO_OP_SELL`      |     1 | Sell               |
| `OO_OP_BUYLIMIT`  |     2 | Pending Buy Limit  |
| `OO_OP_SELLLIMIT` |     3 | Pending Sell Limit |
| `OO_OP_BUYSTOP`   |     4 | Pending Buy Stop   |
| `OO_OP_SELLSTOP`  |     5 | Pending Sell Stop  |

### `pb.SUB_ORDER_OPERATION_TYPE` — order type in stream events (`OnTradeOrderInfo.Type`, `OnOpenedOrdersProfitOrderInfo.Type`)

| Name               | Value | Meaning            |
| ------------------ | ----: | ------------------ |
| `SUB_OP_BUY`       |     0 | Buy                |
| `SUB_OP_SELL`      |     1 | Sell               |
| `SUB_OP_BUYLIMIT`  |     2 | Pending Buy Limit  |
| `SUB_OP_SELLLIMIT` |     3 | Pending Sell Limit |
| `SUB_OP_BUYSTOP`   |     4 | Pending Buy Stop   |
| `SUB_OP_SELLSTOP`  |     5 | Pending Sell Stop  |

### `pb.OrderCloseDeleteMode` — `OrderCloseDeleteData.Mode`

| Name                | Value | Meaning                       |
| ------------------- | ----: | ----------------------------- |
| `OCD_MARKET_ORDER`  |     0 | A market order was closed     |
| `OCD_PENDING_ORDER` |     1 | A pending order was deleted   |

---

## 📈 History & Sorting

### `pb.EnumOrderHistorySortType` — `OrdersHistory(ctx, sortType, ...)`

| Name                                  | Value | Meaning                        |
| ------------------------------------- | ----: | ------------------------------ |
| `HISTORY_SORT_BY_OPEN_TIME_ASC`       |     0 | Sort by open time ascending    |
| `HISTORY_SORT_BY_OPEN_TIME_DESC`      |     1 | Sort by open time descending   |
| `HISTORY_SORT_BY_CLOSE_TIME_ASC`      |     2 | Sort by close time ascending   |
| `HISTORY_SORT_BY_CLOSE_TIME_DESC`     |     3 | Sort by close time descending  |
| `HISTORY_SORT_BY_ORDER_TICKET_ID_ASC` |     4 | Sort by ticket ascending       |
| `HISTORY_SORT_BY_ORDER_TICKET_ID_DESC`|     5 | Sort by ticket descending      |

### `pb.EnumOpenedOrderSortType` — `OpenedOrdersRequest.SortType`

`SORT_BY_OPEN_TIME_ASC` (0), `SORT_BY_OPEN_TIME_DESC` (1), `SORT_BY_ORDER_TICKET_ID_ASC` (2), `SORT_BY_ORDER_TICKET_ID_DESC` (3). `account.OpenedOrders` sends the default (0).

---

## 💹 Quotes & Symbols

### `pb.ENUM_QUOTE_HISTORY_TIMEFRAME` — `QuoteHistory(ctx, symbol, timeframe, ...)`

| Name            | Value | Interval   |
| --------------- | ----: | ---------- |
| `QH_PERIOD_M1`  |     0 | 1 minute   |
| `QH_PERIOD_M5`  |     1 | 5 minutes  |
| `QH_PERIOD_M15` |     2 | 15 minutes |
| `QH_PERIOD_M30` |     3 | 30 minutes |
| `QH_PERIOD_H1`  |     4 | 1 hour     |
| `QH_PERIOD_H4`  |     5 | 4 hours    |
| `QH_PERIOD_D1`  |     6 | 1 day      |
| `QH_PERIOD_W1`  |     7 | 1 week     |
| `QH_PERIOD_MN1` |     8 | 1 month    |

### `pb.SP_ENUM_SYMBOL_TRADE_MODE` — `SymbolParamsManyInfo.TradeMode`

| Name                           | Value | Meaning                         |
| ------------------------------ | ----: | ------------------------------- |
| `SYMBOL_TRADE_MODE_DISABLED`   |     0 | Trading disabled                |
| `SYMBOL_TRADE_MODE_LONGONLY`   |     1 | Only long positions             |
| `SYMBOL_TRADE_MODE_SHORTONLY`  |     2 | Only short positions            |
| `SYMBOL_TRADE_MODE_CLOSEONLY`  |     3 | Only closing positions          |
| `SYMBOL_TRADE_MODE_FULL`       |     4 | No trade restrictions           |

### `pb.SP_ENUM_SYMBOL_TRADE_EXECUTION` — `SymbolParamsManyInfo.TradeExeMode`

`SYMBOL_TRADE_EXECUTION_REQUEST` (0), `SYMBOL_TRADE_EXECUTION_INSTANT` (1), `SYMBOL_TRADE_EXECUTION_MARKET` (2), `SYMBOL_TRADE_EXECUTION_EXCHANGE` (3).

### `pb.SP_ENUM_TRADE_CALC_MODE` — `SymbolParamsManyInfo.TradeCalcMode`

`SYMBOL_TRADE_MARGINE_CALC_MODE_FOREX` (0), `SYMBOL_TRADE_MARGINE_CALC_MODE_FUTURES` (1), `SYMBOL_TRADE_MARGINE_CALC_MODE_CFD` (2).

---

## 🧑‍💻 Account

### `pb.EnumAccountTradeMode` — `AccountSummaryData.AccountTradeMode`

| Name                         | Value | Meaning  |
| ---------------------------- | ----: | -------- |
| `ACCOUNT_TRADE_MODE_DEMO`    |     0 | Demo     |
| `ACCOUNT_TRADE_MODE_CONTEST` |     1 | Contest  |
| `ACCOUNT_TRADE_MODE_REAL`    |     2 | Real     |

---

## 🧮 Streams

### `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE` — `OnTradeData.Type`, `OnOpenedOrdersProfitData.Type`

| Name          | Value | Meaning              |
| ------------- | ----: | -------------------- |
| `OrderProfit` |     0 | Profit update        |
| `OrderUpdate` |     1 | Order/trade update   |

---

## ⚠️ Errors

### `pb.ErrorType` — `Error.Type`

| Name                       | Value | Meaning                                  |
| -------------------------- | ----: | ---------------------------------------- |
| `UNDEFINED`                |     0 | Not specified                            |
| `MRPC`                     |     1 | MetaRPC error                            |
| `MRPC_VALIDATION`          |     2 | Request validation failed                |
| `MRPC_TIMEOUT`             |     3 | MetaRPC timeout                          |
| `TERMINAL_API_TIMEOUT`     |     4 | Terminal did not answer in time          |
| `TERMINAL_API_GET_PARAMS`  |     5 | Terminal could not read the parameters   |
| `TERMINAL_API_SET_RESULT`  |     6 | Terminal could not return the result     |
| `TERMINAL_API_UNEXPECTED`  |     7 | Unexpected terminal error                |
| `MQL_EXECUTION`            |     8 | MT4 refused the operation (see `MqlErrorCode`) |
| `MQL_CUSTOM_EXECUTION`     |     9 | Custom MQL execution error               |

### `pb.MqlErrorCode` — `Error.MqlErrorCode`

The MT4 trade/runtime error codes, e.g. `ERR_INVALID_STOPS` (130), `ERR_INVALID_TRADE_VOLUME` (131), `ERR_MARKET_CLOSED` (132), `ERR_NOT_ENOUGH_MONEY` (134) — Go constants `pb.MqlErrorCode_ERR_INVALID_STOPS`, …. The full list is in `package/mrpc-mt4-error.pb.go`; `Error.MqlErrorIntCode` holds the same value as a number and `Error.MqlErrorDescription` its text.

---

📌 **Tip:** Enums are used in requests and responses. Always check which enum a field expects (`OrderSendOperationType` for sending, `OpenedOrderType` in order lists, `SUB_ORDER_OPERATION_TYPE` in streams). Wrong values will cause broker errors or ignored requests.
