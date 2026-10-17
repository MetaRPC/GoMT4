# 📚 API Reference — Overview (GoMT4)

Short, navigable entry point to the API reference. Use this page to jump to the right place and understand naming rules & conventions.

---

## 🗺️ What’s inside

* **[Messages](./Messages.md)** — request/reply payloads used by the client (`MT4Account`), with field notes.
* **[Enums](./Enums.md)** — the enumerations those messages use, with human meanings.
* **[Streaming](./Streaming.md)** — the server streams (`OnSymbolTick`, `OnTrade`, `OnOpenedOrdersProfit`, `OnOpenedOrdersTickets`) and the client-side history helpers.

> Looking for usage? See **Cookbook** recipes next to this section (e.g. Orders/PlaceMarketOrder, MarketInfo/GetQuote, Streaming/StreamQuotes).

---

## 🏷️ Naming & readability

* All types live in the generated Go package `github.com/MetaRPC/GoMT4/package`, imported as `pb`:

  ```go
  import pb "github.com/MetaRPC/GoMT4/package"
  ```
* Go names follow `protoc-gen-go`: proto `account_balance` → Go field `AccountBalance`, getter `GetAccountBalance()`. Getters are nil-safe; prefer them.
* Enum constants are `<EnumType>_<VALUE>`, e.g. `pb.OrderSendOperationType_OC_OP_BUY`, `pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1`.
* Every reply is `{ Data | Error }` (a proto `oneof`): `reply.GetData()` or `reply.GetError()`. The `MT4Account` methods unwrap this and return the `...Data` message plus a Go `error`.

---

## 🧩 Common type legend

* **`*timestamppb.Timestamp`** ⏰ — a time. Fields without a suffix (`DateTime`, `Time`, `OpenTime`, …) are **broker server time**; the `...Utc` sibling (`DateTimeUtc`, `TimeUtc`, …) is the same moment in UTC (unset while the broker's offset is not known yet). Convert with `.AsTime()`.
* **Pointer scalars (`*float64`, `*int32`, `*string`)** 🎛 — optional fields with presence (`nil` → not set). Use `proto.Float64(...)`, `proto.Int32(...)`, `proto.String(...)` from `google.golang.org/protobuf/proto` to fill them.
* **Money & PnL** 💵 — in **account currency** (`AccountSummaryData.AccountCurrency`).
* **Prices & volumes** 💹 — obey `SymbolParamsManyInfo` (`Digits`, `Point`, `VolumeStep`, `VolumeMin`/`VolumeMax`, `TradeStopsLevel`).
* **Tickets** 🎫 — `int32`.

---

## 🔌 API families → where to read

| Area                     | Client methods (`MT4Account`)                                                       | Start here |
| ------------------------ | ----------------------------------------------------------------------------------- | ---------- |
| **Connection**           | `NewMT4Account`, `ConnectByServerName`, `ConnectByHostPort`, `GetId`, `Disconnect`, `Close` | **Messages** → *Connection* |
| **Account**              | `AccountSummary` (+ `AccountBalance`, `AccountEquity`, … helpers)                   | **Messages** → *Account & Orders* |
| **Orders (sync)**        | `OrderSend`, `OrderModify`, `OrderClose`, `OrderDelete`, `OrderCloseBy`, `OpenedOrders`, `OpenedOrdersTickets`, `OrderSelect`, `OrdersTotal` | **Messages** → *Account & Orders* |
| **Orders (history)**     | `OrdersHistory`, `OrdersHistoryStream`                                              | **Messages** / **Streaming** |
| **Market info & quotes** | `Quote`, `QuoteMany`, `Symbols`, `ShowAllSymbols`, `SymbolParams`, `SymbolParamsMany`, `TickValueWithSize`, `QuoteHistory`, `QuoteHistoryStream` | **Messages** → *Quotes & Market Info* |
| **Streaming**            | `OnSymbolTick`, `OnTrade`, `OnOpenedOrdersProfit`, `OnOpenedOrdersTickets`          | **Streaming** |

The generated package has more services than the client wraps (e.g. `ConnectionClient.ConnectionStatus`, `MarketInfoClient.SymbolSelect`, `TradeStatsClient`). You can call them through the exported clients on `MT4Account` (`ConnectionClient`, `MarketInfoClient`, …) or a client created on its connection (e.g. `pb.NewTradeStatsClient(account.GrpcConn)`), with the `apikey` and `id` metadata — see [Ensure Symbol Visible](../Cookbook/Utils_Helpers/EnsureSymbolVisible.md) for an example.

---

## 🔗 Source of truth

The `.proto` files are owned by the MetaRPC server; the Go code in `package/` is generated from them. When in doubt, the generated `package/*.pb.go` files (field names, types and enum values) are authoritative for Go.

---

## 🚦 Stability notes

* Optional fields (pointer types) may be **omitted** by the server when not applicable.
* New enum values can appear in the future — handle **unknown** values defensively on the client side.
* Streaming: always read the error channel; when a stream ends, the client closes both channels.

Happy building! ✨
