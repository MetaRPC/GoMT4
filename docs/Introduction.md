# Getting Started with MetaTrader 4 in Go

Welcome to the **MetaRPC MT4 Go Documentation** — your guide to integrating with **MetaTrader 4** using **Go** and **gRPC**.

This documentation will help you:

* 📘 Explore all available **account, market, and trading methods**
* 💡 Learn from **Go usage examples** with context and timeout handling
* 🔁 Work with **real-time streaming** for quotes, orders, and trades
* ⚙️ Understand all **input/output types** such as `OpenedOrderInfo`, `QuoteData`, and enums like `OrderSendOperationType`

**How it fits together:** your Go program uses the client in this repository (`examples/mt4`: `MT4Account` and `MT4Service`, import `github.com/MetaRPC/GoMT4/examples/mt4`, used from a clone of the repo) and the generated protobuf messages (`github.com/MetaRPC/GoMT4/package`). The client talks gRPC over TLS to the MetaRPC cloud API at `mt4.mrpc.pro:443`, which runs and manages the MT4 terminal for your account. You do not install or run MetaTrader or any server yourself.

---

## 📚 Main Sections

### Account

* [Show Account Summary](Account/ShowAccountSummary.md)

---

### Market Info

* **Section overview:** [Market_Info_Overview](Market_Info/Market_Info_Overview.md)
* [Show Quote](Market_Info/ShowQuote.md)
* [Show Quotes Many](Market_Info/ShowQuotesMany.md)
* [Show Quote History](Market_Info/ShowQuoteHistory.md)
* [Show Symbol Params](Market_Info/ShowSymbolParams.md)
* [Show Symbols](Market_Info/ShowSymbols.md)
* [Show All Symbols](Market_Info/ShowAllSymbols.md)
* [Show Tick Values](Market_Info/ShowTickValues.md)

---

### Order Operations ⚠️

* **Section overview:** [Order_Operations_Overview](Order_Operations/Order_Operations_Overview.md)
* [Show Opened Orders](Order_Operations/ShowOpenedOrders.md)
* [Show Opened Order Tickets](Order_Operations/ShowOpenedOrderTickets.md)
* [Show Orders History](Order_Operations/ShowOrdersHistory.md)
* [Show Order Send Example](Order_Operations/ShowOrderSendExample.md)
* [Show Order Close Example](Order_Operations/ShowOrderCloseExample.md)
* [Show Order Close By Example](Order_Operations/ShowOrderCloseByExample.md)
* [Show Order Delete Example](Order_Operations/ShowOrderDeleteExample.md)
* [Show Order Modify Example](Order_Operations/ShowOrderModifyExample.md)


---

### Streaming

* **Section overview:** [Streaming_Overview](Streaming/Streaming_Overview.md)
* [Stream Opened Order Profits](Streaming/StreamOpenedOrderProfits.md)
* [Stream Opened Order Tickets](Streaming/StreamOpenedOrderTickets.md)
* [Stream Quotes](Streaming/StreamQuotes.md)
* [Stream Trade Updates](Streaming/StreamTradeUpdates.md)
* [Stream Orders History (Example)](Streaming/StreamOrdersHistoryExample.md)
* [Stream Quote History (Example)](Streaming/StreamQuoteHistoryExample.md)

---

## 🧭 Navigation

* Sections above link **directly** to the markdown files in your repo.
* Methods are organized to mirror the **MT4 API structure**.
* Each **Overview** file contains explanations, best practices, and usage guidance.
