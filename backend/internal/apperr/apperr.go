package apperr

import "errors"

var (
	NotUnique = errors.New("not unique row")
	NoRow     = errors.New("no row")
)
