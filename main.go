package main

import (
	"context"
	"fmt"
	"log"
    
	"github.com/google/uuid"
 
   accountpkg "github.com/MetaRPC/GoMT4/internal/account"
   servicepkg "github.com/MetaRPC/GoMT4/internal/service"
)

func main() {
	// === 🔧 Initialization ===

	cfg, err := LoadConfig("C:/Go/GoMT4/config.json")
	if err != nil {
		log.Fatalf("❌ Failed to load config.json: %v", err)
	}

	account, err := accountpkg.NewMT4Account(uint64(cfg.Login), cfg.Password, "", uuid.Nil)
	if err != nil {
		log.Fatalf("❌ Failed to create MT4 account: %v", err)
	}

	ctx := context.Background()
	err = account.ConnectByServerName(ctx, cfg.Server, cfg.DefaultSymbol, true, 30)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MT4 server: %v", err)
	}
	fmt.Println("✅ Connected to MT4 server")

	svc := servicepkg.NewMT4Service(account)

	if err := svc.ShowSymbolParams(ctx, cfg.DefaultSymbol); err != nil {
		log.Fatalf("❌ Failed to get symbol params: %v", err)

		// === 🚀 Step-by-step execution of methods ===

	// --- 📂 Account Info ---
	svc.ShowAccountSummary(ctx)

	// --- 📂 Order Operations ---
	svc.ShowOpenedOrders(ctx)
	svc.ShowOpenedOrderTickets(ctx)
	svc.ShowOrdersHistory(ctx)

	// ⚠️ Warning: this will open a real trade on demo/real account
	// Opens a market (or pending) order. The order type is set by the operationType field — Buy, Sell, BuyLimit, etc.
	svc.ShowOrderSendExample(ctx, cfg.DefaultSymbol)

	// ⚠️ Make sure the ticket is valid before using these:
	svc.ShowOrderModifyExample(ctx, 12345678)                //  Modify SL/TP
	svc.ShowOrderCloseExample(ctx, 12345678)                 //  Close order by ticket
	svc.ShowOrderCloseByExample(ctx, 12345678, 12345679)     //  Close with opposite order

	// This method only works with pending orders (BuyLimit, SellLimit, etc.)
	// If an active order (Buy/Sell) is passed, deletion will not work — use ShowOrderCloseExample instead.
	svc.ShowOrderDeleteExample(ctx, 12345678)

	// --- 📂 Market Info / Symbol Info ---
	svc.ShowQuote(ctx, cfg.DefaultSymbol)                                          // Get current quote by symbol (Bid, Ask, Time)
	svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})               // Get current quotes for multiple symbols
	svc.ShowQuoteHistory(ctx, cfg.DefaultSymbol)                                  // Get historical candlesticks (OHLC) for the symbol
	svc.ShowAllSymbols(ctx)                                                       // Show all available instruments in MT4
	svc.ShowSymbolParams(ctx, cfg.DefaultSymbol)                                  // Get full symbol info: Digits, Volume, Trade Mode, etc.
	svc.ShowTickValues(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})               // Get tick value, tick size, contract size for each symbol

	// --- 📂 Streaming / Subscriptions ---
	svc.StreamQuotes(ctx)                 //  Live price updates (ticks). Stops after 30 seconds in example.
	svc.StreamOpenedOrderProfits(ctx)    //  Stream open order profits. Interval can be customized.
	svc.StreamOpenedOrderTickets(ctx)    //  Stream order tickets — shows only open ones.
	svc.StreamTradeUpdates(ctx)          //  Stream trade updates — all trading activity in real-time.

	// == 🧪 Demo usage examples — uncomment to test ==
	/*
	   svc.ShowOpenedOrders(ctx)                            // Show active orders
	   svc.ShowOpenedOrderTickets(ctx)                      // List open order tickets
	   svc.ShowOrdersHistory(ctx)                           // Show order history (last 7 days)

	   svc.ShowOrderSendExample(ctx, cfg.DefaultSymbol)     // 🔄 Open Buy order at market

	   svc.ShowOrderModifyExample(ctx, 12345678)            // ✏️ Modify SL/TP — insert real ticket
	   svc.ShowOrderCloseExample(ctx, 12345678)             // ❌ Close order by ticket
	   svc.ShowOrderCloseByExample(ctx, 12345678, 12345679) // ♻️ Close with opposite order
	   svc.ShowOrderDeleteExample(ctx, 12345678)            // 🗑 Delete pending order
	*/


}
}
