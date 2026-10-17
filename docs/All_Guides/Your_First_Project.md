# Your First Project in 10 Minutes (Go)

> **Hands-on Quick Start** - Create a working trading project with MetaTrader 4 and GoMT4 from scratch.

---

## Step 0: Obtain Your API Key

To connect to MetaRPC endpoints (`mt4.mrpc.pro:443`), obtain your API key:
1. Register for free at [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. Generate your API token in your dashboard at [https://mrpc.pro/my](https://mrpc.pro/my).
3. Set your token in your environment or connection config.

---

## Step 1: Create Your Project

Create a new directory for your trading bot:

```bash
mkdir my_gomt4_bot
cd my_gomt4_bot
go mod init my_gomt4_bot
```

Install the package:

```bash
go get github.com/MetaRPC/GoMT4/package
```

---

## Step 2: Write Your Trading Code

Create `main.go`:

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/MetaRPC/GoMT4/package"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	login, _ := strconv.ParseUint(os.Getenv("MT4_LOGIN"), 10, 64)
	password := os.Getenv("MT4_PASSWORD")
	server := os.Getenv("MT4_SERVER") // the broker server name as shown in MetaTrader, e.g. "TeleTrade-Demo"
	apiKey := os.Getenv("MRPC_API_KEY")
	grpcServer := os.Getenv("MRPC_GRPC_HOST")
	if grpcServer == "" {
		grpcServer = "mt4.mrpc.pro:443"
	}

	conn, err := grpc.NewClient(grpcServer, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Connect: every call carries the API key; after connecting, also the terminal id ConnectEx returns.
	timeout := uint32(120)
	reply, err := pb.NewConnectionClient(conn).ConnectEx(metadata.AppendToOutgoingContext(ctx, "apikey", apiKey),
		&pb.ConnectExRequest{User: login, Password: password, MtClusterName: server, TimeoutSeconds: &timeout})
	if err != nil {
		log.Fatal(err)
	}
	if reply.GetError() != nil {
		log.Fatal(reply.GetError().GetErrorMessage())
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "apikey", apiKey, "id", reply.GetData().GetTerminalInstanceGuid())
	defer pb.NewConnectionClient(conn).Disconnect(ctx, &pb.DisconnectRequest{})

	summary, err := pb.NewAccountHelperClient(conn).AccountSummary(ctx, &pb.AccountSummaryRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %.2f, Equity: %.2f %s\n", summary.GetData().GetAccountBalance(), summary.GetData().GetAccountEquity(), summary.GetData().GetAccountCurrency())
}
```

The program calls the generated gRPC clients of `github.com/MetaRPC/GoMT4/package` directly: every call carries your API key, and after
`ConnectEx` also the terminal id it returns.

---

## Step 3: Run the Program

```bash
export MRPC_API_KEY=your-api-key
export MT4_LOGIN=12345678
export MT4_PASSWORD=your-password
export MT4_SERVER=TeleTrade-Demo
go run .
```

```
Balance: 10000.00, Equity: 10000.00 USD
```

---

## 🚀 Next Steps

Congratulations! You have successfully established a direct gRPC connection to MetaTrader 4. Next:
- Explore **[Streaming](../Streaming/Streaming_Overview.md)** to listen to live ticks.
- Check the **[API Reference](../API%20Reference/Overview.md)** for the messages and enums of `github.com/MetaRPC/GoMT4/package`.
- Run the full example app with the `MT4Account` / `MT4Service` client from a clone of the repository: **[Beginner Run Guide](../Beginner_Run_Guide.md)**.
