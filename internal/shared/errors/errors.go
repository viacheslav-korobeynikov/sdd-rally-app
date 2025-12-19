package shared_errors

import "errors"

var ErrUsernameTaken = errors.New("username already taken")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrAccountLocked = errors.New("account is locked")
