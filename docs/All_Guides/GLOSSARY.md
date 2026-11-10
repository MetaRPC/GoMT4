# GoMT4 Glossary

> The terms you'll meet across these docs — architecture, MT4 trading, Go/concurrency, and gRPC — each with a one-paragraph definition and where it shows up in the SDK.

---

## 🏗️ Architecture

### Three-layer architecture
GoMT4 is split into three stacked layers, each wrapping the one below:

| Layer | Type | Role |
|-------|------|------|
| [`MT4Account`](../MT4Account/Overview.md) | low | Raw gRPC — build `…Request`, get `(*…Data, error)`. |
| [`MT4Service`](../MT4Service/Overview.md) | mid | Clean returns (unwrapped slices/scalars), simpler args. |
| [`MT4Sugar`](../MT4Sugar/Overview.md) | high | One-liners, pip math, risk sizing. |

### MT4Account
Layer 1. Owns the gRPC connection, TLS, request metadata (session id + API key), retry and reconnect. Methods build the raw protobuf request and return `(*pb.…Data, error)`. Use it for full control.

### MT4Service
Layer 2. The same 25 operations with **simpler arguments** and **clean returns** — a Go slice or scalar instead of a wrapper struct (`OpenedOrders` → `[]*OpenedOrderInfo`). Streams are renamed `Stream*`.

### MT4Sugar
Layer 3. ~45 convenience methods: buy/sell by side, pip conversions, risk-based lot sizing, trailing/breakeven, watchers, order info. Returns the plain `pb` types and simple values.

### Orchestrator / Preset
An **orchestrator** (`Run*` in [`mt4/strategies/`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/strategies/presets.go)) chains Sugar calls into a full lifecycle (size → enter → manage → close). A **preset** configures it: `StrategyPreset` (what to trade) + `RiskPreset` (how much to risk).

---

## 📈 MT4 trading terms

### Order vs Position
In MT4 there's no separate "position" object: an **open market order** *is* the position. All open orders — market and pending — live in one list, told apart by `GetOrderType()`. `MT4Sugar.Positions()` filters to the market ones; `Orders()` returns everything.

### Order type
An enum **0 Buy, 1 Sell, 2 BuyLimit, 3 SellLimit, 4 BuyStop, 5 SellStop** (`pb.OpenedOrderType`, MQL4 `OP_*` numbering). See [ENUMs Reference](ENUMS_USAGE_REFERENCE.md).

### Pending order
An order that triggers when price reaches a level, rather than filling now: **limit** (buy below / sell above the market) and **stop** (buy above / sell below). Placed with `BuyLimit`/`SellLimit`/`BuyStop`/`SellStop`; deleted with `Close`/`CancelPendings`.

### Lot / Volume
Trade size. 1.0 lot = contract size units (100,000 for FX majors). Brokers constrain it with `VolumeMin`, `VolumeMax`, `VolumeStep` — always pass sizes through `MT4Sugar.NormalizeLot`.

### Point vs Pip
A **point** is the smallest price step (`GetPoint()`, e.g. `0.00001` on a 5-digit pair). A **pip** is the conventional unit: `point × 10` on 3/5-digit symbols, else `point`. `MT4Sugar.PipSize` computes it; `PipsToPrice` / `PriceToPips` convert.

### Digits
Price decimal places (`5` for EURUSD, `3` for USDJPY). Round prices to this with `MT4Sugar.NormalizePrice` to avoid "invalid price" (129).

### Stop Loss (SL) / Take Profit (TP)
Protective exit prices attached to an order — SL caps loss, TP locks profit. Set them as absolute prices (`ModifySLTPByPrice`) or as pip distances from entry (`ModifySLTPByPips`). Must respect `TradeStopsLevel` or the broker rejects with 130.

### Spread
`ask − bid`. In pips via `MT4Sugar.SpreadPips`. Wide spreads (news, low liquidity) are a common reason to skip an entry.

### Swap
Interest credited/debited for holding a position overnight. Part of an order's net result: `GetProfit() + GetSwap() + GetCommision()`.

### Margin / Equity / Free margin / Margin level
**Margin** is collateral reserved by open positions; **equity** = balance + floating P/L; **free margin** = equity − margin; **margin level** = equity / margin × 100. Read from `AccountSummaryData` (`AccountSummary`).

### Magic number
An integer tag you stamp on orders (`SetDefaults(magic, …)`) to identify which bot/strategy owns them. Filter with `OrdersByMagic`.

### Hedging
Holding opposite positions on the same symbol at once. **MT4 is always hedging.** Net two opposite orders with `CloseBy` to save a spread.

### Trailing stop / Breakeven
**Trailing** ratchets the SL behind price as it moves in your favor (`TrailStop`, only ever tightening). **Breakeven** moves the SL to entry (±a few pips) once in profit (`AutoBreakeven`, never loosening). Both are one-shot — call on a timer.

---

## ⚙️ Go & concurrency terms

### `context.Context`
Every SDK method takes a `ctx`. It carries the **deadline** and the **cancel** signal — `context.WithTimeout` bounds a call, `WithCancel` stops a stream. Passing `nil` is accepted (a short default timeout is added). See [The Concurrency Model](SYNC_VS_ASYNC.md).

### Goroutine
A lightweight thread. Go has no `async`/`await`; you run work concurrently by launching `go func(){ … }()`. The gRPC client is safe to share across goroutines.

### Channel (`<-chan`)
A typed pipe for passing values between goroutines. Streaming methods return `(<-chan *pb.…Data, <-chan error)`; you read frames with `select`.

### `error`
Every call returns `(value, error)`. A non-nil `error` means failure — see [Return Codes](RETURN_CODES_REFERENCE.md). There is no separate status code to inspect.

### Server stream
A long-lived subscription that pushes frames. `MT4Account` returns the channel pair directly; `MT4Sugar` wraps three of them as `SubscribeQuotes` / `WatchProfit` / `WatchTickets`. Cancel the context to unsubscribe.

### `pb` (generated types)
The protobuf messages and enums, generated into package `github.com/MetaRPC/GoMT4/package` (imported as `pb`). `MT4Account` works with these directly; every field has a `GetXxx()` getter.

---

## 🔌 gRPC & protocol terms

### gRPC
The high-performance RPC framework the SDK speaks over HTTP/2 + TLS to `mt4.mrpc.pro:443`. Supports both unary (request/response) and server-streaming calls.

### Protobuf (protocol buffers)
The binary message format for requests/replies. The generated Go types live in [`github.com/MetaRPC/GoMT4/package`](https://github.com/MetaRPC/GoMT4/tree/main/package) (via `protoc-gen-go`). Explore them with the [Protobuf Inspector](PROTOBUF_INSPECTOR_GUIDE.md).

### `…Request` / `…Reply` / `…Data`
The call shape: you send a `…Request`, the server answers a `…Reply` whose `response` is a `oneof { Data, Error }`. `MT4Account` unwraps the `Data` payload (`…Data`) or turns `Error` into a Go `error`.

### Session id (terminal GUID)
A per-terminal identifier. `NewMT4Account(…, uuid.Nil, …)` starts with a deterministic id computed from login+password (`ComputeDeterministicTerminalId`); after `ConnectByServerName` the SDK uses the `TerminalInstanceGuid` the server returns, and attaches it (with the API key) to every later call.

### API key (`MRPC_API_KEY`)
Your MetaRPC authentication token, sent on every call. Provide it to `NewMT4Account` (last argument), via the `MRPC_API_KEY` env var, or `ApiKey` in `config.json`. Get one at [mrpc.pro/my](https://mrpc.pro/my).

### Config (`config.json`)
The example loader ([`examples/config/config.go`](https://github.com/MetaRPC/GoMT4/blob/main/examples/config/config.go)): reads `config.json` (fields `Login`, `Password`, `Server`, `DefaultSymbol`, `ApiKey`) and falls back to `MRPC_API_KEY`; `Login: 0` auto-provisions a demo account.

---

## 🧰 SDK conveniences

| Term | What it is |
|------|-----------|
| **Auto-normalization** | `NormalizeLot` / `NormalizePrice` round to the symbol's step/digits so the broker accepts the order. |
| **Risk-based sizing** | `CalcLotByRisk(ctx, symbol, cash, slPips)` → the normalized lot that risks `cash` over the stop. |
| **Pip helpers** | `PipSize`, `PipsToPrice`, `PriceToPips`, `SpreadPips`. |
| **Watchers / Waiters** | Watchers stream values onto a channel; waiters poll until a condition holds (`WaitPrice` / `WaitFilled` / `WaitClosed`). |
| **`EnsureSymbol`** | Select a symbol and wait for a usable two-sided quote before trading it. |

---

**See also:** [Getting Started](GETTING_STARTED.md) · [Project Map](PROJECT_MAP.md) · [ENUMs Reference](ENUMS_USAGE_REFERENCE.md) · [Return Codes](RETURN_CODES_REFERENCE.md)
