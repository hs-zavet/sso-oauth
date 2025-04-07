package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/sso-oauth/internal/api/handlers"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/roles"
	"github.com/sirupsen/logrus"
)

type Api struct {
	cfg    config.Config
	router *chi.Mux
	log    *logrus.Logger
}

func NewAPI(cfg config.Config, log *logrus.Logger) Api {
	return Api{
		cfg:    cfg,
		router: chi.NewRouter(),
		log:    log,
	}
}

func (a *Api) Run(ctx context.Context, app *app.App) {
	auth := tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	admin := tokens.AccessGrant(a.cfg.JWT.AccessToken.SecretKey, roles.Admin, roles.SuperUser)

	h := handlers.NewHandlers(app, a.cfg, a.log)

	a.router.Route("/re-news/sso", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {
				r.Post("/refresh", h.Refresh)
				r.Route("/google", func(r chi.Router) {
					r.Get("/login", h.GoogleLogin)
					r.Get("/callback", h.GoogleCallback)
				})

				r.Route("/account", func(r chi.Router) {
					r.Use(auth)
					r.Route("/sessions", func(r chi.Router) {
						r.Route("/{session_id}", func(r chi.Router) {
							r.Get("/", h.SessionGet)
							r.Delete("/", h.SessionDelete)
						})
						r.Get("/", h.SessionsGet)
						r.Delete("/", h.SessionsTerminate)
					})
					r.Get("/", h.AccountGet)
					r.Delete("/", h.Logout)
				})
			})

			r.Route("/private", func(r chi.Router) {
				r.Route("/account", func(r chi.Router) {
					r.Use(admin)
					r.Route("/{account_id}", func(r chi.Router) {
						r.Route("/sessions", func(r chi.Router) {
							r.Get("/{session_id}", h.AdminSessionGet)
							r.Get("/", h.AdminSessionsGet)
							r.Delete("/", h.AdminSessionsTerminate)
						})
						r.Get("/", h.AdminAccountGet)
						r.Patch("/{role}", h.AdminRoleUpdate)
					})
				})
			})

			r.Post("/test/login", h.LoginSimple)
		})
	})

	server := httpkit.StartServer(ctx, a.cfg.Server.Port, a.router, a.log)

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, a.log)
}
