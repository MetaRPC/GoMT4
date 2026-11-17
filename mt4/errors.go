package mt4

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TradeOutcomeUnknownCode is the error code of TradeOutcomeUnknownError.
const TradeOutcomeUnknownCode = "TRADE_OUTCOME_UNKNOWN"

// TradeOutcomeUnknownError is returned by OrderSend, OrderModify, OrderClose/OrderDelete and
// OrderCloseBy (and the service/sugar methods built on them) when the call failed at the
// transport level with gRPC Unavailable, DeadlineExceeded, Unknown or Internal. The request
// may already have reached the terminal and been executed, so the SDK did NOT re-send it:
// check opened orders before retrying.
//
// It wraps the original gRPC status error: errors.As(err, &*TradeOutcomeUnknownError),
// status.FromError(err) and status.Code(err) all keep working.
type TradeOutcomeUnknownError struct {
	Code      string // always TradeOutcomeUnknownCode
	Operation string // "OrderSend", "OrderModify", "OrderCloseDelete" or "OrderCloseBy"
	Err       error  // the original gRPC status error
}

func newTradeOutcomeUnknownError(operation string, err error) *TradeOutcomeUnknownError {
	return &TradeOutcomeUnknownError{Code: TradeOutcomeUnknownCode, Operation: operation, Err: err}
}

// isTradeOutcomeUnknownCode: transport codes after which a trade may have been executed.
func isTradeOutcomeUnknownCode(c codes.Code) bool {
	switch c {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Unknown, codes.Internal:
		return true
	}
	return false
}

// isApiKeyRefused reports whether a gRPC status message carries a refused API key code.
func isApiKeyRefused(msg string) bool {
	return strings.Contains(msg, "API_KEY_INVALID") || strings.Contains(msg, "API_KEY_REQUIRED")
}

// Error implements the error interface.
func (e *TradeOutcomeUnknownError) Error() string {
	return fmt.Sprintf("%s: %s failed with gRPC %s after the request may have reached the terminal; "+
		"the trade may or may not have been executed. It was NOT re-sent automatically: "+
		"check open positions/orders before retrying. (%v)",
		TradeOutcomeUnknownCode, e.Operation, grpcCodeName(status.Code(e.Err)), e.Err)
}

// ErrorCode returns TradeOutcomeUnknownCode.
func (e *TradeOutcomeUnknownError) ErrorCode() string { return TradeOutcomeUnknownCode }

// Unwrap returns the original gRPC status error.
func (e *TradeOutcomeUnknownError) Unwrap() error { return e.Err }

// GRPCStatus returns the status of the original gRPC error, so status.FromError and
// status.Code see the original code (e.g. codes.Unavailable).
func (e *TradeOutcomeUnknownError) GRPCStatus() *status.Status {
	s, _ := status.FromError(e.Err)
	return s
}

// grpcCodeName returns the upper-case gRPC code name, e.g. "UNAVAILABLE".
func grpcCodeName(c codes.Code) string {
	switch c {
	case codes.Unavailable:
		return "UNAVAILABLE"
	case codes.DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case codes.Unknown:
		return "UNKNOWN"
	case codes.Internal:
		return "INTERNAL"
	}
	return strings.ToUpper(c.String())
}
