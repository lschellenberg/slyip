package slyerrors

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
)

func MapJetError(err error, msg string) error {
	orgError := err
	err = errors.Unwrap(err)
	if err == nil {
		err = orgError
	}
	msg = msg + ": " + err.Error()
	if IsNoRowsError(err) {
		return NotFound("", msg)
	}
	if err, ok := err.(*pq.Error); ok {
		fmt.Println(err.Code)
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	switch pqErr.Code {
	case "23505":
		return Conflict("23505", msg)
	default:
		return err
	}
}

func IsNoRowsError(err error) bool {
	return strings.Contains(err.Error(), "no rows in result set")
}
