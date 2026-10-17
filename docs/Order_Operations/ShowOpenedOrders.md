# Getting Opened Orders

> **Request:** retrieve currently opened orders from MT4
> Fetch all active (non-closed) trade positions on the account.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints all currently opened orders with details.
svc.ShowOpenedOrders(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

ordersData, err := account.OpenedOrders(ctx)
if err != nil {
    log.Fatalf("❌ OpenedOrders error: %v", err)
}

for _, order := range ordersData.GetOrderInfos() {
    fmt.Printf("[%s] Ticket: %d | Symbol: %s | Lots: %.2f | OpenPrice: %.5f | Profit: %.2f\n",
        order.GetOrderType(),
        order.GetTicket(),
        order.GetSymbol(),
        order.GetLots(),
        order.GetOpenPrice(),
        order.GetProfit(),
    )
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowOpenedOrders(ctx context.Context)
func (a *MT4Account) OpenedOrders(ctx context.Context) (*pb.OpenedOrdersData, error)
```

Related: `account.OrderSelect(ctx, ticket)` finds one opened order by ticket (it calls `OpenedOrders` and searches locally), `account.OrdersTotal(ctx)` counts them.

---

## 🔽 Input

Required:

* **`ctx`** (`context.Context`) — context for managing timeout or cancellation.

---

## ⬆️ Output

`ShowOpenedOrders` prints opened order information to stdout and does not return data. `account.OpenedOrders` returns `*pb.OpenedOrdersData`, whose `OrderInfos` is a slice of `*pb.OpenedOrderInfo` (market positions and pending orders):

| Field        | Type                 | Description                                |
| ------------ | -------------------- | ------------------------------------------ |
| `Ticket`     | `int32`              | Unique ticket ID for the order.            |
| `Symbol`     | `string`             | Trading symbol (e.g., "EURUSD").           |
| `OrderType`  | `pb.OpenedOrderType` | Type of the order (Buy, Sell, etc.).       |
| `Lots`       | `float64`            | Volume of the order in lots.               |
| `OpenPrice`  | `float64`            | Price at which the order was opened.       |
| `StopLoss` / `TakeProfit` | `float64` | SL / TP levels (0 = not set).       |
| `Profit`     | `float64`            | Current floating profit/loss of the order. |
| `Swap` / `Commision` | `float64`    | Swap and commission (the field is spelled `Commision`). |
| `MagicNumber`| `int32`              | EA magic number.                           |
| `Comment`    | `string`             | Order comment.                             |
| `OpenTime` / `OpenTimeUtc` | `*timestamppb.Timestamp` | Open time (broker server time / UTC). |
| `ExpirationTime` / `ExpirationTimeUtc` | `*timestamppb.Timestamp` | Pending order expiration. |
| `AccountLogin` | `int64`            | Account the order belongs to.              |

---

### ENUM: `pb.OpenedOrderType`

| Value                         | Description        |
| ----------------------------- | ------------------ |
| `OpenedOrderType_OO_OP_BUY`       | Buy order          |
| `OpenedOrderType_OO_OP_SELL`      | Sell order         |
| `OpenedOrderType_OO_OP_BUYLIMIT`  | Pending Buy Limit  |
| `OpenedOrderType_OO_OP_SELLLIMIT` | Pending Sell Limit |
| `OpenedOrderType_OO_OP_BUYSTOP`   | Pending Buy Stop   |
| `OpenedOrderType_OO_OP_SELLSTOP`  | Pending Sell Stop  |

---

## 🎯 Purpose

Use this method to retrieve and display a list of all currently open orders. Useful for:

* Monitoring active positions
* Building command-line dashboards with real-time order info
* Analyzing trade exposure, floating profit/loss, and position distribution

Provides essential functionality for real-time monitoring in MT4 integrations.
