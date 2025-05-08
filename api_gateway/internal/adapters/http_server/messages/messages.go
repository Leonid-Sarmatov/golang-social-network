package messages

/* Основа ответа */
type BaseResponse struct {
	Status       string
	ErrorMessage string `json:"ErrorMessage,omitempty"`
}