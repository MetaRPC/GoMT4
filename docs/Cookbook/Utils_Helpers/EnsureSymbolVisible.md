# 👁️ EnsureSymbolVisible (GoMT4)

**Goal:** make sure a symbol is **available/visible** in MT4 *Market Watch* **before** you call `Quote`, `OrderSend`, `SymbolParams`, etc.

> Real code refs you already ship:
>
> * Market info demos: `examples/mt4/MT4_service.go` (see *ShowAllSymbols*, *ShowSymbols*, *ShowSymbolParams*)
> * Account layer: `examples/mt4/MT4Account.go` (e.g., `SymbolParams`, `Quote`)
> * Docs: `docs/Market_Info/ShowAllSymbols.md`, `ShowSymbols.md`, `ShowSymbolParams.md`

---

## ❓ Why this matters

* MT4 hides many instruments by default. Hidden symbols often cause errors like **`symbol not found`** or empty quotes.
* Brokers may add **suffixes**: `EURUSD` vs `EURUSD.m` / `.pro` are **different** tools.

---

## ✅ What exists in this repo

Do you already have work calls that you can use to check availability?:

* `SymbolParams(ctx, symbol)` — if the symbol is unavailable/does not exist, an error will be returned.
* `Quote(ctx, symbol)` — similarly, returns an error/empty.
* Demo-*showing* a list of characters: `ShowAllSymbols` / `ShowSymbols` in `MT4_service.go'.

> Direct method `EnsureSymbolVisible(...)`there is no ** in the code**. Below is the **optional convenience** wrapper: not part of the auto‑generated pb, but a small utility that can be added to avoid copying checks.

---

## 🧪 Minimal check (use existing calls)

```go
func CheckSymbolAvailable(ctx context.Context, a *MT4Account, symbol string) error {
    // Fast probe via params; you can also use Quote
    if _, err := a.SymbolParams(ctx, symbol); err != nil {
        return fmt.Errorf("symbol %s is not available (hidden or wrong suffix): %w", symbol, err)
    }
    return nil
}
```

Using:

```go
if err := CheckSymbolAvailable(ctx, account, "EURUSD"); err != nil {
    log.Printf("⚠️ %v", err)
    log.Printf("Hint: open MT4 → Market Watch → Show All, and check actual symbol name (suffix)")
    return err
}
```

---

## 🧰 Optional convenience helper (not in pb; you may add)

This utility tries to pick up common suffixes if the base character is not found. This is an add-on, not part of your current API — add it as desired.

```go
// EnsureSymbolVisible tries base symbol, then common suffixes (.m, .pro, .ecn),
// and returns the first that exists.
func EnsureSymbolVisible(ctx context.Context, a *MT4Account, base string) (string, error) {
    candidates := []string{base, base + ".m", base + ".pro", base + ".ecn"}
    for _, s := range candidates {
        if _, err := a.SymbolParams(ctx, s); err == nil {
            return s, nil // found a visible/real symbol
        }
    }
    return "", fmt.Errorf("symbol %s not visible or not found; open Market Watch → Show All and check the exact name", base)
}
```

Using:

```go
sym, err := EnsureSymbolVisible(ctx, account, "EURUSD")
if err != nil { return err }
q, err := account.Quote(ctx, sym) // safe to proceed
```

---

## 🧭 Manual step (for MT4 user)

1. Open MT4 → **Market Watch**.
2. PCM → **Show All**.
3. Check your broker's exact symbol name (including the suffix).
4. Write it in `examples/config/config.json` → `DefaultSymbol`.

---

## ⚠️ Pitfalls

* **Suffix mismatch**: `EURUSD` ≠ `EURUSD.m`. Always check the exact name.
* **Hidden in Market Watch**: Turn on *Show All* before the first launch.
* **Wrong symbol case**: Names are case-sensitive in terms of suffixes for some brokers.

---

## 🔗 See also

* `SymbolParams.md ` — tool parameters (Digits/Point/LotStep/StopsLevel).
* `GetQuote.md ` is a one—time quote.
* `PlaceMarketOrder.md ` — application after validation of the symbol.
