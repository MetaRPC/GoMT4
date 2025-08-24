# 📌 PlacePendingOrder (GoMT4)

**Goal:** open a **pending order** (Buy Limit, Sell Limit, Buy Stop, Sell Stop) with proper rounding and expiration.

> Uses real code from this repo:
>
> * Account: `examples/mt4/MT4Account.go` (`OrderSend`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderSendExample`)

---

## ✅ 1) Preconditions

* MT4 terminal is running & connected to broker.
* `config.json` filled with valid login/server/symbol.
* Symbol visible in MT4 *Market Watch*.

---

## 🔍 2) Read symbol parameters

```go
info, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

digits    := int(info.GetDigits())
volStep   := info.GetVolumeStep()
volMin    := info.GetVolumeMin()
volMax    := info.GetVolumeMax()
pointSize := math.Pow10(-digits)
```

---

## ⚙️ 3) Align helpers (same as market order)

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

slippage := int32(5) // still required but ignored for pending

// Optional SL/TP
offset := 20 * pointSize
stop   := roundPrice(price-offset, digits)
take   := roundPrice(price+2*offset, digits)

// Expiration (good for 1 day)
expiry := timestamppb.New(time.Now().Add(24 * time.Hour))

comment     := "Go pending order"
magicNumber := int32(123456)

resp, err := account.OrderSend(
    ctx,
    symbol,
    side,
    volume,
    &price,          // required for pending
    &slippage,
    &stop, &take,    // can be nil
    &comment,
    &magicNumber,
    expiry,          // ⬅️ required for pending expiration
)
```

---

## 📊 5) Result

```go
fmt.Printf("✅ Pending order placed! Ticket=%d Type=%s Price=%.5f Expires=%s\n",
    resp.GetTicket(), resp.GetType().String(), resp.GetPrice(), resp.GetExpiration().AsTime())
```

---

## ⚠️ Common pitfalls

* **No expiration** → broker may reject if you omit expiry for pending.
* **Invalid price** → must be correctly rounded and on the correct side (e.g., Buy Limit < current Ask).
* **Suffix mismatch** → always check actual symbol name in MT4.

---

## 🔄 Variations

* `OC_OP_SELLLIMIT`, `OC_OP_BUYSTOP`, `OC_OP_SELLSTOP` → change `side`.
* `expiry=nil` → pending order is good‑till‑cancelled (if broker allows).
* Place multiple pendings in loop (just vary `price` and `comment`).
