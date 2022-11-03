package model

import "fmt"

// ErrorResponse модель ошибки HTTP ответа
type ErrorResponse struct {
	Error string `json:"error"`
}

func (err ErrorResponse) String() string {
	return fmt.Sprintf("{\"error\": %q}", err.Error)
}
