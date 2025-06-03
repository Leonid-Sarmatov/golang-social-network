package messages

/* Основа ответа */
type BaseResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"ErrorMessage,omitempty"`
}