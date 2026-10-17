# ❌ CloseOrder (GoMT4)

**Goal:** correctly close an open market order (BUY/SELL) — fully or partially — with retry logic and proper rounding.

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrderClose`, `OrderSelect`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderCloseExample`)

---

## ✅ 1) Preconditions

* You know the **ticket** of the order you want to close.
* The order is **open** (check via `account.OpenedOrders` / `OrderSelect`).
* `account` is connected.

---

## 🔍 2) Get order & quote

```go
ord, err := account.OrderSelect(ctx, ticket)
if err != nil { return err }

q, err := account.Quote(ctx, ord.GetSymbol())
if err != nil { return err }
```

---

## ⚙️ 3) Prepare close params

* Use **Bid** to close a BUY.
* Use **Ask** to close a SELL.
* Round price to `Digits`.
* Passing `nil` as price closes at the current market price.

```go
info, err := account.SymbolParams(ctx, ord.GetSymbol())
if err != nil { return err }
digits := int(info.GetDigits())

var closePrice float64
if ord.GetOrderType() == pb.OpenedOrderType_OO_OP_BUY {
    closePrice = roundPrice(q.GetBid(), digits)
} else if ord.GetOrderType() == pb.OpenedOrderType_OO_OP_SELL {
    closePrice = roundPrice(q.GetAsk(), digits)
}

lots := ord.GetLots() // whole position
slippage := int32(5)
```

---

## 📝 4) Call `OrderClose`

```go
resp, err := account.OrderClose(
    ctx,
    ticket,
    &lots,
    &closePrice,
    &slippage,
)
if err != nil {
    return fmt.Errorf("close failed: %w", err)
}
fmt.Printf("✅ Closed ticket=%d mode=%s comment=%q\n",
    ticket, resp.GetMode(), resp.GetHistoryOrderComment())
```

`OrderClose` returns `*pb.OrderCloseDeleteData` (`Mode`, `HistoryOrderComment`). For the close price and profit, read the order from history (`OrdersHistory`) or from the `OnTrade` stream (`NewHistoryOrders`).

---

## ⚠️ Pitfalls

* **Invalid price** → always round to `Digits`, or pass `nil` for market price.
* **Wrong side** → BUY closes at Bid, SELL closes at Ask.
* **Partial close** → pass a smaller `lots`; the rest stays open (MT4 gives it a new ticket).
* **Timeout** → 5s default when `ctx` has no deadline.
* **Retries** → transport and terminal-not-found errors are retried inside the client.
* **Pending orders** → the same call deletes a pending order (`Mode = OCD_PENDING_ORDER`).

---

## 🔄 Variations

* **Partial close**: send smaller `lots` than `ord.GetLots()`.
* **Loop over all open orders**: `account.OpenedOrders(ctx)` → close each market order.
* **Close by order**: see [CloseByOrders](CloseByOrders.md).

---

## 🔗 See also

* [ModifyOrder](ModifyOrder.md) — adjust SL/TP instead of closing.
* [DeletePending](DeletePending.md) — remove a pending order.
* [HistoryOrders](HistoryOrders.md) — verify closed orders in history.
