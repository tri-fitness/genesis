package main

import (
	"time"
	"tri-fitness/genesis/api"
	"tri-fitness/genesis/application"
	"tri-fitness/genesis/config"
	"tri-fitness/genesis/infrastructure"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.StartTimeout(time.Second*45),
		api.Module,
		config.Module,
		application.Module,
		infrastructure.Module,
	).Run()
}
