package model

import (
	"fmt"
)

func ErrPartNotFound(uuid string) error {
	return fmt.Errorf("part %s not found", uuid)
}
