package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	apikey "MamangRust/paymentgatewaygrpc/pkg/api-key"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"context"
	"time"
)

type merchantRepository struct {
	db *db.Queries
}

func NewMerchantRepository(db *db.Queries) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

func (r *merchantRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetMerchantsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.GetMerchants(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindAllMerchantsFailed
	}

	return merchant, nil
}

func (r *merchantRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetActiveMerchantsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveMerchants(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindActiveMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*db.GetTrashedMerchantsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedMerchants(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*db.FindAllTransactionsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactions(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsFailed
	}

	return merchant, nil
}

func (r *merchantRepository) FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*db.FindAllTransactionsByMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByMerchant(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsByMerchantFailed
	}

	return merchant, nil
}

func (r *merchantRepository) FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*db.FindAllTransactionsByApikeyRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByApikeyParams{
		ApiKey:  req.ApiKey,
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByApikey(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindAllTransactionsByApiKeyFailed
	}

	return merchant, nil
}

func (r *merchantRepository) FindById(ctx context.Context, merchant_id int) (*db.GetMerchantByIDRow, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByIdFailed
	}

	return res, nil
}

func (r *merchantRepository) FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error) {
	res, err := r.db.GetMerchantByApiKey(ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return res, nil
}

func (r *merchantRepository) FindByName(ctx context.Context, name string) (*db.GetMerchantByNameRow, error) {
	res, err := r.db.GetMerchantByName(ctx, name)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByNameFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*db.GetMonthlyPaymentMethodsMerchantRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodsMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodsMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*db.GetYearlyPaymentMethodMerchantRow, error) {
	res, err := r.db.GetYearlyPaymentMethodMerchant(ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodMerchantFailed
	}

	return res, nil

}

func (r *merchantRepository) GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyAmountMerchantRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmountMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyAmountMerchantRow, error) {
	res, err := r.db.GetYearlyAmountMerchant(ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetMonthlyTotalAmountMerchantRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountMerchant(ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*db.GetYearlyTotalAmountMerchantRow, error) {
	res, err := r.db.GetYearlyTotalAmountMerchant(ctx, int32(year))

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetMonthlyPaymentMethodByMerchantsRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodByMerchants(ctx, db.GetMonthlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*db.GetYearlyPaymentMethodByMerchantsRow, error) {
	res, err := r.db.GetYearlyPaymentMethodByMerchants(ctx, db.GetYearlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetMonthlyAmountByMerchantsRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByMerchants(ctx, db.GetMonthlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*db.GetYearlyAmountByMerchantsRow, error) {
	res, err := r.db.GetYearlyAmountByMerchants(ctx, db.GetYearlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetMonthlyTotalAmountByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByMerchant(ctx, db.GetMonthlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*db.GetYearlyTotalAmountByMerchantRow, error) {
	res, err := r.db.GetYearlyTotalAmountByMerchant(ctx, db.GetYearlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByMerchantsFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetMonthlyPaymentMethodByApikeyRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByApikey(ctx, db.GetMonthlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByApikeyFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*db.GetYearlyPaymentMethodByApikeyRow, error) {
	res, err := r.db.GetYearlyPaymentMethodByApikey(ctx, db.GetYearlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByApikeyFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetMonthlyAmountByApikeyRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByApikey(ctx, db.GetMonthlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByApikeyFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*db.GetYearlyAmountByApikeyRow, error) {
	res, err := r.db.GetYearlyAmountByApikey(ctx, db.GetYearlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByApikeyFailed
	}

	return res, nil
}

func (r *merchantRepository) GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetMonthlyTotalAmountByApikeyRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByApikey(ctx, db.GetMonthlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByApikeyFailed
	}

	return res, nil
}

func (r *merchantRepository) GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*db.GetYearlyTotalAmountByApikeyRow, error) {
	res, err := r.db.GetYearlyTotalAmountByApikey(ctx, db.GetYearlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByApikeyFailed
	}

	return res, nil
}

func (r *merchantRepository) FindByMerchantUserId(ctx context.Context, user_id int) ([]*db.GetMerchantsByUserIDRow, error) {
	res, err := r.db.GetMerchantsByUserID(ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByUserIdFailed
	}

	return res, nil
}

func (r *merchantRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error) {
	req := db.CreateMerchantParams{
		Name:   request.Name,
		ApiKey: apikey.GenerateApiKey(),
		UserID: int32(request.UserID),
		Status: "inactive",
	}

	res, err := r.db.CreateMerchant(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrCreateMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error) {
	req := db.UpdateMerchantParams{
		MerchantID: int32(*request.MerchantID),
		Name:       request.Name,
		UserID:     int32(request.UserID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchant(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatus) (*db.UpdateMerchantStatusRow, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed
	}

	return res, nil
}

func (r *merchantRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.TrashMerchant(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchantFailed
	}

	return res, nil
}

func (r *merchantRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.RestoreMerchant(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchantFailed
	}

	return res, nil
}


func (r *merchantRepository) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(merchant_id))

	if err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanentFailed
	}

	return true, nil
}

func (r *merchantRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchantFailed
	}

	return true, nil
}

func (r *merchantRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanentFailed
	}

	return true, nil
}
