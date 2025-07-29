# Getting All Available Symbols

> **Request:** retrieve a list of all symbols (instruments) available in the terminal
> Returns all symbol names and their corresponding internal indices.

---

### Code Example

```go
// Using service wrapper
service.ShowAllSymbols(context.Background())

// Or directly from MT4Account
data, err := mt4.ShowAllSymbols(context.Background())
if err != nil {
    log.Fatalf("Error fetching all symbols: %v", err)
}

symbols := data.SymbolNameInfos
for _, sym := range symbols {
    fmt.Printf("Symbol: %s, Index: %d\n", sym.GetSymbolName(), sym.GetSymbolIndex())
}
```

---

### Method Signature

```go
func (s *MT4Service) ShowAllSymbols(ctx context.Context)
```

---

## 🔽 Input

No required parameters beyond:

| Field | Type              | Description                          |
| ----- | ----------------- | ------------------------------------ |
| `ctx` | `context.Context` | For timeout and cancellation control |

---

## ⬆️ Output

Returns a `SymbolNamesData` structure with:

| Field             | Type                    | Description                             |
| ----------------- | ----------------------- | --------------------------------------- |
| `SymbolNameInfos` | `[]SymbolNameIndexPair` | List of symbol names with their indices |

Each `SymbolNameIndexPair` includes:

| Field         | Type     | Description                  |
| ------------- | -------- | ---------------------------- |
| `SymbolName`  | `string` | Name of the trading symbol   |
| `SymbolIndex` | `int32`  | Internal index of the symbol |

---

## 🎯 Purpose

Use this method to enumerate all available trading instruments from the MT4 terminal.

This is useful for:

* Populating dropdown menus and symbol lists
* Generating watchlists or market scanners
* Performing bulk operations over instruments
