// Package watchdog implements control of hardware watchdog devices.
package watchdog

import (
	"os"
	"time"
)

// A Device is a hardware watchdog device which can be pinged to keep the system
// from rebooting once the device has been opened.
type Device struct {
	// Identity is the name of the watchdog driver.
	Identity string

	f *os.File
}

// Open opens the primary watchdog device on this system ("/dev/watchdog" on
// Linux, TBD on other platforms). If the device is not found, an error
// compatible with os.ErrNotExist will be returned.
//
// Once a Device is opened, you must call Ping repeatedly to keep the system
// from being rebooted. Call Close to disarm the watchdog device.
func Open() (*Device, error) { return open() }

// Ping pings the watchdog device to keep the device from rebooting the system.
func (d *Device) Ping() error { return d.ping() }

// Timeout returns the configured timeout of the watchdog device.
func (d *Device) Timeout() (time.Duration, error) { return d.timeout() }

// Close closes the handle to the watchdog device and attempts to gracefully
// disarm the device, so no further Ping calls are required to keep the system
// from rebooting.
func (d *Device) Close() error { return d.close() }
