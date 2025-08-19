# Getting Started with MetaTrader 4 in Go

Welcome to the **MetaRPC MT4 Go Documentation** — your guide to integrating with **MetaTrader 4** using **Go** and **gRPC**.

This documentation will help you:

* 📘 Explore all available **account, market, and trading methods**
* 💡 Learn from **Go usage examples** with context and timeout handling
* 🔁 Work with **real-time streaming** for quotes, orders, and trades
* ⚙️ Understand all **input/output types** such as `OrderInfo`, `QuoteData`, and enums like `ENUM_ORDER_TYPE_TF`

---

## 📚 Main Sections

### Account

* [Show Account Summary](Account/ShowAccountSummary.md)

---

### Market Info

* **Overview:** \[Market Info — Overview]\(Market\_Info/Market\_Info — Overview\.md)
* [Show Quote](Market_Info/ShowQuote.md)
* [Show Quotes Many](Market_Info/ShowQuotesMany.md)
* [Show Quote History](Market_Info/ShowQuoteHistory.md)
* [Show Symbol Params](Market_Info/ShowSymbolParams.md)
* [Show Symbols](Market_Info/ShowSymbols.md)
* [Show All Symbols](Market_Info/ShowAllSymbols.md)
* [Show Tick Values](Market_Info/ShowTickValues.md)

---

### Order Operations ⚠️

* **Overview:** \[Order Operations — Overview]\(Order\_Operations/Order\_Operations — Overview\.md)
* [Show Opened Orders](Order_Operations/ShowOpenedOrders.md)
* [Show Opened Order Tickets](Order_Operations/ShowOpenedOrderTickets.md)
* [Show Orders History](Order_Operations/ShowOrdersHistory.md)
* [Order Close Example](Order_Operations/OrderCloseExample.md)
* [Order Close By Example](Order_Operations/ShowOrderCloseByExample.md)
* [Order Delete Example](Order_Operations/ShowOrderDeleteExample.md)
* [Order Modify Example](Order_Operations/ShowOrderModifyExample.md)

---

### Streaming

* **Overview:** \[Streaming — Overview]\(Streaming/Streaming — Overview\.md)
* [Stream Opened Order Profits](Streaming/StreamOpenedOrderProfits.md)
* [Stream Opened Order Tickets](Streaming/StreamOpenedOrderTickets.md)
* [Stream Quotes](Streaming/StreamQuotes.md)
* [Stream Trade Updates](Streaming/StreamTradeUpdates.md)

---

## 🚀 Quick Start

1. Configure your `config.json` with MT4 credentials and connection details.
2. Initialize an `MT4Account` and wrap it with `MT4Service`.
3. Run examples from `main.go` or call the `Show*` helpers.

```go
ctx := context.Background()
svc := mt4.NewMT4Service(account)

// Example: quick account & quote
svc.ShowAccountSummary(ctx)
svc.ShowQuote(ctx, "EURUSD")
```

---

## 🛠 Requirements

* Go 1.21+
* gRPC-Go
* Protobuf Go bindings
* VS Code / GoLand / LiteIDE

---

## 🧭 Navigation

* Sections above link **directly** to the markdown files in your repo.
* Methods are organized to mirror the **MT4 API structure**.
* Each **Overview** file contains explanations, best practices, and usage guidance.
