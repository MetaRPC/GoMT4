# GoMT4 Documentation – Skeleton

<div class="grid cards" markdown>

-   :material-rocket-launch: **Getting Started**
    ---
    Первый запуск, окружение, короткий путь.
    <br>
    [:octicons-play-16: Your Easy Start](Your_Easy_Start.md){ .md-button } 
    [:octicons-gear-16: Setup](setup.md){ .md-button-outline }

-   :material-diagram-project: **Architecture & Reliability**
    ---
    Как всё устроено, потоки данных, таймауты и реконнекты.
    <br>
    [:material-source-branch: Architecture](Architecture_DataFlow.md){ .md-button }
    [:material-refresh-auto: Reliability](ReTimeouts_Reconnects_Backoff.md){ .md-button-outline }

-   :material-chart-line: **Market Info**
    ---
    Котировки, параметры символов, история.
    <br>
    [:material-currency-usd: ShowQuote](Market_Info/ShowQuote.md){ .md-button }
    [:material-format-list-bulleted: QuotesMany](Market_Info/ShowQuotesMany.md){ .md-button-outline }

-   :material-shopping: **Orders**
    ---
    Открытие/модификация/закрытие, примеры.
    <br>
    [:material-cart-plus: Place Market](Cookbook/Orders/PlaceMarketOrder.md){ .md-button }
    [:material-content-save-edit: Modify](Cookbook/Orders/ModifyOrder.md){ .md-button-outline }

-   :material-broadcast: **Streaming**
    ---
    Тики, обновления ордеров, история потоками.
    <br>
    [:material-rss: Stream Quotes](Streaming/StreamQuotes.md){ .md-button }
    [:material-finance: Order Profits](Streaming/StreamOpenedOrderProfits.md){ .md-button-outline }

-   :material-book-open-page-variant: **Cookbook**
    ---
    Рецепты из реального кода.
    <br>
    [:material-notebook: Все рецепты](Cookbook/index.md){ .md-button }

-   :material-api: **API Reference**
    ---
    Типы, сообщения, стримы из .proto.
    <br>
    [:material-format-list-bulleted-type: Enums](API%20Reference/Enums.md){ .md-button }
    [:material-email-fast: Streaming](API%20Reference/Streaming.md){ .md-button-outline }

-   :material-shield-lock: **Ops**
    ---
    Производительность, секреты, логирование.
    <br>
    [:material-speedometer: Performance](Performance_Notes.md){ .md-button }
    [:material-file-lock: Security](Security_Secrets.md){ .md-button-outline }

</div>


---

## 📑 Table of Contents

- [Introduction](Introduction.md)
- [Setup & Environment](setup.md) 
- [Beginner Run Guide](Beginner_Run_Guide.md)
- [Your Easy Start](Your_Easy_Start.md)
- [Architecture & Data Flow](Architecture_DataFlow.md)
- [Reliability: Timeouts, Reconnects, Backoff](ReTimeouts_Reconnects_Backoff.md)
- [Troubleshooting & FAQ](Troubleshooting_FAQ.md)
- [Cookbook (Recipes)](Cookbook/index.md)

  **Account**
  - [Show Account Summary](Cookbook/Account/AccountSummary.md)
  - [Stream Opened Order Profits](Cookbook/Account/StreamOpenedOrderProfits.md)

  **Market Info**
  - [Get Quote](Cookbook/Market_Info/GetQuote.md)
  - [Get Multiple Quotes](Cookbook/Market_Info/GetMultipleQuotes.md)
  - [Stream Quotes](Cookbook/Market_Info/StreamQuotes.md)
  - [Symbol Params](Cookbook/Market_Info/SymbolParams.md)

  **Orders**
  - [Place Market Order](Cookbook/Orders/PlaceMarketOrder.md)
  - [Place Pending Order](Cookbook/Orders/PlacePendingOrder.md)
  - [Modify Order](Cookbook/Orders/ModifyOrder.md)
  - [Close Order](Cookbook/Orders/CloseOrder.md)
  - [Close By Orders](Cookbook/Orders/CloseByOrders.md)
  - [Delete Pending](Cookbook/Orders/DeletePending.md)
  - [History Orders](Cookbook/Orders/HistoryOrders.md)

  **Reliability & Connection**
  - [Handle Reconnect](Cookbook/Reliability_Connection/HandleReconnect.md)
  - [Unary Retries](Cookbook/Reliability_Connection/UnaryRetries.md)
  - [Health Check](Cookbook/Reliability_Connection/HealthCheck.md)

  **Utils & Helpers**
  - [Round Volume/Price](Cookbook/Utils_Helpers/RoundVolumePrice.md)
  - [Ensure Symbol Visible](Cookbook/Utils_Helpers/EnsureSymbolVisible.md)
  - [Config Example](Cookbook/Utils_Helpers/ConfigExample.md)


- [CLI Usage (Playground)](cli_usage.md)
- [API Reference (Types & Enums)](API%20Reference/Overview.md)
  - [Enums](API%20Reference/Enums.md)
  - [Messages](API%20Reference/Messages.md)
  - [Streaming](API%20Reference/Streaming.md)
- [Performance Notes](Performance_Notes.md)
- [Security & Secrets](Security_Secrets.md)
- [Observability (Logs & Metrics)](Observability.md)
- [Glossary (MT4 Terms)](Glossary.md)

---

## Introduction

**What is GoMT4?**  
Small, pragmatic bridge between **MT4 Terminal** and your **Go** code via gRPC.

**Who is it for?**  
Beginners, algo developers, ops teams who need a scriptable MT4 integration.

**After reading you can:**
- Run a local demo.
- Connect to MT4.
- Subscribe to quotes.
- Place & close an order safely.

Quick links:  
👉 [Your Easy Start](Your_Easy_Start.md) · 👉 [Beginner Run Guide](Beginner_Run_Guide.md) · 👉 [CLI Usage](cli_usage.md)

---

## Setup & Environment

**Goal:** Run everything on Windows with MT4 installed.

### Prerequisites
- Windows 10/11, MT4 Terminal.
- Go ≥ 1.21.
- (Optional) VS Code + Go extension.

### Install (draft)
1. Clone the repo, run `go mod tidy` in `examples/`.
2. Configure credentials (see [Config Example](Cookbook/Utils_Helpers/ConfigExample.md)).
3. Open in VS Code and launch the debug profile.

---

## Architecture & Data Flow

Make the system non-magical:
- MT4 Terminal ⇄ **GoMT4 gRPC server** ⇄ client code.
- Lifecycles: connect → use → disconnect.
- Streams: Quotes, Orders; buffering & backpressure.
- Where retries/backoff kick in.

See details: [Architecture & Data Flow](Architecture_DataFlow.md)

---

## Reliability: Timeouts, Reconnects, Backoff

- `context.WithTimeout` for unary calls (2–5s baseline).
- Treat `io.EOF` as transient on streams; reconnect with jitter.
- Ensure cancelation closes goroutines.
- Don’t leak streams; add health checks.

See recipes:
- [Handle Reconnect](Cookbook/Reliability_Connection/HandleReconnect.md)
- [Unary Retries](Cookbook/Reliability_Connection/UnaryRetries.md)
- [Health Check](Cookbook/Reliability_Connection/HealthCheck.md)

---

## Troubleshooting & FAQ

- “**Symbol not found EURUSD**” → Try broker suffix `EURUSD.m`.
- “**Invalid volume**” → Respect `LotStep`/`LotMin` and round.
- “**Digits mismatch**” → Format prices using `Digits`.
- “**No connection**” → Firewall/UAC, terminal path, server reachability.

Full page: [Troubleshooting & FAQ](Troubleshooting_FAQ.md)

---

## Cookbook (Recipes)

Jump into ready-made snippets:

- Watchlist & Quotes → [Get Quote](Cookbook/Market_Info/GetQuote.md), [Multiple Quotes](Cookbook/Market_Info/GetMultipleQuotes.md), [Stream](Cookbook/Market_Info/StreamQuotes.md)  
- Place Order Safely → [Place Market](Cookbook/Orders/PlaceMarketOrder.md) / [Pending](Cookbook/Orders/PlacePendingOrder.md) / [Modify](Cookbook/Orders/ModifyOrder.md) / [Close](Cookbook/Orders/CloseOrder.md)  
- Compute PnL Correctly → [Symbol Params](Cookbook/Market_Info/SymbolParams.md)  
- Stream History to DB → [History Orders](Cookbook/Orders/HistoryOrders.md)

Full list: [Cookbook index](Cookbook/index.md)

---

## CLI Usage (Playground)

- Subscribe to quotes.
- Dump symbol params.
- Close orders by filter.

See: [CLI Usage](cli_usage.md)

---

## API Reference (Types & Enums)

Autogenerated types index with human-readable notes, units, ranges, and gotchas.

- [Overview](API%20Reference/Overview.md)
- [Enums](API%20Reference/Enums.md)
- [Messages](API%20Reference/Messages.md)
- [Streaming](API%20Reference/Streaming.md)

---

## Performance Notes

- Batch calls when possible; avoid per-tick RPCs.
- Track expected latencies; simple load test plan.

[Performance Notes](Performance_Notes.md)

---

## Security & Secrets

- `.env` handling (do **not** commit secrets).
- Windows credentials vault tips.

[Security & Secrets](Security_Secrets.md)

---

## Observability (Logs & Metrics)

- Log format and levels.
- Basic metrics: latency, reconnects, dropped ticks.

[Observability](Observability.md)

---

## Glossary (MT4 Terms)

- Digits, Point, TickSize, TickValue, Lot, ContractSize.
- Hedging vs Netting (MT5 nuances for reference).

[Glossary](Glossary.md)

📊 Whoosh-and-it-works: get a quote, open an order, subscribe to ticks.

???+ tip "Quick start: Get one quote (Go)"
    ```go
    // examples/mt4/MT4_service.go has ShowQuote(ctx, symbol string)
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    price, err := svc.ShowQuote(ctx, "EURUSD")
    if err != nil {
        log.Fatalf("quote error: %v", err)
    }
    fmt.Printf("EURUSD  Bid: %.5f  Ask: %.5f  Time: %s\n",
        price.Bid, price.Ask, price.Time.Format(time.RFC3339))
    ```

???+ example "Place market order safely (uses rounding & symbol params)"
    ```go
    // See Cookbook/Orders/PlaceMarketOrder.md and ShowOrderSendExample in MT4_service.go
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    ticket, err := svc.ShowOrderSendExample(ctx, "EURUSD")
    if err != nil {
        log.Fatalf("order send failed: %v", err)
    }
    fmt.Println("✅ Order ticket:", ticket)
    ```

???+ info "Stream live quotes (local terminal)"
    ```go
    // Based on StreamQuotes in MT4_service.go
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    updates, errCh := svc.StreamQuotes(ctx, []string{"EURUSD","GBPUSD"})
    for {
        select {
        case q := <-updates:
            fmt.Printf("[Tick] %s  Bid: %.5f  Ask: %.5f\n", q.Symbol, q.Bid, q.Ask)
        case err := <-errCh:
            log.Printf("stream error: %v (reconnect logic in cookbook)", err)
            return
        }
    }
    ```

