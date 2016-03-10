package main

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    uint   `json:"code"`
	Message string `json:"msg"`
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
	ErrorCodeItemOutOfLimit
	ErrorCodeRepOutOfLimit
	ErrorCodeItemPriceOutOfLimit
	ErrorCodeRepExistCooperateItem
	ErrorCodeNoLogin
	ErrorCodeRepExistDataitem

	ErrorCodeNoParameter = 1400

	NumErrors = 1401
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
	initError(ErrorCodeQueryDBNotFound, "%s '%s' not found.")
	initError(ErrorCodeOptFile, "file operation")
	initError(ErrorCodeItemOutOfLimit, "dataitem out of limit 50")
	initError(ErrorCodeRepOutOfLimit, "repository out of limit")
	initError(ErrorCodeItemPriceOutOfLimit, "dataitem price out of limit 6")
	initError(ErrorCodeRepExistCooperateItem, "repository exists cooperate dataitem")
	initError(ErrorCodeNoLogin, "no login")
	initError(ErrorCodeRepExistDataitem, "repository exists dataitem, can not delete")

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
		return newError(code, fmt.Sprintf("%s (%s)", e.Message, message))
	}
}

func newError(code uint, message string) *Error {
	return &Error{Code: code, Message: message}
}

func newUnknownError(message string) *Error {
	return &Error{
		Code:    ErrorCodeUnkown,
		Message: message,
	}
}

func ErrInvalidParameter(paramName string) *Error {
	return &Error{
		Code:    ErrorCodeInvalidParameters,
		Message: fmt.Sprintf("%s: %s", E(ErrorCodeInvalidParameters).Message, paramName),
	}
}

func ErrParseJson(e error) *Error {
	return &Error{
		Code:    ErrorCodeJsonBuilding,
		Message: fmt.Sprintf("%s: %s", E(ErrorCodeJsonBuilding).Message, e.Error()),
	}
}

func ErrNoParameter(paramName string) *Error {
	return &Error{
		Code:    ErrorCodeNoParameter,
		Message: fmt.Sprintf("%s: %s", E(ErrorCodeNoParameter).Message, paramName),
	}
}

func ErrRepositoryNotFound(name string) *Error {
	return &Error{
		Code:    ErrorCodeQueryDBNotFound,
		Message: fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, "Repository", name),
	}
}

func ErrDataitemNotFound(name string) *Error {
	return &Error{
		Code:    ErrorCodeQueryDBNotFound,
		Message: fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, "Dataitem", name),
	}
}

func ErrTagNotFound(name string) *Error {
	return &Error{
		Code:    ErrorCodeQueryDBNotFound,
		Message: fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, "Tag", name),
	}
}

func ErrFieldNotFound(field, name string) *Error {
	return &Error{
		Code:    ErrorCodeQueryDBNotFound,
		Message: fmt.Sprintf(E(ErrorCodeQueryDBNotFound).Message, field, name),
	}
}

func ErrDataBase(e error) *Error {
	if e == nil {
		return E(OK)
	} else {
		return &Error{
			Code:    ErrorCodeItemOutOfLimit,
			Message: fmt.Sprintf("%s : %s", E(ErrorCodeDataBase).Message, e.Error()),
		}
	}
}

func ErrFile(e error) *Error {
	return &Error{
		Code:    ErrorCodeOptFile,
		Message: fmt.Sprintf("%s : %s", E(ErrorCodeOptFile).Message, e.Error()),
	}
}

func (e *Error) ErrToString() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func ErrRepOutOfLimit(num int) *Error {
	return &Error{
		Code:    ErrorCodeRepOutOfLimit,
		Message: fmt.Sprintf("%s : %d", E(ErrorCodeRepOutOfLimit).Message, num),
	}
}

func ErrRepExistCooperateItem(repname, itemname string) *Error {
	return &Error{
		Code:    ErrorCodeRepExistCooperateItem,
		Message: fmt.Sprintf("repname=%s exist cooperate item %s", repname, itemname),
	}
}
