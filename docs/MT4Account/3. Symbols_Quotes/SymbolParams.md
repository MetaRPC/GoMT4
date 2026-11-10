# SymbolParams (RECOMMENDED)

> One symbol's contract spec — point, digits, volume step/limits, tick value/size, stops level.

**API Information:**

* **Low-level API:** `MT4Account.SymbolParams(ctx, symbol)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **Built on:** [`SymbolParamsMany`](./SymbolParamsMany.md) (`mt4_term_api.AccountHelper`) — it requests the one symbol and returns the first info.
* **Payload:** `SymbolParamsManyInfo`

## 💬 Just the essentials

* **What it is.** The trading rules for a symbol — everything you need to normalize volume/price and size risk.
* **Why you need it.** Avoid invalid-volume (131) / invalid-stops (130) rejections.
* **When to use.** Before sizing a lot, rounding a price, or setting SL/TP.

---

## 🎯 Purpose

* Read `Point` and `Digits` for pip math and price rounding
* Read `VolumeMin` / `VolumeMax` / `VolumeStep` for lot normalization
* Read `TradeTickValue` / `TradeTickSize` for risk sizing
* Read `TradeStopsLevel` for minimum SL/TP distance

---

## Method Signature

```go
func (a *MT4Account) SymbolParams(ctx context.Context, symbol string) (*pb.SymbolParamsManyInfo, error)
```

**How it works:** it calls [`SymbolParamsMany`](./SymbolParamsMany.md) with a single-element list and returns `GetSymbolInfos()[0]`.

**Payload message (excerpt):**

```protobuf
message SymbolParamsManyInfo {
  string symbol_name = ...;
  int32  digits = ...;
  double point = ...;
  double volume_min = ...;  double volume_max = ...;  double volume_step = ...;
  double trade_tick_value = ...;  double trade_tick_size = ...;
  int32  trade_stops_level = ...;
  double trade_contract_size = ...;
  string currency_base = ...;  string currency_profit = ...;
  double bid = ...;  double ask = ...;
  bool   spread_float = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. |
| `symbol` | `string` | yes | Instrument, e.g. `"EURUSD"`. |

---

## ⬆️ Output — `*pb.SymbolParamsManyInfo`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `symbol_name` | string | `GetSymbolName() string` | Symbol |
| `digits` | int32 | `GetDigits() int32` | Price decimal places |
| `point` | double | `GetPoint() float64` | Point size |
| `volume_min` / `volume_max` | double | `GetVolumeMin()` / `GetVolumeMax()` | Lot limits |
| `volume_step` | double | `GetVolumeStep() float64` | Lot increment |
| `trade_tick_value` | double | `GetTradeTickValue() float64` | Money value of one tick per lot |
| `trade_tick_size` | double | `GetTradeTickSize() float64` | Price size of one tick |
| `trade_stops_level` | int32 | `GetTradeStopsLevel() int32` | Min SL/TP distance (points) |
| `trade_contract_size` | double | `GetTradeContractSize() float64` | Contract size |
| `currency_base` / `currency_profit` | string | `GetCurrencyBase()` / `GetCurrencyProfit()` | Currencies |
| `bid` / `ask` | double | `GetBid()` / `GetAsk()` | Last known prices |
| `spread_float` | bool | `GetSpreadFloat() bool` | Floating vs fixed spread |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Pip size** = `GetPoint()` × 10 on 3/5-digit symbols, else `GetPoint()`.
* **Normalize volume** to `GetVolumeStep()` and clamp to `[VolumeMin, VolumeMax]` — or the broker rejects it (131).
* **Pip value per lot** ≈ `GetTradeTickValue() × (pip / GetTradeTickSize())` — the basis of risk sizing.
* **Respect `GetTradeStopsLevel()`** when placing SL/TP (130 otherwise).
* The high-level [`MT4Sugar`](../../MT4Sugar/3.%20Symbol_Math.md) wraps all of this.

---

## 🔗 Usage Examples

### 1) Print the trading rules

```go
func rules(ctx context.Context, acc *mt4.MT4Account, symbol string) error {
    p, err := acc.SymbolParams(ctx, symbol)
    if err != nil { return err }
    fmt.Printf("%s  digits=%d point=%.5f  vol[%.2f..%.2f step %.2f]  stopsLvl=%d\n",
        p.GetSymbolName(), p.GetDigits(), p.GetPoint(),
        p.GetVolumeMin(), p.GetVolumeMax(), p.GetVolumeStep(), p.GetTradeStopsLevel())
    return nil
}
```

### 2) Normalize a lot to the step and limits

```go
func normalizeLot(ctx context.Context, acc *mt4.MT4Account, symbol string, lot float64) (float64, error) {
    p, err := acc.SymbolParams(ctx, symbol)
    if err != nil { return 0, err }
    step := p.GetVolumeStep()
    if step <= 0 { step = 0.01 }
    v := math.Round(lot/step) * step
    if v < p.GetVolumeMin() { v = p.GetVolumeMin() }
    if p.GetVolumeMax() > 0 && v > p.GetVolumeMax() { v = p.GetVolumeMax() }
    return v, nil
}
```

### 3) Pip value per lot (for risk sizing)

```go
func pipValuePerLot(ctx context.Context, acc *mt4.MT4Account, symbol string) (float64, error) {
    p, err := acc.SymbolParams(ctx, symbol)
    if err != nil { return 0, err }
    pip := p.GetPoint()
    if p.GetDigits() == 3 || p.GetDigits() == 5 { pip *= 10 }
    if p.GetTradeTickSize() > 0 {
        return p.GetTradeTickValue() * (pip / p.GetTradeTickSize()), nil
    }
    return p.GetTradeTickValue(), nil
}
```

---

## ⚠️ Pitfalls

* **Don't hardcode digits/step.** They vary per symbol and broker — always read them.
* **Zero step guard.** Treat a `0` `GetVolumeStep()` as `0.01` to avoid divide-by-zero.
* **Stops level is in points**, not pips — convert when comparing to a pip distance.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** `GetPoint() > 0`, `GetVolumeStep() > 0`, `GetDigits()` in {2,3,5}.
* **Normalize.** `normalizeLot` rounds `0.137 → 0.14` on a `0.01` step and clamps below `VolumeMin`.
* **Pip value.** `pipValuePerLot` is positive and finite for a liquid FX symbol.

---

## 📚 See also

* [SymbolParamsMany](./SymbolParamsMany.md) · [TickValueWithSize](./TickValueWithSize.md) · [Quote](./Quote.md) · [MT4Sugar · Symbol Math](../../MT4Sugar/3.%20Symbol_Math.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
