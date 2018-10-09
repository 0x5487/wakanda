package types

import "fmt"

var (
	ErrDuplicateEntry = AppError{ErrorCode: "duplicate_entry", Message: "duplicate entry"}
)

type AppError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	//Details   map[string]interface{} `json:"details"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s - %s", e.ErrorCode, e.Message)
}

func IsErrDuplicate(err error) bool {
	return err == ErrDuplicateEntry
}
