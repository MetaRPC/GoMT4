# Symbols

> The list of instruments the server exposes — names and indices.

**API Information:**

* **Low-level API:** `MT4Account.Symbols(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.MarketInfo`
* **RPC:** `Symbols(SymbolsRequest) → SymbolsReply`
* **Payload:** `SymbolsData` → `GetSymbolNameInfos()` → `SymbolNameInfo`

## 💬 Just the essentials

* **What it is.** The catalog of tradable symbol names (with their indices).
* **Why you need it.** Discover/validate symbol names before quoting or trading.
* **When to use.** Startup discovery, building pickers, validating config.

---

## 🎯 Purpose

* Enumerate all available symbols
* Validate a configured symbol exists
* Build a symbol picker / autocomplete

---

## Method Signature

```go
func (a *MT4Account) Symbols(ctx context.Context) (*pb.SymbolsData, error)
```

**Request message:**

```protobuf
message SymbolsRequest {}   // no arguments
```

**Payload message:**

```protobuf
message SymbolsData { repeated SymbolNameInfo symbol_name_infos = 1; }

message SymbolNameInfo {
  string symbol_name = 1;
  int32  symbol_index = 2;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. The call takes no other arguments. |

---

## ⬆️ Output — `*pb.SymbolsData`

`GetSymbolNameInfos()` returns `[]*pb.SymbolNameInfo`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `symbol_name` | string | `GetSymbolName() string` | Symbol, e.g. `"EURUSD"` |
| `symbol_index` | int32 | `GetSymbolIndex() int32` | Server-side index |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Names, not prices.** For prices use [`Quote`](./Quote.md) / [`QuoteMany`](./QuoteMany.md); for contract specs use [`SymbolParams`](./SymbolParams.md).
* **Add to Market Watch** before trading a discovered symbol with [`SymbolSelect`](./SymbolSelect.md).
* **Large lists** — brokers may expose hundreds of symbols.

---

## 🔗 Usage Examples

### 1) Count and print the first few

```go
func listSymbols(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.Symbols(ctx)
    if err != nil { return err }
    all := d.GetSymbolNameInfos()
    fmt.Printf("%d symbols\n", len(all))
    for i, s := range all {
        if i >= 10 { break }
        fmt.Println(" -", s.GetSymbolName())
    }
    return nil
}
```

### 2) Validate a configured symbol exists

```go
func symbolExists(ctx context.Context, acc *mt4.MT4Account, want string) (bool, error) {
    d, err := acc.Symbols(ctx)
    if err != nil { return false, err }
    for _, s := range d.GetSymbolNameInfos() {
        if s.GetSymbolName() == want { return true, nil }
    }
    return false, nil
}
```

### 3) Filter by prefix (e.g. all USD crosses)

```go
func withPrefix(ctx context.Context, acc *mt4.MT4Account, prefix string) ([]string, error) {
    d, err := acc.Symbols(ctx)
    if err != nil { return nil, err }
    var out []string
    for _, s := range d.GetSymbolNameInfos() {
        if strings.HasPrefix(s.GetSymbolName(), prefix) {
            out = append(out, s.GetSymbolName())
        }
    }
    return out, nil
}
```

---

## ⚠️ Pitfalls

* **Suffixes vary.** Brokers may use `EURUSD`, `EURUSD.m`, `EURUSDpro`, etc. — validate rather than assume.
* **Not all are tradable** at all times (sessions/holidays).

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `len(GetSymbolNameInfos()) > 0`; your base symbol is present.
* **Validate.** `symbolExists` returns true for a known symbol, false for a bogus one.
* **Prefix.** Filtering by a common prefix returns a non-empty subset.

---

## 📚 See also

* [SymbolSelect](./SymbolSelect.md) · [SymbolParamsMany](./SymbolParamsMany.md) · [Quote](./Quote.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
