# 🗑️ DeletePending (GoMT4)

**Goal:** safely delete a **pending order** (Buy/Sell Limit/Stop) by ticket, with small sanity checks.

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (`OrderDelete`, `OrderSelect`, `OpenedOrders`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderDeleteExample`)

---

## ✅ 1) Preconditions

* You have the **ticket** of a *pending* order (not a market position).
* `account` is connected.

---

## 🔎 2) Locate and verify the order (strongly recommended)

`OrderDelete` sends the same `OrderCloseDelete` request as `OrderClose`: given the ticket of a **market** position, it would **close** it. Check the type first:

```go
ord, err := account.OrderSelect(ctx, ticket)
if err != nil { return err }

switch ord.GetOrderType() {
case pb.OpenedOrderType_OO_OP_BUYLIMIT, pb.OpenedOrderType_OO_OP_SELLLIMIT,
     pb.OpenedOrderType_OO_OP_BUYSTOP,  pb.OpenedOrderType_OO_OP_SELLSTOP:
    // ok, this is a pending order
default:
    return fmt.Errorf("ticket %d is not a pending order (type=%s)", ticket, ord.GetOrderType())
}
```

---

## 🗑️ 3) Call `OrderDelete`

`OrderDelete` takes only the ticket (it sends no lots, price or slippage):

```go
resp, err := account.OrderDelete(
    ctx,
    ticket,
)
if err != nil {
    return fmt.Errorf("delete pending failed: %w", err)
}
fmt.Printf("✅ Pending deleted: ticket=%d symbol=%s mode=%s\n",
    ticket, ord.GetSymbol(), resp.GetMode()) // mode: OCD_PENDING_ORDER
```

> Internally, the account helper retries transient transport errors (`codes.Unavailable`) and terminal-not-found errors.

---

## ⚠️ Pitfalls

* **Not a pending** → for market positions use `OrderClose`; `OrderDelete` on a market ticket closes the position.
* **Ticket not found** → it may have been executed or deleted already; refresh open orders and history.
* **Freeze level** → some brokers refuse actions near market price (`TradeFreezeLevel`); try again later.
* **Timeouts** → 5s default when `ctx` has no deadline.

---

## 🔄 Variations

* **Bulk delete** all pendings on a symbol:

```go
opened, err := account.OpenedOrders(ctx)
if err != nil { return err }
for _, o := range opened.GetOrderInfos() {
    if o.GetSymbol() != symbol { continue }
    switch o.GetOrderType() {
    case pb.OpenedOrderType_OO_OP_BUYLIMIT, pb.OpenedOrderType_OO_OP_SELLLIMIT,
         pb.OpenedOrderType_OO_OP_BUYSTOP,  pb.OpenedOrderType_OO_OP_SELLSTOP:
        if _, err := account.OrderDelete(ctx, o.GetTicket()); err != nil {
            log.Printf("delete %d failed: %v", o.GetTicket(), err)
        }
    }
}
```

* **Conditional cleanup**: if `ord.GetExpirationTime() != nil` and the expiration has passed → delete. (`ExpirationTime` is broker server time; compare with `ExpirationTimeUtc` against `time.Now()`.)

---

## 📎 See also

* [PlacePendingOrder](PlacePendingOrder.md) — how to place a pending with correct price/expiry.
* [ModifyOrder](ModifyOrder.md) — change pending price/expiration instead of deleting.
* [CloseOrder](CloseOrder.md) — for market positions.
