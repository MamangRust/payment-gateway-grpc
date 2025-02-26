package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	apikey "MamangRust/paymentgatewaygrpc/pkg/api-key"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
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

func (r *merchantRepository) FindAllMerchants(search string, page, pageSize int) ([]*record.MerchantRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetMerchantsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.GetMerchants(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find merchants: %w", err)
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsGetAllRecord(merchant), totalCount, nil
}

func (r *merchantRepository) FindById(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) FindByApiKey(api_key string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByApiKey(r.ctx, api_key)

	if err != nil {
		return nil, fmt.Errorf("failed to merchant by api-key :%w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) FindByName(name string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByName(r.ctx, name)

	if err != nil {
		return nil, fmt.Errorf("failed to find merchant by name: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodsMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment methods for merchant: %w", err)
	}
	return r.mapping.ToMerchantMonthlyPaymentMethods(res), nil
}

func (r *merchantRepository) GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodMerchant(r.ctx, year)

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly payment methods for merchant: %w", err)
	}
	return r.mapping.ToMerchantYearlyPaymentMethods(res), nil

}

func (r *merchantRepository) GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountMerchant(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amount for merchant: %w", err)
	}
	return r.mapping.ToMerchantMonthlyAmounts(res), nil
}

func (r *merchantRepository) GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountMerchant(r.ctx, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amount for merchant: %w", err)
	}
	return r.mapping.ToMerchantYearlyAmounts(res), nil
}

func (r *merchantRepository) GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountMerchant(r.ctx, yearStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amount for merchant: %w", err)
	}

	return r.mapping.ToMerchantMonthlyTotalAmounts(res), nil
}

func (r *merchantRepository) GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountMerchant(r.ctx, int32(year))
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amount for merchant: %w", err)
	}
	return r.mapping.ToMerchantYearlyTotalAmounts(res), nil
}

func (r *merchantRepository) FindAllTransactions(search string, page, pageSize int) ([]*record.MerchantTransactionsRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.FindAllTransactionsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactions(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find merchants: %w", err)
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionRecord(merchant), totalCount, nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodByMerchants(merchantID int, year int) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByMerchants(r.ctx, db.GetMonthlyPaymentMethodByMerchantsParams{
		MerchantID: int32(merchantID),
		Column1:    yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment methods for merchant %d: %w", merchantID, err)
	}
	return r.mapping.ToMerchantMonthlyPaymentMethodsByMerchant(res), nil
}

func (r *merchantRepository) GetYearlyPaymentMethodByMerchants(merchantID int, year int) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByMerchants(r.ctx, db.GetYearlyPaymentMethodByMerchantsParams{
		MerchantID: int32(merchantID),
		Column2:    year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly payment methods for merchant %d: %w", merchantID, err)
	}
	return r.mapping.ToMerchantYearlyPaymentMethodsByMerchant(res), nil
}

func (r *merchantRepository) GetMonthlyAmountByMerchants(merchantID int, year int) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByMerchants(r.ctx, db.GetMonthlyAmountByMerchantsParams{
		MerchantID: int32(merchantID),
		Column1:    yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amount for merchant %d: %w", merchantID, err)
	}
	return r.mapping.ToMerchantMonthlyAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetYearlyAmountByMerchants(merchantID int, year int) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByMerchants(r.ctx, db.GetYearlyAmountByMerchantsParams{
		MerchantID: int32(merchantID),
		Column2:    year,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amount for merchant %d: %w", merchantID, err)
	}

	return r.mapping.ToMerchantYearlyAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetMonthlyTotalAmountByMerchants(merchantID int, year int) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByMerchant(r.ctx, db.GetMonthlyTotalAmountByMerchantParams{
		Column2: int32(merchantID),
		Column1: yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amount for merchant %d: %w", merchantID, err)
	}
	return r.mapping.ToMerchantMonthlyTotalAmountsByMerchant(res), nil
}

func (r *merchantRepository) GetYearlyTotalAmountByMerchants(merchantID int, year int) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByMerchant(r.ctx, db.GetYearlyTotalAmountByMerchantParams{
		Column2: int32(merchantID),
		Column1: int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amount for merchant %d: %w", merchantID, err)
	}

	return r.mapping.ToMerchantYearlyTotalAmountsByMerchant(res), nil
}

func (r *merchantRepository) FindAllTransactionsByMerchant(merchant_id int, search string, page, pageSize int) ([]*record.MerchantTransactionsRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.FindAllTransactionsByMerchantParams{
		MerchantID: int32(merchant_id),
		Column2:    search,
		Limit:      int32(pageSize),
		Offset:     int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByMerchant(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find merchantFindAllTransactionsByMerchantRows: %w", err)
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionByMerchantRecord(merchant), totalCount, nil
}

func (r *merchantRepository) GetMonthlyPaymentMethodByApikey(api_key string, year int) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByApikey(r.ctx, db.GetMonthlyPaymentMethodByApikeyParams{
		ApiKey:  api_key,
		Column1: yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly payment methods for merchant %s: %w", api_key, err)
	}
	return r.mapping.ToMerchantMonthlyPaymentMethodsByApikey(res), nil
}

func (r *merchantRepository) GetYearlyPaymentMethodByApikey(api_key string, year int) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByApikey(r.ctx, db.GetYearlyPaymentMethodByApikeyParams{
		ApiKey:  api_key,
		Column2: year,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly payment methods for merchant %s: %w", api_key, err)
	}
	return r.mapping.ToMerchantYearlyPaymentMethodsByApikey(res), nil
}

func (r *merchantRepository) GetMonthlyAmountByApikey(api_key string, year int) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByApikey(r.ctx, db.GetMonthlyAmountByApikeyParams{
		ApiKey:  api_key,
		Column1: yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amount for merchant %s: %w", api_key, err)
	}
	return r.mapping.ToMerchantMonthlyAmountsByApikey(res), nil
}

func (r *merchantRepository) GetYearlyAmountByApikey(api_key string, year int) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByApikey(r.ctx, db.GetYearlyAmountByApikeyParams{
		ApiKey:  api_key,
		Column2: year,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amount for merchant %s: %w", api_key, err)
	}

	return r.mapping.ToMerchantYearlyAmountsByApikey(res), nil
}

func (r *merchantRepository) GetMonthlyTotalAmountByApikey(api_key string, year int) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByApikey(r.ctx, db.GetMonthlyTotalAmountByApikeyParams{
		ApiKey:  api_key,
		Column1: yearStart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly amount for merchant %s: %w", api_key, err)
	}
	return r.mapping.ToMerchantMonthlyTotalAmountsByApikey(res), nil
}

func (r *merchantRepository) GetYearlyTotalAmountByApikey(api_key string, year int) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByApikey(r.ctx, db.GetYearlyTotalAmountByApikeyParams{
		ApiKey:  api_key,
		Column1: int32(year),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get yearly amount for merchant %s: %w", api_key, err)
	}

	return r.mapping.ToMerchantYearlyTotalAmountsByApikey(res), nil
}

func (r *merchantRepository) FindAllTransactionsByApikey(api_key string, search string, page, pageSize int) ([]*record.MerchantTransactionsRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.FindAllTransactionsByApikeyParams{
		ApiKey:  api_key,
		Column2: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByApikey(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find merchantFindAllTransactionsByMerchantRows: %w", err)
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionByApikeyRecord(merchant), totalCount, nil
}

func (r *merchantRepository) FindByMerchantUserId(user_id int) ([]*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantsByUserID(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find merchants by user_id: %w", err)
	}

	return r.mapping.ToMerchantsRecord(res), nil
}

func (r *merchantRepository) FindByActive(search string, page, pageSize int) ([]*record.MerchantRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetActiveMerchantsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveMerchants(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find active merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsActiveRecord(res), totalCount, nil
}

func (r *merchantRepository) FindByTrashed(search string, page, pageSize int) ([]*record.MerchantRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTrashedMerchantsParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedMerchants(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trashed merchant: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsTrashedRecord(res), totalCount, nil
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
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID: int32(request.MerchantID),
		Name:       request.Name,
		UserID:     int32(request.UserID),
		Status:     request.Status,
	}

	err := r.db.UpdateMerchant(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update merchant: %w", err)
	}

	res, err := r.db.GetMerchantByID(r.ctx, int32(request.MerchantID))

	if err != nil {
		return nil, fmt.Errorf("failed to find merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) UpdateMerchantStatus(request *requests.UpdateMerchantStatus) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(request.MerchantID),
		Status:     request.Status,
	}

	err := r.db.UpdateMerchantStatus(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update Merchant amount :%w", err)
	}

	res, err := r.db.GetMerchantByID(r.ctx, req.MerchantID)

	if err != nil {
		return nil, fmt.Errorf("failed to find Merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantRepository) TrashedMerchant(merchantId int) (*record.MerchantRecord, error) {
	err := r.db.TrashMerchant(r.ctx, int32(merchantId))

	if err != nil {
		return nil, fmt.Errorf("failed to trash merchant: %w", err)
	}

	merchant, err := r.db.GetTrashedMerchantByID(r.ctx, int32(merchantId))

	if err != nil {
		return nil, fmt.Errorf("failed to find trashed by id merchant: %w", err)
	}

	return r.mapping.ToMerchantRecord(merchant), nil
}

func (r *merchantRepository) RestoreMerchant(merchant_id int) (*record.MerchantRecord, error) {
	err := r.db.RestoreMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore merchant: %w", err)
	}

	merchant, err := r.db.GetMerchantByID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, fmt.Errorf("failed not found card :%w", err)
	}

	return r.mapping.ToMerchantRecord(merchant), nil
}

func (r *merchantRepository) DeleteMerchantPermanent(merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(merchant_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete merchant permanently: %w", err)
	}

	return true, nil
}

func (r *merchantRepository) RestoreAllMerchant() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all merchants: %w", err)
	}

	return true, nil
}

func (r *merchantRepository) DeleteAllMerchantPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all merchants permanently: %w", err)
	}

	return true, nil
}
