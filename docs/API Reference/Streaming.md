# 📡 API Reference — Streaming (GoMT4)

This page documents **streaming messages & chunks** used in GoMT4. Streaming APIs are long-lived gRPC calls that continuously send updates.

---

## 🔔 Quotes Streaming

### Mt4QuoteUpdate (QuoteUpdate)

| Field                                                     | Type      |
| --------------------------------------------------------- | --------- |
| `symbol`                                                  | string    |
| `bid`                                                     | double    |
| `ask`                                                     | double    |
| `point`                                                   | double    |
| `digits`                                                  | int32     |
| `time`                                                    | Timestamp |
| **Notes:** one tick. Compute spread: `(ask - bid)/point`. |           |

### Mt4QuoteStreamRequest (QuoteStreamRequest)

* Fields: `symbols[]` (list of symbols to subscribe).

### Mt4QuoteStreamChunk (QuoteStreamChunk)

* Fields: `update: QuoteUpdate`, `is_last: bool`.
* `is_last = true` marks stream termination.

---

## 📑 Orders History Streaming

### Mt4OrdersHistoryPagedRequest (OrdersHistoryPagedRequest)

* Fields: `from`, `to`, optional `symbol`, `page_size`.

### Mt4OrdersHistoryPagedChunk (OrdersHistoryPagedChunk)

| Field                                                             | Type         |
| ----------------------------------------------------------------- | ------------ |
| `orders[]`                                                        | OrderHistory |
| `next_page_token`                                                 | string       |
| `is_last`                                                         | bool         |
| **Notes:** chunked history pages. Use token to request next page. |              |

---

## 🔄 Trade Updates

### Mt4TradeUpdate (TradeUpdate)

| Field    | Type             |
| -------- | ---------------- |
| `ticket` | int64            |
| `symbol` | string           |
| `type`   | OrderType        |
| `volume` | double           |
| `price`  | double           |
| `sl`     | double           |
| `tp`     | double           |
| `time`   | Timestamp        |
| `state`  | TradeUpdateState |

### Mt4TradeUpdateStreamRequest (TradeUpdateStreamRequest)

* Fields: `symbols[]` (subscribe per symbol).

### Mt4TradeUpdateStreamChunk (TradeUpdateStreamChunk)

* Fields: `update: TradeUpdate`, `is_last: bool`.

---

## 🎟️ Opened Orders Tickets

### Mt4OpenedOrdersTicketStreamRequest (OpenedOrdersTicketStreamRequest)

* Fields: `symbols[]` (subscribe tickets by symbol).

### Mt4OpenedOrdersTicketStreamChunk (OpenedOrdersTicketStreamChunk)

| Field       | Type               |
| ----------- | ------------------ |
| `tickets[]` | OpenedOrdersTicket |
| `is_last`   | bool               |

---

## 💰 Opened Orders Profit Stream

### Mt4OpenedOrdersProfitOrderInfo (ProfitOrderInfo)

| Field          | Type   |
| -------------- | ------ |
| `ticket`       | int64  |
| `symbol`       | string |
| `order_profit` | double |
| `swap`         | double |
| `commission`   | double |

### Mt4OpenedOrdersProfitSnapshot (ProfitSnapshot)

* Fields: `opened_orders_with_profit_updated[]: ProfitOrderInfo`, `time: Timestamp`.

### Mt4OpenedOrdersProfitStreamRequest (ProfitStreamRequest)

* Field: `buffer_size: int32` (controls channel buffer).

### Mt4OpenedOrdersProfitStreamChunk (ProfitStreamChunk)

* Fields: `snapshot: ProfitSnapshot`, `is_last: bool`.

**Notes:** updates real-time PnL per order.

---

## 🕒 Quote History Streaming

### Mt4QuoteHistoryPoint (QuoteHistoryPoint)

* Fields: `time, bid, ask`.

### Mt4QuoteHistoryChunk (QuoteHistoryChunk)

* Fields: `point: QuoteHistoryPoint`, `is_last: bool`.

---

## 📊 Chart Streaming

### Mt4ChartStreamRequest (ChartStreamRequest)

* Fields: `symbol, period`.

### Mt4ChartStreamChunk (ChartStreamChunk)

* Fields: `bar: ChartBar`, `is_last: bool`.

### Mt4ChartHistoryStreamRequest (ChartHistoryStreamRequest)

* Fields: `symbol, period, chunks: ChartTimeChunks`.

### Mt4ChartTimeChunks (ChartTimeChunks)

* Array of `from/to` ranges.

---

## 📈 Internal Chart Streaming

### Mt4InternalChartStreamRequest (InternalChartStreamRequest)

* Field: `symbol`.

### Mt4InternalChartStreamChunk (InternalChartStreamChunk)

* Fields: `point: InternalChartPoint`, `is_last: bool`.

### Mt4InternalChartHistoryStreamRequest (InternalChartHistoryStreamRequest)

* Fields: `symbol`, `chunks: InternalChartTimeChunks`.

### Mt4InternalChartHistoryResponse (InternalChartHistoryResponse)

* Fields: `series[]: InternalChartSeries`.

---

## ⚠️ Stream Errors

### Mt4StreamErrorReason (StreamErrorReason)

| Name          | Value | Meaning            |
| ------------- | ----: | ------------------ |
| `NONE`        |     0 | No error           |
| `EOF`         |     1 | End of stream      |
| `UNAVAILABLE` |     2 | Stream unavailable |

---

📌 **Tip:**

* Use `context.WithTimeout` or `WithCancel` to control stream lifetime.
* Always handle `is_last = true` and error channels in client code.
* Streaming is ideal for real-time dashboards (quotes, PnL, trades).
