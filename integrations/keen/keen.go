package keen

import (
	"log"

	"github.com/jipiboily/dickspatch/integrations"
)

// Keen integration
type Keen struct {
}

// Identify forwards and identify call to Keen.io
func (Keen) Identify(event integrations.Event) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Mixpanel\n", event)
	return
}

// Enabled returns wether or not the Keen.io integration is enabled/configured
func (Keen) Enabled() bool {
	return false
}

func init() {
	integrations.RegisterIntegration("keen", Keen{})
}
