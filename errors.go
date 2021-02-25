package errors

import (
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
)

const (
	CodeInternal = iota
	CodeBadRequest
	CodeUnauthorized
	CodeForbidden
	CodeNotFound
	CodeConflict
	CodeUnImplemented
	CodeServiceUnavailable
	CodeVersionNotSupported
	CodeUnReferencedID
	CodeNotSupported
)

type Details struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type Error struct {
	Code    int       `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
	Details []Details `json:"details,omitempty"`
}

func (e *Error) Error() string {
	encoded, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("{\"code\":%d, \"message\":\"%s\"}", e.Code, e.Message)
	}
	return string(encoded)
}

func (e *Error) AddDetails(key string, value interface{}) {
	e.Details = append(e.Details, Details{Key: key, Value: value})
}

func (e *Error) HTTPStatus() int {
	switch e.Code {
	case CodeBadRequest:
		return 400
	case CodeUnauthorized:
		return 401
	case CodeForbidden:
		return 403
	case CodeNotFound:
		return 404
	case CodeConflict:
		return 409
	case CodeUnImplemented:
		return 501
	case CodeServiceUnavailable:
		return 503
	case CodeVersionNotSupported:
		return 505
	default:
		return 500
	}
}

func HTTPStatus(err error) int {
	e, _ := Parse(err.Error())
	return e.HTTPStatus()
}

func Parse(str string) (Error, bool) {
	var e Error
	err := json.Unmarshal([]byte(str), &e)
	return e, err == nil && e.Code > 0 && e.Message != ""
}

func IsNotFound(e error) bool {
	if err, ok := e.(*Error); ok {
		return err.Code == CodeNotFound
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)
	return err.Code == CodeNotFound
}

func IsServiceUnavailable(e error) bool {
	if err, ok := e.(*Error); ok {
		return err.Code == CodeServiceUnavailable
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == CodeServiceUnavailable
}

func IsConflict(e error) bool {
	if isPrimaryKeyConstraintError(e) {
		return true
	}

	if err, ok := e.(*Error); ok {
		return err.Code == CodeConflict
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == CodeConflict
}

func IsNotReferencedID(e error) bool {
	if isForeignKeyConstraintError(e) {
		return true
	}

	if err, ok := e.(*Error); ok {
		return err.Code == CodeUnReferencedID
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == CodeUnReferencedID
}

func IsUnauthorized(e error) bool {
	if err, ok := e.(*Error); ok {
		return err.Code == CodeUnauthorized
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == CodeUnauthorized
}

func IsForbidden(e error) bool {

	if err, ok := e.(*Error); ok {
		return err.Code == CodeForbidden
	}

	var err Error
	_ = json.Unmarshal([]byte(e.Error()), &err)

	return err.Code == CodeForbidden
}

func Internal(message string, details ...Details) *Error {
	return &Error{Code: CodeInternal, Message: message, Details: details}
}
func BadRequest(message string, details ...Details) *Error {
	return &Error{Code: CodeBadRequest, Message: message, Details: details}
}
func Unauthorized(message string, details ...Details) *Error {
	return &Error{Code: CodeUnauthorized, Message: message, Details: details}
}
func Forbidden(message string, details ...Details) *Error {
	return &Error{Code: CodeForbidden, Message: message, Details: details}
}
func NotFound(message string, details ...Details) *Error {
	return &Error{Code: CodeNotFound, Message: message, Details: details}
}
func Conflict(message string, details ...Details) *Error {
	return &Error{Code: CodeConflict, Message: message, Details: details}
}
func UnImplemented(message string, details ...Details) *Error {
	return &Error{Code: CodeUnImplemented, Message: message, Details: details}
}
func ServiceUnavailable(message string, details ...Details) *Error {
	return &Error{Code: CodeServiceUnavailable, Message: message, Details: details}
}
func VersionNotSupported(message string, details ...Details) *Error {
	return &Error{Code: CodeVersionNotSupported, Message: message, Details: details}
}
func UnReferencedID(message string, details ...Details) *Error {
	return &Error{Code: CodeUnReferencedID, Message: message, Details: details}
}
func Unsupported(message string, details ...Details) *Error {
	return &Error{Code: CodeNotSupported, Message: message, Details: details}
}

func isPrimaryKeyConstraintError(err error) bool {
	if me, ok := err.(*mysql.MySQLError); ok {
		return me.Number == 1062

	} else if se, ok := err.(sqlite3.Error); ok {
		return se.ExtendedCode == 2067 || se.ExtendedCode == 1555
	}
	return false
}

func isForeignKeyConstraintError(err error) bool {
	if me, ok := err.(*mysql.MySQLError); ok {
		return me.Number == 1216

	} else if se, ok := err.(sqlite3.Error); ok {
		return se.ExtendedCode == sqlite3.ErrConstraintForeignKey
	}
	return false
}