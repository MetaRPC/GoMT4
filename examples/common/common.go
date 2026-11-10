// Package common holds helpers shared by the grouped example packages.
package common

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TradingEnabled gates the live order-sending demos. Set MT4_ENABLE_TRADING=1 to
// run them (demo accounts only!).
var TradingEnabled = os.Getenv("MT4_ENABLE_TRADING") == "1"

// Ptr returns a pointer to v — handy for the pointer-valued optional arguments.
func Ptr[T any](v T) *T { return &v }

// ReadStream prints up to maxFrames frames from a (data, err) channel pair,
// stopping on error, close, or timeout.
func ReadStream[T any](label string, data <-chan *T, errs <-chan error, maxFrames int, timeout time.Duration, show func(*T) string) {
	t := time.After(timeout)
	n := 0
	for n < maxFrames {
		select {
		case d, ok := <-data:
			if !ok {
				fmt.Printf("  [%s] stream closed after %d frames\n", label, n)
				return
			}
			if d != nil {
				fmt.Printf("  [%s] %s\n", label, show(d))
				n++
			}
		case err := <-errs:
			if err != nil {
				// Event streams (OnTrade, OnConnectState) simply have no frames in a
				// quiet window; the context deadline then fires. That's not a failure.
				if status.Code(err) == codes.DeadlineExceeded {
					fmt.Printf("  [%s] no events in window (%d frames)\n", label, n)
				} else {
					fmt.Printf("  [%s] error: %v\n", label, err)
				}
				return
			}
		case <-t:
			fmt.Printf("  [%s] stopped (%d frames)\n", label, n)
			return
		}
	}
	fmt.Printf("  [%s] done (%d frames)\n", label, n)
}
