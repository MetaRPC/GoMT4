# Sending a Market or Pending Order

> **Request:** send a trade order (market or pending)
> Sends a new order using the specified parameters and receives back execution details.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Opens a market BUY of 0.1 lot (slippage 5, comment "Go order test", magic 123456, no SL/TP)
// and prints ticket/price/time.
// ⚠️ This places a real trade — use on demo or with caution.
svc.ShowOrderSendExample(ctx, "EURUSD")

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// proto.Int32 / proto.Float64 / proto.String come from google.golang.org/protobuf/proto
result, err := account.OrderSend(
    ctx,
    "EURUSD",
    pb.OrderSendOperationType_OC_OP_BUY, // order type: Buy
    0.1,                                 // lots
    nil,                                 // price (nil = market)
    proto.Int32(5),                      // slippage
    proto.Float64(1.0500),               // stop loss
    proto.Float64(1.0900),               // take profit
    proto.String("Go order test"),       // comment
    proto.Int32(123456),                 // magic number
    nil /* expiration */,
)
if err != nil {
    log.Fatalf("❌ OrderSend error: %v", err)
}

fmt.Printf("✅ Order opened! Ticket: %d, Price: %.5f, Time: %s\n",
    result.GetTicket(),
    result.GetPrice(),
    result.GetOpenTime().AsTime().Format("2006-01-02 15:04:05"),
)
```

---

### Method Signature

```go
func (s *MT4Service) ShowOrderSendExample(ctx context.Context, symbol string)
func (a *MT4Account) OrderSend(
    ctx context.Context,
    symbol string,
    operationType pb.OrderSendOperationType,
    volume float64,
    price *float64,
    slippage *int32,
    stoploss *float64,
    takeprofit *float64,
    comment *string,
    magicNumber *int32,
    expiration *timestamppb.Timestamp,
) (*pb.OrderSendData, error)
```

---

## 🔽 Input

Required (`ShowOrderSendExample` takes only `ctx` and `symbol`):

| Field       | Type                        | Description                         |
| ----------- | --------------------------- | ----------------------------------- |
| `ctx`       | `context.Context`           | Timeout or cancellation management. |
| `symbol`    | `string`                    | Trading symbol (e.g., "EURUSD").    |
| `operationType` | `pb.OrderSendOperationType` | `OrderSendOperationType_OC_OP_BUY`, `_SELL`, `_BUYLIMIT`, `_SELLLIMIT`, `_BUYSTOP`, `_SELLSTOP`. |
| `volume`    | `float64`                   | Order volume in lots (e.g., 0.1).   |

Optional (when calling `MT4Account.OrderSend` directly):

| Field        | Type                     | Description                                 |
| ------------ | ------------------------ | ------------------------------------------- |
| `price`      | `*float64`               | Price for pending orders; `nil` for market. |
| `slippage`   | `*int32`                 | Max slippage (points).                      |
| `stopLoss`   | `*float64`               | Stop Loss price.                            |
| `takeProfit` | `*float64`               | Take Profit price.                          |
| `comment`    | `*string`                | Optional order comment.                     |
| `magic`      | `*int32`                 | Magic number tag.                           |
| `expiration` | `*timestamppb.Timestamp` | Expiration (pending orders only; broker server time). |

---

## ⬆️ Output

Returns `*pb.OrderSendData`:

| Field         | Type                     | Description                         |
| ------------- | ------------------------ | ----------------------------------- |
| `Ticket`      | `int32`                  | Unique order ID assigned by MT4.    |
| `Volume`      | `float64`                | Volume of the order.                |
| `Price`       | `float64`                | Actual execution price.             |
| `OpenTime`    | `*timestamppb.Timestamp` | Time when order was executed (broker server time). |
| `OpenTimeUtc` | `*timestamppb.Timestamp` | The same moment in UTC.             |

---

## 🎯 Purpose

Place new trade orders (market or pending), controlling volume, price, and risk parameters. The result confirms ticket number, execution price, and open time for tracking/logging.

---

## 🧩 Notes & Tips

* **Timeouts:** The method sets a default 5s timeout if your `ctx` has no deadline — keep calls bounded.
* **Volume validation:** Ensure `volume` respects `VolumeMin/Max` and `VolumeStep` from `SymbolParams` before sending.
* **Pending orders:** `price` must be provided for pending types; for market orders it should be `nil`.
* **Types:** `slippage` is `*int32`, `magic` is `*int32`, `expiration` uses protobuf timestamp.

---

## ⚠️ Pitfalls

* **Not connected:** Without a connected terminal the call fails (terminal-not-found after the retries); after `Disconnect` the client returns `"not connected"`.
* **Rejected by broker:** Invalid SL/TP distances, disabled trading, or wrong price for pending orders will cause API errors.
* **Races:** Price can move between validation and send; expect slippage/requotes depending on broker settings.
