package ccv2

// Organization represents a Cloud Foundry organization.
type Organization struct {
	Metadata `json:"metadata"`

	Entity struct {
		Name               string `json:"name"`
		BillingEnabled     bool   `json:"billing_enabled"`
		QuotaDefinitonGUID string `json:"quota_definiton_guid"`
		Status             string `json:"status"`
	} `json:"entity"`
}
