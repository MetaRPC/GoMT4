# 🔁 HandleReconnect (GoMT4)

**Goal:** robust reconnects for **unary** and **streaming** RPCs using the helpers already present in this repo — and what to do when they give up.

> Real code refs:
>
> * Backoff & helpers: `examples/mt4/MT4Account.go` (`backoffDelay`, `waitWithCtx`, `maxRetries` — unexported)
> * Unary pattern: `examples/mt4/MT4Account.go` (`ExecuteWithReconnect`, exported)
> * Streams: `ExecuteStreamWithReconnect` (exported), used by `OnSymbolTick`, `OnTrade`, `OnOpenedOrdersProfit`, `OnOpenedOrdersTickets`

---

## ✅ 1) Principles

* Retry **only transient** errors: `codes.Unavailable`, `io.EOF` (streams), and the API errors `TERMINAL_INSTANCE_NOT_FOUND` / `TERMINAL_REGISTRY_TERMINAL_NOT_FOUND` (terminal restarting). The built-in helpers already do exactly this.
* Respect **context** (timeouts/cancel) to avoid leaks.
* Use **exponential backoff + jitter** (central constants in `MT4Account.go`).

---

## 🔹 2) Unary RPC with built-in retry

Every `MT4Account` method runs its RPC through `ExecuteWithReconnect`: try → on a transient error wait `backoffDelay(attempt)` → retry, up to `maxRetries` (10) attempts. The excerpt below is the real `AccountSummary` (inside package `mt4`):

```go
grpcCall := func(headers metadata.MD) (*pb.AccountSummaryReply, error) {
    c := metadata.NewOutgoingContext(ctx, headers)
    return a.AccountClient.AccountSummary(c, req)
}
errorSelector := func(reply *pb.AccountSummaryReply) *pb.Error {
    return reply.GetError()
}
reply, err := ExecuteWithReconnect(a, ctx, grpcCall, errorSelector)
```

`ExecuteWithReconnect` is exported, so you can use it for an RPC the client does not wrap (here `ConnectionStatus` from the generated `ConnectionClient`):

```go
reply, err := mt4.ExecuteWithReconnect(account, ctx,
    func(md metadata.MD) (*pb.ConnectionStatusReply, error) {
        return account.ConnectionClient.ConnectionStatus(metadata.NewOutgoingContext(ctx, md), &pb.ConnectionStatusRequest{})
    },
    func(r *pb.ConnectionStatusReply) *pb.Error { return r.GetError() },
)
if err != nil { return err }
fmt.Println("terminal state:", reply.GetData().GetState())
```

---

## 🔸 3) Streaming reconnect loop (structure)

The stream helpers (`OnSymbolTick`, `OnOpenedOrdersProfit`, …) encapsulate the reconnect loop in `ExecuteStreamWithReconnect`. Simplified:

```go
for {
    // (re)open the stream, retrying codes.Unavailable with backoff
    stream, err := streamInvoker(request, a.getHeaders(), ctx)
    // ...
    for {
        reply := newReply()
        err := stream.RecvMsg(reply)
        if err != nil {
            // codes.Unavailable or io.EOF → wait backoffDelay(attempt), break to reopen
            // anything else → errCh <- err; return
        }
        if apiErr := getError(reply); apiErr != nil {
            // TERMINAL_INSTANCE_NOT_FOUND / TERMINAL_REGISTRY_TERMINAL_NOT_FOUND → wait, reopen
            // other API errors → errCh <- wrapAPIError(apiErr); return
        }
        // forward data to dataCh (or stop on ctx.Done())
    }
}
```

After `maxRetries` consecutive failures the helper sends the last error to `errCh` and closes both channels.

---

## ▶️ 4) Consumer pattern with re-subscribe

When a stream ends with an error, the built-in retries are exhausted. Re-subscribe after a pause (or give up):

```go
func watchTicks(ctx context.Context, account *mt4.MT4Account, symbols []string, handle func(*pb.OnSymbolMqlTickInfo)) error {
    for {
        dataCh, errCh := account.OnSymbolTick(ctx, symbols)
    recv:
        for {
            select {
            case <-ctx.Done():
                return ctx.Err() // graceful stop
            case pkt, ok := <-dataCh:
                if !ok {
                    break recv
                }
                handle(pkt.GetSymbolTick()) // keep it fast; offload heavy work
            }
        }
        if err, ok := <-errCh; ok && err != nil {
            log.Printf("stream stopped: %v; re-subscribing in 5s", err)
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(5 * time.Second):
        }
    }
}
```

* Heavy work → send to a worker via a buffered channel.
* Don't forget `ctx.Done()` for a clean finish.

---

## 🔌 5) When the terminal itself is gone

If calls keep failing with terminal-not-found after the retries, the terminal for your account is not connected any more. Connect again:

```go
if err := account.ConnectByServerName(ctx, serverName, "EURUSD", true, 60); err != nil {
    // A failed connect (or Disconnect) closes the account's gRPC connection:
    // create a fresh account and connect that one on the next attempt.
    account, err = mt4.NewMT4Account(login, password, "", uuid.Nil, apiKey)
    if err != nil { return err }
}
```

* After `Disconnect()` / `Close()`, or a connect whose health-check failed, the `MT4Account` cannot be reused; make a new one with `NewMT4Account`.

---

## 🧭 6) Tuning backoff (central knobs)

Constants found in `examples/mt4/MT4Account.go` (unexported — change them in your clone):

```go
const (
    backoffBase = 300 * time.Millisecond
    backoffMax  = 5 * time.Second
    jitterRange = 200 * time.Millisecond
    maxRetries  = 10
)
```

* **Home Wi‑Fi / unstable** → try `backoffMax=8–10s`, `jitterRange=300–400ms`.
* **Fast VPS** → `backoffBase=150ms`, `backoffMax=3–5s`.

---

## ⚠️ Pitfalls

* **Retrying business errors** → not allowed (invalid volume, invalid stops, …); the helpers only retry transport and terminal-not-found errors.
* **Forgot to cancel the context** → goroutine leaks. Always `defer cancel()`.
* **Blocking dataCh** → the stream stalls (the channel is unbuffered). Read fast and hand work off.
* **Endless re-subscribe loops** → add a cap or alerting around your own loop, log the final error.

---

## 🔗 See also

* [Reliability](../../ReTimeouts_Reconnects_Backoff.md) — timeouts, reconnects & backoff summary.
* [StreamQuotes](../Market_Info/StreamQuotes.md), [StreamOpenedOrderProfits](../Account/StreamOpenedOrderProfits.md) — ready-made wrappers.
* [UnaryRetries](UnaryRetries.md) — extra retries around individual calls.
