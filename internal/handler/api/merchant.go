package api

import (
	merchant_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/merchant"
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

type merchantHandleApi struct {
	merchant   pb.MerchantServiceClient
	logger     logger.LoggerInterface
	apiHandler errors.ApiHandler
	mapping    apimapper.MerchantResponseMapper
	cache      merchant_cache.MerchantMencache
}

func NewHandlerMerchant(merchant pb.MerchantServiceClient, router *echo.Echo, logger logger.LoggerInterface, apiHandler errors.ApiHandler, mapping apimapper.MerchantResponseMapper, cache merchant_cache.MerchantMencache) *merchantHandleApi {
	merchantHandler := &merchantHandleApi{
		merchant:   merchant,
		logger:     logger,
		apiHandler: apiHandler,
		mapping:    mapping,
		cache:      cache,
	}

	routerMerchant := router.Group("/api/merchants")

	routerMerchant.GET("", apiHandler.Handle("find-all-merchants", merchantHandler.FindAll))
	routerMerchant.GET("/:id", apiHandler.Handle("find-merchant-by-id", merchantHandler.FindById))
	routerMerchant.GET("/api-key", apiHandler.Handle("find-merchant-by-apikey", merchantHandler.FindByApiKey))
	routerMerchant.GET("/merchant-user", apiHandler.Handle("find-merchant-by-user-id", merchantHandler.FindByMerchantUserId))
	routerMerchant.GET("/active", apiHandler.Handle("find-active-merchants", merchantHandler.FindByActive))
	routerMerchant.GET("/trashed", apiHandler.Handle("find-trashed-merchants", merchantHandler.FindByTrashed))

	routerMerchant.GET("/transactions", apiHandler.Handle("find-all-transactions", merchantHandler.FindAllTransactions))
	routerMerchant.GET("/transactions/:merchant_id", apiHandler.Handle("find-all-transactions-by-merchant", merchantHandler.FindAllTransactionByMerchant))
	routerMerchant.GET("/transactions/api-key/:api_key", apiHandler.Handle("find-all-transactions-by-apikey", merchantHandler.FindAllTransactionByApikey))

	routerMerchant.GET("/monthly-payment-methods", apiHandler.Handle("find-monthly-payment-methods", merchantHandler.FindMonthlyPaymentMethodsMerchant))
	routerMerchant.GET("/yearly-payment-methods", apiHandler.Handle("find-yearly-payment-methods", merchantHandler.FindYearlyPaymentMethodMerchant))
	routerMerchant.GET("/monthly-amount", apiHandler.Handle("find-monthly-amount", merchantHandler.FindMonthlyAmountMerchant))
	routerMerchant.GET("/yearly-amount", apiHandler.Handle("find-yearly-amount", merchantHandler.FindYearlyAmountMerchant))
	routerMerchant.GET("/monthly-total-amount", apiHandler.Handle("find-monthly-total-amount", merchantHandler.FindMonthlyTotalAmountMerchant))
	routerMerchant.GET("/yearly-total-amount", apiHandler.Handle("find-yearly-total-amount", merchantHandler.FindYearlyTotalAmountMerchant))

	routerMerchant.GET("/monthly-payment-methods-by-merchant", apiHandler.Handle("find-monthly-payment-methods-by-merchant", merchantHandler.FindMonthlyPaymentMethodByMerchants))
	routerMerchant.GET("/yearly-payment-methods-by-merchant", apiHandler.Handle("find-yearly-payment-methods-by-merchant", merchantHandler.FindYearlyPaymentMethodByMerchants))
	routerMerchant.GET("/monthly-amount-by-merchant", apiHandler.Handle("find-monthly-amount-by-merchant", merchantHandler.FindMonthlyAmountByMerchants))
	routerMerchant.GET("/yearly-amount-by-merchant", apiHandler.Handle("find-yearly-amount-by-merchant", merchantHandler.FindYearlyAmountByMerchants))
	routerMerchant.GET("/monthly-totalamount-by-merchant", apiHandler.Handle("find-monthly-total-amount-by-merchant", merchantHandler.FindMonthlyTotalAmountByMerchants))
	routerMerchant.GET("/yearly-totalamount-by-merchant", apiHandler.Handle("find-yearly-total-amount-by-merchant", merchantHandler.FindYearlyTotalAmountByMerchants))

	routerMerchant.GET("/monthly-payment-methods-by-apikey", apiHandler.Handle("find-monthly-payment-methods-by-apikey", merchantHandler.FindMonthlyPaymentMethodByApikeys))
	routerMerchant.GET("/yearly-payment-methods-by-apikey", apiHandler.Handle("find-yearly-payment-methods-by-apikey", merchantHandler.FindYearlyPaymentMethodByApikeys))
	routerMerchant.GET("/monthly-amount-by-apikey", apiHandler.Handle("find-monthly-amount-by-apikey", merchantHandler.FindMonthlyAmountByApikeys))
	routerMerchant.GET("/yearly-amount-by-apikey", apiHandler.Handle("find-yearly-amount-by-apikey", merchantHandler.FindYearlyAmountByApikeys))
	routerMerchant.GET("/monthly-totalamount-by-apikey", apiHandler.Handle("find-monthly-total-amount-by-apikey", merchantHandler.FindMonthlyTotalAmountByApikeys))
	routerMerchant.GET("/yearly-totalamount-by-apikey", apiHandler.Handle("find-yearly-total-amount-by-apikey", merchantHandler.FindYearlyTotalAmountByApikeys))

	routerMerchant.POST("/create", apiHandler.Handle("create-merchant", merchantHandler.Create))
	routerMerchant.POST("/updates/:id", apiHandler.Handle("update-merchant", merchantHandler.Update))

	routerMerchant.POST("/trashed/:id", apiHandler.Handle("trash-merchant", merchantHandler.TrashedMerchant))
	routerMerchant.POST("/restore/:id", apiHandler.Handle("restore-merchant", merchantHandler.RestoreMerchant))
	routerMerchant.DELETE("/permanent/:id", apiHandler.Handle("delete-merchant-permanent", merchantHandler.Delete))

	routerMerchant.POST("/restore/all", apiHandler.Handle("restore-all-merchants", merchantHandler.RestoreAllMerchant))
	routerMerchant.POST("/permanent/all", apiHandler.Handle("delete-all-merchants-permanent", merchantHandler.DeleteAllMerchantPermanent))

	return merchantHandler
}

// FindAll godoc
// @Summary Find all merchants
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of all merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchant "List of merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants [get]
func (h *merchantHandleApi) FindAll(c echo.Context) error {
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

	reqCache := &requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedMerchants(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindAllMerchant(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindAllMerchant")
	}

	apiResponse := h.mapping.ToApiResponsesMerchant(res)

	h.cache.SetCachedMerchants(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindAllTransactions godoc
// @Summary Find all transactions
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/merchants/transaction [get]
func (h *merchantHandleApi) FindAllTransactions(c echo.Context) error {
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

	cacheReq := &requests.FindAllMerchantTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCacheAllMerchantTransactions(ctx, cacheReq)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindAllTransactionMerchant(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseMerchantsTransactionResponse(res)

	h.cache.SetCacheAllMerchantTransactions(ctx, cacheReq, so)

	return c.JSON(http.StatusOK, so)
}

// FindAllTransactionByMerchant godoc
// @Summary Find all transactions by merchant ID
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of transactions for a specific merchant
// @Accept json
// @Produce json
// @Param merchant_id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/merchants/transactions/:merchant_id [get]
func (h *merchantHandleApi) FindAllTransactionByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return err
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

	cacheReq := &requests.FindAllMerchantTransactionsById{
		MerchantID: merchantID,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	cachedData, found := h.cache.GetCacheMerchantTransactions(ctx, cacheReq)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllMerchantTransaction{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.merchant.FindAllTransactionByMerchant(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseMerchantsTransactionResponse(res)

	h.cache.SetCacheMerchantTransactions(ctx, cacheReq, so)

	return c.JSON(http.StatusOK, so)
}

// FindById godoc
// @Summary Find a merchant by ID
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a merchant by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/{id} [get]
func (h *merchantHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required and must be an integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMerchant(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindByIdMerchantRequest{
		MerchantId: int32(id),
	}

	res, err := h.merchant.FindByIdMerchant(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByIdMerchant")
	}

	apiResponse := h.mapping.ToApiResponseMerchant(res)

	h.cache.SetCachedMerchant(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyPaymentMethodsMerchant godoc
// @Summary Find monthly payment methods for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchants/monthly-payment-methods [get]
func (h *merchantHandleApi) FindMonthlyPaymentMethodsMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyPaymentMethodsMerchantCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindMonthlyPaymentMethodsMerchant(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyPaymentMethodsMerchant")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyPaymentMethods(res)
	h.cache.SetMonthlyPaymentMethodsMerchantCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyPaymentMethodMerchant godoc.
// @Summary Find yearly payment methods for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchants/monthly-amount [get]
func (h *merchantHandleApi) FindYearlyPaymentMethodMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyPaymentMethodMerchantCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindYearlyPaymentMethodMerchant(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyPaymentMethodMerchant")
	}

	apiResponse := h.mapping.ToApiResponseYearlyPaymentMethods(res)
	h.cache.SetYearlyPaymentMethodMerchantCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountMerchant godoc
// @Summary Find monthly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-amount [get]
func (h *merchantHandleApi) FindMonthlyAmountMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyAmountMerchantCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindMonthlyAmountMerchant(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyAmountMerchant")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyAmountMerchantCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountMerchant godoc.
// @Summary Find yearly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-amount [get]
func (h *merchantHandleApi) FindYearlyAmountMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyAmountMerchantCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindYearlyAmountMerchant(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyAmountMerchant")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyAmountMerchantCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountMerchant godoc
// @Summary Find monthly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-total-amount [get]
func (h *merchantHandleApi) FindMonthlyTotalAmountMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTotalAmountMerchantCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindMonthlyTotalAmountMerchant(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalAmountMerchant")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyTotalAmounts(res)
	h.cache.SetMonthlyTotalAmountMerchantCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountMerchant godoc.
// @Summary Find yearly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-total-amount [get]
func (h *merchantHandleApi) FindYearlyTotalAmountMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTotalAmountMerchantCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindYearlyTotalAmountMerchant(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTotalAmountMerchant")
	}

	apiResponse := h.mapping.ToApiResponseYearlyTotalAmounts(res)
	h.cache.SetYearlyTotalAmountMerchantCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyPaymentMethodByMerchants godoc.
// @Summary Find monthly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchants/monthly-payment-methods-by-merchant [get]
func (h *merchantHandleApi) FindMonthlyPaymentMethodByMerchants(c echo.Context) error {
	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant_id is required and must be a positive integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethodMerchant{
		MerchantID: merchantID,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyPaymentMethodByMerchantsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindMonthlyPaymentMethodByMerchants(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyPaymentMethodByMerchants")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyPaymentMethods(res)
	h.cache.SetMonthlyPaymentMethodByMerchantsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyPaymentMethodByMerchants godoc.
// @Summary Find yearly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchants/yearly-payment-methods-by-merchant [get]
func (h *merchantHandleApi) FindYearlyPaymentMethodByMerchants(c echo.Context) error {
	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant_id is required and must be a positive integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethodMerchant{
		MerchantID: merchantID,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyPaymentMethodByMerchantsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindYearlyPaymentMethodByMerchants(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyPaymentMethodByMerchants")
	}

	apiResponse := h.mapping.ToApiResponseYearlyPaymentMethods(res)
	h.cache.SetYearlyPaymentMethodByMerchantsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountByMerchants godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-amount-by-merchant [get]
func (h *merchantHandleApi) FindMonthlyAmountByMerchants(c echo.Context) error {
	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant_id is required and must be a positive integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearAmountMerchant{
		MerchantID: merchantID,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyAmountByMerchantsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindMonthlyAmountByMerchants(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyAmountByMerchants")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyAmountByMerchantsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountByMerchants godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-amount-by-merchant [get]
func (h *merchantHandleApi) FindYearlyAmountByMerchants(c echo.Context) error {
	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant_id is required and must be a positive integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearAmountMerchant{
		MerchantID: merchantID,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyAmountByMerchantsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindYearlyAmountByMerchants(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyAmountByMerchants")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyAmountByMerchantsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountByMerchants godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-totalamount-by-merchant [get]
func (h *merchantHandleApi) FindMonthlyTotalAmountByMerchants(c echo.Context) error {
	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant_id is required and must be a positive integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearTotalAmountMerchant{
		MerchantID: merchantID,
		Year:       year,
	}

	cachedData, found := h.cache.GetMonthlyTotalAmountByMerchantsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindMonthlyTotalAmountByMerchants(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalAmountByMerchants")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyTotalAmounts(res)
	h.cache.SetMonthlyTotalAmountByMerchantsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountByMerchants godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-totalamount-by-merchant [get]
func (h *merchantHandleApi) FindYearlyTotalAmountByMerchants(c echo.Context) error {
	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("merchant_id is required and must be a positive integer")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearTotalAmountMerchant{
		MerchantID: merchantID,
		Year:       year,
	}

	cachedData, found := h.cache.GetYearlyTotalAmountByMerchantsCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindYearlyTotalAmountByMerchants(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTotalAmountByMerchants")
	}

	apiResponse := h.mapping.ToApiResponseYearlyTotalAmounts(res)
	h.cache.SetYearlyTotalAmountByMerchantsCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindAllTransactionByApikey godoc
// @Summary Find all transactions by api_key
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of transactions for a specific merchant
// @Accept json
// @Produce json
// @Param api_key path string true "Api key"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/merchants/transactions/api-key/:api_key [get]
func (h *merchantHandleApi) FindAllTransactionByApikey(c echo.Context) error {
	api_key := c.Param("api_key")

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

	cacheReq := &requests.FindAllMerchantTransactionsByApiKey{
		ApiKey:   api_key,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCacheMerchantTransactionApikey(ctx, cacheReq)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindAllMerchantApikey{
		ApiKey:   api_key,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindAllTransactionByApikey(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseMerchantsTransactionResponse(res)

	h.cache.SetCacheMerchantTransactionApikey(ctx, cacheReq, so)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodByApikeys godoc.
// @Summary Find monthly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchants/monthly-payment-methods-by-apikey [get]
func (h *merchantHandleApi) FindMonthlyPaymentMethodByApikeys(c echo.Context) error {
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	if api_key == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethodApiKey{
		Apikey: api_key,
		Year:   year,
	}

	cachedData, found := h.cache.GetMonthlyPaymentMethodByApikeysCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindMonthlyPaymentMethodByApikey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyPaymentMethodByApikey")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyPaymentMethods(res)
	h.cache.SetMonthlyPaymentMethodByApikeysCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyPaymentMethodByApikeys godoc.
// @Summary Find yearly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchants/yearly-payment-methods-by-apikey [get]
func (h *merchantHandleApi) FindYearlyPaymentMethodByApikeys(c echo.Context) error {
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	if api_key == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearPaymentMethodApiKey{
		Apikey: api_key,
		Year:   year,
	}

	cachedData, found := h.cache.GetYearlyPaymentMethodByApikeysCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindYearlyPaymentMethodByApikey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyPaymentMethodByApikey")
	}

	apiResponse := h.mapping.ToApiResponseYearlyPaymentMethods(res)
	h.cache.SetYearlyPaymentMethodByApikeysCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountByApikeys godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-amount-by-apikey [get]
func (h *merchantHandleApi) FindMonthlyAmountByApikeys(c echo.Context) error {
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	if api_key == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearAmountApiKey{
		Apikey: api_key,
		Year:   year,
	}

	cachedData, found := h.cache.GetMonthlyAmountByApikeysCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindMonthlyAmountByApikey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyAmountByApikey")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyAmountByApikeysCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountByApikeys godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-amount-by-apikey [get]
func (h *merchantHandleApi) FindYearlyAmountByApikeys(c echo.Context) error {
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	if api_key == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearAmountApiKey{
		Apikey: api_key,
		Year:   year,
	}

	cachedData, found := h.cache.GetYearlyAmountByApikeysCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindYearlyAmountByApikey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyAmountByApikey")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyAmountByApikeysCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyAmountByApikeys godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-totalamount-by-apikey [get]
func (h *merchantHandleApi) FindMonthlyTotalAmountByApikeys(c echo.Context) error {
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	if api_key == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearTotalAmountApiKey{
		Apikey: api_key,
		Year:   year,
	}

	cachedData, found := h.cache.GetMonthlyTotalAmountByApikeysCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindMonthlyTotalAmountByApikey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTotalAmountByApikey")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyTotalAmounts(res)
	h.cache.SetMonthlyTotalAmountByApikeysCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyAmountByApikeys godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-totalamount-by-apikey [get]
func (h *merchantHandleApi) FindYearlyTotalAmountByApikeys(c echo.Context) error {
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	if api_key == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearTotalAmountApiKey{
		Apikey: api_key,
		Year:   year,
	}

	cachedData, found := h.cache.GetYearlyTotalAmountByApikeysCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindYearlyTotalAmountByApikey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTotalAmountByApikey")
	}

	apiResponse := h.mapping.ToApiResponseYearlyTotalAmounts(res)
	h.cache.SetYearlyTotalAmountByApikeysCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindByApiKey godoc
// @Summary Find a merchant by API key
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a merchant by its API key
// @Accept json
// @Produce json
// @Param api_key query string true "API key"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/api-key [get]
func (h *merchantHandleApi) FindByApiKey(c echo.Context) error {
	apiKey := c.QueryParam("api_key")
	if apiKey == "" {
		return errors.NewBadRequestError("api_key is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMerchantByApiKey(ctx, apiKey)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindByApiKeyRequest{
		ApiKey: apiKey,
	}

	res, err := h.merchant.FindByApiKey(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByApiKey")
	}

	apiResponse := h.mapping.ToApiResponseMerchant(res)

	h.cache.SetCachedMerchantByApiKey(ctx, apiKey, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindByMerchantUserId godoc.
// @Summary Find a merchant by user ID
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a merchant by its user ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponsesMerchant "Merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/merchant-user [get]
func (h *merchantHandleApi) FindByMerchantUserId(c echo.Context) error {
	userId, ok := c.Get("user_id").(int32)
	if !ok {
		return errors.NewBadRequestError("user_id is required and must be valid")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedMerchantsByUserId(ctx, int(userId))
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindByMerchantUserIdRequest{
		UserId: userId,
	}

	res, err := h.merchant.FindByMerchantUserId(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByMerchantUserId")
	}

	apiResponse := h.mapping.ToApiResponseMerchants(res)

	h.cache.SetCachedMerchantsByUserId(ctx, int(userId), apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindByActive godoc
// @Summary Find active merchants
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of active merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesMerchant "List of active merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/active [get]
func (h *merchantHandleApi) FindByActive(c echo.Context) error {
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

	reqCache := &requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedMerchantActive(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindByActive(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	apiResponse := h.mapping.ToApiResponsesMerchantDeleteAt(res)

	h.cache.SetCachedMerchantActive(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindByTrashed godoc
// @Summary Find trashed merchants
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of trashed merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesMerchant "List of trashed merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/trashed [get]
func (h *merchantHandleApi) FindByTrashed(c echo.Context) error {
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

	reqCache := &requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedMerchantTrashed(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindByTrashed(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	apiResponse := h.mapping.ToApiResponsesMerchantDeleteAt(res)

	h.cache.SetCachedMerchantTrashed(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// Create godoc
// @Summary Create a new merchant
// @Tags Merchant
// @Security Bearer
// @Description Create a new merchant with the given name and user ID
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantRequest true "Create merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Created merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant"
// @Router /api/merchants/create [post]
func (h *merchantHandleApi) Create(c echo.Context) error {
	var body requests.CreateMerchantRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantRequest{
		Name:   body.Name,
		UserId: int32(body.UserID),
	}

	res, err := h.merchant.CreateMerchant(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// Update godoc
// @Summary Update a merchant
// @Tags Merchant
// @Security Bearer
// @Description Update a merchant with the given ID
// @Accept json
// @Produce json
// @Param body body requests.UpdateMerchantRequest true "Update merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Updated merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant"
// @Router /api/merchants/update/{id} [post]
func (h *merchantHandleApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateMerchantRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()
	req := &pb.UpdateMerchantRequest{
		MerchantId: int32(id),
		Name:       body.Name,
		UserId:     int32(body.UserID),
		Status:     body.Status,
	}

	res, err := h.merchant.UpdateMerchant(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// TrashedMerchant godoc
// @Summary Trashed a merchant
// @Tags Merchant
// @Security Bearer
// @Description Trashed a merchant by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed merchant"
// @Router /api/merchants/trashed/{id} [post]
func (h *merchantHandleApi) TrashedMerchant(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.merchant.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseMerchantDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// RestoreMerchant godoc
// @Summary Restore a merchant
// @Tags Merchant
// @Security Bearer
// @Description Restore a merchant by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Restored merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchants/restore/{id} [post]
func (h *merchantHandleApi) RestoreMerchant(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.merchant.RestoreMerchant(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseMerchantDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// Delete godoc
// @Summary Delete a merchant permanently
// @Tags Merchant
// @Security Bearer
// @Description Delete a merchant by its ID permanently
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Deleted merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant"
// @Router /api/merchants/{id} [delete]
func (h *merchantHandleApi) Delete(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	res, err := h.merchant.DeleteMerchantPermanent(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		return h.handleGrpcError(err, "DeletePermanent")
	}

	so := h.mapping.ToApiResponseMerchantDelete(res)

	return c.JSON(http.StatusOK, so)
}

// RestoreAllMerchant godoc.
// @Summary Restore all merchant records
// @Tags Merchant
// @Security Bearer
// @Description Restore all merchant records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all merchant records"
// @Router /api/merchants/restore/all [post]
func (h *merchantHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.merchant.RestoreAllMerchant(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	h.logger.Debug("Successfully restored all merchant")

	so := h.mapping.ToApiResponseMerchantAll(res)

	return c.JSON(http.StatusOK, so)
}

// DeleteAllMerchantPermanent godoc.
// @Summary Permanently delete all merchant records
// @Tags Merchant
// @Security Bearer
// @Description Permanently delete all merchant records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all merchant records"
// @Router /api/merchants/permanent/all [post]
func (h *merchantHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.merchant.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	h.logger.Debug("Successfully deleted all merchant permanently")

	so := h.mapping.ToApiResponseMerchantAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Merchant").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Merchant already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Merchant service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *merchantHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *merchantHandleApi) getValidationMessage(fe validator.FieldError) string {
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
