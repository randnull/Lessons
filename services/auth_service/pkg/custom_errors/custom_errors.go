package custom_errors

import "errors"

var ErrorGetUser = errors.New("cannot get user")
var ErrorCreateUser = errors.New("cannot create user")
var ErrorInvalidToken = errors.New("invalid token")
var ErrorInvalidRole = errors.New("invalid role")
