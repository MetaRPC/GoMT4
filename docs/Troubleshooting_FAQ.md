# 🛠️ Troubleshooting & FAQ (GoMT4)

Short, practical answers. Each item points to real code paths where relevant.

---

## ❓ 1) “No quotes / symbol not found (EURUSD)”

**Symptoms:** `symbol not found`, empty quotes, or RPC returns OK but payload is empty.

**Likely causes:**

* Broker uses a **suffix** (e.g., `EURUSD.m`, `EURUSD.pro`).
* Symbol is **hidden** in MT4 *Market Watch*.

**Fix:**

* Open MT4 → *Market Watch* → *Show All* → note exact symbol name → put into `examples/config/config.json` → `DefaultSymbol`.
* If available, call `EnsureSymbolVisible(symbol)` before requests. Otherwise, add one (see `examples/mt4/MT4Account.go`).

💡 Tip: print `Digits`, `Point`, `LotStep` in logs when selecting a symbol.

---

## ⏱️ 2) “Timeout / context deadline exceeded” on simple reads

**Symptoms:** `context deadline exceeded` on read-only calls (quotes, account summary).

**Causes:**

* MT4 not fully connected or just launched.
* Network latency spikes.

**Fix:**

* Start MT4 manually once and wait until *connected*.
* Use per-call timeout (2–5s) and retry only transport errors.
* Reference: 3s health-check with `AccountSummary` in `examples/mt4/MT4Account.go`.

---

## 🔁 3) “Max retries reached” / frequent reconnects in streams

**Symptoms:** stream stops with error after reconnect attempts; logs mention `io.EOF` or `codes.Unavailable`.

**Causes:**

* Unstable connection; too aggressive backoff; context canceled early.

**Fix:**

* Increase `backoffMax` or `backoffBase` (see `examples/mt4/MT4Account.go`).
* Ensure parent `ctx` is not canceled prematurely.
* Consumer loop: always select on `dataCh`, `errCh`, `<-ctx.Done()`.

---

## 🚫 4) “Can’t connect to gRPC: connection refused”

**Checklist:**

* Server running? (`go run ./examples/main.go`)
* Listening on correct address? (`127.0.0.1:50051` default)
* Check listener: `netstat -ano | findstr LISTENING | findstr :50051`
* Windows Firewall: allow port if bound externally:

  ```powershell
  New-NetFirewallRule -DisplayName "GoMT4 gRPC" -Direction Inbound -Protocol TCP -LocalPort 50051 -Action Allow
  ```

---

## 📊 5) “Invalid volume / invalid price” when sending orders

**Symptoms:** `invalid volume`, `invalid price`, or broker rejects order.

**Causes:**

* Volume misaligned with `LotStep`.
* Price/SL/TP not aligned to `Digits`/`Point` or too close to market.

**Fix:**

* Query symbol params first, then round:

  * volume → `LotStep` (clamp to `MinLot`…`MaxLot`).
  * prices → round to `Digits` using `Point`.
* Print logs: price, `Digits`, `Point`, `LotStep`.
* Create helper in `examples/mt4/MT4Account.go` for rounding.

---

## 💤 6) “Quotes freeze after a while”

**Symptoms:** stream stops emitting data; no errors.

**Causes:**

* Consumer stopped reading from `dataCh`.
* `ctx` canceled elsewhere.

**Fix:**

* Ensure consumer loop never blocks (use bounded queue/backpressure).
* Always monitor `errCh` and `<-ctx.Done()>`.

---

## 📦 7) “module … not found / checksum mismatch” (Go modules)

**Symptoms:** during `go mod tidy` or build.

**Fix:**

* Ensure pb import path matches module path:

  ```go
  import pb "git.mtapi.io/root/mrpc-proto/mt4/libraries/go"
  ```
* Update/pin module:

  ```powershell
  go get -u git.mtapi.io/root/mrpc-proto/mt4/libraries/go@latest
  go mod tidy
  ```
* For offline CI: `go mod vendor` + build with `-mod=vendor`.

---

## 🔐 8) “TLS handshake / certificate” issues

**Fix:**

* For local dev, prefer plaintext on `127.0.0.1`.
* If TLS, use consistent creds: `grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ ... }))`.
* Ensure CN/SAN matches host.

---

## 📉 9) “No history / partial history returned”

**Symptoms:** history calls return fewer records.

**Causes:**

* MT4 hasn’t downloaded range yet.
* Request range too large.

**Fix:**

* Open symbol chart in MT4 to preload.
* Use paging (day-by-day / month-by-month).

---

## 🔥 10) “High CPU / goroutine leak”

**Symptoms:** CPU climbs; goroutines accumulate.

**Causes:**

* Missing `defer cancel()`.
* Streams not canceled.
* Consumer loop busy-spins.

**Fix:**

* Always `defer cancel()`.
* On shutdown: cancel parent `ctx`, wait for goroutines, then `Disconnect()`.

---

## 📝 11) Quick reference (commands)

* Verify listener:

  ```powershell
  netstat -ano | findstr LISTENING | findstr :50051
  Test-NetConnection -ComputerName 127.0.0.1 -Port 50051
  ```
* Refresh deps:

  ```powershell
  go mod tidy
  go get -u git.mtapi.io/root/mrpc-proto/mt4/libraries/go@latest
  ```
* Vendor offline builds:

  ```powershell
  go mod vendor
  go build -mod=vendor ./...
  ```

---

## 📂 12) Where in code

* Retry/backoff & helpers → `examples/mt4/MT4Account.go`
* Streaming wrappers (ticks/orders/history) → `examples/mt4/MT4Account.go`
* Entrypoint & cleanup → `examples/main.go`
* Config shape → `examples/config/config.json`
