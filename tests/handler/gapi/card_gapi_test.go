package gapi_test

import (
	"MamangRust/paymentgatewaygrpc/internal/cache"
	card_cache "MamangRust/paymentgatewaygrpc/internal/cache/card"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/handler/gapi"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	"MamangRust/paymentgatewaygrpc/internal/service"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"MamangRust/paymentgatewaygrpc/pkg/observability"
	"MamangRust/paymentgatewaygrpc/tests"
	"context"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CardGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.CardServiceClient
	conn        *grpc.ClientConn
	repos       *repository.Repositories
	userID      int32
}

func (s *CardGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	obs, _ := observability.NewObservability("test", log)
	cacheStore := cache.NewCacheStore(s.redisClient, log, nil)
	cCard := card_cache.NewCardMencache(cacheStore)

	cardService := service.NewCardService(service.CardServiceDeps{
		CardRepo:      s.repos.Card,
		UserRepo:      s.repos.User,
		Logger:        log,
		Observability: obs,
		Cache:         cCard,
	})

	// Seed User
	user, err := s.repos.User.CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Card",
		LastName:  "Gapi",
		Email:     "card.gapi@example.com",
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = user.UserID

	// Start gRPC Server
	cardHandler := gapi.NewCardHandleGrpc(cardService)
	server := grpc.NewServer()
	pb.RegisterCardServiceServer(server, cardHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewCardServiceClient(conn)
}

func (s *CardGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	s.dbPool.Close()
	s.ts.Teardown()
}

func (s *CardGapiTestSuite) TestCardLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateCardRequest{
		UserId:       int32(s.userID),
		CardType:     "debit",
		CardProvider: "visa",
		Cvv:          "123",
		ExpireDate:   timestamppb.New(time.Now().AddDate(1, 0, 0)),
	}
	res, err := s.client.CreateCard(ctx, createReq)
	s.NoError(err)
	s.NotEmpty(res.Data.CardNumber)

	cardID := res.Data.Id

	// 2. Find By ID
	found, err := s.client.FindByIdCard(ctx, &pb.FindByIdCardRequest{CardId: cardID})
	s.NoError(err)
	s.Equal(cardID, found.Data.Id)

	// 3. Update
	updateReq := &pb.UpdateCardRequest{
		CardId:       cardID,
		UserId:       int32(s.userID),
		CardType:     "credit",
		CardProvider: "mastercard",
		Cvv:          "456",
		ExpireDate:   timestamppb.New(time.Now().AddDate(1, 0, 0)),
	}
	updated, err := s.client.UpdateCard(ctx, updateReq)
	s.NoError(err)
	s.Equal("credit", updated.Data.CardType)

	// 4. Delete (Trash)
	_, err = s.client.TrashedCard(ctx, &pb.FindByIdCardRequest{CardId: cardID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreCard(ctx, &pb.FindByIdCardRequest{CardId: cardID})
	s.NoError(err)

	// 6. Permanent Delete
	_, err = s.client.DeleteCardPermanent(ctx, &pb.FindByIdCardRequest{CardId: cardID})
	s.NoError(err)
}

func TestCardGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardGapiTestSuite))
}
