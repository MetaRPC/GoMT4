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
```

---

## 🔽 Input

Required:

* **`ctx`** (`context.Context`) — context for managing timeout or cancellation.

---

## ⬆️ Output

This method prints opened order information to stdout and does not return data directly. The printed output includes:

| Field       | Type                 | Description                                |
| ----------- | -------------------- | ------------------------------------------ |
| `Ticket`    | `int32`              | Unique ticket ID for the order.            |
| `Symbol`    | `string`             | Trading symbol (e.g., "EURUSD").           |
| `Lots`      | `float64`            | Volume of the order in lots.               |
| `OpenPrice` | `float64`            | Price at which the order was opened.       |
| `Profit`    | `float64`            | Current floating profit/loss of the order. |
| `OrderType` | `ENUM_ORDER_TYPE_TF` | Type of the order (Buy, Sell, etc.).       |

---

### ENUM: `ENUM_ORDER_TYPE_TF`

| Value                  | Description        |
| ---------------------- | ------------------ |
| `OrderTypeTfBuy`       | Buy order          |
| `OrderTypeTfSell`      | Sell order         |
| `OrderTypeTfBuyLimit`  | Pending Buy Limit  |
| `OrderTypeTfSellLimit` | Pending Sell Limit |
| `OrderTypeTfBuyStop`   | Pending Buy Stop   |
| `OrderTypeTfSellStop`  | Pending Sell Stop  |

---

## 🎯 Purpose

Use this method to retrieve and display a list of all currently open orders. Useful for:

* Monitoring active positions
* Building command-line dashboards with real-time order info
* Analyzing trade exposure, floating profit/loss, and position distribution

Provides essential functionality for real-time monitoring in MT4 integrations.
