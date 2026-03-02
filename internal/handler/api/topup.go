package api

import (
	topup_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/topup"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	apimapper "MamangRust/paymentgatewaygrpc/internal/mapper"
	"fmt"

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
)

type topupHandleApi struct {
	client     pb.TopupServiceClient
	logger     logger.LoggerInterface
	mapping    apimapper.TopupResponseMapper
	apihandler errors.ApiHandler
	cache      topup_cache.TopupMencach
}

func NewHandlerTopup(client pb.TopupServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.TopupResponseMapper, apiHandler errors.ApiHandler, cache topup_cache.TopupMencach) *topupHandleApi {
	topupHandler := &topupHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apihandler: apiHandler,
		cache:      cache,
	}
	routerTopup := router.Group("/api/topups")

	routerTopup.GET("", apiHandler.Handle("find-all-topups", topupHandler.FindAll))
	routerTopup.GET("/card-number/:card_number", apiHandler.Handle("find-all-topups-by-card-number", topupHandler.FindAllByCardNumber))
	routerTopup.GET("/:id", apiHandler.Handle("find-topup-by-id", topupHandler.FindById))
	routerTopup.GET("/active", apiHandler.Handle("find-active-topups", topupHandler.FindByActive))
	routerTopup.GET("/trashed", apiHandler.Handle("find-trashed-topups", topupHandler.FindByTrashed))

	routerTopup.GET("/monthly-success", apiHandler.Handle("find-monthly-topup-status-success", topupHandler.FindMonthlyTopupStatusSuccess))
	routerTopup.GET("/yearly-success", apiHandler.Handle("find-yearly-topup-status-success", topupHandler.FindYearlyTopupStatusSuccess))
	routerTopup.GET("/monthly-failed", apiHandler.Handle("find-monthly-topup-status-failed", topupHandler.FindMonthlyTopupStatusFailed))
	routerTopup.GET("/yearly-failed", apiHandler.Handle("find-yearly-topup-status-failed", topupHandler.FindYearlyTopupStatusFailed))

	routerTopup.GET("/monthly-success-by-card", apiHandler.Handle("find-monthly-topup-status-success-by-card", topupHandler.FindMonthlyTopupStatusSuccessByCardNumber))
	routerTopup.GET("/yearly-success-by-card", apiHandler.Handle("find-yearly-topup-status-success-by-card", topupHandler.FindYearlyTopupStatusSuccessByCardNumber))
	routerTopup.GET("/monthly-failed-by-card", apiHandler.Handle("find-monthly-topup-status-failed-by-card", topupHandler.FindMonthlyTopupStatusFailedByCardNumber))
	routerTopup.GET("/yearly-failed-by-card", apiHandler.Handle("find-yearly-topup-status-failed-by-card", topupHandler.FindYearlyTopupStatusFailedByCardNumber))

	routerTopup.GET("/monthly-methods", apiHandler.Handle("find-monthly-topup-methods", topupHandler.FindMonthlyTopupMethods))
	routerTopup.GET("/yearly-methods", apiHandler.Handle("find-yearly-topup-methods", topupHandler.FindYearlyTopupMethods))
	routerTopup.GET("/monthly-methods-by-card", apiHandler.Handle("find-monthly-topup-methods-by-card", topupHandler.FindMonthlyTopupMethodsByCardNumber))
	routerTopup.GET("/yearly-methods-by-card", apiHandler.Handle("find-yearly-topup-methods-by-card", topupHandler.FindYearlyTopupMethodsByCardNumber))

	routerTopup.GET("/monthly-amounts", apiHandler.Handle("find-monthly-topup-amounts", topupHandler.FindMonthlyTopupAmounts))
	routerTopup.GET("/yearly-amounts", apiHandler.Handle("find-yearly-topup-amounts", topupHandler.FindYearlyTopupAmounts))

	routerTopup.GET("/monthly-amounts-by-card", apiHandler.Handle("find-monthly-topup-amounts-by-card", topupHandler.FindMonthlyTopupAmountsByCardNumber))
	routerTopup.GET("/yearly-amounts-by-card", apiHandler.Handle("find-yearly-topup-amounts-by-card", topupHandler.FindYearlyTopupAmountsByCardNumber))

	routerTopup.POST("/create", apiHandler.Handle("create-topup", topupHandler.Create))
	routerTopup.POST("/update/:id", apiHandler.Handle("update-topup", topupHandler.Update))
	routerTopup.POST("/trashed/:id", apiHandler.Handle("trash-topup", topupHandler.TrashTopup))
	routerTopup.POST("/restore/:id", apiHandler.Handle("restore-topup", topupHandler.RestoreTopup))
	routerTopup.DELETE("/permanent/:id", apiHandler.Handle("delete-topup-permanent", topupHandler.DeleteTopupPermanent))

	routerTopup.POST("/delete/all", apiHandler.Handle("delete-all-topups-permanent", topupHandler.DeleteAllTopupPermanent))
	routerTopup.POST("/restore/all", apiHandler.Handle("restore-all-topups", topupHandler.RestoreAllTopup))

	return topupHandler

}

// @Tags Topup
// @Security Bearer
// @Description Retrieve a list of all topup data with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTopup "List of topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups [get]
func (h *topupHandleApi) FindAll(c echo.Context) error {
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

	reqCache := &requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTopupsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTopup(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))
		return h.handleGrpcError(err, "FindAll")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTopup(res)
	h.cache.SetCachedTopupsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find all topup by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of transactions for a specific card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTopup "List of topups"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topups data"
// @Router /api/topups/card-number/{card_number} [get]
func (h *topupHandleApi) FindAllByCardNumber(c echo.Context) error {
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

	reqCache := &requests.FindAllTopupsByCardNumber{
		CardNumber: cardNumber,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	cachedData, found := h.cache.GetCacheTopupByCardCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTopupByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllTopupByCardNumber(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))
		return h.handleGrpcError(err, "FindAllByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTopup(res)
	h.cache.SetCacheTopupByCardCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find a topup by ID
// @Tags Topup
// @Security Bearer
// @Description Retrieve a topup record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups/{id} [get]
func (h *topupHandleApi) FindById(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.NewBadRequestError("invalid id parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedTopupCache(ctx, idInt)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByIdTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))
		return h.handleGrpcError(err, "FindById")
	}

	apiResponse := h.mapping.ToApiResponseTopup(res)
	h.cache.SetCachedTopupCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Find active topups
// @Tags Topup
// @Security Bearer
// @Description Retrieve a list of active topup records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesTopup "Active topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups/active [get]
func (h *topupHandleApi) FindByActive(c echo.Context) error {
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

	reqCache := &requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTopupActiveCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTopupDeleteAt(res)
	h.cache.SetCachedTopupActiveCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Retrieve trashed topups
// @Tags Topup
// @Security Bearer
// @Description Retrieve a list of trashed topup records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesTopup "List of trashed topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups/trashed [get]
func (h *topupHandleApi) FindByTrashed(c echo.Context) error {
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

	reqCache := &requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedTopupTrashedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, reqGrpc)
	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsePaginationTopupDeleteAt(res)
	h.cache.SetCachedTopupTrashedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupStatusSuccess retrieves the monthly top-up status for successful transactions.
// @Summary Get monthly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess "Monthly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for successful transactions"
// @Router /api/topups/monthly-success [get]
func (h *topupHandleApi) FindMonthlyTopupStatusSuccess(c echo.Context) error {
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

	reqCache := &requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetMonthTopupStatusSuccessCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupStatusSuccess(ctx, &pb.FindMonthlyTopupStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthStatusSuccess(res)
	h.cache.SetMonthTopupStatusSuccessCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupStatusSuccess retrieves the yearly top-up status for successful transactions.
// @Summary Get yearly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusSuccess "Yearly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for successful transactions"
// @Router /api/topups/yearly-success [get]
func (h *topupHandleApi) FindYearlyTopupStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTopupStatusSuccessCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupStatusSuccess(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupStatusSuccess")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearStatusSuccess(res)
	h.cache.SetYearlyTopupStatusSuccessCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupStatusFailed retrieves the monthly top-up status for failed transactions.
// @Summary Get monthly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusFailed "Monthly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for failed transactions"
// @Router /api/topups/monthly-failed [get]
func (h *topupHandleApi) FindMonthlyTopupStatusFailed(c echo.Context) error {
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

	reqCache := &requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetMonthTopupStatusFailedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupStatusFailed(ctx, &pb.FindMonthlyTopupStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthStatusFailed(res)
	h.cache.SetMonthTopupStatusFailedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupStatusFailed retrieves the yearly top-up status for failed transactions.
// @Summary Get yearly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed "Yearly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for failed transactions"
// @Router /api/topups/yearly-failed [get]
func (h *topupHandleApi) FindYearlyTopupStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTopupStatusFailedCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupStatusFailed(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupStatusFailed")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearStatusFailed(res)
	h.cache.SetYearlyTopupStatusFailedCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupStatusSuccess retrieves the monthly top-up status for successful transactions.
// @Summary Get monthly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess "Monthly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for successful transactions"
// @Router /api/topups/monthly-success [get]
func (h *topupHandleApi) FindMonthlyTopupStatusSuccessByCardNumber(c echo.Context) error {
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

	reqCache := &requests.MonthTopupStatusCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetMonthTopupStatusSuccessByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTopupStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status success", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthStatusSuccess(res)
	h.cache.SetMonthTopupStatusSuccessByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupStatusSuccess retrieves the yearly top-up status for successful transactions.
// @Summary Get yearly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupYearStatusSuccess "Yearly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for successful transactions"
// @Router /api/topups/yearly-success [get]
func (h *topupHandleApi) FindYearlyTopupStatusSuccessByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearTopupStatusCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTopupStatusSuccessByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupStatusSuccessByCardNumber(ctx, &pb.FindYearTopupStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status success", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupStatusSuccessByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearStatusSuccess(res)
	h.cache.SetYearlyTopupStatusSuccessByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupStatusFailed retrieves the monthly top-up status for failed transactions.
// @Summary Get monthly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupMonthStatusFailed "Monthly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for failed transactions"
// @Router /api/topups/monthly-failed [get]
func (h *topupHandleApi) FindMonthlyTopupStatusFailedByCardNumber(c echo.Context) error {
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

	reqCache := &requests.MonthTopupStatusCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	cachedData, found := h.cache.GetMonthTopupStatusFailedByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupStatusFailedByCardNumber(ctx, &pb.FindMonthlyTopupStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthStatusFailed(res)
	h.cache.SetMonthTopupStatusFailedByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupStatusFailedByCardNumber retrieves the yearly top-up status for failed transactions.
// @Summary Get yearly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed "Yearly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for failed transactions"
// @Router /api/topups/yearly-failed [get]
func (h *topupHandleApi) FindYearlyTopupStatusFailedByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearTopupStatusCardNumber{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTopupStatusFailedByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupStatusFailedByCardNumber(ctx, &pb.FindYearTopupStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status failed", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupStatusFailedByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearStatusFailed(res)
	h.cache.SetYearlyTopupStatusFailedByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupMethods retrieves the monthly top-up methods for a specific year.
// @Summary Get monthly top-up methods
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up methods for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod "Monthly top-up methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up methods"
// @Router /api/topups/monthly-methods [get]
func (h *topupHandleApi) FindMonthlyTopupMethods(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTopupMethodsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupMethods(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup methods", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupMethods")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthMethod(res)
	h.cache.SetMonthlyTopupMethodsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupMethods retrieves the yearly top-up methods for a specific year.
// @Summary Get yearly top-up methods
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up methods for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod "Yearly top-up methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up methods"
// @Router /api/topups/yearly-methods [get]
func (h *topupHandleApi) FindYearlyTopupMethods(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTopupMethodsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupMethods(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup methods", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupMethods")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearMethod(res)
	h.cache.SetYearlyTopupMethodsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupAmounts retrieves the monthly top-up amounts for a specific year.
// @Summary Get monthly top-up amounts
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount "Monthly top-up amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up amounts"
// @Router /api/topup/monthly-amounts [get]
func (h *topupHandleApi) FindMonthlyTopupAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTopupAmountsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupAmounts(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup amounts", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupAmounts")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthAmount(res)
	h.cache.SetMonthlyTopupAmountsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupAmounts retrieves the yearly top-up amounts for a specific year.
// @Summary Get yearly top-up amounts
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount "Yearly top-up amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up amounts"
// @Router /api/topups/yearly-amounts [get]
func (h *topupHandleApi) FindYearlyTopupAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTopupAmountsCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupAmounts(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup amounts", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupAmounts")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearAmount(res)
	h.cache.SetYearlyTopupAmountsCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupMethodsByCardNumber retrieves the monthly top-up methods for a specific card number and year.
// @Summary Get monthly top-up methods by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up methods for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod "Monthly top-up methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up methods by card number"
// @Router /api/topups/monthly-methods-by-card [get]
func (h *topupHandleApi) FindMonthlyTopupMethodsByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearMonthMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyTopupMethodsByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupMethodsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup methods by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupMethodsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthMethod(res)
	h.cache.SetMonthlyTopupMethodsByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupMethodsByCardNumber retrieves the yearly top-up methods for a specific card number and year.
// @Summary Get yearly top-up methods by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up methods for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod "Yearly top-up methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up methods by card number"
// @Router /api/topups/yearly-methods-by-card [get]
func (h *topupHandleApi) FindYearlyTopupMethodsByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearMonthMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTopupMethodsByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupMethodsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup methods by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupMethodsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearMethod(res)
	h.cache.SetYearlyTopupMethodsByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupAmountsByCardNumber retrieves the monthly top-up amounts for a specific card number and year.
// @Summary Get monthly top-up amounts by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount "Monthly top-up amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up amounts by card number"
// @Router /api/topups/monthly-amounts-by-card [get]
func (h *topupHandleApi) FindMonthlyTopupAmountsByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearMonthMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyTopupAmountsByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTopupAmountsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup amounts by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTopupAmountsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupMonthAmount(res)
	h.cache.SetMonthlyTopupAmountsByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupAmountsByCardNumber retrieves the yearly top-up amounts for a specific card number and year.
// @Summary Get yearly top-up amounts by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount "Yearly top-up amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up amounts by card number"
// @Router /api/topups/yearly-amounts-by-card [get]
func (h *topupHandleApi) FindYearlyTopupAmountsByCardNumber(c echo.Context) error {
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

	reqCache := &requests.YearMonthMethod{
		CardNumber: cardNumber,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTopupAmountsByCardNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTopupAmountsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup amounts by card number", zap.Error(err))
		return h.handleGrpcError(err, "FindYearlyTopupAmountsByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseTopupYearAmount(res)
	h.cache.SetYearlyTopupAmountsByCardNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Create topup
// @Tags Topup
// @Security Bearer
// @Description Create a new topup record
// @Accept json
// @Produce json
// @Param CreateTopupRequest body requests.CreateTopupRequest true "Create topup request"
// @Success 200 {object} response.ApiResponseTopup "Created topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: "
// @Failure 500 {object} response.ErrorResponse "Failed to create topup: "
// @Router /api/topups/create [post]
func (h *topupHandleApi) Create(c echo.Context) error {
	var body requests.CreateTopupRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	res, err := h.client.CreateTopup(ctx, &pb.CreateTopupRequest{
		CardNumber:  body.CardNumber,
		TopupAmount: int32(body.TopupAmount),
		TopupMethod: body.TopupMethod,
	})

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseTopup(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update topup
// @Tags Topup
// @Security Bearer
// @Description Update an existing topup record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Param UpdateTopupRequest body requests.UpdateTopupRequest true "Update topup request"
// @Success 200 {object} response.ApiResponseTopup "Updated topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid input data"
// @Failure 500 {object} response.ErrorResponse "Failed to update topup: "
// @Router /api/topups/update/{id} [post]
func (h *topupHandleApi) Update(c echo.Context) error {
	idint, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateTopupRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	res, err := h.client.UpdateTopup(ctx, &pb.UpdateTopupRequest{
		TopupId:     int32(idint),
		CardNumber:  body.CardNumber,
		TopupAmount: int32(body.TopupAmount),
		TopupMethod: body.TopupMethod,
	})

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseTopup(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a topup
// @Tags Topup
// @Security Bearer
// @Description Trash a topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Successfully trashed topup record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trash topup:"
// @Router /api/topups/trash/{id} [post]
func (h *topupHandleApi) TrashTopup(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.client.TrashedTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseTopupDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed topup
// @Tags Topup
// @Security Bearer
// @Description Restore a trashed topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Successfully restored topup record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore topup:"
// @Router /api/topups/restore/{id} [post]
func (h *topupHandleApi) RestoreTopup(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.client.RestoreTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseTopupDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a topup
// @Tags Topup
// @Security Bearer
// @Description Permanently delete a topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopupDelete "Successfully deleted topup record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete topup:"
// @Router /api/topups/permanent/{id} [delete]
func (h *topupHandleApi) DeleteTopupPermanent(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.client.DeleteTopupPermanent(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "DeleteTopup")
	}

	so := h.mapping.ToApiResponseTopupDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore all topup records
// @Tags Topup
// @Security Bearer
// @Description Restore all topup records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTopupAll "Successfully restored all topup records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all topup records"
// @Router /api/topups/restore/all [post]
func (h *topupHandleApi) RestoreAllTopup(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllTopup(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	h.logger.Debug("Successfully restored all topup")

	so := h.mapping.ToApiResponseTopupAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete all topup records
// @Tags Topup
// @Security Bearer
// @Description Permanently delete all topup records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTopupAll "Successfully deleted all topup records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all topup records"
// @Router /api/topups/permanent/all [post]
func (h *topupHandleApi) DeleteAllTopupPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllTopupPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	h.logger.Debug("Successfully deleted all topup permanently")

	so := h.mapping.ToApiResponseTopupAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *topupHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Topup").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Topup already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Topup service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *topupHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *topupHandleApi) getValidationMessage(fe validator.FieldError) string {
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
