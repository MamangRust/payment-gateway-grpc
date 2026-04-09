package gapi

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/internal/service"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/errors/card_errors"
	"context"
	"fmt"
	"math"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type cardHandleGrpc struct {
	pb.UnimplementedCardServiceServer
	cardService service.CardService
}

func NewCardHandleGrpc(card service.CardService) *cardHandleGrpc {
	return &cardHandleGrpc{cardService: card}
}

func (s *cardHandleGrpc) FindAllCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCard, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cards, totalRecords, err := s.cardService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCards := make([]*pb.CardResponse, len(cards))
	for i, card := range cards {
		protoCards[i] = &pb.CardResponse{
			Id:         int32(card.CardID),
			UserId:     int32(card.UserID),
			CardNumber: card.CardNumber,
			CardType:   card.CardType,
			Cvv:        card.Cvv,
			ExpireDate: card.ExpireDate.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:  card.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  card.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCard{
		Status:     "success",
		Message:    "Successfully fetched card records",
		Data:       protoCards,
		Pagination: paginationMeta,
	}, nil
}

func (s *cardHandleGrpc) FindByIdCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetCardId())

	if id == 0 {

		return nil, card_errors.ErrGrpcInvalidCardID
	}

	card, err := s.cardService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	res := &pb.ApiResponseCard{
		Message: "successfully",
		Status:  "success",
		Data: &pb.CardResponse{
			Id:           int32(card.CardID),
			UserId:       int32(card.UserID),
			CardNumber:   card.CardNumber,
			CardType:     card.CardType,
			CardProvider: card.CardProvider,
			Cvv:          card.Cvv,
			ExpireDate:   card.ExpireDate.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:    card.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		},
	}

	return res, nil
}

func (s *cardHandleGrpc) FindByUserIdCard(ctx context.Context, req *pb.FindByUserIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetUserId())

	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidUserID
	}
	res, err := s.cardService.FindByUserID(ctx, id)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	Pbres := &pb.ApiResponseCard{
		Message: "successfully",
		Status:  "success",
		Data: &pb.CardResponse{
			Id:           int32(res.CardID),
			UserId:       int32(res.UserID),
			CardNumber:   res.CardNumber,
			CardType:     res.CardType,
			CardProvider: res.CardProvider,
			Cvv:          res.Cvv,
			ExpireDate:   res.ExpireDate.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:    res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:    res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		},
	}

	return Pbres, nil
}

func (s *cardHandleGrpc) DashboardCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseDashboardCard, error) {
	dashboardCard, err := s.cardService.DashboardCard(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbRes := &pb.ApiResponseDashboardCard{
		Status:  "success",
		Message: "Dashboard card retrieved successfully",
		Data: &pb.CardResponseDashboard{
			TotalBalance:     *dashboardCard.TotalBalance,
			TotalTopup:       *dashboardCard.TotalTopup,
			TotalTransaction: *dashboardCard.TotalTransaction,
			TotalTransfer:    *dashboardCard.TotalTransfer,
			TotalWithdraw:    *dashboardCard.TotalWithdraw,
		},
	}

	return pbRes, nil
}

func (s *cardHandleGrpc) DashboardCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseDashboardCardNumber, error) {
	card_number := req.GetCardNumber()

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	dashboardCard, err := s.cardService.DashboardCardCardNumber(ctx, card_number)

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbRes := &pb.ApiResponseDashboardCardNumber{
		Status:  "success",
		Message: "Dashboard card for card number retrieved successfully",
		Data: &pb.CardResponseDashboardCardNumber{
			TotalBalance:          *dashboardCard.TotalBalance,
			TotalTopup:            *dashboardCard.TotalTopup,
			TotalTransaction:      *dashboardCard.TotalTransaction,
			TotalTransferSend:     *dashboardCard.TotalTransferSend,
			TotalTransferReceiver: *dashboardCard.TotalTransferReceiver,
			TotalWithdraw:         *dashboardCard.TotalWithdraw,
		},
	}

	return pbRes, nil
}

func (s *cardHandleGrpc) FindMonthlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseMonthlyBalance, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindMonthlyBalance(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyBalance, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyBalance{
			Month:        item.Month,
			TotalBalance: int64(item.TotalBalance),
		}
	}

	return &pb.ApiResponseMonthlyBalance{
		Status:  "success",
		Message: "Monthly balance retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseYearlyBalance, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindYearlyBalance(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyBalance, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyBalance{
			Year:         item.Year.Int.String(),
			TotalBalance: item.TotalBalance,
		}
	}

	return &pb.ApiResponseYearlyBalance{
		Status:  "success",
		Message: "Yearly balance retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindMonthlyTopupAmount(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalTopupAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly topup amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindYearlyTopupAmount(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalTopupAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly topup amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindMonthlyWithdrawAmount(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalWithdrawAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly withdraw amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindYearlyWithdrawAmount(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalWithdrawAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly withdraw amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindMonthlyTransactionAmount(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalTransactionAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly transaction amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindYearlyTransactionAmount(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalTransactionAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly transaction amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindMonthlyTransferAmountSender(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalSentAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly transfer sender amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindYearlyTransferAmountSender(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalSentAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "transfer sender amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindMonthlyTransferAmountReceiver(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalReceivedAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly transfer receiver amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())
	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardService.FindYearlyTransferAmountReceiver(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalReceivedAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly transfer receiver amount retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseMonthlyBalance, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindMonthlyBalancesByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyBalance, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyBalance{
			Month:        item.Month,
			TotalBalance: int64(item.TotalBalance),
		}
	}

	return &pb.ApiResponseMonthlyBalance{
		Status:  "success",
		Message: "Monthly balance retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseYearlyBalance, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindYearlyBalanceByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyBalance, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyBalance{
			Year:         item.Year.Int.String(),
			TotalBalance: item.TotalBalance,
		}
	}

	return &pb.ApiResponseYearlyBalance{
		Status:  "success",
		Message: "Yearly balance retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindMonthlyTopupAmountByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalTopupAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly topup amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindYearlyTopupAmountByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalTopupAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly topup amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindMonthlyWithdrawAmountByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalWithdrawAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly withdraw amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindYearlyWithdrawAmountByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalWithdrawAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly withdraw amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindMonthlyTransactionAmountByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalTransactionAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly transaction amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindYearlyTransactionAmountByCardNumber(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalTransactionAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly transaction amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindMonthlyTransferAmountBySender(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalSentAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly transfer sender amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindYearlyTransferAmountBySender(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalSentAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly transfer sender amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindMonthlyTransferAmountByReceiver(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseMonthlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseMonthlyAmount{
			Month:       item.Month,
			TotalAmount: int64(item.TotalReceivedAmount),
		}
	}

	return &pb.ApiResponseMonthlyAmount{
		Status:  "success",
		Message: "Monthly transfer receiver amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindYearlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardService.FindYearlyTransferAmountByReceiver(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.CardResponseYearlyAmount, len(res))
	for i, item := range res {
		protoData[i] = &pb.CardResponseYearlyAmount{
			Year:        item.Year.Int.String(),
			TotalAmount: item.TotalReceivedAmount,
		}
	}

	return &pb.ApiResponseYearlyAmount{
		Status:  "success",
		Message: "Yearly transfer receiver amount by card number retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *cardHandleGrpc) FindByActiveCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.cardService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCards := make([]*pb.CardResponseDeleteAt, len(res))
	for i, card := range res {
		protoCards[i] = &pb.CardResponseDeleteAt{
			Id:           int32(card.CardID),
			UserId:       int32(card.UserID),
			CardNumber:   card.CardNumber,
			CardType:     card.CardType,
			Cvv:          card.Cvv,
			CardProvider: card.CardProvider,
			ExpireDate:   card.ExpireDate.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:    card.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:    &wrapperspb.StringValue{Value: card.DeletedAt.Time.Format("2006-01-02 15:04:05")},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCardDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched card record",
		Data:       protoCards,
		Pagination: paginationMeta,
	}, nil
}

func (s *cardHandleGrpc) FindByTrashedCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.cardService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCards := make([]*pb.CardResponseDeleteAt, len(res))
	for i, card := range res {
		protoCards[i] = &pb.CardResponseDeleteAt{
			Id:           int32(card.CardID),
			UserId:       int32(card.UserID),
			CardNumber:   card.CardNumber,
			CardType:     card.CardType,
			Cvv:          card.Cvv,
			CardProvider: card.CardProvider,
			ExpireDate:   card.ExpireDate.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:    card.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:    card.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:    &wrapperspb.StringValue{Value: card.DeletedAt.Time.Format("2006-01-02 15:04:05")},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationCardDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched card record",
		Data:       protoCards,
		Pagination: paginationMeta,
	}, nil
}

func (s *cardHandleGrpc) FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseCard, error) {
	card_number := req.GetCardNumber()
	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	res, err := s.cardService.FindByCardNumber(ctx, card_number)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCard := &pb.CardResponse{
		Id:         int32(res.CardID),
		UserId:     int32(res.UserID),
		CardNumber: res.CardNumber,
		CardType:   res.CardType,
		Cvv:        res.Cvv,
		ExpireDate: res.ExpireDate.Time.Format("2006-01-02 15:04:05"),
		CreatedAt:  res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully fetched card record",
		Data:    protoCard,
	}, nil
}

func (s *cardHandleGrpc) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.ApiResponseCard, error) {
	fmt.Println("DEBUG: CreateCard called")
	if req == nil {
		return nil, errors.ToGrpcError(card_errors.ErrGrpcValidateCreateCardRequest)
	}

	expireDate := time.Now().AddDate(1, 0, 0)
	if req.ExpireDate != nil {
		expireDate = req.ExpireDate.AsTime()
	}

	request := requests.CreateCardRequest{
		UserID:       int(req.UserId),
		CardType:     req.CardType,
		ExpireDate:   expireDate,
		CVV:          req.Cvv,
		CardProvider: req.CardProvider,
	}

	if err := request.Validate(); err != nil {
		return nil, card_errors.ErrGrpcValidateCreateCardRequest
	}

	res, err := s.cardService.CreateCard(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	if res == nil {
		return nil, errors.ToGrpcError(card_errors.ErrCardNotFound) // Or suitable error
	}

	protoCard := &pb.CardResponse{
		Id:         int32(res.CardID),
		UserId:     int32(res.UserID),
		CardNumber: res.CardNumber,
		CardType:   res.CardType,
		Cvv:        res.Cvv,
		ExpireDate: res.ExpireDate.Time.Format("2006-01-02 15:04:05"),
		CreatedAt:  res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully created card",
		Data:    protoCard,
	}, nil
}

func (s *cardHandleGrpc) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.ApiResponseCard, error) {
	if req == nil {
		return nil, errors.ToGrpcError(card_errors.ErrGrpcValidateUpdateCardRequest)
	}

	expireDate := time.Now().AddDate(1, 0, 0)
	if req.ExpireDate != nil {
		expireDate = req.ExpireDate.AsTime()
	}

	request := requests.UpdateCardRequest{
		CardID:       int(req.CardId),
		UserID:       int(req.UserId),
		CardType:     req.CardType,
		ExpireDate:   expireDate,
		CVV:          req.Cvv,
		CardProvider: req.CardProvider,
	}

	if err := request.Validate(); err != nil {
		return nil, card_errors.ErrGrpcValidateUpdateCardRequest
	}

	res, err := s.cardService.UpdateCard(ctx, &request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	if res == nil {
		return nil, errors.ToGrpcError(card_errors.ErrCardNotFound)
	}

	protoCard := &pb.CardResponse{
		Id:         int32(res.CardID),
		UserId:     int32(res.UserID),
		CardNumber: res.CardNumber,
		CardType:   res.CardType,
		Cvv:        res.Cvv,
		ExpireDate: res.ExpireDate.Time.Format("2006-01-02 15:04:05"),
		CreatedAt:  res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	return &pb.ApiResponseCard{
		Status:  "success",
		Message: "Successfully updated card",
		Data:    protoCard,
	}, nil
}

func (s *cardHandleGrpc) TrashedCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDeleteAt, error) {
	id := int(req.GetCardId())
	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	res, err := s.cardService.TrashedCard(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCard := &pb.CardResponseDeleteAt{
		Id:         int32(res.CardID),
		UserId:     int32(res.UserID),
		CardNumber: res.CardNumber,
		CardType:   res.CardType,
		Cvv:        res.Cvv,
		ExpireDate: res.ExpireDate.Time.Format("2006-01-02 15:04:05"),
		CreatedAt:  res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:  wrapperspb.String(res.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseCardDeleteAt{
		Status:  "success",
		Message: "Successfully trashed card",
		Data:    protoCard,
	}, nil
}

func (s *cardHandleGrpc) RestoreCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDeleteAt, error) {
	id := int(req.GetCardId())
	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	res, err := s.cardService.RestoreCard(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoCard := &pb.CardResponseDeleteAt{
		Id:         int32(res.CardID),
		UserId:     int32(res.UserID),
		CardNumber: res.CardNumber,
		CardType:   res.CardType,
		Cvv:        res.Cvv,
		ExpireDate: res.ExpireDate.Time.Format("2006-01-02 15:04:05"),
		CreatedAt:  res.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  res.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:  wrapperspb.String(res.DeletedAt.Time.Format("2006-01-02 15:04:05")),
	}

	return &pb.ApiResponseCardDeleteAt{
		Status:  "success",
		Message: "Successfully restored card",
		Data:    protoCard,
	}, nil
}

func (s *cardHandleGrpc) DeleteCardPermanent(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDelete, error) {
	id := int(req.GetCardId())
	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	_, err := s.cardService.DeleteCardPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCardDelete{
		Status:  "success",
		Message: "Successfully deleted card",
	}, nil
}

func (s *cardHandleGrpc) RestoreAllCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCardAll, error) {
	_, err := s.cardService.RestoreAllCard(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCardAll{
		Status:  "success",
		Message: "Successfully restore card",
	}, nil
}

func (s *cardHandleGrpc) DeleteAllCardPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCardAll, error) {
	_, err := s.cardService.DeleteAllCardPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseCardAll{
		Status:  "success",
		Message: "Successfully delete card permanent",
	}, nil
}
