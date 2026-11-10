package lowlevel

import (
	"context"
	"fmt"
	"time"

	"github.com/MetaRPC/GoMT4/examples/common"
	"github.com/MetaRPC/GoMT4/mt4"
	pb "github.com/MetaRPC/GoMT4/package"
)

// Streaming exercises the MT4Account server streams, each bounded by a short
// timeout / frame cap so the demo returns.
func Streaming(ctx context.Context, acc *mt4.MT4Account, symbol string) {
	fmt.Println("\n=== 📂 MT4Account · Streaming (bounded) ===")

	sctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	d1, e1 := acc.OnSymbolTick(sctx, []string{symbol})
	common.ReadStream("OnSymbolTick", d1, e1, 3, 6*time.Second, func(d *pb.OnSymbolTickData) string {
		t := d.GetSymbolTick()
		return fmt.Sprintf("%s bid %.5f ask %.5f", t.GetSymbol(), t.GetBid(), t.GetAsk())
	})

	d2, e2 := acc.OnOpenedOrdersProfit(sctx, 1000)
	common.ReadStream("OnOpenedOrdersProfit", d2, e2, 2, 5*time.Second, func(d *pb.OnOpenedOrdersProfitData) string {
		return fmt.Sprintf("%d orders", len(d.GetOpenedOrdersWithProfitUpdated()))
	})

	d3, e3 := acc.OnOpenedOrdersTickets(sctx, 1000)
	common.ReadStream("OnOpenedOrdersTickets", d3, e3, 2, 5*time.Second, func(d *pb.OnOpenedOrdersTicketsData) string {
		return fmt.Sprintf("%d position tickets", len(d.GetPositionTickets()))
	})

	d4, e4 := acc.OnTrade(sctx)
	common.ReadStream("OnTrade", d4, e4, 1, 3*time.Second, func(d *pb.OnTradeData) string { return "trade event" })

	d5, e5 := acc.OnConnectState(sctx, nil)
	common.ReadStream("OnConnectState", d5, e5, 1, 3*time.Second, func(d *pb.ConnectStateData) string { return d.GetState() })
}
