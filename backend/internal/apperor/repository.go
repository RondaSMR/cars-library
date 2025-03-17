package apperor

import "errors"

var (
	ErrNoEffect     = errors.New("no effect")
	ErrRepoNotFound = errors.New("not found")
)
