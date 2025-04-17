package main

import (
	_ "api-gateway/internal/adapter/input/api/courier/request"
	"api-gateway/internal/infrastructure/di"

	"go.uber.org/fx"
)

//	@title			Clean DDD App API Gateway
//	@version		1.0
//	@description	This is the API Gateway for the Clean DDD application.
//	@BasePath		/

//	@securityDefinitions.apikey	CourierBearerAuth
//	@in							header
//	@name						Authorization
//	@description			    Courier's JWT token.

//	@securityDefinitions.apikey	CustomerBearerAuth
//	@in							header
//	@name						Authorization
//	@description				Customer's JWT token.

//	@securityDefinitions.apikey	AdminAccessToken
//	@in							header
//	@name						X-Access-Token
//	@description				Admin's access token.

func main() {
	fx.New(
		di.ConfigModule,
		di.ClientsModule,
		di.AuthModule,
		di.UseCasesModule,
		di.LoggerModule,
		di.ApiModule,
	).Run()
}
