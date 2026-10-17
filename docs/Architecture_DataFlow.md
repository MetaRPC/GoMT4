# 🏗️ Architecture & Data Flow (GoMT4)

This section describes the overall structure of the GoMT4 project and how data flows between components.

---

##  General Diagram

```
   ┌──────────────────────────────────────┐
   │ 👩‍💻 Your Go program                   │
   │  examples/main.go or your own code   │
   │  MT4Service → MT4Account → pb stubs  │
   └──────────────────┬───────────────────┘
                      │ gRPC over TLS
                      │ metadata: apikey + id
                      ▼
   ┌──────────────────────────────────────┐
   │ ☁️ MetaRPC cloud API                  │
   │  mt4.mrpc.pro:443                    │
   │  (checks the API key, routes the     │
   │   call to the account's terminal)    │
   └──────────────────┬───────────────────┘
                      │
                      ▼
   ┌──────────────────────────────────────┐
   │ 💻 MT4 terminal managed by MetaRPC    │
   │  logged in to your account, talks    │
   │  to the broker, executes the call    │
   └──────────────────────────────────────┘

📄 config.json → login, password, server name, default symbol, API key for the example.
📦 package/    → Go module with the generated messages, enums and gRPC clients.
```

---

## ⚙️ Components

* **👩‍💻 Your program**
  `examples/main.go` (or your own code) creates an `MT4Account`, connects it, and calls methods on it or on `MT4Service`.

* **🧩 `MT4Account`** (`examples/mt4/MT4Account.go`)
  The client. Opens one TLS gRPC connection to `mt4.mrpc.pro:443` (`NewMT4Account`), connects a terminal (`ConnectByServerName` / `ConnectByHostPort`), and wraps every RPC in a typed Go method (`AccountSummary`, `Quote`, `OrderSend`, `OnSymbolTick`, …). Unary calls go through `ExecuteWithReconnect`, streams through `ExecuteStreamWithReconnect` (retry with backoff on `Unavailable` and on terminal-not-found errors).

* **🧰 `MT4Service`** (`examples/mt4/MT4_service.go`)
  Demo layer on top of `MT4Account`: `Show...` methods print results, `Stream...` methods print a few stream messages and stop.

* **📦 `package/`** (`github.com/MetaRPC/GoMT4/package`)
  Generated protobuf messages, enums and gRPC clients (`ConnectionClient`, `AccountHelperClient`, `TradingHelperClient`, `MarketInfoClient`, `SubscriptionServiceClient`, …).

* **☁️ MetaRPC cloud API** (`mt4.mrpc.pro:443`)
  Authenticates the call by its `apikey` metadata (no key = trial; an unknown key is refused with `API_KEY_INVALID`) and forwards it to the terminal named by the `id` metadata.

* **💻 MT4 terminal**
  Started and managed by MetaRPC for your account. Connects to the broker, holds history, executes trading operations.

* **📄 `examples/config/config.json`**
  Login, password, server name, default symbol and API key used by the example.

---

## 🔀 Data Flow

1. **🔌 Connect**
   `NewMT4Account(user, password, "", uuid.Nil, apiKey)` dials `mt4.mrpc.pro:443` (the default when the server argument is empty). `ConnectByServerName` sends `ConnectEx` with your credentials; MetaRPC brings up a terminal logged in to the account and returns its id (`TerminalInstanceGuid`). The client then checks it with an `AccountSummary` call.

2. **📡 RPC call**
   Every call carries two metadata headers: `apikey` (your key, `TRIAL` if none) and `id` (the terminal id).

3. **💻 MT4 terminal**
   Executes the operation (e.g., get a quote or send an order) and returns the result.

4. **⬅️ Return path**
   The reply comes back as a protobuf message. API-level failures arrive in its `Error` field (`ErrorCode`, `ErrorMessage`, `MqlErrorCode`); the client turns them into Go errors.

5. **🔄 Streaming calls**
   Subscriptions (`OnSymbolTick`, `OnTrade`, `OnOpenedOrdersProfit`, `OnOpenedOrdersTickets`) are server-streaming RPCs. The client returns two channels (data, error) and keeps reading until the context is cancelled.

6. **📴 Disconnect**
   `Disconnect` sends `Disconnect` to the terminal (with `true`, the terminal instance is deleted) and closes the gRPC connection.

---

## 🆔 Terminal id

* Without an explicit id, `NewMT4Account` uses a deterministic id computed from login + password (`ComputeDeterministicTerminalId`, the same formula the server's `GetId` uses).
* After a successful connect the client uses the id returned by the server.
* A call only reaches a terminal that is connected with that id; otherwise it fails with a terminal-not-found error.

---

## ✨ Highlights

* Endpoint: `mt4.mrpc.pro:443` (gRPC over TLS).
* Streaming methods allow real-time subscriptions.
* Supported domains: account, orders, history, quotes, symbols.
* The API itself is defined by the server (`.proto` files); the Go stubs in `package/` are generated from them.

---

## 🛠️ Developer Notes

* Main entry logic: `examples/main.go`.
* Account config: `examples/config/config.json`.
* Client code: `examples/mt4/MT4Account.go`, `examples/mt4/MT4_service.go`.
* Debugging: use VS Code with `launch.json` (see [Setup](setup.md)).
