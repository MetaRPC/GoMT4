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
   │ New Order │
   └─────┬─────┘
         │ (executed at market / placed pending)
         ▼
   ┌───────────┐
   │   Open    │
   └───┬───┬───┘
       │   │
       │   │ SL hit (loss)
       │   ▼
       │ ┌───────┐
       │ │Closed │
       │ └───────┘
       │
       │ TP hit (profit)
       │
       ▼
   ┌───────────┐
   │ Delete/   │
   │ Cancelled │ (for pending orders)
   └───────────┘
```

This diagram shows a typical order lifecycle: creation → open → closure (SL/TP or manual) or cancellation (if pending).

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
* **Lot Step** → Minimum increment allowed when specifying volume.
* **Contract Size** → Amount of base currency per lot.
* **Stops Level** → Minimum distance (in points) required for SL/TP from current price.

---

## 📦 Orders

* **Order** → Instruction to buy or sell a symbol.
* **Market Order** → Executed immediately at current market price.
* **Pending Order** → Placed to execute in the future at a specific price (Limit or Stop).
* **Ticket** → Unique ID of an order (int64).
* **SL (Stop Loss)** → Protective level to cap loss.
* **TP (Take Profit)** → Target level to close with profit.
* **Magic Number** → User-defined integer to tag EAs/orders.
* **Comment** → Free text attached to an order.

---

## 🔄 Order Types (MT4)

* `OP_BUY` → Buy at market.
* `OP_SELL` → Sell at market.
* `OP_BUYLIMIT` → Pending: buy if price drops to X.
* `OP_SELLLIMIT` → Pending: sell if price rises to X.
* `OP_BUYSTOP` → Pending: buy if price rises to X.
* `OP_SELLSTOP` → Pending: sell if price drops to X.

---

## 🔌 Connection & RPC

* **gRPC** → Protocol used by GoMT4 to expose MT4 functions.
* **Port 50051** → Default local address where server listens (`127.0.0.1:50051`).
* **Stream** → Long‑lived connection pushing updates (quotes, orders, history).
* **Unary RPC** → One‑shot request/response (e.g., GetQuote).

---

## 📊 History & Streaming

* **Quote History** → Past bid/ask points (OHLC bars or ticks).
* **Order History** → List of closed trades for a period.
* **StreamQuotes** → Continuous tick updates.
* **StreamTradeUpdates** → Real-time feed of order lifecycle events.

---

## 🛡️ Errors & Codes

* **MrpcError** → Generic RPC error (code + message).
* **OrderError** → Error tied to a specific order (invalid volume, price, etc.).
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
