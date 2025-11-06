package model

import (
	"fmt"
)

type PartNotFoundError struct {
	PartUUID string
}

func (e *PartNotFoundError) Error() string {
	return fmt.Sprintf("part with UUID %q not found", e.PartUUID)
}
