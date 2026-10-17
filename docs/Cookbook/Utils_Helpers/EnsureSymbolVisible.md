# 👁️ EnsureSymbolVisible (GoMT4)

**Goal:** make sure a symbol is **available/visible** in the terminal's *Market Watch* **before** you call `Quote`, `OrderSend`, `SymbolParams`, etc.

> Real code refs you already ship:
>
> * Market info demos: `examples/mt4/MT4_service.go` (`ShowAllSymbols`, `ShowSymbols`, `ShowSymbolParams`)
> * Account layer: `examples/mt4/MT4Account.go` (`Symbols`, `SymbolParams`, `Quote`, exported `MarketInfoClient`)
> * Docs: [ShowAllSymbols](../../Market_Info/ShowAllSymbols.md), [ShowSymbols](../../Market_Info/ShowSymbols.md), [ShowSymbolParams](../../Market_Info/ShowSymbolParams.md)

---

## ❓ Why this matters

* MT4 hides many instruments from *Market Watch* by default. Hidden symbols can cause errors like **`symbol not found`** or empty quotes.
* Brokers may add **suffixes**: `EURUSD` vs `EURUSD.m` / `.pro` are **different** instruments.

---

## ✅ What exists in this repo

Calls you can use to check availability:

* `account.Symbols(ctx)` — all symbol names of the broker (visible or not).
* `account.SymbolParams(ctx, symbol)` — returns an error if the symbol is unknown; its `IsSelected` / `Visible` fields tell whether it is in Market Watch.
* `account.Quote(ctx, symbol)` — similarly, returns an error for unknown symbols.
* The generated `MarketInfo.SymbolSelect` RPC (MQL4 `SymbolSelect`) adds a symbol to Market Watch. `MT4Account` has no wrapper for it, but exposes `MarketInfoClient` and the exported `ExecuteWithReconnect`, which attaches the `apikey` / `id` metadata.

> There is **no** `EnsureSymbolVisible(...)` method in the code. Below are optional helpers you can add to your own code.

---

## 🧪 Minimal check (use existing calls)

```go
func CheckSymbolAvailable(ctx context.Context, a *mt4.MT4Account, symbol string) error {
    // Fast probe via params; you can also use Quote
    if _, err := a.SymbolParams(ctx, symbol); err != nil {
        return fmt.Errorf("symbol %s is not available (wrong name or suffix): %w", symbol, err)
    }
    return nil
}
```

Using:

```go
if err := CheckSymbolAvailable(ctx, account, "EURUSD"); err != nil {
    log.Printf("⚠️ %v", err)
    log.Printf("Hint: list the exact names with account.Symbols(ctx)")
    return err
}
```

---

## 🧰 Optional helper: select the symbol in Market Watch

```go
// SelectSymbol adds the symbol to Market Watch via the generated SymbolSelect RPC.
func SelectSymbol(ctx context.Context, a *mt4.MT4Account, symbol string) error {
    reply, err := mt4.ExecuteWithReconnect(a, ctx,
        func(md metadata.MD) (*pb.SymbolSelectReply, error) {
            return a.MarketInfoClient.SymbolSelect(metadata.NewOutgoingContext(ctx, md),
                &pb.SymbolSelectRequest{Symbol: symbol, Select: true})
        },
        func(r *pb.SymbolSelectReply) *pb.Error { return r.GetError() },
    )
    if err != nil {
        return err
    }
    if !reply.GetData().GetSuccess() {
        return fmt.Errorf("SymbolSelect(%s) returned false", symbol)
    }
    return nil
}
```

---

## 🧰 Optional helper: find the broker's name (suffixes)

This utility looks the base name up in the broker's symbol list, accepting common suffixes:

```go
// EnsureSymbolVisible finds base (or base+suffix) in the broker's symbol list,
// selects it in Market Watch and returns the exact name.
func EnsureSymbolVisible(ctx context.Context, a *mt4.MT4Account, base string) (string, error) {
    syms, err := a.Symbols(ctx)
    if err != nil {
        return "", err
    }
    names := map[string]bool{}
    for _, s := range syms.GetSymbolNameInfos() {
        names[s.GetSymbolName()] = true
    }
    for _, s := range []string{base, base + ".m", base + ".pro", base + ".ecn"} {
        if names[s] {
            if err := SelectSymbol(ctx, a, s); err != nil {
                return "", err
            }
            return s, nil
        }
    }
    return "", fmt.Errorf("symbol %s not found; check the exact name in account.Symbols(ctx)", base)
}
```

Using:

```go
sym, err := EnsureSymbolVisible(ctx, account, "EURUSD")
if err != nil { return err }
q, err := account.Quote(ctx, sym) // safe to proceed
```

---

## 🧭 Manual step

1. Print the broker's names: `svc.ShowAllSymbols(ctx)`.
2. Check your broker's exact symbol name (including the suffix).
3. Write it in `examples/config/config.json` → `DefaultSymbol`.

---

## ⚠️ Pitfalls

* **Suffix mismatch**: `EURUSD` ≠ `EURUSD.m`. Always check the exact name.
* **Case**: use the name exactly as the broker lists it.
* **Market Watch limit**: MT4's Market Watch holds a limited number of symbols; select only what you need.

---

## 🔗 See also

* [SymbolParams](../Market_Info/SymbolParams.md) — instrument parameters (Digits/Point/VolumeStep/StopsLevel).
* [GetQuote](../Market_Info/GetQuote.md) — a one‑time quote.
* [PlaceMarketOrder](../Orders/PlaceMarketOrder.md) — application after validation of the symbol.
