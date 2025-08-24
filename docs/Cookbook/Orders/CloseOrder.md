# ❌ CloseOrder (GoMT4)

**Goal:** correctly close an open market order (BUY/SELL) with retry logic and proper rounding.

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrderClose`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderCloseExample`)

---

## ✅ 1) Preconditions

* You know the **ticket** of the order you want to close.
* The order is **open** (check via `ShowOpenedOrders` or similar).
* MT4 terminal is running & connected.

---

## 🔍 2) Get order & quote

```go
ord, err := account.OrderByTicket(ctx, ticket)
if err != nil { return err }

q, err := account.Quote(ctx, ord.GetSymbol())
if err != nil { return err }
```

---

## ⚙️ 3) Prepare close params

* Use **Bid** to close a BUY.
* Use **Ask** to close a SELL.
* Round price to `Digits`.

```go
info, _ := account.SymbolParams(ctx, ord.GetSymbol())
digits := int(info.GetDigits())

var closePrice float64
if ord.GetType() == pb.OrderType_OP_BUY {
    closePrice = roundPrice(q.GetBid(), digits)
} else if ord.GetType() == pb.OrderType_OP_SELL {
    closePrice = roundPrice(q.GetAsk(), digits)
}

slippage := int32(5)
```

---

## 📝 4) Call `OrderClose`

```go
resp, err := account.OrderClose(
    ctx,
    ticket,
    ord.GetVolume(),
    &closePrice,
    &slippage,
)
if err != nil {
    return fmt.Errorf("close failed: %w", err)
}
fmt.Printf("✅ Closed ticket=%d at price=%.5f time=%s\n",
    ticket, resp.GetPrice(), resp.GetCloseTime().AsTime())
```

---

## ⚠️ Pitfalls

* **Invalid price** → always round to `Digits`.
* **Wrong side** → BUY closes at Bid, SELL closes at Ask.
* **Partial close** → adjust `volume` parameter if closing partially.
* **Timeout** → wrap in `context.WithTimeout(ctx, 5*time.Second)`.
* **Retries** → transport errors auto‑retried inside helper.

---

## 🔄 Variations

* **Partial close**: send smaller volume than `ord.GetVolume()`.
* **Loop over all open orders**: call `ShowOpenedOrders` → close each.
* **Close by order**: see `CloseByOrders.md` recipe.

---

## 🔗 See also

* `ModifyOrder.md` — adjust SL/TP instead of closing.
* `DeletePending.md` — remove a pending order.
* `HistoryOrders.md` — verify closed orders in history.
