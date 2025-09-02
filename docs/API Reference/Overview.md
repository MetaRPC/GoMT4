# 📚 API Reference — Overview (GoMT4)

Short, navigable entry point to the API reference. Use this page to jump to the right place and understand naming rules & conventions.

---

## 🗺️ What’s inside

* **[Messages](./Messages.md)** — payload structures (requests, responses, snapshots) with field notes.
* **[Enums](./Enums.md)** — all enumerations with human meanings.
* **[Streaming](./Streaming.md)** — long‑lived gRPC streams and their chunk types.

> Looking for usage? See **Cookbook** recipes next to this section (e.g. Orders/PlaceMarketOrder, MarketInfo/GetQuote, Streaming/StreamQuotes).

---

## 🏷️ Naming & readability

* Original proto names are **prefixed with `Mt4`** (e.g., `Mt4AccountSummary`).
* In headings we show **both**: full name and a short alias in parentheses — e.g. *Mt4AccountSummary (AccountSummary)*.
* Inside tables and notes we use **short names** for easier reading.

**Why `Mt4`?** It scopes types (MT4 vs MT5) and avoids name collisions across modules/languages.

---

## 🧩 Common type legend

* **`Timestamp`** ⏰ — UTC time. Log in RFC3339.
* **`DoubleValue` / `StringValue` / `Int32Value`** 🎛 — optional fields with presence (omit → not set).
* **Money & PnL** 💵 — in **account currency** (`AccountSummary.currency`).
* **Prices & volumes** 💹 — obey `SymbolParams` (`digits`, `point`, `lot_step`, `lot_min/max`).

---

## 🔌 API families → where to read

| Area                     | Start here                                                                                                                        |
| ------------------------ | --------------------------------------------------------------------------------------------------------------------------------- |
| **Connection & Health**  | `Connect/Disconnect`, `ConnectionState`, `Ping`, `ServerInfo`, `Time`, **HealthCheck** → see **Messages** → *Connection & Health* |
| **Orders (sync)**        | `OrderSendRequest/Result`, `OrderModify/Close/Delete/CloseBy`, `OrderActionResult` → **Messages** → *Account & Orders*            |
| **Orders (history)**     | `OrderHistory`, paged/streaming history → **Streaming** → *Orders History Streaming*                                              |
| **Market info & quotes** | `Quote`, `SymbolParams`, tick values, quote history → **Messages** → *Quotes & Market Info*                                       |
| **Streaming**            | Quotes, trade updates, opened tickets, PnL snapshots, charts → **Streaming**                                                      |

---

## 🔗 Source of truth

This reference is generated from your `.proto` files (mrpc‑proto). When proto changes, **Messages/Enums/Streaming** are updated to match field order and enum values exactly.

---

## 🚦 Stability notes

* Fields marked optional (wrapper types) may be **omitted** by the server when not applicable.
* New enum values can appear in the future — handle **unknown** values defensively on the client side.
* Streaming: always process `is_last = true` and surface transport errors to your retry logic.

Happy building! ✨
