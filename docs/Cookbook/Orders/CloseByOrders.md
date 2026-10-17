# 🔁 CloseByOrders (GoMT4)

**Goal:** close two **opposite** positions of the **same symbol** against each other (MT4 "Close By").

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrderCloseBy`, `OpenedOrders`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderCloseByExample`)

---

## ✅ 1) Preconditions

* Both orders are **open market positions** and belong to the **same symbol** (e.g., `EURUSD`).
* Orders are of **opposite types**: one **BUY**, one **SELL**.
* The broker allows Close By (MT4 accounts are hedging accounts, so usually yes).

**Notes about volume:**

* If volumes **equal** → both positions are fully closed.
* If volumes **differ** → the overlapping volume is closed; the **remainder** of the larger position stays open (MT4 gives it a new ticket).

---

## 🔍 2) Find a pair to close-by

```go
// Find an opposite order for the same symbol.
func pickCloseByPair(ctx context.Context, a *mt4.MT4Account, sym string) (buy, sell *pb.OpenedOrderInfo, err error) {
    opened, err := a.OpenedOrders(ctx)
    if err != nil { return nil, nil, err }
    for _, o := range opened.GetOrderInfos() {
        if o.GetSymbol() != sym { continue }
        switch o.GetOrderType() {
        case pb.OpenedOrderType_OO_OP_BUY:
            if buy == nil { buy = o }
        case pb.OpenedOrderType_OO_OP_SELL:
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

> You can refine selection by `MagicNumber`, `Comment`, or `Lots`.

---

## ⚙️ 3) Prepare params

```go
buy, sell, err := pickCloseByPair(ctx, account, symbol)
if err != nil { return err }

// Choose which ticket goes first (the order to close vs. the opposite one)
src := buy.GetTicket()  // ticketToClose
dst := sell.GetTicket() // oppositeTicket
```

---

## 📝 4) Call `OrderCloseBy`

```go
resp, err := account.OrderCloseBy(
    ctx,
    src,
    dst,
)
if err != nil {
    return fmt.Errorf("close-by failed: %w", err)
}
fmt.Printf("✅ CloseBy done: src=%d dst=%d price=%.5f profit=%.2f time=%s\n",
    src, dst, resp.GetClosePrice(), resp.GetProfit(), resp.GetCloseTime().AsTime().Format(time.RFC3339))
```

> If lots differ, the overlapping volume is matched and closed; the remainder of the larger position stays open. Handle that in your post‑logic if you want to loop and fully flatten.

---

## 🔄 5) Full flatten loop (if volumes differ)

```go
func closeByAll(ctx context.Context, a *mt4.MT4Account, sym string) error {
    for {
        buy, sell, err := pickCloseByPair(ctx, a, sym)
        if err != nil { return nil } // nothing more to pair
        if _, err := a.OrderCloseBy(ctx, buy.GetTicket(), sell.GetTicket()); err != nil {
            return err
        }
        // loop continues until no opposite pair remains
    }
}
```

---

## ⚠️ Pitfalls

* **Different symbols/suffixes** → `EURUSD` vs `EURUSD.m` are *not* the same.
* **Same direction** → Buy+Buy or Sell+Sell cannot be closed by.
* **Pending orders** → only market positions can be closed by.
* **Timeouts** → 5s default when `ctx` has no deadline; transport errors are retried inside the client.

---

## 📎 See also

* [CloseOrder](CloseOrder.md) — standard market close.
* [ModifyOrder](ModifyOrder.md) — adjust SL/TP.
* [HistoryOrders](HistoryOrders.md) — verify results in account history.
