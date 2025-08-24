# 📘 API Reference — Messages (GoMT4)

This page documents **message types** used by GoMT4 (from your `.proto`).
Numbers, names and field order match the source. Enums are on a separate page.

---

## 🔖 Conventions

* **`google.protobuf.Timestamp`** → ⏰ UTC time; format with `time.RFC3339` for logs.
* **`google.protobuf.{Double,String,Int32}Value`** → 🎛 optional fields with presence (omit if not set).
* **Prices & volumes** → 💹 respect `SymbolParams` (`digits`, `point`, `lot_step`, `lot_min/max`).
* **Currency** → 💵 monetary fields are in **account currency** (see `AccountSummary.currency`).

---

## 🧾 Account & Orders

### 📊 Mt4AccountSummary

|                                                                      # | Field          | Type   |
| ---------------------------------------------------------------------: | -------------- | ------ |
|                                                                      1 | `login`        | int64  |
|                                                                      2 | `name`         | string |
|                                                                      3 | `server`       | string |
|                                                                      4 | `currency`     | string |
|                                                                      5 | `balance`      | double |
|                                                                      6 | `equity`       | double |
|                                                                      7 | `margin`       | double |
|                                                                      8 | `margin_free`  | double |
|                                                                      9 | `margin_level` | double |
|                                                                     10 | `leverage`     | int32  |
|                                                                     11 | `profit`       | double |
| **Notes:** snapshot for health/risk checks. `margin_level` in percent. |                |        |

### 📌 Mt4OpenedOrder

|                                                           # | Field        | Type           |
| ----------------------------------------------------------: | ------------ | -------------- |
|                                                           1 | `ticket`     | int64          |
|                                                           2 | `symbol`     | string         |
|                                                           3 | `type`       | `Mt4OrderType` |
|                                                           4 | `volume`     | double         |
|                                                           5 | `open_price` | double         |
|                                                           6 | `sl`         | double         |
|                                                           7 | `tp`         | double         |
|                                                           8 | `open_time`  | Timestamp      |
|                                                           9 | `commission` | double         |
|                                                          10 | `swap`       | double         |
|                                                          11 | `magic`      | int32          |
|                                                          12 | `comment`    | string         |
| **Notes:** live positions; use with `OpenedOrders`/streams. |              |                |

### 📜 Mt4OrderHistory

|                                                   # | Field         | Type           |
| --------------------------------------------------: | ------------- | -------------- |
|                                                   1 | `ticket`      | int64          |
|                                                   2 | `symbol`      | string         |
|                                                   3 | `type`        | `Mt4OrderType` |
|                                                   4 | `volume`      | double         |
|                                                   5 | `open_price`  | double         |
|                                                   6 | `close_price` | double         |
|                                                   7 | `sl`          | double         |
|                                                   8 | `tp`          | double         |
|                                                   9 | `open_time`   | Timestamp      |
|                                                  10 | `close_time`  | Timestamp      |
|                                                  11 | `commission`  | double         |
|                                                  12 | `swap`        | double         |
|                                                  13 | `profit`      | double         |
|                                                  14 | `magic`       | int32          |
|                                                  15 | `comment`     | string         |
| **Notes:** closed deals; combine with history APIs. |               |                |

### 📨 Mt4OrderSendRequest

|          # | Field        | Type                        | Optional |
| ---------: | ------------ | --------------------------- | :------: |
|          1 | `symbol`     | string                      |          |
|          2 | `operation`  | `Mt4OrderSendOperationType` |          |
|          3 | `volume`     | double                      |          |
|          4 | `price`      | DoubleValue                 |     ✓    |
|          5 | `slippage`   | int32                       |          |
|          6 | `sl`         | DoubleValue                 |     ✓    |
|          7 | `tp`         | DoubleValue                 |     ✓    |
|          8 | `comment`    | StringValue                 |     ✓    |
|          9 | `magic`      | Int32Value                  |     ✓    |
|         10 | `expiration` | Timestamp                   |     ✓    |
| **Notes:** |              |                             |          |

* Market orders: `price` usually omitted (server uses current Bid/Ask, with `slippage`).
* Pending orders: set `price` and (optionally) `expiration`.

### 📬 Mt4OrderSendResult

|                                                   # | Field       | Type            |
| --------------------------------------------------: | ----------- | --------------- |
|                                                   1 | `ticket`    | int64           |
|                                                   2 | `price`     | double          |
|                                                   3 | `open_time` | Timestamp       |
|                                                   4 | `error`     | `Mt4OrderError` |
| **Notes:** `error` filled on broker-side rejection. |             |                 |

### ✏️ Mt4OrderModifyRequest

|  # | Field        | Type        | Optional |
| -: | ------------ | ----------- | :------: |
|  1 | `ticket`     | int64       |          |
|  2 | `price`      | DoubleValue |     ✓    |
|  3 | `sl`         | DoubleValue |     ✓    |
|  4 | `tp`         | DoubleValue |     ✓    |
|  5 | `expiration` | Timestamp   |     ✓    |

### ❌ Mt4OrderCloseRequest

|  # | Field      | Type        | Optional |
| -: | ---------- | ----------- | :------: |
|  1 | `ticket`   | int64       |          |
|  2 | `volume`   | double      |          |
|  3 | `price`    | DoubleValue |     ✓    |
|  4 | `slippage` | int32       |          |

### 🗑️ Mt4OrderDeleteRequest

|  # | Field      | Type  |
| -: | ---------- | ----- |
|  1 | `ticket`   | int64 |
|  2 | `slippage` | int32 |

### 🔄 Mt4OrderCloseByRequest

|  # | Field        | Type  |
| -: | ------------ | ----- |
|  1 | `ticket_src` | int64 |
|  2 | `ticket_dst` | int64 |
|  3 | `slippage`   | int32 |

### 📑 Mt4OrderActionResult

|  # | Field    | Type            |
| -: | -------- | --------------- |
|  1 | `ticket` | int64           |
|  2 | `price`  | double          |
|  3 | `time`   | Timestamp       |
|  4 | `error`  | `Mt4OrderError` |

---

## 💹 Quotes & Market Info

### 💱 Mt4Quote

|                                                           # | Field    | Type      |
| ----------------------------------------------------------: | -------- | --------- |
|                                                           1 | `symbol` | string    |
|                                                           2 | `bid`    | double    |
|                                                           3 | `ask`    | double    |
|                                                           4 | `point`  | double    |
|                                                           5 | `digits` | int32     |
|                                                           6 | `time`   | Timestamp |
| **Notes:** compute spread in points: `(ask - bid) / point`. |          |           |

### 📊 Mt4SymbolParams

|                                                                      # | Field           | Type   |
| ---------------------------------------------------------------------: | --------------- | ------ |
|                                                                      1 | `symbol`        | string |
|                                                                      2 | `digits`        | int32  |
|                                                                      3 | `point`         | double |
|                                                                      4 | `lot_step`      | double |
|                                                                      5 | `lot_min`       | double |
|                                                                      6 | `lot_max`       | double |
|                                                                      7 | `stops_level`   | int32  |
|                                                                      8 | `freeze_level`  | int32  |
|                                                                      9 | `contract_size` | double |
| **Notes:** use for rounding, SL/TP distance checks and lot validation. |                 |        |

---

## 🔌 Connection & Health

### 🔐 Connect/Disconnect

* **Mt4ConnectRequest**: `login,password,server` → **Mt4ConnectResponse**: `ok,error`.
* **Mt4DisconnectRequest**: `reason` → **Mt4DisconnectResponse**: `ok,error`.

### 📡 State & Heartbeat

* **Mt4ConnectionState**: `connected,login,server`.
* **Mt4Heartbeat**: `time` (keep-alive).

### 🏓 Ping & Version

* **Mt4Ping{Request,Response}**: echo `payload`.
* **Mt4Version**: `plugin_version`, `api_version`. Request is dummy.

### 🖥️ Server info & time

* **Mt4ServerInfo**: `name,address,timezone`.
* **Mt4Time**: `server_time`. Request is dummy.

### ❤️ Health-check API

* **Mt4HealthSummary**: `ok`, `error`.
* \*\*Mt4HealthCheck{Request,Response}`**: wraps `summary\`.

---

## 📡 Streaming payload helpers (overview)

> Detailed stream methods & chunk types are documented on the **Streaming** page.

* **Quotes**: `Mt4QuoteUpdate`, `Mt4QuoteStream{Request,Chunk}`.
* **Orders history (paged/stream)**: `Mt4OrdersHistoryPaged{Request,Chunk}`.
* **Trade updates**: `Mt4TradeUpdate`, `Mt4TradeUpdateStream{Request,Chunk}`.
* **Opened order tickets**: `Mt4OpenedOrdersTicketStream{Request,Chunk}`.
* **Opened order profits**: `Mt4OpenedOrdersProfit{StreamRequest,Snapshot,StreamChunk}`.
* **Chart streams**: `Mt4Chart{StreamRequest,StreamChunk}`, `Mt4ChartHistoryStreamRequest`.
* **Internal charts**: `Mt4InternalChart{StreamRequest,StreamChunk,HistoryStreamRequest}`.

---

## ⚠️ Errors (messages only)

* **MrpcError**: `code:int32, message, details`.
* **Mt4OrderError**: `ticket:int64, reason:int32, text`.
* **Mt4BatchOrderError**: `errors[]` of `MrpcError`.

**See also:** the **Enums** page for error/enum value tables.
