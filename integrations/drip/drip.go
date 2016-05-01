package drip

import (
	"log"

	"github.com/jipiboily/forwardlytics/integrations"
)

// Drip integration
type Drip struct {
}

// Identify forwards and identify call to Drip
func (Drip) Identify(identification integrations.Identification) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Drip\n", identification)
	return
}

// Track forwards the event to Drip
func (Drip) Track(event integrations.Event) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Drip\n", event)
	return
}

// Enabled returns wether or not the Drip integration is enabled/configured
func (Drip) Enabled() bool {
	return false
}

func init() {
	integrations.RegisterIntegration("drip", Drip{})
}
