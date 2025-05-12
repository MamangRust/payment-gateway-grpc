package server_test

import (
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	protomapper "MamangRust/paymentgatewaygrpc/internal/mapper/proto"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	responsemapper "MamangRust/paymentgatewaygrpc/internal/mapper/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/auth"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/database/seeder"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"net"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type ServerTestSuite struct {
	suite.Suite
	client            *grpc.ClientConn
	ctx               context.Context
	authClient        pb.AuthServiceClient
	userClient        pb.UserServiceClient
	cardClient        pb.CardServiceClient
	merchantClient    pb.MerchantServiceClient
	saldoClient       pb.SaldoServiceClient
	topupClient       pb.TopupServiceClient
	txClient          pb.TransactionServiceClient
	transferClient    pb.TransferServiceClient
	transactionClient pb.TransactionServiceClient
	withdrawClient    pb.WithdrawServiceClient
	logger            logger.LoggerInterface
	cleanupFunc       func()
	listener          *bufconn.Listener
	testDB            *sql.DB
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	s.ctx = context.Background()

	logger, err := logger.NewLogger()
	require.NoError(s.T(), err)

	err = s.setupMigrationDB()
	if err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	testDB, err := s.setupTestDB()
	require.NoError(s.T(), err)

	err = s.setupSeeder()
	require.NoError(s.T(), err)

	listener := bufconn.Listen(bufSize)
	server := s.setupGRPCServer(logger, testDB)

	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(s.T(), err)

	s.client = conn
	s.listener = listener
	s.testDB = testDB
	s.logger = logger

	s.authClient = pb.NewAuthServiceClient(conn)
	s.userClient = pb.NewUserServiceClient(conn)
	s.cardClient = pb.NewCardServiceClient(conn)
	s.merchantClient = pb.NewMerchantServiceClient(conn)
	s.saldoClient = pb.NewSaldoServiceClient(conn)
	s.topupClient = pb.NewTopupServiceClient(conn)
	s.txClient = pb.NewTransactionServiceClient(conn)
	s.transactionClient = pb.NewTransactionServiceClient(conn)
	s.transferClient = pb.NewTransferServiceClient(conn)
	s.withdrawClient = pb.NewWithdrawServiceClient(conn)

	s.cleanupFunc = func() {
		if s.client != nil {
			s.client.Close()
		}
		if s.listener != nil {
			s.listener.Close()
		}
		if s.testDB != nil {
			s.testDB.Close()
		}
	}
}

func (s *ServerTestSuite) TearDownSuite() {
	if s.cleanupFunc != nil {
		s.cleanupFunc()
	}
}

func (s *ServerTestSuite) setupTestDB() (*sql.DB, error) {
	conn := "host=172.17.0.2 port=5432 user=postgres password=postgres dbname=paymentgateway_test sslmode=disable"

	db, err := sql.Open("postgres", conn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *ServerTestSuite) setupMigrationDB() error {
	connStr := "host=172.17.0.2 port=5432 user=postgres password=postgres dbname=paymentgateway_test sslmode=disable"

	migrationDir := "../../pkg/database/migrations"

	db, err := goose.OpenDBWithDriver("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (s *ServerTestSuite) setupSeeder() error {
	DB := db.New(s.testDB)

	seeder := seeder.NewSeeder(seeder.Deps{
		DB:     DB,
		Ctx:    s.ctx,
		Logger: s.logger,
	})

	if err := seeder.Run(); err != nil {
		s.logger.Fatal("Failed to run seeder", zap.Error(err))
	}

	return nil
}

func (s *ServerTestSuite) setupGRPCServer(logger logger.LoggerInterface, testDB *sql.DB) *grpc.Server {
	token, err := auth.NewManager("test_secret_key")
	require.NoError(s.T(), err)

	hash := hash.NewHashingPassword()
	DB := db.New(testDB)

	mapperRecord := recordmapper.NewRecordMapper()
	mapperResponse := responsemapper.NewResponseMapper()
	mapperProto := protomapper.NewProtoMapper()

	depsRepo := repository.Deps{
		DB:           DB,
		Ctx:          s.ctx,
		MapperRecord: mapperRecord,
	}

	repository := repository.NewRepositories(depsRepo)

	service := service.NewService(service.Deps{
		Repositories: repository,
		Hash:         hash,
		Token:        token,
		Logger:       logger,
		Mapper:       *mapperResponse,
	})

	handlerAuth := gapi.NewAuthHandleGrpc(service.Auth, mapperProto.AuthProtoMapper)
	handlerCard := gapi.NewCardHandleGrpc(service.Card, mapperProto.CardProtoMapper)
	handleMerchant := gapi.NewMerchantHandleGrpc(service.Merchant, mapperProto.MerchantProtoMapper)
	handlerUser := gapi.NewUserHandleGrpc(service.User, mapperProto.UserProtoMapper)
	handlerSaldo := gapi.NewSaldoHandleGrpc(service.Saldo, mapperProto.SaldoProtoMapper)
	handlerTopup := gapi.NewTopupHandleGrpc(service.Topup, mapperProto.TopupProtoMapper)
	handlerTransaction := gapi.NewTransactionHandleGrpc(service.Transaction, mapperProto.TransactionProtoMapper)
	handlerTransfer := gapi.NewTransferHandleGrpc(service.Transfer, mapperProto.TransferProtoMapper)
	handlerWithraw := gapi.NewWithdrawHandleGrpc(service.Withdraw, mapperProto.WithdrawalProtoMapper)

	c := grpc.NewServer()

	pb.RegisterAuthServiceServer(c, handlerAuth)
	pb.RegisterUserServiceServer(c, handlerUser)
	pb.RegisterCardServiceServer(c, handlerCard)
	pb.RegisterMerchantServiceServer(c, handleMerchant)
	pb.RegisterSaldoServiceServer(c, handlerSaldo)
	pb.RegisterTopupServiceServer(c, handlerTopup)
	pb.RegisterTransactionServiceServer(c, handlerTransaction)
	pb.RegisterTransferServiceServer(c, handlerTransfer)
	pb.RegisterWithdrawServiceServer(c, handlerWithraw)

	return c
}
