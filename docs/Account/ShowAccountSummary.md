# Getting an Account Summary

> **Request:** full account summary (`AccountSummaryData`) from MT4
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

Prints selected fields from `AccountSummaryData` to console:

| Field               | Type     | Description                                       |
| ------------------- | -------- | ------------------------------------------------- |
| `AccountBalance`    | `double` | Account balance excluding open positions.         |
| `AccountEquity`     | `double` | Equity — balance including floating P/L.          |
| `AccountMargin`     | `double` | Currently used margin.                            |
| `AccountFreeMargin` | `double` | Free margin available for opening new trades.     |
| `AccountCurrency`   | `string` | Account deposit currency (e.g. `"USD"`, `"EUR"`). |
| `AccountLeverage`   | `int`    | Leverage applied to the account.                  |
| `AccountName`       | `string` | Account holder's name.                            |
| `AccountNumber`     | `int`    | Account number (login ID).                        |
| `Company`           | `string` | Broker's name or company.                         |

---

## 🎯Purpose

This method is used to retrieve and display key real-time account information. It is typically used for:

* Showing account status in dashboards or CLI output
* Checking available margin and equity before placing trades
* Monitoring general account health and exposure

It is a fundamental method for any MT4 integration dealing with account monitoring or diagnostics.

---

## 🧩 Notes & Tips

* **Bounded context:** Always call with a timeout (e.g., 3–5s) to avoid hanging RPCs.
* **Currency vs P/L base:** `AccountCurrency` is the deposit currency; P/L reporting may use symbol quote currency. Do not mix them in formatting or math.
* **Equity is volatile:** Treat equity as a snapshot; for risk checks, use a short rolling window or re-query just-in-time before order placement.
* **Leverage source of truth:** Use leverage from the summary for margin math; do not hardcode per-broker assumptions.
* **Rounding/formatting:** For UI, round to instrument-appropriate precision. Keep raw doubles for calculations.

---

## ⚠️ Pitfalls

* **Stale terminal:** If the terminal is disconnected, you may get old values without an explicit error. Log the terminal connection state alongside the numbers.
* **Commissions/swaps timing:** Some brokers apply swaps/commissions at roll-over; short intervals may show counterintuitive equity vs. balance deltas.
* **Type drift:** Ensure your struct field types match proto definitions (e.g., `double` vs `float`, `int32` vs `int64`).
* **Time zones:** Values are server-side; when correlating with events, log both server time and UTC.

