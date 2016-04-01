package intercom

import (
	"log"
	"os"

	"github.com/jipiboily/forwardlytics/integrations"
	"gopkg.in/intercom/intercom-go.v1"
)

// Intercom integration
type Intercom struct {
}

// Identify forwards and identify call to Intercom
func (i Intercom) Identify(user integrations.User) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Intercom\n", user)

	// what to do if there is an error? Log in BugSnag, but besides that....retry?
	ic := intercom.NewClient(appID(), apiKey())
	icUser, err := ic.Users.FindByEmail("z4@metrics.watch")
	if err != nil {
		log.Fatalf("Error fetching the user on Intercom: %s", err)
	}
	log.Printf("User on Intercom: %#v\n", icUser)

	return
}

// Enabled returns wether or not the Intercom integration is enabled/configured
func (i Intercom) Enabled() bool {
	return apiKey() != "" && appID() != ""
}

func apiKey() string {
	return os.Getenv("INTERCOM_API_KEY")
}

func appID() string {
	return os.Getenv("INTERCOM_APP_ID")
}

func init() {
	integrations.RegisterIntegration("intercom", Intercom{})
}
