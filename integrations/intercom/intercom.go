package intercom

import (
	"log"
	"os"
	"strings"

	"github.com/jipiboily/dickspatch/integrations"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

// Intercom integration
type Intercom struct {
	*intercom.Client
	Service
}

// Service defines the interface for working with the Intercom API
type Service interface {
	FindByUserID(userID string) (intercom.User, error)
	Save(user intercom.User) (intercom.User, error)
}

// API implements IntercomService
type API struct {
	*intercom.Client
}

// FindByUserID gets the user by UserID on Intercom
func (api API) FindByUserID(userID string) (user intercom.User, err error) {
	user, err = api.Client.Users.FindByUserID(userID)
	return
}

// Save the user on Intercom
func (api API) Save(user intercom.User) (savedUser intercom.User, err error) {
	savedUser, err = api.Client.Users.Save(&user)
	return
}

// Identify forwards and identify call to Intercom
func (i Intercom) Identify(event integrations.Event) (err error) {
	icUser, err := i.Service.FindByUserID(event.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not_found") {
			// The user doesn't exist, we just need to create it.
			icUser = intercom.User{UserID: event.UserID}
		} else {
			log.Println("Error fetching the Intercom user:", err)
			return
		}
	}

	icUser.CustomAttributes = event.UserTraits

	if event.UserTraits["email"] != nil {
		icUser.Email = event.UserTraits["email"].(string)
	}

	if event.UserTraits["name"] != nil {
		icUser.Name = event.UserTraits["name"].(string)
	}

	if event.UserTraits["createdAt"] != nil {
		// TODO: this is horrible, there must be a better way...
		icUser.CreatedAt = int64(event.UserTraits["createdAt"].(float64))
		icUser.SignedUpAt = int64(event.UserTraits["createdAt"].(float64))
	}

	savedUser, err := i.Service.Save(icUser)
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
	ic := Intercom{}
	ic.Client = intercom.NewClient(appID(), apiKey())
	ic.Service = API{ic.Client}
	integrations.RegisterIntegration("intercom", ic)
}
