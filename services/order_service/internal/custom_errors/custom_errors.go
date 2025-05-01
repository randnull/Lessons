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
var ErrorSetStatus = errors.New("cannot set status")
var ErrorServiceError = errors.New("error with service")
var ErrorSelectTutor = errors.New("error with select tutor")
var ErrorGetResponse = errors.New("error with get response")
var ErrorBadStatus = errors.New("error bad status")
var ErrorLowTimeFromResponse = errors.New("from response to review less than 3 days")
var ErrorNotActiveReview = errors.New("error not active review")
var ErrorNotFound = errors.New("not found")
