package verifier

import (
	"context"
	"net/http"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

type TokenVerifierMiddleware struct {
	Verify func(token string) (*Principal, error)
}

func NewTokenVerifierMiddleware(Verify func(token string) (*Principal, error)) TokenVerifierMiddleware {
	return TokenVerifierMiddleware{Verify: Verify}
}

type ctxPrincipalType int

const ctxPrincipal = ctxPrincipalType(0) // quicker and safer by ensuring that we never run into any kind of collision

func GetPrincipal(ctx context.Context) (Principal, error) {
	if t := ctx.Value(ctxPrincipal); t == nil {
		return Principal{}, slyerrors.Unexpectedf("token not in context")
	} else {
		return t.(Principal), nil
	}
}

func (v TokenVerifierMiddleware) PrincipalCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, err := ParseAuthorizationBearer(r)
		if err != nil {
			httpx.RespondWithJSON(w, httpx.Unauthorized(err.Error()))
			return
		}
		principal, err := v.Verify(rawToken)
		if err != nil {
			e := slyerrors.Cause(err)
			switch e.Kind {
			case slyerrors.KindBadRequest:
				httpx.RespondWithJSON(w, httpx.BadRequest(e.Details))
			case slyerrors.KindUnauthorized:
				httpx.RespondWithJSON(w, httpx.Unauthorized(e.Details))
			}

			return
		}

		ctx := context.WithValue(r.Context(), ctxPrincipal, *principal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	httpx.RespondWithError(w, http.StatusUnauthorized, msg, msg)
}
