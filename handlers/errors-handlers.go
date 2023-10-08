package handlers

import "fmt"

type HTTPError struct {
	Status      int    `json:"status"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error - Status %d -%q -%q\n", h.Status, h.Code, h.Description)
}

const (
	ErrNotFound            = "NotFound"
	ErrInternalServerError = "InternalServerError"
	ErrBadRequest          = "BadRequest"
)
