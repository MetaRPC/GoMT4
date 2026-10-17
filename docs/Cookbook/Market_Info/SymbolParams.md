# 🧩 SymbolParams (GoMT4)

**Goal:** read **symbol parameters** (Digits, Point, VolumeStep, Min/Max volume, Stops/Freeze levels, etc.) and use them for rounding & validations.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`SymbolParams`, `SymbolParamsMany`)
> * Demo: `examples/mt4/MT4_service.go` (`ShowSymbolParams`)

---

## ✅ 1) Preconditions

* `account` is connected; the symbol exists on the broker.
* `config.json` has a known default symbol or you pass one explicitly.

---

## 🔎 2) Read params

```go
p, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

fmt.Printf("%s: Digits=%d Point=%.10f VolumeStep=%.2f VolumeMin=%.2f VolumeMax=%.2f\n",
    symbol,
    p.GetDigits(),
    p.GetPoint(),
    p.GetVolumeStep(),
    p.GetVolumeMin(),
    p.GetVolumeMax(),
)
fmt.Printf("StopsLevel=%d FreezeLevel=%d ContractSize=%.2f\n",
    p.GetTradeStopsLevel(), p.GetTradeFreezeLevel(), p.GetTradeContractSize(),
)
```

**Common fields** (`*pb.SymbolParamsManyInfo`):

* `Digits` — the number of decimal places in the price.
* `Point` — the smallest price step (usually `10^-Digits`).
* `VolumeStep` — the volume step in lots.
* `VolumeMin`/`VolumeMax` — acceptable volume range.
* `TradeStopsLevel` — minimum distance (in points) of SL/TP/pending price from the current price.
* `TradeFreezeLevel` — the zone around the current price where modification/deletion may be refused.
* `TradeContractSize` — the size of the contract (for calculating position value).
* `TradeMode` — whether trading is allowed (`SYMBOL_TRADE_MODE_FULL`, `_CLOSEONLY`, …).

---

## 🧮 3) Helpers: rounding by params

```go
func alignVolume(v, step, min, max float64) float64 {
    if step <= 0 { return v }
    v = math.Floor(v/step+0.5) * step
    if v < min { v = min }
    if v > max { v = max }
    return v
}

func roundPrice(p float64, digits int) float64 {
    mul := math.Pow10(digits)
    return math.Round(p*mul) / mul
}
```

Usage:

```go
vol := alignVolume(0.13, p.GetVolumeStep(), p.GetVolumeMin(), p.GetVolumeMax())
price := roundPrice(1.092345, int(p.GetDigits()))
```

---

## 🛡️ 4) Validate SL/TP and pending price

```go
// Distance in points from current price (Bid for SELL, Ask for BUY)
q, err := account.Quote(ctx, symbol)
if err != nil { return err }

entry := q.GetAsk() // for BUY; use Bid for SELL
point := p.GetPoint()

// Desired SL/TP distances (points)
slDist := 20.0
tpDist := 40.0

sl := roundPrice(entry - slDist*point, int(p.GetDigits()))
tp := roundPrice(entry + tpDist*point, int(p.GetDigits()))

// Respect StopsLevel
if slDist < float64(p.GetTradeStopsLevel()) || tpDist < float64(p.GetTradeStopsLevel()) {
    return fmt.Errorf("SL/TP too close: StopsLevel=%d", p.GetTradeStopsLevel())
}
```

For **pending** orders:

```go
pendingPrice := roundPrice(entry - 10*point, int(p.GetDigits())) // e.g., Buy Limit
if math.Abs((entry - pendingPrice)/point) < float64(p.GetTradeStopsLevel()) {
    return fmt.Errorf("pending too close to market: StopsLevel=%d", p.GetTradeStopsLevel())
}
```

---

## ⚠️ Pitfalls

* **Wrong precision** → always round prices to `Digits` and volumes to `VolumeStep`.
* **Too close SL/TP** → compare distances in **points** with `TradeStopsLevel`.
* **FreezeLevel** → the broker may refuse modifications close to the price; try a little further from the market.
* **Suffix mismatch** → `EURUSD` vs `EURUSD.m` — different instruments.
* **Many symbols** → `SymbolParamsMany(ctx, nil)` returns all symbols; with a non-empty slice it filters by the **first** name only.

---

## 🔗 See also

* [RoundVolumePrice](../Utils_Helpers/RoundVolumePrice.md) — ready helpers for rounding.
* [PlaceMarketOrder](../Orders/PlaceMarketOrder.md) / [PlacePendingOrder](../Orders/PlacePendingOrder.md) — using Digits/VolumeStep/StopsLevel in orders.
* [GetQuote](GetQuote.md) — get the current price for calculations.
