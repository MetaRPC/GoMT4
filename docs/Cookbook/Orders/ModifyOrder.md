# ✏️ ModifyOrder (GoMT4)

**Goal:** change **SL/TP** for market orders or **price/expiration** for pending orders — with proper rounding and validation.

> Real code references:
>
> * Account methods: `examples/mt4/MT4Account.go` (e.g., `OrderModify`, `SymbolParams`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderModifyExample`)

---

## ✅ 1) Preconditions

* You know the **ticket** of the order you want to modify.
* Symbol is visible and parameters are known (Digits/LotStep/Frozen levels if applicable).
* For **market orders**: you can modify **SL/TP** only.
* For **pending orders**: you can modify **price** and **expiration** (and SL/TP too).

---

## 🔍 2) Get current order & symbol parameters

```go
// 1) Read symbol params (for rounding)
info, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }
digits := int(info.GetDigits())
point  := math.Pow10(-digits)

// 2) (Optional) Read current order data to decide new SL/TP
ord, err := account.OrderByTicket(ctx, ticket)
if err != nil { return err }
side := ord.GetType() // BUY/SELL or pending types
```

---

## 🧮 3) Helpers (reuse)

```go
func roundPrice(p float64, digits int) float64 {
    mul := math.Pow10(digits)
    return math.Round(p*mul) / mul
}
```

---

## 📝 4) Modify SL/TP for **market** order

* Compute SL/TP relative to current market price (Bid for SELL, Ask for BUY).
* Round to symbol `Digits`.

```go
q, err := account.Quote(ctx, symbol)
if err != nil { return err }
var sl, tp *float64

switch ord.GetType() {
case pb.OrderType_OP_BUY:
    entry := q.GetAsk()
    s := roundPrice(entry - 20*point, digits)
    t := roundPrice(entry + 40*point, digits)
    sl, tp = &s, &t
case pb.OrderType_OP_SELL:
    entry := q.GetBid()
    s := roundPrice(entry + 20*point, digits)
    t := roundPrice(entry - 40*point, digits)
    sl, tp = &s, &t
}

resp, err := account.OrderModify(
    ctx,
    ticket,
    nil,      // price stays the same for market order
    sl, tp,   // new SL/TP
    nil,      // expiration (not used for market)
)
if err != nil {
    return fmt.Errorf("modify failed: %w", err)
}
fmt.Printf("✅ Modified SL/TP for ticket %d at %s\n", ticket, time.Now().Format(time.RFC3339))
```

---

## ⏰ 5) Modify **pending**: price and expiration

* For pendings, you can change the **entry price** and **expiration**.
* Ensure price is on the correct side of the market (e.g., Buy Limit < Ask, Buy Stop > Ask).

```go
// New pending price
desired := 1.09500
price   := roundPrice(desired, digits)

// Move expiry +12 hours
expiry := timestamppb.New(time.Now().Add(12 * time.Hour))

resp, err := account.OrderModify(
    ctx,
    ticket,
    &price, // ⬅️ new pending price
    nil, nil,
    expiry, // ⬅️ new expiration
)
if err != nil {
    return fmt.Errorf("modify pending failed: %w", err)
}
fmt.Printf("✅ Pending modified! Ticket=%d NewPrice=%.5f NewExpiry=%s\n",
    ticket, price, expiry.AsTime().Format("2006-01-02 15:04:05"))
```

---

## ⚠️ Pitfalls & checks

* **Rounding:** always round to `Digits` (SL/TP/price).
* **Side rules:**

  * Buy Limit < current Ask, Sell Limit > current Bid.
  * Buy Stop  > current Ask, Sell Stop  < current Bid.
* **Freeze level:** some brokers restrict how close to price you can modify. If you hit freeze, increase distance.
* **Context timeouts:** wrap modify calls with short timeout (3–8s). Retry only **transport** errors.

---

## 🔗 See also

* `CloseOrder.md` — how to close market orders safely.
* `PlacePendingOrder.md` — how to place pending with correct price/expiry.
* `RoundVolumePrice.md` — helpers for volume/price alignment.
