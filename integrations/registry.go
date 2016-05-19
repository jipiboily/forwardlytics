package integrations

import (
	"sort"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
)

var integrationsMu sync.Mutex
var integrations = make(map[string]Integration)

// RegisterIntegration registers a integration so it can be created from its name. Integrations should
// call this from an init() function so that they registers themselvse on
// import
func RegisterIntegration(name string, integration Integration) {
	integrationsMu.Lock()
	defer integrationsMu.Unlock()
	if integration == nil {
		panic("integration: Register integration is nil")
	}
	if _, dup := integrations[name]; dup {
		panic("sql: Register called twice for integration " + name)
	}
	integrations[name] = integration
}

// GetIntegration retrieves a registered integration by name
func GetIntegration(name string) Integration {
	integrationsMu.Lock()
	defer integrationsMu.Unlock()
	integration := integrations[name]
	return integration
}

// IntegrationList returns a sorted list of the names of the registered integrations.
func IntegrationList() []string {
	integrationsMu.Lock()
	defer integrationsMu.Unlock()
	var list []string
	for name := range integrations {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// RemoveIntegration is ONLY used for tests, to be able remove a test integration
func RemoveIntegration(name string) {
	if !strings.HasPrefix(name, "test-only-integration-") {
		logrus.Error("THIS IS ONLY FOR TESTING!!")
	}
	integrationsMu.Lock()
	defer integrationsMu.Unlock()
	delete(integrations, name)
}
