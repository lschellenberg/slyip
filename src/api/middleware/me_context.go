package middleware

import (
	"context"
	"fmt"
	"net/http"
	"yip/src/api/auth/verifier"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

type MeMiddleware[T any] struct {
	fetcher    func(context.Context, string) (T, error)
	keyContext string
}

func NewMeMiddleWare[T any](fetcher func(context.Context, string) (T, error), keyContext string) MeMiddleware[T] {
	return MeMiddleware[T]{
		fetcher:    fetcher,
		keyContext: keyContext,
	}
}

func (a MeMiddleware[T]) EntityFromCtx(r *http.Request) T {
	return r.Context().Value(a.keyContext).(T)
}

func (a MeMiddleware[T]) EntityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := verifier.GetPrincipal(r.Context())
		if err != nil {
			httpx.RespondWithError(w, 401, "unauthenticated", err.Error())
		}

		entity, err := a.fetcher(r.Context(), p.ID)
		if err != nil {
			if slyerrors.IsNoRowsError(err) {
				httpx.RespondWithJSON(w, httpx.NotFound(fmt.Sprintf("could not find profile slywallet: for context %s", a.keyContext)))
			}
			httpx.RespondWithJSON(w, httpx.MapServiceError(err))
			return
		}
		ctx := context.WithValue(r.Context(), a.keyContext, entity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
