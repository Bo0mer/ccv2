package ccv2

// Info represents a Cloud Controller info.
type Info struct {
	Name        string `json:"name"`
	Build       string `json:"build"`
	Support     string `json:"support"`
	Version     int    `json:"version"`
	Description string `json:"description"`

	MinCliVersion string `json:"min_cli_version"`
	APIVersion    string `json:"api_version"`

	AppSSHEndpoint        string `json:"app_ssh_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	RoutingEndpoint       string `json:"routing_endpoint"`
	LoggingEndpoint       string `json:"logging_endpoint"`
	DopplerEndpoint       string `json:"doppler_logging_endpoint"`
}
