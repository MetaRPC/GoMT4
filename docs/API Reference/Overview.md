# API Reference (Types & Enums)

Generated from your provided **.proto** files. Field order and enum values exactly match the source.

---

## mrpc-mt4-error.proto

**Messages (3):**

### MrpcError

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `code`    | `int32`  |    no    |
|  2 | `message` | `string` |    no    |
|  3 | `details` | `string` |    no    |

### OrderError

|  # | Field    | Type     | Repeated |
| -: | -------- | -------- | :------: |
|  1 | `ticket` | `int64`  |    no    |
|  2 | `reason` | `int32`  |    no    |
|  3 | `text`   | `string` |    no    |

### BatchOrderError

|  # | Field    | Type        | Repeated |
| -: | -------- | ----------- | :------: |
|  1 | `errors` | `MrpcError` |    yes   |

**Enums (3):**

### MrpcErrorCode

| Name                          | Value |
| ----------------------------- | ----: |
| `MRPC_ERROR_CODE_NONE`        |     0 |
| `MRPC_ERROR_CODE_UNKNOWN`     |     1 |
| `MRPC_ERROR_CODE_TIMEOUT`     |     2 |
| `MRPC_ERROR_CODE_CONNECTION`  |     3 |
| `MRPC_ERROR_CODE_UNAVAILABLE` |     4 |

### OrderErrorReason

| Name                                    | Value |
| --------------------------------------- | ----: |
| `MT4_ORDER_ERROR_REASON_NONE`           |     0 |
| `MT4_ORDER_ERROR_REASON_INVALID_VOLUME` |     1 |
| `MT4_ORDER_ERROR_REASON_INVALID_PRICE`  |     2 |
| `MT4_ORDER_ERROR_REASON_MARKET_CLOSED`  |     3 |
| `MT4_ORDER_ERROR_REASON_SERVER_BUSY`    |     4 |

### StreamErrorReason

| Name                                  | Value |
| ------------------------------------- | ----: |
| `MT4_STREAM_ERROR_REASON_NONE`        |     0 |
| `MT4_STREAM_ERROR_REASON_EOF`         |     1 |
| `MT4_STREAM_ERROR_REASON_UNAVAILABLE` |     2 |

---

## mt4-term-api-account-helper.proto

**Messages (22):**

### LoginRequest

|  # | Field      | Type     | Repeated |
| -: | ---------- | -------- | :------: |
|  1 | `login`    | `int64`  |    no    |
|  2 | `password` | `string` |    no    |
|  3 | `server`   | `string` |    no    |

### AccountSummary

|  # | Field          | Type     | Repeated |
| -: | -------------- | -------- | :------: |
|  1 | `login`        | `int64`  |    no    |
|  2 | `name`         | `string` |    no    |
|  3 | `server`       | `string` |    no    |
|  4 | `currency`     | `string` |    no    |
|  5 | `balance`      | `double` |    no    |
|  6 | `equity`       | `double` |    no    |
|  7 | `margin`       | `double` |    no    |
|  8 | `margin_free`  | `double` |    no    |
|  9 | `margin_level` | `double` |    no    |
| 10 | `leverage`     | `int32`  |    no    |
| 11 | `profit`       | `double` |    no    |

### OpenedOrder

|  # | Field        | Type                        | Repeated |
| -: | ------------ | --------------------------- | :------: |
|  1 | `ticket`     | `int64`                     |    no    |
|  2 | `symbol`     | `string`                    |    no    |
|  3 | `type`       | `Mt4OrderType`              |    no    |
|  4 | `volume`     | `double`                    |    no    |
|  5 | `open_price` | `double`                    |    no    |
|  6 | `sl`         | `double`                    |    no    |
|  7 | `tp`         | `double`                    |    no    |
|  8 | `open_time`  | `google.protobuf.Timestamp` |    no    |
|  9 | `commission` | `double`                    |    no    |
| 10 | `swap`       | `double`                    |    no    |
| 11 | `magic`      | `int32`                     |    no    |
| 12 | `comment`    | `string`                    |    no    |

### OrderHistory

|  # | Field         | Type                        | Repeated |
| -: | ------------- | --------------------------- | :------: |
|  1 | `ticket`      | `int64`                     |    no    |
|  2 | `symbol`      | `string`                    |    no    |
|  3 | `type`        | `Mt4OrderType`              |    no    |
|  4 | `volume`      | `double`                    |    no    |
|  5 | `open_price`  | `double`                    |    no    |
|  6 | `close_price` | `double`                    |    no    |
|  7 | `sl`          | `double`                    |    no    |
|  8 | `tp`          | `double`                    |    no    |
|  9 | `open_time`   | `google.protobuf.Timestamp` |    no    |
| 10 | `close_time`  | `google.protobuf.Timestamp` |    no    |
| 11 | `commission`  | `double`                    |    no    |
| 12 | `swap`        | `double`                    |    no    |
| 13 | `profit`      | `double`                    |    no    |
| 14 | `magic`       | `int32`                     |    no    |
| 15 | `comment`     | `string`                    |    no    |

### OrdersHistoryResponse

|  # | Field    | Type              | Repeated |
| -: | -------- | ----------------- | :------: |
|  1 | `orders` | `Mt4OrderHistory` |    yes   |

### Mt4OpenedOrdersResponse

|  # | Field    | Type             | Repeated |
| -: | -------- | ---------------- | :------: |
|  1 | `orders` | `Mt4OpenedOrder` |    yes   |

### OrderSendRequest

|  # | Field        | Type                          | Repeated |
| -: | ------------ | ----------------------------- | :------: |
|  1 | `symbol`     | `string`                      |    no    |
|  2 | `operation`  | `Mt4OrderSendOperationType`   |    no    |
|  3 | `volume`     | `double`                      |    no    |
|  4 | `price`      | `google.protobuf.DoubleValue` |    no    |
|  5 | `slippage`   | `int32`                       |    no    |
|  6 | `sl`         | `google.protobuf.DoubleValue` |    no    |
|  7 | `tp`         | `google.protobuf.DoubleValue` |    no    |
|  8 | `comment`    | `google.protobuf.StringValue` |    no    |
|  9 | `magic`      | `google.protobuf.Int32Value`  |    no    |
| 10 | `expiration` | `google.protobuf.Timestamp`   |    no    |

### OrderSendResult

|  # | Field       | Type                        | Repeated |
| -: | ----------- | --------------------------- | :------: |
|  1 | `ticket`    | `int64`                     |    no    |
|  2 | `price`     | `double`                    |    no    |
|  3 | `open_time` | `google.protobuf.Timestamp` |    no    |
|  4 | `error`     | `Mt4OrderError`             |    no    |

### OrderModifyRequest

|  # | Field        | Type                          | Repeated |
| -: | ------------ | ----------------------------- | :------: |
|  1 | `ticket`     | `int64`                       |    no    |
|  2 | `price`      | `google.protobuf.DoubleValue` |    no    |
|  3 | `sl`         | `google.protobuf.DoubleValue` |    no    |
|  4 | `tp`         | `google.protobuf.DoubleValue` |    no    |
|  5 | `expiration` | `google.protobuf.Timestamp`   |    no    |

### OrderCloseRequest

|  # | Field      | Type                          | Repeated |
| -: | ---------- | ----------------------------- | :------: |
|  1 | `ticket`   | `int64`                       |    no    |
|  2 | `volume`   | `double`                      |    no    |
|  3 | `price`    | `google.protobuf.DoubleValue` |    no    |
|  4 | `slippage` | `int32`                       |    no    |

### OrderDeleteRequest

|  # | Field      | Type    | Repeated |
| -: | ---------- | ------- | :------: |
|  1 | `ticket`   | `int64` |    no    |
|  2 | `slippage` | `int32` |    no    |

### OrderCloseByRequest

|  # | Field        | Type    | Repeated |
| -: | ------------ | ------- | :------: |
|  1 | `ticket_src` | `int64` |    no    |
|  2 | `ticket_dst` | `int64` |    no    |
|  3 | `slippage`   | `int32` |    no    |

### OrderActionResult

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `ticket` | `int64`                     |    no    |
|  2 | `price`  | `double`                    |    no    |
|  3 | `time`   | `google.protobuf.Timestamp` |    no    |
|  4 | `error`  | `Mt4OrderError`             |    no    |

### GetHistoryRequest

|  # | Field    | Type                          | Repeated |
| -: | -------- | ----------------------------- | :------: |
|  1 | `from`   | `google.protobuf.Timestamp`   |    no    |
|  2 | `to`     | `google.protobuf.Timestamp`   |    no    |
|  3 | `symbol` | `google.protobuf.StringValue` |    no    |

### GetOpenedOrdersRequest

|  # | Field    | Type                          | Repeated |
| -: | -------- | ----------------------------- | :------: |
|  1 | `symbol` | `google.protobuf.StringValue` |    no    |

### OrdersHistoryStreamRequest

|  # | Field    | Type                          | Repeated |
| -: | -------- | ----------------------------- | :------: |
|  1 | `from`   | `google.protobuf.Timestamp`   |    no    |
|  2 | `to`     | `google.protobuf.Timestamp`   |    no    |
|  3 | `symbol` | `google.protobuf.StringValue` |    no    |

### OrdersHistoryStreamChunk

|  # | Field     | Type              | Repeated |
| -: | --------- | ----------------- | :------: |
|  1 | `orders`  | `Mt4OrderHistory` |    yes   |
|  2 | `is_last` | `bool`            |    no    |

### OpenedOrdersTicket

|  # | Field    | Type    | Repeated |
| -: | -------- | ------- | :------: |
|  1 | `ticket` | `int64` |    no    |

### OpenedOrderTicketsResponse

|  # | Field     | Type                    | Repeated |
| -: | --------- | ----------------------- | :------: |
|  1 | `tickets` | `Mt4OpenedOrdersTicket` |    yes   |

### OpenedOrderProfitOrderInfo

|  # | Field          | Type     | Repeated |
| -: | -------------- | -------- | :------: |
|  1 | `ticket`       | `int64`  |    no    |
|  2 | `symbol`       | `string` |    no    |
|  3 | `order_profit` | `double` |    no    |
|  4 | `swap`         | `double` |    no    |
|  5 | `commission`   | `double` |    no    |

### OpenedOrdersProfitSnapshot

|  # | Field                               | Type                            | Repeated |
| -: | ----------------------------------- | ------------------------------- | :------: |
|  1 | `opened_orders_with_profit_updated` | `Mt4OpenedOrderProfitOrderInfo` |    yes   |
|  2 | `time`                              | `google.protobuf.Timestamp`     |    no    |

### OpenedOrdersProfitStreamRequest

|  # | Field         | Type    | Repeated |
| -: | ------------- | ------- | :------: |
|  1 | `buffer_size` | `int32` |    no    |

**Enums (7):**

### OrderType

| Name           | Value |
| -------------- | ----: |
| `OP_BUY`       |     0 |
| `OP_SELL`      |     1 |
| `OP_BUYLIMIT`  |     2 |
| `OP_SELLLIMIT` |     3 |
| `OP_BUYSTOP`   |     4 |
| `OP_SELLSTOP`  |     5 |

### OrderSendOperationType

| Name              | Value |
| ----------------- | ----: |
| `OC_OP_BUY`       |     0 |
| `OC_OP_SELL`      |     1 |
| `OC_OP_BUYLIMIT`  |     2 |
| `OC_OP_SELLLIMIT` |     3 |
| `OC_OP_BUYSTOP`   |     4 |
| `OC_OP_SELLSTOP`  |     5 |

### OrderResultCode

| Name                 | Value |
| -------------------- | ----: |
| `MT4_ORDER_OK`       |     0 |
| `MT4_ORDER_REJECTED` |     1 |
| `MT4_ORDER_PARTIAL`  |     2 |

### HistorySort

| Name                               | Value |
| ---------------------------------- | ----: |
| `MT4_HISTORY_SORT_OPEN_TIME_ASC`   |     0 |
| `MT4_HISTORY_SORT_OPEN_TIME_DESC`  |     1 |
| `MT4_HISTORY_SORT_CLOSE_TIME_ASC`  |     2 |
| `MT4_HISTORY_SORT_CLOSE_TIME_DESC` |     3 |

### OrdersFilter

| Name                             | Value |
| -------------------------------- | ----: |
| `MT4_ORDERS_FILTER_ALL`          |     0 |
| `MT4_ORDERS_FILTER_ONLY_MARKET`  |     1 |
| `MT4_ORDERS_FILTER_ONLY_PENDING` |     2 |

### ProfitStreamMode

| Name                               | Value |
| ---------------------------------- | ----: |
| `MT4_PROFIT_STREAM_MODE_UPDATES`   |     0 |
| `MT4_PROFIT_STREAM_MODE_SNAPSHOTS` |     1 |

### OrderAction

| Name                      | Value |
| ------------------------- | ----: |
| `MT4_ORDER_ACTION_CLOSE`  |     0 |
| `MT4_ORDER_ACTION_DELETE` |     1 |
| `MT4_ORDER_ACTION_MODIFY` |     2 |

---

## mt4-term-api-charts.proto

**Messages (8):**

### ChartBar

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `time`   | `google.protobuf.Timestamp` |    no    |
|  2 | `open`   | `double`                    |    no    |
|  3 | `high`   | `double`                    |    no    |
|  4 | `low`    | `double`                    |    no    |
|  5 | `close`  | `double`                    |    no    |
|  6 | `volume` | `int64`                     |    no    |

### ChartHistoryRequest

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `symbol` | `string`                    |    no    |
|  2 | `period` | `Mt4ChartPeriod`            |    no    |
|  3 | `from`   | `google.protobuf.Timestamp` |    no    |
|  4 | `to`     | `google.protobuf.Timestamp` |    no    |

### ChartHistoryResponse

|  # | Field  | Type          | Repeated |
| -: | ------ | ------------- | :------: |
|  1 | `bars` | `Mt4ChartBar` |    yes   |

### ChartStreamRequest

|  # | Field    | Type             | Repeated |
| -: | -------- | ---------------- | :------: |
|  1 | `symbol` | `string`         |    no    |
|  2 | `period` | `Mt4ChartPeriod` |    no    |

### ChartStreamChunk

|  # | Field     | Type          | Repeated |
| -: | --------- | ------------- | :------: |
|  1 | `bar`     | `Mt4ChartBar` |    no    |
|  2 | `is_last` | `bool`        |    no    |

### ChartTimeChunk

|  # | Field  | Type                        | Repeated |
| -: | ------ | --------------------------- | :------: |
|  1 | `from` | `google.protobuf.Timestamp` |    no    |
|  2 | `to`   | `google.protobuf.Timestamp` |    no    |

### ChartTimeChunks

|  # | Field    | Type                | Repeated |
| -: | -------- | ------------------- | :------: |
|  1 | `chunks` | `Mt4ChartTimeChunk` |    yes   |

### ChartHistoryStreamRequest

|  # | Field    | Type                 | Repeated |
| -: | -------- | -------------------- | :------: |
|  1 | `symbol` | `string`             |    no    |
|  2 | `period` | `Mt4ChartPeriod`     |    no    |
|  3 | `chunks` | `Mt4ChartTimeChunks` |    no    |

**Enums (3):**

### ChartPeriod

| Name         | Value |
| ------------ | ----: |
| `PERIOD_M1`  |     0 |
| `PERIOD_M5`  |     1 |
| `PERIOD_M15` |     2 |
| `PERIOD_M30` |     3 |
| `PERIOD_H1`  |     4 |
| `PERIOD_H4`  |     5 |
| `PERIOD_D1`  |     6 |
| `PERIOD_W1`  |     7 |
| `PERIOD_MN1` |     8 |

### ChartStreamMode

| Name                          | Value |
| ----------------------------- | ----: |
| `MT4_CHART_STREAM_MODE_BARS`  |     0 |
| `MT4_CHART_STREAM_MODE_TICKS` |     1 |

### ChartAggregation

| Name                 | Value |
| -------------------- | ----: |
| `MT4_CHART_AGG_NONE` |     0 |
| `MT4_CHART_AGG_OHLC` |     1 |

---

## mt4-term-api-connection.proto

**Messages (19):**

### ConnectRequest

|  # | Field      | Type     | Repeated |
| -: | ---------- | -------- | :------: |
|  1 | `login`    | `int64`  |    no    |
|  2 | `password` | `string` |    no    |
|  3 | `server`   | `string` |    no    |

### ConnectResponse

|  # | Field   | Type        | Repeated |
| -: | ------- | ----------- | :------: |
|  1 | `ok`    | `bool`      |    no    |
|  2 | `error` | `MrpcError` |    no    |

### DisconnectRequest

|  # | Field    | Type     | Repeated |
| -: | -------- | -------- | :------: |
|  1 | `reason` | `string` |    no    |

### DisconnectResponse

|  # | Field   | Type        | Repeated |
| -: | ------- | ----------- | :------: |
|  1 | `ok`    | `bool`      |    no    |
|  2 | `error` | `MrpcError` |    no    |

### ConnectionState

|  # | Field       | Type     | Repeated |
| -: | ----------- | -------- | :------: |
|  1 | `connected` | `bool`   |    no    |
|  2 | `login`     | `int64`  |    no    |
|  3 | `server`    | `string` |    no    |

### Heartbeat

|  # | Field  | Type                        | Repeated |
| -: | ------ | --------------------------- | :------: |
|  1 | `time` | `google.protobuf.Timestamp` |    no    |

### PingRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `payload` | `string` |    no    |

### PingResponse

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `payload` | `string` |    no    |

### ServerInfo

|  # | Field      | Type     | Repeated |
| -: | ---------- | -------- | :------: |
|  1 | `name`     | `string` |    no    |
|  2 | `address`  | `string` |    no    |
|  3 | `timezone` | `string` |    no    |

### GetServerInfoRequest

|  # | Field   | Type   | Repeated |
| -: | ------- | ------ | :------: |
|  1 | `dummy` | `bool` |    no    |

### Quote

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `symbol` | `string`                    |    no    |
|  2 | `bid`    | `double`                    |    no    |
|  3 | `ask`    | `double`                    |    no    |
|  4 | `point`  | `double`                    |    no    |
|  5 | `digits` | `int32`                     |    no    |
|  6 | `time`   | `google.protobuf.Timestamp` |    no    |

### QuoteRequest

|  # | Field    | Type     | Repeated |
| -: | -------- | -------- | :------: |
|  1 | `symbol` | `string` |    no    |

### QuoteBatchRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### QuoteBatchResponse

|  # | Field    | Type       | Repeated |
| -: | -------- | ---------- | :------: |
|  1 | `quotes` | `Mt4Quote` |    yes   |

### Time

|  # | Field         | Type                        | Repeated |
| -: | ------------- | --------------------------- | :------: |
|  1 | `server_time` | `google.protobuf.Timestamp` |    no    |

### ServerTimeRequest

|  # | Field   | Type   | Repeated |
| -: | ------- | ------ | :------: |
|  1 | `dummy` | `bool` |    no    |

### Version

|  # | Field            | Type     | Repeated |
| -: | ---------------- | -------- | :------: |
|  1 | `plugin_version` | `string` |    no    |
|  2 | `api_version`    | `string` |    no    |

### GetVersionRequest

|  # | Field   | Type   | Repeated |
| -: | ------- | ------ | :------: |
|  1 | `dummy` | `bool` |    no    |

**Enums (2):**

### ConnectionStateReason

| Name                                       | Value |
| ------------------------------------------ | ----: |
| `MT4_CONNECTION_STATE_REASON_UNKNOWN`      |     0 |
| `MT4_CONNECTION_STATE_REASON_DISCONNECTED` |     1 |
| `MT4_CONNECTION_STATE_REASON_CONNECTED`    |     2 |

### DisconnectReason

| Name                            | Value |
| ------------------------------- | ----: |
| `MT4_DISCONNECT_REASON_UNKNOWN` |     0 |
| `MT4_DISCONNECT_REASON_USER`    |     1 |
| `MT4_DISCONNECT_REASON_TIMEOUT` |     2 |

---

## mt4-term-api-health-check.proto

**Messages (4):**

### HealthSummary

|  # | Field   | Type        | Repeated |
| -: | ------- | ----------- | :------: |
|  1 | `ok`    | `bool`      |    no    |
|  2 | `error` | `MrpcError` |    no    |

### HealthCheckRequest

|  # | Field   | Type   | Repeated |
| -: | ------- | ------ | :------: |
|  1 | `dummy` | `bool` |    no    |

### HealthCheckResponse

|  # | Field     | Type               | Repeated |
| -: | --------- | ------------------ | :------: |
|  1 | `summary` | `Mt4HealthSummary` |    no    |

### HealthHeartbeat

|  # | Field  | Type                        | Repeated |
| -: | ------ | --------------------------- | :------: |
|  1 | `time` | `google.protobuf.Timestamp` |    no    |

---

## mt4-term-api-internal-charts.proto

**Messages (10):**

### InternalChartPoint

|  # | Field   | Type                        | Repeated |
| -: | ------- | --------------------------- | :------: |
|  1 | `time`  | `google.protobuf.Timestamp` |    no    |
|  2 | `value` | `double`                    |    no    |

### InternalChartSeries

|  # | Field    | Type                    | Repeated |
| -: | -------- | ----------------------- | :------: |
|  1 | `name`   | `string`                |    no    |
|  2 | `points` | `Mt4InternalChartPoint` |    yes   |

### InternalChartRequest

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `symbol` | `string`                    |    no    |
|  2 | `from`   | `google.protobuf.Timestamp` |    no    |
|  3 | `to`     | `google.protobuf.Timestamp` |    no    |

### InternalChartResponse

|  # | Field    | Type                     | Repeated |
| -: | -------- | ------------------------ | :------: |
|  1 | `series` | `Mt4InternalChartSeries` |    yes   |

### InternalChartStreamRequest

|  # | Field    | Type     | Repeated |
| -: | -------- | -------- | :------: |
|  1 | `symbol` | `string` |    no    |

### InternalChartStreamChunk

|  # | Field     | Type                    | Repeated |
| -: | --------- | ----------------------- | :------: |
|  1 | `point`   | `Mt4InternalChartPoint` |    no    |
|  2 | `is_last` | `bool`                  |    no    |

### InternalChartTimeChunk

|  # | Field  | Type                        | Repeated |
| -: | ------ | --------------------------- | :------: |
|  1 | `from` | `google.protobuf.Timestamp` |    no    |
|  2 | `to`   | `google.protobuf.Timestamp` |    no    |

### InternalChartTimeChunks

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `chunks` | `Mt4InternalChartTimeChunk` |    yes   |

### InternalChartHistoryStreamRequest

|  # | Field    | Type                         | Repeated |
| -: | -------- | ---------------------------- | :------: |
|  1 | `symbol` | `string`                     |    no    |
|  2 | `chunks` | `Mt4InternalChartTimeChunks` |    no    |

### InternalChartHistoryResponse

|  # | Field    | Type                     | Repeated |
| -: | -------- | ------------------------ | :------: |
|  1 | `series` | `Mt4InternalChartSeries` |    yes   |

**Enums (1):**

### InternalChartMode

| Name                             | Value |
| -------------------------------- | ----: |
| `MT4_INTERNAL_CHART_MODE_LINEAR` |     0 |
| `MT4_INTERNAL_CHART_MODE_STEP`   |     1 |

---

## mt4-term-api-market-info.proto

**Messages (14):**

### SymbolParams

|  # | Field           | Type     | Repeated |
| -: | --------------- | -------- | :------: |
|  1 | `symbol`        | `string` |    no    |
|  2 | `digits`        | `int32`  |    no    |
|  3 | `point`         | `double` |    no    |
|  4 | `lot_step`      | `double` |    no    |
|  5 | `lot_min`       | `double` |    no    |
|  6 | `lot_max`       | `double` |    no    |
|  7 | `stops_level`   | `int32`  |    no    |
|  8 | `freeze_level`  | `int32`  |    no    |
|  9 | `contract_size` | `double` |    no    |

### TickValue

|  # | Field        | Type     | Repeated |
| -: | ------------ | -------- | :------: |
|  1 | `symbol`     | `string` |    no    |
|  2 | `tick_value` | `double` |    no    |
|  3 | `tick_size`  | `double` |    no    |

### AllSymbolsResponse

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### SymbolsRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `pattern` | `string` |    no    |

### SymbolsResponse

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### TickValuesRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### TickValuesResponse

|  # | Field   | Type           | Repeated |
| -: | ------- | -------------- | :------: |
|  1 | `items` | `Mt4TickValue` |    yes   |

### SymbolParamsRequest

|  # | Field    | Type     | Repeated |
| -: | -------- | -------- | :------: |
|  1 | `symbol` | `string` |    no    |

### SymbolParamsResponse

|  # | Field    | Type              | Repeated |
| -: | -------- | ----------------- | :------: |
|  1 | `params` | `Mt4SymbolParams` |    no    |

### QuoteHistoryPoint

|  # | Field  | Type                        | Repeated |
| -: | ------ | --------------------------- | :------: |
|  1 | `time` | `google.protobuf.Timestamp` |    no    |
|  2 | `bid`  | `double`                    |    no    |
|  3 | `ask`  | `double`                    |    no    |

### QuoteHistoryRequest

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `symbol` | `string`                    |    no    |
|  2 | `from`   | `google.protobuf.Timestamp` |    no    |
|  3 | `to`     | `google.protobuf.Timestamp` |    no    |

### QuoteHistoryResponse

|  # | Field    | Type                   | Repeated |
| -: | -------- | ---------------------- | :------: |
|  1 | `quotes` | `Mt4QuoteHistoryPoint` |    yes   |

### SymbolsPagedRequest

|  # | Field       | Type    | Repeated |
| -: | ----------- | ------- | :------: |
|  1 | `page`      | `int32` |    no    |
|  2 | `page_size` | `int32` |    no    |

### SymbolsPagedResponse

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |
|  2 | `total`   | `int32`  |    no    |

**Enums (1):**

### MarketInfoError

| Name                         | Value |
| ---------------------------- | ----: |
| `MT4_MARKET_INFO_ERROR_NONE` |     0 |

---

## mt4-term-api-subscriptions.proto

**Messages (18):**

### SymbolsSubscribeRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### SymbolsUnsubscribeRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### SymbolsSubscribed

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### QuoteUpdate

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `symbol` | `string`                    |    no    |
|  2 | `bid`    | `double`                    |    no    |
|  3 | `ask`    | `double`                    |    no    |
|  4 | `point`  | `double`                    |    no    |
|  5 | `digits` | `int32`                     |    no    |
|  6 | `time`   | `google.protobuf.Timestamp` |    no    |

### QuoteStreamRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### QuoteStreamChunk

|  # | Field     | Type             | Repeated |
| -: | --------- | ---------------- | :------: |
|  1 | `update`  | `Mt4QuoteUpdate` |    no    |
|  2 | `is_last` | `bool`           |    no    |

### OrdersHistoryPagedRequest

|  # | Field       | Type                          | Repeated |
| -: | ----------- | ----------------------------- | :------: |
|  1 | `from`      | `google.protobuf.Timestamp`   |    no    |
|  2 | `to`        | `google.protobuf.Timestamp`   |    no    |
|  3 | `symbol`    | `google.protobuf.StringValue` |    no    |
|  4 | `page_size` | `int32`                       |    no    |

### OrdersHistoryPagedChunk

|  # | Field             | Type              | Repeated |
| -: | ----------------- | ----------------- | :------: |
|  1 | `orders`          | `Mt4OrderHistory` |    yes   |
|  2 | `next_page_token` | `string`          |    no    |
|  3 | `is_last`         | `bool`            |    no    |

### TradeUpdate

|  # | Field    | Type                        | Repeated |
| -: | -------- | --------------------------- | :------: |
|  1 | `ticket` | `int64`                     |    no    |
|  2 | `symbol` | `string`                    |    no    |
|  3 | `type`   | `Mt4OrderType`              |    no    |
|  4 | `volume` | `double`                    |    no    |
|  5 | `price`  | `double`                    |    no    |
|  6 | `sl`     | `double`                    |    no    |
|  7 | `tp`     | `double`                    |    no    |
|  8 | `time`   | `google.protobuf.Timestamp` |    no    |
|  9 | `state`  | `Mt4TradeUpdateState`       |    no    |

### TradeUpdateStreamRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### TradeUpdateStreamChunk

|  # | Field     | Type             | Repeated |
| -: | --------- | ---------------- | :------: |
|  1 | `update`  | `Mt4TradeUpdate` |    no    |
|  2 | `is_last` | `bool`           |    no    |

### OpenedOrdersTicketStreamRequest

|  # | Field     | Type     | Repeated |
| -: | --------- | -------- | :------: |
|  1 | `symbols` | `string` |    yes   |

### OpenedOrdersTicketStreamChunk

|  # | Field     | Type                    | Repeated |
| -: | --------- | ----------------------- | :------: |
|  1 | `tickets` | `Mt4OpenedOrdersTicket` |    yes   |
|  2 | `is_last` | `bool`                  |    no    |

### OpenedOrdersProfitStreamChunk

|  # | Field      | Type                            | Repeated |
| -: | ---------- | ------------------------------- | :------: |
|  1 | `snapshot` | `Mt4OpenedOrdersProfitSnapshot` |    no    |
|  2 | `is_last`  | `bool`                          |    no    |

### QuoteHistoryChunk

|  # | Field     | Type                   | Repeated |
| -: | --------- | ---------------------- | :------: |
|  1 | `point`   | `Mt4QuoteHistoryPoint` |    no    |
|  2 | `is_last` | `bool`                 |    no    |

**Enums (2):**

### TradeUpdateState

| Name                              | Value |
| --------------------------------- | ----: |
| `MT4_TRADE_UPDATE_STATE_OPENED`   |     0 |
| `MT4_TRADE_UPDATE_STATE_MODIFIED` |     1 |
| `MT4_TRADE_UPDATE_STATE_CLOSED`   |     2 |

### SubscriptionError

| Name                                        | Value |
| ------------------------------------------- | ----: |
| `MT4_SUBSCRIPTION_ERROR_NONE`               |     0 |
| `MT4_SUBSCRIPTION_ERROR_ALREADY_SUBSCRIBED` |     1 |

---

## mt4-term-api-trading-helper.proto

**Messages (12):**

### TradeRequest

|  # | Field        | Type                          | Repeated |
| -: | ------------ | ----------------------------- | :------: |
|  1 | `symbol`     | `string`                      |    no    |
|  2 | `type`       | `Mt4OrderType`                |    no    |
|  3 | `volume`     | `double`                      |    no    |
|  4 | `price`      | `google.protobuf.DoubleValue` |    no    |
|  5 | `sl`         | `google.protobuf.DoubleValue` |    no    |
|  6 | `tp`         | `google.protobuf.DoubleValue` |    no    |
|  7 | `deviation`  | `int32`                       |    no    |
|  8 | `comment`    | `google.protobuf.StringValue` |    no    |
|  9 | `magic`      | `google.protobuf.Int32Value`  |    no    |
| 10 | `expiration` | `google.protobuf.Timestamp`   |    no    |

### TradeResult

|  # | Field     | Type                        | Repeated |
| -: | --------- | --------------------------- | :------: |
|  1 | `retcode` | `Mt4OrderResultCode`        |    no    |
|  2 | `ticket`  | `int64`                     |    no    |
|  3 | `price`   | `double`                    |    no    |
|  4 | `time`    | `google.protobuf.Timestamp` |    no    |
|  5 | `comment` | `string`                    |    no    |
|  6 | `error`   | `Mt4OrderError`             |    no    |

### TradeActionRequest

|  # | Field        | Type                          | Repeated |
| -: | ------------ | ----------------------------- | :------: |
|  1 | `action`     | `Mt4OrderAction`              |    no    |
|  2 | `ticket`     | `int64`                       |    no    |
|  3 | `symbol`     | `string`                      |    no    |
|  4 | `volume`     | `double`                      |    no    |
|  5 | `price`      | `google.protobuf.DoubleValue` |    no    |
|  6 | `sl`         | `google.protobuf.DoubleValue` |    no    |
|  7 | `tp`         | `google.protobuf.DoubleValue` |    no    |
|  8 | `deviation`  | `int32`                       |    no    |
|  9 | `expiration` | `google.protobuf.Timestamp`   |    no    |

### TradeActionResponse

|  # | Field    | Type             | Repeated |
| -: | -------- | ---------------- | :------: |
|  1 | `result` | `Mt4TradeResult` |    no    |

### TradeHistoryRequest

|  # | Field    | Type                          | Repeated |
| -: | -------- | ----------------------------- | :------: |
|  1 | `from`   | `google.protobuf.Timestamp`   |    no    |
|  2 | `to`     | `google.protobuf.Timestamp`   |    no    |
|  3 | `symbol` | `google.protobuf.StringValue` |    no    |

### TradeHistoryResponse

|  # | Field    | Type              | Repeated |
| -: | -------- | ----------------- | :------: |
|  1 | `orders` | `Mt4OrderHistory` |    yes   |

### TradeHistoryPagedRequest

|  # | Field       | Type                          | Repeated |
| -: | ----------- | ----------------------------- | :------: |
|  1 | `from`      | `google.protobuf.Timestamp`   |    no    |
|  2 | `to`        | `google.protobuf.Timestamp`   |    no    |
|  3 | `symbol`    | `google.protobuf.StringValue` |    no    |
|  4 | `page_size` | `int32`                       |    no    |

### TradeHistoryPagedChunk

|  # | Field             | Type              | Repeated |
| -: | ----------------- | ----------------- | :------: |
|  1 | `orders`          | `Mt4OrderHistory` |    yes   |
|  2 | `next_page_token` | `string`          |    no    |
|  3 | `is_last`         | `bool`            |    no    |

### OpenedOrdersRequest

|  # | Field    | Type                          | Repeated |
| -: | -------- | ----------------------------- | :------: |
|  1 | `symbol` | `google.protobuf.StringValue` |    no    |

### OpenedOrdersResponse

|  # | Field    | Type             | Repeated |
| -: | -------- | ---------------- | :------: |
|  1 | `orders` | `Mt4OpenedOrder` |    yes   |

### OpenedOrdersPagedRequest

|  # | Field       | Type                          | Repeated |
| -: | ----------- | ----------------------------- | :------: |
|  1 | `symbol`    | `google.protobuf.StringValue` |    no    |
|  2 | `page_size` | `int32`                       |    no    |

### OpenedOrdersPagedChunk

|  # | Field             | Type             | Repeated |
| -: | ----------------- | ---------------- | :------: |
|  1 | `orders`          | `Mt4OpenedOrder` |    yes   |
|  2 | `next_page_token` | `string`         |    no    |
|  3 | `is_last`         | `bool`           |    no    |

**Enums (2):**

### TradeError

| Name                             | Value |
| -------------------------------- | ----: |
| `MT4_TRADE_ERROR_NONE`           |     0 |
| `MT4_TRADE_ERROR_TRADE_DISABLED` |     1 |

### TradeActionError

| Name                                    | Value |
| --------------------------------------- | ----: |
| `MT4_TRADE_ACTION_ERROR_NONE`           |     0 |
| `MT4_TRADE_ACTION_ERROR_INVALID_ACTION` |     1 |

---
