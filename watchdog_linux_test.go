//+build linux

package watchdog_test

import (
	"errors"
	"os"
	"testing"

	"github.com/mdlayher/watchdog"
)

func TestIntegrationDevice(t *testing.T) {
	// Since this test requires the presence of specific hardware, it's highly
	// likely to not pass on regular machines. Make sure to check the errors
	// and skip when necessary.
	d, err := watchdog.Open()
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			t.Skipf("skipping, watchdog device does not exist: %v", err)
		case errors.Is(err, os.ErrPermission):
			t.Skipf("skipping, permission denied (try running as root): %v", err)
		default:
			t.Fatalf("failed to open device: %v", err)
		}
	}

	// Double-close to ensure we always disarm the watchdog and don't reboot
	// the host.
	defer d.Close()

	t.Logf("device: %q", d.Identity)

	if err := d.Ping(); err != nil {
		t.Fatalf("failed to ping: %v", err)
	}

	if err := d.Close(); err != nil {
		t.Fatalf("failed to disarm: %v", err)
	}
}
