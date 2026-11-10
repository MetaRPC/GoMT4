# 2. Orders & History — Overview

> The read side of trading: what's open now, and what has closed. (The write side — send/modify/close — is [4. Trading Operations](../4.%20Trading_Operations/Trading_Operations.Overview.md).)

In MT4 an open **market** order *is* the position; pending orders sit in the same list, told apart by `GetOrderType()` ([OpenedOrderType](../../All_Guides/ENUMS_USAGE_REFERENCE.md): 0 Buy, 1 Sell = market; 2–5 = pending).

---

## Methods in this group

| Method | Returns | One-liner |
|--------|---------|-----------|
| [`OpenedOrders(ctx)`](./OpenedOrders.md) | `*pb.OpenedOrdersData` | Open market + pending orders, with prices, SL/TP, profit. |
| [`OpenedOrdersTickets(ctx)`](./OpenedOrdersTickets.md) | `*pb.OpenedOrdersTicketsData` | Just the open ticket numbers. |
| [`OrdersHistory(ctx, sort, from, to, page, items)`](./OrdersHistory.md) | `*pb.OrdersHistoryData` | Closed/cancelled orders over a window (paged). |
| [`OrderSelect(ctx, ticket)`](./OrderSelect.md) | `*pb.OpenedOrderInfo` | One open order by ticket. |
| [`OrdersTotal(ctx)`](./OrdersTotal.md) | `int32` | Count of open orders. |

---

**Back to:** [MT4Account — Overview](../Overview.md)
