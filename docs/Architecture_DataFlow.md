# 🏗️ Architecture & Data Flow (GoMT4)

This section describes the overall structure of the GoMT4 project and how data flows between components.

---

##  General Diagram

```
          ┌─────────────────────────────┐
          │        💻 MT4 Terminal      │
          │ (local, broker connection, │
          │  quotes, orders handling)  │
          └──────────────┬──────────────┘
                         │
                         ▼
          ┌─────────────────────────────┐
          │   🚀 GoMT4 gRPC Server      │
          │ (examples/main.go + pb API)│
          └───────┬───────────┬────────┘
                  │           │
                  ▼           ▼
       ┌────────────────┐   ┌───────────────────┐
       │ 👩‍💻 Client Apps │   │ 🔄 Streaming      │
       │ (Go, C#, etc.) │   │ Handlers (quotes, │
       │                │   │ orders, updates)  │
       └────────────────┘   └───────────────────┘

📄 Config.json → used by GoMT4 to log into account and select symbol.  
📦 pb module  → external Go module with generated structures and services.
```

---

## ⚙️ Components

* **💻 MT4 Terminal**
  Runs locally. Connects to broker, stores history, handles trading operations.

* **🚀 GoMT4 gRPC Server**
  Proxy between MT4 and external apps. Implemented in `examples/main.go` and code that uses pb module.

* **📦 pb module**
  Contains generated structures and services from `.proto` files (`mrpc-proto` repository).

* **📂 examples/**
  Contains entrypoint and usage examples.

* **📑 docs/**
  Documentation for each feature.

* **📄 config.json**
  Stores login, password, server and default symbol.

---

## 🔀 Data Flow

1. **📡 RPC call**
   A client (Go, C#, etc.) sends an RPC to the gRPC server (`127.0.0.1:50051`).

2. **⚙️ GoMT4 server**
   Receives the request, translates it into MT4 calls, processes the response.

3. **💻 MT4 Terminal**
   Executes the operation (e.g., get a quote or send an order) and returns the result.

4. **⬅️ Return path**
   Result goes back to the client through GoMT4.

5. **🔄 Streaming calls**
   If the client subscribed (quotes, orders updates), GoMT4 keeps the connection open and pushes updates in real time.

---

## ✨ Highlights

* Default gRPC port: `127.0.0.1:50051`.
* To extend the API, edit `.proto` files in `mrpc-proto` repo.
* Streaming methods allow real-time subscriptions.
* Supported domains: account, orders, history, quotes.

---

## 🛠️ Developer Notes

* Main entry logic: `examples/main.go`.
* Account config: `examples/config/config.json`.
* New functions: edit `.proto` and rebuild pb module.
* Debugging: use VS Code with `launch.json`.
