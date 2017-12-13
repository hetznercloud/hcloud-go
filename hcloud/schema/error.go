package schema

// Error represents the schema of an error response.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse defines the schema of a response containing an error.
type ErrorResponse struct {
	Error Error `json:"error"`
}
