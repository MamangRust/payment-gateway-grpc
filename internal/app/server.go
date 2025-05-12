package app

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/service"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	"MamangRust/paymentgatewaygrpc/pkg/database"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/database/seeder"
	"MamangRust/paymentgatewaygrpc/pkg/dotenv"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
)

var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type Server struct {
	Logger       logger.LoggerInterface
	DB           *db.Queries
	TokenManager *auth.Manager
	Services     *service.Service
	Handlers     *gapi.Handler
	Ctx          context.Context
}

func NewServer() (*Server, error) {
	flag.Parse()

	logger, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		logger.Fatal("Failed to create token manager", zap.Error(err))
	}

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	DB := db.New(conn)

	ctx := context.Background()

	hash := hash.NewHashingPassword()
	mapperRecord := recordmapper.NewRecordMapper()
	mapperResponse := responseservice.NewResponseServiceMapper()

	depsRepo := repository.Deps{
		DB:           DB,
		Ctx:          ctx,
		MapperRecord: mapperRecord,
	}

	repositories := repository.NewRepositories(depsRepo)

	services := service.NewService(service.Deps{
		Repositories: repositories,
		Hash:         hash,
		Token:        tokenManager,
		Logger:       logger,
		Mapper:       *mapperResponse,
	})

	mapperProto := protomapper.NewProtoMapper()

	handlers := gapi.NewHandler(gapi.Deps{
		Service: *services,
		Mapper:  *mapperProto,
	})

	db_seeder := viper.GetString("DB_SEEDER")

	if db_seeder == "true" {
		logger.Debug("Seeding database", zap.String("seeder", db_seeder))

		seeder := seeder.NewSeeder(seeder.Deps{
			DB:     DB,
			Ctx:    ctx,
			Logger: logger,
		})

		if err := seeder.Run(); err != nil {
			logger.Fatal("Failed to run seeder", zap.Error(err))
		}

	}

	return &Server{
		Logger:       logger,
		DB:           DB,
		TokenManager: tokenManager,
		Services:     services,
		Handlers:     handlers,
		Ctx:          ctx,
	}, nil
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, s.Handlers.Auth)
	pb.RegisterUserServiceServer(grpcServer, s.Handlers.User)
	pb.RegisterRoleServiceServer(grpcServer, s.Handlers.Role)
	pb.RegisterCardServiceServer(grpcServer, s.Handlers.Card)
	pb.RegisterMerchantServiceServer(grpcServer, s.Handlers.Merchant)
	pb.RegisterSaldoServiceServer(grpcServer, s.Handlers.Saldo)
	pb.RegisterTopupServiceServer(grpcServer, s.Handlers.Topup)
	pb.RegisterTransactionServiceServer(grpcServer, s.Handlers.Transaction)
	pb.RegisterTransferServiceServer(grpcServer, s.Handlers.Transfer)
	pb.RegisterWithdrawServiceServer(grpcServer, s.Handlers.Withdraw)

	s.Logger.Info(fmt.Sprintf("Server running on port %d", *port))

	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}
