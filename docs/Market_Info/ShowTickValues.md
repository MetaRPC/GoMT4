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
    fmt.Printf("💹 Symbol: %s\n  TickValue: %.5f\n  TickSize: %.5f\n  ContractSize: %.2f\n\n",
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

Returns `*pb.TickValueWithSizeData` containing:

| Field   | Type                                | Description                   |
| ------- | ----------------------------------- | ----------------------------- |
| `Infos` | `[]*pb.TickValueWithSizeSymbolInfo` | Tick-related info per symbol. |

Each `*pb.TickValueWithSizeSymbolInfo` includes:

| Field               | Type      | Description                                |
| ------------------- | --------- | ------------------------------------------ |
| `SymbolName`        | `string`  | Trading symbol (e.g., "EURUSD").           |
| `TradeTickValue`    | `float64` | Value of one tick in account currency.     |
| `TradeTickSize`     | `float64` | Smallest possible price movement.          |
| `TradeContractSize` | `float64` | Units per lot (e.g., 100000 for major FX). |

---

## 🎯 Purpose

Access **core trading calculations** such as:

* Estimating profit/loss per tick movement
* Determining pip/tick monetary value
* Building accurate position sizing formulas

---

## 🧩 Notes & Tips

* **Currency context:** `TradeTickValue` is in the **account currency**. Cross-currency symbols will factor broker conversions.
* **From tick to pip:** If you need pip value, convert via symbol `Point/Digits` from `SymbolParams`.
* **Vector use:** Prefer this batched call over per-symbol queries when working with lists.

---

## ⚠️ Pitfalls

* **Empty input:** The API requires at least one symbol — calling with an empty slice returns an error.
* **Mixed asset classes:** Contract sizes differ across FX, metals, indices — don’t assume `100000` universally.
* **Precision:** Use appropriate formatting; keep raw floats for calculations.
