# Order Operations — Overview

This section groups together methods for **managing orders**: creating, modifying, closing, deleting, and retrieving their details.
It’s the practical toolkit for working with active and historical trades.

---

## 📂 Methods in this Section

* [ShowOpenedOrders.md](ShowOpenedOrders.md)
  Get full details of all currently opened orders.

* [ShowOpenedOrderTickets.md](ShowOpenedOrderTickets.md)
  Stream only the **ticket IDs** of active orders.

* [ShowOrdersHistory.md](ShowOrdersHistory.md)
  Retrieve closed trades within a selected date range.

* [ShowOrderModifyExample.md](ShowOrderModifyExample.md)
  Example of updating SL/TP or other order parameters.

* [ShowOrderCloseByExample.md](ShowOrderCloseByExample.md)
  Demonstrates closing one order against another (`CloseBy`).

* [OrderCloseExample.md](OrderCloseExample.md)
  Standard order close example by ticket.

* [ShowOrderDeleteExample.md](ShowOrderDeleteExample.md)
  Example of deleting a pending order.

---

## ⚡ Example Workflow

```go
// Example: lifecycle of an order

// 1. Place an order (via Send) -> returns ticket
ticket := 123456

// 2. Modify order parameters (SL/TP)
svc.ShowOrderModifyExample(ctx, ticket)

// 3. Monitor open orders
svc.ShowOpenedOrders(ctx)

// 4. Option A: Close directly
svc.ShowOrderCloseExample(ctx, ticket)

// 5. Option B: Use CloseBy to offset positions
svc.ShowOrderCloseByExample(ctx, ticket)

// 6. Review historical performance
svc.ShowOrdersHistory(ctx)
```

---

## ✅ Best Practices

1. Always fetch **OpenedOrders** before attempting modifications/closures.
2. Use `ShowOrderModifyExample` instead of re-sending orders for SL/TP updates.
3. Delete pending orders you no longer need to keep account clean.
4. Store historical results using `ShowOrdersHistory` for analytics and compliance.

---

## 🎯 Purpose

The methods in this section are for **full lifecycle management of trades**:

* Open → Modify → Monitor → Close/Delete → Audit.
* Simplifies automation logic.
* Ensures robust error handling for every order action.

---

👉 Use this overview as a **map** and follow links to each `.md` file for complete details.
