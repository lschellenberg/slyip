package info

import (
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/httpx"
)

func AdminCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, err := verifier.GetPrincipal(r.Context())
		if err != nil {
			httpx.RespondWithJSON(w, httpx.Unauthorized(err.Error()))
			return
		}

		if principal.Role != verifier.RoleAdmin {
			httpx.RespondWithJSON(w, httpx.Forbidden("user has no admin"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
