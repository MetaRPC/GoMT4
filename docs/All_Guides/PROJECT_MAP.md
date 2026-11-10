# 🗺️ GoMT4 Project Map

> The whole project as one annotated tree, with links. **Links point to section *overviews*, never to individual method pages** — open an overview to drill down into its methods.

---

## 🧱 The three layers

GoMT4 is built in three decoupled layers (composition, not inheritance). Your code can enter at any level.

| Layer | Source | What it does | Overview |
| ----- | ------ | ------------ | -------- |
| **3 · MT4Sugar** | `mt4/mt4sugar.go` | Convenience on top of Service: risk-based sizing, pip-based SL/TP, price/volume helpers, trailing & breakeven, watchers, order queries. | [MT4Sugar overview](../MT4Sugar/Overview.md) |
| **2 · MT4Service** | `mt4/mt4service.go` | Typed facade: unwrapped returns (slices/scalars) and simpler arguments over MT4Account. What most code uses. | [MT4Service overview](../MT4Service/Overview.md) |
| **1 · MT4Account** | `mt4/MT4Account.go` | Low-level gRPC: builds each proto `…Request`, attaches auth metadata, returns `(*pb.…Data, error)`; connection, retry & streaming. Everything builds on it. | [MT4Account overview](../MT4Account/Overview.md) |

```text
Your code → MT4Sugar → MT4Service → MT4Account → gRPC → MetaRPC gateway → MT4 terminal
```

Mix freely: `sugar.Service()` borrows the Service, `svc.Account()` borrows the Account.

---

## 🌳 Repository tree

GoMT4 is a Go **workspace** (`go.work`) of three modules:

- **`GoMT4/`** — repo root (`go.work` ties the modules together)
  - **`package/`** — module `github.com/MetaRPC/GoMT4/package`
    - `*.pb.go` — vendored protobuf / gRPC stubs (`AccountHelper`, `TradingHelper`, `MarketInfo`, `Subscriptions`, …)
  - **`mt4/`** — module `github.com/MetaRPC/GoMT4/mt4` — the SDK
    - `MT4Account.go` — **Layer 1** → [MT4Account overview](../MT4Account/Overview.md)
    - `mt4service.go` — **Layer 2** → [MT4Service overview](../MT4Service/Overview.md)
    - `mt4sugar.go` — **Layer 3** → [MT4Sugar overview](../MT4Sugar/Overview.md)
    - **`strategies/`** — orchestrators + presets
      - `presets.go` — `StrategyPreset`, `RiskPreset`, catalog, `sizeLot`, `manage`
      - `market_one_shot.go` · `bracket_trailing.go` · `pending_bracket.go`
  - **`examples/`** — module `github.com/MetaRPC/GoMT4/examples` — the runnable app
    - `main.go` — the **interactive menu** (connects once, dispatches to the demos)
    - `config/` — `config.json` loader
    - `common/` — shared demo helpers (`ReadStream`, `Ptr`, trading flag)
    - **`1_lowlevel/`** — MT4Account demos (`general`, `trading`, `streaming`)
    - **`2_service/`** — MT4Service demo
    - **`3_sugar/`** — MT4Sugar demos (`basics`, `trading`, `positions`)
    - **`4_strategies/`** — orchestrator demo
    - **`protobuf_inspector/`** — the offline [type explorer](PROTOBUF_INSPECTOR_GUIDE.md)

---

## 📚 Documentation tree

- **`docs/`**
  - [Home](../index.md) — landing page
  - **`All_Guides/`** — the guides ([Getting Started](GETTING_STARTED.md), [Your First Project](Your_First_Project.md), [MT4 for Beginners](MT4_For_Beginners.md), [Glossary](GLOSSARY.md), [Concurrency Model](SYNC_VS_ASYNC.md), [gRPC Streaming](GRPC_STREAM_MANAGEMENT.md), [ENUMs](ENUMS_USAGE_REFERENCE.md), [Return Codes](RETURN_CODES_REFERENCE.md), [Protobuf Inspector](PROTOBUF_INSPECTOR_GUIDE.md))
  - **`MT4Account/`** · **`MT4Service/`** · **`MT4Sugar/`** — the three layer overviews
  - `styles/` · `javascripts/` · `includes/` — theme assets

---

## 🎯 Which layer should you use?

- **Most trading apps → [MT4Sugar](../MT4Sugar/Overview.md).** Risk sizing, pip math and bulk ops out of the box.
- **Custom execution → [MT4Service](../MT4Service/Overview.md).** Clean returns when you manage your own request shapes.
- **Protocol-level control → [MT4Account](../MT4Account/Overview.md).** Raw proto structures and streams.
