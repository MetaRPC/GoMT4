# 🛠️ Troubleshooting & FAQ (GoMT4)

Short, practical answers. Each item points to real code paths where relevant.

---

## ❓ 1) “No quotes / symbol not found (EURUSD)”

**Symptoms:** `symbol not found`, empty quotes, or RPC returns OK but payload is empty.

**Likely causes:**

* Broker uses a **suffix** (e.g., `EURUSD.m`, `EURUSD.pro`).
* Symbol is **hidden** in the terminal's *Market Watch*.

**Fix:**

* List the exact names with `svc.ShowAllSymbols(ctx)` (or `account.Symbols(ctx)`) → put the right one into `examples/config/config.json` → `DefaultSymbol`.
* To show a hidden symbol, call the generated `SymbolSelect` RPC (the client has no wrapper for it) — see [Ensure Symbol Visible](Cookbook/Utils_Helpers/EnsureSymbolVisible.md).

💡 Tip: print `Digits`, `Point`, `VolumeStep` (from `account.SymbolParams`) in logs when selecting a symbol.

---

## ⏱️ 2) “Timeout / context deadline exceeded” on simple reads

**Symptoms:** `context deadline exceeded` on read-only calls (quotes, account summary).

**Causes:**

* The terminal for the account is still starting or connecting to the broker (first connect, or after a restart on the MetaRPC side).
* Network latency spikes.

**Fix:**

* Connect first (`ConnectByServerName` / `ConnectByHostPort`) and only call methods after it returns without error.
* Pass a `ctx` with a longer deadline if the default per-call timeout (3s for reads) is too short for you.
* Reference: 3s health-check with `AccountSummary` in `examples/mt4/MT4Account.go`.

---

## 🔁 3) “exceeded retries” / frequent reconnects in streams

**Symptoms:** stream stops with `exceeded retries ...` after reconnect attempts; logs mention `io.EOF`, `codes.Unavailable` or `TERMINAL_INSTANCE_NOT_FOUND`.

**Causes:**

* Unstable connection; too aggressive backoff; context canceled early.
* The terminal is not connected (any more) under the id the client sends.

**Fix:**

* Increase `backoffMax` or `backoffBase` (see `examples/mt4/MT4Account.go`).
* Ensure parent `ctx` is not canceled prematurely.
* Consumer loop: always select on `dataCh`, `errCh`, `<-ctx.Done()`.
* If the terminal is gone, connect again (a new `MT4Account` + `ConnectByServerName`).

---

## 🚫 4) “Can’t connect to gRPC” / “unavailable” / “API_KEY_INVALID”

**Checklist:**

* The client always talks to the MetaRPC cloud API at `mt4.mrpc.pro:443` (gRPC over TLS); there is no local server to start and no port to open.
* Network: outbound HTTPS (443) to `mt4.mrpc.pro` must be allowed (proxy/firewall):

  ```powershell
  Test-NetConnection -ComputerName mt4.mrpc.pro -Port 443
  ```
* API key: `API_KEY_INVALID` means the key is not known — copy it again from [https://mrpc.pro/my](https://mrpc.pro/my). Without a key the calls run as a trial (`TRIAL`).
* Credentials: wrong login/password/server name make `ConnectByServerName` return an error.

---

## 📊 5) “Invalid volume / invalid price” when sending orders

**Symptoms:** `invalid volume`, `invalid price`, or broker rejects order.

**Causes:**

* Volume misaligned with `VolumeStep`.
* Price/SL/TP not aligned to `Digits`/`Point` or too close to market (`TradeStopsLevel`).

**Fix:**

* Query symbol params first (`account.SymbolParams(ctx, symbol)`), then round:

  * volume → `VolumeStep` (clamp to `VolumeMin`…`VolumeMax`).
  * prices → round to `Digits`.
* Print logs: price, `Digits`, `Point`, `VolumeStep`.
* See [Round Volume/Price](Cookbook/Utils_Helpers/RoundVolumePrice.md) for helpers.

---

## 💤 6) “Quotes freeze after a while”

**Symptoms:** stream stops emitting data; no errors.

**Causes:**

* Consumer stopped reading from `dataCh`.
* `ctx` canceled elsewhere.
* The market is closed (no new ticks arrive).

**Fix:**

* Ensure consumer loop never blocks (use bounded queue/backpressure).
* Always monitor `errCh` and `<-ctx.Done()`.

---

## 📦 7) “module … not found / checksum mismatch” (Go modules)

**Symptoms:** during `go mod tidy` or build.

**Fix:**

* Ensure pb import path matches module path:

  ```go
  import pb "github.com/MetaRPC/GoMT4/package"
  ```
* The client (`github.com/MetaRPC/GoMT4/examples/mt4`) is not fetched by `go get`; use it from a clone (run inside `examples/`, or add a `replace` directive in your module — see [Setup](setup.md)).
* Update/pin module:

  ```powershell
  go get -u github.com/MetaRPC/GoMT4/package@latest
  go mod tidy
  ```
* For offline CI: `go mod vendor` + build with `-mod=vendor`.

---

## 🔐 8) “TLS handshake / certificate” issues

**Fix:**

* The connection is always TLS; `NewMT4Account` uses the system certificate store (`credentials.NewTLS(&tls.Config{})`).
* A corporate proxy that rewrites TLS needs its root certificate in the system store.
* Check the system clock (certificate validity).

---

## 📉 9) “No history / partial history returned”

**Symptoms:** history calls return fewer records.

**Causes:**

* The terminal hasn't downloaded the range from the broker yet.
* Request range too large.

**Fix:**

* Repeat the request after a short pause.
* Use paging: `OrdersHistoryStream` (pages) and `QuoteHistoryStream` (time chunks).

---

## 🔥 10) “High CPU / goroutine leak”

**Symptoms:** CPU climbs; goroutines accumulate.

**Causes:**

* Missing `defer cancel()`.
* Streams not canceled.
* Consumer loop busy-spins (e.g. keeps reading from a closed channel without checking `ok`).

**Fix:**

* Always `defer cancel()`.
* On shutdown: cancel parent `ctx`, wait for goroutines, then `Disconnect()`.

---

## 📝 11) Quick reference (commands)

* Check reachability of the API:

  ```powershell
  Test-NetConnection -ComputerName mt4.mrpc.pro -Port 443
  ```
* Refresh deps:

  ```powershell
  go mod tidy
  go get -u github.com/MetaRPC/GoMT4/package@latest
  ```
* Vendor offline builds (own module):

  ```powershell
  go mod vendor
  go build -mod=vendor ./...
  ```

---

## 📂 12) Where in code

* Retry/backoff & helpers → `examples/mt4/MT4Account.go`
* Streaming wrappers (ticks/orders/history) → `examples/mt4/MT4Account.go`
* Entrypoint & cleanup → `examples/main.go`
* Config shape → `examples/config/config.go`, `examples/config/config.json`
