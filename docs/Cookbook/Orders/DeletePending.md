# 🗑️ DeletePending (GoMT4)

**Goal:** safely delete a **pending order** (Buy/Sell Limit/Stop) by ticket, with small sanity checks.

> Real code refs:
>
> * Account methods: `examples/mt4/MT4Account.go` (e.g., `OrderDelete`, `OrderByTicket`, `ShowOpenedOrders`)
> * Example: `examples/mt4/MT4_service.go` (`ShowOrderDeleteExample`)

---

## ✅ 1) Preconditions

* You have the **ticket** of a *pending* order (not a market position).
* MT4 terminal is running & connected.
* Symbol is visible in MT4 (helps avoid “symbol not found” edge cases).

---

## 🔎 2) Locate and verify the order (optional but safer)

```go
ord, err := account.OrderByTicket(ctx, ticket)
if err != nil { return err }

switch ord.GetType() {
case pb.OrderType_OP_BUYLIMIT, pb.OrderType_OP_SELLLIMIT,
     pb.OrderType_OP_BUYSTOP,  pb.OrderType_OP_SELLSTOP:
    // ok, this is a pending order
default:
    return fmt.Errorf("ticket %d is not a pending order (type=%s)", ticket, ord.GetType())
}
```

---

## 🧰 3) Prepare params

Some brokers still check *slippage* field (points) even for delete; keep a small value for consistency.

```go
slippage := int32(5) // points
```

---

## 🗑️ 4) Call `OrderDelete`

```go
resp, err := account.OrderDelete(
    ctx,
    ticket,
    &slippage, // optional; broker-dependent
)
if err != nil {
    return fmt.Errorf("delete pending failed: %w", err)
}
fmt.Printf("✅ Pending deleted: ticket=%d symbol=%s time=%s\n",
    ticket, ord.GetSymbol(), time.Now().Format(time.RFC3339))
```

> Internally, the account helper wraps retries on transient transport errors (`codes.Unavailable`).

---

## ⚠️ Pitfalls

* **Not a pending** → `OrderDelete` works only for pending orders; for market positions use `OrderClose`.
* **Ticket not found** → it may have been executed or canceled already; refresh open orders and history.
* **Freeze level** → some brokers restrict actions near market price; try again later or move price (modify) before delete.
* **Timeouts** → wrap in `context.WithTimeout` (3–6s) and retry only transport errors.

---

## 🔄 Variations

* **Bulk delete** all pendings on a symbol:

```go
opened, err := account.ShowOpenedOrders(ctx)
if err != nil { return err }
for _, o := range opened.GetOrders() {
    if o.GetSymbol() != symbol { continue }
    switch o.GetType() {
    case pb.OrderType_OP_BUYLIMIT, pb.OrderType_OP_SELLLIMIT,
         pb.OrderType_OP_BUYSTOP,  pb.OrderType_OP_SELLSTOP:
        if _, err := account.OrderDelete(ctx, o.GetTicket(), &slippage); err != nil {
            log.Printf("delete %d failed: %v", o.GetTicket(), err)
        }
    }
}
```

* **Conditional cleanup** (good‑till‑time): if `ord.GetExpiration().AsTime().Before(time.Now())` → delete.

---

## 📎 See also

* `PlacePendingOrder.md` — how to place a pending with correct price/expiry.
* `ModifyOrder.md` — change pending price/expiration instead of deleting.
* `CloseOrder.md` — for market positions.
