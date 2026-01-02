package contracts

// AmqpMessage
type AmqpMessage struct {
	RoomID string `json:"room_id"`
	Type   string `json:"type"`
	Data   []byte `json:"data"`
}

// Routing Keys (Events)
const (
	TodoEventCreated = "todo.event.created"
	TodoEventUpdated = "todo.event.updated"
	TodoEventDeleted = "todo.event.deleted"
)
