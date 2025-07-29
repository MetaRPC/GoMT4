# Getting Symbol Parameters

> **Request:** retrieve extended trading parameters for a symbol
> Provides detailed attributes such as precision, volume limits, currencies, and trade modes.

---

### Code Example

```go
// Using service wrapper
service.ShowSymbolParams(context.Background(), "EURUSD")

// Or directly from MT4Account
info, err := mt4.SymbolParams(context.Background(), "EURUSD")
if err != nil {
    log.Fatalf("Error retrieving symbol parameters: %v", err)
}

fmt.Println("\uD83D\uDCCA Symbol Parameters:")
fmt.Printf("• Symbol: %s\n", info.GetSymbolName())
fmt.Printf("• Description: %s\n", info.GetSymDescription())
fmt.Printf("• Digits: %d\n", info.GetDigits())
fmt.Printf("• Volume Min: %.2f\n", info.GetVolumeMin())
fmt.Printf("• Volume Max: %.2f\n", info.GetVolumeMax())
fmt.Printf("• Volume Step: %.2f\n", info.GetVolumeStep())
fmt.Printf("• Trade Mode: %s\n", tradeModeToString(info.GetTradeMode()))
fmt.Printf("• Currency Base: %s\n", info.GetCurrencyBase())
fmt.Printf("• Currency Profit: %s\n", info.GetCurrencyProfit())
fmt.Printf("• Currency Margin: %s\n", info.GetCurrencyMargin())
```

---

### Method Signature

```go
func (s *MT4Service) ShowSymbolParams(ctx context.Context, symbol string) error
```

---

## 🔽 Input

| Field    | Type              | Description                           |
| -------- | ----------------- | ------------------------------------- |
| `symbol` | `string`          | Trading symbol (e.g., "EURUSD").      |
| `ctx`    | `context.Context` | For timeout and cancellation control. |

---

## ⬆️ Output

Returns a single `SymbolParamsInfo` object containing:

| Field            | Type      | Description                                      |
| ---------------- | --------- | ------------------------------------------------ |
| `SymbolName`     | `string`  | Name of the symbol                               |
| `SymDescription` | `string`  | Descriptive name or label for the symbol         |
| `Digits`         | `int32`   | Number of decimal places                         |
| `VolumeMin`      | `float64` | Minimum allowed lot volume                       |
| `VolumeMax`      | `float64` | Maximum allowed lot volume                       |
| `VolumeStep`     | `float64` | Minimum lot increment                            |
| `SpreadFloat`    | `float64` | Current floating spread in points                |
| `Bid`            | `float64` | Current bid price                                |
| `CurrencyBase`   | `string`  | Base currency of the symbol                      |
| `CurrencyProfit` | `string`  | Profit currency for trades in this symbol        |
| `CurrencyMargin` | `string`  | Margin currency used for this symbol             |
| `TradeMode`      | `int32`   | Trade mode (e.g., disabled, long-only, etc.)     |
| `TradeExeMode`   | `int32`   | Execution mode (e.g., market, instant execution) |

---

## 🎯 Purpose

Use this method to retrieve a **comprehensive profile** of a trading instrument, including trading rules, volume constraints, and precision. Useful for:

* Validating trade requests and constraints
* Displaying full instrument configurations
* Enabling condition-aware trading decisions

