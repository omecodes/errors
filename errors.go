package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	Internal            = 0
	BadRequest          = 1
	Unauthorized        = 2
	Forbidden           = 3
	NotFound            = 4
	Conflict            = 5
	NotImplemented      = 6
	ServiceUnavailable  = 7
	VersionNotSupported = 8
	DuplicateResource   = 9
)

type Details map[string]string

type Info struct {
	Name    string
	Details string
}

type Error struct {
	Code    int     `json:"code,omitempty"`
	Message string  `json:"message,omitempty"`
	Details Details `json:"details,omitempty"`
}

func (e Error) Error() string {
	var entries []string
	if e.Details != nil && len(e.Details) > 0 {
		for name, details := range e.Details {
			entries = append(entries, fmt.Sprintf("\"%s\": \"%s\"", name, details))
		}
		return fmt.Sprintf("{\"code\":%d, \"message\":\"%s\", \"details\": {%s}}", e.Code, e.Message, strings.Join(entries, ","))
	}
	return fmt.Sprintf("{\"code\":%d, \"message\":\"%s\"}", e.Code, e.Message)
}

func (e Error) SetDetails(details Details) {
	e.Details = details
}

func (e Error) AddDetails(name, value string) {
	if e.Details == nil {
		e.Details = Details{}
	}
	e.Details[name] = value
}

func (e Error) HTTPStatus() int {
	switch e.Code {
	case BadRequest:
		return 400
	case Unauthorized:
		return 401
	case Forbidden:
		return 403
	case NotFound:
		return 404
	case Conflict:
		return 409
	case NotImplemented:
		return 501
	case ServiceUnavailable:
		return 503
	case VersionNotSupported:
		return 505
	default:
		return 500
	}
}

func Create(code int, message string, info ...Info) *Error {
	e := &Error{
		Code:    code,
		Message: message,
		Details: Details{},
	}

	for _, i := range info {
		e.Details[i.Name] = i.Details
	}
	return e
}

func New(msg string) Error {
	return Error{Message: msg}
}

func IsNotFound(e error) bool {
	if err, ok := e.(Error); ok {
		return err.Code == NotFound || err.Code == ServiceUnavailable
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == NotFound || err.Code == ServiceUnavailable
}

func IsDuplicate(e error) bool {
	if err, ok := e.(Error); ok {
		return err.Code == DuplicateResource
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == DuplicateResource
}

func IsPermissionDenied(e error) bool {
	if err, ok := e.(Error); ok {
		return err.Code == Forbidden || err.Code == Unauthorized
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == Forbidden || err.Code == Unauthorized
}
