package token

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/api/services/dto"
	"yip/src/common"
	"yip/src/httpx"
)

type Controller struct {
	tokenService       *services.TokenService
	userService        *services.UserService
	yipAdminMiddleware verifier.TokenVerifierMiddleware
}

func NewController(
	service *services.TokenService,
	userService *services.UserService,
	tokenMiddleware *verifier.TokenVerifierMiddleware) Controller {
	return Controller{
		tokenService:       service,
		userService:        userService,
		yipAdminMiddleware: *tokenMiddleware,
	}
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/verify", c.VerifyToken)
		r.Post("/refresh", c.RefreshToken)

		r.Group(func(r chi.Router) {
			r.Use(c.yipAdminMiddleware.PrincipalCtx)
			r.Get("/userinfo", c.UserInfo)
		})
	}
}

// swagger:route POST /auth/token/refresh Token refresh
// Refreshes A Token
//
// # When Token expires get refresh it with refresh token
//
// Responses:
//
//	200: Token
func (a Controller) RefreshToken(w http.ResponseWriter, r *http.Request) {
	data := &dto.RefreshTokenRequestDTO{}

	if err := common.ReadAndValidate(data, r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	result, err := a.tokenService.RefreshToken(data.RefreshToken)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
	}

	httpx.RespondWithJSON(w, httpx.OK(result))
}

// swagger:parameters verifyToken
type verifyToken struct {
	// in:body
	Body dto.VerifyTokenRequestDTO
}

// swagger:route POST /auth/token/verify Token verifyToken
// Verifies Token
//
// Responses:
//
//	200: Principal
func (a Controller) VerifyToken(w http.ResponseWriter, r *http.Request) {
	data := &dto.VerifyTokenRequestDTO{}

	if err := common.ReadAndValidate(data, r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	token, err := a.tokenService.VerifyToken(data.Token)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
	}

	httpx.RespondWithJSON(w, httpx.OK(token))
}

// swagger:route GET /auth/token/userinfo Token
// returns the all userinfos
//
// Responses:
//
//	200: User
func (a Controller) UserInfo(w http.ResponseWriter, r *http.Request) {
	principal, err := verifier.GetPrincipal(r.Context())
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}

	uu, err := uuid.Parse(principal.ID)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.BadRequest(err.Error()))
		return
	}

	account, err := a.userService.GetSLYWalletsByAccountId(r.Context(), uu)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(account))
}
