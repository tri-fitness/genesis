package api

import (
	"context"
	"tri-fitness/genesis/api/middleware"
	r "tri-fitness/genesis/api/representations"
	j "tri-fitness/genesis/api/representations/json"
	"tri-fitness/genesis/api/resources"
	"tri-fitness/genesis/api/server"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(resources.NewAccountResource),
	fx.Provide(server.New),
	fx.Provide(zap.NewDevelopment),
	fx.Provide(middleware.NewAuthenticator),
	fx.Provide(func() map[string]r.RepresentationFactory {
		factories := make(map[string]r.RepresentationFactory)
		factories["application/json"] = j.NewJSONRepresentationFactory()
		return factories
	}),
	fx.Invoke(Start),
)

func Start(lc fx.Lifecycle, s server.Server) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Stop(ctx)
		},
	})
}
