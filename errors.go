package errors

import (
	"encoding/json"
	"io"
)

const (
	badInput           = "bad input"
	forbidden          = "forbidden"
	unauthorized       = "unauthorized"
	internal           = "internal"
	notSupported       = "unsupported"
	serviceUnavailable = "service unavailable"
	unImplemented      = "unimplemented"
	notFound           = "not found"
	conflict           = "conflict"
	timeout            = "timeout"
)

// D is a map that contains details
type D map[string]interface{}

func New(details ...D) error {
	d := &e{
		kind: internal,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

type e struct {
	kind    string
	details map[string]interface{}
}

func (d *e) Error() string {
	return d.kind
}

func Write(w io.Writer, err error) (int, error) {
	var (
		data []byte
		mErr error
	)

	d, ok := err.(*e)
	if ok {
		data, mErr = json.Marshal(D{
			"code":    GetHttpStatusCode(d),
			"details": d.details,
		})
	} else {
		data, mErr = json.Marshal(D{
			"error": err.Error(),
		})
	}

	if mErr != nil {
		return 0, err
	}
	return w.Write(data)
}

func GetHttpStatusCode(err error) int {
	d, ok := err.(*e)
	if ok {
		switch d.kind {
		case badInput:
			return 400
		case unauthorized:
			return 401
		case forbidden:
			return 403
		case notFound:
			return 404
		case conflict:
			return 409
		case unImplemented:
			return 501
		case serviceUnavailable:
			return 503
		case notSupported:
			return 505
		default:
			return 500
		}
	}
	return 500

}

func FromHttpStatus(status int, details D) error {
	switch status {
	case 400:
		return BadInput(details)
	case 401:
		return Unauthorized(details)
	case 403:
		return Forbidden(details)
	case 404:
		return NotFound(details)
	case 409:
		return Conflict(details)
	case 501:
		return Unimplemented(details)
	case 503:
		return ServiceUnavailable(details)
	case 505:
		return NotSupported(details)
	default:
		return New(details)
	}
}

func Timeout(details ...D) error {
	d := &e{
		kind: timeout,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func BadInput(details ...D) error {
	d := &e{
		kind: badInput,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func Forbidden(details ...D) error {
	d := &e{
		kind: forbidden,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func Unauthorized(details ...D) error {
	d := &e{
		kind: unauthorized,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func NotSupported(details ...D) error {
	d := &e{
		kind: notSupported,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func ServiceUnavailable(details ...D) error {
	d := &e{
		kind: serviceUnavailable,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func Unimplemented(details ...D) error {
	d := &e{
		kind: unImplemented,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func NotFound(details ...D) error {
	d := &e{
		kind: notFound,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func Conflict(details ...D) error {
	d := &e{
		kind: conflict,
	}
	for _, i := range details {
		for k, v := range i {
			d.details[k] = v
		}
	}
	return d
}

func IsTimeout(err error) bool {
	return timeout == err.Error()
}

func IsBadInput(err error) bool {
	return badInput == err.Error()
}

func IsForbidden(err error) bool {
	return forbidden == err.Error()
}

func IsUnauthorized(err error) bool {
	return unauthorized == err.Error()
}

func IsNotSupported(err error) bool {
	return notSupported == err.Error()
}

func IsServiceUnavailable(err error) bool {
	return serviceUnavailable == err.Error()
}

func IsUnimplemented(err error) bool {
	return unImplemented == err.Error()
}

func IsNotFound(err error) bool {
	return notFound == err.Error()
}

func IsConflict(err error) bool {
	return conflict == err.Error()
}
