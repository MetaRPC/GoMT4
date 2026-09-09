# Getting Started with GoMT4

> **Quick Setup & Overview** - Start automating your trading with GoMT4 (Go for MetaTrader 4) in minutes.

---

## 🎯 What is GoMT4?

**GoMT4** is an industrial-grade, type-safe Go client library for interacting with **MetaTrader 4** terminals via high-performance **gRPC**. It eliminates complex C++ DLL wrappers and provides direct, reliable programmatic trading.

### 🌟 Key Advantages
- 🚀 **High Throughput**: Native gRPC streaming for sub-millisecond price ticks and trade execution.
- 🛡️ **Three-Layer Architecture**: Low-level gRPC (`MT4Account`), typed wrapper methods (`MT4Service`), and high-level convenience (`MT4Sugar`).
- 🔄 **Resilient Connection**: Auto-reconnect, exponential backoff, and transparent channel healing.
- 💼 **Production Ready**: Fully verified across institutional accounts, hedge fund systems, and automated retail bots.

---

## 📦 Installation

Install GoMT4 via your standard Go package manager:

```bash
go get github.com/MetaRPC/GoMT4
```

---

## 🔑 API Key & Authentication

Connecting to MetaRPC production endpoints (`mt4.mrpc.pro:443`) requires an API key:

1. **Sign Up**: Create an account for free at [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. **Generate API Key**: In your MetaRPC Portal dashboard at [https://mrpc.pro/my](https://mrpc.pro/my), go to **API Keys** to generate and copy your personal API token.
3. **Configure Connection**: Pass your API key / token along with the server address (`mt4.mrpc.pro:443`) in your connection settings.

---


---

## 🆔 Automatic Account ID & Authentication

MetaRPC endpoints require two credentials for all terminal operations:
1. **`APIKey`**: Your personal authentication token from [https://mrpc.pro/my](https://mrpc.pro/my) (obtained by registering at [https://mrpc.pro/signup](https://mrpc.pro/signup)). Sent in the `APIKey` header.
2. **`id`**: A deterministic account GUID derived from your MetaTrader `user` (login number) and `password`.

> 💡 **Seamless Automation**: You do not need to call `GetId` manually. The SDK automatically derives your deterministic account ID from your credentials upon initialization and attaches both the `id` and `APIKey` headers to all requests and streaming subscriptions.

## 🔌 Minimal Connection Example

Here is how easy it is to initialize `MT4Account`, connect to your MetaTrader terminal, and retrieve your account balance:

```
client, err := mt.NewMT4Account(user, password, grpcServer)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err = client.ConnectByServerName(ctx, serverName, "EURUSD")
summary, err := client.AccountSummary(ctx)
fmt.Printf("Balance: %.2f, Equity: %.2f\n", summary.Balance, summary.Equity)
```

---

## 🗺️ Documentation Road Map

To get the most out of GoMT4, follow this suggested reading order:

1. 🚀 **[Your First Project](Your_First_Project.md)** - Build and run a working project in 10 minutes.
2. 🗺️ **[Project Map](PROJECT_MAP.md)** - Understand the 3 architectural layers and interaction flow.
3. 📖 **[Glossary](GLOSSARY.md)** - Essential MetaTrader 4 and algorithmic trading terminology.
4. 🐣 **[MT4 for Beginners](MT4_For_Beginners.md)** - Step-by-step terminal setup and demo account guide.
5. 📡 **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)** - Subscribe to ticks, trades, DOM, and position updates.
6. 📊 **[Return Codes](RETURN_CODES_REFERENCE.md)** - Complete reference of broker and gateway return codes.
