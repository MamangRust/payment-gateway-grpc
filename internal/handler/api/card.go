package api

import (
	card_cache "MamangRust/paymentgatewaygrpc/internal/cache/api/card"
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

type cardHandleApi struct {
	card       pb.CardServiceClient
	logger     logger.LoggerInterface
	apiHandler errors.ApiHandler
	mapping    apimapper.CardResponseMapper
	cache      card_cache.CardMencache
}

func NewHandlerCard(card pb.CardServiceClient, router *echo.Echo, logger logger.LoggerInterface, apiHandler errors.ApiHandler, mapping apimapper.CardResponseMapper, cache card_cache.CardMencache) *cardHandleApi {
	cardHandler := &cardHandleApi{
		card:       card,
		logger:     logger,
		apiHandler: apiHandler,
		mapping:    mapping,
		cache:      cache,
	}

	routerCard := router.Group("/api/card")

	routerCard.GET("", apiHandler.Handle("find-all", cardHandler.FindAll))
	routerCard.GET("/:id", apiHandler.Handle("find-by-id", cardHandler.FindById))
	routerCard.GET("/user", apiHandler.Handle("find-by-user-id", cardHandler.FindByUserID))
	routerCard.GET("/active", apiHandler.Handle("find-by-active", cardHandler.FindByActive))
	routerCard.GET("/trashed", apiHandler.Handle("find-by-trashed", cardHandler.FindByTrashed))
	routerCard.GET("/card_number/:card_number", apiHandler.Handle("find-by-card-number", cardHandler.FindByCardNumber))

	routerCard.GET("/dashboard", apiHandler.Handle("dashboard-card", cardHandler.DashboardCard))
	routerCard.GET("/dashboard/:cardNumber", apiHandler.Handle("dashboard-card-by-card-number", cardHandler.DashboardCardCardNumber))

	routerCard.GET("/monthly-balance", apiHandler.Handle("find-monthly-balance", cardHandler.FindMonthlyBalance))
	routerCard.GET("/yearly-balance", apiHandler.Handle("find-yearly-balance", cardHandler.FindYearlyBalance))

	routerCard.GET("/monthly-topup-amount", apiHandler.Handle("find-monthly-topup-amount", cardHandler.FindMonthlyTopupAmount))
	routerCard.GET("/yearly-topup-amount", apiHandler.Handle("find-yearly-topup-amount", cardHandler.FindYearlyTopupAmount))
	routerCard.GET("/monthly-withdraw-amount", apiHandler.Handle("find-monthly-withdraw-amount", cardHandler.FindMonthlyWithdrawAmount))
	routerCard.GET("/yearly-withdraw-amount", apiHandler.Handle("find-yearly-withdraw-amount", cardHandler.FindYearlyWithdrawAmount))

	routerCard.GET("/monthly-transaction-amount", apiHandler.Handle("find-monthly-transaction-amount", cardHandler.FindMonthlyTransactionAmount))
	routerCard.GET("/yearly-transaction-amount", apiHandler.Handle("find-yearly-transaction-amount", cardHandler.FindYearlyTransactionAmount))

	routerCard.GET("/monthly-transfer-sender-amount", apiHandler.Handle("find-monthly-transfer-sender-amount", cardHandler.FindMonthlyTransferSenderAmount))
	routerCard.GET("/yearly-transfer-sender-amount", apiHandler.Handle("find-yearly-transfer-sender-amount", cardHandler.FindYearlyTransferSenderAmount))
	routerCard.GET("/monthly-transfer-receiver-amount", apiHandler.Handle("find-monthly-transfer-receiver-amount", cardHandler.FindMonthlyTransferReceiverAmount))
	routerCard.GET("/yearly-transfer-receiver-amount", apiHandler.Handle("find-yearly-transfer-receiver-amount", cardHandler.FindYearlyTransferReceiverAmount))

	routerCard.GET("/monthly-balance-by-card", apiHandler.Handle("find-monthly-balance-by-card", cardHandler.FindMonthlyBalanceByCardNumber))
	routerCard.GET("/yearly-balance-by-card", apiHandler.Handle("find-yearly-balance-by-card", cardHandler.FindYearlyBalanceByCardNumber))
	routerCard.GET("/monthly-topup-amount-by-card", apiHandler.Handle("find-monthly-topup-amount-by-card", cardHandler.FindMonthlyTopupAmountByCardNumber))
	routerCard.GET("/yearly-topup-amount-by-card", apiHandler.Handle("find-yearly-topup-amount-by-card", cardHandler.FindYearlyTopupAmountByCardNumber))

	routerCard.GET("/monthly-withdraw-amount-by-card", apiHandler.Handle("find-monthly-withdraw-amount-by-card", cardHandler.FindMonthlyWithdrawAmountByCardNumber))
	routerCard.GET("/yearly-withdraw-amount-by-card", apiHandler.Handle("find-yearly-withdraw-amount-by-card", cardHandler.FindYearlyWithdrawAmountByCardNumber))
	routerCard.GET("/monthly-transaction-amount-by-card", apiHandler.Handle("find-monthly-transaction-amount-by-card", cardHandler.FindMonthlyTransactionAmountByCardNumber))
	routerCard.GET("/yearly-transaction-amount-by-card", apiHandler.Handle("find-yearly-transaction-amount-by-card", cardHandler.FindYearlyTransactionAmountByCardNumber))

	routerCard.GET("/monthly-transfer-sender-amount-by-card", apiHandler.Handle("find-monthly-transfer-sender-amount-by-card", cardHandler.FindMonthlyTransferSenderAmountByCardNumber))
	routerCard.GET("/yearly-transfer-sender-amount-by-card", apiHandler.Handle("find-yearly-transfer-sender-amount-by-card", cardHandler.FindYearlyTransferSenderAmountByCardNumber))
	routerCard.GET("/monthly-transfer-receiver-amount-by-card", apiHandler.Handle("find-monthly-transfer-receiver-amount-by-card", cardHandler.FindMonthlyTransferReceiverAmountByCardNumber))
	routerCard.GET("/yearly-transfer-receiver-amount-by-card", apiHandler.Handle("find-yearly-transfer-receiver-amount-by-card", cardHandler.FindYearlyTransferReceiverAmountByCardNumber))

	routerCard.POST("/create", apiHandler.Handle("create-card", cardHandler.CreateCard))
	routerCard.POST("/update/:id", apiHandler.Handle("update-card", cardHandler.UpdateCard))
	routerCard.POST("/trashed/:id", apiHandler.Handle("trashed-card", cardHandler.TrashedCard))
	routerCard.POST("/restore/:id", apiHandler.Handle("restore-card", cardHandler.RestoreCard))
	routerCard.DELETE("/permanent/:id", apiHandler.Handle("delete-card-permanent", cardHandler.DeleteCardPermanent))
	routerCard.POST("/restore/all", apiHandler.Handle("restore-all-card", cardHandler.RestoreAllCard))
	routerCard.POST("/permanent/all", apiHandler.Handle("delete-all-card-permanent", cardHandler.DeleteAllCardPermanent))

	return cardHandler
}

// FindAll godoc
// @Summary Retrieve all cards
// @Tags Card
// @Security Bearer
// @Description Retrieve all cards with pagination
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Number of data per page"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCard "Card data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card data"
// @Router /api/card [get]
func (h *cardHandleApi) FindAll(c echo.Context) error {
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

	reqCache := &requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetFindAllCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	cards, err := h.card.FindAllCard(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindAllCard")
	}

	apiResponse := h.mapping.ToApiResponsesCard(cards)
	h.cache.SetFindAllCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindById godoc
// @Summary Retrieve card by ID
// @Tags Card
// @Security Bearer
// @Description Retrieve a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid card ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/{id} [get]
func (h *cardHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required and must be an integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetByIdCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindByIdCardRequest{
		CardId: int32(id),
	}

	card, err := h.card.FindByIdCard(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByIdCard")
	}

	apiResponse := h.mapping.ToApiResponseCard(card)
	h.cache.SetByIdCache(ctx, id, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindByUserID godoc
// @Summary Retrieve cards by user ID
// @Tags Card
// @Security Bearer
// @Description Retrieve a list of cards associated with a user by their ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/user [get]
func (h *cardHandleApi) FindByUserID(c echo.Context) error {
	userIDStr, ok := c.Get("userID").(string)
	if !ok {
		return errors.NewBadRequestError("user_id is required")
	}

	uid, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		return errors.NewBadRequestError("invalid user ID format")
	}
	userID := int(uid)

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetByUserIDCache(ctx, userID)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindByUserIdCardRequest{
		UserId: int32(userID),
	}

	card, err := h.card.FindByUserIdCard(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByUserIdCard")
	}

	apiResponse := h.mapping.ToApiResponseCard(card)
	h.cache.SetByUserIDCache(ctx, userID, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// DashboardCard godoc
// @Summary Get dashboard card data
// @Description Retrieve dashboard card data
// @Tags Card
// @Security Bearer
// @Produce json
// @Success 200 {object} response.ApiResponseDashboardCard
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/dashboard [get]
func (h *cardHandleApi) DashboardCard(c echo.Context) error {
	ctx := c.Request().Context()

	cachedData, found := h.cache.GetDashboardCardCache(ctx)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.card.DashboardCard(ctx, &emptypb.Empty{})
	if err != nil {
		return h.handleGrpcError(err, "DashboardCard")
	}

	apiResponse := h.mapping.ToApiResponseDashboardCard(res)
	h.cache.SetDashboardCardCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// DashboardCardCardNumber godoc
// @Summary Get dashboard card data by card number
// @Description Retrieve dashboard card data for a specific card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param cardNumber path string true "Card Number"
// @Success 200 {object} response.ApiResponseDashboardCardNumber
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/dashboard/{cardNumber} [get]
func (h *cardHandleApi) DashboardCardCardNumber(c echo.Context) error {
	ctx := c.Request().Context()

	cardNumber := c.Param("cardNumber")
	if cardNumber == "" {
		return errors.NewBadRequestError("cardNumber is required")
	}

	cachedData, found := h.cache.GetDashboardCardCardNumberCache(ctx, cardNumber)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.card.DashboardCardNumber(ctx, req)
	if err != nil {
		return h.handleGrpcError(err, "DashboardCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseDashboardCardCardNumber(res)
	h.cache.SetDashboardCardCardNumberCache(ctx, cardNumber, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyBalance godoc
// @Summary Get monthly balance data
// @Description Retrieve monthly balance data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-balance [get]
func (h *cardHandleApi) FindMonthlyBalance(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyBalanceCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearBalance{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyBalance(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyBalance")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyBalances(res)
	h.cache.SetMonthlyBalanceCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyBalance godoc
// @Summary Get yearly balance data
// @Description Retrieve yearly balance data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-balance [get]
func (h *cardHandleApi) FindYearlyBalance(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyBalanceCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearBalance{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyBalance(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyBalance")
	}

	apiResponse := h.mapping.ToApiResponseYearlyBalances(res)
	h.cache.SetYearlyBalanceCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupAmount godoc
// @Summary Get monthly topup amount data
// @Description Retrieve monthly topup amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-topup-amount [get]
func (h *cardHandleApi) FindMonthlyTopupAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTopupCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTopupAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTopupAmount")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTopupCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupAmount godoc
// @Summary Get yearly topup amount data
// @Description Retrieve yearly topup amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/topup/yearly-topup-amount [get]
func (h *cardHandleApi) FindYearlyTopupAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTopupCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTopupAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTopupAmount")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTopupCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawAmount godoc
// @Summary Get monthly withdraw amount data
// @Description Retrieve monthly withdraw amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-withdraw-amount [get]
func (h *cardHandleApi) FindMonthlyWithdrawAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyWithdrawCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyWithdrawAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyWithdrawAmount")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyWithdrawCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawAmount godoc
// @Summary Get yearly withdraw amount data
// @Description Retrieve yearly withdraw amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-withdraw-amount [get]
func (h *cardHandleApi) FindYearlyWithdrawAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyWithdrawCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyWithdrawAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyWithdrawAmount")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyWithdrawCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransactionAmount godoc
// @Summary Get monthly transaction amount data
// @Description Retrieve monthly transaction amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transaction-amount [get]
func (h *cardHandleApi) FindMonthlyTransactionAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTransactionCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransactionAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTransactionAmount")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTransactionCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransactionAmount godoc
// @Summary Get yearly transaction amount data
// @Description Retrieve yearly transaction amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transaction-amount [get]
func (h *cardHandleApi) FindYearlyTransactionAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTransactionCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransactionAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTransactionAmount")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTransactionCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferSenderAmount godoc
// @Summary Get monthly transfer sender amount data
// @Description Retrieve monthly transfer sender amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-sender-amount [get]
func (h *cardHandleApi) FindMonthlyTransferSenderAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTransferSenderCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransferSenderAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTransferSenderAmount")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTransferSenderCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferSenderAmount godoc
// @Summary Get yearly transfer sender amount data
// @Description Retrieve yearly transfer sender amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transfer/yearly-transfer-sender-amount [get]
func (h *cardHandleApi) FindYearlyTransferSenderAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTransferSenderCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransferSenderAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTransferSenderAmount")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTransferSenderCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferReceiverAmount godoc
// @Summary Get monthly transfer receiver amount data
// @Description Retrieve monthly transfer receiver amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-receiver-amount [get]
func (h *cardHandleApi) FindMonthlyTransferReceiverAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetMonthlyTransferReceiverCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransferReceiverAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTransferReceiverAmount")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTransferReceiverCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferReceiverAmount godoc
// @Summary Get yearly transfer receiver amount data
// @Description Retrieve yearly transfer receiver amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transfer-receiver-amount [get]
func (h *cardHandleApi) FindYearlyTransferReceiverAmount(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearlyTransferReceiverCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransferReceiverAmount(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTransferReceiverAmount")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTransferReceiverCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyBalanceByCardNumber godoc
// @Summary Get monthly balance data by card number
// @Description Retrieve monthly balance data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-balance-by-card [get]
func (h *cardHandleApi) FindMonthlyBalanceByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetMonthlyBalanceByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearBalanceCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyBalanceByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyBalanceByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyBalances(res)
	h.cache.SetMonthlyBalanceByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyBalanceByCardNumber godoc
// @Summary Get yearly balance data by card number
// @Description Retrieve yearly balance data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-balance-by-card [get]
func (h *cardHandleApi) FindYearlyBalanceByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetYearlyBalanceByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearBalanceCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyBalanceByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyBalanceByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseYearlyBalances(res)
	h.cache.SetYearlyBalanceByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTopupAmountByCardNumber godoc
// @Summary Get monthly topup amount data by card number
// @Description Retrieve monthly topup amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-topup-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTopupAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetMonthlyTopupByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTopupAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTopupAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTopupByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTopupAmountByCardNumber godoc
// @Summary Get yearly topup amount data by card number
// @Description Retrieve yearly topup amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-topup-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTopupAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetYearlyTopupByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTopupAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTopupAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTopupByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyWithdrawAmountByCardNumber godoc
// @Summary Get monthly withdraw amount data by card number
// @Description Retrieve monthly withdraw amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-withdraw-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyWithdrawAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetMonthlyWithdrawByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyWithdrawAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyWithdrawAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyWithdrawByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyWithdrawAmountByCardNumber godoc
// @Summary Get yearly withdraw amount data by card number
// @Description Retrieve yearly withdraw amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-withdraw-amount-by-card [get]
func (h *cardHandleApi) FindYearlyWithdrawAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetYearlyWithdrawByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyWithdrawAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyWithdrawAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyWithdrawByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransactionAmountByCardNumber godoc
// @Summary Get monthly transaction amount data by card number
// @Description Retrieve monthly transaction amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transaction-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTransactionAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetMonthlyTransactionByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransactionAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTransactionAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTransactionByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransactionAmountByCardNumber godoc
// @Summary Get yearly transaction amount data by card number
// @Description Retrieve yearly transaction amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transaction-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTransactionAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetYearlyTransactionByNumberCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransactionAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTransactionAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTransactionByNumberCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferSenderAmountByCardNumber godoc
// @Summary Get monthly transfer sender amount data by card number
// @Description Retrieve monthly transfer sender amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-sender-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTransferSenderAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetMonthlyTransferBySenderCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransferSenderAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTransferSenderAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTransferBySenderCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferSenderAmountByCardNumber godoc
// @Summary Get yearly transfer sender amount by card number
// @Description Retrieve the total amount sent by a specific card number in a given year
// @Tags Card
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transfer-sender-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTransferSenderAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetYearlyTransferBySenderCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransferSenderAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTransferSenderAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTransferBySenderCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindMonthlyTransferReceiverAmountByCardNumber godoc
// @Summary Get monthly transfer receiver amount by card number
// @Description Retrieve the total amount received by a specific card number in a given year, broken down by month
// @Tags Card
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-receiver-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTransferReceiverAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetMonthlyTransferByReceiverCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransferReceiverAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindMonthlyTransferReceiverAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseMonthlyAmounts(res)
	h.cache.SetMonthlyTransferByReceiverCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearlyTransferReceiverAmountByCardNumber godoc
// @Summary Get yearly transfer receiver amount by card number
// @Description Retrieve the total amount received by a specific card number in a given year
// @Tags Card
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transfer-receiver-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTransferReceiverAmountByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("year is required and must be a positive integer")
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		return errors.NewBadRequestError("card_number is required")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthYearCardNumberCard{
		Year:       year,
		CardNumber: cardNumber,
	}

	cachedData, found := h.cache.GetYearlyTransferByReceiverCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransferReceiverAmountByCardNumber(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindYearlyTransferReceiverAmountByCardNumber")
	}

	apiResponse := h.mapping.ToApiResponseYearlyAmounts(res)
	h.cache.SetYearlyTransferByReceiverCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active card by Saldo ID
// @Tags Card
// @Description Retrieve an active card associated with a Saldo ID
// @Accept json
// @Produce json
// @Success 200 {object} pb.ApiResponsePaginationCardDeleteAt "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid Saldo ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/active [get]
func (h *cardHandleApi) FindByActive(c echo.Context) error {
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

	reqCache := &requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetByActiveCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.card.FindByActiveCard(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByActiveCard")
	}

	apiResponse := h.mapping.ToApiResponsesCardDeletedAt(res)
	h.cache.SetByActiveCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Summary Retrieve trashed cards
// @Tags Card
// @Security Bearer
// @Description Retrieve a list of trashed cards
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCardDeleteAt "Card data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/trashed [get]
func (h *cardHandleApi) FindByTrashed(c echo.Context) error {
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

	reqCache := &requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetByTrashedCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	reqGrpc := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.card.FindByTrashedCard(ctx, reqGrpc)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashedCard")
	}

	apiResponse := h.mapping.ToApiResponsesCardDeletedAt(res)
	h.cache.SetByTrashedCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve card by card number
// @Tags Card
// @Description Retrieve a card by its card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card number"
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Failed to fetch card record"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/{card_number} [get]
func (h *cardHandleApi) FindByCardNumber(c echo.Context) error {
	cardNumber := c.Param("card_number")

	ctx := c.Request().Context()

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.card.FindByCardNumber(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to fetch card record", zap.Error(err))
		return err
	}

	so := h.mapping.ToApiResponseCard(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new card
// @Tags Card
// @Description Create a new card for a user
// @Accept json
// @Produce json
// @Param CreateCardRequest body requests.CreateCardRequest true "Create card request"
// @Success 200 {object} response.ApiResponseCard "Created card"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create card"
// @Router /api/card/create [post]
func (h *cardHandleApi) CreateCard(c echo.Context) error {
	var body requests.CreateCardRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.CreateCardRequest{
		UserId:       int32(body.UserID),
		CardType:     body.CardType,
		ExpireDate:   timestamppb.New(body.ExpireDate),
		Cvv:          body.CVV,
		CardProvider: body.CardProvider,
	}

	res, err := h.card.CreateCard(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	so := h.mapping.ToApiResponseCard(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update a card
// @Tags Card
// @Description Update a card for a user
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param UpdateCardRequest body requests.UpdateCardRequest true "Update card request"
// @Success 200 {object} response.ApiResponseCard "Updated card"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update card"
// @Router /api/card/update/{id} [post]
func (h *cardHandleApi) UpdateCard(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateCardRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request")
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateCardRequest{
		CardId:       int32(idInt),
		UserId:       int32(body.UserID),
		CardType:     body.CardType,
		ExpireDate:   timestamppb.New(body.ExpireDate),
		Cvv:          body.CVV,
		CardProvider: body.CardProvider,
	}

	res, err := h.card.UpdateCard(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	so := h.mapping.ToApiResponseCard(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Trashed a card
// @Tags Card
// @Description Trashed a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Trashed card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed card"
// @Router /api/card/trashed/{id} [post]
func (h *cardHandleApi) TrashedCard(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.TrashedCard(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Trashed")
	}

	so := h.mapping.ToApiResponseCardDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore a card
// @Tags Card
// @Description Restore a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Restored card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore card"
// @Router /api/card/restore/{id} [post]
func (h *cardHandleApi) RestoreCard(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.RestoreCard(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	so := h.mapping.ToApiResponseCardDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete a card permanently
// @Tags Card
// @Description Delete a card by its ID permanently
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCardDelete "Deleted card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete card"
// @Router /api/card/permanent/{id} [delete]
func (h *cardHandleApi) DeleteCardPermanent(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.DeleteCardPermanent(ctx, req)

	if err != nil {
		return h.handleGrpcError(err, "DeleteCard")
	}

	so := h.mapping.ToApiResponseCardDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore all card records
// @Tags Card
// @Description Restore all card records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCardAll "Successfully restored all card records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all card records"
// @Router /api/card/restore/all [post]
func (h *cardHandleApi) RestoreAllCard(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.card.RestoreAllCard(ctx, &emptypb.Empty{})
	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	h.logger.Debug("Successfully restored all cards")

	so := h.mapping.ToApiResponseCardAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer.
// @Summary Permanently delete all card records
// @Tags Card
// @Description Permanently delete all card records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCardAll "Successfully deleted all card records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all card records"
// @Router /api/card/permanent/all [post]
func (h *cardHandleApi) DeleteAllCardPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.card.DeleteAllCardPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	h.logger.Debug("Successfully deleted all cards permanently")

	so := h.mapping.ToApiResponseCardAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *cardHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Card").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Card already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Card service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *cardHandleApi) parseValidationErrors(err error) []errors.ValidationError {
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

func (h *cardHandleApi) getValidationMessage(fe validator.FieldError) string {
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
