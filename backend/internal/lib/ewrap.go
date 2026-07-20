package lib

import (
	"fmt"
)

func Ewrap(prefix string, err error) error {
	return fmt.Errorf("%s -> %w", prefix, err)
}
