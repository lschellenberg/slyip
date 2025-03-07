package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Response is meant to be returned by a modified version
type Response struct {
	Payload     interface{}
	StatusCode  int
	contentType string
	header      map[string]string
}

func (r *Response) Error() string {
	return fmt.Sprintf("%d: %v", r.StatusCode, r.Payload)
}

// SetContentType should be used only when you're using wrappers streaming raw data.
// It returns the response itself for convenience in handlers.
func (r *Response) SetContentType(ct string) *Response {
	r.contentType = ct
	return r
}

func (r *Response) Pagination(count string, resume string) *Response {
	if r.header == nil {
		r.header = make(map[string]string, 0)
	}
	r.header["X-TotalCount"] = count
	r.header["X-PageResume"] = resume
	return r
}

func (r *Response) AddHeader(key string, value string) *Response {
	if r.header == nil {
		r.header = make(map[string]string, 0)
	}
	r.header[key] = value
	return r
}

func OK(content interface{}) *Response {
	return &Response{
		Payload:    content,
		StatusCode: http.StatusOK,
	}
}

func Created(content interface{}) *Response {
	return &Response{
		Payload:    content,
		StatusCode: http.StatusCreated,
	}
}

func BadRequest(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "bad request",
			"details": details,
		},
		StatusCode: http.StatusBadRequest,
	}
}

func NotFound(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "not found",
			"details": details,
		},
		StatusCode: http.StatusNotFound,
	}
}

func Forbidden(reason string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "forbidden",
			"details": reason,
		},
		StatusCode: http.StatusForbidden,
	}
}

func Unauthorized(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "unauthorized",
			"details": details,
		},
		StatusCode: http.StatusUnauthorized,
	}
}

func InternalError(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "internal server error",
			"details": details,
		},
		StatusCode: http.StatusInternalServerError,
	}
}

func NotImplemented(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "not implemented",
			"details": details,
		},
		StatusCode: http.StatusNotImplemented,
	}
}

func Conflict(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "conflict",
			"details": details,
		},
		StatusCode: http.StatusConflict,
	}
}

func ServiceUnavailable(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "service unavailable",
			"details": details,
		},
		StatusCode: http.StatusServiceUnavailable,
	}
}

func NoContent() *Response {
	return &Response{
		StatusCode: http.StatusNoContent,
	}
}

func UnprocessableEntity(details string) *Response {
	return &Response{
		Payload: map[string]string{
			"message": "unprocessable entity",
			"details": details,
		},
		StatusCode: http.StatusUnprocessableEntity,
	}
}

type Handler func(w http.ResponseWriter, r *http.Request) *Response

func JSON(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, handler(w, r))
	}
}

// Stream uses the content type defined by the user calling SetContentType().
// Response.Payload is expected to implement io.Reader.
// Data streaming is preformed with io.Copy().
//
// In case Response.Payload is not io.Reader, RespondWithJSON is called instead.
func Stream(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := handler(w, r)
		v, ok := resp.Payload.(io.Reader)
		if ok {
			w.Header().Set("content-type", resp.contentType)
			if resp.header != nil {
				for k, v := range resp.header {
					w.Header().Set(k, v)
				}
			}
			w.WriteHeader(resp.StatusCode)
			if _, err := io.Copy(w, v); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "encoding error", err.Error())
			}
		} else {
			RespondWithJSON(w, resp)
		}
	}
}

// RespondWithError will retrieve a ResponseWriter, code and a payload
// It will marshal the payload and set the Content-MediaType to json.
// If there are extra headers, it will also set the headers.
// If json encoding or response write fails, it panics. You should then handle
// this panic in your router's middleware.
func RespondWithJSON(w http.ResponseWriter, response *Response) {
	buff := bytes.Buffer{}

	if err := json.NewEncoder(&buff).Encode(response.Payload); err != nil {
		// either the struct has a broken MarshalJSON implementation, or the json encoder cannot do its job.
		// so, we don't call it again at all, as this is a real, unrecoverable 500
		http.Error(w, "encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	if response.header != nil {
		for k, v := range response.header {
			w.Header().Set(k, v)
		}
	}

	w.WriteHeader(response.StatusCode)

	written, err := w.Write(buff.Bytes())
	if err != nil {
		log.Println("error writing response: " + err.Error())
	} else if written != len(buff.Bytes()) {
		log.Println("did not write the same amount of data")
	}

}

// RespondWithError will retrieve a ResponseWriter, code, msg and details.
// It writes a json response to the ResponseWriter with:
//
//	{
//		"message": YOUR_MESSAGE
//		"details": YOUR_DETAILS
//	}
//
// Note that it is preferred to use helpers like RespondWithJSON(BadRequest(details))
func RespondWithError(w http.ResponseWriter, code int, msg string, details string) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	RespondWithJSON(w, &Response{
		Payload:    map[string]string{"message": msg, "details": details},
		StatusCode: code,
	})
}

func OKHealthCheck(w http.ResponseWriter, r *http.Request) *Response {
	return OK("ok")
}

type StatusResponse struct {
	StatusCode int  `json:"statusCode"`
	Success    bool `json:"success"`
}
