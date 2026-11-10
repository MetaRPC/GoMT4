package mt4

// MT4Sugar is the high-level convenience API over MT4Service (Layer 3).
//
// It mirrors the three-layer SDK design: one-liners and business logic — quick
// prices, symbol math (point/pip/normalize), risk-based lot sizing, order entry
// by side, order management, trailing/breakeven, watchers/waiters and order
// info. It composes an MT4Service; borrow it with Service().

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MetaRPC/GoMT4/package"
)

// MT4Sugar is the high-level convenience layer. It wraps an MT4Service and adds
// one-liner helpers so trading code reads like intent rather than plumbing.
//
// Typical lifecycle — build the three layers, set defaults, trade:
//
//	acc := mt4.NewMT4Account(login, password, "", uuid.Nil, apiKey)
//	sugar := mt4.NewMT4Sugar(mt4.NewMT4Service(acc)).SetDefaults(777, 20, "my-bot")
//	_ = sugar.ConnectByServerName(ctx, "MetaQuotes-Demo", "EURUSD", true, 30)
//
//	// size a position that risks 1% of balance over a 20-pip stop, then buy it
//	bal := must(sugar.GetAccountInfo(ctx)).GetAccountBalance()
//	lots, _ := sugar.CalcLotByRisk(ctx, "EURUSD", bal*0.01, 20)
//	r, _ := sugar.BuyMarket(ctx, "EURUSD", lots, nil, nil)
//	_ = sugar.ModifySLTPByPips(ctx, r.GetTicket(), ptr(20.0), ptr(40.0))
type MT4Sugar struct {
	svc      *MT4Service
	magic    int32  // stamped on order-entry helpers when non-zero
	slippage int32  // max price deviation in points (default 20)
	comment  string // stamped on order-entry helpers when non-empty
}

// NewMT4Sugar wraps a service with sensible defaults (magic 0, slippage 20,
// empty comment). Chain SetDefaults to change them:
//
//	sugar := mt4.NewMT4Sugar(svc).SetDefaults(777, 30, "scalper")
func NewMT4Sugar(svc *MT4Service) *MT4Sugar {
	return &MT4Sugar{svc: svc, slippage: 20}
}

// SetDefaults sets the magic / slippage (points) / comment used by the
// order-entry helpers. Returns the receiver for chaining.
func (s *MT4Sugar) SetDefaults(magic, slippage int32, comment string) *MT4Sugar {
	s.magic = magic
	s.slippage = slippage
	s.comment = comment
	return s
}

// Service returns the underlying MT4Service.
func (s *MT4Sugar) Service() *MT4Service { return s.svc }

// sym fetches the contract parameters for one symbol.
func (s *MT4Sugar) sym(ctx context.Context, symbol string) (*pb.SymbolParamsManyInfo, error) {
	return s.svc.Account().SymbolParams(ctx, symbol)
}

// ===== 1. Connection & defaults =====

// ConnectByServerName resolves the broker by name and logs in.
func (s *MT4Sugar) ConnectByServerName(ctx context.Context, server, defaultSymbol string, autoSelect bool, timeoutSeconds int) error {
	return s.svc.Account().ConnectByServerName(ctx, server, defaultSymbol, autoSelect, timeoutSeconds)
}

// Disconnect closes the session (delete=true also removes the terminal instance).
func (s *MT4Sugar) Disconnect(delete bool) error {
	return s.svc.Account().Disconnect(delete)
}

// Ping is a liveness probe — true if an account snapshot succeeds.
func (s *MT4Sugar) Ping(ctx context.Context) bool {
	_, err := s.svc.AccountSummary(ctx)
	return err == nil
}

// ===== 2. Prices & quotes =====

// LastQuote returns the latest quote for a symbol.
func (s *MT4Sugar) LastQuote(ctx context.Context, symbol string) (*pb.QuoteData, error) {
	return s.svc.Quote(ctx, symbol)
}

// Quotes returns the latest quotes for several symbols.
func (s *MT4Sugar) Quotes(ctx context.Context, symbols []string) ([]*pb.QuoteData, error) {
	data, err := s.svc.QuoteMany(ctx, symbols)
	if err != nil {
		return nil, err
	}
	return data.GetQuotes(), nil
}

// MidPrice returns (bid + ask) / 2.
func (s *MT4Sugar) MidPrice(ctx context.Context, symbol string) (float64, error) {
	q, err := s.svc.Quote(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return (q.GetBid() + q.GetAsk()) / 2.0, nil
}

// SpreadPips returns the current spread (ask − bid) expressed in pips. A common
// use is to skip trading when the spread is wide (news, thin liquidity).
//
// Example:
//
//	if sp, _ := sugar.SpreadPips(ctx, "EURUSD"); sp > 2.0 {
//	    log.Println("spread too wide, skipping")
//	}
func (s *MT4Sugar) SpreadPips(ctx context.Context, symbol string) (float64, error) {
	q, err := s.svc.Quote(ctx, symbol)
	if err != nil {
		return 0, err
	}
	pip, err := s.PipSize(ctx, symbol)
	if err != nil {
		return 0, err
	}
	if pip <= 0 {
		return 0, nil
	}
	return (q.GetAsk() - q.GetBid()) / pip, nil
}

// EnsureSymbol adds a symbol to Market Watch and then polls until it publishes a
// usable two-sided quote (bid > 0 and ask > bid), or the wait elapses. Call it
// once before trading a symbol you haven't touched — fresh/demo feeds often send
// bid == ask or zeros for the first moment.
//
// Example:
//
//	q, err := sugar.EnsureSymbol(ctx, "EURUSD", 10*time.Second)
//	if err != nil { return err } // symbol never produced a usable quote
//	fmt.Println("ready:", q.GetBid(), q.GetAsk())
func (s *MT4Sugar) EnsureSymbol(ctx context.Context, symbol string, wait time.Duration) (*pb.QuoteData, error) {
	if _, err := s.svc.SymbolSelect(ctx, symbol, true); err != nil {
		return nil, err
	}
	deadline := time.Now().Add(wait)
	for {
		if q, err := s.svc.Quote(ctx, symbol); err == nil {
			if q.GetBid() > 0 && q.GetAsk() > q.GetBid() {
				return q, nil
			}
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("%s: no usable quote within %s", symbol, wait)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// Bars returns OHLC history covering roughly the newest `count` bars of the
// given timeframe (it computes the time window for you).
//
// Example — last 100 H1 candles:
//
//	h, _ := sugar.Bars(ctx, "EURUSD", pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1, 100)
//	fmt.Println("bars:", len(h.GetHistoricalQuotes()))
func (s *MT4Sugar) Bars(ctx context.Context, symbol string, timeframe pb.ENUM_QUOTE_HISTORY_TIMEFRAME, count int) (*pb.QuoteHistoryData, error) {
	now := time.Now()
	from := now.Add(-time.Duration(tfSeconds(timeframe)*int64(count)*3) * time.Second)
	return s.svc.QuoteHistory(ctx, symbol, timeframe, from, now.Add(time.Hour))
}

// ===== 3. Symbol math =====

// Point returns the point size (SYMBOL_POINT).
func (s *MT4Sugar) Point(ctx context.Context, symbol string) (float64, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return si.GetPoint(), nil
}

// Digits returns the price decimal places.
func (s *MT4Sugar) Digits(ctx context.Context, symbol string) (int32, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return si.GetDigits(), nil
}

// PipSize returns point*10 on 3/5-digit symbols, else point.
func (s *MT4Sugar) PipSize(ctx context.Context, symbol string) (float64, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	if si.GetDigits() == 3 || si.GetDigits() == 5 {
		return si.GetPoint() * 10.0, nil
	}
	return si.GetPoint(), nil
}

// PipsToPrice converts a pip distance into a price distance.
func (s *MT4Sugar) PipsToPrice(ctx context.Context, symbol string, pips float64) (float64, error) {
	pip, err := s.PipSize(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return pips * pip, nil
}

// PriceToPips converts a price distance into pips.
func (s *MT4Sugar) PriceToPips(ctx context.Context, symbol string, priceDelta float64) (float64, error) {
	pip, err := s.PipSize(ctx, symbol)
	if err != nil {
		return 0, err
	}
	if pip <= 0 {
		return 0, nil
	}
	return priceDelta / pip, nil
}

// NormalizeLot rounds a lot to the symbol's volume step and clamps it to
// [volumeMin, volumeMax]. Always run a computed volume through this before
// sending it, or the broker may reject it as an invalid volume (error 131).
//
// Example:
//
//	lots, _ := sugar.NormalizeLot(ctx, "EURUSD", 0.137) // e.g. -> 0.14 on a 0.01 step
func (s *MT4Sugar) NormalizeLot(ctx context.Context, symbol string, lot float64) (float64, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	step := si.GetVolumeStep()
	if step <= 0 {
		step = 0.01
	}
	v := roundTo(lot/step) * step
	if v < si.GetVolumeMin() {
		v = si.GetVolumeMin()
	} else if si.GetVolumeMax() > 0 && v > si.GetVolumeMax() {
		v = si.GetVolumeMax()
	}
	return v, nil
}

// NormalizePrice rounds a price to the symbol's digits.
func (s *MT4Sugar) NormalizePrice(ctx context.Context, symbol string, price float64) (float64, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	d := si.GetDigits()
	if d < 0 {
		d = 0
	}
	f := pow10(int(d))
	return roundTo(price*f) / f, nil
}

// TickValue returns the money value of one tick per lot.
func (s *MT4Sugar) TickValue(ctx context.Context, symbol string) (float64, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return si.GetTradeTickValue(), nil
}

// ===== 4. Risk sizing =====

// CalcCashRisk returns the cash at risk for `lots` over an `slPips` stop.
func (s *MT4Sugar) CalcCashRisk(ctx context.Context, symbol string, lots, slPips float64) (float64, error) {
	pv, err := s.pipValuePerLot(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return lots * slPips * pv, nil
}

// CalcLotByRisk returns the normalized lot size that risks riskCash (in the
// account currency) over an slPips stop. It converts cash → lots using the
// symbol's pip value, then rounds to the volume step and clamps to [min, max].
// Guards: if slPips <= 0 or the pip value is unknown it returns the symbol's
// minimum lot instead of erroring.
//
// Example — risk 1% of balance on a 25-pip stop:
//
//	bal := must(sugar.GetAccountInfo(ctx)).GetAccountBalance()
//	lots, err := sugar.CalcLotByRisk(ctx, "EURUSD", bal*0.01, 25)
//	// e.g. lots == 0.12
func (s *MT4Sugar) CalcLotByRisk(ctx context.Context, symbol string, riskCash, slPips float64) (float64, error) {
	pv, err := s.pipValuePerLot(ctx, symbol)
	if err != nil {
		return 0, err
	}
	if slPips <= 0 || pv <= 0 {
		return s.NormalizeLot(ctx, symbol, 0)
	}
	return s.NormalizeLot(ctx, symbol, riskCash/(slPips*pv))
}

func (s *MT4Sugar) pipValuePerLot(ctx context.Context, symbol string) (float64, error) {
	si, err := s.sym(ctx, symbol)
	if err != nil {
		return 0, err
	}
	pip := si.GetPoint()
	if si.GetDigits() == 3 || si.GetDigits() == 5 {
		pip = si.GetPoint() * 10.0
	}
	if si.GetTradeTickSize() > 0 {
		return si.GetTradeTickValue() * (pip / si.GetTradeTickSize()), nil
	}
	return si.GetTradeTickValue(), nil
}

// ===== 5. Order entry =====

// send builds the OrderSend call with the configured defaults.
func (s *MT4Sugar) send(ctx context.Context, symbol string, op pb.OrderSendOperationType, volume float64, price, sl, tp *float64) (*pb.OrderSendData, error) {
	var slip *int32
	if s.slippage != 0 {
		v := s.slippage
		slip = &v
	}
	var magic *int32
	if s.magic != 0 {
		v := s.magic
		magic = &v
	}
	var comment *string
	if s.comment != "" {
		c := s.comment
		comment = &c
	}
	return s.svc.OrderSend(ctx, symbol, op, volume, price, slip, sl, tp, comment, magic, nil)
}

// BuyMarket opens a market buy at the current ask. sl/tp are absolute prices
// (pass nil to skip, or set them later by pips with ModifySLTPByPips). The
// default magic / slippage / comment from SetDefaults are applied.
//
// Example — buy 0.10 lot, then attach a 20-pip stop and 40-pip target:
//
//	r, err := sugar.BuyMarket(ctx, "EURUSD", 0.10, nil, nil)
//	if err != nil { return err }
//	_ = sugar.ModifySLTPByPips(ctx, r.GetTicket(), ptr(20.0), ptr(40.0))
func (s *MT4Sugar) BuyMarket(ctx context.Context, symbol string, lots float64, sl, tp *float64) (*pb.OrderSendData, error) {
	return s.send(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUY, lots, nil, sl, tp)
}

// SellMarket opens a market sell.
func (s *MT4Sugar) SellMarket(ctx context.Context, symbol string, lots float64, sl, tp *float64) (*pb.OrderSendData, error) {
	return s.send(ctx, symbol, pb.OrderSendOperationType_OC_OP_SELL, lots, nil, sl, tp)
}

// BuyLimit places a buy-limit pending order at price (which must be below the
// current ask — it fills when price drops to it). The SellLimit/BuyStop/SellStop
// siblings follow the same shape.
//
// Example — buy the dip 15 pips below the current ask:
//
//	ask := must(sugar.LastQuote(ctx, "EURUSD")).GetAsk()
//	dist, _ := sugar.PipsToPrice(ctx, "EURUSD", 15)
//	r, _ := sugar.BuyLimit(ctx, "EURUSD", 0.10, ask-dist, nil, nil)
func (s *MT4Sugar) BuyLimit(ctx context.Context, symbol string, lots, price float64, sl, tp *float64) (*pb.OrderSendData, error) {
	return s.send(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUYLIMIT, lots, &price, sl, tp)
}

// SellLimit places a sell-limit pending order at price.
func (s *MT4Sugar) SellLimit(ctx context.Context, symbol string, lots, price float64, sl, tp *float64) (*pb.OrderSendData, error) {
	return s.send(ctx, symbol, pb.OrderSendOperationType_OC_OP_SELLLIMIT, lots, &price, sl, tp)
}

// BuyStop places a buy-stop pending order at price.
func (s *MT4Sugar) BuyStop(ctx context.Context, symbol string, lots, price float64, sl, tp *float64) (*pb.OrderSendData, error) {
	return s.send(ctx, symbol, pb.OrderSendOperationType_OC_OP_BUYSTOP, lots, &price, sl, tp)
}

// SellStop places a sell-stop pending order at price.
func (s *MT4Sugar) SellStop(ctx context.Context, symbol string, lots, price float64, sl, tp *float64) (*pb.OrderSendData, error) {
	return s.send(ctx, symbol, pb.OrderSendOperationType_OC_OP_SELLSTOP, lots, &price, sl, tp)
}

// ===== 6. Order management =====

// Close closes a market order fully (lots=nil) or partially (lots=&volume), or
// deletes a pending order. It uses the configured default slippage.
//
// Example:
//
//	sugar.Close(ctx, ticket, nil)          // full close / delete
//	sugar.Close(ctx, ticket, ptr(0.05))    // close 0.05 lot of it
func (s *MT4Sugar) Close(ctx context.Context, ticket int32, lots *float64) (*pb.OrderCloseDeleteData, error) {
	slip := s.slippage
	return s.svc.OrderClose(ctx, ticket, lots, nil, &slip)
}

// ClosePartial closes part of a market order by volume.
func (s *MT4Sugar) ClosePartial(ctx context.Context, ticket int32, lots float64) (*pb.OrderCloseDeleteData, error) {
	return s.Close(ctx, ticket, &lots)
}

// CloseBy closes an order by an opposite one.
func (s *MT4Sugar) CloseBy(ctx context.Context, ticket, oppositeTicket int32) (*pb.OrderCloseByData, error) {
	return s.svc.OrderCloseBy(ctx, ticket, oppositeTicket)
}

// CloseAll closes every open market order, or only those on symbol when it is
// non-empty. Best-effort: a single failure does not abort the rest. Returns how
// many were closed.
//
// Example:
//
//	n, _ := sugar.CloseAll(ctx, "")        // flatten everything
//	n, _ := sugar.CloseAll(ctx, "EURUSD")  // flatten EURUSD only
func (s *MT4Sugar) CloseAll(ctx context.Context, symbol string) (int, error) {
	positions, err := s.Positions(ctx)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, o := range positions {
		if symbol != "" && o.GetSymbol() != symbol {
			continue
		}
		if _, err := s.Close(ctx, o.GetTicket(), nil); err == nil {
			n++
		}
	}
	return n, nil
}

// CancelPendings deletes every pending order (optionally only on `symbol`);
// returns how many deleted.
func (s *MT4Sugar) CancelPendings(ctx context.Context, symbol string) (int, error) {
	orders, err := s.Orders(ctx)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, o := range orders {
		if !isPending(o.GetOrderType()) {
			continue
		}
		if symbol != "" && o.GetSymbol() != symbol {
			continue
		}
		if _, err := s.Close(ctx, o.GetTicket(), nil); err == nil {
			n++
		}
	}
	return n, nil
}

// ModifySLTPByPrice sets a position's SL/TP to absolute prices (nil leaves a side).
func (s *MT4Sugar) ModifySLTPByPrice(ctx context.Context, ticket int32, sl, tp *float64) error {
	_, err := s.svc.OrderModify(ctx, ticket, nil, sl, tp, nil)
	return err
}

// ModifySLTPByPips sets SL/TP as pip distances from the order's open price,
// handling the buy/sell direction for you (for a BUY the stop goes below and the
// target above the entry; for a SELL it is mirrored). Pass nil for a side to
// leave it unchanged. It looks the order up to read its side and open price.
//
// Example — 20-pip stop, 40-pip target, regardless of side:
//
//	err := sugar.ModifySLTPByPips(ctx, ticket, ptr(20.0), ptr(40.0))
//
// Move only the stop (leave the target):
//
//	err := sugar.ModifySLTPByPips(ctx, ticket, ptr(15.0), nil)
func (s *MT4Sugar) ModifySLTPByPips(ctx context.Context, ticket int32, slPips, tpPips *float64) error {
	o, err := s.FindTicket(ctx, ticket)
	if err != nil {
		return err
	}
	if o == nil {
		return fmt.Errorf("order %d not open", ticket)
	}
	pip, err := s.PipSize(ctx, o.GetSymbol())
	if err != nil {
		return err
	}
	sign := 1.0
	if !isBuy(o.GetOrderType()) {
		sign = -1.0
	}
	var sl, tp *float64
	if slPips != nil {
		v := o.GetOpenPrice() - sign*(*slPips)*pip
		sl = &v
	}
	if tpPips != nil {
		v := o.GetOpenPrice() + sign*(*tpPips)*pip
		tp = &v
	}
	return s.ModifySLTPByPrice(ctx, ticket, sl, tp)
}

// ===== 7. Trailing & breakeven =====

// AutoBreakeven moves the stop to entry ± plusPips once the position is at least
// triggerPips in profit, and never loosens an existing stop. Like TrailStop it is
// a one-shot primitive — call it on a timer.
//
// Example — arm at +10 pips, lock +2 pips of profit:
//
//	moved, _ := sugar.AutoBreakeven(ctx, ticket, 10, 2)
//	if moved { log.Println("stop moved to breakeven+2") }
func (s *MT4Sugar) AutoBreakeven(ctx context.Context, ticket int32, triggerPips, plusPips float64) (bool, error) {
	o, err := s.FindTicket(ctx, ticket)
	if err != nil {
		return false, err
	}
	if o == nil {
		return false, fmt.Errorf("order %d not open", ticket)
	}
	pip, err := s.PipSize(ctx, o.GetSymbol())
	if err != nil {
		return false, err
	}
	q, err := s.svc.Quote(ctx, o.GetSymbol())
	if err != nil {
		return false, err
	}
	buy := isBuy(o.GetOrderType())
	var profitPips float64
	if buy {
		profitPips = (q.GetBid() - o.GetOpenPrice()) / pip
	} else {
		profitPips = (o.GetOpenPrice() - q.GetAsk()) / pip
	}
	if profitPips < triggerPips {
		return false, nil
	}
	var newSL float64
	if buy {
		newSL = o.GetOpenPrice() + plusPips*pip
	} else {
		newSL = o.GetOpenPrice() - plusPips*pip
	}
	if !stopIsBetter(buy, newSL, o.GetStopLoss()) {
		return false, nil
	}
	if err := s.ModifySLTPByPrice(ctx, ticket, &newSL, nil); err != nil {
		return false, err
	}
	return true, nil
}

// TrailStop ratchets the stop to distancePips behind the current price, only
// ever tightening (a call that would loosen the stop returns false and does
// nothing). It is a one-shot primitive — call it on a timer to trail live:
//
//	for {
//	    if o, _ := sugar.FindTicket(ctx, ticket); o == nil { break } // closed
//	    sugar.TrailStop(ctx, ticket, 15)                              // trail 15 pips
//	    time.Sleep(2 * time.Second)
//	}
//
// Returns whether the stop moved on this call.
func (s *MT4Sugar) TrailStop(ctx context.Context, ticket int32, distancePips float64) (bool, error) {
	o, err := s.FindTicket(ctx, ticket)
	if err != nil {
		return false, err
	}
	if o == nil {
		return false, fmt.Errorf("order %d not open", ticket)
	}
	pip, err := s.PipSize(ctx, o.GetSymbol())
	if err != nil {
		return false, err
	}
	q, err := s.svc.Quote(ctx, o.GetSymbol())
	if err != nil {
		return false, err
	}
	buy := isBuy(o.GetOrderType())
	var desired float64
	if buy {
		desired = q.GetBid() - distancePips*pip
	} else {
		desired = q.GetAsk() + distancePips*pip
	}
	if !stopIsBetter(buy, desired, o.GetStopLoss()) {
		return false, nil
	}
	if err := s.ModifySLTPByPrice(ctx, ticket, &desired, nil); err != nil {
		return false, err
	}
	return true, nil
}

// ===== 8. Watchers & waiters =====

// SubscribeQuotes streams live ticks for the given symbols on a pair of
// channels. Read frames off the data channel; the error channel delivers a
// terminal error. Cancel the context (or let it time out) to stop.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
//	defer cancel()
//	data, errs := sugar.SubscribeQuotes(ctx, []string{"EURUSD"})
//	for {
//	    select {
//	    case d := <-data:
//	        t := d.GetSymbolTick(); fmt.Println(t.GetSymbol(), t.GetBid(), t.GetAsk())
//	    case err := <-errs:
//	        if err != nil { return err }
//	    case <-ctx.Done():
//	        return nil
//	    }
//	}
func (s *MT4Sugar) SubscribeQuotes(ctx context.Context, symbols []string) (<-chan *pb.OnSymbolTickData, <-chan error) {
	return s.svc.StreamTicks(ctx, symbols)
}

// WatchProfit streams per-order floating profit every intervalMs.
func (s *MT4Sugar) WatchProfit(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersProfitData, <-chan error) {
	return s.svc.StreamOpenedProfit(ctx, intervalMs)
}

// WatchTickets streams the open-order ticket list every intervalMs.
func (s *MT4Sugar) WatchTickets(ctx context.Context, intervalMs int32) (<-chan *pb.OnOpenedOrdersTicketsData, <-chan error) {
	return s.svc.StreamOpenedTickets(ctx, intervalMs)
}

// WaitPrice polls until a symbol's price reaches target — with above=true it
// waits for bid >= target (price rose to it), with above=false for ask <= target
// (price fell to it). Returns the triggering quote, or an error on timeout.
//
// Example — wait for EURUSD to rise to 1.0950 (max 1 minute):
//
//	q, err := sugar.WaitPrice(ctx, "EURUSD", 1.0950, true, time.Minute)
//	if err != nil { return err } // did not reach the target in time
//	fmt.Println("reached, bid", q.GetBid())
func (s *MT4Sugar) WaitPrice(ctx context.Context, symbol string, target float64, above bool, timeout time.Duration) (*pb.QuoteData, error) {
	deadline := time.Now().Add(timeout)
	for {
		q, err := s.svc.Quote(ctx, symbol)
		if err != nil {
			return nil, err
		}
		hit := q.GetAsk() <= target
		if above {
			hit = q.GetBid() >= target
		}
		if hit {
			return q, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("%s: price did not reach %g within %s", symbol, target, timeout)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// WaitFilled polls (every 500 ms) until a pending ticket becomes a market order,
// i.e. it filled. Returns true on fill, false on timeout.
//
// Example — place a pending, then cancel it if it doesn't fill in a minute:
//
//	r, _ := sugar.BuyLimit(ctx, "EURUSD", 0.10, price, nil, nil)
//	if filled, _ := sugar.WaitFilled(ctx, r.GetTicket(), time.Minute); !filled {
//	    sugar.Close(ctx, r.GetTicket(), nil) // delete the resting pending
//	}
func (s *MT4Sugar) WaitFilled(ctx context.Context, ticket int32, timeout time.Duration) (bool, error) {
	deadline := time.Now().Add(timeout)
	for {
		o, err := s.FindTicket(ctx, ticket)
		if err != nil {
			return false, err
		}
		if o != nil && isMarket(o.GetOrderType()) {
			return true, nil
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// WaitClosed polls until `ticket` is no longer open (closed/deleted).
func (s *MT4Sugar) WaitClosed(ctx context.Context, ticket int32, timeout time.Duration) (bool, error) {
	deadline := time.Now().Add(timeout)
	for {
		o, err := s.FindTicket(ctx, ticket)
		if err != nil {
			return false, err
		}
		if o == nil {
			return true, nil
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// ===== 9. Order info =====

// Orders returns all open market + pending orders.
func (s *MT4Sugar) Orders(ctx context.Context) ([]*pb.OpenedOrderInfo, error) {
	return s.svc.OpenedOrders(ctx)
}

// Positions returns the open market orders only (order type Buy/Sell) — MT4's
// equivalent of "positions". Use Orders for everything including pendings.
//
// Example — sum floating P/L across positions:
//
//	var pl float64
//	for _, o := range must(sugar.Positions(ctx)) {
//	    pl += o.GetProfit() + o.GetSwap() + o.GetCommision()
//	}
func (s *MT4Sugar) Positions(ctx context.Context) ([]*pb.OpenedOrderInfo, error) {
	all, err := s.Orders(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*pb.OpenedOrderInfo, 0, len(all))
	for _, o := range all {
		if isMarket(o.GetOrderType()) {
			out = append(out, o)
		}
	}
	return out, nil
}

// OrdersForSymbol returns open orders on a symbol.
func (s *MT4Sugar) OrdersForSymbol(ctx context.Context, symbol string) ([]*pb.OpenedOrderInfo, error) {
	all, err := s.Orders(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*pb.OpenedOrderInfo, 0)
	for _, o := range all {
		if o.GetSymbol() == symbol {
			out = append(out, o)
		}
	}
	return out, nil
}

// OrdersByMagic returns open orders tagged with a magic number.
func (s *MT4Sugar) OrdersByMagic(ctx context.Context, magic int32) ([]*pb.OpenedOrderInfo, error) {
	all, err := s.Orders(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*pb.OpenedOrderInfo, 0)
	for _, o := range all {
		if o.GetMagicNumber() == magic {
			out = append(out, o)
		}
	}
	return out, nil
}

// FindTicket returns one open order by ticket, or nil if not present.
func (s *MT4Sugar) FindTicket(ctx context.Context, ticket int32) (*pb.OpenedOrderInfo, error) {
	all, err := s.Orders(ctx)
	if err != nil {
		return nil, err
	}
	for _, o := range all {
		if o.GetTicket() == ticket {
			return o, nil
		}
	}
	return nil, nil
}

// TotalProfit returns total floating profit across open positions (optionally on
// a symbol), including swap and commission.
func (s *MT4Sugar) TotalProfit(ctx context.Context, symbol string) (float64, error) {
	positions, err := s.Positions(ctx)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, o := range positions {
		if symbol != "" && o.GetSymbol() != symbol {
			continue
		}
		sum += o.GetProfit() + o.GetSwap() + o.GetCommision()
	}
	return sum, nil
}

// PositionsValue returns the floating value of positions filtered by symbol
// and/or magic (empty symbol / magic<0 = no filter on that axis).
func (s *MT4Sugar) PositionsValue(ctx context.Context, symbol string, magic int32) (float64, error) {
	positions, err := s.Positions(ctx)
	if err != nil {
		return 0, err
	}
	var sum float64
	for _, o := range positions {
		if symbol != "" && o.GetSymbol() != symbol {
			continue
		}
		if magic >= 0 && o.GetMagicNumber() != magic {
			continue
		}
		sum += o.GetProfit() + o.GetSwap() + o.GetCommision()
	}
	return sum, nil
}

// ExposureSummary returns net lots per symbol across open market positions, with
// buys positive and sells negative (so a hedged symbol nets toward zero).
//
// Example:
//
//	for sym, net := range must(sugar.ExposureSummary(ctx)) {
//	    fmt.Printf("%s: net %+.2f lots\n", sym, net)
//	}
func (s *MT4Sugar) ExposureSummary(ctx context.Context) (map[string]float64, error) {
	positions, err := s.Positions(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]float64)
	for _, o := range positions {
		if isBuy(o.GetOrderType()) {
			m[o.GetSymbol()] += o.GetLots()
		} else {
			m[o.GetSymbol()] -= o.GetLots()
		}
	}
	return m, nil
}

// GetAccountInfo returns the account snapshot.
func (s *MT4Sugar) GetAccountInfo(ctx context.Context) (*pb.AccountSummaryData, error) {
	return s.svc.AccountSummary(ctx)
}

// ===== free helpers =====

func isMarket(t pb.OpenedOrderType) bool {
	return t == pb.OpenedOrderType_OO_OP_BUY || t == pb.OpenedOrderType_OO_OP_SELL
}

func isPending(t pb.OpenedOrderType) bool {
	switch t {
	case pb.OpenedOrderType_OO_OP_BUYLIMIT, pb.OpenedOrderType_OO_OP_SELLLIMIT,
		pb.OpenedOrderType_OO_OP_BUYSTOP, pb.OpenedOrderType_OO_OP_SELLSTOP:
		return true
	}
	return false
}

func isBuy(t pb.OpenedOrderType) bool {
	return t == pb.OpenedOrderType_OO_OP_BUY || t == pb.OpenedOrderType_OO_OP_BUYLIMIT || t == pb.OpenedOrderType_OO_OP_BUYSTOP
}

// stopIsBetter reports whether newSL is a tighter (safer) stop than current.
// A current of 0 means "no stop yet", so any stop is better.
func stopIsBetter(buy bool, newSL, current float64) bool {
	if current == 0 {
		return true
	}
	if buy {
		return newSL > current
	}
	return newSL < current
}

func tfSeconds(tf pb.ENUM_QUOTE_HISTORY_TIMEFRAME) int64 {
	switch tf {
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_M1:
		return 60
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_M5:
		return 300
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_M15:
		return 900
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_M30:
		return 1800
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H1:
		return 3600
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_H4:
		return 14400
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_D1:
		return 86400
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_W1:
		return 604800
	case pb.ENUM_QUOTE_HISTORY_TIMEFRAME_QH_PERIOD_MN1:
		return 2592000
	}
	return 60
}

func roundTo(v float64) float64 {
	if v < 0 {
		return float64(int64(v - 0.5))
	}
	return float64(int64(v + 0.5))
}

func pow10(n int) float64 {
	r := 1.0
	for i := 0; i < n; i++ {
		r *= 10
	}
	return r
}
