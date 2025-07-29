# Getting a Quote for Symbol

> **Request:** fetch the latest quote for a given symbol
> Returns current bid/ask prices, spread, and time for a specified trading instrument.

---

### Code Example

```go
// Using service wrapper
service.ShowQuote(context.Background(), "EURUSD")

// Or directly from MT4Account
data, err := mt4.Quote(context.Background(), "EURUSD")
if err != nil {
    log.Fatalf("Error fetching quote: %v", err)
}

fmt.Printf("Bid: %.5f, Ask: %.5f, Time: %s\n",
    data.GetBid(),
    data.GetAsk(),
    data.GetDateTime().AsTime().Format("2006-01-02 15:04:05"),
)
```

---

### Method Signature

```go
func (s *MT4Service) ShowQuote(ctx context.Context, symbol string)
```

---

## 🔽 Input

Required:

| Field    | Type              | Description                          |
| -------- | ----------------- | ------------------------------------ |
| `ctx`    | `context.Context` | Context for timeout or cancellation. |
| `symbol` | `string`          | Trading symbol (e.g., "EURUSD").     |

---

## ⬆️ Output

Returns a **QuoteData** object with the following fields:

| Field      | Type        | Description                         |
| ---------- | ----------- | ----------------------------------- |
| `Bid`      | `float64`   | Current bid price                   |
| `Ask`      | `float64`   | Current ask price                   |
| `Spread`   | `float64`   | Spread between ask and bid (points) |
| `DateTime` | `timestamp` | UTC timestamp of the quote          |

---

## 🎯 Purpose

Use this method to retrieve **live market pricing** for a specific symbol.
Useful for:

* Displaying real-time bid/ask prices
* Building pricing widgets or dashboards
* Monitoring spreads and triggering alerts

Provides an up-to-date snapshot of market conditions for a trading instrument.
