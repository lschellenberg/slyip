package user

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"yip/src/httpx"
	"yip/src/repositories"
)

const KeyAccountContext = "account"

func getAccountFromCtx(r *http.Request) repositories.UserAccount {
	return r.Context().Value(KeyAccountContext).(repositories.UserAccount)
}

func (a Controller) AccountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accountId := chi.URLParam(r, "accountId"); accountId != "" {
			uu, err := uuid.Parse(accountId)
			if err != nil {
				httpx.RespondWithJSON(w, httpx.BadRequest("account id is no uuid"))
				return
			}
			account, err := a.service.GetAccountById(r.Context(), uu)
			if err != nil {
				httpx.RespondWithJSON(w, httpx.MapServiceError(err))
				return
			}
			fmt.Println("Account in middleware is", account)
			ctx := context.WithValue(r.Context(), KeyAccountContext, account)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httpx.RespondWithError(w, 400, "no accountId provided", "")
		}
	})
}
