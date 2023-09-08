package response

// ErrorResponse ...
type ErrorResponse struct {
	Message string `json:"Message"`
	Code    int    `json:"Code"`
	Status  bool   `json:"Status"`
}
