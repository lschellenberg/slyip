package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"yip/src/httpx"
	"yip/src/slyerrors"
)

type EntityMiddleware[T any] struct {
	fetcher    func(context.Context, string) (T, error)
	keyContext string
	path       string
}

func NewEntityMiddleWare[T any](fetcher func(context.Context, string) (T, error), keyContext string, path string) EntityMiddleware[T] {
	return EntityMiddleware[T]{
		fetcher:    fetcher,
		keyContext: keyContext,
		path:       path,
	}
}

func (a EntityMiddleware[T]) EntityFromCtx(r *http.Request) T {
	return r.Context().Value(a.keyContext).(T)
}

func (a EntityMiddleware[T]) EntityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pathValue := chi.URLParam(r, a.path); pathValue != "" {
			entity, err := a.fetcher(r.Context(), pathValue)
			if err != nil {
				if slyerrors.IsNoRowsError(err) {
					httpx.RespondWithJSON(w, httpx.NotFound(fmt.Sprintf("could not find entity: %s for context %s", pathValue, a.keyContext)))
				}
				httpx.RespondWithJSON(w, httpx.MapServiceError(err))
				return
			}
			ctx := context.WithValue(r.Context(), a.keyContext, entity)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httpx.RespondWithError(w, 400, "no path value detected", "")
		}
	})
}
