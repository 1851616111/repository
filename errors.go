package main

import (
	"fmt"
)

type Error struct {
	code    uint
	message string
}

var (
	Errors            [NumErrors]*Error
	ErrorNone         *Error
	ErrorUnkown       *Error
	ErrorJsonBuilding *Error
)

const (
	OK = iota + 1000
	ErrorCodeUnkown
	ErrorCodeJsonBuilding

	ErrorCodeUrlNotSupported
	ErrorCodeDbNotInitlized
	ErrorCodeAuthFailed
	ErrorCodePermissionDenied
	ErrorCodeInvalidParameters
	ErrorCodeGetDataItem
	ErrorCodeCreateSubscription
	ErrorCodeGetSubscription
	ErrorCodeCancelSubscription
	ErrorCodeQuerySubscription
	ErrorCodeSubscriptionNotFound
	ErrorCodeCreateTransaction
	ErrorCodeGetTransaction
	ErrorCodeQueryTransaction
	ErrorCodeNoParameter
	ErrorCodeDataBase

	NumErrors
)

func init() {
	initError(OK, "OK")
	initError(ErrorCodeUnkown, "unknown error")
	initError(ErrorCodeJsonBuilding, "json building error")

	initError(ErrorCodeUrlNotSupported, "unsupported url")
	initError(ErrorCodeDbNotInitlized, "db is not inited")
	initError(ErrorCodeAuthFailed, "auth failed")
	initError(ErrorCodePermissionDenied, "permission denied")
	initError(ErrorCodeInvalidParameters, "invalid parameters")
	initError(ErrorCodeNoParameter, "no parameter")
	initError(ErrorCodeGetDataItem, "failed to get data item")

	ErrorNone = E(OK)
	ErrorUnkown = E(ErrorCodeUnkown)
	ErrorJsonBuilding = E(ErrorCodeJsonBuilding)
}

func initError(code uint, message string) {
	if code < NumErrors {
		Errors[code] = newError(code, message)
	}
}

func E(code uint) *Error {
	if code > NumErrors {
		return Errors[ErrorCodeUnkown]
	}

	return Errors[code]
}

func GetError2(code uint, message string) *Error {
	e := E(code)
	if e == nil {
		return newError(code, message)
	} else {
		return newError(code, fmt.Sprintf("%s (%s)", e.message, message))
	}
}

func newError(code uint, message string) *Error {
	return &Error{code: code, message: message}
}

func newUnknownError(message string) *Error {
	return &Error{
		code:    ErrorCodeUnkown,
		message: message,
	}
}

func ErrInvalidParameter(paramName string) *Error {
	return &Error{
		code:    ErrorCodeInvalidParameters,
		message: fmt.Sprint("%s: %s", E(ErrorCodeInvalidParameters).message, paramName),
	}
}

func ErrNoParameter(paramName string) *Error {
	return &Error{
		code:    ErrorCodeNoParameter,
		message: fmt.Sprint("%s : %s", E(ErrorCodeNoParameter).message, paramName),
	}
}

func ErrDataBase(e error) *Error {
	return &Error{
		code:    ErrorCodeInvalidParameters,
		message: fmt.Sprint("%s: %s", E(ErrorCodeInvalidParameters).message, e.Error()),
	}
}
