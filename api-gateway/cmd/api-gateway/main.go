package main

import (
	_ "api-gateway/internal/adapter/input/api/courier/request"
	"api-gateway/internal/infrastructure/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		di.ClientsModule,
		di.AuthModule,
		di.UseCasesModule,
		di.ApiModule,
	).Run()
}
