//+build linux

package watchdog

import (
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

// All system calls in this code are part of the Linux watchdog API. For
// reference, see:
// https://www.kernel.org/doc/html/latest/watchdog/watchdog-api.html.

func open() (*Device, error) {
	// TODO(mdlayher): determine the significance of the "/dev/watchdogN" nodes
	// on Linux. It appears that my machine with only one device exposes both
	// "/dev/watchdog" and "/dev/watchdog0".
	//
	// According to Terin, /dev/watchdog is an alias for /dev/watchdog0 on
	// modern machines. It's possible there could be more than one device, so
	// we'll eventually want to support that.
	f, err := os.OpenFile("/dev/watchdog", os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}

	// Immediately fetch the device's information to return to the caller.
	info, err := unix.IoctlGetWatchdogInfo(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	return &Device{
		// Clean up any trailing NULL bytes.
		Identity: strings.TrimRight(string(info.Identity[:]), "\x00"),

		f: f,
	}, nil
}

func (d *Device) ping() error { return unix.IoctlWatchdogKeepalive(int(d.f.Fd())) }

func (d *Device) close() error {
	// Attempt a Magic Close to disarm the watchdog device, since any call to
	// Close would be intentional and it's unlikely the user would want a system
	// reboot. Reference:
	// https://www.kernel.org/doc/html/latest/watchdog/watchdog-api.html#magic-close-feature
	if _, err := d.f.Write([]byte("V")); err != nil {
		// Make sure the file descriptor is closed even if Magic Close fails.
		_ = d.f.Close()
		return err
	}

	return d.f.Close()
}
