package integrations

// Integration defines what an integration is made of.
// Each integrations is responsible to register it self to the registry (see
// RegisterIntegration for details).
type Integration interface {
	// Identify is responsible of forwarding the identify call to the integration
	Identify(identification Identification) error

	// Track forwards the event to the integration
	Track(event Event) error

	// Page forwards the page-view to the integration
	Page(page Page) error

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

// Event defines the structure for the incoming event data from the API
type Event struct {
	// Name is the name of the event
	Name string `json:"name"`

	// Unique user ID. Should not change, ever.
	UserID string `json:"userID"`

	// Properties are custom variables you can send with the event
	Properties map[string]interface{} `json:"properties"`

	// Timestamp of when the identifiaction originally triggered
	Timestamp int64 `json:"timestamp"`

	// ReceivedAt of when Forwardlytics received the identifiaction.
	ReceivedAt int64 `json:"receivedAt"`
}

// Validate the content of the event to be sure it has everything that's needed
func (e Event) Validate() (missingParameters []string) {
	if e.Name == "" {
		missingParameters = append(missingParameters, "name")
	}

	if e.UserID == "" {
		missingParameters = append(missingParameters, "userID")
	}

	if e.Timestamp == 0 {
		missingParameters = append(missingParameters, "timestamp")
	}
	return
}

// Page defines the structure for the incoming page-view data
type Page struct {
	// Name is the name of the page
	Name string `json:"name"`

	// Unique user ID. Should not change, ever.
	UserID string `json:"userID"`

	// Unique user ID. Should not change, ever.
	Url string `json:"url"`

	// Properties are custom variables you can send with the page
	Properties map[string]interface{} `json:"properties"`

	// Timestamp of when the page-call originally triggered
	Timestamp int64 `json:"timestamp"`

	// ReceivedAt of when Forwardlytics received the page-call.
	ReceivedAt int64 `json:"receivedAt"`
}

// Validate the content of the page to be sure it has everything that's needed
func (p Page) Validate() (missingParameters []string) {
	if p.Name == "" {
		missingParameters = append(missingParameters, "name")
	}

	if p.Url == "" {
		missingParameters = append(missingParameters, "url")
	}

	if p.UserID == "" {
		missingParameters = append(missingParameters, "userID")
	}

	if p.Timestamp == 0 {
		missingParameters = append(missingParameters, "timestamp")
	}
	return
}
