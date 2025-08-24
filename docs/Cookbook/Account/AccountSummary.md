# 📒 AccountSummary (GoMT4)

**Goal:** fetch **account summary** (balance, equity, margin, free margin, leverage, currency) and use it for basic checks.

> Real code refs:
>
> * Account: `examples/mt4/MT4Account.go` (`AccountSummary`)
> * Demo: `examples/mt4/MT4_service.go` (`ShowAccountSummary`)

---

## ✅ 1) Preconditions

* MT4 terminal is running & connected to broker.
* Valid credentials in `examples/config/config.json`.

---

## 📝 2) Request summary (with short timeout)

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

sum, err := account.AccountSummary(ctx)
if err != nil {
    return fmt.Errorf("account summary failed: %w", err)
}

fmt.Printf("Login=%d Currency=%s Leverage=1:%d\n",
    sum.GetLogin(), sum.GetCurrency(), sum.GetLeverage())
fmt.Printf("Balance=%.2f Equity=%.2f Margin=%.2f FreeMargin=%.2f\n",
    sum.GetBalance(), sum.GetEquity(), sum.GetMargin(), sum.GetMarginFree())
```

---

## 🔍 3) Inspect fields

Common fields on `pb.AccountSummary`:

* `Login`, `Name`, `Server`, `Currency`
* `Balance`, `Equity`, `Margin`, `MarginFree`, `MarginLevel`
* `Leverage` (1\:x)
* `Profit` (floating P/L)

> Use `MarginLevel` to trigger risk warnings; below broker thresholds, position operations may be limited.

---

## 🧪 4) Simple health/risk checks

```go
// 1) Terminal healthy?
if sum.GetBalance() <= 0 && sum.GetEquity() <= 0 {
    return fmt.Errorf("account looks inactive (balance/equity <= 0)")
}

// 2) Risk guard: margin level warning
if ml := sum.GetMarginLevel(); ml > 0 && ml < 100.0 {
    log.Printf("⚠️ Low margin level: %.1f%%", ml)
}
```

---

## 🔄 5) Periodic polling (no stream)

If you need periodic updates without streaming:

```go
ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()
for {
    select {
    case <-ctx.Done():
        return
    case <-ticker.C:
        s, err := account.AccountSummary(ctx)
        if err != nil { log.Printf("summary err: %v", err); continue }
        log.Printf("Equity=%.2f MarginFree=%.2f", s.GetEquity(), s.GetMarginFree())
    }
}
```

---

## ⚠️ Pitfalls

* **Fresh login** → just after connect, give MT4 a second to initialize.
* **Timeouts** → use 2–5s per call; retry only transport errors (see *Reliability*).
* **Currency conversion** → values are in **account currency**; convert before comparing across accounts.

---

## 📎 See also

* `StreamAccountProfits.md` — live profit updates for open orders.
* `HistoryOrders.md` — check closed orders against balance changes.
* `Reliability (en)` — per-call timeouts & retries.
