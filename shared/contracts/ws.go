package contracts

// WSMessage pesan yang dikirim ke Frontend
type WSMessage struct {
	Type string `json:"type"` // e.g., "TODO_UPDATED"
	Data any    `json:"data"` // Payload objek Todo
}
