# 🧩 SymbolParams (GoMT4)

**Goal:** read **symbol parameters** (Digits, Point, LotStep, Min/Max lot, Stops/Freeze levels, etc.) and use them for rounding & validations.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`SymbolParams`)
> * Demo: `examples/mt4/MT4_service.go` (see *ShowSymbolParams* flow)

---

## ✅ 1) Preconditions

* Symbol exists and is visible in MT4 (*Market Watch* → Show All).
* `config.json` has a known default symbol or you pass one explicitly.

---

## 🔎 2) Read params

```go
p, err := account.SymbolParams(ctx, symbol)
if err != nil { return err }

fmt.Printf("%s: Digits=%d Point=%.10f LotStep=%.2f MinLot=%.2f MaxLot=%.2f\n",
    symbol,
    p.GetDigits(),
    p.GetPoint(),
    p.GetVolumeStep(),
    p.GetVolumeMin(),
    p.GetVolumeMax(),
)
fmt.Printf("StopsLevel=%d FreezeLevel=%d ContractSize=%.2f\n",
    p.GetStopsLevel(), p.GetFreezeLevel(), p.GetContractSize(),
)
```

**Common fields:**

* `Digits` —the number of decimal places in the price.
* `Point` — the tick value (usually `10^-Digits').
* `VolumeStep` — the volume step in lots.
* `VolumeMin`/`VolumeMax` — acceptable volume range.
* `StopsLevel` — minimum distance (in points) for SL/TP/postponement from the current price.
* `FreezeLevel` — the area around the current price, where modification/removal of the ban may be limited.
* `ContractSize` — the size of the contract (for calculating the cost of the item).

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
q, _ := account.Quote(ctx, symbol)

entry := q.GetAsk() // for BUY; use Bid for SELL
point := p.GetPoint()

// Desired SL/TP distances (points)
slDist := 20.0
tpDist := 40.0

sl := roundPrice(entry - slDist*point, int(p.GetDigits()))
tp := roundPrice(entry + tpDist*point, int(p.GetDigits()))

// Respect StopsLevel
if slDist < float64(p.GetStopsLevel()) || tpDist < float64(p.GetStopsLevel()) {
    return fmt.Errorf("SL/TP too close: StopsLevel=%d", p.GetStopsLevel())
}
```

For **pending** orders:

```go
pendingPrice := roundPrice(entry - 10*point, int(p.GetDigits())) // e.g., Buy Limit
if math.Abs((entry - pendingPrice)/point) < float64(p.GetStopsLevel()) {
    return fmt.Errorf("pending too close to market: StopsLevel=%d", p.GetStopsLevel())
}
```

---

## ⚠️ Pitfalls

* **Wrong precision** → always round prices to `Digits` and volumes to `VolumeStep`.
* **Too close SL/TP** → compare distances in **points** with `StopsLevel`.
* **FreezeLevel** → The broker may prohibit modifications close to the price; try a little further from the market.
* **Suffix mismatch** → `EURUSD` vs `EURUSD.m` — different tools.

---

## 🔗 See also

* `RoundVolumePrice.md` — rendered helpers for rounding.
* `PlaceMarketOrder.md` / `PlacePendingOrder.md` — using Digits/LotStep/StopsLevel in orders.
* `GetQuote.md` — get the current price for calculations.
