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

	// --- Reads (safe) ---
	if n, err := acc.OrdersTotal(ctx); err == nil {
		fmt.Printf("  OrdersTotal: %d\n", n) // count of open orders
	}

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

	// Closed-order history over a window. The sort enum is passed as a value;
	// from/to are *time.Time, page/items are nil to use defaults.
	now := time.Now()
	from := now.Add(-30 * 24 * time.Hour)
	if h, err := acc.OrdersHistory(ctx, pb.EnumOrderHistorySortType(0), &from, &now, nil, nil); err == nil {
		fmt.Printf("  OrdersHistory (30d): %d orders\n", len(h.GetOrdersInfo()))
	}

	// --- Writes (dangerous — opt in) ---
	if !common.TradingEnabled {
		fmt.Println("  (trading demos skipped — set MT4_ENABLE_TRADING=1 to run)")
		return
	}
	fmt.Println("  ⚠️  Trading: OrderSend 0.01 BUY → modify → close")

	// OrderSend is the one write that creates orders. The signature is explicit:
	// symbol, operation type, volume, then optional pointers (price, slippage,
	// stoploss, takeprofit, comment, magic, expiration) — nil means "omit".
	// Here: a plain market BUY of 0.01 lot with everything else defaulted.
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

	// OrderModify changes price/SL/TP/expiry (nil pointers = leave unchanged).
	// Passing all nil is a no-op here, shown only for the call shape.
	if _, err := acc.OrderModify(ctx, t, nil, nil, nil, nil); err != nil {
		fmt.Println("  OrderModify error:", err)
	}

	// OrderClose closes a market order (lots nil = full) or deletes a pending.
	if _, err := acc.OrderClose(ctx, t, nil, nil, nil); err != nil {
		fmt.Println("  OrderClose error:", err)
	} else {
		fmt.Printf("  OrderClose #%d ok\n", t)
	}
}
