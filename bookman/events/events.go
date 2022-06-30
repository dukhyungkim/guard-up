package events

type EventType string

func (t EventType) String() string {
	return string(t)
}

const (
	EventAddedBook   EventType = "ADDED_BOOK"
	EventUpdatedBook EventType = "UPDATED_BOOK"
	EventDeletedBook EventType = "DELETED_BOOK"

	EventAddedUser   EventType = "ADDED_USER"
	EventUpdatedUser EventType = "UPDATED_USER"
	EventDeletedUser EventType = "DELETED_USER"

	EventStartRental EventType = "START_RENTAL"
	EventEndRental   EventType = "END_RENTAL"
)

type EventSender func(eventType EventType, eventData any)
