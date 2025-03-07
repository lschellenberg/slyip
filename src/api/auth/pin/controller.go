package pin

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"yip/src/httpx"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) Controller {
	return Controller{
		service: service,
	}
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", c.RequestPin)
		r.Post("/redeem", c.Redeem)

	}
}

// swagger:parameters requestPin
type requestPin struct {
	// in:body
	Body PinRequestDTO
}

// swagger:route POST /auth/pin Pin requestPin
// Requests a pin
//
// # Request a pin
//
// Responses:
//
//	200: ok
func (a Controller) RequestPin(w http.ResponseWriter, r *http.Request) {
	data := &PinRequestDTO{}

	if err := data.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	pin, err := a.service.RequestPin(r.Context(), data)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError(err.Error()))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(pin))
}

// swagger:parameters redeemPin
type redeemPin struct {
	// in:body
	Body PinRedeemDTO
}

// swagger:route POST /auth/pin/redeem Pin redeemPin
// Redeems a given pin
//
// Responses:
//
//	200: Principal
func (a Controller) Redeem(w http.ResponseWriter, r *http.Request) {
	data := &PinRedeemDTO{}

	if err := data.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}

	token, err := a.service.Redeem(r.Context(), data)

	if err != nil {
		fmt.Println(err.Error())
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(token))
}
