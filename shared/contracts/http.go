package contracts

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

// APIError is the error structure for the API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
