# Getting Basic Symbol List

> **Request:** retrieve symbol names and indices available in the terminal
> Returns a simplified list of available trading symbols and their corresponding internal indices.

---

### Code Example

```go
// Using service wrapper
service.ShowSymbols(context.Background())

// Or directly from MT4Account
data, err := mt4.Symbols(context.Background())
if err != nil {
    log.Fatalf("Error fetching symbols: %v", err)
}

fmt.Println("=== Available Symbols ===")
for _, symbolInfo := range data.GetSymbolNameInfos() {
    fmt.Printf("Symbol: %s, Index: %d\n",
        symbolInfo.GetSymbolName(),
        symbolInfo.GetSymbolIndex())
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowSymbols(ctx context.Context)
```

---

## 🔽 Input

| Field | Type              | Description                   |
| ----- | ----------------- | ----------------------------- |
| `ctx` | `context.Context` | For timeout and cancellation. |

---

## ⬆️ Output

Returns a list of `SymbolNameIndexPair` structures:

| Field         | Type     | Description                         |
| ------------- | -------- | ----------------------------------- |
| `SymbolName`  | `string` | The name of the trading symbol.     |
| `SymbolIndex` | `int32`  | The internal index for that symbol. |

---

## 🎯 Purpose

Use this method to fetch a clean list of symbols from the terminal for:

* Populating symbol dropdowns and selectors
* Filtering instruments by index
* Lightweight symbol enumeration for setup and diagnostics
