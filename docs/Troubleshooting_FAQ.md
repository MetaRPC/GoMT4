# Troubleshooting & FAQ (GoMT4)

Short, practical answers. Each item points to real code paths where relevant.

---

## 1) “No quotes / symbol not found (EURUSD)”

**Symptoms:** `symbol not found`, empty quotes, or RPC returns OK but payload is empty.

**Likely causes:**

* Broker uses a **suffix** (e.g., `EURUSD.m`, `EURUSD.pro`).
* Symbol is **hidden** in MT4 *Market Watch*.

**Fix:**

* Open MT4 → *Market Watch* → *Show All* → note the exact symbol name and put it into `examples/config/config.json` → `DefaultSymbol`.
* If your code has an `EnsureSymbolVisible(symbol)` helper, call it before requests; otherwise, add one. (Typical place: `examples/mt4/MT4Account.go`.)

**Tip:** print `Digits`, `Point`, `LotStep` in logs when selecting a symbol—this quickly reveals mismatches.

---

## 2) “Timeout / context deadline exceeded” on simple reads

**Symptoms:** `context deadline exceeded` on read‑only calls (quotes, account summary).

**Causes:**

* MT4 not fully connected to broker or just launched.
* Network latency spikes.

**Fix:**

* Start MT4 **manually once** and wait for *connected* state.
* Use a short per‑call timeout (2–5s) and retry only **transport** errors (see §3 in `Docs/reliability (en)`).
* Code reference: a 3s health‑check with `AccountSummary` in `examples/mt4/MT4Account.go`.

---

## 3) “Max retries reached” / frequent reconnects in streams

**Symptoms:** stream stops with an error after several reconnect attempts; logs mention `io.EOF` or `codes.Unavailable`.

**Causes:**

* Unstable connection; too aggressive backoff; context canceled early.

**Fix:**

* Increase `backoffMax` or `backoffBase` in retry settings (see `examples/mt4/MT4Account.go`).
* Ensure the parent `ctx` is not canceled by your app prematurely.
* Consumer pattern: always select on `dataCh`, `errCh`, and `<-ctx.Done()`; exit cleanly when channels close.

---

## 4) “Can’t connect to gRPC: connection refused”

**Symptoms:** client can’t dial the server; `connection refused` or hangs.

**Checklist:**

* Is server running? (`go run ./examples/main.go`)
* Is it listening on the expected address? (`127.0.0.1:50051` by default)
* `netstat -ano | findstr LISTENING | findstr :50051` — do you see a listener?
* Windows Firewall: if you bind to `0.0.0.0` or external IP, allow the port:

  ```powershell
  New-NetFirewallRule -DisplayName "GoMT4 gRPC" -Direction Inbound -Protocol TCP -LocalPort 50051 -Action Allow
  ```

---

## 5) “Invalid volume / invalid price” when sending orders

**Symptoms:** `invalid volume`, `invalid price`, or broker rejects the order.

**Causes:**

* Volume not aligned to `LotStep`.
* Price/SL/TP not aligned to `Digits`/`Point` / too close to the market.

**Fix:**

* Query symbol params first, then round:

  * volume → to `LotStep` (clamp to `MinLot`…`MaxLot`).
  * prices → to `Digits` using `Point` (or 10^Digits helper).
* In logs, print the calculated price, `Digits`, `Point`, `LotStep`.
* Put the rounding into a small helper in `examples/mt4/MT4Account.go` (or your order module) and reuse it.

---

## 6) “Quotes freeze after a while”

**Symptoms:** stream was active but stopped emitting data; no errors printed.

**Causes:**

* Consumer stopped reading from `dataCh` (blocked).
* `ctx` canceled elsewhere.

**Fix:**

* Ensure your consumer loop **never blocks** (use a bounded queue or backpressure strategy).
* Always monitor `errCh` and `<-ctx.Done()>` and exit cleanly. The helper closes channels on terminal errors.

---

## 7) “module … not found / checksum mismatch” (Go modules)

**Symptoms:** during `go mod tidy` or build.

**Fix:**

* Make sure the pb import path matches the module path (no `.git` suffix):

  ```go
  import pb "git.mtapi.io/root/mrpc-proto/mt4/libraries/go"
  ```
* Update or pin the module:

  ```powershell
  go get -u git.mtapi.io/root/mrpc-proto/mt4/libraries/go@latest
  go mod tidy
  ```
* If CI needs offline builds: `go mod vendor` and build with `-mod=vendor`.

---

## 8) “TLS handshake / certificate” issues

**Symptoms:** errors around TLS when using secure channels.

**Fix:**

* For local dev, prefer plaintext on `127.0.0.1` (no TLS) unless you explicitly configured creds.
* If using TLS, verify you pass `grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ /* … */ }))` consistently on both sides and that CN/SAN match the host you dial.

---

## 9) “No history / partial history returned”

**Symptoms:** history calls return fewer records than expected.

**Causes:**

* MT4 hasn’t downloaded that range yet.
* Too big range in a single request.

**Fix:**

* Open the symbol chart in MT4 once to let it preload more history.
* Use **paging/batching** for long ranges (day‑by‑day or month‑by‑month) instead of one massive call.

---

## 10) “High CPU / goroutine leak”

**Symptoms:** CPU climbs or goroutines accumulate.

**Causes:**

* Missing `defer cancel()`; streams not canceled; consumer loop busy‑spins.

**Fix:**

* After every `WithTimeout/WithCancel`, add `defer cancel()`.
* On shutdown: cancel parent `ctx` first, then wait for goroutines to exit, then `Disconnect()`.

---

## 11) Quick reference (commands)

* Verify port listener:

  ```powershell
  netstat -ano | findstr LISTENING | findstr :50051
  Test-NetConnection -ComputerName 127.0.0.1 -Port 50051
  ```
* Refresh deps:

  ```powershell
  go mod tidy
  go get -u git.mtapi.io/root/mrpc-proto/mt4/libraries/go@latest
  ```
* Vendor for offline builds:

  ```powershell
  go mod vendor
  go build -mod=vendor ./...
  ```

---

## 12) Where in code

* Retry/backoff & helpers: `examples/mt4/MT4Account.go`
* Streaming wrappers (ticks/orders/history): `examples/mt4/MT4Account.go`
* Entrypoint & cleanup: `examples/main.go`
* Config shape: `examples/config/config.json`
