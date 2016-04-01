package mixpanel

import (
	"log"
	"os"

	"github.com/jipiboily/forwardlytics/integrations"
)

// Mixpanel integration
type Mixpanel struct {
}

// Identify forwards and identify call to Mixpanel
func (Mixpanel) Identify(event integrations.Event) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Mixpanel\n", event)
	return
}

// Enabled returns wether or not the Mixpanel integration is enabled/configured
func (Mixpanel) Enabled() bool {
	return apiKey() != "" && token() != ""
}

func apiKey() string {
	return os.Getenv("MIXPANEL_API_KEY")
}

func token() string {
	return os.Getenv("MIXPANEL_TOKEN")
}

func init() {
	integrations.RegisterIntegration("mixpanel", Mixpanel{})
}
