//+build !linux

package watchdog

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// errNotImplemented is a sentinel which indicates package watchdog does not
// support this OS.
var errNotImplemented = fmt.Errorf("watchdog: not implemented on %s: %w", runtime.GOOS, os.ErrNotExist)

func open() (*Device, error)                    { return nil, errNotImplemented }
func (*Device) ping() error                     { return errNotImplemented }
func (*Device) timeout() (time.Duration, error) { return 0, errNotImplemented }
func (*Device) close() error                    { return errNotImplemented }
