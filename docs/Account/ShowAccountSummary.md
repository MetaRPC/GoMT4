# Getting an Account Summary

> **Request:** full account summary (`*pb.AccountSummaryData`) from MT4
> Fetch all core account metrics in a single call.

---

### Code Example

```go
summary, err := s.account.AccountSummary(ctx)
if err != nil {
    log.Printf("❌ AccountSummary error: %v", err)
    return
}
fmt.Printf("Account Summary: Balance=%.2f, Equity=%.2f, Currency=%s\n",
    summary.GetAccountBalance(),
    summary.GetAccountEquity(),
    summary.GetAccountCurrency())
```

---

### Method Signature

```go
func (s *MT4Service) ShowAccountSummary(ctx context.Context)
```

---

## 🔽Input

No required input parameters.

| Parameter | Type              | Description                                  |
| --------- | ----------------- | -------------------------------------------- |
| `ctx`     | `context.Context` | Context for timeout or cancellation control. |

---

## ⬆️Output

Prints selected fields from `*pb.AccountSummaryData` to console:

| Field                | Type     | Description                                   |
| -------------------- | -------- | --------------------------------------------- |
| `AccountBalance`     | `double` | Balance excluding open positions.             |
| `AccountEquity`      | `double` | Equity — balance including floating P/L.      |
| `AccountMargin`      | `double` | Currently used margin.                        |
| `AccountFreeMargin`  | `double` | Free margin available for opening new trades. |
| `AccountCurrency`    | `string` | Deposit currency (e.g. `"USD"`, `"EUR"`).     |
| `AccountLeverage`    | `int64`  | Leverage applied to the account.              |
| `AccountUserName`    | `string` | Account holder's name.                        |
| `AccountLogin`       | `int64`  | Account login ID.                             |
| `AccountCompanyName` | `string` | Broker's name or company.                     |

---

## 🎯Purpose

Retrieve and display key real-time account information. Typical uses:

* Showing account status in dashboards or CLI output
* Checking available margin and equity before placing trades
* Monitoring general account health and exposure

---

## 🧩 Notes & Tips

* **Bounded context:** Your implementation sets a 3s timeout if none is provided — keep it; bump to 5s on slow terminals.
* **Connection required:** Returns `"not connected"` if neither `Host` nor `ServerName` is set. Call `ConnectByHostPort`/`ConnectByServerName` first.
* **Reconnect behavior:** `ExecuteWithReconnect` retries on `codes.Unavailable` and `TERMINAL_INSTANCE_NOT_FOUND` with \~500ms backoff + jitter — short pauses are expected during restarts.
* **Currency vs P/L base:** `AccountCurrency` is the deposit currency; P/L can be in the quote currency of instruments — don’t mix in calculations.
* **Equity is a snapshot:** Re-query right before risk checks or order placement.
* **Leverage source of truth:** Use `AccountLeverage` from summary for margin math; avoid hardcoded broker values.
* **Formatting:** Round for display only; keep raw doubles for math.

---

## ⚠️ Pitfalls

* **Stale terminal:** With a disconnected terminal, values may be old without a hard error. Log connection state along with numbers.
* **Roll-over effects:** Swaps/commissions at roll-over can cause brief equity/balance divergences.
* **Type drift:** Match pb types exactly (`int64` for leverage/login). Mixing `int32`/`uint64` can bite later.

---

## 🧪 Testing Suggestions

* **Happy path:** Values are non-negative; equity ≈ balance on flat accounts.
* **Edge cases:** With open positions, equity ≠ balance; currency non-empty.
* **Failure path:** Simulate terminal down; expect error logged and no panic.
