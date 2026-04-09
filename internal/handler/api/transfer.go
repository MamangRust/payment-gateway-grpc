package api

import (
	transfer_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/transfer"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"MamangRust/paymentgatewaygrpc/internal/pb"
	"MamangRust/paymentgatewaygrpc/pkg/errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transferHandleApi struct {
	client     pb.TransferServiceClient
	logger     logger.LoggerInterface
	mapping    apimapper.TransferResponseMapper
	apihandler errors.ApiHandler
	cache      transfer_cache.TransferMencache
}

func NewHandlerTransfer(client pb.TransferServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.TransferResponseMapper, apiHandler errors.ApiHandler, cache transfer_cache.TransferMencache) *transferHandleApi {
	transferHandler := &transferHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apihandler: apiHandler,
		cache:      cache,
	}

	routerTransfer := router.Group("/api/transfers")

	routerTransfer.GET("", apiHandler.Handle("find-all-transfers", transferHandler.FindAll))
	routerTransfer.GET("/:id", apiHandler.Handle("find-transfer-by-id", transferHandler.FindById))
	routerTransfer.GET("/transfer_from/:transfer_from", apiHandler.Handle("find-transfers-by-transfer-from", transferHandler.FindByTransferByTransferFrom))
	routerTransfer.GET("/transfer_to/:transfer_to", apiHandler.Handle("find-transfers-by-transfer-to", transferHandler.FindByTransferByTransferTo))

	routerTransfer.GET("/active", apiHandler.Handle("find-active-transfers", transferHandler.FindByActiveTransfer))
	routerTransfer.GET("/trashed", apiHandler.Handle("find-trashed-transfers", transferHandler.FindByTrashedTransfer))

	routerTransfer.GET("/monthly-success", apiHandler.Handle("find-monthly-transfer-status-success", transferHandler.FindMonthlyTransferStatusSuccess))
	routerTransfer.GET("/yearly-success", apiHandler.Handle("find-yearly-transfer-status-success", transferHandler.FindYearlyTransferStatusSuccess))
	routerTransfer.GET("/monthly-failed", apiHandler.Handle("find-monthly-transfer-status-failed", transferHandler.FindMonthlyTransferStatusFailed))
	routerTransfer.GET("/yearly-failed", apiHandler.Handle("find-yearly-transfer-status-failed", transferHandler.FindYearlyTransferStatusFailed))

	routerTransfer.GET("/monthly-success-by-card", apiHandler.Handle("find-monthly-transfer-status-success-by-card", transferHandler.FindMonthlyTransferStatusSuccessByCardNumber))
	routerTransfer.GET("/yearly-success-by-card", apiHandler.Handle("find-yearly-transfer-status-success-by-card", transferHandler.FindYearlyTransferStatusSuccessByCardNumber))
	routerTransfer.GET("/monthly-failed-by-card", apiHandler.Handle("find-monthly-transfer-status-failed-by-card", transferHandler.FindMonthlyTransferStatusFailedByCardNumber))
	routerTransfer.GET("/yearly-failed-by-card", apiHandler.Handle("find-yearly-transfer-status-failed-by-card", transferHandler.FindYearlyTransferStatusFailedByCardNumber))

	routerTransfer.GET("/monthly-amount", apiHandler.Handle("find-monthly-transfer-amounts", transferHandler.FindMonthlyTransferAmounts))
	routerTransfer.GET("/yearly-amount", apiHandler.Handle("find-yearly-transfer-amounts", transferHandler.FindYearlyTransferAmounts))
	routerTransfer.GET("/monthly-by-sender", apiHandler.Handle("find-monthly-transfer-amounts-by-sender", transferHandler.FindMonthlyTransferAmountsBySenderCardNumber))
	routerTransfer.GET("/monthly-by-receiver", apiHandler.Handle("find-monthly-transfer-amounts-by-receiver", transferHandler.FindMonthlyTransferAmountsByReceiverCardNumber))
	routerTransfer.GET("/yearly-by-sender", apiHandler.Handle("find-yearly-transfer-amounts-by-sender", transferHandler.FindYearlyTransferAmountsBySenderCardNumber))
	routerTransfer.GET("/yearly-by-receiver", apiHandler.Handle("find-yearly-transfer-amounts-by-receiver", transferHandler.FindYearlyTransferAmountsByReceiverCardNumber))

	routerTransfer.POST("/create", apiHandler.Handle("create-transfer", transferHandler.CreateTransfer))
	routerTransfer.POST("/update/:id", apiHandler.Handle("update-transfer", transferHandler.UpdateTransfer))
	routerTransfer.POST("/trashed/:id", apiHandler.Handle("trash-transfer", transferHandler.TrashTransfer))
	routerTransfer.POST("/restore/:id", apiHandler.Handle("restore-transfer", transferHandler.RestoreTransfer))
	routerTransfer.DELETE("/permanent/:id", apiHandler.Handle("delete-transfer-permanent", transferHandler.DeleteTransferPermanent))

	routerTransfer.POST("/restore/all", apiHandler.Handle("restore-all-transfers", transferHandler.RestoreAllTransfer))
	routerTransfer.POST("/permanent/all", apiHandler.Handle("delete-all-transfers-permanent", transferHandler.DeleteAllTransferPermanent))

	return transferHandler
}

// @Summary Find all transfer records
// @Tags Transfer
// @Security Bearer
// @Description Retrieve a list of all transfer records with pagination
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransfer "List of transfer records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer [get]
func (h *transferHandleApi) FindAll(c echo.Context) error {
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

	reqCache := &requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransfersCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTransfer(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transfer data", zap.Error(err))
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransfer(res)
	h.cache.SetCachedTransfersCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find a transfer by ID
// @Tags Transfer
// @Security Bearer
// @Description Retrieve a transfer record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer/{id} [get]
func (h *transferHandleApi) FindById(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))
		return errors.NewBadRequestError("invalid id parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTransferCache(ctx, idInt)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByIdTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve transfer data", zap.Error(err))
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseTransfer(res)
	h.cache.SetCachedTransferCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferStatusSuccess retrieves the monthly transfer status for successful transactions.
// @Summary Get monthly transfer status for successful transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusSuccess "Monthly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for successful transactions"
// @Router /api/transfers/monthly-success [get]
func (h *transferHandleApi) FindMonthlyTransferStatusSuccess(c echo.Context) error {
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

	reqCache := &requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetCachedMonthTransferStatusSuccess(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferStatusSuccess(ctx, &pb.FindMonthlyTransferStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthStatusSuccess(res)
	h.cache.SetCachedMonthTransferStatusSuccess(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferStatusSuccess retrieves the yearly transfer status for successful transactions.
// @Summary Get yearly transfer status for successful transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusSuccess "Yearly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for successful transactions"
// @Router /api/transfers/yearly-success [get]
func (h *transferHandleApi) FindYearlyTransferStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearlyTransferStatusSuccess(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferStatusSuccess(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearStatusSuccess(res)
	h.cache.SetCachedYearlyTransferStatusSuccess(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferStatusFailed retrieves the monthly transfer status for failed transactions.
// @Summary Get monthly transfer status for failed transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusFailed "Monthly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for failed transactions"
// @Router /api/transfers/monthly-failed [get]
func (h *transferHandleApi) FindMonthlyTransferStatusFailed(c echo.Context) error {
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

	reqCache := &requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetCachedMonthTransferStatusFailed(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferStatusFailed(ctx, &pb.FindMonthlyTransferStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthStatusFailed(res)
	h.cache.SetCachedMonthTransferStatusFailed(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferStatusFailed retrieves the yearly transfer status for failed transactions.
// @Summary Get yearly transfer status for failed transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusFailed "Yearly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for failed transactions"
// @Router /api/transfers/yearly-failed [get]
func (h *transferHandleApi) FindYearlyTransferStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearlyTransferStatusFailed(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferStatusFailed(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearStatusFailed(res)
	h.cache.SetCachedYearlyTransferStatusFailed(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferStatusSuccessByCardNumber retrieves the monthly transfer status for successful transactions.
// @Summary Get monthly transfer status for successful transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusSuccess "Monthly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for successful transactions"
// @Router /api/transfers/monthly-success-by-card [get]
func (h *transferHandleApi) FindMonthlyTransferStatusSuccessByCardNumber(c echo.Context) error {
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

	reqCache := &requests.MonthStatusTransferCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetMonthTransferStatusSuccessByCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTransferStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthStatusSuccess(res)
	h.cache.SetMonthTransferStatusSuccessByCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferStatusSuccessByCardNumber retrieves the yearly transfer status for successful transactions.
// @Summary Get yearly transfer status for successful transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransferYearStatusSuccess "Yearly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for successful transactions"
// @Router /api/transfers/yearly-success-by-card [get]
func (h *transferHandleApi) FindYearlyTransferStatusSuccessByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearStatusTransferCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTransferStatusSuccessByCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferStatusSuccessByCardNumber(ctx, &pb.FindYearTransferStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearStatusSuccess(res)
	h.cache.SetYearlyTransferStatusSuccessByCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferStatusFailedByCardNumber retrieves the monthly transfer status for failed transactions.
// @Summary Get monthly transfer status for failed transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransferMonthStatusFailed "Monthly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for failed transactions"
// @Router /api/transfers/monthly-failed-by-card [get]
func (h *transferHandleApi) FindMonthlyTransferStatusFailedByCardNumber(c echo.Context) error {
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

	reqCache := &requests.MonthStatusTransferCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetMonthTransferStatusFailedByCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferStatusFailedByCardNumber(ctx, &pb.FindMonthlyTransferStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthStatusFailed(res)
	h.cache.SetMonthTransferStatusFailedByCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferStatusFailedByCardNumber retrieves the yearly transfer status for failed transactions.
// @Summary Get yearly transfer status for failed transactions
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransferYearStatusFailed "Yearly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for failed transactions"
// @Router /api/transfers/yearly-failed-by-card [get]
func (h *transferHandleApi) FindYearlyTransferStatusFailedByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearStatusTransferCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTransferStatusFailedByCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferStatusFailedByCardNumber(ctx, &pb.FindYearTransferStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearStatusFailed(res)
	h.cache.SetYearlyTransferStatusFailedByCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferAmounts retrieves the monthly transfer amounts for a specific year.
// @Summary Get monthly transfer amounts
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount "Monthly transfer amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer amounts"
// @Router /api/transfers/monthly-amounts [get]
func (h *transferHandleApi) FindMonthlyTransferAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMonthTransferAmounts(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferAmounts(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer amounts", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferAmounts")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthAmount(res)
	h.cache.SetCachedMonthTransferAmounts(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferAmounts retrieves the yearly transfer amounts for a specific year.
// @Summary Get yearly transfer amounts
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount "Yearly transfer amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer amounts"
// @Router /api/transfers/yearly-amounts [get]
func (h *transferHandleApi) FindYearlyTransferAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearlyTransferAmounts(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferAmounts(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer amounts", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferAmounts")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearAmount(res)
	h.cache.SetCachedYearlyTransferAmounts(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferAmountsBySenderCardNumber retrieves the monthly transfer amounts for a specific sender card number and year.
// @Summary Get monthly transfer amounts by sender card number
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer amounts for a specific sender card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Sender Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount "Monthly transfer amounts by sender card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer amounts by sender card number"
// @Router /api/transfers/monthly-amounts-by-sender-card [get]
func (h *transferHandleApi) FindMonthlyTransferAmountsBySenderCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyTransferAmountsBySenderCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferAmountsBySenderCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer amounts by sender card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferAmountsBySenderCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthAmount(res)
	h.cache.SetMonthlyTransferAmountsBySenderCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferAmountsByReceiverCardNumber retrieves the monthly transfer amounts for a specific receiver card number and year.
// @Summary Get monthly transfer amounts by receiver card number
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the monthly transfer amounts for a specific receiver card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Receiver Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount "Monthly transfer amounts by receiver card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer amounts by receiver card number"
// @Router /api/transfers/monthly-amounts-by-receiver-card [get]
func (h *transferHandleApi) FindMonthlyTransferAmountsByReceiverCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyTransferAmountsByReceiverCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTransferAmountsByReceiverCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transfer amounts by receiver card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTransferAmountsByReceiverCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferMonthAmount(res)
	h.cache.SetMonthlyTransferAmountsByReceiverCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferAmountsBySenderCardNumber retrieves the yearly transfer amounts for a specific sender card number and year.
// @Summary Get yearly transfer amounts by sender card number
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer amounts for a specific sender card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Sender Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount "Yearly transfer amounts by sender card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer amounts by sender card number"
// @Router /api/transfers/yearly-amounts-by-sender-card [get]
func (h *transferHandleApi) FindYearlyTransferAmountsBySenderCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTransferAmountsBySenderCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferAmountsBySenderCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer amounts by sender card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferAmountsBySenderCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearAmount(res)
	h.cache.SetYearlyTransferAmountsBySenderCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferAmountsByReceiverCardNumber retrieves the yearly transfer amounts for a specific receiver card number and year.
// @Summary Get yearly transfer amounts by receiver card number
// @Tags Transfer
// @Security Bearer
// @Description Retrieve the yearly transfer amounts for a specific receiver card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Receiver Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount "Yearly transfer amounts by receiver card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer amounts by receiver card number"
// @Router /api/transfers/yearly-amounts-by-receiver-card [get]
func (h *transferHandleApi) FindYearlyTransferAmountsByReceiverCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTransferAmountsByReceiverCard(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTransferAmountsByReceiverCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transfer amounts by receiver card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTransferAmountsByReceiverCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTransferYearAmount(res)
	h.cache.SetYearlyTransferAmountsByReceiverCard(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find transfers by transfer_from
// @Tags Transfer
// @Security Bearer
// @Description Retrieve a list of transfer records using the transfer_from parameter
// @Accept json
// @Produce json
// @Param transfer_from path string true "Transfer From"
// @Success 200 {object} response.ApiResponseTransfers "Transfer data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer/transfer_from/{transfer_from} [get]
func (h *transferHandleApi) FindByTransferByTransferFrom(c echo.Context) error {
	transferFrom := c.Param("transfer_from")
	if transferFrom == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTransferByFrom(ctx, transferFrom)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindTransferByTransferFrom(ctx, &pb.FindTransferByTransferFromRequest{
		TransferFrom: transferFrom,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve transfer data", zap.Error(err))
		return h.handleGrpcError(err, "FindByTransferByTransferFrom")
	}

	apiResponse := h.mapping.ToApiResponseTransfers(res)
	h.cache.SetCachedTransferByFrom(ctx, transferFrom, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find transfers by transfer_to
// @Tags Transfer
// @Security Bearer
// @Description Retrieve a list of transfer records using the transfer_to parameter
// @Accept json
// @Produce json
// @Param transfer_to path string true "Transfer To"
// @Success 200 {object} response.ApiResponseTransfers "Transfer data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer/transfer_to/{transfer_to} [get]
func (h *transferHandleApi) FindByTransferByTransferTo(c echo.Context) error {
	transferTo := c.Param("transfer_to")
	if transferTo == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTransferByTo(ctx, transferTo)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindTransferByTransferTo(ctx, &pb.FindTransferByTransferToRequest{
		TransferTo: transferTo,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve transfer data", zap.Error(err))
		return h.handleGrpcError(err, "FindByTransferByTransferTo")
	}

	apiResponse := h.mapping.ToApiResponseTransfers(res)
	h.cache.SetCachedTransferByTo(ctx, transferTo, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find active transfers
// @Tags Transfer
// @Security Bearer
// @Description Retrieve a list of active transfer records
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransfers "Active transfer data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer/active [get]
func (h *transferHandleApi) FindByActiveTransfer(c echo.Context) error {
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

	reqCache := &requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransferActiveCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActiveTransfer(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transfer data", zap.Error(err))
		return h.handleGrpcError(err, "FindByActiveTransfer")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransferDeleteAt(res)
	h.cache.SetCachedTransferActiveCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Retrieve trashed transfers
// @Tags Transfer
// @Security Bearer
// @Description Retrieve a list of trashed transfer records
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransfers "List of trashed transfer records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer/trashed [get]
func (h *transferHandleApi) FindByTrashedTransfer(c echo.Context) error {
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

	reqCache := &requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTransferTrashedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashedTransfer(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve transfer data", zap.Error(err))
		return h.handleGrpcError(err, "FindByTrashedTransfer")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTransferDeleteAt(res)
	h.cache.SetCachedTransferTrashedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Create a transfer
// @Tags Transfer
// @Security Bearer
// @Description Create a new transfer record
// @Accept json
// @Produce json
// @Param body body requests.CreateTransferRequest true "Transfer request"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Validation Error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transfer"
// @Router /api/transfer/create [post]
func (h *transferHandleApi) CreateTransfer(c echo.Context) error {
	var body requests.CreateTransferRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request body: ", zap.Error(err))

		return err
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error: ", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.CreateTransfer(ctx, &pb.CreateTransferRequest{
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	})

	if err != nil {
		h.logger.Debug("Failed to create transfer: ", zap.Error(err))

		return err
	}

	so := h.mapping.ToApiResponseTransfer(res)

	h.cache.SetCachedTransferCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update a transfer
// @Tags Transfer
// @Security Bearer
// @Description Update an existing transfer record
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param body body requests.UpdateTransferRequest true "Transfer request"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Validation Error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transfer"
// @Router /api/transfer/update/{id} [post]
func (h *transferHandleApi) UpdateTransfer(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))
		return err
	}

	var body requests.UpdateTransferRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request body: ", zap.Error(err))

		return err
	}

	body.TransferID = &idInt

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error: ", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.UpdateTransfer(ctx, &pb.UpdateTransferRequest{
		TransferId:     int32(idInt),
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	})

	if err != nil {
		h.logger.Debug("Failed to update transfer: ", zap.Error(err))

		return err
	}

	so := h.mapping.ToApiResponseTransfer(res)

	h.cache.DeleteTransferCache(ctx, idInt)
	h.cache.SetCachedTransferCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Soft delete a transfer
// @Tags Transfer
// @Security Bearer
// @Description Soft delete a transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Successfully trashed transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transfer"
// @Router /api/transfer/trash/{id} [post]
func (h *transferHandleApi) TrashTransfer(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.TrashedTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to trash transfer: ", zap.Error(err))

		return err
	}

	so := h.mapping.ToApiResponseTransferDeleteAt(res)

	h.cache.DeleteTransferCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transfer
// @Tags Transfer
// @Security Bearer
// @Description Restore a trashed transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Successfully restored transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transfer:"
// @Router /api/transfer/restore/{id} [post]
func (h *transferHandleApi) RestoreTransfer(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.RestoreTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to restore transfer: ", zap.Error(err))

		return err
	}

	so := h.mapping.ToApiResponseTransferDeleteAt(res)

	h.cache.DeleteTransferCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transfer
// @Tags Transfer
// @Security Bearer
// @Description Permanently delete a transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransferDelete "Successfully deleted transfer record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transfer:"
// @Router /api/transfer/permanent/{id} [delete]
func (h *transferHandleApi) DeleteTransferPermanent(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.DeleteTransferPermanent(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		return err
	}

	so := h.mapping.ToApiResponseTransferDelete(res)

	h.cache.DeleteTransferCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transfer
// @Tags Transfer
// @Security Bearer
// @Description Restore a trashed transfer all
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransferAll "Successfully restored transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transfer:"
// @Router /api/transfer/restore/all [post]
func (h *transferHandleApi) RestoreAllTransfer(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllTransfer(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all transfer", zap.Error(err))
		return err
	}

	h.logger.Debug("Successfully restored all transfer")

	so := h.mapping.ToApiResponseTransferAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transfer
// @Tags Transfer
// @Security Bearer
// @Description Permanently delete a transfer record all.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransferAll "Successfully deleted transfer all"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transfer:"
// @Router /api/transfer/permanent/all [post]
func (h *transferHandleApi) DeleteAllTransferPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllTransferPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all transfer", zap.Error(err))

		return err
	}

	h.logger.Debug("Successfully deleted all transfer permanently")

	so := h.mapping.ToApiResponseTransferAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *transferHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Transfer").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Transfer already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Transfer service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *transferHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *transferHandleApi) getValidationMessage(fe validator.FieldError) string {
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
