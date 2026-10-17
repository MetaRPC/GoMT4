# Getting Basic Symbol List

> **Request:** retrieve symbol names and indices available in the terminal
> Returns a simplified list of available trading symbols and their corresponding internal indices.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints all available symbols with indices.
svc.ShowSymbols(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

data, err := account.Symbols(ctx)
if err != nil {
    log.Fatalf("❌ Symbols error: %v", err)
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
func (a *MT4Account) Symbols(ctx context.Context) (*pb.SymbolsData, error)
```

`MT4Account.Symbols` and `MT4Account.ShowAllSymbols` send the same request and return the same data.

---

## 🔽 Input

| Field | Type              | Description                   |
| ----- | ----------------- | ----------------------------- |
| `ctx` | `context.Context` | For timeout and cancellation. |

---

## ⬆️ Output

Returns `*pb.SymbolsData` with:

| Field             | Type                        | Description                     |
| ----------------- | --------------------------- | ------------------------------- |
| `SymbolNameInfos` | `[]*pb.SymbolNameInfo` | Pairs of symbol name and index. |

Each `*pb.SymbolNameInfo` includes:

| Field         | Type     | Description                         |
| ------------- | -------- | ----------------------------------- |
| `SymbolName`  | `string` | The name of the trading symbol.     |
| `SymbolIndex` | `int32`  | The internal index for that symbol. |

---

## 🎯 Purpose

Fetch a clean list of symbols from the terminal for:

* Populating symbol dropdowns and selectors
* Filtering instruments by index
* Lightweight enumeration for setup/diagnostics

---

## 🧩 Notes & Tips

* **Index stability:** `SymbolIndex` can change after terminal restarts/updates. Key by `SymbolName` if you need persistence.
* **Sorting:** The API doesn’t guarantee order. Sort client-side for deterministic UI.
* **Suffixes:** Treat broker suffixes (e.g., `EURUSD.m`) as distinct symbols; don’t auto-normalize.

---

## ⚠️ Pitfalls

* **Large lists:** Avoid printing thousands of lines to stdout; prefer paging or file output.
* **Empty result:** Treat an empty list as suspicious (e.g. the terminal has just started) — add a sanity check (expect common pairs).
