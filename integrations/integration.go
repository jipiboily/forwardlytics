package integrations

// Integration defines what an integration is made of.
// Each integrations is responsible to register it self to the registry (see
// RegisterIntegration for details).
type Integration interface {
	// Identify is responsible of forwarding the identify call to the integration
	Identify(user User) error

	// Enabled returns wether or not the integration is enabled/configured
	Enabled() bool
}

// User defines the structure of the data we receive from the API
type User struct {
	ID     string            `json:"user_id"`
	Traits map[string]string `json:"traits"`
}

// New will return the properly instanciated integration, or an error
func New(name string) (integration Integration, err error) {
	return
}
