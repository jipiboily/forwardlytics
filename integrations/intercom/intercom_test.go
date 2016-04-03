package intercom

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/jipiboily/forwardlytics/integrations"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

func TestIdentifySuccessWhenCreate(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	ic.Service = FakeIntercomAPIWhenCreate{t: t}
	event := integrations.Event{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
	}
	err := ic.Identify(event)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIdentifySuccessWhenUpdate(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	ic.Service = FakeIntercomAPISuccess{t: t}
	event := integrations.Event{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": float64(123),
		},
	}
	err := ic.Identify(event)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIdentifyWhenFail(t *testing.T) {
	ic := Intercom{}
	ic.Client = intercom.NewClient("", "")
	ic.Service = FakeIntercomAPIFailSave{}
	err := ic.Identify(integrations.Event{})
	if err == nil {
		t.Fatal("Expecting an error.")
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
	t *testing.T
}

func (api FakeIntercomAPISuccess) FindByUserID(userID string) (user intercom.User, err error) {
	return
}

func (api FakeIntercomAPISuccess) Save(user intercom.User) (savedUser intercom.User, err error) {
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
}

func (api FakeIntercomAPIFailSave) Save(user intercom.User) (savedUser intercom.User, err error) {
	err = errors.New("Some API error")
	return
}

type FakeIntercomAPIWhenCreate struct {
	t *testing.T
	FakeIntercomAPISuccess
}

func (api FakeIntercomAPIWhenCreate) FindByUserID(userID string) (user intercom.User, err error) {
	err = errors.New("404: not_found, User Not Found")
	return
}

func (api FakeIntercomAPIWhenCreate) Save(user intercom.User) (savedUser intercom.User, err error) {
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
