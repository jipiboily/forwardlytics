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
	service := &FakeIntercomAPIWhenCreate{}
	ic.Service = service
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
		ReceivedAt: 3344,
	}

	err := ic.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if !service.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	expectedUser := intercom.User{
		UserID:     "123",
		Name:       "John Doe",
		Email:      "john@example.com",
		CreatedAt:  123,
		SignedUpAt: 123,
		CustomAttributes: map[string]interface{}{
			"name":                    "John Doe",
			"email":                   "john@example.com",
			"createdAt":               float64(123),
			"forwardlyticsReceivedAt": int64(3344),
		},
	}

	if !reflect.DeepEqual(service.ReceivedUser, expectedUser) {
		t.Errorf("Wrong user. Expected \n%#v\n but got \n%#v\n", expectedUser, service.ReceivedUser)
	}
}

func TestIdentifySuccessWhenUpdate(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	service := &FakeIntercomAPISuccess{}
	ic.Service = service
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
		ReceivedAt: 234,
	}

	err := ic.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if !service.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	expectedUser := intercom.User{
		Name:       "John Doe",
		Email:      "john@example.com",
		CreatedAt:  123,
		SignedUpAt: 123,
		CustomAttributes: map[string]interface{}{
			"name":                    "John Doe",
			"email":                   "john@example.com",
			"createdAt":               float64(123),
			"forwardlyticsReceivedAt": int64(234),
		},
	}

	if !reflect.DeepEqual(service.ReceivedUser, expectedUser) {
		t.Errorf("Wrong user. Expected \n%#v\n but got \n%#v\n", expectedUser, service.ReceivedUser)
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
	ic.Service = &FakeIntercomAPIWhenCreate{}
	es := &FakeIntercomEventsService{t: t}
	ic.EventRepository = es

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := ic.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	expectedEvent := &intercom.Event{
		UserID:    "123",
		EventName: "account.created",
		Metadata: map[string]interface{}{
			"email":                   "john@example.com",
			"forwardlyticsReceivedAt": int64(65),
		},
		Email:     "john@example.com",
		CreatedAt: 1234567,
	}

	if !reflect.DeepEqual(es.ReceivedEvent, expectedEvent) {
		es.t.Errorf("Wrong event. Expected \n%#v\n but got \n%#v\n", expectedEvent, es.ReceivedEvent)
	}
}

func TestTrackWhenFail(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	ic.Service = &FakeIntercomAPIWhenCreate{}
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
		herr, ok := err.(intercom.IntercomError)
		if !ok || herr.GetCode() != "not_found" {
			t.Fatal(err)
		}
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	if service.SaveCalled {
		t.Error("New Intercom user was created, and it should not")
	}
}

func TestTrackWhenAPropertyIsAMap(t *testing.T) {
	// It removes the properties that are map, as they are not supported
	// by Intercom. See the `Metadata support` section of
	// https://docs.intercom.io/the-intercom-platform/tracking-events-in-intercom
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	ic.Service = &FakeIntercomAPIWhenCreate{}
	es := &FakeIntercomEventsService{t: t}
	ic.EventRepository = es

	settings := map[string]interface{}{
		"metric_name": "rt:activeUsers",
	}

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email":    "john@example.com",
			"settings": settings,
		},
		Timestamp:  1234567,
		ReceivedAt: 33445,
	}

	err := ic.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	expectedEvent := &intercom.Event{
		UserID:    "123",
		EventName: "account.created",
		Metadata: map[string]interface{}{
			"email":                   "john@example.com",
			"forwardlyticsReceivedAt": int64(33445),
		},
		Email:     "john@example.com",
		CreatedAt: 1234567,
	}

	if !reflect.DeepEqual(es.ReceivedEvent, expectedEvent) {
		es.t.Errorf("Wrong event. Expected \n%#v\n but got \n%#v\n", expectedEvent, es.ReceivedEvent)
	}
}

func TestPage(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsService{t: t}
	ic.EventRepository = es

	event := integrations.Page{
		Name:   "Homepage",
		UserID: "123",
		Url:    "http://www.example.com/homepage",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := ic.Page(event)
	if err != nil {
		t.Fatal(err)
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	expectedPageEvent := &intercom.Event{
		UserID:    "123",
		EventName: "Page visited",
		Metadata: map[string]interface{}{
			"email": "john@example.com",
			"url":   "http://www.example.com/homepage",
			"forwardlyticsReceivedAt": int64(65),
			"forwardlyticsName":       "Homepage",
		},
		Email:     "john@example.com",
		CreatedAt: 1234567,
	}

	if !reflect.DeepEqual(es.ReceivedEvent, expectedPageEvent) {
		es.t.Errorf("Wrong page-event. Expected \n%#v\n but got \n%#v\n", expectedPageEvent, es.ReceivedEvent)
	}
}

func TestPageWhenFail(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsServiceFailing{t: t}
	ic.EventRepository = es

	event := integrations.Page{
		Name:   "Homepage",
		UserID: "123",
		Url:    "http://www.example.com/homepage",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := ic.Page(event)
	if err == nil {
		t.Fatal("Expecting an error")
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}
}

func TestPageWhenUserDoesNotExists(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsServiceFailingWithNotFound{t: t}
	ic.EventRepository = es
	service := &FakeIntercomAPISuccessCreateNewUserFromTrack{t: t}
	ic.Service = service

	event := integrations.Page{
		Name:   "Homepage",
		UserID: "123",
		Url:    "http://www.example.com/homepage",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := ic.Page(event)
	if err != nil {
		herr, ok := err.(intercom.IntercomError)
		if !ok || herr.GetCode() != "not_found" {
			t.Fatal(err)
		}
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	if service.SaveCalled {
		t.Error("New Intercom user was created, and it should not")
	}
}

func TestPageWhenAPropertyIsAMap(t *testing.T) {
	// It removes the properties that are map, as they are not supported
	// by Intercom. See the `Metadata support` section of
	// https://docs.intercom.io/the-intercom-platform/tracking-events-in-intercom
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	es := &FakeIntercomEventsService{t: t}
	ic.EventRepository = es

	settings := map[string]interface{}{
		"metric_name": "rt:activeUsers",
	}

	event := integrations.Page{
		Name:   "Homepage",
		UserID: "123",
		Url:    "http://www.example.com/homepage",
		Properties: map[string]interface{}{
			"email":    "john@example.com",
			"settings": settings,
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := ic.Page(event)
	if err != nil {
		t.Fatal(err)
	}

	if !es.SaveCalled {
		t.Error("Save was NOT called on the event")
	}

	expectedPageEvent := &intercom.Event{
		UserID:    "123",
		EventName: "Page visited",
		Metadata: map[string]interface{}{
			"email": "john@example.com",
			"url":   "http://www.example.com/homepage",
			"forwardlyticsReceivedAt": int64(65),
			"forwardlyticsName":       "Homepage",
		},
		Email:     "john@example.com",
		CreatedAt: 1234567,
	}

	if !reflect.DeepEqual(es.ReceivedEvent, expectedPageEvent) {
		es.t.Errorf("Wrong event. Expected \n%#v\n but got \n%#v\n", expectedPageEvent, es.ReceivedEvent)
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
	SaveCalled   bool
	ReceivedUser intercom.User
}

func (api FakeIntercomAPISuccess) FindByUserID(userID string) (user intercom.User, err error) {
	return
}

func (api *FakeIntercomAPISuccess) Save(user intercom.User) (savedUser intercom.User, err error) {
	api.SaveCalled = true
	api.ReceivedUser = user
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
	FakeIntercomAPISuccess
	SaveCalled   bool
	ReceivedUser intercom.User
}

func (api FakeIntercomAPIWhenCreate) FindByUserID(userID string) (user intercom.User, err error) {
	err = errors.New("404: not_found, User Not Found")
	return
}

func (api *FakeIntercomAPIWhenCreate) Save(user intercom.User) (savedUser intercom.User, err error) {
	api.SaveCalled = true
	api.ReceivedUser = user
	return
}

type FakeIntercomEventsService struct {
	t             *testing.T
	SaveCalled    bool
	ReceivedEvent *intercom.Event
}

func (es *FakeIntercomEventsService) Save(event *intercom.Event) error {
	es.SaveCalled = true
	es.ReceivedEvent = event
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
