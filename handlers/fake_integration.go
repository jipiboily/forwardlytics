package handlers

import "github.com/jipiboily/forwardlytics/integrations"

// FakeIntegration is the base of a fake integration, used for testing.
type FakeIntegration struct {
}

// Identify is responsible of forwarding the identify call to the integration
func (fi FakeIntegration) Identify(identification integrations.Identification) error {
	return nil
}

// Track forwards the event to the integration
func (fi FakeIntegration) Track(event integrations.Event) error {
	return nil
}

// Enabled returns wether or not the integration is enabled/configured
func (fi FakeIntegration) Page(page integrations.Page) error {
	return nil
}

// Enabled returns wether or not the integration is enabled/configured
func (fi FakeIntegration) Enabled() bool {
	return true
}
