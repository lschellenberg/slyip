package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"yip/src/api/services/dto"
	"yip/src/httpx"
)

const KeyProfileContext = "profile"

type ProfileMiddleware struct {
	ProfileFetcher func(context.Context, uuid.UUID) (dto.SLYBase, error)
	ShowLog        bool
}

func NewProfileMiddleWare(fetcher func(context.Context, uuid.UUID) (dto.SLYBase, error)) ProfileMiddleware {
	return ProfileMiddleware{
		ProfileFetcher: fetcher,
		ShowLog:        true,
	}
}

func GetProfileFromCtx(r *http.Request) dto.SLYBase {
	return r.Context().Value(KeyProfileContext).(dto.SLYBase)
}

func (a ProfileMiddleware) ProfileCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accountId := chi.URLParam(r, "accountId"); accountId != "" {
			uu, err := uuid.Parse(accountId)
			if err != nil {
				httpx.RespondWithJSON(w, httpx.BadRequest("account id is no uuid"))
				return
			}
			profile, err := a.ProfileFetcher(r.Context(), uu)
			if err != nil {
				httpx.RespondWithJSON(w, httpx.MapServiceError(err))
				return
			}
			fmt.Println("Profile in middleware is", profile)
			ctx := context.WithValue(r.Context(), KeyProfileContext, profile)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httpx.RespondWithError(w, 400, "no accountId provided", "")
		}
	})
}
