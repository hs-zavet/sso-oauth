package cli

import (
	"context"
	"sync"

	"github.com/hs-zavet/sso-oauth/internal/api"
	"github.com/hs-zavet/sso-oauth/internal/app"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/sirupsen/logrus"
)

func runServices(ctx context.Context, cfg config.Config, log *logrus.Logger, wg *sync.WaitGroup, app *app.App) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	API := api.NewAPI(cfg, log, app)
	run(func() { API.Run(ctx, log) })

	//run(func() { eventlistener.NewListener(cfg, &app) })
}
