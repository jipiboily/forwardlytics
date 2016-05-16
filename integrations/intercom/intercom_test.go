package intercom

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/jipiboily/forwardlytics/integrations"
	intercom "gopkg.in/intercom/intercom-go.v2"
	intercomInterfaces "gopkg.in/intercom/intercom-go.v2/interfaces"
)

func TestIdentifySuccessWhenCreate(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	service := &FakeIntercomAPIWhenCreate{t: t}
	ic.Service = service
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
	}

	err := ic.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if !service.SaveCalled {
		t.Error("Save was NOT called on the event")
	}
}

func TestIdentifySuccessWhenUpdate(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	service := &FakeIntercomAPISuccess{t: t}
	ic.Service = service
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
	}

	err := ic.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if !service.SaveCalled {
		t.Error("Save was NOT called on the event")
	}
}

func TestIdentifyWhenFail(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	service := &FakeIntercomAPIFailSave{}
	ic.Service = service
	err := ic.Identify(integrations.Identification{})
	if err == nil {
		t.Fatal("Expecting an error.")
	}

	if !service.SaveCalled {
		t.Error("Save was NOT called on the event")
	}
}

func TestTrack(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsServiceSuccess{t: t}
	ic.EventRepository = es

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp: 1234567,
	}

	err := ic.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}
}

func TestTrackWhenFail(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsServiceFailing{t: t}
	ic.EventRepository = es

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp: 1234567,
	}

	err := ic.Track(event)
	if err == nil {
		t.Fatal("Expecting an error")
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}
}

func TestTrackWhenUserDoesNotExists(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsServiceFailingWithNotFound{t: t}
	ic.EventRepository = es
	service := &FakeIntercomAPISuccessCreateNewUserFromTrack{t: t}
	ic.Service = service

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp: 1234567,
	}

	err := ic.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	if !service.SaveCalled {
		t.Error("New Intercom user was NOT created as expected")
	}
}

func TestEnabledWhenConfigured(t *testing.T) {
	if err := os.Setenv("INTERCOM_API_KEY", "ABC"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("INTERCOM_APP_ID", "XYZ"); err != nil {
		t.Fatal(err)
	}

	ic := Intercom{}
	if ic.Enabled() != true {
		t.Error("Wrong value. The intergraiton should be enabled.")
	}
}

func TestEnabledWhenNotConfigured(t *testing.T) {
	if err := os.Unsetenv("INTERCOM_API_KEY"); err != nil {
		t.Fatal(err)
	}
	if err := os.Unsetenv("INTERCOM_APP_ID"); err != nil {
		t.Fatal(err)
	}

	ic := Intercom{}
	if ic.Enabled() == true {
		t.Error("Wrong value. The intergraiton should NOT be enabled.")
	}
}

type FakeIntercomAPISuccess struct {
	t          *testing.T
	SaveCalled bool
}

func (api FakeIntercomAPISuccess) FindByUserID(userID string) (user intercom.User, err error) {
	return
}

func (api *FakeIntercomAPISuccess) Save(user intercom.User) (savedUser intercom.User, err error) {
	api.SaveCalled = true
	expectedUser := intercom.User{
		Name:       "John Doe",
		Email:      "john@example.com",
		CreatedAt:  123,
		SignedUpAt: 123,
		CustomAttributes: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
	}
	if !reflect.DeepEqual(user, expectedUser) {
		api.t.Errorf("Wrong user. Expected \n%#v\n but got \n%#v\n", expectedUser, user)
	}
	return
}

type FakeIntercomAPIFailSave struct {
	FakeIntercomAPISuccess
	SaveCalled bool
}

func (api *FakeIntercomAPIFailSave) Save(user intercom.User) (savedUser intercom.User, err error) {
	api.SaveCalled = true
	err = errors.New("Some API error")
	return
}

type FakeIntercomAPIWhenCreate struct {
	t *testing.T
	FakeIntercomAPISuccess
	SaveCalled bool
}

func (api FakeIntercomAPIWhenCreate) FindByUserID(userID string) (user intercom.User, err error) {
	err = errors.New("404: not_found, User Not Found")
	return
}

func (api *FakeIntercomAPIWhenCreate) Save(user intercom.User) (savedUser intercom.User, err error) {
	api.SaveCalled = true
	expectedUser := intercom.User{
		UserID:     "123",
		Name:       "John Doe",
		Email:      "john@example.com",
		CreatedAt:  123,
		SignedUpAt: 123,
		CustomAttributes: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
	}
	if !reflect.DeepEqual(user, expectedUser) {
		api.t.Errorf("Wrong user. Expected \n%#v\n but got \n%#v\n", expectedUser, user)
	}
	return
}

type FakeIntercomEventsServiceSuccess struct {
	t          *testing.T
	SaveCalled bool
}

func (es *FakeIntercomEventsServiceSuccess) Save(event *intercom.Event) error {
	es.SaveCalled = true
	expectedEvent := &intercom.Event{
		UserID:    "123",
		EventName: "account.created",
		Metadata: map[string]interface{}{
			"email": "john@example.com",
		},
		Email:     "john@example.com",
		CreatedAt: 1234567,
	}

	if !reflect.DeepEqual(event, expectedEvent) {
		es.t.Errorf("Wrong event. Expected \n%#v\n but got \n%#v\n", expectedEvent, event)
	}

	return nil
}

type FakeIntercomEventsServiceFailing struct {
	t          *testing.T
	SaveCalled bool
}

func (es *FakeIntercomEventsServiceFailing) Save(event *intercom.Event) error {
	es.SaveCalled = true
	return errors.New("some error")
}

type FakeIntercomEventsServiceFailingWithNotFound struct {
	t          *testing.T
	SaveCalled bool
}

func (es *FakeIntercomEventsServiceFailingWithNotFound) Save(event *intercom.Event) error {
	es.SaveCalled = true
	err := intercomInterfaces.HTTPError{
		StatusCode: 404,
		Code:       "not_found",
		Message:    "Some message for not found",
	}
	return err
}

type FakeIntercomAPISuccessCreateNewUserFromTrack struct {
	t          *testing.T
	SaveCalled bool
}

func (api FakeIntercomAPISuccessCreateNewUserFromTrack) FindByUserID(userID string) (user intercom.User, err error) {
	return
}

func (api *FakeIntercomAPISuccessCreateNewUserFromTrack) Save(user intercom.User) (savedUser intercom.User, err error) {
	api.SaveCalled = true
	expectedUser := intercom.User{
		UserID: "123",
	}
	if !reflect.DeepEqual(user, expectedUser) {
		api.t.Errorf("Wrong user. Expected \n%#v\n but got \n%#v\n", expectedUser, user)
	}
	return
}
