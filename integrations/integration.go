package integrations

// Integration defines what an integration is made of.
// Each integrations is responsible to register it self to the registry (see
// RegisterIntegration for details).
type Integration interface {
	// Identify is responsible of forwarding the identify call to the integration
	Identify(identification Identification) error

	// Enabled returns wether or not the integration is enabled/configured
	Enabled() bool
}

// Identification defines the structure of the data we receive from the API
type Identification struct {
	// Unique user ID. Should not change, ever.
	UserID string `json:"userID"`
	// Set of custom traits sent to the integrations. Some might be required, on
	// a per integration basis.
	UserTraits map[string]interface{} `json:"userTraits"`
	// Timestamp of when the identifiaction originally triggered
	Timestamp int64 `json:"timestamp"`
	// Timestamp of when Forwardlytics received the identifiaction.
	ReceivedAt int64 `json:"receivedAt"`
}

// Validate the content of the identifiaction to be sure it has everything that's needed
func (i Identification) Validate() (missingParameters []string) {
	if i.UserID == "" {
		missingParameters = append(missingParameters, "userID")
	}

	if i.Timestamp == 0 {
		missingParameters = append(missingParameters, "timestamp")
	}
	return
}
