package main

import (
	"encoding/json"
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
	OK              = iota
	ErrorCodeUnkown = iota + 1000
	ErrorCodeJsonBuilding

	ErrorCodeUrlNotSupported
	ErrorCodeDbNotInitlized
	ErrorCodeUnauthorized
	ErrorCodePermissionDenied
	ErrorCodeInvalidParameters
	ErrorCodeDataBase
	ErrorCodeQueryDBNotFound
	ErrorCodeOptFile

	NumErrors

	ErrorCodeNoParameter = 1400
)

func init() {
	initError(OK, "OK")
	initError(ErrorCodeUnkown, "unknown error")
	initError(ErrorCodeJsonBuilding, "json building error")

	initError(ErrorCodeUrlNotSupported, "unsupported url")
	initError(ErrorCodeDbNotInitlized, "db is not inited")
	initError(ErrorCodeUnauthorized, "unauthorized")
	initError(ErrorCodePermissionDenied, "permission denied")
	initError(ErrorCodeInvalidParameters, "invalid parameters")
	initError(ErrorCodeNoParameter, "no parameter")
	initError(ErrorCodeDataBase, "database operate")
	initError(ErrorCodeQueryDBNotFound, "query %s no found")
	initError(ErrorCodeOptFile, "file operation")

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
		message: fmt.Sprintf("%s: %s", E(ErrorCodeInvalidParameters).message, paramName),
	}
}

func ErrParseJson(e error) *Error {
	return &Error{
		code:    ErrorCodeJsonBuilding,
		message: fmt.Sprintf("%s: %s", E(ErrorCodeJsonBuilding).message, e.Error()),
	}
}

func ErrNoParameter(paramName string) *Error {
	return &Error{
		code:    ErrorCodeNoParameter,
		message: fmt.Sprintf("%s: %s", E(ErrorCodeNoParameter).message, paramName),
	}
}

func ErrQueryNotFound(paramName string) *Error {
	return &Error{
		code:    ErrorCodeQueryDBNotFound,
		message: fmt.Sprintf(E(ErrorCodeQueryDBNotFound).message, paramName),
	}
}

func ErrDataBase(e error) *Error {
	if e == nil {
		return E(OK)
	} else {
		return &Error{
			code:    ErrorCodeDataBase,
			message: fmt.Sprintf("%s : %s", E(ErrorCodeDataBase).message, e.Error()),
		}
	}
}

func ErrFile(e error) *Error {
	return &Error{
		code:    ErrorCodeOptFile,
		message: fmt.Sprintf("%s : %s", E(ErrorCodeOptFile).message, e.Error()),
	}
}

func (e *Error) ErrToString() string {
	b, _ := json.Marshal(e)
	return string(b)
}
