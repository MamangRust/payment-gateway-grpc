package protomapper

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"MamangRust/paymentgatewaygrpc/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type merchantProtoMapper struct{}

func NewMerchantProtoMapper() *merchantProtoMapper {
	return &merchantProtoMapper{}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchant(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantResponse) *pb.ApiResponsePaginationMerchant {
	return &pb.ApiResponsePaginationMerchant{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantResponses(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchantDeleteAt(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantResponseDeleteAt) *pb.ApiResponsePaginationMerchantDeleteAt {
	return &pb.ApiResponsePaginationMerchantDeleteAt{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantResponsesDeleteAt(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) ToProtoResponsePaginationMerchantTransaction(pagination *pb.PaginationMeta, status string, message string, merchants []*response.MerchantTransactionResponse) *pb.ApiResponsePaginationMerchantTransaction {

	return &pb.ApiResponsePaginationMerchantTransaction{
		Status:     status,
		Message:    message,
		Data:       m.mapMerchantTransactionResponses(merchants),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMonthlyPaymentMethods(status string, message string, ms []*response.MerchantResponseMonthlyPaymentMethod) *pb.ApiResponseMerchantMonthlyPaymentMethod {
	return &pb.ApiResponseMerchantMonthlyPaymentMethod{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMonthlyPaymentMethod(ms),
	}
}

func (m *merchantProtoMapper) ToProtoResponseYearlyPaymentMethods(status string, message string, ms []*response.MerchantResponseYearlyPaymentMethod) *pb.ApiResponseMerchantYearlyPaymentMethod {
	return &pb.ApiResponseMerchantYearlyPaymentMethod{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesYearlyPaymentMethod(ms),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMonthlyAmounts(status string, message string, ms []*response.MerchantResponseMonthlyAmount) *pb.ApiResponseMerchantMonthlyAmount {
	return &pb.ApiResponseMerchantMonthlyAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMonthlyAmount(ms),
	}
}

func (m *merchantProtoMapper) ToProtoResponseYearlyAmounts(status string, message string, ms []*response.MerchantResponseYearlyAmount) *pb.ApiResponseMerchantYearlyAmount {
	return &pb.ApiResponseMerchantYearlyAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesYearlyAmount(ms),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMonthlyTotalAmounts(status string, message string, ms []*response.MerchantResponseMonthlyTotalAmount) *pb.ApiResponseMerchantMonthlyTotalAmount {
	return &pb.ApiResponseMerchantMonthlyTotalAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesMonthlyTotalAmount(ms),
	}
}

func (m *merchantProtoMapper) ToProtoResponseYearlyTotalAmounts(status string, message string, ms []*response.MerchantResponseYearlyTotalAmount) *pb.ApiResponseMerchantYearlyTotalAmount {
	return &pb.ApiResponseMerchantYearlyTotalAmount{
		Status:  status,
		Message: message,
		Data:    m.mapResponsesYearlyTotalAmount(ms),
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchant(status string, message string, res *response.MerchantResponse) *pb.ApiResponseMerchant {
	return &pb.ApiResponseMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantResponse(res),
	}

}

func (m *merchantProtoMapper) ToProtoResponseMerchants(status string, message string, res []*response.MerchantResponse) *pb.ApiResponsesMerchant {
	return &pb.ApiResponsesMerchant{
		Status:  status,
		Message: message,
		Data:    m.mapMerchantResponses(res),
	}

}

func (m *merchantProtoMapper) ToProtoResponseMerchantAll(status string, message string) *pb.ApiResponseMerchantAll {
	return &pb.ApiResponseMerchantAll{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) ToProtoResponseMerchantDelete(status string, message string) *pb.ApiResponseMerchantDelete {
	return &pb.ApiResponseMerchantDelete{
		Status:  status,
		Message: message,
	}
}

func (m *merchantProtoMapper) mapMerchantResponse(merchant *response.MerchantResponse) *pb.MerchantResponse {
	return &pb.MerchantResponse{
		Id:        int32(merchant.ID),
		Name:      merchant.Name,
		Status:    merchant.Status,
		ApiKey:    merchant.ApiKey,
		UserId:    int32(merchant.UserID),
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
	}
}

func (s *merchantProtoMapper) mapMerchantResponses(roles []*response.MerchantResponse) []*pb.MerchantResponse {
	var responseRoles []*pb.MerchantResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapMerchantResponse(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapMerchantResponseDeleteAt(merchant *response.MerchantResponseDeleteAt) *pb.MerchantResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if merchant.DeletedAt != nil {
		deletedAt = wrapperspb.String(*merchant.DeletedAt)
	}

	return &pb.MerchantResponseDeleteAt{
		Id:        int32(merchant.ID),
		Name:      merchant.Name,
		Status:    merchant.Status,
		UserId:    int32(merchant.UserID),
		ApiKey:    merchant.ApiKey,
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (s *merchantProtoMapper) mapMerchantResponsesDeleteAt(roles []*response.MerchantResponseDeleteAt) []*pb.MerchantResponseDeleteAt {
	var responseRoles []*pb.MerchantResponseDeleteAt

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapMerchantResponseDeleteAt(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapMerchantTransactionResponse(merchant *response.MerchantTransactionResponse) *pb.MerchantTransactionResponse {
	return &pb.MerchantTransactionResponse{
		Id:              int32(merchant.ID),
		CardNumber:      merchant.CardNumber,
		Amount:          merchant.Amount,
		PaymentMethod:   merchant.PaymentMethod,
		MerchantId:      merchant.MerchantID,
		MerchantName:    merchant.MerchantName,
		TransactionTime: merchant.TransactionTime,
		CreatedAt:       merchant.CreatedAt,
		UpdatedAt:       merchant.UpdatedAt,
	}
}

func (s *merchantProtoMapper) mapMerchantTransactionResponses(roles []*response.MerchantTransactionResponse) []*pb.MerchantTransactionResponse {
	var responseRoles []*pb.MerchantTransactionResponse

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapMerchantTransactionResponse(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapResponseMonthlyPaymentMethod(ms *response.MerchantResponseMonthlyPaymentMethod) *pb.MerchantResponseMonthlyPaymentMethod {
	return &pb.MerchantResponseMonthlyPaymentMethod{
		Month:         ms.Month,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int64(ms.TotalAmount),
	}
}

func (s *merchantProtoMapper) mapResponsesMonthlyPaymentMethod(roles []*response.MerchantResponseMonthlyPaymentMethod) []*pb.MerchantResponseMonthlyPaymentMethod {
	var responseRoles []*pb.MerchantResponseMonthlyPaymentMethod

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseMonthlyPaymentMethod(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapResponseYearlyPaymentMethod(ms *response.MerchantResponseYearlyPaymentMethod) *pb.MerchantResponseYearlyPaymentMethod {
	return &pb.MerchantResponseYearlyPaymentMethod{
		Year:          ms.Year,
		PaymentMethod: ms.PaymentMethod,
		TotalAmount:   int64(ms.TotalAmount),
	}
}

func (s *merchantProtoMapper) mapResponsesYearlyPaymentMethod(roles []*response.MerchantResponseYearlyPaymentMethod) []*pb.MerchantResponseYearlyPaymentMethod {
	var responseRoles []*pb.MerchantResponseYearlyPaymentMethod

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseYearlyPaymentMethod(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapResponseMonthlyAmount(ms *response.MerchantResponseMonthlyAmount) *pb.MerchantResponseMonthlyAmount {
	return &pb.MerchantResponseMonthlyAmount{
		Month:       ms.Month,
		TotalAmount: int64(ms.TotalAmount),
	}
}

func (s *merchantProtoMapper) mapResponsesMonthlyAmount(roles []*response.MerchantResponseMonthlyAmount) []*pb.MerchantResponseMonthlyAmount {
	var responseRoles []*pb.MerchantResponseMonthlyAmount

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseMonthlyAmount(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapResponseYearlyAmount(ms *response.MerchantResponseYearlyAmount) *pb.MerchantResponseYearlyAmount {
	return &pb.MerchantResponseYearlyAmount{
		Year:        ms.Year,
		TotalAmount: int64(ms.TotalAmount),
	}
}

func (s *merchantProtoMapper) mapResponsesYearlyAmount(roles []*response.MerchantResponseYearlyAmount) []*pb.MerchantResponseYearlyAmount {
	var responseRoles []*pb.MerchantResponseYearlyAmount

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseYearlyAmount(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapResponseMonthlyTotalAmount(ms *response.MerchantResponseMonthlyTotalAmount) *pb.MerchantResponseMonthlyTotalAmount {
	return &pb.MerchantResponseMonthlyTotalAmount{
		Month:       ms.Month,
		Year:        ms.Year,
		TotalAmount: int64(ms.TotalAmount),
	}
}

func (s *merchantProtoMapper) mapResponsesMonthlyTotalAmount(roles []*response.MerchantResponseMonthlyTotalAmount) []*pb.MerchantResponseMonthlyTotalAmount {
	var responseRoles []*pb.MerchantResponseMonthlyTotalAmount

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseMonthlyTotalAmount(role))
	}

	return responseRoles
}

func (m *merchantProtoMapper) mapResponseYearlyTotalAmount(ms *response.MerchantResponseYearlyTotalAmount) *pb.MerchantResponseYearlyTotalAmount {
	return &pb.MerchantResponseYearlyTotalAmount{
		Year:        ms.Year,
		TotalAmount: int64(ms.TotalAmount),
	}
}

func (s *merchantProtoMapper) mapResponsesYearlyTotalAmount(roles []*response.MerchantResponseYearlyTotalAmount) []*pb.MerchantResponseYearlyTotalAmount {
	var responseRoles []*pb.MerchantResponseYearlyTotalAmount

	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseYearlyTotalAmount(role))
	}

	return responseRoles
}
