package apperr

import "errors"

var (
	NotUnique = errors.New("not unique row")
	NoRow     = errors.New("no row")

	ShortString = errors.New("short string")

	TokenExpired = errors.New("token expired")
)
