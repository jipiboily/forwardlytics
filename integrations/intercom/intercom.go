package intercom

import (
	"log"
	"os"
	"strings"

	"github.com/jipiboily/forwardlytics/integrations"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

// Intercom integration
type Intercom struct {
	*intercom.Client
	Service
	EventRepository
}

// Identify forwards and identify call to Intercom
func (i Intercom) Identify(identification integrations.Identification) (err error) {
	icUser, err := i.Service.FindByUserID(identification.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not_found") {
			// The user doesn't exist, we just need to create it.
			icUser = intercom.User{UserID: identification.UserID}
		} else {
			log.Println("Error fetching the Intercom user:", err)
			return
		}
	}

	icUser.CustomAttributes = identification.UserTraits

	if identification.UserTraits["email"] != nil {
		icUser.Email = identification.UserTraits["email"].(string)
	}

	if identification.UserTraits["name"] != nil {
		icUser.Name = identification.UserTraits["name"].(string)
	}

	if identification.UserTraits["createdAt"] != nil {
		// TODO: this is horrible, there must be a better way...
		icUser.CreatedAt = int64(identification.UserTraits["createdAt"].(float64))
		icUser.SignedUpAt = int64(identification.UserTraits["createdAt"].(float64))
	}

	savedUser, err := i.Service.Save(icUser)
	if err == nil {
		log.Printf("User saved on Intercom: %#v\n", savedUser)
	} else {
		log.Println("Error while saving on Intercom:", err)
	}
	return
}

// Track forwards the event to Intercom
func (i Intercom) Track(event integrations.Event) (err error) {
	icEvent := intercom.Event{}
	icEvent.UserID = event.UserID
	icEvent.EventName = event.Name
	icEvent.CreatedAt = event.Timestamp
	icEvent.Metadata = event.Properties
	if event.Properties["email"] != nil {
		icEvent.Email = event.Properties["email"].(string)
	}

	err = i.EventRepository.Save(&icEvent)

	if herr, ok := err.(intercom.IntercomError); ok && herr.GetCode() == "not_found" {
		log.Println("User not found, we need to create it first, then try again", herr)
		icUser := intercom.User{UserID: icEvent.UserID}
		_, err = i.Service.Save(icUser)
		if err == nil {
			log.Printf("User created on Intercom as part of sending an event")
		} else {
			log.Println("Error while creating user on Intercom as part of sending an event:", err)
		}
	}

	if err != nil {
		log.Println("Error while saving event on Intercom:", err)
	}
	return
}

// Enabled returns wether or not the Intercom integration is enabled/configured
func (i Intercom) Enabled() bool {
	return apiKey() != "" && appID() != ""
}

// EventRepository defines the interface for tracking events on Intercom
type EventRepository interface {
	Save(event *intercom.Event) error
}

// EventService tracks services on Intercom
type EventService struct {
	*intercom.Client
}

// Save the event on Intercom
func (es EventService) Save(event *intercom.Event) error {
	return es.Client.Events.Save(event)
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
	ic.EventRepository = EventService{ic.Client}

	// Useful for debugging, keeping it around to avoid remembering how to use it
	// ic.Client.Option(intercom.TraceHTTP(true))

	integrations.RegisterIntegration("intercom", ic)
}
