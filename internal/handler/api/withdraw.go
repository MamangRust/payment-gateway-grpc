package api

import (
	withdraw_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/withdraw"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type withdrawHandleApi struct {
	client     pb.WithdrawServiceClient
	logger     logger.LoggerInterface
	mapping    apimapper.WithdrawResponseMapper
	apihandler errors.ApiHandler
	cache      withdraw_cache.WithdrawMencache
}

func NewHandlerWithdraw(client pb.WithdrawServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.WithdrawResponseMapper, apiHandler errors.ApiHandler, cache withdraw_cache.WithdrawMencache) *withdrawHandleApi {
	withdrawHandler := &withdrawHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apihandler: apiHandler,
		cache:      cache,
	}
	routerWithdraw := router.Group("/api/withdraws")

	routerWithdraw.GET("", apiHandler.Handle("find-all-withdraws", withdrawHandler.FindAll))
	routerWithdraw.GET("/card-number/:card_number", apiHandler.Handle("find-all-withdraws-by-card-number", withdrawHandler.FindAllByCardNumber))
	routerWithdraw.GET("/:id", apiHandler.Handle("find-withdraw-by-id", withdrawHandler.FindById))
	routerWithdraw.GET("/active", apiHandler.Handle("find-active-withdraws", withdrawHandler.FindByActive))
	routerWithdraw.GET("/trashed", apiHandler.Handle("find-trashed-withdraws", withdrawHandler.FindByTrashed))

	routerWithdraw.GET("/monthly-success", apiHandler.Handle("find-monthly-withdraw-status-success", withdrawHandler.FindMonthlyWithdrawStatusSuccess))
	routerWithdraw.GET("/yearly-success", apiHandler.Handle("find-yearly-withdraw-status-success", withdrawHandler.FindYearlyWithdrawStatusSuccess))
	routerWithdraw.GET("/monthly-failed", apiHandler.Handle("find-monthly-withdraw-status-failed", withdrawHandler.FindMonthlyWithdrawStatusFailed))
	routerWithdraw.GET("/yearly-failed", apiHandler.Handle("find-yearly-withdraw-status-failed", withdrawHandler.FindYearlyWithdrawStatusFailed))

	routerWithdraw.GET("/monthly-success-by-card", apiHandler.Handle("find-monthly-withdraw-status-success-by-card", withdrawHandler.FindMonthlyWithdrawStatusSuccessByCardNumber))
	routerWithdraw.GET("/yearly-success-by-card", apiHandler.Handle("find-yearly-withdraw-status-success-by-card", withdrawHandler.FindYearlyWithdrawStatusSuccessByCardNumber))
	routerWithdraw.GET("/monthly-failed-by-card", apiHandler.Handle("find-monthly-withdraw-status-failed-by-card", withdrawHandler.FindMonthlyWithdrawStatusFailedByCardNumber))
	routerWithdraw.GET("/yearly-failed-by-card", apiHandler.Handle("find-yearly-withdraw-status-failed-by-card", withdrawHandler.FindYearlyWithdrawStatusFailedByCardNumber))

	routerWithdraw.GET("/monthly-amount", apiHandler.Handle("find-monthly-withdraw-amounts", withdrawHandler.FindMonthlyWithdraws))
	routerWithdraw.GET("/yearly-amount", apiHandler.Handle("find-yearly-withdraw-amounts", withdrawHandler.FindYearlyWithdraws))

	routerWithdraw.GET("/monthly-amount-card", apiHandler.Handle("find-monthly-withdraw-amounts-by-card", withdrawHandler.FindMonthlyWithdrawsByCardNumber))
	routerWithdraw.GET("/yearly-amount-card", apiHandler.Handle("find-yearly-withdraw-amounts-by-card", withdrawHandler.FindYearlyWithdrawsByCardNumber))

	routerWithdraw.POST("/create", apiHandler.Handle("create-withdraw", withdrawHandler.Create))
	routerWithdraw.POST("/update/:id", apiHandler.Handle("update-withdraw", withdrawHandler.Update))

	routerWithdraw.POST("/trashed/:id", apiHandler.Handle("trash-withdraw", withdrawHandler.TrashWithdraw))
	routerWithdraw.POST("/restore/:id", apiHandler.Handle("restore-withdraw", withdrawHandler.RestoreWithdraw))
	routerWithdraw.DELETE("/permanent/:id", apiHandler.Handle("delete-withdraw-permanent", withdrawHandler.DeleteWithdrawPermanent))

	routerWithdraw.POST("/restore/all", apiHandler.Handle("restore-all-withdraws", withdrawHandler.RestoreAllWithdraw))
	routerWithdraw.POST("/permanent/all", apiHandler.Handle("delete-all-withdraws-permanent", withdrawHandler.DeleteAllWithdrawPermanent))

	return withdrawHandler
}

// @Summary Find all withdraw records
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of all withdraw records with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationWithdraw "List of withdraw records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw [get]
func (h *withdrawHandleApi) FindAll(c echo.Context) error {
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

	reqCache := &requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedWithdrawsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllWithdraw(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationWithdraw(res)
	h.cache.SetCachedWithdrawsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find all withdraw records by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of withdraw records for a specific card number with pagination and search
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationWithdraw "List of withdraw records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw/card-number/{card_number} [get]
func (h *withdrawHandleApi) FindAllByCardNumber(c echo.Context) error {
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

	reqCache := &requests.FindAllWithdrawCardNumber{
		CardNumber: cardNumber,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	cachedData, found := h.cache.GetCachedWithdrawByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllWithdrawByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllWithdrawByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindAllByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponsePaginationWithdraw(res)
	h.cache.SetCachedWithdrawByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a withdraw record using its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw/{id} [get]
func (h *withdrawHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedWithdrawCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	}

	withdraw, err := h.client.FindByIdWithdraw(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseWithdraw(withdraw)
	h.cache.SetCachedWithdrawCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawStatusSuccess retrieves the monthly withdraw status for successful transactions.
// @Summary Get monthly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess "Monthly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for successful transactions"
// @Router /api/withdraws/monthly-success [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusSuccess(c echo.Context) error {
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

	reqCache := &requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetCachedMonthWithdrawStatusSuccessCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyWithdrawStatusSuccess(ctx, &pb.FindMonthlyWithdrawStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraw status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyWithdrawStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawMonthStatusSuccess(res)
	h.cache.SetCachedMonthWithdrawStatusSuccessCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawStatusSuccess retrieves the yearly withdraw status for successful transactions.
// @Summary Get yearly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for successful transactions"
// @Router /api/withdraws/yearly-success [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearlyWithdrawStatusSuccessCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyWithdrawStatusSuccess(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraw status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyWithdrawStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawYearStatusSuccess(res)
	h.cache.SetCachedYearlyWithdrawStatusSuccessCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawStatusFailed retrieves the monthly withdraw status for failed transactions.
// @Summary Get monthly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed "Monthly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for failed transactions"
// @Router /api/withdraws/monthly-failed [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusFailed(c echo.Context) error {
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

	reqCache := &requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetCachedMonthWithdrawStatusFailedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyWithdrawStatusFailed(ctx, &pb.FindMonthlyWithdrawStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraw status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyWithdrawStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawMonthStatusFailed(res)
	h.cache.SetCachedMonthWithdrawStatusFailedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawStatusFailed retrieves the yearly withdraw status for failed transactions.
// @Summary Get yearly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for failed transactions"
// @Router /api/withdraws/yearly-failed [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearlyWithdrawStatusFailedCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyWithdrawStatusFailed(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraw status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyWithdrawStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawYearStatusFailed(res)
	h.cache.SetCachedYearlyWithdrawStatusFailedCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawStatusSuccessByCardNumber retrieves the monthly withdraw status for successful transactions.
// @Summary Get monthly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess "Monthly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for successful transactions"
// @Router /api/withdraws/monthly-success-by-card [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthStatusWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyWithdrawStatusSuccessCardNumber(ctx, &pb.FindMonthlyWithdrawStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraw status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyWithdrawStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawMonthStatusSuccess(res)
	h.cache.SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawStatusSuccessByCardNumber retrieves the yearly withdraw status for successful transactions.
// @Summary Get yearly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for successful transactions"
// @Router /api/withdraws/yearly-success-by-card-number [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.YearStatusWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyWithdrawStatusSuccessCardNumber(ctx, &pb.FindYearWithdrawStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraw status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyWithdrawStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawYearStatusSuccess(res)
	h.cache.SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawStatusFailedByCardNumber retrieves the monthly withdraw status for failed transactions.
// @Summary Get monthly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed "Monthly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for failed transactions"
// @Router /api/withdraws/monthly-failed-by-card [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthStatusWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetCachedMonthWithdrawStatusFailedByCardNumber(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyWithdrawStatusFailedCardNumber(ctx, &pb.FindMonthlyWithdrawStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraw status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyWithdrawStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawMonthStatusFailed(res)
	h.cache.SetCachedMonthWithdrawStatusFailedByCardNumber(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawStatusFailedByCardNumber retrieves the yearly withdraw status for failed transactions.
// @Summary Get yearly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for failed transactions"
// @Router /api/withdraws/yearly-failed-by-card [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.YearStatusWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyWithdrawStatusFailedCardNumber(ctx, &pb.FindYearWithdrawStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraw status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyWithdrawStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawYearStatusFailed(res)
	h.cache.SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdraws retrieves the monthly withdraws for a specific year.
// @Summary Get monthly withdraws
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraws for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount "Monthly withdraws"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraws"
// @Router /api/withdraws/monthly [get]
func (h *withdrawHandleApi) FindMonthlyWithdraws(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMonthlyWithdraws(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyWithdraws(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraws", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyWithdraws")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawMonthAmount(res)
	h.cache.SetCachedMonthlyWithdraws(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdraws retrieves the yearly withdraws for a specific year.
// @Summary Get yearly withdraws
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraws for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount "Yearly withdraws"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraws"
// @Router /api/withdraws/yearly [get]
func (h *withdrawHandleApi) FindYearlyWithdraws(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedYearlyWithdraws(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyWithdraws(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraws", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyWithdraws")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawYearAmount(res)
	h.cache.SetCachedYearlyWithdraws(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawsByCardNumber retrieves the monthly withdraws for a specific card number and year.
// @Summary Get monthly withdraws by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraws for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount "Monthly withdraws by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraws by card number"
// @Router /api/withdraws/monthly-by-card [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.YearMonthCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedMonthlyWithdrawsByCardNumber(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyWithdrawsByCardNumber(ctx, &pb.FindYearWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraws by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyWithdrawsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawMonthAmount(res)
	h.cache.SetCachedMonthlyWithdrawsByCardNumber(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawsByCardNumber retrieves the yearly withdraws for a specific card number and year.
// @Summary Get yearly withdraws by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraws for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount "Yearly withdraws by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraws by card number"
// @Router /api/withdraws/yearly-by-card [get]
func (h *withdrawHandleApi) FindYearlyWithdrawsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	if cardNumber == "" {
		return errors.NewBadRequestError("invalid card_number parameter")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.YearMonthCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetCachedYearlyWithdrawsByCardNumber(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyWithdrawsByCardNumber(ctx, &pb.FindYearWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraws by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyWithdrawsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseWithdrawYearAmount(res)
	h.cache.SetCachedYearlyWithdrawsByCardNumber(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find a withdraw by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a withdraw record using its card number
// @Accept json
// @Produce json
// @Param card_number query string true "Card number"
// @Success 200 {object} response.ApiResponsesWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid card number"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraws/card/{card_number} [get]
func (h *withdrawHandleApi) FindByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")

	ctx := c.Request().Context()

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	withdraw, err := h.client.FindByCardNumber(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponsesWithdraw(withdraw)

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve all active withdraw data
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of all active withdraw data
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsesWithdraw "List of withdraw data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraws/active [get]
func (h *withdrawHandleApi) FindByActive(c echo.Context) error {
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

	reqCache := &requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedWithdrawActiveCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))
		return err
	}

	apiResponse := h.mapping.ToApiResponsePaginationWithdrawDeleteAt(res)
	h.cache.SetCachedWithdrawActiveCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Retrieve trashed withdraw data
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of trashed withdraw data
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsesWithdraw "List of trashed withdraw data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraws/trashed [get]
func (h *withdrawHandleApi) FindByTrashed(c echo.Context) error {
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

	reqCache := &requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedWithdrawTrashedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))
		return err
	}

	apiResponse := h.mapping.ToApiResponsePaginationWithdrawDeleteAt(res)
	h.cache.SetCachedWithdrawTrashedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Create a new withdraw
// @Tags Withdraw
// @Security Bearer
// @Description Create a new withdraw record with the provided details.
// @Accept json
// @Produce json
// @Param CreateWithdrawRequest body requests.CreateWithdrawRequest true "Create Withdraw Request"
// @Success 200 {object} response.ApiResponseWithdraw "Successfully created withdraw record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create withdraw"
// @Router /api/withdraws/create [post]
func (h *withdrawHandleApi) Create(c echo.Context) error {
	var body requests.CreateWithdrawRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	res, err := h.client.CreateWithdraw(ctx, &pb.CreateWithdrawRequest{
		CardNumber:     body.CardNumber,
		WithdrawAmount: int32(body.WithdrawAmount),
		WithdrawTime:   timestamppb.New(body.WithdrawTime),
	})

	if err != nil {
		h.logger.Debug("Failed to create withdraw", zap.Error(err))

		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseWithdraw(res)
	h.cache.SetCachedWithdrawCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update an existing withdraw
// @Tags Withdraw
// @Security Bearer
// @Description Update an existing withdraw record with the provided details.
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Param UpdateWithdrawRequest body requests.UpdateWithdrawRequest true "Update Withdraw Request"
// @Success 200 {object} response.ApiResponseWithdraw "Successfully updated withdraw record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update withdraw"
// @Router /api/withdraws/update/{id} [post]
func (h *withdrawHandleApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateWithdrawRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	res, err := h.client.UpdateWithdraw(ctx, &pb.UpdateWithdrawRequest{
		WithdrawId:     int32(id),
		CardNumber:     body.CardNumber,
		WithdrawAmount: int32(body.WithdrawAmount),
		WithdrawTime:   timestamppb.New(body.WithdrawTime),
	})

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseWithdraw(res)

	h.cache.DeleteCachedWithdrawCache(ctx, id)
	h.cache.SetCachedWithdrawCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Trash a withdraw using its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdaw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trash withdraw"
// @Router /api/withdraws/trashed/{id} [post]
func (h *withdrawHandleApi) TrashWithdraw(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.client.TrashedWithdraw(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseWithdrawDeleteAt(res)

	h.cache.DeleteCachedWithdrawCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Restore a withdraw by its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore withdraw"
// @Router /api/withdraws/restore/{id} [post]
func (h *withdrawHandleApi) RestoreWithdraw(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.client.RestoreWithdraw(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseWithdrawDeleteAt(res)

	h.cache.DeleteCachedWithdrawCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Permanently delete a withdraw by its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdrawDelete "Successfully deleted withdraw permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete withdraw permanently:"
// @Router /api/withdraws/permanent/{id} [delete]
func (h *withdrawHandleApi) DeleteWithdrawPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.client.DeleteWithdrawPermanent(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		return h.handleGrpcError(err, "DeleteWithdraw")
	}

	so := h.mapping.ToApiResponseWithdrawDelete(res)

	h.cache.DeleteCachedWithdrawCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a withdraw all
// @Tags Withdraw
// @Security Bearer
// @Description Restore a withdraw all
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseWithdrawAll "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore withdraw"
// @Router /api/withdraws/restore/all [post]
func (h *withdrawHandleApi) RestoreAllWithdraw(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllWithdraw(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	h.logger.Debug("Successfully restored all withdraw")

	so := h.mapping.ToApiResponseWithdrawAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Permanently delete a withdraw by its ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseWithdrawAll "Successfully deleted withdraw permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete withdraw permanently:"
// @Router /api/withdraws/permanent/all [post]
func (h *withdrawHandleApi) DeleteAllWithdrawPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllWithdrawPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	h.logger.Debug("Successfully deleted all withdraw permanently")

	so := h.mapping.ToApiResponseWithdrawAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *withdrawHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Withdraw").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Withdraw already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Withdraw service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *withdrawHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *withdrawHandleApi) getValidationMessage(fe validator.FieldError) string {
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
