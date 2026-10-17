# Getting Symbol Parameters

> **Request:** retrieve extended trading parameters for a symbol
> Provides detailed attributes such as precision, volume limits, currencies, and trade modes.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints full symbol parameters in a readable format.
svc.ShowSymbolParams(ctx, "EURUSD")

// --- Low-level (direct account call) ---
// Preconditions: account is already connected.

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

info, err := account.SymbolParams(ctx, "EURUSD")
if err != nil {
    log.Fatalf("❌ SymbolParams error: %v", err)
}

fmt.Println("📊 Symbol Parameters:")
fmt.Printf("• Symbol: %s\n", info.GetSymbolName())
fmt.Printf("• Description: %s\n", info.GetSymDescription())
fmt.Printf("• Digits: %d\n", info.GetDigits())
fmt.Printf("• Volume Min: %.2f\n", info.GetVolumeMin())
fmt.Printf("• Volume Max: %.2f\n", info.GetVolumeMax())
fmt.Printf("• Volume Step: %.2f\n", info.GetVolumeStep())
fmt.Printf("• Trade Mode: %s\n", info.GetTradeMode()) // e.g. SYMBOL_TRADE_MODE_FULL
fmt.Printf("• Currency Base: %s\n", info.GetCurrencyBase())
fmt.Printf("• Currency Profit: %s\n", info.GetCurrencyProfit())
fmt.Printf("• Currency Margin: %s\n", info.GetCurrencyMargin())
```

---

### Method Signature

```go
func (s *MT4Service) ShowSymbolParams(ctx context.Context, symbol string) error
func (a *MT4Account) SymbolParams(ctx context.Context, symbol string) (*pb.SymbolParamsManyInfo, error)
func (a *MT4Account) SymbolParamsMany(ctx context.Context, symbols []string) (*pb.SymbolParamsManyData, error)
```

`SymbolParams` sends `SymbolParamsMany` with `SymbolName` set and returns the first entry. `SymbolParamsMany` returns all symbols when `symbols` is empty; otherwise it filters by the **first** symbol of the slice only (the request takes one name).

---

## 🔽 Input

| Field    | Type              | Description                           |
| -------- | ----------------- | ------------------------------------- |
| `symbol` | `string`          | Trading symbol (e.g., "EURUSD").      |
| `ctx`    | `context.Context` | For timeout and cancellation control. |

---

## ⬆️ Output

Returns `*pb.SymbolParamsManyInfo` with fields (key subset):

| Field            | Type      | Description                                 |
| ---------------- | --------- | ------------------------------------------- |
| `SymbolName`     | `string`  | Name of the symbol.                         |
| `SymDescription` | `string`  | Descriptive name/label.                     |
| `Digits`         | `int32`   | Number of decimal places.                   |
| `VolumeMin`      | `float64` | Minimum allowed lot volume.                 |
| `VolumeMax`      | `float64` | Maximum allowed lot volume.                 |
| `VolumeStep`     | `float64` | Minimum lot increment.                      |
| `CurrencyBase`   | `string`  | Base currency of the symbol.                |
| `CurrencyProfit` | `string`  | Profit currency for trades in this symbol.  |
| `CurrencyMargin` | `string`  | Margin currency for this symbol.            |
| `TradeMode`      | `pb.SP_ENUM_SYMBOL_TRADE_MODE` | Trade mode (`SYMBOL_TRADE_MODE_DISABLED`, `_LONGONLY`, `_SHORTONLY`, `_CLOSEONLY`, `_FULL`). |
| `Point`          | `float64` | Point size (smallest price step).           |
| `TradeStopsLevel`| `int32`   | Minimum SL/TP distance in points.           |
| `TradeFreezeLevel`| `int32`  | Freeze distance in points.                  |
| `TradeTickValue` / `TradeTickSize` / `TradeContractSize` | `float64` | Tick value, tick size, contract size. |
| `Bid` / `Ask` …  | `float64` | Current prices and many more MT4 symbol properties (see `pb.SymbolParamsManyInfo`). |

---

## 🎯 Purpose

Obtain a **comprehensive profile** of an instrument: precision, volume constraints, currencies, and trade mode — to validate orders and display instrument config.

---

## 🧩 Notes & Tips

* **Order validation:** Use `VolumeMin/Max/Step` and `Digits` to validate user inputs **before** `OrderSend`.
* **Rounding rule:** Round order volume to the nearest `VolumeStep` (never exceed `VolumeMax`).
* **Precision:** Format prices using `Digits`; do not hardcode decimals per symbol.
* **TradeMode usage:** If `TradeMode` indicates disabled/restricted, surface a clear message and skip order placement.

---

## ⚠️ Pitfalls

* **Broker differences:** Parameters may vary across accounts/servers for the same symbol.
* **Stale cache:** Don’t cache forever — refresh on reconnect or at session start.
* **Step mismatch:** Floating arithmetic can break step checks; compare with a small epsilon when validating steps.

---

## 🧪 Testing Suggestions

* **Happy path:** `EURUSD` returns non-empty description; digits match expected (e.g., 5).
* **Volume bounds:** Try `VolumeMin - ε` and `VolumeMax + ε` → validation rejects.
* **Mode edge:** Force a symbol with restricted `TradeMode` → UI/action must block placing orders.
