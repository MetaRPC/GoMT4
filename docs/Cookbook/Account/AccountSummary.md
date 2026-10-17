# 📒 AccountSummary (GoMT4)

**Goal:** fetch the **account summary** (login, balance, equity, credit, leverage, currency, broker) and use it for basic checks.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`AccountSummary`, plus `AccountBalance`, `AccountEquity`, `AccountLeverage`, … helpers)
> * Demo: `examples/mt4/MT4_service.go` (`ShowAccountSummary`)

---

## ✅ 1) Preconditions

* `account` is connected (`ConnectByServerName` / `ConnectByHostPort` returned no error).
* Valid credentials in `examples/config/config.json` (or your own config).

---

## 📝 2) Request summary (with short timeout)

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

sum, err := account.AccountSummary(ctx)
if err != nil {
    return fmt.Errorf("account summary failed: %w", err)
}

fmt.Printf("Login=%d Currency=%s Leverage=1:%d Broker=%s\n",
    sum.GetAccountLogin(), sum.GetAccountCurrency(), sum.GetAccountLeverage(), sum.GetAccountCompanyName())
fmt.Printf("Balance=%.2f Equity=%.2f Credit=%.2f\n",
    sum.GetAccountBalance(), sum.GetAccountEquity(), sum.GetAccountCredit())
```

---

## 🔍 3) Inspect fields

Fields on `*pb.AccountSummaryData`:

* `AccountLogin`, `AccountUserName`, `AccountCompanyName`, `AccountCurrency`
* `AccountBalance`, `AccountEquity`, `AccountCredit`
* `AccountLeverage` (1\:x), `AccountTradeMode` (demo / contest / real), `IsInvestor`
* `ServerTime` / `ServerTimeUtc`, `UtcServerTimeShiftMinutes`, `ServerUtcOffsetSeconds`

> Margin, free margin and margin level are **not** in the summary. They arrive with the profit stream: `OnOpenedOrdersProfit(...)` → `AccountInfo` (`Margin`, `FreeMargin`, `MarginLevel`) — see [StreamOpenedOrderProfits](StreamOpenedOrderProfits.md).

---

## 🧪 4) Simple health/risk checks

```go
// 1) Account active?
if sum.GetAccountBalance() <= 0 && sum.GetAccountEquity() <= 0 {
    return fmt.Errorf("account looks inactive (balance/equity <= 0)")
}

// 2) Read-only login?
if sum.GetIsInvestor() {
    log.Printf("⚠️ investor password: trading calls will be refused")
}

// 3) Floating loss guard: equity far below balance
if bal := sum.GetAccountBalance(); bal > 0 && sum.GetAccountEquity() < 0.8*bal {
    log.Printf("⚠️ Equity is %.1f%% of balance", 100*sum.GetAccountEquity()/bal)
}
```

---

## 🔄 5) Periodic polling (no stream)

If you need periodic updates without streaming (each call has its own 3s timeout when `ctx` has no deadline):

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()
for {
    select {
    case <-ctx.Done():
        return
    case <-ticker.C:
        s, err := account.AccountSummary(ctx)
        if err != nil { log.Printf("summary err: %v", err); continue }
        log.Printf("Equity=%.2f Balance=%.2f", s.GetAccountEquity(), s.GetAccountBalance())
    }
}
```

---

## ⚠️ Pitfalls

* **Fresh login** → just after connect the server time/offset fields can still be unset (`ServerUtcOffsetSeconds == nil`).
* **Timeouts** → use 2–5s per call; the client already retries transport and terminal-not-found errors (see *Reliability*).
* **Currency conversion** → values are in **account currency**; convert before comparing across accounts.

---

## 📎 See also

* [StreamOpenedOrderProfits](StreamOpenedOrderProfits.md) — live profit updates for open orders (with margin figures).
* [HistoryOrders](../Orders/HistoryOrders.md) — check closed orders against balance changes.
* [Reliability](../../ReTimeouts_Reconnects_Backoff.md) — per-call timeouts & retries.
