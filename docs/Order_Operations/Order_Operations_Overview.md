# Order Operations — Overview

This section groups together methods for **managing orders**: creating, modifying, closing, deleting, and retrieving their details.
It’s the practical toolkit for working with active and historical trades.

---

## 📂 Methods in this Section

* [ShowOrderSendExample.md](ShowOrderSendExample.md)  
  Example of placing a new order.

* [ShowOrderModifyExample.md](ShowOrderModifyExample.md)  
  Example of updating SL/TP or other order parameters.

* [ShowOpenedOrders.md](ShowOpenedOrders.md)  
  Get full details of all currently opened orders.

* [ShowOpenedOrderTickets.md](ShowOpenedOrderTickets.md)  
  Get only the **ticket IDs** of active orders.

* [ShowOrderCloseExample.md](ShowOrderCloseExample.md)  
  Standard order close example by ticket.

* [ShowOrderCloseByExample.md](ShowOrderCloseByExample.md)  
  Demonstrates closing one order against another (`CloseBy`).

* [ShowOrderDeleteExample.md](ShowOrderDeleteExample.md)  
  Example of deleting a pending order.

* [ShowOrdersHistory.md](ShowOrdersHistory.md)  
  Retrieve historical (closed/deleted) orders within a selected date range.


---

## ⚡ Example Workflow

```go
// Example: lifecycle of an order

// 1. Place an order (via Send) -> returns ticket
res, err := account.OrderSend(ctx, "EURUSD", pb.OrderSendOperationType_OC_OP_BUY, 0.10,
    nil, nil, nil, nil, nil, nil, nil)
if err != nil {
    log.Fatal(err)
}
ticket := res.GetTicket() // int32

// 2. Modify order parameters (SL/TP; the demo sets SL=1.0500, TP=1.0900)
svc.ShowOrderModifyExample(ctx, ticket)

// 3. Monitor open orders
svc.ShowOpenedOrders(ctx)

// 4. Option A: Close directly
svc.ShowOrderCloseExample(ctx, ticket)

// 5. Option B: Use CloseBy to offset positions (needs an opposite order's ticket)
var oppositeTicket int32 = 654321
svc.ShowOrderCloseByExample(ctx, ticket, oppositeTicket)

// 6. Review historical performance
svc.ShowOrdersHistory(ctx)
```

---

## ✅ Best Practices

1. Always fetch **OpenedOrders** before attempting modifications/closures.
2. Use `OrderModify` instead of re-sending orders for SL/TP updates.
3. Delete pending orders you no longer need to keep account clean.
4. Store historical results using `OrdersHistory` / `OrdersHistoryStream` for analytics and compliance.

---

## 🎯 Purpose

The methods in this section are for **full lifecycle management of trades**:

* Open → Modify → Monitor → Close/Delete → Audit.
* Simplifies automation logic.
* Every `MT4Account` order method returns a Go error (transport or API error) you can handle.

---

👉 Use this overview as a **map** and follow links to each `.md` file for complete details.
