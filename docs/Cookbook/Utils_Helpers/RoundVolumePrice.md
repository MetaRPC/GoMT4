# 🔢 RoundVolumePrice (GoMT4)

**Goal:** correctly round **volume (lots)** and **price** values according to symbol parameters (Digits, VolumeStep, Min/Max volume). This avoids broker rejections when sending/modifying orders.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`SymbolParams` provides `Digits`, `Point`, `VolumeStep`, `VolumeMin`, `VolumeMax`)
> * Demos: `examples/mt4/MT4_service.go` (the order examples use fixed demo values and do **not** round)

---

## ✅ Why we need this

* Brokers only accept volumes that are multiples of `VolumeStep` and within `[VolumeMin, VolumeMax]`.
* Prices must be rounded to the symbol’s `Digits` (e.g., 5‑digit EURUSD → 1.23456).
* If you don’t align values, `OrderSend` or `OrderModify` returns an error (invalid volume/price).

---

## 🔎 1) Read params first

```go
p, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

fmt.Printf("%s Digits=%d VolumeStep=%.2f VolumeMin=%.2f VolumeMax=%.2f\n",
    symbol, p.GetDigits(), p.GetVolumeStep(), p.GetVolumeMin(), p.GetVolumeMax())
```

---

## 🧮 2) Helpers (convenience functions)

⚠️ Note: these helpers are not part of the repo or the generated `pb` package. Add them to your own code (e.g. `utils.go`); they simply apply the parameters you fetch with `SymbolParams`.

```go
func AlignVolume(v, step, min, max float64) float64 {
    if step <= 0 { return v }
    v = math.Floor(v/step+0.5) * step
    if v < min { v = min }
    if v > max { v = max }
    return v
}

func RoundPrice(p float64, digits int) float64 {
    mul := math.Pow10(digits)
    return math.Round(p*mul) / mul
}
```

Usage:

```go
vol := AlignVolume(0.13, p.GetVolumeStep(), p.GetVolumeMin(), p.GetVolumeMax())
price := RoundPrice(1.092345, int(p.GetDigits()))
```

---

## 📝 3) Use in OrderSend

```go
side := pb.OrderSendOperationType_OC_OP_BUY

q, err := account.Quote(ctx, symbol)
if err != nil { return err }

vol := AlignVolume(0.13, p.GetVolumeStep(), p.GetVolumeMin(), p.GetVolumeMax())
price := RoundPrice(q.GetAsk(), int(p.GetDigits()))
slip := int32(5)
comment := "rounded order"
magic := int32(123456)

resp, err := account.OrderSend(ctx, symbol, side, vol, &price, &slip, nil, nil, &comment, &magic, nil)
if err != nil {
    return fmt.Errorf("OrderSend failed: %w", err)
}
fmt.Println("ticket:", resp.GetTicket())
```

---

## ⚠️ Pitfalls

* **Skipping rounding** → broker rejects with “invalid price/volume”.
* **Wrong Digits** → always take from `SymbolParams`, not from hardcoded assumptions.
* **Different brokers** → `VolumeStep` can differ (0.01 vs 0.1) and `VolumeMin`/`VolumeMax` vary.

---

## 📎 See also

* [SymbolParams](../Market_Info/SymbolParams.md) — explains where Digits/VolumeStep come from.
* [PlaceMarketOrder](../Orders/PlaceMarketOrder.md), [PlacePendingOrder](../Orders/PlacePendingOrder.md) — show order placement using these helpers.
* [ModifyOrder](../Orders/ModifyOrder.md) — reuses `roundPrice` for SL/TP adjustments.
