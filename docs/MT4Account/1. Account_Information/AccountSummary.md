# AccountSummary (RECOMMENDED)

> Complete account snapshot (balance, equity, margin, leverage, currency, server time) in a single call.

**API Information:**

* **Low-level API:** `MT4Account.AccountSummary(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.AccountHelper`
* **RPC:** `AccountSummary(AccountSummaryRequest) → AccountSummaryReply`
* **Payload:** `AccountSummaryData`

## 💬 Just the essentials

* **What it is.** One round-trip that returns every essential account property.
* **Why you need it.** The most efficient way to read balance/equity/margin at once — one call instead of many single-field reads.
* **When to use.** At session start (banner / sanity check), and right before any risk decision or order.

---

## 🎯 Purpose

* RECOMMENDED way to read account state
* Balance, equity, credit in one round-trip
* Login, name, server, company
* Leverage and trade mode (DEMO/REAL)
* Currency, margin fields and server time

---

## Method Signature

```go
func (a *MT4Account) AccountSummary(ctx context.Context) (*pb.AccountSummaryData, error)
```

**Request message:**

```protobuf
message AccountSummaryRequest {}   // no arguments
```

**Payload message:**

```protobuf
message AccountSummaryData {
  int64  account_login = 1;
  double account_balance = 2;
  double account_equity = 3;
  string account_user_name = 4;
  int64  account_leverage = 5;
  EnumAccountTradeMode account_trade_mode = 6;
  string account_company_name = 7;
  string account_currency = 8;
  Timestamp server_time = 9;
  Timestamp server_time_utc = 15;
  double account_credit = 10;
  bool   is_investor = 11;
  int32  utc_server_time_shift_minutes = 12;
  int32  server_utc_offset_seconds = 14;
  double account_margin = 13;
  double account_margin_free = 16;
  double account_margin_level = 17;
  EnumAccountMarginMode account_margin_mode = 18;
  bool   account_trade_allowed = 19;
  bool   account_trade_expert = 20;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. The call takes no other arguments. |

---

## ⬆️ Output — `*pb.AccountSummaryData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `account_login` | int64 | `GetAccountLogin() int64` | Account number |
| `account_balance` | double | `GetAccountBalance() float64` | Balance (excludes floating P/L) |
| `account_equity` | double | `GetAccountEquity() float64` | Equity = balance + floating P/L |
| `account_credit` | double | `GetAccountCredit() float64` | Credit in deposit currency |
| `account_margin` | double | `GetAccountMargin() float64` | Margin used by open orders |
| `account_margin_free` | double | `GetAccountMarginFree() float64` | Free margin for new orders |
| `account_margin_level` | double | `GetAccountMarginLevel() float64` | `equity / margin × 100` (0 if no margin) |
| `account_margin_mode` | EnumAccountMarginMode | `GetAccountMarginMode()` | Margin mode (MT4 = hedging) |
| `account_leverage` | int64 | `GetAccountLeverage() int64` | Leverage (e.g. 100 = 1:100) |
| `account_currency` | string | `GetAccountCurrency() string` | Deposit currency |
| `account_user_name` | string | `GetAccountUserName() string` | Account owner |
| `account_company_name` | string | `GetAccountCompanyName() string` | Broker company |
| `account_trade_mode` | EnumAccountTradeMode | `GetAccountTradeMode()` | Demo / Contest / Real |
| `account_trade_allowed` | bool | `GetAccountTradeAllowed() bool` | Manual trading permitted |
| `account_trade_expert` | bool | `GetAccountTradeExpert() bool` | EA / API trading permitted |
| `is_investor` | bool | `GetIsInvestor() bool` | Read-only (investor) login |
| `server_time` / `server_time_utc` | Timestamp | `GetServerTime()` / `GetServerTimeUtc()` | Server time at the reply (server / UTC) |
| `server_utc_offset_seconds` | int32 | `GetServerUtcOffsetSeconds() int32` | Broker UTC offset (seconds) |

> On failure the call returns a non-nil `error` carrying the broker's code + message as text (`"API error (code=…): …"`); map it via the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧱 Related enums (from proto)

### `EnumAccountTradeMode` — field `account_trade_mode`

| Constant | Value | Meaning |
|----------|-------|---------|
| `ACCOUNT_TRADE_MODE_DEMO` | 0 | Demo / practice account |
| `ACCOUNT_TRADE_MODE_CONTEST` | 1 | Contest account |
| `ACCOUNT_TRADE_MODE_REAL` | 2 | Real (live) account |

### `EnumAccountMarginMode` — field `account_margin_mode`

| Constant | Value | Meaning |
|----------|-------|---------|
| `ACCOUNT_MARGIN_MODE_RETAIL_NETTING` | 0 | Retail, netting |
| `ACCOUNT_MARGIN_MODE_EXCHANGE` | 1 | Exchange margin model |
| `ACCOUNT_MARGIN_MODE_RETAIL_HEDGING` | 2 | Retail, hedging (typical MT4) |

> In Go, compare against the typed constant: `s.GetAccountTradeMode() == pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_REAL`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **One call beats many.** Balance, equity, margin, free margin and margin level come together — don't fan out into several calls.
* **`balance` vs `equity`.** Balance ignores open orders; equity includes their floating P/L. On a flat account they're equal; `equity − balance` is your current floating P/L.
* **Money is in the deposit currency.** Every monetary field is in `GetAccountCurrency()`.
* **`margin_level` is 0 with no open orders** (division by zero avoided) — guard before comparing to a stop-out threshold.
* **Read-only sessions.** `GetIsInvestor()` true or `GetAccountTradeAllowed()` false means trading calls will be rejected — check before placing orders.

---

## 🔗 Usage Examples

### 1) Session-start banner

```go
import pb "github.com/MetaRPC/GoMT4/package"

func printBanner(ctx context.Context, acc *mt4.MT4Account) error {
    s, err := acc.AccountSummary(ctx)
    if err != nil { return err }
    mode := map[pb.EnumAccountTradeMode]string{
        pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_DEMO:    "DEMO",
        pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_CONTEST: "CONTEST",
        pb.EnumAccountTradeMode_ACCOUNT_TRADE_MODE_REAL:    "REAL",
    }[s.GetAccountTradeMode()]
    floating := s.GetAccountEquity() - s.GetAccountBalance()
    fmt.Printf("Account #%d (%s) @ %s\n", s.GetAccountLogin(), mode, s.GetAccountCompanyName())
    fmt.Printf("Balance %.2f %s  Equity %.2f (floating %+.2f)  1:%d\n",
        s.GetAccountBalance(), s.GetAccountCurrency(),
        s.GetAccountEquity(), floating, s.GetAccountLeverage())
    return nil
}
```

### 2) Pre-trade risk gate

```go
func canOpen(ctx context.Context, acc *mt4.MT4Account, minFreeMargin, minMarginLevel float64) (bool, error) {
    s, err := acc.AccountSummary(ctx)
    if err != nil { return false, err }

    if !s.GetAccountTradeAllowed() || s.GetIsInvestor() {
        return false, nil // read-only / disabled account
    }
    if s.GetAccountMarginFree() < minFreeMargin {
        return false, nil
    }
    // margin level is 0 when nothing is open — only enforce once margin is in use
    if s.GetAccountMargin() > 0 && s.GetAccountMarginLevel() < minMarginLevel {
        return false, nil
    }
    return true, nil
}
```

### 3) Equity / drawdown monitor

```go
func watchDrawdown(ctx context.Context, acc *mt4.MT4Account, maxDDPercent float64) error {
    var peak float64
    for {
        s, err := acc.AccountSummary(ctx)
        if err != nil { return err }
        eq := s.GetAccountEquity()
        if eq > peak { peak = eq }
        dd := 0.0
        if peak > 0 { dd = (peak - eq) / peak * 100 }
        fmt.Printf("equity=%.2f peak=%.2f drawdown=%.2f%%\n", eq, peak, dd)
        if dd >= maxDDPercent {
            return fmt.Errorf("drawdown %.2f%% ≥ %.2f%% — de-risk", dd, maxDDPercent)
        }
        time.Sleep(5 * time.Second)
    }
}
```

---

## ⚠️ Pitfalls

* **Snapshot, not a stream.** Re-query right before a risk check or order; don't cache across time. For live equity use [`OnOpenedOrdersProfit`](../5.%20Streaming_Methods/OnOpenedOrdersProfit.md).
* **`margin_level == 0` with nothing open.** Guard `GetAccountMargin() > 0` before comparing to a stop-out threshold.
* **Currency mismatch.** All money fields are in `GetAccountCurrency()`; converting instrument P/L is on you.
* **Read-only / disabled accounts.** `GetIsInvestor()` or `GetAccountTradeAllowed() == false` means orders are rejected even though this call succeeds.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** Fields populated; `GetAccountCurrency()` non-empty; on a flat account `equity == balance` and `GetAccountMargin() == 0`.
* **With open orders.** `equity != balance`, `GetAccountMargin() > 0`, and `GetAccountMarginLevel()` is a sensible percentage.
* **Risk gate.** Exercise `canOpen()` against low free margin / investor session / trading disabled and assert it refuses.

---

## 📚 See also

* [TradeStats](./TradeStats.md) · [TradeUnrealizedPL](./TradeUnrealizedPL.md) · [EquityHistory](./EquityHistory.md) · [OnOpenedOrdersProfit](../5.%20Streaming_Methods/OnOpenedOrdersProfit.md)
* **Back to:** [1. Account Information — Overview](./Account_Information.Overview.md) · [MT4Account — Overview](../Overview.md)
