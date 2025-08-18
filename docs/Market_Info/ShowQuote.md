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

Returns `*pb.QuoteData` with fields:

| Field      | Type        | Description                 |
| ---------- | ----------- | --------------------------- |
| `Bid`      | `float64`   | Current bid price.          |
| `Ask`      | `float64`   | Current ask price.          |
| `DateTime` | `timestamp` | UTC timestamp of the quote. |

> **Spread:** If not exposed directly, compute as `Ask - Bid`. For points/pips, divide by the symbol’s `Point` (from `SymbolParams`).

---

## 🎯 Purpose

Retrieve **live market pricing** for a specific symbol. Useful for:

* Displaying real-time bid/ask prices
* Building dashboards or widgets
* Monitoring spreads and triggering alerts

---

## 🧩 Notes & Tips

* **Precision:** Print with instrument-specific decimals (e.g., 5 for EURUSD). Keep raw values for calculations.
* **Timestamp:** `DateTime` is UTC — format for display, log in UTC.

---

## ⚠️ Pitfalls

* **Zero/invalid values:** Check `Bid > 0 && Ask >= Bid`. Otherwise treat as stale.
* **Wrong symbol string:** Use the exact broker symbol (including suffixes).

---

## 🧪 Testing Suggestions

* **Happy path:** `EURUSD` → `Ask > Bid`, timestamp recent.
* **Error path:** Unknown/disabled symbol → return error or empty data handled gracefully.
