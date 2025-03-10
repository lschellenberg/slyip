package siwe

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/api/services/dto"
	"yip/src/config"
	"yip/src/httpx"
)

type Controller struct {
	service     *services.SIWEService
	userService *services.UserService
	config      *config.Config
}

func NewController(c *config.Config, service *services.SIWEService, userService *services.UserService) Controller {
	return Controller{
		service:     service,
		userService: userService,
		config:      c,
	}
}

func (c Controller) Routes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/challenge", c.Challenge)
		r.Get("/nonce", c.Nonce)
		r.Post("/verify", c.Verify)
		r.Post("/submit", c.Submit)
	}
}

// swagger:parameters siweChallenge
type siweChallenge struct {
	// in:body
	Body dto.ChallengeRequestDTO
}

// swagger:route POST /siwe/challenge SIWE siweChallenge
// Requests a SIWE challenge
//
// Requests a challenge by providing a domain
// Responses:
//
//	200: ChallengeResponse
func (a Controller) Challenge(w http.ResponseWriter, r *http.Request) {
	data := &dto.ChallengeRequestDTO{}
	if err := data.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	if a.config.EthConfig.Chain.ID != data.ChainId {
		httpx.RespondWithJSON(w, httpx.BadRequest(fmt.Sprintf("unknown chain: %s", data.ChainId)))
		return
	}

	m, err := a.service.Challenge(data)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	fmt.Println(m.Challenge, m.ChainId, m.Address, m.Domain)

	httpx.RespondWithJSON(w, httpx.OK(m))
}

// swagger:route GET /siwe/nonce SIWE
// Requests a SIWE nonce
//
// Requests a nonce
// Responses:
//
//	200: NonceResponse
func (a Controller) Nonce(w http.ResponseWriter, r *http.Request) {
	httpx.RespondWithJSON(w, httpx.OK(a.service.Nonce()))
}

// swagger:parameters siweSubmission
type siweSubmission struct {
	// in:body
	Body dto.SubmitRequestDTO
}

// swagger:route POST /siwe/verify SIWE siweSubmission
// Verifies a SIWE signature
//
// Verifies a signature by recovering and comparing to original message
// Responses:
//
//	200: ChallengeResponse
func (a Controller) Verify(w http.ResponseWriter, r *http.Request) {
	data := &dto.SubmitRequestDTO{}
	fmt.Println(data)
	if err := data.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	m, err := a.service.Verify(data)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(m))
}

// swagger:route POST /siwe/submit SIWE siweSubmission
// Requests a token by  a SIWE signature
//
// Verifies a signature by recovering and comparing to original message
// Responses:
//
//	200: ChallengeResponse
func (a Controller) Submit(w http.ResponseWriter, r *http.Request) {
	data := &dto.SubmitRequestDTO{}
	if err := data.ReadAndValidate(r); err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	m, err := a.service.Verify(data)

	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(err))
		return
	}

	// TODO check for EIP1271
	if m.RecoveredAddress != m.OriginalAddress {
		httpx.RespondWithJSON(w, httpx.BadRequest(fmt.Sprintf("recovered address not recognized [recovered: %s, original: %s]", m.RecoveredAddress, m.OriginalAddress)))
		return
	}

	auds := make([]string, len(a.config.Audiences))
	for k, v := range a.config.Audiences {
		auds[k] = v.URL
	}

	fmt.Println("GetOrCreateAccount")
	ecdsa, err := a.service.GetOrCreateAccount(r.Context(), m.OriginalAddress)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(httpx.MapAuthError(err)))
		return
	}

	uu, err := uuid.Parse(ecdsa.AccountId)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError(err.Error()))
		return
	}

	fmt.Println("GetAccountById")
	account, err := a.userService.GetAccountById(r.Context(), uu)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.InternalError(err.Error()))
		return
	}

	fmt.Println("CreateToken")
	token, err := a.service.CreateToken(auds, ecdsa.AccountId, ecdsa.Address, account.LastUsedSLYWallet, verifier.RoleBasic)
	if err != nil {
		httpx.RespondWithJSON(w, httpx.MapServiceError(httpx.MapAuthError(err)))
		return
	}

	httpx.RespondWithJSON(w, httpx.OK(token))
}
