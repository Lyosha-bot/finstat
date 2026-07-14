package ewrap

import (
	"fmt"
)

func Wrap(prefix string, err error) error {
	return fmt.Errorf("%s -> %w", prefix, err)
}
