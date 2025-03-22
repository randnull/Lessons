package custom_errors

import "errors"

var ErrorGetUser = errors.New("cannot get user")
var ErrorCreateUser = errors.New("cannot create user")
var ErrStudentByOrderNotFound = errors.New("student`s id not found")
var ErrGetOrder = errors.New("cannot get order")
var ErrResponseAlredyExist = errors.New("response already exist")
var ErrGetResponse = errors.New("cannot get response")
var ErrNotAllowed = errors.New("not allowed")
var ErrorInvalidToken = errors.New("invalid token")
var ErrorAlreadySetTutor = errors.New("tutor is already set on this order")
