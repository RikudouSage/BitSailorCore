package result

import "github.com/google/uuid"

type ProfileOrganization struct {
	ID  uuid.UUID `json:"id"`
	Key string    `json:"key"`
}
