# Return Codes & Errors Reference — GoMT4

> How failures reach you in GoMT4, and the MQL4 trade-error codes behind rejected orders.

> ⚠️ MT4 uses **MQL4 error codes** (0–150 for the trade server, 4000+ for MQL runtime) — **not** the MT5 `ENUM_TRADE_RETCODE` numbers (10004, 10018, …). If you've seen those, they belong to MT5.

---

## How errors surface

Every SDK method returns Go's idiomatic `(value, error)`. **A failed call returns a non-nil `error`** — there is no separate status code to inspect:

```go
r, err := sugar.BuyMarket(ctx, "EURUSD", 0.10, nil, nil)
if err != nil {
    log.Println("order failed:", err) // e.g. "API error (code=TERMINAL_ERROR): Not enough money"
    return
}
```

Two kinds of error reach you:

* **Transport / gRPC errors** — network, TLS, deadline. `DeadlineExceeded` / `Unavailable` are retried automatically by the client's reconnect wrapper.
* **Application (terminal) errors** — the gateway answered with its `Error` payload. The SDK wraps it as a Go error of the form `API error (code=<CODE>): <message>`, where the MQL failure (invalid stops, market closed, …) is in the message/description.

> **No error ≠ nothing to check** only for *reads*. For *trades*, `err == nil` means it executed; a rejection comes back as a non-nil `error`, so always handle it.

---

## The gateway `Error` payload

When a call fails server-side, the reply carries a `pb.Error` (before the SDK wraps it):

| Getter | Meaning |
|--------|---------|
| `GetErrorCode()` | Gateway code string, e.g. `TERMINAL_INSTANCE_NOT_FOUND` |
| `GetErrorMessage()` | Human-readable message |
| `GetMqlErrorIntCode()` | The **MQL4 error number** (see tables below) |
| `GetMqlErrorCode()` | The same as a typed `pb.MqlErrorCode` enum |
| `GetMqlErrorDescription()` | MQL4 error text |

The convenience methods flatten this to a Go `error`. If you need the raw MQL code, call the method on the **low-level `MT4Account`** and read the reply's `GetError()` before it is wrapped.

---

## MQL4 trade-server error codes

The codes you'll actually hit from trading calls. Source: the generated `pb.MqlErrorCode` enum (`pb.MqlErrorCode_ERR_…`).

| Code | Name | Meaning / fix |
|------|------|---------------|
| 0 | `ERR_NO_ERROR` | Success |
| 1 | `ERR_NO_RESULT` | No change (e.g. modify to identical values) |
| 4 | `ERR_SERVER_BUSY` | Trade server busy — back off |
| 6 | `ERR_NO_CONNECTION` | No link to the trade server |
| 128 | `ERR_TRADE_TIMEOUT` | Trade timed out — **may still have executed**; reconcile |
| 129 | `ERR_INVALID_PRICE` | Price off/stale — normalize price, re-quote |
| **130** | `ERR_INVALID_STOPS` | SL/TP too close — respect `TradeStopsLevel` |
| 131 | `ERR_INVALID_TRADE_VOLUME` | Lot invalid — `NormalizeLot` to step/min/max |
| **132** | `ERR_MARKET_CLOSED` | Market closed — wait for the session |
| 133 | `ERR_TRADE_DISABLED` | Trading disabled for the account/symbol |
| **134** | `ERR_NOT_ENOUGH_MONEY` | Insufficient margin — smaller lot |
| 135 | `ERR_PRICE_CHANGED` | Price moved — re-quote and resend |
| 136 | `ERR_OFF_QUOTES` | No prices — symbol not live / illiquid |
| 137 | `ERR_BROKER_BUSY` | Broker busy — retry |
| **138** | `ERR_REQUOTE` | Requote — resend at the new price, or raise slippage |
| 145 | `ERR_TRADE_MODIFY_DENIED` | Modify denied (too close to market / frozen) |
| 146 | `ERR_TRADE_CONTEXT_BUSY` | Trade context busy — serialize trade calls |
| 148 | `ERR_TRADE_TOO_MANY_ORDERS` | Order count limit reached |

> Codes **4000+** (`ERR_NO_MQLERROR`, `ERR_ARRAY_INDEX_OUT_OF_RANGE`, …) are MQL **runtime** errors — rare for API trading. List them all with the inspector (`enum MqlErrorCode`).

---

## Handling errors in practice

```go
r, err := sugar.BuyMarket(ctx, "EURUSD", 0.10, nil, nil)
if err != nil {
    msg := err.Error()
    switch {
    case strings.Contains(msg, "Market closed"):
        log.Println("market closed — try later")
    case strings.Contains(msg, "Not enough money"):
        log.Println("size down")
    default:
        log.Println("trade failed:", err)
    }
    return
}
```

Need the exact MQL number? Drop to the low-level layer:

```go
// pseudo: inspect the raw reply error on MT4Account for the mql int code
// (the convenience helpers wrap it into a string error)
```

### What's retried for you

The client's reconnect wrapper **auto-retries** transient transport errors and `TERMINAL_INSTANCE_NOT_FOUND` (one reconnect). **Trade calls are not retried** — a timed-out `OrderSend` may have filled, so on error you reconcile with `OpenedOrders` before resending.

---

**See also:** [ENUMs Reference](ENUMS_USAGE_REFERENCE.md) · [Protobuf Inspector](PROTOBUF_INSPECTOR_GUIDE.md) · [MT4Account overview](../MT4Account/Overview.md)
