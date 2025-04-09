package participants

import "github.com/google/uuid"

// EventParticipant describes an event participant entity.
type EventParticipant struct {
	EventID uuid.UUID
	UserID  uuid.UUID
}
