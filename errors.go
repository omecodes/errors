package errors

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrBadInput           = errors.New("bad input")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthorized       = errors.New("forbidden")
	ErrInternal           = errors.New("internal")
	ErrNotSupported       = errors.New("unsupported")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrUnImplemented      = errors.New("unimplemented")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
)

func New(message string) error {
	return Detailed(ErrInternal, With("message", message))
}

type jsonObject map[string]interface{}

type detailed struct {
	err     error
	details map[string]interface{}
}

func (d *detailed) Error() string {
	return d.err.Error()
}

func (d *detailed) Is(err error) bool {
	return errors.Is(d.err, err)
}

type Info func(err *detailed)

func With(key string, value interface{}) Info {
	return func(err *detailed) {
		if err.details == nil {
			err.details = make(map[string]interface{})
		}
		err.details[key] = value
	}
}

func Detailed(err error, info ...Info) error {
	d := &detailed{
		err: err,
	}
	for _, i := range info {
		i(d)
	}
	return d
}

func Write(w io.Writer, err error) (int, error) {
	var (
		data []byte
		mErr error
	)

	d, ok := err.(*detailed)
	if ok {
		data, mErr = json.Marshal(jsonObject{
			"code":    GetHttpStatusCode(d.err),
			"details": d.details,
		})
	} else {
		data, mErr = json.Marshal(jsonObject{
			"error": err.Error(),
		})
	}

	if mErr != nil {
		return 0, err
	}
	return w.Write(data)
}

func GetHttpStatusCode(err error) int {
	d, ok := err.(*detailed)
	if ok {
		return GetHttpStatusCode(d.err)
	}

	switch err {
	case ErrBadInput:
		return 400
	case ErrUnauthorized:
		return 401
	case ErrForbidden:
		return 403
	case ErrNotFound:
		return 404
	case ErrConflict:
		return 409
	case ErrUnImplemented:
		return 501
	case ErrServiceUnavailable:
		return 503
	case ErrNotSupported:
		return 505
	default:
		return 500
	}
}
