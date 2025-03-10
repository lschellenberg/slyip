package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"net/http"
	"yip/src/api/admin"
	"yip/src/api/auth"
	"yip/src/api/auth/verifier"
	"yip/src/api/services"
	"yip/src/api/slywallet"
	"yip/src/app"
)

const apiVersionURL = "/api/v1"

type Api struct {
	Router  *chi.Mux
	App     *app.App
	Modules *Modules
}

type Modules struct {
	AuthModule      auth.Module
	AdminModule     admin.AdminModule
	SLYWalletModule slywallet.Module
}

func NewApi(app *app.App) Api {
	api := Api{
		nil,
		app,
		&Modules{},
	}

	apiServices := services.GenerateApiServices(app)

	tokenMiddleware := initMiddleware(app.Verifier)

	api.Modules.AuthModule = auth.NewAuthModule(app.Config, &apiServices, &tokenMiddleware)
	api.Modules.AdminModule = admin.NewAdminModule(app.Config, &apiServices, &tokenMiddleware, app.EthProvider)
	api.Modules.SLYWalletModule = slywallet.NewModule(&apiServices, &tokenMiddleware)

	api.Router = newRouter(&api)
	return api
}

func newRouter(api *Api) *chi.Mux {
	r := chi.NewRouter()

	origins := make([]string, 0)

	for _, a := range api.App.Config.Audiences {
		origins = append(origins, a.Clients...)
	}

	// TODO make sure to delete this
	origins = []string{"*"}

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-MediaType", "X-CSRF-Auth"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	routes := chi.NewRouter()
	api.defineRoutes(routes)

	r.Mount(apiVersionURL, routes)

	oidcRoutes := chi.NewRouter()

	oidcRoutes.Route("/.well-known", func(r chi.Router) {
		r.Get("/jwks", api.WellKnownJWKS)
		r.Get("/openid-configuration", api.WellKnownConfiguration)
	})
	r.Mount("/", oidcRoutes)

	return r
}

func (api Api) defineRoutes(r *chi.Mux) {
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("works"))
	})

	r.Route("/auth", api.Modules.AuthModule.Routes())
	r.Route("/admin", api.Modules.AdminModule.Routes())
	r.Route("/sly", api.Modules.SLYWalletModule.Routes())
}

func initMiddleware(verf *verifier.Verifier) verifier.TokenVerifierMiddleware {
	v := func(token string) (*verifier.Principal, error) {
		return verf.VerifyToken(context.Background(), token)
	}

	return verifier.NewTokenVerifierMiddleware(v)
}
