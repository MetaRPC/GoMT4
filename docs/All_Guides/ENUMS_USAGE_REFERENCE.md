# ENUMs Usage Reference — GoMT4

> Every enumeration the GoMT4 API uses, with its Go constants and values. All live in the generated package `pb "github.com/MetaRPC/GoMT4/package"`.

## How enums work here

`protoc-gen-go` generates each proto enum as a typed Go constant set, named `pb.<EnumType>_<VALUE>`. **The SDK methods take the typed enum directly** — no `int32` cast needed:

```go
import pb "github.com/MetaRPC/GoMT4/package"

// typed enum passed straight through
r, _ := acc.OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.10,
    nil, nil, nil, nil, nil, nil, nil)

h, _ := acc.OrdersHistory(ctx, pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC, &from, &to, nil, nil)
```

> To explore any enum interactively, use the [Protobuf Inspector](PROTOBUF_INSPECTOR_GUIDE.md): `enum OrderSendOperationType`, `whereenum OpenedOrderType`.

---

## `OrderSendOperationType` — order direction/type (requests)

**Used by:** `OrderSend` (the `operationType` argument). `MT4Sugar` wraps these as `BuyMarket` / `BuyLimit` / `BuyStop` / … .

| Constant | Value | Meaning |
|----------|-------|---------|
| `pb.OrderSendOperationType_OC_OP_BUY` | 0 | Market buy |
| `pb.OrderSendOperationType_OC_OP_SELL` | 1 | Market sell |
| `pb.OrderSendOperationType_OC_OP_BUYLIMIT` | 2 | Buy limit |
| `pb.OrderSendOperationType_OC_OP_SELLLIMIT` | 3 | Sell limit |
| `pb.OrderSendOperationType_OC_OP_BUYSTOP` | 4 | Buy stop |
| `pb.OrderSendOperationType_OC_OP_SELLSTOP` | 5 | Sell stop |

---

## `OpenedOrderType` — order type on open/history orders (replies)

**Used by:** `OpenedOrderInfo.GetOrderType()`, `HistoryOrderInfo` — tells market orders from pendings and buy from sell.

| Constant | Value | Meaning |
|----------|-------|---------|
| `pb.OpenedOrderType_OO_OP_BUY` | 0 | Market buy (a "position") |
| `pb.OpenedOrderType_OO_OP_SELL` | 1 | Market sell (a "position") |
| `pb.OpenedOrderType_OO_OP_BUYLIMIT` | 2 | Buy limit (pending) |
| `pb.OpenedOrderType_OO_OP_SELLLIMIT` | 3 | Sell limit (pending) |
| `pb.OpenedOrderType_OO_OP_BUYSTOP` | 4 | Buy stop (pending) |
| `pb.OpenedOrderType_OO_OP_SELLSTOP` | 5 | Sell stop (pending) |
| `pb.OpenedOrderType_OO_OP_BALANCE` | 6 | Balance entry (deposit/withdrawal) |
| `pb.OpenedOrderType_OO_OP_CREDIT` | 7 | Credit entry |

```go
for _, o := range orders {
    if o.GetOrderType() == pb.OpenedOrderType_OO_OP_BUY || o.GetOrderType() == pb.OpenedOrderType_OO_OP_SELL {
        // it's an open market position
    }
}
```

---

## `EnumOpenedOrderSortType` — open-order sort

**Used by:** the raw `OpenedOrders` request path (MT4Sugar/Service default to open-time ascending).

| Constant | Value |
|----------|-------|
| `pb.EnumOpenedOrderSortType_SORT_BY_OPEN_TIME_ASC` | 0 |
| `pb.EnumOpenedOrderSortType_SORT_BY_OPEN_TIME_DESC` | 1 |
| `pb.EnumOpenedOrderSortType_SORT_BY_ORDER_TICKET_ID_ASC` | 2 |
| `pb.EnumOpenedOrderSortType_SORT_BY_ORDER_TICKET_ID_DESC` | 3 |

---

## `EnumOrderHistorySortType` — history sort

**Used by:** `OrdersHistory` / `OrdersHistoryStream` (`sortType`).

| Constant | Value |
|----------|-------|
| `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_OPEN_TIME_ASC` | 0 |
| `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_OPEN_TIME_DESC` | 1 |
| `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_ASC` | 2 |
| `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_CLOSE_TIME_DESC` | 3 |
| `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_ORDER_TICKET_ID_ASC` | 4 |
| `pb.EnumOrderHistorySortType_HISTORY_SORT_BY_ORDER_TICKET_ID_DESC` | 5 |

---

## `ENUM_QUOTE_HISTORY_TIMEFRAME` — bar timeframe

**Used by:** `QuoteHistory` / `QuoteHistoryStream` / `MT4Sugar.Bars`.

| Constant | Value | Period |
|----------|-------|--------|
| `pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_M1` | 0 | 1 minute |
| `…_QH_PERIOD_M5` | 1 | 5 minutes |
| `…_QH_PERIOD_M15` | 2 | 15 minutes |
| `…_QH_PERIOD_M30` | 3 | 30 minutes |
| `…_QH_PERIOD_H1` | 4 | 1 hour |
| `…_QH_PERIOD_H4` | 5 | 4 hours |
| `…_QH_PERIOD_D1` | 6 | 1 day |
| `…_QH_PERIOD_W1` | 7 | 1 week |
| `…_QH_PERIOD_MN1` | 8 | 1 month |

> ⚠️ **Ordinal 0–8, not MQL minute counts.** `H1` is `4`, not `60`.

---

## `EnumAccountTradeMode` — account type

**Used by:** `AccountSummaryData.GetAccountTradeMode()`.

| Constant | Value | Meaning |
|----------|-------|---------|
| `pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_DEMO` | 0 | Demo account |
| `pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_CONTEST` | 1 | Contest account |
| `pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_REAL` | 2 | Live account |

---

## `EnumAccountMarginMode` — margin mode

**Used by:** `AccountSummaryData.GetAccountMarginMode()`.

| Constant | Value | Meaning |
|----------|-------|---------|
| `pb.EnumAccountMarginMode_ACCOUNT_MARGIN_MODE_RETAIL_NETTING` | 0 | Netting (parity with MT5) |
| `pb.EnumAccountMarginMode_ACCOUNT_MARGIN_MODE_EXCHANGE` | 1 | Exchange |
| `pb.EnumAccountMarginMode_ACCOUNT_MARGIN_MODE_RETAIL_HEDGING` | 2 | Hedging |

> **MT4 is always hedging.** The other values exist only for parity with MT5.

---

## `OrderCloseDeleteMode` — what a close/delete did (reply)

**Used by:** `OrderClose` / `OrderCloseDelete` → `OrderCloseDeleteData.GetMode()`.

| Constant | Value | Meaning |
|----------|-------|---------|
| `pb.OrderCloseDeleteMode_OCD_MARKET_ORDER` | 0 | A market order was **closed** |
| `pb.OrderCloseDeleteMode_OCD_PENDING_ORDER` | 1 | A pending order was **deleted** |

---

## `SERVER_TIMEZONE_SOURCE` — origin of the UTC offset

**Used by:** `ServerTimezoneData.GetSource()`.

| Constant | Value | Meaning |
|----------|-------|---------|
| `pb.SERVER_TIMEZONE_SOURCE_SERVER_TIMEZONE_SOURCE_UNKNOWN` | 0 | Not read yet (offset stays unset) |
| `pb.SERVER_TIMEZONE_SOURCE_SERVER_TIMEZONE_SOURCE_LIVE` | 1 | Read by this terminal |
| `pb.SERVER_TIMEZONE_SOURCE_SERVER_TIMEZONE_SOURCE_LAST_VALIDATED` | 2 | Last validated offset for this broker |

---

## `MT4_SUB_ENUM_EVENT_GROUP_TYPE` — streaming event group

**Used by:** `OnTradeData` / `OnOpenedOrdersProfitData` (the event `type`).

| Constant | Value |
|----------|-------|
| `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE_OrderProfit` | 0 |
| `pb.MT4_SUB_ENUM_EVENT_GROUP_TYPE_OrderUpdate` | 1 |

---

## `SUB_ORDER_OPERATION_TYPE` — order type in stream frames

**Used by:** the order rows inside `OnTrade` / `OnOpenedOrdersProfit` frames. Same 0–5 numbering as `OrderSendOperationType` (`SUB_OP_BUY` … `SUB_OP_SELLSTOP`).

---

## Symbol-params enums (on `SymbolParamsManyInfo`)

Returned by `SymbolParamsMany` — contract/trading rules per symbol.

**`SP_ENUM_SYMBOL_TRADE_MODE`** — `GetTradeMode()`: `…_DISABLED` 0 · `…_LONGONLY` 1 · `…_SHORTONLY` 2 · `…_CLOSEONLY` 3 · `…_FULL` 4.
**`SP_ENUM_SYMBOL_TRADE_EXECUTION`** — `GetTradeExeMode()`: `…_REQUEST` 0 · `…_INSTANT` 1 · `…_MARKET` 2 · `…_EXCHANGE` 3.
**`SP_ENUM_TRADE_CALC_MODE`** — `GetTradeCalcMode()`: `…_FOREX` 0 · `…_FUTURES` 1 · `…_CFD` 2.

---

## `ErrorType` — error classification

**Used by:** `pb.Error.GetType()` (see [Return Codes](RETURN_CODES_REFERENCE.md)). Values: `UNDEFINED` 0, `MRPC` 1, `MRPC_VALIDATION` 2, `MRPC_TIMEOUT` 3, `TERMINAL_API_TIMEOUT` 4, `MQL_EXECUTION` 8, … (10 total).

---

> **Not covered here** (not part of the trading surface): admin / GUI / EA / chart enums (`ChartTemplateMode`, `EA_PARAM_TYPE`, `GuiFormFieldType`, `ProxyTypes`, `LogType`, …) and the large `MqlErrorCode` set (see [Return Codes](RETURN_CODES_REFERENCE.md)). List **all 26** enums with the inspector: `go run ./protobuf_inspector`, then `enum <Name>` — see the [Protobuf Inspector](PROTOBUF_INSPECTOR_GUIDE.md).

**See also:** [Return Codes](RETURN_CODES_REFERENCE.md) · [MT4Account overview](../MT4Account/Overview.md)
