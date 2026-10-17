# Getting All Available Symbols

> **Request:** retrieve a list of all symbols (instruments) available in the terminal
> Returns all symbol names and their corresponding internal indices.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints all available symbols with indices.
svc.ShowAllSymbols(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

data, err := account.ShowAllSymbols(ctx)
if err != nil {
    log.Fatalf("❌ ShowAllSymbols error: %v", err)
}

for _, sym := range data.SymbolNameInfos {
    fmt.Printf("Symbol: %s, Index: %d\n", sym.GetSymbolName(), sym.GetSymbolIndex())
}

```

---

### Method Signature

```go
func (s *MT4Service) ShowAllSymbols(ctx context.Context)
func (a *MT4Account) ShowAllSymbols(ctx context.Context) (*pb.SymbolsData, error)
```

`MT4Account.ShowAllSymbols` and `MT4Account.Symbols` send the same `Symbols` request and return the same data.

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
| `SymbolNameInfos` | `[]*pb.SymbolNameInfo` | List of symbol names with their indices |

Each `*pb.SymbolNameInfo` includes:

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
* **Empty responses:** Treat an empty list as suspicious (e.g. the terminal has just started) and retry before acting on it.

---

## 🧪 Testing Suggestions

* **Happy path:** List is non‑empty and contains common pairs like `EURUSD`.
* **Edge:** Include known disabled symbols and verify they don’t break downstream logic.
* **Failure path:** Simulate no connection — expect error or empty slice handled gracefully.
