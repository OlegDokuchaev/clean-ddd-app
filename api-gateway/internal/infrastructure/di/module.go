package di

import (
	"api-gateway/internal/adapter/input/api"
	courierApi "api-gateway/internal/adapter/input/api/courier"
	customerApi "api-gateway/internal/adapter/input/api/customer"
	orderApi "api-gateway/internal/adapter/input/api/order"
	warehouseApi "api-gateway/internal/adapter/input/api/warehouse"
	"api-gateway/internal/adapter/output/auth/admin"
	courierClient "api-gateway/internal/adapter/output/clients/courier"
	customerClient "api-gateway/internal/adapter/output/clients/customer"
	orderClient "api-gateway/internal/adapter/output/clients/order"
	warehouseClient "api-gateway/internal/adapter/output/clients/warehouse"
	courierUseCase "api-gateway/internal/domain/usecases/courier"
	customerUseCase "api-gateway/internal/domain/usecases/customer"
	orderUseCase "api-gateway/internal/domain/usecases/order"
	warehouseUseCase "api-gateway/internal/domain/usecases/warehouse"
	"api-gateway/internal/infrastructure/logger"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"go.uber.org/fx"
)

var ClientsModule = fx.Options(
	fx.Provide(
		// Courier Client
		courierClient.NewConfig,
		courierClient.NewGRPCClient,
		courierClient.NewClient,

		// Order Client
		orderClient.NewConfig,
		orderClient.NewGRPCClient,
		orderClient.NewClient,

		// Customer Client
		customerClient.NewConfig,
		customerClient.NewGRPCClient,
		customerClient.NewClient,

		// Warehouse Client
		warehouseClient.NewConfig,
		warehouseClient.NewGRPCClient,
		warehouseClient.NewClient,
	),
)

var AuthModule = fx.Options(
	// Admin Auth
	fx.Provide(
		admin.NewConfig,
		admin.NewAuth,
	),
)

var UseCasesModule = fx.Options(
	fx.Provide(
		// Courier Use Case
		courierUseCase.NewUseCase,

		// Order Use Case
		orderUseCase.NewUseCase,

		// Customer Use Case
		customerUseCase.NewUseCase,

		// Warehouse Use Case
		warehouseUseCase.NewUseCase,
	),
)

var ApiModule = fx.Options(
	fx.Provide(
		// Courier Handler
		courierApi.NewHandler,

		// Order Handler
		orderApi.NewHandler,

		// Customer Handler
		customerApi.NewHandler,

		// Warehouse Handler
		warehouseApi.NewHandler,

		// API
		api.NewConfig,
		api.NewAPI,
	),

	fx.Invoke(RunServer),
)

var LoggerModule = fx.Provide(
	// Config
	logger.NewConfig,

	// Logstash
	logger.NewLogstash,

	// Logrus
	logger.NewLogrus,

	// Logger
	logger.NewLogger,
)

func RunServer(lc fx.Lifecycle, router *gin.Engine, config *api.Config, log logger.Logger) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("Could not start HTTP server", map[string]any{"error": err.Error()})
					panic(err)
				}
			}()
			log.Println("HTTP server started on ", config.Port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")
			return srv.Shutdown(ctx)
		},
	})
}
