# SymbolParamsMany

> Contract specs for several symbols in one call.

**API Information:**

* **Low-level API:** `MT4Account.SymbolParamsMany(ctx, symbols)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **RPC:** `SymbolParamsMany(SymbolParamsManyRequest) → SymbolParamsManyReply`
* **Payload:** `SymbolParamsManyData` → `GetSymbolInfos()` → `SymbolParamsManyInfo`

## 💬 Just the essentials

* **What it is.** The batch version of [`SymbolParams`](./SymbolParams.md) — many symbols' trading rules at once.
* **Why you need it.** Pre-load specs for a basket without N round-trips.
* **When to use.** Startup caching of point/digits/step/tick-value for your traded symbols.

---

## 🎯 Purpose

* Fetch contract specs for a list of symbols
* Pre-cache normalization/sizing inputs
* Compare specs across a basket

---

## Method Signature

```go
func (a *MT4Account) SymbolParamsMany(ctx context.Context, symbols []string) (*pb.SymbolParamsManyData, error)
```

**Request message:**

```protobuf
message SymbolParamsManyRequest { repeated string symbols = 1; }
```

**Payload message:**

```protobuf
message SymbolParamsManyData { repeated SymbolParamsManyInfo symbol_infos = 1; }
```

Each `SymbolParamsManyInfo` has the same getters as [`SymbolParams`](./SymbolParams.md#output-pbsymbolparamsmanyinfo).

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbols` | `[]string` | yes | Symbols to describe, e.g. `[]string{"EURUSD","USDJPY"}`. |

---

## ⬆️ Output — `*pb.SymbolParamsManyData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `symbol_infos` | repeated SymbolParamsManyInfo | `GetSymbolInfos()` | One spec per symbol — see [`SymbolParams`](./SymbolParams.md) for the fields |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Match by `GetSymbolName()`**, not slice position — unknown symbols may be omitted.
* **Cache it.** Specs rarely change within a session; read once and reuse.
* **`SymbolParams` is this** with a one-element list.

---

## 🔗 Usage Examples

### 1) Cache specs for a basket

```go
func cacheSpecs(ctx context.Context, acc *mt4.MT4Account, symbols []string) (map[string]*pb.SymbolParamsManyInfo, error) {
    d, err := acc.SymbolParamsMany(ctx, symbols)
    if err != nil { return nil, err }
    out := make(map[string]*pb.SymbolParamsManyInfo, len(d.GetSymbolInfos()))
    for _, s := range d.GetSymbolInfos() {
        out[s.GetSymbolName()] = s
    }
    return out, nil
}
```

### 2) Print step/digits per symbol

```go
func specTable(ctx context.Context, acc *mt4.MT4Account, symbols []string) error {
    d, err := acc.SymbolParamsMany(ctx, symbols)
    if err != nil { return err }
    for _, s := range d.GetSymbolInfos() {
        fmt.Printf("%-8s digits=%d step=%.2f tickVal=%.5f\n",
            s.GetSymbolName(), s.GetDigits(), s.GetVolumeStep(), s.GetTradeTickValue())
    }
    return nil
}
```

### 3) Find the tightest min-volume symbol

```go
func smallestMinLot(ctx context.Context, acc *mt4.MT4Account, symbols []string) (string, float64, error) {
    d, err := acc.SymbolParamsMany(ctx, symbols)
    if err != nil { return "", 0, err }
    var sym string
    best := math.MaxFloat64
    for _, s := range d.GetSymbolInfos() {
        if s.GetVolumeMin() < best { best, sym = s.GetVolumeMin(), s.GetSymbolName() }
    }
    return sym, best, nil
}
```

---

## ⚠️ Pitfalls

* **Missing symbols.** Unknown/unselected symbols may not appear — don't index by position.
* **Large baskets** are heavier; request only what you need.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `len(GetSymbolInfos())` matches the valid symbol count; each has `GetPoint() > 0`.
* **Map.** The cached map has an entry per returned info.
* **Parity.** The entry for `"EURUSD"` matches a single [`SymbolParams`](./SymbolParams.md) call.

---

## 📚 See also

* [SymbolParams](./SymbolParams.md) · [TickValueWithSize](./TickValueWithSize.md) · [Symbols](./Symbols.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
