# 🧰 Using GoMT4 via CLI (No GUI)

This section demonstrates how to **use GoMT4 directly from the terminal**, without any graphical user interface (GUI). Ideal for developers, DevOps, and command-line enthusiasts who prefer full control.

---

## 🔧 Requirements

| Tool          | Purpose                                      |
| ------------- | -------------------------------------------- |
| Go 1.20+      | For building and running the project         |
| MetaTrader 4  | Terminal with gRPC plugin enabled            |
| `config.json` | Login credentials and default symbol         |
| Terminal      | All operations are executed via command-line |

---

## 📁 Project Structure

```bash
GoMT4/
├── examples/
│   └── main.go         // Usage example MT4Account
├── internal/
│   ├── account/
│   │   └── MT4Account.go
│   └── service/
│       └── MT4_service.go
├── pb/
│   └── ...             // gRPC-generated files
├── config.json
├── go.mod
└── go.sum
 
```

---

## 🧩 Example `config.json`

```json
{
  "Login": 501401178,
  "Password": "v8gctta",
  "Server": "RoboForex-Demo",
  "DefaultSymbol": "EURUSD"
}
```

---

## 🚀 Running the App

```bash
go run main.go
```

> If all goes well, you’ll see:
> `✅ Connected to MT4 server`
> and other output depending on the enabled functions.

---

## 🧪 Available Functions

### 🧾 Account Information

| Function                 | Description                         |
| ------------------------ | ----------------------------------- |
| `ShowAccountSummary`     | Print balance, equity, and currency |
| `ShowOpenedOrders`       | List current open orders            |
| `ShowOrdersHistory`      | View closed trades from last 7 days |
| `ShowOpenedOrderTickets` | Print open order ticket numbers     |

---

### 📦 Order Operations

| Function                          | Description                       |
| --------------------------------- | --------------------------------- |
| `ShowOrderSendExample("EURUSD")`  | Submit a sample Buy order         |
| `ShowOrderModifyExample(ticket)`  | Modify SL/TP for a ticket         |
| `ShowOrderCloseExample(ticket)`   | Close an order by ticket          |
| `ShowOrderDeleteExample(ticket)`  | Delete a pending order            |
| `ShowOrderCloseByExample(t1, t2)` | Close one order with its opposite |

⚠️ Real order execution (even on demo) — use carefully.

---

### 📈 Market Info & Symbols

| Function                     | Description                            |
| ---------------------------- | -------------------------------------- |
| `ShowQuote("EURUSD")`        | Get live bid/ask quote                 |
| `ShowQuotesMany([...])`      | Get quotes for multiple symbols        |
| `ShowQuoteHistory("EURUSD")` | Get last 5 days of OHLC candles        |
| `ShowAllSymbols()`           | List all available trading instruments |
| `ShowSymbolParams("EURUSD")` | Get full symbol metadata               |
| `ShowTickValues([...])`      | Get tick/contract values for symbols   |

---

### 🔄 Streaming / Subscriptions

| Function                     | Description                             |
| ---------------------------- | --------------------------------------- |
| `StreamQuotes()`             | Subscribe to live tick updates          |
| `StreamOpenedOrderProfits()` | Real-time profit updates per open order |
| `StreamOpenedOrderTickets()` | Monitor currently open order tickets    |
| `StreamTradeUpdates()`       | Subscribe to all trading events         |

> Example output:

```
[Tick] EURUSD | Bid: 1.09876 | Ask: 1.09889 | Time: 2025-07-29 18:00:01
```

---

## 🧑‍💻 How to Enable a Function

In `main.go`, uncomment the desired method:

```go
svc.ShowQuote(ctx, "EURUSD")
svc.StreamQuotes(ctx)
```

You can call multiple methods one after another — for example, open an order and immediately monitor it.

---

## 🧠 Tips

* Use `context.WithTimeout(...)` to limit long operations.
* Stop the MT4 terminal gracefully to avoid lingering gRPC connections.
* Even on demo, actions like sending orders are real.

---

## 📎 Quick Example

```go
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, "EURUSD")
svc.ShowOrderSendExample(ctx, "EURUSD")
svc.ShowOpenedOrders(ctx)
svc.StreamQuotes(ctx)
```

---

This is your terminal-powered trading dashboard — minimal, fast, and fully controlled by code.
