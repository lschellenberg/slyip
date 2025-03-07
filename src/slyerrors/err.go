package slyerrors

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Kind can be used to convert from lperr.Error to an HTTP error code for example.
//
// Do not remove the comments: https://pkg.go.dev/golang.org/x/tools/cmd/stringer
type Kind int

var idxKindMap map[string]Kind

func init() {
	idxKindMap = make(map[string]Kind)
	for intKind := 0; intKind < len(_Kind_index)-1; intKind++ {
		idxKindMap[Kind(intKind).String()] = Kind(intKind)
	}
}

func (i Kind) MarshalJSON() ([]byte, error) {
	buff := &bytes.Buffer{}
	buff.WriteRune('"')
	buff.WriteString(i.String())
	buff.WriteRune('"')
	return buff.Bytes(), nil
}

func (i *Kind) UnmarshalJSON(in []byte) error {
	if len(in) < 3 {
		return fmt.Errorf("invalid error kind %s", string(in))
	}
	skind := string(in)[1:][:len(in)-2]
	kind, exists := idxKindMap[skind]
	if !exists {
		return fmt.Errorf("unknown error kind %s", skind)
	}
	*i = kind
	return nil
}

const (
	KindUnknown             Kind = iota // Unknown
	KindUnexpected                      // Unexpected
	KindUnauthorized                    // Unauthorized
	KindForbidden                       // Forbidden
	KindValidation                      // Validation
	KindBadRequest                      // BadRequest
	KindNotFound                        // NotFound
	KindConflict                        // Conflict
	KindUnprocessableEntity             // UnprocessabeEntity
)

// Error can be used as error or *Error, and supports JSON encoding.
type Error struct {
	Kind       Kind                 `json:"kind"`
	Code       string               `json:"code"`
	Message    string               `json:"message"`                 // contains the root message only
	Details    string               `json:"details"`                 // contains the whole chain
	Validation map[string][2]string `json:"invalidFields,omitempty"` // key is a field name, value is a validation code and a message
}

// ValidationCodes are a middleware but basic collection of codes that you can use in most scenarios.
// For "too long/short/large/small" and "out of range" codes, you may want to append something to the code for a
// specific translation.
const (
	ValidationCodeCannotValidate               = "cannotValidate"  // ex. used to handle and error during validation where no precise information on the failure can be used.
	ValidationCodeUnexpectedValue              = "unexpectedValue" // ex. supports "DESC" and "ASC" but got "DERP"
	ValidationCodeInvalidUUID                  = "invalidUUID"
	ValidationCodeStringEmpty                  = "stringEmpty"
	ValidationCodeNotEthAddress                = "notEthAddress"
	ValidationCodeListEmpty                    = "listEmpty"
	ValidationCodeStringNotInList              = "stringNotInList"
	ValidationCodeStringTooLong                = "stringTooLong"
	ValidationCodeStringTooShort               = "stringTooShort"
	ValidationCodeNumberOutOfRange             = "numberOutOfRange" // ex. integer expected to be between X and Y
	ValidationCodeNumberTooSmall               = "numberTooSmall"
	ValidationCodeNumberTooLarge               = "numberTooLarge"
	ValidationCodeMalformedBearerAuthorization = "malformedBearerAuthorization"
)

// AddValidation will panic if not of class KindValidation.
// field is the actual field of the request that is not valid.
// code is a middleware code to be used by clients to translate to a user-friendly message.
//
// Example:
// AddValidation("shipmentId", "invalid_uuid")
func (e *Error) AddValidation(field, code, msg string, args ...interface{}) *Error {
	if e.Kind != KindValidation {
		panic("error is not of class Validation")
	}

	if e.Validation == nil {
		e.Validation = make(map[string][2]string)
	}

	e.Validation[field] = [2]string{code, fmt.Sprintf(msg, args...)}

	return e
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("error: kind %s (%d), code: %s, message: %s", e.Kind.String(), e.Kind, e.Code, e.Message)
	if e.Kind == KindValidation {
		chunks := make([]string, 0, len(e.Validation))
		for k, v := range e.Validation {
			chunks = append(chunks, fmt.Sprintf("[%s: code: %s | msg: %s]", k, v[0], v[1]))
		}

		sort.Strings(chunks) // stable output

		if len(chunks) > 0 {
			msg += ", invalid fields: " + strings.Join(chunks, "|")
		}
	}
	return msg
}

// New *Error with the values. args is used with fmt.Sprintf(msg, args...) to format the message.
// A code must be a value that can be used for translating to a user-friendly message for API consumers.
// The message is a directly human readable message that may contain details about the error.
// The kind of error is exposed but should only be used to convert to other slyerrors, such as HTTP codes.
func New(kind Kind, code string, msg string, details string, args ...interface{}) *Error {
	details = fmt.Sprintf(details, args...)
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: msg,
		Details: details,
	}
}

// Wrap err if not nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}

// Cause returns the first Error found in an error chain created by fmt.Errorf("... %w ...", err)
// If it cannot find an Error in the chain, it will return un Unknown error class with the last wrapped error
// as message, and "unknown" code.
//
// Passing a nil error will panic.
func Cause(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}

	details := err.Error()

	for {
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			err = unwrapped
		} else {
			break
		}
	}

	cerr := New(KindUnknown, ErrCodeUnknown, "Unknown", err.Error())
	cerr.Details = details
	return cerr
}

// Creation helpers

func Unexpected(code string, msg string, args ...interface{}) *Error {
	return New(KindUnexpected, code, "unexpected", msg, args...)
}

func Unauthorized(code string, msg string, args ...interface{}) *Error {
	return New(KindUnauthorized, code, "unauthorised", msg, args...)
}

func Forbidden(code string, msg string, args ...interface{}) *Error {
	return New(KindForbidden, code, "forbidden", msg, args...)
}

func ValidationErr(code string, msg string, args ...interface{}) *Error {
	e := New(KindValidation, code, "Validation", msg, args...)
	e.Validation = make(map[string][2]string)
	return e
}

func BadRequest(code string, details string, args ...interface{}) *Error {
	return New(KindBadRequest, code, "bad request", details, args...)
}

func NotFound(code string, msg string, args ...interface{}) *Error {
	return New(KindNotFound, code, "not found", msg, args...)
}

func Conflict(code string, msg string, args ...interface{}) *Error {
	return New(KindConflict, code, "conflict", msg, args...)
}
