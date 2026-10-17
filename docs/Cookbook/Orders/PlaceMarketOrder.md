# PlaceMarketOrder (GoMT4)

**Goal:** open a market order safely (BUY/SELL) with proper rounding for volume and prices.

> This recipe references real code in this repo:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrderSend`, `SymbolParams`, `Quote`)
> * Service example: `examples/mt4/MT4_service.go` (`ShowOrderSendExample`)
> * Config: `examples/config/config.json`

---

## 1) Preconditions

* `account` is connected (`ConnectByServerName` / `ConnectByHostPort`).
* `examples/config/config.json` contains valid credentials and `DefaultSymbol`.
* The symbol exists on the broker (exact name; see [Ensure Symbol Visible](../Utils_Helpers/EnsureSymbolVisible.md)).
* The account is not logged in with the investor (read‑only) password.

---

## 2) Read symbol parameters (Digits, Volume limits)

```go
info, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

digits    := int(info.GetDigits())
volMin    := info.GetVolumeMin()
volMax    := info.GetVolumeMax()
volStep   := info.GetVolumeStep()
pointSize := info.GetPoint() // usually 10^-Digits
```

> Why: brokers require volume and price aligned to the symbol settings.

---

## 3) Helpers: align volume & round prices

Put these small helpers next to your order code (or reuse existing ones).

```go
func alignVolume(v, step, min, max float64) float64 {
    if step <= 0 { return v }
    v = math.Floor(v/step+0.5) * step
    if v < min { v = min }
    if v > max { v = max }
    return v
}

func roundPrice(p float64, digits int) float64 {
    mul := math.Pow10(digits)
    return math.Round(p*mul) / mul
}
```

---

## 4) Build inputs (market order)

* For a **market** order, pass `price=nil` and set a small `slippage` (in *points*).
* Optionally pre‑compute SL/TP relative to current price.

```go
side := pb.OrderSendOperationType_OC_OP_BUY // or OC_OP_SELL
rawVolume := 0.10
volume    := alignVolume(rawVolume, volStep, volMin, volMax)

var price *float64 = nil
slippage := int32(5) // 5 points

// Optional SL/TP (rounded to Digits)
var sl, tp *float64
if wantSLTP {
    bidAsk, err := account.Quote(ctx, symbol)
    if err != nil { return err }
    entry := bidAsk.GetAsk() // for BUY; for SELL use Bid
    stop  := roundPrice(entry-20*pointSize, digits)
    take  := roundPrice(entry+40*pointSize, digits)
    sl, tp = &stop, &take
}

comment     := "Go order test"
magicNumber := int32(123456)
```

---

## 5) Send the order (with retries handled inside)

`OrderSend` already wraps gRPC with reconnect/retry (see `ExecuteWithReconnect`) and a 5s default timeout. You only supply clean inputs.

```go
resp, err := account.OrderSend(
    ctx,
    symbol,
    side,
    volume,
    price,             // nil for market
    &slippage,         // points
    sl, tp,            // can be nil
    &comment,
    &magicNumber,
    nil /* expiration (for pending only) */,
)
if err != nil {
    return fmt.Errorf("OrderSend failed: %w", err)
}
fmt.Printf("✅ Order opened! Ticket=%d Price=%.5f Time=%s\n",
    resp.GetTicket(), resp.GetPrice(), resp.GetOpenTime().AsTime().Format("2006-01-02 15:04:05"))
```

> See a minimal working example in `examples/mt4/MT4_service.go` → `ShowOrderSendExample` (it uses small unexported pointer helpers `ptrInt32`, `ptrString`; in your own code use `proto.Int32`, `proto.String` from `google.golang.org/protobuf/proto` or take addresses of variables as above).

---

## 6) Common pitfalls

* **Invalid volume/price** → make sure to align volume to `VolumeStep` and round prices to `Digits`.
* **`symbol not found` / empty quotes** → wrong name or broker suffix (e.g., `EURUSD.m`). List names with `account.Symbols(ctx)`.
* **`context deadline exceeded`** → the terminal is not ready or the broker is slow; pass a longer deadline in `ctx` (see *Reliability* chapter).
* **Trade refused** → check `MqlErrorCode` in the error (e.g. `ERR_INVALID_STOPS`, `ERR_NOT_ENOUGH_MONEY`, `ERR_MARKET_CLOSED`).

---

## 7) Variations

* **SELL**: use `pb.OrderSendOperationType_OC_OP_SELL` and entry price = `Bid` when computing SL/TP.
* **Immediate TP/SL omitted**: pass `nil` and modify later via `OrderModify`.
* **Pending order**: use the same call with a pending type, `price!=nil` and optionally `expiration` (see [PlacePendingOrder](PlacePendingOrder.md)).
