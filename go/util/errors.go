package util

import "fmt"

type CustomError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
