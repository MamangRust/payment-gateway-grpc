package api

import (
	transaction_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transaction"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"fmt"

	"MamangRust/paymentgatewaygrpc/internal/middlewares"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type transactionHandleApi struct {
	transaction pb.TransactionServiceClient
	logger      logger.LoggerInterface
	mapping     apimapper.TransactionResponseMapper
	apihandler  errors.ApiHandler
	cache       transaction_cache.TransactionMencache
}

func NewHandlerTransaction(transaction pb.TransactionServiceClient, merchant pb.MerchantServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.TransactionResponseMapper, apiHandler errors.ApiHandler, cache transaction_cache.TransactionMencache) *transactionHandleApi {
	transactionHandler := transactionHandleApi{
		transaction: transaction,
		mapping:     mapping,
		logger:      logger,
		apihandler:  apiHandler,
		cache:       cache,
	}

	routerTransaction := router.Group("/api/transactions")

	routerTransaction.GET("", apiHandler.Handle("find-all-transactions", transactionHandler.FindAll))
	routerTransaction.GET("/card-number/:card_number", apiHandler.Handle("find-all-transactions-by-card-number", transactionHandler.FindAllTransactionByCardNumber))
	routerTransaction.GET("/:id", apiHandler.Handle("find-transaction-by-id", transactionHandler.FindById))
	routerTransaction.GET("/merchant/:merchant_id", apiHandler.Handle("find-transactions-by-merchant-id", transactionHandler.FindByTransactionMerchantId))
	routerTransaction.GET("/active", apiHandler.Handle("find-active-transactions", transactionHandler.FindByActiveTransaction))
	routerTransaction.GET("/trashed", apiHandler.Handle("find-trashed-transactions", transactionHandler.FindByTrashedTransaction))

	routerTransaction.GET("/monthly-success", apiHandler.Handle("find-monthly-transaction-status-success", transactionHandler.FindMonthlyTransactionStatusSuccess))
	routerTransaction.GET("/yearly-success", apiHandler.Handle("find-yearly-transaction-status-success", transactionHandler.FindYearlyTransactionStatusSuccess))
	routerTransaction.GET("/monthly-failed", apiHandler.Handle("find-monthly-transaction-status-failed", transactionHandler.FindMonthlyTransactionStatusFailed))
	routerTransaction.GET("/yearly-failed", apiHandler.Handle("find-yearly-transaction-status-failed", transactionHandler.FindYearlyTransactionStatusFailed))

	routerTransaction.GET("/monthly-success-by-card", apiHandler.Handle("find-monthly-transaction-status-success-by-card", transactionHandler.FindMonthlyTransactionStatusSuccessByCardNumber))
	routerTransaction.GET("/yearly-success-by-card", apiHandler.Handle("find-yearly-transaction-status-success-by-card", transactionHandler.FindYearlyTransactionStatusSuccessByCardNumber))
	routerTransaction.GET("/monthly-failed-by-card", apiHandler.Handle("find-monthly-transaction-status-failed-by-card", transactionHandler.FindMonthlyTransactionStatusFailedByCardNumber))
	routerTransaction.GET("/yearly-failed-by-card", apiHandler.Handle("find-yearly-transaction-status-failed-by-card", transactionHandler.FindYearlyTransactionStatusFailedByCardNumber))

	routerTransaction.GET("/monthly-methods", apiHandler.Handle("find-monthly-payment-methods", transactionHandler.FindMonthlyPaymentMethods))
	routerTransaction.GET("/yearly-methods", apiHandler.Handle("find-yearly-payment-methods", transactionHandler.FindYearlyPaymentMethods))
	routerTransaction.GET("/monthly-methods-by-card", apiHandler.Handle("find-monthly-payment-methods-by-card", transactionHandler.FindMonthlyPaymentMethodsByCardNumber))
	routerTransaction.GET("/yearly-methods-by-card", apiHandler.Handle("find-yearly-payment-methods-by-card", transactionHandler.FindYearlyPaymentMethodsByCardNumber))

	routerTransaction.GET("/monthly-amounts-by-card", apiHandler.Handle("find-monthly-amounts-by-card", transactionHandler.FindMonthlyAmountsByCardNumber))
	routerTransaction.GET("/yearly-amounts-by-card", apiHandler.Handle("find-yearly-amounts-by-card", transactionHandler.FindYearlyAmountsByCardNumber))
	routerTransaction.GET("/monthly-amounts", apiHandler.Handle("find-monthly-amounts", transactionHandler.FindMonthlyAmounts))
	routerTransaction.GET("/yearly-amounts", apiHandler.Handle("find-yearly-amounts", transactionHandler.FindYearlyAmounts))

	routerTransaction.POST("/create", middlewares.ApiKeyMiddleware(merchant)(apiHandler.Handle("create-transaction", transactionHandler.Create)))
	routerTransaction.POST("/update/:id", middlewares.ApiKeyMiddleware(merchant)(apiHandler.Handle("update-transaction", transactionHandler.Update)))

	routerTransaction.POST("/restore/:id", apiHandler.Handle("restore-transaction", transactionHandler.RestoreTransaction))
	routerTransaction.POST("/trashed/:id", apiHandler.Handle("trash-transaction", transactionHandler.TrashedTransaction))
	routerTransaction.DELETE("/permanent/:id", apiHandler.Handle("delete-transaction-permanent", transactionHandler.DeletePermanent))

	routerTransaction.POST("/restore/all", apiHandler.Handle("restore-all-transactions", transactionHandler.RestoreAllTransaction))
	routerTransaction.POST("/permanent/all", apiHandler.Handle("delete-all-transactions-permanent", transactionHandler.DeleteAllTransactionPermanent))

	return &transactionHandler
}

// @Summary Find all
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions [get]
func (h *transactionHandleApi) FindAll(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	reqCache := &requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransactionsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.transaction.FindAllTransaction(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransaction(res)
	h.cache.SetCachedTransactionsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find all transactions by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of transactions for a specific card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/card-number/{card_number} [get]
func (h *transactionHandleApi) FindAllTransactionByCardNumber(c echo.Context) error {
	cardNumber := c.Param("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	reqCache := &requests.FindAllTransactionCardNumber{
		CardNumber: cardNumber,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	cachedData, found := h.cache.GetCachedTransactionByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransactionCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.transaction.FindAllTransactionByCardNumber(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return h.handleGrpcError(err, "FindAllTransactionByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransaction(res)
	h.cache.SetCachedTransactionByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find a transaction by ID
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a transaction record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/{id} [get]
func (h *transactionHandleApi) FindById(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Debug("Invalid transaction ID", zap.Error(err))
		return errors.NewBadRequestError("invalid id parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTransactionCache(ctx, idInt)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindByIdTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseTransaction(res)
	h.cache.SetCachedTransactionCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransactionStatusSuccess retrieves the monthly transaction status for successful transactions.
// @Summary Get monthly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess "Monthly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for successful transactions"
// @Router /api/transactions/monthly-success [get]
func (h *transactionHandleApi) FindMonthlyTransactionStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetMonthTransactionStatusSuccessCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransactionStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthStatusSuccess(res)
	h.cache.SetMonthTransactionStatusSuccessCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
// @Summary Get yearly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusSuccess "Yearly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for successful transactions"
// @Router /api/transactions/yearly-success [get]
func (h *transactionHandleApi) FindYearlyTransactionStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearTransactionStatusSuccessCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyTransactionStatusSuccess(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransactionStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearStatusSuccess(res)
	h.cache.SetYearTransactionStatusSuccessCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransactionStatusFailed retrieves the monthly transaction status for failed transactions.
// @Summary Get monthly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed "Monthly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for failed transactions"
// @Router /api/transactions/monthly-failed [get]
func (h *transactionHandleApi) FindMonthlyTransactionStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetMonthTransactionStatusFailedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransactionStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthStatusFailed(res)
	h.cache.SetMonthTransactionStatusFailedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransactionStatusFailed retrieves the yearly transaction status for failed transactions.
// @Summary Get yearly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed "Yearly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for failed transactions"
// @Router /api/transactions/yearly-failed [get]
func (h *transactionHandleApi) FindYearlyTransactionStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearTransactionStatusFailedCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyTransactionStatusFailed(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransactionStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearStatusFailed(res)
	h.cache.SetYearTransactionStatusFailedCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransactionStatusSuccess retrieves the monthly transaction status for successful transactions.
// @Summary Get monthly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess "Monthly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for successful transactions"
// @Router /api/transactions/monthly-success-by-card [get]
func (h *transactionHandleApi) FindMonthlyTransactionStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetMonthTransactionStatusSuccessByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
		Month:      int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransactionStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthStatusSuccess(res)
	h.cache.SetMonthTransactionStatusSuccessByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
// @Summary Get yearly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param cardNumber query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionYearStatusSuccess "Yearly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for successful transactions"
// @Router /api/transactions/yearly-success-by-card [get]
func (h *transactionHandleApi) FindYearlyTransactionStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.YearStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearTransactionStatusSuccessByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyTransactionStatusSuccessByCardNumber(ctx, &pb.FindYearTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransactionStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearStatusSuccess(res)
	h.cache.SetYearTransactionStatusSuccessByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransactionStatusFailed retrieves the monthly transaction status for failed transactions.
// @Summary Get monthly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param cardNumber query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed "Monthly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for failed transactions"
// @Router /api/transactions/monthly-failed-by-card [get]
func (h *transactionHandleApi) FindMonthlyTransactionStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number paramater")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetMonthTransactionStatusFailedByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusFailedByCardNumber(ctx, &pb.FindMonthlyTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
		Month:      int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransactionStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthStatusFailed(res)
	h.cache.SetMonthTransactionStatusFailedByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransactionStatusFailedByCardNumber retrieves the yearly transaction status for failed transactions.
// @Summary Get yearly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for failed transactions by year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed "Yearly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for failed transactions"
// @Router /api/transactions/yearly-failed-by-card [get]
func (h *transactionHandleApi) FindYearlyTransactionStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.YearStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearTransactionStatusFailedByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyTransactionStatusFailedByCardNumber(ctx, &pb.FindYearTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransactionStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearStatusFailed(res)
	h.cache.SetYearTransactionStatusFailedByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyPaymentMethods retrieves the monthly payment methods for transactions.
// @Summary Get monthly payment methods
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly payment methods for transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/transactions/monthly-payment-methods [get]
func (h *transactionHandleApi) FindMonthlyPaymentMethods(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyPaymentMethodsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyPaymentMethods(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly payment methods", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyPaymentMethods")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthMethod(res)
	h.cache.SetMonthlyPaymentMethodsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyPaymentMethods retrieves the yearly payment methods for transactions.
// @Summary Get yearly payment methods
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly payment methods for transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/transactions/yearly-payment-methods [get]
func (h *transactionHandleApi) FindYearlyPaymentMethods(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyPaymentMethodsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyPaymentMethods(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly payment methods", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyPaymentMethods")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearMethod(res)
	h.cache.SetYearlyPaymentMethodsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmounts retrieves the monthly transaction amounts for a specific year.
// @Summary Get monthly transaction amounts
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/transactions/monthly-amounts [get]
func (h *transactionHandleApi) FindMonthlyAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyAmountsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyAmounts(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly amounts", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyAmounts")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthAmount(res)
	h.cache.SetMonthlyAmountsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmounts retrieves the yearly transaction amounts for a specific year.
// @Summary Get yearly transaction amounts
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/transactions/yearly-amounts [get]
func (h *transactionHandleApi) FindYearlyAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyAmountsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyAmounts(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly amounts", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyAmounts")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearAmount(res)
	h.cache.SetYearlyAmountsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyPaymentMethodsByCardNumber retrieves the monthly payment methods for transactions by card number and year.
// @Summary Get monthly payment methods by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly payment methods for transactions by card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod "Monthly payment methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods by card number"
// @Router /api/transactions/monthly-payment-methods-by-card [get]
func (h *transactionHandleApi) FindMonthlyPaymentMethodsByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyPaymentMethodsByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyPaymentMethodsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly payment methods by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyPaymentMethodsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthMethod(res)
	h.cache.SetMonthlyPaymentMethodsByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyPaymentMethodsByCardNumber retrieves the yearly payment methods for transactions by card number and year.
// @Summary Get yearly payment methods by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly payment methods for transactions by card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod "Yearly payment methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods by card number"
// @Router /api/transactions/yearly-payment-methods-by-card [get]
func (h *transactionHandleApi) FindYearlyPaymentMethodsByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyPaymentMethodsByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyPaymentMethodsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly payment methods by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyPaymentMethodsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearMethod(res)
	h.cache.SetYearlyPaymentMethodsByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountsByCardNumber retrieves the monthly transaction amounts for a specific card number and year.
// @Summary Get monthly transaction amounts by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount "Monthly transaction amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts by card number"
// @Router /api/transactions/monthly-amounts-by-card [get]
func (h *transactionHandleApi) FindMonthlyAmountsByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyAmountsByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindMonthlyAmountsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly amounts by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyAmountsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionMonthAmount(res)
	h.cache.SetMonthlyAmountsByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountsByCardNumber retrieves the yearly transaction amounts for a specific card number and year.
// @Summary Get yearly transaction amounts by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount "Yearly transaction amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts by card number"
// @Router /api/transactions/yearly-amounts-by-card [get]
func (h *transactionHandleApi) FindYearlyAmountsByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyAmountsByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.transaction.FindYearlyAmountsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly amounts by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyAmountsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransactionYearAmount(res)
	h.cache.SetYearlyAmountsByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find transactions by merchant ID
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of transactions using the merchant ID
// @Accept json
// @Produce json
// @Param merchant_id query string true "Merchant ID"
// @Success 200 {object} response.ApiResponseTransactions "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/merchant/{merchant_id} [get]
func (h *transactionHandleApi) FindByTransactionMerchantId(c echo.Context) error {
	merchantIdStr := c.Param("merchant_id")
	merchantIdInt, err := strconv.Atoi(merchantIdStr)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTransactionByMerchantIdCache(ctx, merchantIdInt)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindTransactionByMerchantIdRequest{
		MerchantId: int32(merchantIdInt),
	}

	res, err := h.transaction.FindTransactionByMerchantId(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return h.handleGrpcError(err, "FindByTransactionMerchantId")
	}

	apiResponse := h.mapping.ToApiResponseTransactions(res)
	h.cache.SetCachedTransactionByMerchantIdCache(ctx, merchantIdInt, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find active transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransactions "List of active transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/active [get]
func (h *transactionHandleApi) FindByActiveTransaction(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	reqCache := &requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransactionActiveCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.transaction.FindByActiveTransaction(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return h.handleGrpcError(err, "FindByActiveTransaction")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)
	h.cache.SetCachedTransactionActiveCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Retrieve trashed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of trashed transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransactions "List of trashed transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/trashed [get]
func (h *transactionHandleApi) FindByTrashedTransaction(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	reqCache := &requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransactionTrashedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.transaction.FindByTrashedTransaction(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return h.handleGrpcError(err, "FindByTrashedTransaction")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)
	h.cache.SetCachedTransactionTrashedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Create a new transaction
// @Tags Transaction
// @Security Bearer
// @Description Create a new transaction record with the provided details.
// @Accept json
// @Produce json
// @Param CreateTransactionRequest body requests.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transactions/create [post]
func (h *transactionHandleApi) Create(c echo.Context) error {
	var body requests.CreateTransactionRequest

	apiKey := c.Get("apiKey").(string)

	if apiKey == "" {
		return errors.NewBadRequestError("apiKey is required")
	}

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	res, err := h.transaction.CreateTransaction(ctx, &pb.CreateTransactionRequest{
		ApiKey:          apiKey,
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseTransaction(res)

	h.cache.SetCachedTransactionCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update a transaction
// @Tags Transaction
// @Security Bearer
// @Description Update an existing transaction record using its ID
// @Accept json
// @Produce json
// @Param transaction body requests.UpdateTransactionRequest true "Transaction data"
// @Success 200 {object} response.ApiResponseTransaction "Updated transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transactions/update [post]
func (h *transactionHandleApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateTransactionRequest

	body.MerchantID = &id

	apiKey, ok := c.Get("apiKey").(string)
	if !ok {
		return errors.NewBadRequestError("api-key is required")
	}

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	res, err := h.transaction.UpdateTransaction(ctx, &pb.UpdateTransactionRequest{
		TransactionId:   int32(id),
		CardNumber:      body.CardNumber,
		ApiKey:          apiKey,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseTransaction(res)

	h.cache.DeleteTransactionCache(ctx, id)
	h.cache.SetCachedTransactionCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a transaction
// @Tags Transaction
// @Security Bearer
// @Description Trash a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully trashed transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transaction"
// @Router /api/transactions/trashed/{id} [post]
func (h *transactionHandleApi) TrashedTransaction(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.transaction.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	h.cache.DeleteTransactionCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction
// @Security Bearer
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transactions/restore/{id} [post]
func (h *transactionHandleApi) RestoreTransaction(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.transaction.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	h.cache.DeleteTransactionCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction
// @Security Bearer
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transactions/permanent/{id} [delete]
func (h *transactionHandleApi) DeletePermanent(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.transaction.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "DeleteTransaction")
	}

	so := h.mapping.ToApiResponseTransactionDelete(res)

	h.cache.DeleteTransactionCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction
// @Security Bearer
// @Description Restore a trashed transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transactions/restore/all [post]
func (h *transactionHandleApi) RestoreAllTransaction(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.transaction.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all transaction", zap.Error(err))
		return h.handleGrpcError(err, "RestoreAll")
	}

	h.logger.Debug("Successfully restored all transaction")

	so := h.mapping.ToApiResponseTransactionAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction
// @Security Bearer
// @Description Permanently delete a transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transactions/delete/all [post]
func (h *transactionHandleApi) DeleteAllTransactionPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.transaction.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all transaction", zap.Error(err))

		return h.handleGrpcError(err, "DeleteAll")
	}

	h.logger.Debug("Successfully deleted all transaction permanently")

	so := h.mapping.ToApiResponseTransactionAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Transaction").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Transaction already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Transaction service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *transactionHandleApi) parseValidationErrors(err error) []errors.ValidationError {
	var validationErrs []errors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, errors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []errors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *transactionHandleApi) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
