# 🔁 CloseByOrders (GoMT4)

**Goal:** close two **opposite** positions of the **same symbol** against each other (MT4 "Close By").

> Real code refs (by convention in this repo):
>
> * Account methods: `examples/mt4/MT4Account.go` (e.g., `OrderCloseBy`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderCloseByExample`)

---

## ✅ 1) Preconditions

* Both orders are **open** and belong to the **same symbol** (e.g., `EURUSD`).
* Orders are of **opposite types**: one **BUY**, one **SELL**.
* Broker **allows hedging** (MT4 usually does; on netting accounts this may be unavailable).

**Notes about volume:**

* If volumes **equal** → both positions are fully closed.
* If volumes **differ** → the **smaller** volume is closed out; the **remainder** stays open with reduced volume.

---

## 🔍 2) Find a pair to close-by

```go
// Find an opposite order for the same symbol.
func pickCloseByPair(ctx context.Context, a *MT4Account, sym string) (buy, sell *pb.Order, err error) {
    opened, err := a.ShowOpenedOrders(ctx)
    if err != nil { return nil, nil, err }
    for _, o := range opened.GetOrders() {
        if o.GetSymbol() != sym { continue }
        switch o.GetType() {
        case pb.OrderType_OP_BUY:
            if buy == nil { buy = o }
        case pb.OrderType_OP_SELL:
            if sell == nil { sell = o }
        }
        if buy != nil && sell != nil { break }
    }
    if buy == nil || sell == nil {
        return nil, nil, fmt.Errorf("no opposite pair for %s", sym)
    }
    return buy, sell, nil
}
```

> You can refine selection by `MagicNumber`, `Comment`, or minimal slippage tolerance.

---

## ⚙️ 3) Prepare params

```go
buy, sell, err := pickCloseByPair(ctx, account, symbol)
if err != nil { return err }

// Optional: choose which ticket goes first (source vs. target)
src := buy.GetTicket()  // any order can be src
dst := sell.GetTicket() // the other one is dst

slippage := int32(5) // points (broker-dependent)
```

---

## 📝 4) Call `OrderCloseBy`

```go
resp, err := account.OrderCloseBy(
    ctx,
    src,
    dst,
    &slippage,
)
if err != nil {
    return fmt.Errorf("close-by failed: %w", err)
}
fmt.Printf("✅ CloseBy done: src=%d dst=%d result=%s time=%s\n",
    src, dst, resp.GetResult().String(), time.Now().Format(time.RFC3339))
```

> If lots differ, the smaller lot is matched and closed; the remainder of the larger position stays open (volume reduced). Handle that in your post‑logic if you want to loop and fully flatten.

---

## 🔄 5) Full flatten loop (if volumes differ)

```go
func closeByAll(ctx context.Context, a *MT4Account, sym string, slip int32) error {
    for {
        buy, sell, err := pickCloseByPair(ctx, a, sym)
        if err != nil { return nil } // nothing more to pair
        _, err = a.OrderCloseBy(ctx, buy.GetTicket(), sell.GetTicket(), &slip)
        if err != nil { return err }
        // loop continues until no opposite pair remains
    }
}
```

---

## ⚠️ Pitfalls

* **Different symbols/suffixes** → `EURUSD` vs `EURUSD.m` are *not* the same.
* **Same direction** → Buy+Buy or Sell+Sell cannot be closed by.
* **Netting account** → feature may be unavailable; use `OrderClose` instead.
* **Freeze/Stops level** → some brokers restrict actions near price; CloseBy usually ignores price, but slippage rules may still apply.
* **Timeouts** → wrap in `context.WithTimeout` (3–8s) and retry only transport errors.

---

## 📎 See also

* `CloseOrder.md` — standard market close.
* `ModifyOrder.md` — adjust SL/TP.
* `HistoryOrders.md` — verify results in account history.
