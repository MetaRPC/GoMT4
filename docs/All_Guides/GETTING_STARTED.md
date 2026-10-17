# Getting Started with GoMT4

> **Quick Setup & Overview** - Start automating your trading with GoMT4 (Go for MetaTrader 4) in minutes.

---

## 🎯 What is GoMT4?

**GoMT4** is a Go client for **MetaTrader 4** accounts through the **MetaRPC** cloud API. Your program talks gRPC over TLS to `mt4.mrpc.pro:443`; MetaRPC runs and manages an MT4 terminal logged in to your account. You do not install MetaTrader or run any server yourself.

### 🌟 What you get
- 📦 **Generated API package** `github.com/MetaRPC/GoMT4/package` (`pb`): all messages, enums and gRPC clients.
- 🧩 **Client** in the repository's examples (`examples/mt4`): `MT4Account` (one typed Go method per call: account, quotes, symbols, orders, history, streams) and `MT4Service` (demo methods that print results).
- 🔄 **Retries built in**: unary calls and streams retry transient errors with exponential backoff and jitter.
- ▶️ **Runnable example** `examples/main.go` with a JSON config.

---

## 📦 Installation

The generated API package can be fetched with Go:

```bash
go get github.com/MetaRPC/GoMT4/package
```

The `MT4Account` / `MT4Service` client is part of the examples module and is used from a clone:

```bash
git clone https://github.com/MetaRPC/GoMT4.git
cd GoMT4/examples
go run .
```

---

## 🔑 API Key & Authentication

Connecting to MetaRPC production endpoints (`mt4.mrpc.pro:443`) uses an API key:

1. **Sign Up**: Create an account for free at [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. **Generate API Key**: In your MetaRPC Portal dashboard at [https://mrpc.pro/my](https://mrpc.pro/my), go to **API Keys** to generate and copy your personal API token.
3. **Configure Connection**: Pass the key to `NewMT4Account` (last argument), set `MRPC_API_KEY`, or put it into `examples/config/config.json` → `ApiKey`.

Without a key, calls run as a trial (`TRIAL`); an unknown key is refused with `API_KEY_INVALID`.

---

## 🆔 Account ID & Authentication

Every call carries two gRPC metadata headers:
1. **`apikey`**: your API key (`TRIAL` when none is given).
2. **`id`**: the id of the terminal that serves your account.

> 💡 **How the client handles the id**: `NewMT4Account(..., uuid.Nil, ...)` starts with a deterministic id computed from login + password (`ComputeDeterministicTerminalId`, the same value the server's `GetId` returns). `ConnectByServerName` / `ConnectByHostPort` connect the terminal, and the client then uses the `TerminalInstanceGuid` the server returns. A call only reaches a terminal that is connected under that id, so connect first.

## 🔌 Minimal Connection Example

Initialize `MT4Account`, connect to your MetaTrader account, and read the balance (run inside a clone, e.g. as a file in `examples/`):

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MetaRPC/GoMT4/examples/mt4"
	"github.com/google/uuid"
)

func main() {
	user := uint64(12345678)
	password := "your-password"
	serverName := "MetaQuotes-Demo" // broker server name as shown in MetaTrader
	apiKey := "your-api-key"

	// "" = default endpoint mt4.mrpc.pro:443; uuid.Nil = deterministic id from login+password
	client, err := mt4.NewMT4Account(user, password, "", uuid.Nil, apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := client.ConnectByServerName(ctx, serverName, "EURUSD", true, 30); err != nil {
		log.Fatal(err)
	}
	summary, err := client.AccountSummary(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %.2f, Equity: %.2f\n", summary.GetAccountBalance(), summary.GetAccountEquity())
}
```

---

## 🗺️ Documentation Road Map

To get the most out of GoMT4, follow this suggested reading order:

1. 🚀 **[Your First Project](Your_First_Project.md)** - Build and run a working project in 10 minutes (generated clients only).
2. 🚦 **[Beginner Run Guide](../Beginner_Run_Guide.md)** - Run the example app from a clone.
3. 🗺️ **[Architecture & Data Flow](../Architecture_DataFlow.md)** - How the client, the MetaRPC API and the terminal fit together.
4. 📖 **[Glossary](../Glossary.md)** - Essential MetaTrader 4 terminology.
5. 📡 **[Streaming](../Streaming/Streaming_Overview.md)** - Subscribe to ticks, trades, and position updates.
6. 📊 **[API Reference](../API%20Reference/Overview.md)** - Messages, enums and error codes.
