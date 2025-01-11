package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	apikey "MamangRust/paymentgatewaygrpc/pkg/api-key"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"fmt"
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

// func (r *merchantRepository) GetMonthlyPaymentMethodsMerchant() {
// 	res, err := r.db.GetMonthlyPaymentMethodsMerchant(r.ctx)

// }

// func (r *merchantRepository) GetYearlyPaymentMethodMerchant() {
// 	res, err := r.db.GetYearlyPaymentMethodMerchant(r.ctx)
// }

// func (r *merchantRepository) GetMonthlyAmountMerchant() {
// 	res, err := r.db.GetMonthlyAmountMerchant(r.ctx)
// }

// func (r *merchantRepository) GetYearlyAmountMerchant() {
// 	res, err := r.db.GetYearlyAmountMerchant(r.ctx)

// }

// func (r *merchantRepository) GetMonthlyPaymentMethodByMerchants() {
// 	res, err := r.db.GetMonthlyPaymentMethodByMerchants(r.ctx)

// }

// func (r *merchantRepository) GetYearlyPaymentMethodByMerchants() {
// 	res, err := r.db.GetYearlyPaymentMethodByMerchants(r.ctx)
// }

// func (r *merchantRepository) GetMonthlyAmountByMerchants() {
// 	res, err := r.db.GetMonthlyAmountByMerchants(r.ctx)
// }

// func (r *merchantRepository) GetYearlyAmountByMerchants() {
// 	res, err := r.db.GetYearlyAmountByMerchants(r.ctx)
// }

// func (r *merchantRepository) FindAllTransactionsByMerchant(){
// 	res, err := r.db.FindAllTransactionsByMerchant(r.ctx, )
// }

func (r *merchantRepository) CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		Name:   request.Name,
		ApiKey: apikey.GenerateApiKey(),
		UserID: int32(request.UserID),
		Status: "active",
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
