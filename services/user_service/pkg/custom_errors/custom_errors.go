package custom_errors

import "errors"

var UserNotFound = errors.New("user not found")
var ErrorWithCreate = errors.New("cannot create user")
var ErrorIncorrectRole = errors.New("incorrect role")
var ErrorAfterRowScan = errors.New("error after rows scan")
var ErrorNotFound = errors.New("not found")
var ErrorCountLessZero = errors.New("count less than zero")
var ErrorServiceError = errors.New("server error")
