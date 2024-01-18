package common_errors

import "errors"

var ErrTooManyRequests error = errors.New("too many requests")

func Is(err error, target error) bool {
	return err.Error() == target.Error()
}
