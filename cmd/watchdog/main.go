// Command watchdog is a demo application which shows usage of the package
// watchdog APIs.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mdlayher/watchdog"
)

func main() {
	d, err := watchdog.Open()
	if err != nil {
		log.Fatalf("failed to open watchdog: %v", err)
	}

	// We purposely double-close the file to ensure that the explicit Close
	// later on also disarms the device as the program exits. Otherwise it's
	// possible we may exit early or with a subtle error and leave the system
	// in a doomed state.
	defer d.Close()

	timeout, err := d.Timeout()
	if err != nil {
		log.Fatalf("failed to fetch watchdog timeout: %v", err)
	}

	fmt.Printf("device: %q, timeout: %s\n", d.Identity, timeout)

	for i := 0; i < 3; i++ {
		if err := d.Ping(); err != nil {
			log.Fatalf("failed to ping watchdog: %v", err)
		}

		// Note progress to the caller.
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}

	fmt.Println()

	// Safely disarm the device before exiting.
	if err := d.Close(); err != nil {
		log.Fatalf("failed to disarm watchdog: %v", err)
	}
}
