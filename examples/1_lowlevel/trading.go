/*══════════════════════════════════════════════════════════════════════════════
 FILE: examples/1_lowlevel/trading.go - LOW-LEVEL MT4 TRADING & ORDERS DEMO
 PURPOSE:
   Exercises the MT4Account order methods at the RAW layer: you pass the protobuf
   operation type and optional pointer fields yourself (MT4Sugar wraps all of this
   in BuyMarket/ModifySLTPByPips/…). Reads always run; the write path is gated.

 🎯 WHO SHOULD USE THIS:
   • Developers who need full control over order construction
   • Building a custom trading layer on top of MT4Account
   • Learning the raw OrderSend / OrderModify / OrderClose flow

 📚 WHAT THIS DEMO COVERS (2 sections):

   1. READS (always safe)
      • OrdersTotal() / OpenedOrders() / OpenedOrdersTickets()
      • OrderSelect() / OrdersHistory()

   2. WRITES (gated by MT4_ENABLE_TRADING=1)
      • OrderSend()               - open a market BUY (0.01 lot)
      • OrderModify()             - change price / SL / TP / expiry
      • OrderClose()              - close a market order / delete a pending

 🔄 API LEVELS:
   LOW-LEVEL  (MT4Account) - THIS FILE: raw OrderSend / OrderModify / OrderClose
   MID-LEVEL  (MT4Service) - typed wrappers, cleaner errors
   HIGH-LEVEL (MT4Sugar)   - one-liner BuyMarket / ModifySLTPByPips / Close

 ⚠️  WARNING: sends real orders on a demo account (gated by MT4_ENABLE_TRADING=1).

 🚀 HOW TO RUN:
   cd examples && go run .        (then choose 12 from the menu)
══════════════════════════════════════════════════════════════════════════════*/

package lowlevel

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// Trading exercises the MT4Account order methods. This is the RAW layer: you pass
// the protobuf operation type and optional pointer fields yourself (MT4Sugar wraps
// all of this in BuyMarket/ModifySLTPByPips/…). Reads always run; the write path
// is gated behind MT4_ENABLE_TRADING=1.
func Trading(ctx context.Context, acc *mt4.MT4Account, symbol string) {
	fmt.Println("\n=== 📂 MT4Account · Trading & Orders ===")

	//#region 1. READS (safe)
	// ══════════════════════════════════════════════════════════════
	// 1. READS (always safe)
	// ══════════════════════════════════════════════════════════════

	// ──────────────────────────────────────────────────────────────
	// OrdersTotal()
	//   Count of currently open orders.
	// ──────────────────────────────────────────────────────────────
	if n, err := acc.OrdersTotal(ctx); err == nil {
		fmt.Printf("  OrdersTotal: %d\n", n) // count of open orders
	}

	// ──────────────────────────────────────────────────────────────
	// OpenedOrders() / OpenedOrdersTickets() / OrderSelect()
	//   The whole open book, just its tickets, and one order's detail.
	// ──────────────────────────────────────────────────────────────
	// OpenedOrders returns the whole open book (market + pending) as a wrapper;
	// at this layer you unwrap it yourself with GetOrderInfos().
	var firstTicket int32
	if d, err := acc.OpenedOrders(ctx); err != nil {
		fmt.Println("  OpenedOrders error:", err)
	} else {
		fmt.Printf("  OpenedOrders: %d\n", len(d.GetOrderInfos()))
		for _, o := range d.GetOrderInfos() {
			if firstTicket == 0 {
				firstTicket = o.GetTicket() // remember one ticket for OrderSelect below
			}
			fmt.Printf("    #%d %s %.2f lot  P/L %+.2f\n", o.GetTicket(), o.GetSymbol(), o.GetLots(), o.GetProfit())
		}
	}
	if d, err := acc.OpenedOrdersTickets(ctx); err == nil {
		fmt.Printf("  OpenedOrdersTickets: %v\n", d.GetTickets()) // just the ticket numbers
	}
	if firstTicket != 0 {
		// OrderSelect fetches the full detail of one order by ticket.
		if o, err := acc.OrderSelect(ctx, firstTicket); err == nil {
			fmt.Printf("  OrderSelect #%d: %s %.2f lot\n", o.GetTicket(), o.GetSymbol(), o.GetLots())
		}
	}

	// ──────────────────────────────────────────────────────────────
	// OrdersHistory()
	//   Closed-order history over a window (last 30 days).
	// ──────────────────────────────────────────────────────────────
	// Closed-order history over a window. The sort enum is passed as a value;
	// from/to are *time.Time, page/items are nil to use defaults.
	now := time.Now()
	from := now.Add(-30 * 24 * time.Hour)
	if h, err := acc.OrdersHistory(ctx, pb.EnumOrderHistorySortType(0), &from, &now, nil, nil); err == nil {
		fmt.Printf("  OrdersHistory (30d): %d orders\n", len(h.GetOrdersInfo()))
	}
	//#endregion 1. READS (safe)

	//#region 2. WRITES (gated by MT4_ENABLE_TRADING=1)
	// ══════════════════════════════════════════════════════════════
	// 2. WRITES (dangerous — opt in via MT4_ENABLE_TRADING=1)
	// ══════════════════════════════════════════════════════════════
	if !common.TradingEnabled {
		fmt.Println("  (trading demos skipped — set MT4_ENABLE_TRADING=1 to run)")
		return
	}
	fmt.Println("  ⚠️  Trading: OrderSend 0.01 BUY → modify → close")

	// ──────────────────────────────────────────────────────────────
	// OrderSend()
	//   The one write that creates orders. Explicit signature: symbol,
	//   operation type, volume, then optional pointers (price, slippage,
	//   stoploss, takeprofit, comment, magic, expiration) — nil = omit.
	//   Here: a plain market BUY of 0.01 lot with everything else defaulted.
	// ──────────────────────────────────────────────────────────────
	r, err := acc.OrderSend(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUY, 0.01, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		// Trade rejections (market closed, no money, invalid stops…) arrive as an
		// error here — the low-level layer does NOT auto-retry trade calls, because
		// a timed-out send may actually have filled.
		fmt.Println("  OrderSend error:", err)
		return
	}
	t := r.GetTicket()
	fmt.Printf("  OrderSend ok: #%d @ %.5f\n", t, r.GetPrice())

	// ──────────────────────────────────────────────────────────────
	// OrderModify()
	//   Changes price / SL / TP / expiry (nil pointers = leave unchanged).
	//   Passing all nil is a no-op here, shown only for the call shape.
	// ──────────────────────────────────────────────────────────────
	if _, err := acc.OrderModify(ctx, t, nil, nil, nil, nil); err != nil {
		fmt.Println("  OrderModify error:", err)
	}

	// ──────────────────────────────────────────────────────────────
	// OrderClose()
	//   Closes a market order (lots nil = full) or deletes a pending.
	// ──────────────────────────────────────────────────────────────
	if _, err := acc.OrderClose(ctx, t, nil, nil, nil); err != nil {
		fmt.Println("  OrderClose error:", err)
	} else {
		fmt.Printf("  OrderClose #%d ok\n", t)
	}
	//#endregion 2. WRITES (gated by MT4_ENABLE_TRADING=1)
}
