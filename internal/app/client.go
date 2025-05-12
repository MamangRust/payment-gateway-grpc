package app

import (
	"MamangRust/paymentgatewaygrpc/internal/handler/api"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper/response/api"
	"MamangRust/paymentgatewaygrpc/internal/middlewares"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	"MamangRust/paymentgatewaygrpc/pkg/dotenv"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "MamangRust/paymentgatewaygrpc/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// @title PaymentGateway gRPC
// @version 1.0
// @description gRPC based Payment Gateway service

// @host localhost:5000
// @BasePath /api/

// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
func RunClient() {
	flag.Parse()

	logger, err := logger.NewLogger()
	limiter := middlewares.NewRateLimiter(20, 50)

	if err != nil {
		logger.Fatal("Failed to create logger", zap.Error(err))
	}

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Fatal("Failed to connect to server", zap.Error(err))
	}

	err = dotenv.Viper()

	if err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	e := echo.New()

	e.Use(limiter.Limit)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1420", "http://localhost:33451"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-API-Key",
		},
		AllowCredentials: true,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	middlewares.WebSecurityConfig(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	token, err := auth.NewManager(viper.GetString("SECRET_KEY"))

	if err != nil {
		logger.Fatal("Failed to create token manager", zap.Error(err))
	}

	mapping := apimapper.NewResponseApiMapper()

	depsHandler := api.Deps{
		Conn:    conn,
		Token:   token,
		E:       e,
		Logger:  logger,
		Mapping: *mapping,
	}

	api.NewHandler(depsHandler)

	go func() {
		if err := e.Start(":5000"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Server.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
