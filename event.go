package ccv2

import "time"

// Event represents a Cloud Foundry application event.
type Event struct {
	Metadata `json:"metadata"`

	Entity struct {
		Type             string    `json:"type"`
		Actor            string    `json:"actor"`
		ActorType        string    `json:"actor_type"`
		ActorName        string    `json:"actor_name"`
		Actee            string    `json:"actee"`
		ActeeType        string    `json:"actee_type"`
		ActeeName        string    `json:"actee_name"`
		Timestamp        time.Time `json:"timestamp"`
		SpaceGUID        string    `json:"space_guid"`
		OrganizationGUID string    `json:"organization_guid"`
	} `json:"entity"`
}
