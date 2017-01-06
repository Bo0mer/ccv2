package ccv2

// Application represents a Cloud Foundry application.
type Application struct {
	Metadata `json:"metadata"`

	Entity struct {
		Name               string `json:"name"`
		SpaceGUID          string `json:"space_guid"`
		StackGUID          string `json:"stack_guid"`
		Memory             int    `json:"memory"`
		Instances          int    `json:"instances"`
		DiskQuota          int    `json:"disk_quota"`
		State              string `json:"state"`
		Version            string `json:"version"`
		PackageState       string `json:"package_state"`
		HealthCheckType    string `json:"health_check_type"`
		HealthCheckTimeout int    `json:"health_check_timeout"`
		Buildpack          string `json:"buildpack"`
		Command            string `json:"command"`
		DetectedBuildpack  string `json:"detected_buildpack"`
		DetectedCommand    string `json:"detected_start_command"`
		Diego              bool   `json:"diego"`
		EnableSSH          bool   `json:"enable_ssh"`
	} `json:"entity"`
}

// ApplicationSummary represents summary about an application.
type ApplicationSummary struct {
	GUID               string `json:"guid"`
	Name               string `json:"name"`
	SpaceGUID          string `json:"space_guid"`
	StackGUID          string `json:"stack_guid"`
	Memory             int    `json:"memory"`
	Instances          int    `json:"instances"`
	DiskQuota          int    `json:"disk_quota"`
	State              string `json:"state"`
	Version            string `json:"version"`
	PackageState       string `json:"package_state"`
	HealthCheckType    string `json:"health_check_type"`
	HealthCheckTimeout int    `json:"health_check_timeout"`
	Buildpack          string `json:"buildpack"`
	Command            string `json:"command"`
	DetectedBuildpack  string `json:"detected_buildpack"`
	DetectedCommand    string `json:"detected_start_command"`
	Diego              bool   `json:"diego"`
	EnableSSH          bool   `json:"enable_ssh"`
	RunningInstances   int    `json:"running_instances"`
}
