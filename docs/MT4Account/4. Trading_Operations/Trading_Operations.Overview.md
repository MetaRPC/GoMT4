# 4. Trading Operations — Overview

> The write path: send, modify, close and delete orders. These are the **only methods that change account state**.

They live on `mt4_term_api.TradingHelper`. Each takes explicit arguments (optional fields are pointers — pass `nil` to omit). Side-based one-liners with defaults (`BuyMarket`, `BuyLimit`, …) are one layer up in [`MT4Sugar`](../../MT4Sugar/Overview.md).

> ⚠️ **Not retried.** Unlike reads, trade calls are not auto-retried on transient errors — a timed-out send may have filled. On error, reconcile with [`OpenedOrders`](../2.%20Orders_History/OpenedOrders.md). See [Return Codes](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## Methods in this group

| Method | Returns | One-liner |
|--------|---------|-----------|
| [`OrderSend(ctx, symbol, op, volume, …)`](./OrderSend.md) | `*pb.OrderSendData` | Open a market or pending order. |
| [`OrderModify(ctx, ticket, price, sl, tp, exp)`](./OrderModify.md) | `bool` | Change price / SL / TP / expiry. |
| [`OrderClose(ctx, ticket, lots, price, slippage)`](./OrderClose.md) | `*pb.OrderCloseDeleteData` | Close a market order (full/partial) or delete a pending. |
| [`OrderDelete(ctx, ticket)`](./OrderDelete.md) | `*pb.OrderCloseDeleteData` | Delete a pending order. |
| [`OrderCloseBy(ctx, ticket, opposite)`](./OrderCloseBy.md) | `*pb.OrderCloseByData` | Close an order by an opposite one. |

---

**Back to:** [MT4Account — Overview](../Overview.md)
