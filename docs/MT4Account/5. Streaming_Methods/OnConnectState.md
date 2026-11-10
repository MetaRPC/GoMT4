# OnConnectState

> Terminal connection-state events — know when a terminal goes up, down or reconnects.

**API Information:**

* **Low-level API:** `MT4Account.OnConnectState(ctx, terminalIds)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.Connection`
* **RPC:** `OnConnectState(OnConnectStateRequest) → stream OnConnectStateReply`
* **Payload:** `ConnectStateData`

## 💬 Just the essentials

* **What it is.** Health/heartbeat stream for the MT4 terminal(s) behind your session.
* **Why you need it.** Surface disconnects to the user; gate trading on a live link.
* **When to use.** Dashboards, watchdogs, "is the bridge alive?" checks.

---

## 🎯 Purpose

* Watch terminal connection state in real time
* Detect disconnects and read the reason
* Gate trading on a confirmed-live link
* Watch specific terminals or just your own

---

## Method Signature

```go
func (a *MT4Account) OnConnectState(ctx context.Context, terminalIds []string) (<-chan *pb.ConnectStateData, <-chan error)
```

**Request message:**

```protobuf
message OnConnectStateRequest { repeated string terminal_ids = 1; }
```

**Payload message:**

```protobuf
message ConnectStateData {
  string state = ...;
  bool   is_connected = ...;
  bool   is_alive = ...;  bool api_is_alive = ...;
  string server = ...;
  uint64 account = ...;
  string error_message = ...;
  Timestamp timestamp = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancel it to unsubscribe and close the channels. |
| `terminalIds` | `[]string` | no | Specific terminals to watch; `nil` = this account's own terminal. |

---

## ⬆️ Frame — `*pb.ConnectStateData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `state` | string | `GetState() string` | Human-readable state |
| `is_connected` | bool | `GetIsConnected() bool` | Link up to the broker |
| `is_alive` / `api_is_alive` | bool | `GetIsAlive()` / `GetApiIsAlive()` | Terminal / API process alive |
| `server` | string | `GetServer() string` | Broker server name |
| `account` | uint64 | `GetAccount() uint64` | Login |
| `error_message` | string | `GetErrorMessage() string` | Reason, when down |
| `timestamp` | Timestamp | `GetTimestamp()` | Event time |

> A non-nil value on the **error channel** ends the stream; both channels close on `ctx` cancel.

---

## 🧩 Notes & Tips

* **Lives on `ConnectionClient`** — not `SubscriptionClient` like the other `On*` streams.
* **`nil` terminalIds** = just watch your own terminal (the common case).
* **Check `GetErrorMessage()`** when `GetIsConnected()` is false for the reason.

---

## 🔗 Usage Examples

### 1) Watch the terminal's state

```go
func watchState(ctx context.Context, acc *mt4.MT4Account) error {
    ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
    defer cancel()
    data, errs := acc.OnConnectState(ctx, nil)
    for {
        select {
        case d := <-data:
            fmt.Printf("%s connected=%v alive=%v\n",
                d.GetState(), d.GetIsConnected(), d.GetIsAlive())
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 2) Flag disconnects with the reason

```go
func flagDisconnects(ctx context.Context, acc *mt4.MT4Account) error {
    data, errs := acc.OnConnectState(ctx, nil)
    for {
        select {
        case d := <-data:
            if !d.GetIsConnected() {
                fmt.Printf("⚠️ %s disconnected: %s\n", d.GetServer(), d.GetErrorMessage())
            }
        case err := <-errs:
            if err != nil { return err }
        case <-ctx.Done():
            return nil
        }
    }
}
```

### 3) Gate: only trade while connected

```go
func waitConnected(ctx context.Context, acc *mt4.MT4Account) (bool, error) {
    data, errs := acc.OnConnectState(ctx, nil)
    for {
        select {
        case d := <-data:
            if d.GetIsConnected() && d.GetApiIsAlive() { return true, nil }
        case err := <-errs:
            if err != nil { return false, err }
        case <-ctx.Done():
            return false, nil
        }
    }
}
```

---

## ⚠️ Pitfalls

* **Different client.** Needs a live `ConnectionClient` — it won't stream before connect.
* **Watchdog, not a trigger** — use it to gate, but also handle errors on the actual calls.
* **Leaking the context** leaks the goroutine — `defer cancel()`.

---

## 🧪 Testing Suggestions

* **Happy path (demo).** On a healthy session, `GetIsConnected()` is true and `GetState()` is a normal value.
* **Reason.** When down, `GetErrorMessage()` is populated.
* **Scoped.** Passing a specific terminal id streams just that terminal.

---

## 📚 See also

* [OnTrade](./OnTrade.md) · [gRPC Stream Management](../../All_Guides/GRPC_STREAM_MANAGEMENT.md)
* **Back to:** [5. Streaming Methods — Overview](./Streaming_Methods.Overview.md) · [MT4Account — Overview](../Overview.md)
