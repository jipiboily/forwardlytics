package drip

import (
	"log"

	"github.com/jipiboily/forwardlytics/integrations"
)

// Drip integration
type Drip struct {
}

// Identify forwards and identify call to Drip.io
func (Drip) Identify(identification integrations.Identification) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Drip\n", identification)
	return
}

// Enabled returns wether or not the Drip.io integration is enabled/configured
func (Drip) Enabled() bool {
	return false
}

func init() {
	integrations.RegisterIntegration("drip", Drip{})
}
