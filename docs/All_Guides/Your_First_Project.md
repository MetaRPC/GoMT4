# Your First Project in 10 Minutes

> **Hands-on quick start** — connect to MetaTrader 4 and place a demo order with GoMT4, from scratch.

---

## Step 0 — API key

Connecting to `mt4.mrpc.pro:443` needs a key:

1. Register (free): [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. Generate a token: [https://mrpc.pro/my](https://mrpc.pro/my).
3. Export it: `set MRPC_API_KEY=...` (Windows) / `export MRPC_API_KEY=...` (bash).

---

## Step 1 — Get the SDK

```bash
git clone https://github.com/MetaRPC/GoMT4.git
cd GoMT4
```

It's a Go workspace (`go.work`) of three modules — the generated `package`, the `mt4` SDK, and the runnable `examples`. The fastest path is to run an example.

---

## Step 2 — The code

The three-layer client. `MT4Sugar` gives you the shortest path to a trade:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MetaRPC/GoMT4/mt4"
	"github.com/google/uuid"
)

func main() {
	// "" = default endpoint mt4.mrpc.pro:443; uuid.Nil = deterministic id from login+password
	acc, err := mt4.NewMT4Account(12_345_678, "password", "", uuid.Nil, "your-api-key")
	if err != nil {
		log.Fatal(err)
	}
	defer acc.Disconnect(true)

	ctx := context.Background()
	// login done at construction; connect the terminal by broker server name
	if err := acc.ConnectByServerName(ctx, "MetaQuotes-Demo", "EURUSD", true, 30); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MT4!")

	// magic, slippage (points), comment stamped on order-entry helpers
	sugar := mt4.NewMT4Sugar(mt4.NewMT4Service(acc)).SetDefaults(1001, 20, "my-first-bot")

	a, _ := sugar.GetAccountInfo(ctx)
	fmt.Printf("Balance: %.2f %s\n", a.GetAccountBalance(), a.GetAccountCurrency())

	// make sure the symbol is live, then buy 0.01 lot at market
	if _, err := sugar.EnsureSymbol(ctx, "EURUSD", 10*time.Second); err != nil {
		log.Fatal(err)
	}
	r, err := sugar.BuyMarket(ctx, "EURUSD", 0.01, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order sent! Ticket #%d @ %.5f\n", r.GetTicket(), r.GetPrice())

	// 20-pip stop, 40-pip target (direction-aware, by pips from entry)
	sl, tp := 20.0, 40.0
	_ = sugar.ModifySLTPByPips(ctx, r.GetTicket(), &sl, &tp)
}
```

> **Demo accounts only** while you learn — `BuyMarket` sends a real order on whatever account you connect.

---

## Step 3 — Run the bundled examples

The repo ships the same flow (and much more) as a runnable app:

```bash
cd examples
go run .            # interactive menu of grouped examples
```

Pick `31` for the MT4Sugar basics, `11` for low-level account info, and so on. Credentials come from `examples/config/config.json` (or the `MRPC_API_KEY` env var) — see the [config loader](PROJECT_MAP.md). Trading demos are off by default; enable them (demo only!) with `MT4_ENABLE_TRADING=1`.

---

## 🚀 Next steps

* Stream live prices → **[gRPC Streaming](GRPC_STREAM_MANAGEMENT.md)**.
* Every terminal method → **[MT4Account API](../MT4Account/Overview.md)**.
* Risk sizing, trailing, one-liners → **[MT4Sugar API](../MT4Sugar/Overview.md)**.
* Explore the types → **[Protobuf Inspector](PROTOBUF_INSPECTOR_GUIDE.md)**.
