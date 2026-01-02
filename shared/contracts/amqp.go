package contracts

// AmqpMessage structure standar untuk event bus
type AmqpMessage struct {
	RoomID string `json:"room_id"` // Siapa yang harus terima update ini? (atau GroupID)
	Type   string `json:"type"`    // Jenis Event
	Data   []byte `json:"data"`    // Payload JSON asli (Todo yang diupdate)
}

// Routing Keys (Events)
const (
	TodoEventCreated = "todo.event.created"
	TodoEventUpdated = "todo.event.updated"
	TodoEventDeleted = "todo.event.deleted"
)
