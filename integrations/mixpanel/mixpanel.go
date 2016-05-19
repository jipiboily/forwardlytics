package mixpanel

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Mixpanel integration
type Mixpanel struct {
}

// Identify forwards and identify call to Mixpanel
func (Mixpanel) Identify(identification integrations.Identification) (err error) {
	logrus.Errorf("NOT IMPLEMENTED: will send %#v to Mixpanel\n", identification)
	return
}

// Track forwards the event to Mixpanel
func (Mixpanel) Track(event integrations.Event) (err error) {
	logrus.Errorf("NOT IMPLEMENTED: will send %#v to Mixpanel\n", event)
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
