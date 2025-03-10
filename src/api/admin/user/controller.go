package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"yip/src/api/auth/pin"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/api/services/dto"
	"yip/src/common"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

type Controller struct {
	service            *services.UserService
	pinService         *pin.Service
	yipAdminMiddleware *verifier.TokenVerifierMiddleware
}

func NewController(service *services.UserService, pinService *pin.Service, tokenMiddleware *verifier.TokenVerifierMiddleware) Controller {
	return Controller{
		service:            service,
		pinService:         pinService,
		yipAdminMiddleware: tokenMiddleware,
	}
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/token", c.SignInUser)

		r.Group(func(r chi.Router) {
			r.Use(c.yipAdminMiddleware.PrincipalCtx)
			r.Get("/", c.GetUsers)
			r.Put("/role", c.SetRole)
			r.Post("/register", c.RegisterUser)
			r.Put("/email", c.SetEmail)
			r.Get("/pins", c.GetPins)

			r.Route("/{accountId}", func(r chi.Router) {
				r.Use(c.AccountCtx)
				r.Get("/", c.GetAccount)
			})
		})
	}
}

// swagger:parameters register
type register struct {
	// in:body
	Body dto.RegisterRequest
}

// swagger:route POST /admin/accounts/register OIDC register
// Registers A User
//
// Registers a user by providing email credentials
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: SignedMessage
func (c Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	data := &dto.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	validationError := slyerrors.NewValidation("400").
		ValidateNotEmpty("email", data.Email).
		ValidateNotEmpty("password", data.Password).
		Error()

	if validationError != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(validationError))
		return
	}

	result, err := c.service.RegisterUser(r.Context(), data)

	if err != nil {
		fmt.Println(err)
		fmt.Println(slyerrors.MapJetError(err, ""))
		fmt.Println(httpx.MapServiceError(slyerrors.MapJetError(err, "")))
		httpx.RespondWithJSON(w, httpx.MapServiceError(slyerrors.MapJetError(err, "")))
		return
	}

	fmt.Println(result)

	httpx.RespondWithJSON(w, httpx.Created(result))
}

// swagger:parameters signin
type signin struct {
	// in:body
	Body dto.SignInRequest
}

// swagger:route POST /admin/accounts/token OIDC signin
// Signs In A User
//
// Responds with token when email credentials are given
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: Token
func (c Controller) SignInUser(w http.ResponseWriter, r *http.Request) {
	data := &dto.SignInRequest{}
	if err := data.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	isYIPAdminRequest := false
	for _, a := range data.Audiences {
		if a == c.service.Config.JWT.Issuer {
			isYIPAdminRequest = true
		}
	}

	// trying to sign in for YIP API
	if isYIPAdminRequest {
		result, err := c.service.SignInYIPAdmin(r.Context(), data)

		if err != nil {
			httpx.RespondWithJSON(w, httpx.MapServiceError(err))
			return
		}

		httpx.RespondWithJSON(w, httpx.OK(result))
		return
	}

	// trying to sign in for Resource API -> check for correct audience
	audiences := c.service.Config.Audiences
	for _, a := range data.Audiences {
		isInside := false
		for _, configuredAudiences := range audiences {
			if configuredAudiences.URL == a {
				isInside = true
			}
		}
		if !isInside {
			httpx.RespondWithJSON(w, httpx.BadRequest(fmt.Sprintf("%s not found. Wrong audience", a)))
			return
		}
	}

	result, err := c.service.SignInUser(r.Context(), data)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(slyerrors.MapJetError(err, "")))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(result))
}

// swagger:route PUT /admin/accounts/role setrole
// Sets A Role To A Given User
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: Token
func (c Controller) SetRole(w http.ResponseWriter, r *http.Request) {
	data := &dto.SetRoleRequest{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	uu, err := uuid.Parse(data.UserID)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest("account id is no uuid"))
		return
	}

	validationError := slyerrors.NewValidation("400").
		ValidateNotEmpty("userId", data.UserID).
		//ValidateInList("role", data.Role, []string{auth.RoleAdmin, auth.RoleBasic}).
		Error()

	if validationError != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(validationError))
		return
	}

	user, err := verifier.GetPrincipal(r.Context())

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	if !user.IsAdmin() {
		http.Error(w, slyerrors.Forbidden("403", "access forbidden").Error(), 403)
		return
	}

	result, err := c.service.SetRole(r.Context(), uu, data.Role)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(slyerrors.MapJetError(err, "")))
	}

	httpx.RespondWithJSON(w, httpx.OK(result))
}

// swagger:route GET /admin/accounts users
// Returns all uses
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: []User
func (c Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetUsers(r.Context(), common.PaginationQuery{
		QueryFilters: common.QueryFilters{},
		PageSize:     1000,
		Offset:       0,
		SortBy:       "",
		SortOrder:    "",
	})

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(slyerrors.MapJetError(err, "")))
	}

	httpx.RespondWithJSON(w, httpx.OK(result))
}

// swagger:route GET /admin/accounts/{accountId} users
// Return specific account
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: []User
func (c Controller) GetAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Controller#GetAccount", getAccountFromCtx(r))
	httpx.RespondWithJSON(w, httpx.OK(getAccountFromCtx(r)))
}

// swagger:parameters setemail
type setemail struct {
	// in:body
	Body dto.SetEmailRequest
}

// swagger:route PUT /admin/accounts/email admin setemail
// Sets An Email Of A Given User
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: UserAccount
func (c Controller) SetEmail(w http.ResponseWriter, r *http.Request) {
	data := &dto.SetEmailRequest{}
	if err := data.ReadAndValidate(r); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	uu, err := uuid.Parse(data.UserID)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest("account id is no uuid"))
		return
	}

	user, err := verifier.GetPrincipal(r.Context())

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	if !user.IsAdmin() {
		http.Error(w, slyerrors.Forbidden("403", "access forbidden").Error(), 403)
		return
	}

	result, err := c.service.SetEmail(r.Context(), uu, data.Email)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(slyerrors.MapJetError(err, "")))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(result))
}

// swagger:route GET /admin/accounts/pins users
// Returns all pins
//
// Security:
//   - Bearer: []
//
// Responses:
//
//	200: []Pin
func (c Controller) GetPins(w http.ResponseWriter, r *http.Request) {
	pins, err := c.pinService.ListPins(r.Context())
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(pins))
}
