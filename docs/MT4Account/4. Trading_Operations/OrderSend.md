# OrderSend (RECOMMENDED)

> Open a market or pending order — the low-level entry point that creates positions and pendings.

**API Information:**

* **Low-level API:** `MT4Account.OrderSend(...)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.TradingHelper`
* **RPC:** `OrderSend(OrderSendRequest) → OrderSendReply`
* **Payload:** `OrderSendData`

## 💬 Just the essentials

* **What it is.** MQL `OrderSend()` — a market buy/sell, or any of the four pendings, with optional SL/TP, comment and magic.
* **Why you need it.** Full control over every order field at the lowest level. It's how you enter the market.
* **When to use.** When you need raw control. For pip-based SL/TP, lot normalization and risk sizing, prefer [`MT4Sugar`](../../MT4Sugar/5.%20Order_Entry.md) (`BuyMarket`, `BuyLimit`, …) which wraps this.

---

## 🎯 Purpose

* Open market orders (BUY / SELL)
* Place pending orders (BUY/SELL · LIMIT/STOP)
* Attach SL / TP, comment and magic number
* Set slippage and a pending-order expiration

---

## Method Signature

```go
func (a *MT4Account) OrderSend(
    ctx context.Context,
    symbol string,
    operationType pb.OrderSendOperationType,
    volume float64,
    price *float64,       // required for pendings; nil for market
    slippage *int32,      // points
    stoploss *float64,    // price; nil/0 = none
    takeprofit *float64,  // price; nil/0 = none
    comment *string,
    magicNumber *int32,
    expiration *timestamppb.Timestamp, // pending expiry
) (*pb.OrderSendData, error)
```

**Request message:**

```protobuf
message OrderSendRequest {
  string symbol = 1;                           // required, non-empty
  OrderSendOperationType operation_type = 2;   // required
  double volume = 3;                           // required (lots)
  optional double price = 4;                   // required for *LIMIT / *STOP; omit for market
  optional int32  slippage = 5;                // points; default 0
  optional double stoploss = 6;                // price; default 0 (none)
  optional double takeprofit = 7;              // price; default 0 (none)
  optional string comment = 8;
  optional int32  magic_number = 9;
  optional Timestamp expiration = 10;          // (or expiration_utc — send one, not both)
  optional Timestamp expiration_utc = 11;
}
```

**Payload message:**

```protobuf
message OrderSendData {
  int32 ticket = 1;
  double volume = 2;
  double price = 3;
  Timestamp open_time = 4;
  Timestamp open_time_utc = 5;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbol` | `string` | yes | Instrument, e.g. `"EURUSD"`. |
| `operationType` | `pb.OrderSendOperationType` | yes | Buy/sell, market or pending (see enum below). |
| `volume` | `float64` | yes | Lots — normalize to the symbol's `VolumeStep` first. |
| `price` | `*float64` | pendings | Entry/trigger price; **required for pendings**, `nil` for market. |
| `slippage` | `*int32` | no | Max deviation in **points**; `nil` = broker default. |
| `stoploss` | `*float64` | no | SL **price** (not pips); `nil`/`0` = none. |
| `takeprofit` | `*float64` | no | TP **price** (not pips); `nil`/`0` = none. |
| `comment` | `*string` | no | Free-text order comment. |
| `magicNumber` | `*int32` | no | EA / strategy id to tag the order. |
| `expiration` | `*timestamppb.Timestamp` | no | Pending expiry; `nil` = GTC. |

---

## ⬆️ Output — `*pb.OrderSendData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `ticket` | int32 | `GetTicket() int32` | Ticket of the new order |
| `volume` | double | `GetVolume() float64` | Filled / placed volume |
| `price` | double | `GetPrice() float64` | Actual fill price (market) or pending price |
| `open_time` / `open_time_utc` | Timestamp | `GetOpenTime()` / `GetOpenTimeUtc() *Timestamp` | Open/placement time (server / UTC) |

> On failure the call returns a non-nil `error` carrying the broker's code and message as text (`"API error (code=…): …"`). There is no numeric return code to read off the result — map the message via the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md). Trade calls are **not** auto-retried.

---

## 🧱 Related enums (from proto)

### `OrderSendOperationType` — input `operationType`

| Constant | Value | Meaning |
|----------|-------|---------|
| `OC_OP_BUY` | 0 | Market buy |
| `OC_OP_SELL` | 1 | Market sell |
| `OC_OP_BUYLIMIT` | 2 | Buy limit (pending, below price) |
| `OC_OP_SELLLIMIT` | 3 | Sell limit (pending, above price) |
| `OC_OP_BUYSTOP` | 4 | Buy stop (pending, above price) |
| `OC_OP_SELLSTOP` | 5 | Sell stop (pending, below price) |

> In Go, pass the typed constant directly: `pb.OrderSendOperationType_OC_OP_BUY`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Market vs pending.** `OC_OP_BUY`/`OC_OP_SELL` fill immediately (`price` may be `nil`); the `*LIMIT`/`*STOP` variants are pending and **require** `price`.
* **Absolute SL/TP prices.** `stoploss`/`takeprofit` are prices, not pip distances. Respect the symbol's `GetTradeStopsLevel()` (from [`SymbolParams`](../3.%20Symbols_Quotes/SymbolParams.md)) or you'll get `ERR_INVALID_STOPS` (130).
* **Normalize first.** Round `volume` to `VolumeStep` and prices to `Digits`; [`MT4Sugar`](../../MT4Sugar/3.%20Symbol_Math.md) does this automatically.
* **Tag with `magicNumber`.** It's how you later pick out "my" orders in [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md).
* **Returns a ticket** on success — errors are returned values, not a numeric result field.

---

## 🔗 Usage Examples

> ⚠️ These place **real orders on the connected (demo) account**. Use `0.01` lots and clean up.

### 1) Open a market BUY

```go
import pb "github.com/MetaRPC/GoMT4/package"

func buy(ctx context.Context, acc *mt4.MT4Account, symbol string, lots float64) (int32, error) {
    r, err := acc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUY, lots,
        nil, nil, nil, nil, nil, nil, nil)
    if err != nil {
        return 0, err
    }
    fmt.Printf("opened #%d %.2f lot @ %.5f\n", r.GetTicket(), r.GetVolume(), r.GetPrice())
    return r.GetTicket(), nil
}
```

### 2) Market SELL with SL/TP, comment and magic

```go
func sellWithStops(ctx context.Context, acc *mt4.MT4Account, symbol string, lots float64) (int32, error) {
    q, err := acc.Quote(ctx, symbol)
    if err != nil { return 0, err }
    sl := q.GetBid() + 0.0020 // 20 pips above (stop for a short)
    tp := q.GetBid() - 0.0030 // 30 pips below
    slip := int32(20)
    cmt := "demo sell"
    magic := int32(777)
    r, err := acc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_SELL, lots,
        nil, &slip, &sl, &tp, &cmt, &magic, nil)
    if err != nil { return 0, err }
    fmt.Printf("opened #%d SL=%.5f TP=%.5f\n", r.GetTicket(), sl, tp)
    return r.GetTicket(), nil
}
```

### 3) Place a BUY LIMIT below the market (with expiry)

```go
import (
    pb "github.com/MetaRPC/GoMT4/package"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func buyLimit(ctx context.Context, acc *mt4.MT4Account, symbol string, lots float64) (int32, error) {
    q, err := acc.Quote(ctx, symbol)
    if err != nil { return 0, err }
    price := q.GetAsk() - 0.0030 // 30 pips below current ask
    exp := timestamppb.New(time.Now().Add(2 * time.Hour))
    r, err := acc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUYLIMIT, lots,
        &price, nil, nil, nil, nil, nil, exp)
    if err != nil { return 0, err }
    fmt.Printf("placed BUY LIMIT #%d @ %.5f (expires 2h)\n", r.GetTicket(), price)
    return r.GetTicket(), nil
}
```

### 4) Robust send — reconcile on error

Go surfaces the broker's rejection as a text error (code + message); a *timed-out* send may still have filled, so reconcile with `OpenedOrders` before retrying.

```go
func safeBuy(ctx context.Context, acc *mt4.MT4Account, symbol string, lots float64) (int32, error) {
    r, err := acc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUY, lots,
        nil, nil, nil, nil, nil, nil, nil)
    if err == nil {
        return r.GetTicket(), nil
    }

    // Did it actually fill despite the error? Check before re-sending.
    if open, e2 := acc.OpenedOrders(ctx); e2 == nil {
        for _, o := range open.GetOrderInfos() {
            if o.GetSymbol() == symbol && o.GetComment() == "" {
                log.Printf("send errored (%v) but an order on %s exists (#%d) — not retrying",
                    err, symbol, o.GetTicket())
                return o.GetTicket(), nil
            }
        }
    }
    return 0, fmt.Errorf("order_send failed: %w", err) // see Return Codes Reference
}
```

---

## ⚠️ Pitfalls

* **Pending orders need `price`.** Omitting it for a `*LIMIT`/`*STOP` fails.
* **SL/TP are prices, not pips** — and must respect `GetTradeStopsLevel()`; too-close stops are rejected with `ERR_INVALID_STOPS` (130).
* **Un-normalized volume/price.** Values off the broker's step/digits are rejected (`ERR_INVALID_TRADE_VOLUME` 131 / `ERR_INVALID_PRICE` 129).
* **Don't double-send on timeout.** A timed-out call may already have filled — check `OpenedOrders` before retrying.
* **Clean up demos.** Close what you open with [`OrderClose`](./OrderClose.md).

---

## 🧪 Testing Suggestions

* **Happy path (demo).** A `0.01` market BUY returns a positive `GetTicket()`; the ticket then appears in [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md).
* **Round-trip.** Open, then [`OrderClose`](./OrderClose.md) the ticket; it leaves `OpenedOrders`.
* **Pending + expiry.** Place a far BUY LIMIT with a short expiry; confirm it rests, then lapses.
* **Failure path.** Send an oversized volume and assert a non-nil error whose message mentions an invalid-volume code (131).

---

## 📚 See also

* [OrderModify](./OrderModify.md) · [OrderClose](./OrderClose.md) · [OrderCloseBy](./OrderCloseBy.md) · [OpenedOrders](../2.%20Orders_History/OpenedOrders.md)
* **Back to:** [4. Trading Operations — Overview](./Trading_Operations.Overview.md) · [MT4Account — Overview](../Overview.md)
