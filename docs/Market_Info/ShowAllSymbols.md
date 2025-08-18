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

| Field | Type              | Description                          |
| ----- | ----------------- | ------------------------------------ |
| `ctx` | `context.Context` | For timeout and cancellation control |

---

## ⬆️ Output

Returns `*pb.SymbolsData`:

| Field             | Type                        | Description                             |
| ----------------- | --------------------------- | --------------------------------------- |
| `SymbolNameInfos` | `[]*pb.SymbolNameIndexPair` | List of symbol names with their indices |

Each `*pb.SymbolNameIndexPair` includes:

| Field         | Type     | Description                  |
| ------------- | -------- | ---------------------------- |
| `SymbolName`  | `string` | Name of the trading symbol   |
| `SymbolIndex` | `int32`  | Internal index of the symbol |

---

## 🎯 Purpose

Enumerate all available trading instruments from the MT4 terminal. Useful for:

* Populating dropdown menus and symbol lists
* Building watchlists or market scanners
* Performing bulk operations across instruments

---

## 🧩 Notes & Tips

* **Indices are not stable:** `SymbolIndex` can change after terminal restarts or broker updates. Always use `SymbolName` as the key.
* **Broker suffixes:** Symbols may have suffixes like `EURUSD.m` or `DE40.cash`. Treat each as distinct — no auto-stripping.
* **Sorting:** The API does not guarantee order. Sort client‑side if you need deterministic lists.

---

## ⚠️ Pitfalls

* **Large catalogs:** Brokers may expose hundreds/thousands of instruments. Printing/logging all at once can flood output.
* **Disabled instruments:** Some returned symbols may not be tradable on your account type — check permissions before using.
* **Empty responses:** A stale/disconnected terminal can return an empty list without error. Add sanity checks.

---

## 🧪 Testing Suggestions

* **Happy path:** List is non‑empty and contains common pairs like `EURUSD`.
* **Edge:** Include known disabled symbols and verify they don’t break downstream logic.
* **Failure path:** Simulate no connection — expect error or empty slice handled gracefully.
