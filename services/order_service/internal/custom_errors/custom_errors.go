package custom_errors

import "errors"

var ErrorGetUser = errors.New("cannot get user")
var ErrorCreateUser = errors.New("cannot create user")
