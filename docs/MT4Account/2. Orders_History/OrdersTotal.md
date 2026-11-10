# OrdersTotal

> How many orders are currently open — a single integer count.

**API Information:**

* **Low-level API:** `MT4Account.OrdersTotal(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **Built on:** [`OpenedOrders`](./OpenedOrders.md) — it is a convenience, **not** a separate RPC: it returns `len(OpenedOrders().GetOrderInfos())`.
* **Returns:** `int32`

## 💬 Just the essentials

* **What it is.** MQL `OrdersTotal()` — the count of currently open orders (positions + pendings).
* **Why you need it.** A quick "do I have anything open?" check without pulling the list.
* **When to use.** Gates and caps: skip work when flat, refuse new entries past a limit.

---

## 🎯 Purpose

* Count open orders in one tiny call
* Gate logic on "flat vs has-positions"
* Enforce a max-open-orders cap

---

## Method Signature

```go
func (a *MT4Account) OrdersTotal(ctx context.Context) (int32, error)
```

It takes no request of its own — it counts the rows that [`OpenedOrders`](./OpenedOrders.md) returns.

**Returns:** `int32` — the number of open orders (and an `error`).

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. The call takes no other arguments. |

---

## ⬆️ Output

| Return | Type | Description |
|--------|------|-------------|
| count | `int32` | Number of currently open orders (positions + pendings) |

> On failure the call returns `0` and a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Positions + pendings** are both counted — it doesn't distinguish them (use [`OpenedOrders`](./OpenedOrders.md) for that).
* **Cheaper than the list** when you only need the number.
* **0 means flat** — a clean signal to start fresh.

---

## 🔗 Usage Examples

### 1) Skip work when flat

```go
func hasOpen(ctx context.Context, acc *mt4.MT4Account) (bool, error) {
    n, err := acc.OrdersTotal(ctx)
    if err != nil { return false, err }
    return n > 0, nil
}
```

### 2) Enforce a max-open cap before entering

```go
func underCap(ctx context.Context, acc *mt4.MT4Account, max int32) (bool, error) {
    n, err := acc.OrdersTotal(ctx)
    if err != nil { return false, err }
    if n >= max {
        fmt.Printf("at cap: %d/%d open — not entering\n", n, max)
        return false, nil
    }
    return true, nil
}
```

### 3) Wait until flat

```go
func waitFlat(ctx context.Context, acc *mt4.MT4Account, timeout time.Duration) (bool, error) {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        n, err := acc.OrdersTotal(ctx)
        if err != nil { return false, err }
        if n == 0 { return true, nil }
        time.Sleep(time.Second)
    }
    return false, nil
}
```

---

## ⚠️ Pitfalls

* **Count only.** It won't tell you *which* orders — pair with [`OpenedOrders`](./OpenedOrders.md) when you need detail.
* **Snapshot.** The number can change immediately after the read.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** On a flat account it returns `0`; after opening one order it returns `1`.
* **Cap.** Open N orders and assert `underCap` flips at the threshold.
* **Consistency.** It equals `len(OpenedOrders().GetOrderInfos())`.

---

## 📚 See also

* [OpenedOrders](./OpenedOrders.md) · [OpenedOrdersTickets](./OpenedOrdersTickets.md)
* **Back to:** [2. Orders & History — Overview](./Orders_History.Overview.md) · [MT4Account — Overview](../Overview.md)
