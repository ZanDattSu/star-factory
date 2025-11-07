package v1

import (
	"fmt"
	"strings"
)

type PartsNotFoundError struct {
	Code      int      `json:"code"`
	PartsUUID []string `json:"parts_uuid"`
}

func (e *PartsNotFoundError) Error() string {
	return fmt.Sprintf("parts with UUIDs [%s] not found", strings.Join(e.PartsUUID, ", "))
}

func NewPartsNotFoundError(uuids []string) *PartsNotFoundError {
	return &PartsNotFoundError{
		Code:      404,
		PartsUUID: uuids,
	}
}
