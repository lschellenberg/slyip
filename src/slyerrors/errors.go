package slyerrors

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

const (
	MsgUnexpected        = "UnexpectedError"
	MsgInvalidDateFormat = "InvalidDateFormat"
)

func Unexpectedf(msg string, args ...interface{}) *Error {
	return Unexpected(MsgUnexpected, msg, args...)
}

func NoPrincipal() *Error {
	return Unexpectedf("not authenticated")
}

func HTTPReq(err error) *Error {
	return Unexpectedf("request: %v", err)
}

func HTTPDial(err error) *Error {
	return Unexpectedf("dial: %v", err)
}

func HTTPDecode(err error) *Error {
	return Unexpectedf("decode: %v", err)
}

func HTTPRead(err error) *Error {
	return Unexpectedf("read: %v", err)
}

type Validation struct {
	err  *Error
	code string
}

func NewValidation(code string) *Validation {
	return &Validation{
		code: code,
	}
}

func (v *Validation) Add(field, code, msg string, args ...interface{}) *Validation {
	if v.err == nil {
		v.err = ValidationErr(code, "invalid request")
	}

	v.err.AddValidation(field, code, msg, args...)
	return v
}

func (v *Validation) Merge(err error, sub string) *Validation {
	if err == nil {
		return v
	}

	cerr := Cause(err)
	if cerr.Kind != KindValidation {
		v.Add(sub, v.code, cerr.Error())
	} else {
		if sub != "" {
			sub = sub + "."
		}

		for f, cv := range cerr.Validation {
			v.Add(sub+f, cv[0], cv[1])
		}
	}

	return v
}

func (v *Validation) ValidateUUID(field, value string) *Validation {
	if _, err := uuid.Parse(value); err != nil {
		v.Add(field, ValidationCodeInvalidUUID, "")
	}
	return v
}

func (v *Validation) ValidateRange(field string, min, value, max int64) *Validation {
	if value < min && value > max {
		v.Add(field, ValidationCodeNumberOutOfRange, "equal or between %d and %d", min, max)
	}
	return v
}

func (v *Validation) ValidateNotEmpty(field, value string) *Validation {
	if len(value) == 0 {
		v.Add(field, ValidationCodeStringEmpty, "")
	}
	return v
}

func (v *Validation) ValidateAtLeastOneElement(field string, value []string) *Validation {
	if len(value) == 0 {
		v.Add(field, ValidationCodeStringEmpty, "")
	}
	return v
}

func (v *Validation) ValidateEthAddress(field, value string) *Validation {
	if !IsValidEthAddress(value) {
		v.Add(field, ValidationCodeNotEthAddress, "")
	}
	return v
}

func (v *Validation) ValidateInList(field, value string, list []string) *Validation {
	isInList := false
	for _, l := range list {
		if l == value {
			isInList = true
			break
		}
	}
	if !isInList {
		v.Add(field, ValidationCodeStringNotInList, "")
	}
	return v
}

func (v *Validation) ValidateDateLayout(field, value, layout string) *Validation {
	_, err := time.Parse(layout, value)
	if err != nil {
		v.Add(field, MsgInvalidDateFormat, err.Error())
	}
	return v
}

func (v *Validation) Error() error {
	if v.err == nil { // even if v.err is nil, returning it directly doesn't count as a nil pointer once it's casted into error interface.
		return nil
	}
	return v.err
}

func IsValidEthAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}
