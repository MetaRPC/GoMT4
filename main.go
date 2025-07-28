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
	svc.ShowOrderModifyExample(ctx, 12345678)
	 svc.ShowOrderCloseExample(ctx, 12345678)
	  svc.ShowOrderCloseByExample(ctx, 12345678, 12345679)

	//This method only works with pending orders (BuyLimit, SellLimit, etc.).
   //If an active order (Buy/Sell) is passed, deletion will not work — use ShowOrderCloseExample() instead.
	svc.ShowOrderDeleteExample(ctx, 12345678)

	// --- 📂 Market Info / Symbol Info ---
    svc.ShowQuote(ctx, cfg.DefaultSymbol) 
     svc.ShowQuotesMany(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"})         //This method requests the current quotes for several symbols at once and displays them.
      svc.ShowQuoteHistory(ctx, cfg.DefaultSymbol)                          // Get historical candlesticks by symbol (Open, High, Low, Close, Time) for the specified period and timeframe (for example, H1).
       svc.ShowAllSymbols(ctx)                                             // Show all available symbols (instruments) that are available in the MT4 terminal: EURUSD, GBPUSD, XAUUSD, etc.
        svc.ShowSymbolParams(ctx, cfg.DefaultSymbol)                      // Get advanced parameters by symbol: Digits, minimum and maximum volume, currencies, trading mode, execution type, floating spread, etc.
         svc.ShowTickValues(ctx, []string{"EURUSD", "GBPUSD", "USDJPY"}) // Show for each character: cost per tick (TickValue); tick size (tickSize); the size of the contract (ContractSize).

    // --- 📂 Streaming / Subscriptions ---
	svc.StreamQuotes(ctx)                 // This is a live subscription, and it won't end until you manually or timer it off (in the example, after 30 seconds).
     svc.StreamOpenedOrderProfits(ctx)   // The interval of 1000 ms can be reduced for a higher frequency. If there are no orders, the stream will be empty, but the channel is still active.
      svc.StreamOpenedOrderTickets(ctx) // Live subscription to the list of open order tickets. Updated when: the order is open — a ticket appears; the order is closed, and the ticket disappears.
       svc.StreamTradeUpdates(ctx)     // Live subscription to any trading events. 



	// === 🧪 Demo usage examples — uncomment to test ===

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