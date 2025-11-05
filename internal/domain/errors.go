package domain

import "errors"

var ErrForbidden = errors.New("forbidden: user does not have the required permissions")
