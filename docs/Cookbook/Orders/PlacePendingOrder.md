# 📌 PlacePendingOrder (GoMT4)

**Goal:** open a **pending order** (Buy Limit, Sell Limit, Buy Stop, Sell Stop) with proper rounding and expiration.

> Uses real code from this repo:
>
> * Account: `examples/mt4/MT4Account.go` (`OrderSend`, `SymbolParams`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderSendExample`, a market BUY)

---

## ✅ 1) Preconditions

* `account` is connected.
* `config.json` filled with valid login/server/symbol.
* The symbol exists on the broker (exact name).

---

## 🔍 2) Read symbol parameters

```go
info, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

digits    := int(info.GetDigits())
volStep   := info.GetVolumeStep()
volMin    := info.GetVolumeMin()
volMax    := info.GetVolumeMax()
pointSize := info.GetPoint()
```

---

## ⚙️ 3) Align helpers (same as market order)

`alignVolume` and `roundPrice` are defined in [PlaceMarketOrder](PlaceMarketOrder.md#3-helpers-align-volume-round-prices).

```go
volume := alignVolume(rawVolume, volStep, volMin, volMax)
price  := roundPrice(desiredPrice, digits)
```

---

## 📝 4) Build inputs (example: Buy Limit)

```go
side    := pb.OrderSendOperationType_OC_OP_BUYLIMIT
volume  := alignVolume(0.10, volStep, volMin, volMax)
price   := roundPrice(1.09500, digits) // target entry price

// Optional SL/TP
offset := 20 * pointSize
stop   := roundPrice(price-offset, digits)
take   := roundPrice(price+2*offset, digits)

// Expiration (good for 1 day). The server reads it as broker server time.
expiry := timestamppb.New(time.Now().Add(24 * time.Hour))

comment     := "Go pending order"
magicNumber := int32(123456)

resp, err := account.OrderSend(
    ctx,
    symbol,
    side,
    volume,
    &price,          // required for pending
    nil,             // slippage: not needed for pending orders
    &stop, &take,    // can be nil
    &comment,
    &magicNumber,
    expiry /* optional expiration */,
)
if err != nil {
    return fmt.Errorf("pending order failed: %w", err)
}
```

---

## 📊 5) Result

`*pb.OrderSendData` has `Ticket`, `Volume`, `Price`, `OpenTime` / `OpenTimeUtc`:

```go
fmt.Printf("✅ Pending order placed! Ticket=%d Type=%s Price=%.5f Expires=%s\n",
    resp.GetTicket(), side.String(), resp.GetPrice(), expiry.AsTime().Format("2006-01-02 15:04"))
```

---

## ⚠️ Common pitfalls

* **Expiration refused** → some brokers do not allow expiration on pending orders; pass `nil`.
* **Invalid price** → must be correctly rounded and on the correct side (e.g., Buy Limit < current Ask) and at least `TradeStopsLevel` points away.
* **Suffix mismatch** → always check the actual symbol name (`account.Symbols(ctx)`).

---

## 🔄 Variations

* `OC_OP_SELLLIMIT`, `OC_OP_BUYSTOP`, `OC_OP_SELLSTOP` → change `side`.
* `expiry=nil` → pending order stays until filled or deleted.
* Place multiple pendings in loop (just vary `price` and `comment`).
