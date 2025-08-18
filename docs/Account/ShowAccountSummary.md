# Getting an Account Summary

> **Request:** retrieve a full account summary (`AccountSummaryData`) from MT4 in one call.

---

### Code Example

```go
// High-level helper: prints a formatted account summary
total, err := s.account.AccountSummary(ctx)
if err != nil {
    log.Printf("❌ AccountSummary error: %v", err)
    return
}
fmt.Printf("Account Summary: Balance=%.2f, Equity=%.2f, Currency=%s\n",
    total.GetAccountBalance(),
    total.GetAccountEquity(),
    total.GetAccountCurrency())
```

---

### Method Signature

```go
func (s *MT4Service) ShowAccountSummary(ctx context.Context)
```

---

## 🔽 Input

No required parameters apart from context.

| Parameter | Type              | Description                       |
| --------- | ----------------- | --------------------------------- |
| `ctx`     | `context.Context` | Controls timeout or cancellation. |

---

## ⬆️ Output

Prints selected fields from `AccountSummaryData`:

| Field               | Type     | Description                            |
| ------------------- | -------- | -------------------------------------- |
| `AccountBalance`    | `double` | Balance excluding open positions.      |
| `AccountEquity`     | `double` | Equity = balance + floating P/L.       |
| `AccountMargin`     | `double` | Currently used margin.                 |
| `AccountFreeMargin` | `double` | Free margin available for new trades.  |
| `AccountCurrency`   | `string` | Deposit currency (`USD`, `EUR`, etc.). |
| `AccountLeverage`   | `int`    | Leverage applied to account.           |
| `AccountName`       | `string` | Account holder’s name.                 |
| `AccountNumber`     | `int`    | Account login ID.                      |
| `Company`           | `string` | Broker’s company name.                 |

---

## 🎯 Purpose

Retrieve and display real-time account information. Useful for:

* Displaying account status in CLI or dashboards
* Checking margin/equity before trade placement
* Monitoring account health and exposure

---

## 🧩 Notes & Tips

* Always wrap calls in a bounded context to avoid blocking RPCs.
* Currency field reflects the deposit currency; do not assume P/L base currency.
* Equity can fluctuate rapidly with market prices; never rely on a single snapshot for risk calculations.
* Leverage is broker-defined; cross-check when using for margin math.

---

## ⚠️ Pitfalls

* Some brokers report slightly different equity/margin values due to swaps or commissions.
* Values are server-provided; delays may exist if the terminal connection is unstable.
* AccountName and Company are mostly informational; do not rely on them programmatically.
