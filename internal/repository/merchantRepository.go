package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	apikey "MamangRust/paymentgatewaygrpc/pkg/api-key"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/merchant_errors"
	"context"
	"time"
)

type merchantRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantRepository {
	return &merchantRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantRepository) FindAllMerchants(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.GetMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllMerchantsFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsGetAllRecord(merchant), &totalCount, nil
}

func (r *merchantRepository) FindByActive(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindActiveMerchantsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsActiveRecord(res), &totalCount, nil
}

func (r *merchantRepository) FindByTrashed(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindTrashedMerchantsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsTrashedRecord(res), &totalCount, nil
}

func (r *merchantRepository) FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionRecord(merchant), &totalCount, nil
}

func (r *merchantRepository) FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByMerchant(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsByMerchantFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsTransactionByMerchantRecord(merchant), &totalCount, nil
}

func (r *merchantRepository) FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByApikeyParams{
		ApiKey:  req.ApiKey,
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByApikey(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsByApiKeyFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionByApikeyRecord(merchant), &totalCount, nil
}

func (r *merchantRepository) FindById(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByIdFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) FindByApiKey(api_key string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByApiKey(r.ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) FindByName(name string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByName(r.ctx, name)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByNameFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodsMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodsMerchantFailed
	}

	return r.mapping.ToMerchantMonthlyPaymentMethods(res), nil
}

func (r *merchantRepository) GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodMerchant(r.ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodMerchantFailed
	}

	return r.mapping.ToMerchantYearlyPaymentMethods(res), nil

}

func (r *merchantRepository) GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmountMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountMerchantFailed
	}

	return r.mapping.ToMerchantMonthlyAmounts(res), nil
}

func (r *merchantRepository) GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountMerchant(r.ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountMerchantFailed
	}

	return r.mapping.ToMerchantYearlyAmounts(res), nil
}

func (r *merchantRepository) GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountMerchantFailed
	}

	return r.mapping.ToMerchantMonthlyTotalAmounts(res), nil
}

func (r *merchantRepository) GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountMerchant(r.ctx, int32(year))

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountMerchantFailed
	}

	return r.mapping.ToMerchantYearlyTotalAmounts(res), nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodByMerchants(r.ctx, db.GetMonthlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByMerchantsFailed
	}

	return r.mapping.ToMerchantMonthlyPaymentMethodsByMerchant(res), nil
}

func (r *merchantRepository) GetYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByMerchants(r.ctx, db.GetYearlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByMerchantsFailed
	}

	return r.mapping.ToMerchantYearlyPaymentMethodsByMerchant(res), nil
}

func (r *merchantRepository) GetMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByMerchants(r.ctx, db.GetMonthlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantMonthlyAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByMerchants(r.ctx, db.GetYearlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantYearlyAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByMerchant(r.ctx, db.GetMonthlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantMonthlyTotalAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByMerchant(r.ctx, db.GetYearlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantYearlyTotalAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByApikey(r.ctx, db.GetMonthlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByApikeyFailed
	}

	return r.mapping.ToMerchantMonthlyPaymentMethodsByApikey(res), nil
}

func (r *merchantRepository) GetYearlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByApikey(r.ctx, db.GetYearlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByApikeyFailed
	}

	return r.mapping.ToMerchantYearlyPaymentMethodsByApikey(res), nil
}

func (r *merchantRepository) GetMonthlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByApikey(r.ctx, db.GetMonthlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByApikeyFailed
	}

	return r.mapping.ToMerchantMonthlyAmountsByApikey(res), nil
}

func (r *merchantRepository) GetYearlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByApikey(r.ctx, db.GetYearlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByApikeyFailed
	}

	return r.mapping.ToMerchantYearlyAmountsByApikey(res), nil
}

func (r *merchantRepository) GetMonthlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByApikey(r.ctx, db.GetMonthlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByApikeyFailed
	}

	return r.mapping.ToMerchantMonthlyTotalAmountsByApikey(res), nil
}

func (r *merchantRepository) GetYearlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByApikey(r.ctx, db.GetYearlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByApikeyFailed
	}

	return r.mapping.ToMerchantYearlyTotalAmountsByApikey(res), nil
}

func (r *merchantRepository) FindByMerchantUserId(user_id int) ([]*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantsByUserID(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByUserIdFailed
	}

	return r.mapping.ToMerchantsRecord(res), nil
}

func (r *merchantRepository) CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		Name:   request.Name,
		ApiKey: apikey.GenerateApiKey(),
		UserID: int32(request.UserID),
		Status: "inactive",
	}

	res, err := r.db.CreateMerchant(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrCreateMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID: int32(*request.MerchantID),
		Name:       request.Name,
		UserID:     int32(request.UserID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchant(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) UpdateMerchantStatus(request *requests.UpdateMerchantStatus) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) TrashedMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.TrashMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) RestoreMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.RestoreMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) DeleteMerchantPermanent(merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(merchant_id))

	if err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanentFailed
	}

	return true, nil
}

func (r *merchantRepository) RestoreAllMerchant() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchantFailed
	}

	return true, nil
}

func (r *merchantRepository) DeleteAllMerchantPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanentFailed
	}

	return true, nil
}
