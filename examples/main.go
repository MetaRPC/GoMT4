package main

import (
	"context"
	"fmt"
	"github.com/MetaRPC/GoMT4/config"
	"github.com/MetaRPC/GoMT4/mt4"
	"github.com/google/uuid"
	"log"
)

func main() {
	// Loading the config
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("❌ Failed to load config.json: %v", err)
	}

	// Creating an account
	account, err := mt4.NewMT4Account(uint64(cfg.Login), cfg.Password, "", uuid.Nil)
	if err != nil {
		log.Fatalf("❌ Failed to create MT4 account: %v", err)
	}
	defer account.Disconnect() // ensure cleanup on exit

	// Connecting to the server
	ctx := context.Background()
	err = account.ConnectByServerName(ctx, cfg.Server, cfg.DefaultSymbol, true, 30)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MT4 server: %v", err)
	}
	fmt.Println("✅ Connected to MT4 server")

	// Creating a service and calling the method
	svc := mt4.NewMT4Service(account)

	if err := svc.ShowSymbolParams(ctx, cfg.DefaultSymbol); err != nil {
		log.Fatalf("❌ Failed to get symbol params: %v", err)
	}

	//=======================================
	//===-------- Feature toggles --------===
	//=======================================
	enableTradingExamples := false // ⚠️ Real trading operations (OrderSend/Modify/Close...)
	enableStreams := true          // Streaming quotes/profit/tickets
	enableHistoryStreams := true   // Our new story streams (pages/chunks)

	// === 🚀 Step-by-step execution of methods ===

	// --- 📂 Account Info ---
	svc.ShowAccountSummary(ctx)

	// --- 📂 Order Operations (safe reading) ---
	svc.ShowOpenedOrders(ctx)
	svc.ShowOpenedOrderTickets(ctx)
	svc.ShowOrdersHistory(ctx)

	// --- ⚠️ Trading (DANGEROUS) ---
	if enableTradingExamples {
		// ATTENTION: opens a real deal (demo/real depending on login)
		svc.ShowOrderSendExample(ctx, cfg.DefaultSymbol)

		// These require valid tickets:
		svc.ShowOrderModifyExample(ctx, 12345678)            // Modify SL/TP
		svc.ShowOrderCloseExample(ctx, 12345678)             // Close by ticket
		svc.ShowOrderCloseByExample(ctx, 12345678, 12345679) // Close with opposite

		// Removes ONLY deposits (use Close for Buy/Sell)
		svc.ShowOrderDeleteExample(ctx, 12345678)
	}

	// --- 📂 Market Info / Symbol Info ---
	svc.ShowQuote(ctx, cfg.DefaultSymbol)
	svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})
	svc.ShowQuoteHistory(ctx, cfg.DefaultSymbol)
	svc.ShowAllSymbols(ctx)
	svc.ShowSymbolParams(ctx, cfg.DefaultSymbol)
	svc.ShowTickValues(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})

	// --- 📂 Streaming / Subscriptions ---
	if enableStreams {
		svc.StreamQuotes(ctx)             // live ticks (in the example, it will stop by itself)
		svc.StreamOpenedOrderProfits(ctx) // stream profits on open orders
		svc.StreamOpenedOrderTickets(ctx) // ticket streaming
		svc.StreamTradeUpdates(ctx)       // stream trade events
	}

	// --- 🧾 improved story streams (pages/chunks) ---
	if enableHistoryStreams {
		svc.StreamOrdersHistoryExample(ctx)                   // page-by-page history (30 days)
		svc.StreamQuoteHistoryExample(ctx, cfg.DefaultSymbol) // candle chunks (90 days)
	}
}
