package result

import (
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	RevisionDate time.Time `json:"revisionDate"`
}
