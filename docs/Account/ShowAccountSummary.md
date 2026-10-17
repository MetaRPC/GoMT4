# Getting an Account Summary

> **Request:** full account summary (`*pb.AccountSummaryData`) from MT4
> Fetch all core account metrics in a single call.

---

### Code Example

```go
// --- Quick use (service wrapper) ---
// Prints balance/equity/currency inside the method.
svc.ShowAccountSummary(ctx)

// --- Low-level (direct account call) ---
// Preconditions: account is connected via ConnectByServerName/ConnectByHostPort.

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // keep short; bump if your terminal is slow
defer cancel()

summary, err := account.AccountSummary(ctx)
if err != nil {
    log.Printf("❌ AccountSummary error: %v", err)
    return
}

fmt.Printf("Account Summary: Balance=%.2f, Equity=%.2f, Currency=%s\n",
    summary.GetAccountBalance(),
    summary.GetAccountEquity(),
    summary.GetAccountCurrency(),
)

```

---

### Method Signature

```go
func (s *MT4Service) ShowAccountSummary(ctx context.Context)
func (a *MT4Account) AccountSummary(ctx context.Context) (*pb.AccountSummaryData, error)
```

Single-value helpers on `MT4Account` (each calls `AccountSummary`): `AccountLogin`, `AccountBalance`, `AccountCredit`, `AccountEquity`, `AccountLeverage`, `AccountName`, `AccountCompany`, `AccountCurrency`.

---

## 🔽Input

No required input parameters.

| Parameter | Type              | Description                                  |
| --------- | ----------------- | -------------------------------------------- |
| `ctx`     | `context.Context` | Context for timeout or cancellation control. |

---

## ⬆️Output

`ShowAccountSummary` prints `AccountBalance`, `AccountEquity` and `AccountCurrency`. `account.AccountSummary` returns the whole `*pb.AccountSummaryData`:

| Field                       | Type                     | Description                                   |
| --------------------------- | ------------------------ | --------------------------------------------- |
| `AccountLogin`              | `int64`                  | Account login ID.                             |
| `AccountBalance`            | `float64`                | Balance excluding open positions.             |
| `AccountEquity`             | `float64`                | Equity — balance including floating P/L.      |
| `AccountCredit`             | `float64`                | Credit funds.                                 |
| `AccountCurrency`           | `string`                 | Deposit currency (e.g. `"USD"`, `"EUR"`).     |
| `AccountLeverage`           | `int64`                  | Leverage applied to the account.              |
| `AccountUserName`           | `string`                 | Account holder's name.                        |
| `AccountCompanyName`        | `string`                 | Broker's name or company.                     |
| `AccountTradeMode`          | `pb.EnumAccountTradeMode`| `ACCOUNT_TRADE_MODE_DEMO` / `_CONTEST` / `_REAL`. |
| `IsInvestor`                | `bool`                   | Logged in with the investor (read-only) password. |
| `ServerTime` / `ServerTimeUtc` | `*timestamppb.Timestamp` | Trade server time (broker time) and the same moment in UTC. |
| `UtcServerTimeShiftMinutes` | `int32`                  | Broker time offset from UTC, in minutes.      |
| `ServerUtcOffsetSeconds`    | `*int32`                 | Broker time offset from UTC, in seconds; `nil` while not known yet (then the `*Utc` fields are unset too). |

Margin and free margin are not part of the summary; the `OnOpenedOrdersProfit` stream reports them in `AccountInfo` (`Margin`, `FreeMargin`, `MarginLevel`).

---

## 🎯Purpose

Retrieve and display key real-time account information. Typical uses:

* Showing account status in dashboards or CLI output
* Checking balance and equity before placing trades
* Monitoring general account health and exposure

---

## 🧩 Notes & Tips

* **Bounded context:** The method sets a 3s timeout if your `ctx` has no deadline — keep it; pass a longer deadline if needed.
* **Connection required:** Call `ConnectByServerName`/`ConnectByHostPort` first. The call is routed to the terminal named by the account's id; if no terminal is connected under it, the call fails after the retries. After `Disconnect` the account returns `"not connected"`.
* **Reconnect behavior:** `ExecuteWithReconnect` retries on `codes.Unavailable`, `TERMINAL_INSTANCE_NOT_FOUND` and `TERMINAL_REGISTRY_TERMINAL_NOT_FOUND` with exponential backoff (300ms doubling up to 5s, ±200ms jitter, up to 10 attempts) — short pauses are expected during restarts.
* **Currency vs P/L base:** `AccountCurrency` is the deposit currency; P/L can be in the quote currency of instruments — don’t mix in calculations.
* **Equity is a snapshot:** Re-query right before risk checks or order placement.
* **Leverage source of truth:** Use `AccountLeverage` from summary for margin math; avoid hardcoded broker values.
* **Formatting:** Round for display only; keep raw doubles for math.

---

## ⚠️ Pitfalls

* **Stale terminal:** If the terminal lost its broker connection, values can be old. Compare `ServerTime` with the current time when freshness matters.
* **Roll-over effects:** Swaps/commissions at roll-over can cause brief equity/balance divergences.
* **Type drift:** Match pb types exactly (`int64` for leverage/login). Mixing `int32`/`uint64` can bite later.

---

## 🧪 Testing Suggestions

* **Happy path:** Values are non-negative; equity ≈ balance on flat accounts.
* **Edge cases:** With open positions, equity ≠ balance; currency non-empty.
* **Failure path:** Simulate terminal down; expect error logged and no panic.
