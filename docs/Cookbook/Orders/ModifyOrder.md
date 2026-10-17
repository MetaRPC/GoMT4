# ✏️ ModifyOrder (GoMT4)

**Goal:** change **SL/TP** for market orders or **price/expiration** for pending orders — with proper rounding and validation.

> Real code references:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrderModify`, `OrderSelect`, `SymbolParams`, `Quote`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderModifyExample`)

---

## ✅ 1) Preconditions

* You know the **ticket** of the order you want to modify.
* Symbol parameters are known (Digits / StopsLevel / FreezeLevel).
* For **market orders**: you can modify **SL/TP** only.
* For **pending orders**: you can modify **price** and **expiration** (and SL/TP too).

---

## 🔍 2) Get current order & symbol parameters

```go
// 1) Read current order data (searches the opened orders by ticket)
ord, err := account.OrderSelect(ctx, ticket)
if err != nil { return err }
symbol := ord.GetSymbol()

// 2) Read symbol params (for rounding)
info, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }
digits := int(info.GetDigits())
point  := info.GetPoint()
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

switch ord.GetOrderType() {
case pb.OpenedOrderType_OO_OP_BUY:
    entry := q.GetAsk()
    s := roundPrice(entry - 20*point, digits)
    t := roundPrice(entry + 40*point, digits)
    sl, tp = &s, &t
case pb.OpenedOrderType_OO_OP_SELL:
    entry := q.GetBid()
    s := roundPrice(entry + 20*point, digits)
    t := roundPrice(entry - 40*point, digits)
    sl, tp = &s, &t
}

modified, err := account.OrderModify(
    ctx,
    ticket,
    nil,      // price stays the same for market order
    sl, tp,   // new SL/TP
    nil /* expiration (not used for market) */,
)
if err != nil {
    return fmt.Errorf("modify failed: %w", err)
}
fmt.Printf("✅ Modified=%v SL/TP for ticket %d at %s\n", modified, ticket, time.Now().Format(time.RFC3339))
```

---

## ⏰ 5) Modify **pending**: price and expiration

* For pendings, you can change the **entry price** and **expiration** (broker server time).
* Ensure price is on the correct side of the market (e.g., Buy Limit < Ask, Buy Stop > Ask).

```go
// New pending price
desired := 1.09500
price   := roundPrice(desired, digits)

// Move expiry +12 hours
expiry := timestamppb.New(time.Now().Add(12 * time.Hour))

modified, err := account.OrderModify(
    ctx,
    ticket,
    &price, // ⬅️ new pending price
    nil, nil,
    expiry /* ⬅️ new expiration */,
)
if err != nil {
    return fmt.Errorf("modify pending failed: %w", err)
}
fmt.Printf("✅ Pending modified=%v! Ticket=%d NewPrice=%.5f NewExpiry=%s\n",
    modified, ticket, price, expiry.AsTime().Format("2006-01-02 15:04:05"))
```

---

## ⚠️ Pitfalls & checks

* **Rounding:** always round to `Digits` (SL/TP/price).
* **Side rules:**

  * Buy Limit < current Ask, Sell Limit > current Bid.
  * Buy Stop  > current Ask, Sell Stop  < current Bid.
* **Freeze level:** some brokers restrict how close to price you can modify (`TradeFreezeLevel`). If you hit freeze, increase distance.
* **Context timeouts:** the call has a 5s default timeout; transport errors are retried inside the client.
* **`OrderSelect` scope:** it searches only currently opened orders (market and pending), not history.

---

## 🔗 See also

* [CloseOrder](CloseOrder.md) — how to close market orders safely.
* [PlacePendingOrder](PlacePendingOrder.md) — how to place pending with correct price/expiry.
* [RoundVolumePrice](../Utils_Helpers/RoundVolumePrice.md) — helpers for volume/price alignment.
