# TickValueWithSize

> Tick value, tick size and contract size for one or more symbols — the money-math inputs for risk sizing.

**API Information:**

* **Low-level API:** `MT4Account.TickValueWithSize(ctx, symbolNames)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **RPC:** `TickValueWithSize(TickValueWithSizeRequest) → TickValueWithSizeReply`
* **Payload:** `TickValueWithSizeData` → `GetInfos()` → `SymbolTickValueWithSizeInfo`

## 💬 Just the essentials

* **What it is.** The per-symbol tick value / tick size / contract size — a focused subset of [`SymbolParams`](./SymbolParams.md).
* **Why you need it.** Convert a pip distance into money for position sizing.
* **When to use.** Risk calculators, exposure/value math.

---

## 🎯 Purpose

* Read tick value (money per tick per lot)
* Read tick size (price per tick)
* Read contract size
* Compute pip value → size a lot by cash risk

---

## Method Signature

```go
func (a *MT4Account) TickValueWithSize(ctx context.Context, symbolNames []string) (*pb.TickValueWithSizeData, error)
```

**Request message:**

```protobuf
message TickValueWithSizeRequest { repeated string symbol_names = 1; }
```

**Payload message:**

```protobuf
message TickValueWithSizeData { repeated SymbolTickValueWithSizeInfo infos = 1; }

message SymbolTickValueWithSizeInfo {
  int32  index = ...;
  string symbol_name = ...;
  double trade_tick_value = ...;
  double trade_tick_size = ...;
  double trade_contract_size = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbolNames` | `[]string` | yes | Symbols to describe, e.g. `[]string{"EURUSD"}`. |

---

## ⬆️ Output — `*pb.TickValueWithSizeData`

`GetInfos()` returns `[]*pb.SymbolTickValueWithSizeInfo`:

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `symbol_name` | string | `GetSymbolName() string` | Symbol |
| `trade_tick_value` | double | `GetTradeTickValue() float64` | Money value of one tick per lot |
| `trade_tick_size` | double | `GetTradeTickSize() float64` | Price size of one tick |
| `trade_contract_size` | double | `GetTradeContractSize() float64` | Contract size |
| `index` | int32 | `GetIndex() int32` | Server index |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Pip value per lot** ≈ `GetTradeTickValue() × (pipSize / GetTradeTickSize())`.
* **Lighter than `SymbolParams`** if money math is all you need — but `SymbolParams` already includes these.
* **Match by `GetSymbolName()`**, not slice index.

---

## 🔗 Usage Examples

### 1) Print the money-math inputs

```go
func tickInfo(ctx context.Context, acc *mt4.MT4Account, symbols []string) error {
    d, err := acc.TickValueWithSize(ctx, symbols)
    if err != nil { return err }
    for _, i := range d.GetInfos() {
        fmt.Printf("%-8s tickVal=%.5f tickSize=%.5f contract=%.0f\n",
            i.GetSymbolName(), i.GetTradeTickValue(), i.GetTradeTickSize(), i.GetTradeContractSize())
    }
    return nil
}
```

### 2) Pip value per lot

```go
func pipValue(ctx context.Context, acc *mt4.MT4Account, symbol string, pipSize float64) (float64, error) {
    d, err := acc.TickValueWithSize(ctx, []string{symbol})
    if err != nil { return 0, err }
    for _, i := range d.GetInfos() {
        if i.GetSymbolName() == symbol && i.GetTradeTickSize() > 0 {
            return i.GetTradeTickValue() * (pipSize / i.GetTradeTickSize()), nil
        }
    }
    return 0, fmt.Errorf("no tick info for %s", symbol)
}
```

### 3) Lot for a cash risk over a stop

```go
func lotForRisk(ctx context.Context, acc *mt4.MT4Account, symbol string, pipSize, riskCash, slPips float64) (float64, error) {
    pv, err := pipValue(ctx, acc, symbol, pipSize)
    if err != nil || pv <= 0 || slPips <= 0 { return 0, err }
    return riskCash / (slPips * pv), nil
}
```

---

## ⚠️ Pitfalls

* **Divide-by-zero.** Guard `GetTradeTickSize() > 0` before dividing.
* **Missing symbols.** Unknown names may be omitted — match by name.
* **Still normalize** the resulting lot with [`SymbolParams`](./SymbolParams.md)'s volume step.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** For a liquid FX symbol, `GetTradeTickValue() > 0` and `GetTradeTickSize() > 0`.
* **Pip value.** `pipValue` is positive and matches the figure derived from `SymbolParams`.
* **Risk lot.** `lotForRisk` scales inversely with `slPips`.

---

## 📚 See also

* [SymbolParams](./SymbolParams.md) · [SymbolParamsMany](./SymbolParamsMany.md) · [MT4Sugar · Risk Sizing](../../MT4Sugar/4.%20Risk_Sizing.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
