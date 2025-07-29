# Getting Started with MetaTrader 4 in Go

Welcome to the **MetaRPC MT4 Go Documentation** — your guide to building integrations with **MetaTrader 4** using **Go** and **gRPC**.

This documentation provides everything you need to:

This documentation provides everything you need to:

- 📘 Understand all available trading, account, and market methods  
- 💡 See **Go usage examples** with context, timeouts, and channels  
- 🔁 Work with **real-time streaming** of quotes, trades, profits, and tickets  
- ⚙️ Explore all **input/output types** including OrderInfo, TradeInfo, QuoteData, and enums like `ENUM_ORDER_TYPE_TF`

---

## 📚 Main Sections

* **Account Information** — balance, equity, margin, leverage, and more
* **Order Operations** — send, modify, close orders and track history
* **Market Info** — available instruments, trading conditions, tick sizes
* **Streaming** — subscribe to real-time updates for trades, profits, and prices

---

## 🚀 Quick Start

To get started with Go + MetaTrader 4:

1. **Configure your `config.json`** with MT4 credentials and connection details.
2. Use the `MT4Account` or `MT4Service` structs to access functionality.
3. Run method examples via `main.go` or through helper files like `Show*.go`.

---

## 🛠 Requirements

* Go 1.20+
* gRPC-Go
* Protobuf-generated Go bindings for MT4 `.proto` definitions
* LiteIDE, VS Code, or GoLand for editing

---

With this documentation, you can:

* Build terminal dashboards
* Automate trade flows
* Create monitoring services
* Analyze real-time market data

Ready to trade? Let’s Go 🟢
