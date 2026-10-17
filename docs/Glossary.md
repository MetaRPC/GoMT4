# 📖 Glossary (MT4 Terms)

A quick reference for common MT4/GoMT4 terms used throughout the docs and code.

---

## 📝 Quick Cheat Sheet

| Term   | Example          | Meaning                                     |
| ------ | ---------------- | ------------------------------------------- |
| Symbol | `EURUSD`         | Instrument identifier                       |
| Lot    | `1.0` → 100,000  | Standard trading volume                     |
| SL     | `1.09500`        | Stop Loss (protect from loss)               |
| TP     | `1.10500`        | Take Profit (close with gain)               |
| Ticket | `12345678`       | Unique order ID                             |
| Digits | `5`              | Quote precision (1.23456)                   |
| Margin | `100.00`         | Locked funds for position                   |
| Equity | `1000.00`        | Balance ± open positions PnL                |
| Stream | `StreamQuotes()` | Continuous updates (ticks, orders, profits) |

---

## 📊 Order Lifecycle Diagram

```text
   ┌───────────┐
   │ New Order │  OrderSend
   └─────┬─────┘
         │
   ┌─────┴──────────────────────┐
   │ market (BUY/SELL)          │ pending (LIMIT/STOP)
   ▼                            ▼
┌───────────┐  price reached ┌───────────┐
│   Open    │◄───────────────│  Pending  │
└─────┬─────┘                └─────┬─────┘
      │ SL/TP hit or               │ OrderDelete
      │ OrderClose / OrderCloseBy  │ (or expiration)
      ▼                            ▼
┌───────────┐                ┌───────────┐
│  Closed   │                │  Deleted  │
└───────────┘                └───────────┘
```

This diagram shows a typical order lifecycle: creation → open → closure (SL/TP or manual), or for pending orders: trigger → open, or deletion.

---

## 🧑‍💻 Account

* **Login** → Numeric ID of your trading account.
* **Password** → Investor or trader password. Grants access to terminal.
* **Server** → Broker server name (e.g., `RoboForex-Demo`).
* **Balance** → Money currently on the account.
* **Equity** → Balance + open positions PnL.
* **Margin** → Funds locked for open positions.
* **Free Margin** → Equity − Margin.
* **Leverage** → Ratio (e.g., 1:500) showing how much borrowed funds you can use.

---

## 📈 Market Info

* **Symbol** → Instrument identifier (e.g., `EURUSD`).
* **Quote** → Current bid/ask prices for a symbol.
* **Digits** → Decimal precision of the quote (e.g., 5 digits → 1.23456).
* **Point** → Smallest price step for the symbol (e.g., 0.00001 for EURUSD).
* **Lot** → Standard trade size (usually 100,000 units base currency).
* **Volume Step** (`VolumeStep`, "lot step") → Minimum increment allowed when specifying volume; `VolumeMin` / `VolumeMax` are the limits.
* **Contract Size** → Amount of base currency per lot.
* **Stops Level** (`TradeStopsLevel`) → Minimum distance (in points) required for SL/TP from current price.

---

## 📦 Orders

* **Order** → Instruction to buy or sell a symbol.
* **Market Order** → Executed immediately at current market price.
* **Pending Order** → Placed to execute in the future at a specific price (Limit or Stop).
* **Ticket** → Unique ID of an order (`int32` in this API).
* **SL (Stop Loss)** → Protective level to cap loss.
* **TP (Take Profit)** → Target level to close with profit.
* **Magic Number** → User-defined integer to tag EAs/orders.
* **Comment** → Free text attached to an order.

---

## 🔄 Order Types (MT4)

| MT4 | Sending (`pb.OrderSendOperationType`) | In opened/history orders (`pb.OpenedOrderType`) | Meaning |
| --- | --- | --- | --- |
| `OP_BUY` | `OrderSendOperationType_OC_OP_BUY` | `OpenedOrderType_OO_OP_BUY` | Buy at market |
| `OP_SELL` | `OrderSendOperationType_OC_OP_SELL` | `OpenedOrderType_OO_OP_SELL` | Sell at market |
| `OP_BUYLIMIT` | `OrderSendOperationType_OC_OP_BUYLIMIT` | `OpenedOrderType_OO_OP_BUYLIMIT` | Pending: buy if price drops to X |
| `OP_SELLLIMIT` | `OrderSendOperationType_OC_OP_SELLLIMIT` | `OpenedOrderType_OO_OP_SELLLIMIT` | Pending: sell if price rises to X |
| `OP_BUYSTOP` | `OrderSendOperationType_OC_OP_BUYSTOP` | `OpenedOrderType_OO_OP_BUYSTOP` | Pending: buy if price rises to X |
| `OP_SELLSTOP` | `OrderSendOperationType_OC_OP_SELLSTOP` | `OpenedOrderType_OO_OP_SELLSTOP` | Pending: sell if price drops to X |

Stream events (`OnTrade`, `OnOpenedOrdersProfit`) use `pb.SUB_ORDER_OPERATION_TYPE` (`SUB_ORDER_OPERATION_TYPE_SUB_OP_BUY`, …) with the same six values.

---

## 🔌 Connection & RPC

* **gRPC** → Protocol used between GoMT4 and the MetaRPC API.
* **`mt4.mrpc.pro:443`** → The MetaRPC cloud API (gRPC over TLS) the client connects to by default. MetaRPC runs the MT4 terminal for your account.
* **API key** → Your MetaRPC key, sent as `apikey` metadata on every call (`TRIAL` when none is given).
* **Terminal id** → GUID of the terminal that serves your account, sent as `id` metadata (`TerminalInstanceGuid` after connect).
* **Stream** → Long‑lived server‑streaming call pushing updates (ticks, trade events, profits, tickets).
* **Unary RPC** → One‑shot request/response (e.g., `Quote`).

---

## 📊 History & Streaming

* **Quote History** → Past OHLC bars for a timeframe (`QuoteHistory`).
* **Order History** → List of closed trades for a period.
* **StreamQuotes** → Continuous tick updates.
* **StreamTradeUpdates** → Real-time feed of order lifecycle events.

---

## 🛡️ Errors & Codes

* **`pb.Error`** → API error inside a reply: `ErrorCode` + `ErrorMessage`, plus `MqlErrorCode` / `MqlErrorIntCode` when MT4 itself refused (invalid volume, price, etc.). The client returns it as a Go error `API error (code=...): ...`.
* **gRPC status error** → Transport-level failure (e.g., `Unavailable`, `DeadlineExceeded`).
* **Slippage** → Max price deviation allowed when executing orders.
* **Requote** → Server rejects trade at requested price, offers new one.

---

## ✅ Cheat Sheet (summary)

* Account → who you are.
* Market Info → what you trade.
* Orders → how you trade.
* Connection → how GoMT4 talks to MT4.
* History/Streaming → how you monitor trades and quotes.
* Errors → what can go wrong.
