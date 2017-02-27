package intercom

import (
	"os"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
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
			logrus.WithError(err).WithField("identification", identification).Error("Error fetching the Intercom user")
			return
		}
	}

	if identification.UserTraits != nil {
		icUser.CustomAttributes = identification.UserTraits
	} else {
		icUser.CustomAttributes = make(map[string]interface{})
	}
	icUser.CustomAttributes["forwardlyticsReceivedAt"] = identification.ReceivedAt

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
		logrus.WithField("savedUser", savedUser).Info("User saved on Intercom")
	} else {
		logrus.WithError(err).WithField("identification", identification).WithField("icUser", icUser).Error("Error while saving on Intercom")
	}
	return
}

// Track forwards the event to Intercom
func (i Intercom) Track(event integrations.Event) (err error) {
	icEvent := intercom.Event{}
	icEvent.UserID = event.UserID
	icEvent.EventName = event.Name
	icEvent.CreatedAt = event.Timestamp

	// It removes the properties that are map, as they are not supported
	// by Intercom. See the `Metadata support` section of
	// https://docs.intercom.io/the-intercom-platform/tracking-events-in-intercom
	metaData := make(map[string]interface{})
	for k, v := range event.Properties {
		if !(reflect.TypeOf(event.Properties[k]).Kind() == reflect.Map) {
			metaData[k] = v
		}
	}

	metaData["forwardlyticsReceivedAt"] = event.ReceivedAt

	icEvent.Metadata = metaData

	// User needs to exist before we send event
	userAutoCreated := false
	icUser, err := i.Service.FindByUserID(event.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not_found") {
			// The user doesn't exist, we just need to create it.
			userAutoCreated = true
			icUser = intercom.User{UserID: icEvent.UserID}
		} else {
			logrus.WithError(err).WithField("event", event).Error("Error fetching the Intercom user")
		}
	}

	if event.Properties["email"] != nil {
		icEvent.Email = event.Properties["email"].(string)
		icUser.Email = icEvent.Email
	}

	// If we auto created the user, save it now
	if userAutoCreated {
		savedUser, err := i.Service.Save(icUser)
		if err == nil {
			logrus.WithField("savedUser", savedUser).Info("User doesn't exist and was auto-created on Intercom")
		} else {
			logrus.WithError(err).WithField("event", event).WithField("icUser", icUser).Error("Error while auto-creating user on Intercom")
		}
	}

	err = i.EventRepository.Save(&icEvent)

	if err != nil {
		logrus.WithError(err).WithField("event", event).WithField("icEvent", icEvent).Error("Error while saving event on Intercom")
	}

	return
}

// Page tracks page views. Intercom doesn't support a special type for
// this, so it's implemented as a special type of event
func (p Intercom) Page(page integrations.Page) (err error) {
	icPage := intercom.Event{}
	icPage.UserID = page.UserID
	icPage.EventName = "Page visited"
	icPage.CreatedAt = page.Timestamp

	// It removes the properties that are map, as they are not supported
	// by Intercom. See the `Metadata support` section of
	// https://docs.intercom.io/the-intercom-platform/tracking-events-in-intercom
	metaData := make(map[string]interface{})
	for k, v := range page.Properties {
		if !(reflect.TypeOf(page.Properties[k]).Kind() == reflect.Map) {
			metaData[k] = v
		}
	}

	metaData["forwardlyticsReceivedAt"] = page.ReceivedAt
	metaData["forwardlyticsName"] = page.Name
	metaData["url"] = page.Url

	icPage.Metadata = metaData

	if page.Properties["email"] != nil {
		icPage.Email = page.Properties["email"].(string)
	}

	err = p.EventRepository.Save(&icPage)

	if err != nil {
		logrus.WithError(err).WithField("event", page).WithField("icPage", icPage).Error("Error while saving event on Intercom")
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
