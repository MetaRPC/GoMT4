# Getting Tick Value, Size, and Contract Size

> **Request:** retrieve tick value, tick size, and contract size for multiple symbols
> Useful for calculating profit/loss and position sizing.

---

### Code Example

```go
// Using service wrapper
symbols := []string{"EURUSD", "XAUUSD"}
service.ShowTickValues(context.Background(), symbols)

// Or directly from MT4Account
data, err := mt4.TickValueWithSize(context.Background(), symbols)
if err != nil {
    log.Fatalf("Error retrieving tick values: %v", err)
}

for _, info := range data.Infos {
    fmt.Printf("\U0001F4B9 Symbol: %s\n  TickValue: %.5f\n  TickSize: %.5f\n  ContractSize: %.2f\n\n",
        info.GetSymbolName(),
        info.GetTradeTickValue(),
        info.GetTradeTickSize(),
        info.GetTradeContractSize(),
    )
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowTickValues(ctx context.Context, symbols []string)
```

---

## 🔽 Input

| Field     | Type              | Description                               |
| --------- | ----------------- | ----------------------------------------- |
| `symbols` | `[]string`        | List of trading symbols (e.g., "EURUSD"). |
| `ctx`     | `context.Context` | For timeout and cancellation management.  |

---

## ⬆️ Output

Returns a `TickValueWithSizeData` object containing:

| Field   | Type                            | Description                   |
| ------- | ------------------------------- | ----------------------------- |
| `Infos` | `[]TickValueWithSizeSymbolInfo` | Tick-related info per symbol. |

Each `TickValueWithSizeSymbolInfo` includes:

| Field               | Type      | Description                                        |
| ------------------- | --------- | -------------------------------------------------- |
| `SymbolName`        | `string`  | Name of the trading symbol (e.g., "EURUSD").       |
| `TradeTickValue`    | `float64` | Value of one tick in account currency.             |
| `TradeTickSize`     | `float64` | Smallest possible price movement.                  |
| `TradeContractSize` | `float64` | Number of units per lot (e.g., 100,000 for Forex). |

---

## 🎯 Purpose

Use this method to access **core trading calculations** such as:

* Estimating profit/loss per tick movement
* Determining pip/tick monetary value
* Building accurate position sizing formulas

Essential for both manual and automated trading systems.
