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

---

## 🧩 Notes & Tips

* **All vs Market Watch:** This method should return the broker's *full* symbol catalog, not only currently visible symbols in Market Watch. For later quotes/orders, ensure visibility per symbol (use a helper like `EnsureSymbolVisible(name)`; it is safe to call repeatedly).
* **Indices are volatile:** `SymbolIndex` is an internal, session-scoped position. Do not persist it or use it as a stable ID. Always key by `SymbolName`.
* **Name normalization:** Brokers often append suffixes (e.g., `EURUSD.m`, `XAUUSD-RAW`). Treat them as distinct instruments. If you need grouping ("base" symbol), implement an explicit normalization function and document the rule.
* **No sort guarantee:** Do not assume returned order. If deterministic output matters, sort by `SymbolName` (case-insensitive) on the client.
* **Visibility pre-checks:** Before bulk price requests or order placement, pre-warm Market Watch by ensuring visibility for the target set to reduce first-quote latency.
* **UTF‑8 safety:** Some symbols include non-ASCII characters (indices, exotics). Ensure your logging and files are UTF‑8.

---

## ⚠️ Pitfalls

* **Very large catalogs:** Some brokers expose 1000+ instruments. Avoid printing them unbounded to stdout; page results or write to a file. Throttle any follow-up per-symbol RPCs.
* **Disabled/non-tradable entries:** The catalog can include symbols that are currently disabled for trading or quoting in your account type. Verify trade permissions and contract specifications before relying on them.
* **Empty-but-ok responses:** If the terminal connection is stale, you might receive an empty slice without a hard error. Add sanity checks (expect `EURUSD`/`GBPUSD` etc.) and log a warning.
* **Changing catalogs:** Brokers can add/remove symbols after server maintenance. Treat the result as a snapshot and refresh before long-running jobs.

---

## ⚡ Performance Notes

* **Batch follow-ups:** For metadata or quotes, prefer batch endpoints when available (e.g., `QuoteMany`) or cap concurrency to \~8–16 workers to avoid saturating the terminal.
* **Context discipline:** Use a bounded context (3–5s). On first call after terminal launch, allow a slightly higher timeout.
* **Memory hygiene:** Iterate directly over the returned slice; avoid unnecessary copies unless you sort/filter.

---

## 🧪 Testing Suggestions

* **Happy path:** Assert the list is non-empty and contains anchors (e.g., `EURUSD`).
* **Determinism check:** Apply a stable sort and snapshot to a golden file for UI components.
* **Failure path:** Simulate terminal down / broker maintenance; verify you handle empty result + warning gracefully.
* **Permission edge:** Include a symbol you know is disabled and ensure downstream code skips it safely.


