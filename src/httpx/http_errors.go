package httpx

import (
	"net/http"
	"yip/src/slyerrors"
)

func MapAuthError(err error) *Response {
	rootErr := err
	response := Response{
		Payload:    rootErr,
		StatusCode: 401,
	}

	return &response
}

func MapServiceError(err error) *Response {
	if response, ok := err.(*Response); ok {
		return response
	}

	rootErr := slyerrors.Cause(err)
	response := Response{
		Payload:    rootErr,
		StatusCode: 0,
	}

	switch rootErr.Kind {
	case slyerrors.KindBadRequest:
		response.StatusCode = http.StatusBadRequest

	case slyerrors.KindNotFound:
		response.StatusCode = http.StatusNotFound

	case slyerrors.KindForbidden:
		response.StatusCode = http.StatusForbidden

	case slyerrors.KindConflict:
		response.StatusCode = http.StatusConflict

	case slyerrors.KindValidation:
		response.StatusCode = http.StatusBadRequest

	case slyerrors.KindUnauthorized:
		response.StatusCode = http.StatusUnauthorized

		// TODO: handle 502 and 503 slyerrors with dedicated kinds

	default:
		response.StatusCode = http.StatusInternalServerError
	}

	return &response
}
