/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/main.go - GoMT4 EXAMPLES LAUNCHER
 PURPOSE:
   Interactive entry point for every GoMT4 example. Loads config, connects once,
   builds the three client layers (Account -> Service -> Sugar) and dispatches
   the chosen demo from a menu.

 🎯 WHO SHOULD USE THIS:
   • Anyone running the examples for the first time
   • Seeing how the three layers are built and shared across demos

 📚 WHAT THIS DEMO COVERS (menu):
   1_lowlevel (MT4Account):  11) General   12) Trading   13) Streaming
   2_service  (MT4Service):  21) Service demo
   3_sugar    (MT4Sugar):    31) Basics    32) Trading   33) Positions
   4_strategies:             41) Orchestrators + presets
   Suites:                    1) all low-level (read-only)   0) everything   q) quit

 🔄 API LEVELS:
   LOW-LEVEL (MT4Account) -> MID-LEVEL (MT4Service) -> HIGH-LEVEL (MT4Sugar)

 ℹ️  SAFE: launcher only. The trading demos it dispatches are gated behind
     MT4_ENABLE_TRADING=1 (skipped otherwise).

 🚀 HOW TO RUN:
   cd examples && go run .                 (interactive menu)
   cd examples && echo 11 | go run .        (pick one non-interactively)
   go run . --api-key <key>                 (override the config API key)
══════════════════════════════════════════════════════════════════════════════*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	lowlevel "github.com/MetaRPC/GoMT4/examples/1_lowlevel"
	service "github.com/MetaRPC/GoMT4/examples/2_service"
	sugardemo "github.com/MetaRPC/GoMT4/examples/3_sugar"
	strategydemo "github.com/MetaRPC/GoMT4/examples/4_strategies"
	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/examples/config"
	"github.com/MetaRPC/GoMT4/mt4"
	"github.com/google/uuid"
)

func main() {
	//#region CONFIG & API KEY
	// ══════════════════════════════════════════════════════════════
	// CONFIG & API KEY
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// LoadConfig() + optional API-key override
	//   Read config/config.json, then let `go run . <api-key>` or
	//   `--api-key <key>` override the key.
	// ──────────────────────────────────────────────────────────────
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("❌ Failed to load config.json: %v", err)
	}

	// Optional API key override: `go run . <api-key>` or `--api-key <key>`.
	apiKey := cfg.ApiKey
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case strings.HasPrefix(arg, "--api-key="):
			apiKey = strings.TrimPrefix(arg, "--api-key=")
		case arg == "--api-key" && i+1 < len(os.Args):
			apiKey = os.Args[i+1]
			i++
		case !strings.HasPrefix(arg, "-") && i == 1:
			apiKey = arg
		}
	}
	if apiKey != "" {
		cfg.ApiKey = apiKey
	}
	//#endregion CONFIG & API KEY

	//#region CONNECT & BUILD LAYERS
	// ══════════════════════════════════════════════════════════════
	// CONNECT & BUILD LAYERS
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// NewMT4Account() + ConnectByServerName()
	//   Create the low-level client and connect via the broker server
	//   name; disconnect is deferred.
	// ──────────────────────────────────────────────────────────────
	account, err := mt4.NewMT4Account(uint64(cfg.Login), cfg.Password, "", uuid.Nil, cfg.ApiKey)
	if err != nil {
		log.Fatalf("❌ Failed to create MT4 account: %v", err)
	}
	defer account.Disconnect(true)

	ctx := context.Background()
	if err := account.ConnectByServerName(ctx, cfg.Server, cfg.DefaultSymbol, true, 30); err != nil {
		log.Fatalf("❌ Failed to connect to MT4 server: %v", err)
	}
	fmt.Println("✅ Connected to MT4 server")

	// ──────────────────────────────────────────────────────────────
	// Build the three layers once (Account -> Service -> Sugar)
	//   Each grouped demo picks whichever layer it shows.
	// ──────────────────────────────────────────────────────────────
	svc := mt4.NewMT4Service(account)
	sugar := mt4.NewMT4Sugar(svc).SetDefaults(777, 20, "gomt4-demo")
	symbol := cfg.DefaultSymbol

	menu(ctx, account, svc, sugar, symbol)
	//#endregion CONNECT & BUILD LAYERS

	//#region DISCONNECT
	// ══════════════════════════════════════════════════════════════
	// DISCONNECT
	// ══════════════════════════════════════════════════════════════
	fmt.Println("\n📴 Disconnecting from MT4...")
	if err := account.Disconnect(true); err != nil {
		fmt.Printf("⚠️ Disconnect notice: %v\n", err)
	} else {
		fmt.Println("✅ Disconnected successfully")
	}
	//#endregion DISCONNECT
}

//#region MENU LOOP
// ══════════════════════════════════════════════════════════════════
// MENU LOOP - print the menu and dispatch the chosen demo
// ══════════════════════════════════════════════════════════════════

// menu prints the interactive menu and runs the selected example until quit/EOF.
func menu(ctx context.Context, acc *mt4.MT4Account, svc *mt4.MT4Service, sugar *mt4.MT4Sugar, symbol string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(`
================= GoMT4 Examples =================
  1_lowlevel (MT4Account)
   11) General (account info + market)
   12) Trading & orders
   13) Streaming
  2_service (MT4Service facade)
   21) Service demo
  3_sugar (MT4Sugar)
   31) Basics (prices · math · risk)
   32) Trading (entry · manage · trail)
   33) Positions & watchers
  4_strategies
   41) Orchestrators + presets
  Suites
    1) Run all low-level (read-only)
    0) Run everything (read-only)
    q) Quit
`)
		if common.TradingEnabled {
			fmt.Println("  ⚠️  Trading demos ENABLED (MT4_ENABLE_TRADING=1)")
		}
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return // EOF (non-interactive) → exit
		}
		switch strings.TrimSpace(line) {
		case "11":
			lowlevel.General(ctx, acc, symbol)
		case "12":
			lowlevel.Trading(ctx, acc, symbol)
		case "13":
			lowlevel.Streaming(ctx, acc, symbol)
		case "21":
			service.Demo(ctx, svc, symbol)
		case "31":
			sugardemo.Basics(ctx, sugar, symbol)
		case "32":
			sugardemo.Trading(ctx, sugar, symbol)
		case "33":
			sugardemo.Positions(ctx, sugar, symbol)
		case "41":
			strategydemo.Demo(ctx, sugar, symbol)
		case "1":
			lowlevel.General(ctx, acc, symbol)
			lowlevel.Trading(ctx, acc, symbol)
			lowlevel.Streaming(ctx, acc, symbol)
		case "0":
			lowlevel.General(ctx, acc, symbol)
			lowlevel.Trading(ctx, acc, symbol)
			lowlevel.Streaming(ctx, acc, symbol)
			service.Demo(ctx, svc, symbol)
			sugardemo.Basics(ctx, sugar, symbol)
			sugardemo.Trading(ctx, sugar, symbol)
			sugardemo.Positions(ctx, sugar, symbol)
			strategydemo.Demo(ctx, sugar, symbol)
		case "q", "Q", "quit", "exit":
			return
		case "":
		default:
			fmt.Println("Unknown choice.")
		}
	}
}

//#endregion MENU LOOP
