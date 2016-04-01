package keen

import (
	"log"

	"github.com/jipiboily/forwardlytics/integrations"
)

// Keen integration
type Keen struct {
}

// Identify forwards and identify call to Keen.io
func (Keen) Identify(user integrations.User) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Mixpanel\n", user)
	return
}

// Enabled returns wether or not the Keen.io integration is enabled/configured
func (Keen) Enabled() bool {
	return false
}

func init() {
	integrations.RegisterIntegration("keen", Keen{})
}
