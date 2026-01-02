package contracts

// WSMessage
type WSMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
