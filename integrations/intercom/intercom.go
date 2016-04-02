package intercom

import (
	"log"
	"os"

	"github.com/intercom/intercom-go"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Intercom integration
type Intercom struct {
}

// Identify forwards and identify call to Intercom
func (i Intercom) Identify(event integrations.Event) (err error) {
	log.Printf("NOT IMPLEMENTED: will send %#v to Intercom\n", event)

	ic := intercom.NewClient(appID(), apiKey())
	icUser, err := ic.Users.FindByUserID(event.UserID)
	if err != nil {
		log.Println("Error fetching the Intercom user:", err)
		return
	}

	icUser.CustomAttributes = event.UserTraits

	if event.UserTraits["email"] != nil {
		icUser.Email = event.UserTraits["email"].(string)
	}

	if event.UserTraits["name"] != nil {
		icUser.Name = event.UserTraits["name"].(string)
	}

	if event.UserTraits["createdAt"] != nil {
		icUser.CreatedAt = int32(event.UserTraits["createdAt"].(float64))
		icUser.SignedUpAt = int32(event.UserTraits["createdAt"].(float64))
	}

	savedUser, err := ic.Users.Save(&icUser)
	if err == nil {
		log.Printf("User saved on Intercom: %#v\n", savedUser)
	} else {
		log.Println("Error while saving on Intercom:", err)
	}
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
