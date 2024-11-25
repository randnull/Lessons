package errors

import "errors"

// Ошибки доступа
var InvalidToken = errors.New("Invalid token")
var NotRoots = errors.New("No roots")
