package drip

import (
	"log"

	"github.com/jipiboily/forwardlytics/integrations"
)

// Drip integration
type Drip struct {
}

// Identify forwards and identify call to Drip.io
func (Drip) Identify(user integrations.User) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Mixpanel\n", user)
	return
}

// Enabled returns wether or not the Drip.io integration is enabled/configured
func (Drip) Enabled() bool {
	return false
}

func init() {
	integrations.RegisterIntegration("drip", Drip{})
}
