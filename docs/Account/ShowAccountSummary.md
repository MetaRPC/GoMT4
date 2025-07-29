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

