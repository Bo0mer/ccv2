package ccv2

// Metadata represents metadata for a resource.
type Metadata struct {
	GUID      string `json:"guid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Space represents a Cloud Foundry space.
type Space struct {
	Metadata `json:"metadata"`

	Entity struct {
		Name                    string `json:"name"`
		OrganizationGUID        string `json:"organization_guid"`
		SpaceQuotaDefinitonGUID string `json:"space_quota_definiton_guid"`
		AllowSSH                bool   `json:"allow_ssh"`
	} `json:"entity"`
}
