# 🎛️ API Reference — Enums (GoMT4)

This page documents all **enumerations** used by GoMT4. Names and values match `.proto`. For readability we show the shorter alias in parentheses.

---

## 📊 Orders

### Mt4OrderType (OrderType)

| Name           | Value | Meaning            |
| -------------- | ----: | ------------------ |
| `OP_BUY`       |     0 | Market Buy         |
| `OP_SELL`      |     1 | Market Sell        |
| `OP_BUYLIMIT`  |     2 | Pending Buy Limit  |
| `OP_SELLLIMIT` |     3 | Pending Sell Limit |
| `OP_BUYSTOP`   |     4 | Pending Buy Stop   |
| `OP_SELLSTOP`  |     5 | Pending Sell Stop  |

### Mt4OrderSendOperationType (OrderSendOp)

Same values/meaning as `OrderType`, but used in `OrderSendRequest.operation`.

### Mt4OrderResultCode (OrderResultCode)

| Name                 | Value | Meaning               |
| -------------------- | ----: | --------------------- |
| `MT4_ORDER_OK`       |     0 | Order executed OK     |
| `MT4_ORDER_REJECTED` |     1 | Broker rejected order |
| `MT4_ORDER_PARTIAL`  |     2 | Partially filled      |

### Mt4OrderAction (OrderAction)

| Name                      | Value | Meaning               |
| ------------------------- | ----: | --------------------- |
| `MT4_ORDER_ACTION_CLOSE`  |     0 | Close market order    |
| `MT4_ORDER_ACTION_DELETE` |     1 | Delete pending order  |
| `MT4_ORDER_ACTION_MODIFY` |     2 | Modify SL/TP or price |

---

## 📈 History & Filters

### Mt4HistorySort (HistorySort)

| Name              | Value | Meaning                       |
| ----------------- | ----: | ----------------------------- |
| `OPEN_TIME_ASC`   |     0 | Sort by open time ascending   |
| `OPEN_TIME_DESC`  |     1 | Sort by open time descending  |
| `CLOSE_TIME_ASC`  |     2 | Sort by close time ascending  |
| `CLOSE_TIME_DESC` |     3 | Sort by close time descending |

### Mt4OrdersFilter (OrdersFilter)

| Name           | Value | Meaning             |
| -------------- | ----: | ------------------- |
| `ALL`          |     0 | All orders          |
| `ONLY_MARKET`  |     1 | Only market orders  |
| `ONLY_PENDING` |     2 | Only pending orders |

---

## 💹 Quotes & Charts

### Mt4ChartPeriod (ChartPeriod)

| Name         | Value | Interval   |
| ------------ | ----: | ---------- |
| `PERIOD_M1`  |     0 | 1 minute   |
| `PERIOD_M5`  |     1 | 5 minutes  |
| `PERIOD_M15` |     2 | 15 minutes |
| `PERIOD_M30` |     3 | 30 minutes |
| `PERIOD_H1`  |     4 | 1 hour     |
| `PERIOD_H4`  |     5 | 4 hours    |
| `PERIOD_D1`  |     6 | 1 day      |
| `PERIOD_W1`  |     7 | 1 week     |
| `PERIOD_MN1` |     8 | 1 month    |

### Mt4ChartStreamMode (ChartStreamMode)

| Name    | Value | Meaning          |
| ------- | ----: | ---------------- |
| `BARS`  |     0 | Stream OHLC bars |
| `TICKS` |     1 | Stream ticks     |

### Mt4ChartAggregation (ChartAgg)

| Name   | Value | Meaning           |
| ------ | ----: | ----------------- |
| `NONE` |     0 | Raw data          |
| `OHLC` |     1 | Aggregate to OHLC |

---

## 🔌 Connection

### Mt4ConnectionStateReason (ConnStateReason)

| Name           | Value | Meaning       |
| -------------- | ----: | ------------- |
| `UNKNOWN`      |     0 | Unknown state |
| `DISCONNECTED` |     1 | Disconnected  |
| `CONNECTED`    |     2 | Connected     |

### Mt4DisconnectReason (DisconnectReason)

| Name      | Value | Meaning        |
| --------- | ----: | -------------- |
| `UNKNOWN` |     0 | Unknown reason |
| `USER`    |     1 | User-requested |
| `TIMEOUT` |     2 | Timeout        |

---

## 🧮 Profits & Streams

### Mt4ProfitStreamMode (ProfitStreamMode)

| Name        | Value | Meaning                 |
| ----------- | ----: | ----------------------- |
| `UPDATES`   |     0 | Send only updates       |
| `SNAPSHOTS` |     1 | Send periodic snapshots |

### Mt4TradeUpdateState (TradeUpdateState)

| Name       | Value | Meaning        |
| ---------- | ----: | -------------- |
| `OPENED`   |     0 | Order opened   |
| `MODIFIED` |     1 | Order modified |
| `CLOSED`   |     2 | Order closed   |

---

## ⚠️ Errors

### MrpcErrorCode (ErrorCode)

| Name          | Value | Meaning             |
| ------------- | ----: | ------------------- |
| `NONE`        |     0 | No error            |
| `UNKNOWN`     |     1 | Unknown error       |
| `TIMEOUT`     |     2 | Timeout occurred    |
| `CONNECTION`  |     3 | Connection lost     |
| `UNAVAILABLE` |     4 | Service unavailable |

### Mt4OrderErrorReason (OrderErrorReason)

| Name             | Value | Meaning           |
| ---------------- | ----: | ----------------- |
| `NONE`           |     0 | No error          |
| `INVALID_VOLUME` |     1 | Wrong lot size    |
| `INVALID_PRICE`  |     2 | Wrong price       |
| `MARKET_CLOSED`  |     3 | Market closed     |
| `SERVER_BUSY`    |     4 | Trade server busy |

### Mt4StreamErrorReason (StreamErrorReason)

| Name          | Value | Meaning            |
| ------------- | ----: | ------------------ |
| `NONE`        |     0 | No error           |
| `EOF`         |     1 | End of stream      |
| `UNAVAILABLE` |     2 | Stream unavailable |

---

## 📊 Misc

### Mt4MarketInfoError (MarketInfoError)

| Name   | Value | Meaning  |
| ------ | ----: | -------- |
| `NONE` |     0 | No error |

### Mt4TradeError (TradeError)

| Name             | Value | Meaning          |
| ---------------- | ----: | ---------------- |
| `NONE`           |     0 | No error         |
| `TRADE_DISABLED` |     1 | Trading disabled |

### Mt4TradeActionError (TradeActionError)

| Name             | Value | Meaning            |
| ---------------- | ----: | ------------------ |
| `NONE`           |     0 | No error           |
| `INVALID_ACTION` |     1 | Action not allowed |

### Mt4InternalChartMode (InternalChartMode)

| Name     | Value | Meaning     |
| -------- | ----: | ----------- |
| `LINEAR` |     0 | Linear mode |
| `STEP`   |     1 | Step mode   |

### Mt4SubscriptionError (SubscriptionError)

| Name                 | Value | Meaning            |
| -------------------- | ----: | ------------------ |
| `NONE`               |     0 | No error           |
| `ALREADY_SUBSCRIBED` |     1 | Already subscribed |

---

📌 **Tip:** Enums are used in requests and responses. Always check which enum a field expects (`OrderType`, `ChartPeriod`, etc.). Wrong values will cause broker errors or ignored requests.
