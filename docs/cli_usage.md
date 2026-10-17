# 🧰 Using GoMT4 via CLI (No GUI)

This section demonstrates how to **use GoMT4 directly from the terminal**, without any graphical user interface (GUI). Ideal for developers, DevOps, and command-line enthusiasts who prefer full control.

---

## 🔧 Requirements

| Tool          | Purpose                                      |
| ------------- | -------------------------------------------- |
| Go 1.23+      | For building and running the project         |
| MT4 account   | Login, password and broker server name (MetaRPC runs the MT4 terminal for you in the cloud) |
| API key       | From [https://mrpc.pro/signup](https://mrpc.pro/signup); without one the calls run as a trial |
| `config.json` | Login credentials, default symbol and API key |
| Terminal      | All operations are executed via command-line |

---

## 📁 Project Structure

```bash
GoMT4/
├── examples/                      # Example usage and main executable code
│   ├── config/
│   │   ├── config.go              # Loads configuration from a JSON file
│   │   └── config.json            # Connection settings: login, password, server, etc.
│   │
│   ├── mt4/
│   │   ├── MT4Account.go          # Handles MT4 connection and account-level operations
│   │   └── MT4_service.go         # Provides service methods for orders, streaming, etc.
│   │
│   ├── go.mod                     # Go module definition for the examples package
│   ├── go.sum                     # Dependency checksums
│   └── main.go                    # Entry point demonstrating MT4Service usage
│
├── package/                       # Generated protobuf/gRPC Go code (module github.com/MetaRPC/GoMT4/package)
├── docs/                          # This documentation
├── go.work                        # Go workspace tying examples/ and package/ together
├── .gitignore                     # Files/folders to be excluded from Git
└── mkdocs.yml                     # Documentation config file (for mkdocs or similar tools)


```

---

## 🧩 Example `config.json`

```json
{
  "Login": 12345678,
  "Password": "<your-mt4-password>",
  "Server": "MetaQuotes-Demo",
  "DefaultSymbol": "EURUSD",
  "ApiKey": "<your-mrpc-api-key>"
}
```

With `"Login": 0` the example opens a new demo account on `Server` and uses it.

---

## 🚀 Running the App

```bash
cd GoMT4/examples
go run .
# or pass the API key on the command line
go run . --api-key=YOUR_KEY
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
| `ShowAllSymbols()`           | List all available trading instruments with their index |
| `ShowSymbolParams("EURUSD")` | Print key symbol parameters (digits, volume limits, trade mode, currencies) |
| `ShowSymbols()`              | List symbol names with their index     |
| `ShowTickValues([...])`      | Get tick/contract values for symbols   |

---

### 🔄 Streaming / Subscriptions

| Function                         | Description                                                              |
| -------------------------------- | ------------------------------------------------------------------------ |
| `StreamQuotes()`                 | Live ticks for EURUSD and GBPUSD (stops after 3 ticks or 5s without one) |
| `StreamOpenedOrderProfits()`     | Real-time profit updates per open order (every 1000 ms)                  |
| `StreamOpenedOrderTickets()`     | Monitor currently open order tickets (every 1000 ms)                     |
| `StreamTradeUpdates()`           | Trade events; prints the first new order of each event                   |
| `StreamOrdersHistoryExample()`   | Page-by-page stream of order history (demo; last 30 days, close-time DESC) |
| `StreamQuoteHistoryExample()`    | Time-chunked OHLC history stream (demo; H1, last 90 days, weekly chunks) |

> Example output:

```
[Tick] EURUSD | Bid: 1.09876 | Ask: 1.09889 | Time: 2025-07-29 18:00:01
```

---

## 🧑‍💻 How to Enable a Function

In `examples/main.go`, add (or comment out) the method calls you want:

```go
svc.ShowQuote(ctx, "EURUSD")
svc.StreamQuotes(ctx)
```

You can call multiple methods one after another — for example, open an order and immediately monitor it.

---

## 🧠 Tips

* Use `context.WithTimeout(...)` to limit long operations.
* End with `account.Disconnect(...)` (the example does `Disconnect(true)`, which also deletes the terminal instance on the MetaRPC side) so no terminal is left running for the account.
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
