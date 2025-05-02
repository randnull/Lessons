package custom_errors

import "errors"

// Ошибки доступа
var ErrInvalidToken = errors.New("Invalid token")
var ErrNotRoots = errors.New("No roots")
var ErrorGetUser = errors.New("cannot get user")
var ErrorCreateUser = errors.New("cannot create user")
var ErrorInvalidToken = errors.New("invalid token")
var ErrorInvalidRole = errors.New("invalid role")
