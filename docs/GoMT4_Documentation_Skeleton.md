# GoMT4 The documentary BASE

<div class="grid cards" markdown>

-   :material-rocket-launch: **Getting Started**
    ---
    First launch, environment, shortcut.  
    <br>
    [:octicons-play-16: Your Easy Start](Your_Easy_Start.md){ .md-button } 
    [:octicons-gear-16: Setup](setup.md){ .md-button-outline }

-   :material-source-branch: **Architecture**
    ---
    How everything works, data flows, timeouts, and reconnects.  
    <br>
    [:material-source-branch: Architecture](Architecture_DataFlow.md){ .md-button }
    [:material-refresh-auto: Reliability](ReTimeouts_Reconnects_Backoff.md){ .md-button-outline }

-   :material-chart-line: **Market Info**
    ---
    Quotes, symbol parameters, history.  
    <br>
    [:material-currency-usd: ShowQuote](Market_Info/ShowQuote.md){ .md-button }
    [:material-format-list-bulleted: QuotesMany](Market_Info/ShowQuotesMany.md){ .md-button-outline }

-   :material-shopping: **Orders**
    ---
    Opening/modification/closing, examples.  
    <br>
    [:material-cart-plus: Place Market](Cookbook/Orders/PlaceMarketOrder.md){ .md-button }
    [:material-content-save-edit: Modify](Cookbook/Orders/ModifyOrder.md){ .md-button-outline }

-   :material-broadcast: **Streaming**
    ---
    Ticks, order updates, flow history.  
    <br>
    [:material-rss: Stream Quotes](Streaming/StreamQuotes.md){ .md-button }
    [:material-finance: Order Profits](Streaming/StreamOpenedOrderProfits.md){ .md-button-outline }

-   :material-book-open-page-variant: **Cookbook**
    ---
    Recipes from real code.  
    <br>
    [:material-notebook: All Recipes](Cookbook/index.md){ .md-button }

-   :material-api: **API Reference**
    ---
    Types, messages, and streams from .proto.  
    <br>
    [:material-format-list-bulleted-type: Enums](API%20Reference/Enums.md){ .md-button }
    [:material-email-fast: Streaming](API%20Reference/Streaming.md){ .md-button-outline }

-   :material-shield-lock: **Ops**
    ---
    Performance, secrets, logging.  
    <br>
    [:material-speedometer: Performance](Performance_Notes.md){ .md-button }
    [:material-file-lock: Security](Security_Secrets.md){ .md-button-outline }

</div>





---

## 📑 Table of Contents

* [ Introduction](Introduction.md)

* [ Setup & Environment](setup.md)

* [ Beginner Run Guide](Beginner_Run_Guide.md)

* [ Your Easy Start](Your_Easy_Start.md)

* [ Architecture & Data Flow](Architecture_DataFlow.md)

* [ Reliability: Timeouts, Reconnects, Backoff](ReTimeouts_Reconnects_Backoff.md)

* [ Troubleshooting & FAQ](Troubleshooting_FAQ.md)

* [Cookbook (Recipes)](Cookbook/index.md)

  **📊 Account**

  * [📄 Show Account Summary](Cookbook/Account/AccountSummary.md)
  * [💹 Stream Opened Order Profits](Cookbook/Account/StreamOpenedOrderProfits.md)

  **📈 Market Info**

  * [💱 Get Quote](Cookbook/Market_Info/GetQuote.md)
  * [📊 Get Multiple Quotes](Cookbook/Market_Info/GetMultipleQuotes.md)
  * [🔄 Stream Quotes](Cookbook/Market_Info/StreamQuotes.md)
  * [ℹ️ Symbol Params](Cookbook/Market_Info/SymbolParams.md)

  **📦 Orders**

  * [🛒 Place Market Order](Cookbook/Orders/PlaceMarketOrder.md)
  * [📌 Place Pending Order](Cookbook/Orders/PlacePendingOrder.md)
  * [✏️ Modify Order](Cookbook/Orders/ModifyOrder.md)
  * [❌ Close Order](Cookbook/Orders/CloseOrder.md)
  * [🔗 Close By Orders](Cookbook/Orders/CloseByOrders.md)
  * [🗑️ Delete Pending](Cookbook/Orders/DeletePending.md)
  * [📜 History Orders](Cookbook/Orders/HistoryOrders.md)

  **🔧 Reliability & Connection**

  * [♻️ Handle Reconnect](Cookbook/Reliability_Connection/HandleReconnect.md)
  * [🔁 Unary Retries](Cookbook/Reliability_Connection/UnaryRetries.md)
  * [❤️ Health Check](Cookbook/Reliability_Connection/HealthCheck.md)

  **🧰 Utils & Helpers**

  * [ Round Volume/Price](Cookbook/Utils_Helpers/RoundVolumePrice.md)
  * [ Ensure Symbol Visible](Cookbook/Utils_Helpers/EnsureSymbolVisible.md)
  * [ Config Example](Cookbook/Utils_Helpers/ConfigExample.md)

* [ CLI Usage (Playground)](cli_usage.md)

* [ API Reference (Types & Enums)](API%20Reference/Overview.md)

  * [ Enums](API%20Reference/Enums.md)
  * [ Messages](API%20Reference/Messages.md)
  * [ Streaming](API%20Reference/Streaming.md)

* [ Performance Notes](Performance_Notes.md)

* [ Security & Secrets](Security_Secrets.md)

* [ Observability (Logs & Metrics)](Observability.md)

* [ Glossary (MT4 Terms)](Glossary.md)

---

## ✨ Introduction

**What is GoMT4?**
A small, pragmatic Go client for **MetaTrader 4** through the MetaRPC cloud API (`mt4.mrpc.pro:443`, gRPC over TLS). MetaRPC runs the MT4 terminal for your account; your Go code only needs the account credentials and an API key.

**Who is it for?**
Beginners , algo developers, ops teams who need a scriptable MT4 integration.

**After reading you can:**

* Run the example app.
* Connect to an MT4 account.
* Subscribe to quotes.
* Place & close an order safely.

Quick links:
👉 [Your Easy Start](Your_Easy_Start.md) · 👉 [Beginner Run Guide](Beginner_Run_Guide.md) · 👉 [CLI Usage](cli_usage.md)

---

## ⚙️ Setup & Environment

**Goal:** Run the example on any OS with Go; no MetaTrader installation needed.

### Prerequisites

*  Go ≥ 1.23, Git.
*  An MT4 account (login, password, broker server name) and a MetaRPC API key ([https://mrpc.pro/signup](https://mrpc.pro/signup)).
*  (Optional) VS Code + Go extension.

### Install

1. Clone the repo, run `go mod download` in `examples/`.
2. Configure credentials in `examples/config/config.json` (see [Config Example](Cookbook/Utils_Helpers/ConfigExample.md)).
3. `go run .` in `examples/`, or open in VS Code and launch the debug profile.

---

## 🏗️ Architecture & Data Flow

Make the system non-magical:

* 🧑‍💻 your Go code (`MT4Account` / `MT4Service`) ⇄ ☁️ **MetaRPC API** (`mt4.mrpc.pro:443`) ⇄ 💻 MT4 terminal run by MetaRPC.
* Lifecycles: connect → use → disconnect.
* Streams: ticks, trade events, order profits/tickets; consumer loops and cancellation.
* Where retries/backoff kick in.

See details: [Architecture & Data Flow](Architecture_DataFlow.md)

---

## 🔒 Reliability: Timeouts, Reconnects, Backoff

* `context.WithTimeout` for unary calls (2–5s baseline; the client adds 3s/5s/8s defaults when your `ctx` has none).
* Streams treat `io.EOF` / `Unavailable` as transient and reconnect with jitter (built in).
* Ensure cancelation closes goroutines.
* Don’t leak streams; add health checks.

See recipes:

* [♻️ Handle Reconnect](Cookbook/Reliability_Connection/HandleReconnect.md)
* [🔁 Unary Retries](Cookbook/Reliability_Connection/UnaryRetries.md)
* [❤️ Health Check](Cookbook/Reliability_Connection/HealthCheck.md)

---

## 🛠️ Troubleshooting & FAQ

* “**Symbol not found EURUSD**” → Try broker suffix `EURUSD.m`.
* “**Invalid volume**” → Respect `VolumeStep`/`VolumeMin` and round.
* “**Digits mismatch**” → Format prices using `Digits`.
* “**No connection**” → Credentials, API key, reachability of `mt4.mrpc.pro:443`.

Full page: [Troubleshooting & FAQ](Troubleshooting_FAQ.md)

---

## 📚 Cookbook (Recipes)

Jump into ready-made snippets:

* 📈 Watchlist & Quotes → [Get Quote](Cookbook/Market_Info/GetQuote.md), [Multiple Quotes](Cookbook/Market_Info/GetMultipleQuotes.md), [Stream](Cookbook/Market_Info/StreamQuotes.md)
* 🛒 Place Order Safely → [Place Market](Cookbook/Orders/PlaceMarketOrder.md) / [Pending](Cookbook/Orders/PlacePendingOrder.md) / [Modify](Cookbook/Orders/ModifyOrder.md) / [Close](Cookbook/Orders/CloseOrder.md)
* 💹 Compute PnL Correctly → [Symbol Params](Cookbook/Market_Info/SymbolParams.md)
* 🗄️ Stream History to DB → [History Orders](Cookbook/Orders/HistoryOrders.md)

Full list: [Cookbook index](Cookbook/index.md)

---

## 🖥️ CLI Usage (Playground)

* Subscribe to quotes 
* Dump symbol params 
* Send, modify and close orders with the example helpers 

See: [CLI Usage](cli_usage.md)

---

## 📖 API Reference (Types & Enums)

Types, enums and streams of the generated `pb` package (`github.com/MetaRPC/GoMT4/package`), with notes and gotchas ⚠️.

* [📖 Overview](API%20Reference/Overview.md)
* [🔢 Enums](API%20Reference/Enums.md)
* [📬 Messages](API%20Reference/Messages.md)
* [📡 Streaming](API%20Reference/Streaming.md)

---

## ⚡ Performance Notes

* Batch calls when possible 
* Avoid per-tick RPCs 
* Track expected latencies 
* Cache symbol metadata 

[Performance Notes](Performance_Notes.md)

---

## 🔑 Security & Secrets

* `.env` handling (do **not** commit secrets) 
* Windows credentials vault tips 

[Security & Secrets](Security_Secrets.md)

---

## 📊 Observability (Logs & Metrics)

* Log format and levels 
* Metrics: latency, reconnects, stream rates 

[Observability](Observability.md)

---

## 📘 Glossary (MT4 Terms)

*  Digits, Point, Lot, Volume Step, Contract Size, Stops Level
*  MT4 order types and their Go enums

[Glossary](Glossary.md)


📊 Whoosh-and-it-works: get a quote, open an order, subscribe to ticks.

???+ tip "Quick start: Get one quote (Go)"
    ```go
    // examples/mt4/MT4Account.go has Quote(ctx, symbol string) (*pb.QuoteData, error)
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    q, err := account.Quote(ctx, "EURUSD")
    if err != nil {
        log.Fatalf("quote error: %v", err)
    }
    fmt.Printf("EURUSD  Bid: %.5f  Ask: %.5f  Time: %s\n",
        q.GetBid(), q.GetAsk(), q.GetDateTime().AsTime().Format(time.RFC3339))
    ```

???+ example "Place a market order (rounding & symbol params: see the Cookbook)"
    ```go
    // See Cookbook/Orders/PlaceMarketOrder.md; ShowOrderSendExample in MT4_service.go is the printing demo
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    res, err := account.OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.10,
        nil, nil, nil, nil, nil, nil, nil) // market price, no slippage/SL/TP/comment/magic/expiration
    if err != nil {
        log.Fatalf("order send failed: %v", err)
    }
    fmt.Println("✅ Order ticket:", res.GetTicket())
    ```

???+ info "Stream live quotes"
    ```go
    // Based on StreamQuotes in MT4_service.go (which uses OnSymbolTick)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    updates, errCh := account.OnSymbolTick(ctx, []string{"EURUSD", "GBPUSD"})
    for {
        select {
        case u, ok := <-updates:
            if !ok {
                return
            }
            t := u.GetSymbolTick()
            fmt.Printf("[Tick] %s  Bid: %.5f  Ask: %.5f\n", t.GetSymbol(), t.GetBid(), t.GetAsk())
        case err := <-errCh:
            log.Printf("stream error: %v (reconnect logic in cookbook)", err)
            return
        }
    }
    ```
