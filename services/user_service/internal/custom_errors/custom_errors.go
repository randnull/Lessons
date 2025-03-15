package custom_errors

import "errors"

// Ошибки доступа

// Ошибки запроса
var UserNotFound = errors.New("user not found")
var ErrorWithCreate = errors.New("cannot create user")
var ErrorIncorrectRole = errors.New("incorrect role")
