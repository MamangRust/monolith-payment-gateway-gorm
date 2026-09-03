package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardAuthTransactionRepository struct {
	db *gorm.DB
}

func NewCardAuthTransactionRepository(db *gorm.DB) CardAuthTransactionRepository {
	return &cardAuthTransactionRepository{db: db}
}

func (r *cardAuthTransactionRepository) InsertPending(ctx context.Context, req *requests.AuthorizeCardRequest) (*models.CardAuthTransaction, error) {
	res := models.CardAuthTransaction{
		TxnID:          req.IdempotencyKey + "-txn",
		CardNumber:     req.CardNumber,
		MerchantID:     int32(req.MerchantID),
		Amount:         req.Amount,
		Currency:       req.Currency,
		Mcc:            req.Mcc,
		PosEntryMode:   req.PosEntryMode,
		IdempotencyKey: req.IdempotencyKey,
		Status:         "pending",
	}
	if err := r.db.WithContext(ctx).Create(&res).Error; err != nil {
		return nil, card_errors.ErrInsertAuthTransactionFailed.WithInternal(err)
	}
	return &res, nil
}

func (r *cardAuthTransactionRepository) Approve(ctx context.Context, txnID string) (*models.CardAuthTransaction, error) {
	var res models.CardAuthTransaction
	err := r.db.WithContext(ctx).Model(&models.CardAuthTransaction{}).
		Where("txn_id = ? AND status = 'pending'", txnID).
		Updates(map[string]interface{}{
			"status":     "approved",
			"updated_at": time.Now(),
		}).First(&res).Error
	if err != nil {
		return nil, card_errors.ErrApproveAuthTransactionFailed.WithInternal(err)
	}

	// Re-fetch to get updated fields
	if err := r.db.WithContext(ctx).Where("txn_id = ?", txnID).First(&res).Error; err != nil {
		return nil, card_errors.ErrApproveAuthTransactionFailed.WithInternal(err)
	}

	return &res, nil
}

func (r *cardAuthTransactionRepository) Decline(ctx context.Context, txnID string) (*models.CardAuthTransaction, error) {
	var res models.CardAuthTransaction
	err := r.db.WithContext(ctx).Model(&models.CardAuthTransaction{}).
		Where("txn_id = ?", txnID).
		Updates(map[string]interface{}{
			"status":     "declined",
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return nil, card_errors.ErrDeclineAuthTransactionFailed.WithInternal(err)
	}

	if err := r.db.WithContext(ctx).Where("txn_id = ?", txnID).First(&res).Error; err != nil {
		return nil, card_errors.ErrDeclineAuthTransactionFailed.WithInternal(err)
	}

	return &res, nil
}

func (r *cardAuthTransactionRepository) Reverse(ctx context.Context, txnID string) (*models.CardAuthTransaction, error) {
	var res models.CardAuthTransaction
	err := r.db.WithContext(ctx).Model(&models.CardAuthTransaction{}).
		Where("txn_id = ?", txnID).
		Updates(map[string]interface{}{
			"status":     "reversed",
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return nil, card_errors.ErrReverseAuthTransactionFailed.WithInternal(err)
	}

	if err := r.db.WithContext(ctx).Where("txn_id = ?", txnID).First(&res).Error; err != nil {
		return nil, card_errors.ErrReverseAuthTransactionFailed.WithInternal(err)
	}

	return &res, nil
}

func (r *cardAuthTransactionRepository) FindByIdempotencyKey(ctx context.Context, key string) (*models.CardAuthTransaction, error) {
	var res models.CardAuthTransaction
	err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetAuthTransactionFailed.WithInternal(err)
	}
	return &res, nil
}

func (r *cardAuthTransactionRepository) FindByTxnID(ctx context.Context, txnID string) (*models.CardAuthTransaction, error) {
	var res models.CardAuthTransaction
	err := r.db.WithContext(ctx).Where("txn_id = ?", txnID).First(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetAuthTransactionFailed.WithInternal(err)
	}
	return &res, nil
}

func (r *cardAuthTransactionRepository) FindByCardNumber(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.GetAuthTxnByCardNumberRow, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var res []*models.GetAuthTxnByCardNumberRow
	err := r.db.WithContext(ctx).Table("card_auth_transactions").
		Where("card_number = ?", cardNumber).
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetAuthTransactionFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardAuthTransactionRepository) CountRecentByCardNumber(ctx context.Context, cardNumber string, since time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("card_auth_transactions").
		Where("card_number = ? AND created_at >= ?", cardNumber, since).
		Count(&count).Error
	if err != nil {
		return 0, card_errors.ErrGetAuthTransactionFailed.WithInternal(err)
	}
	return int(count), nil
}

func (r *cardAuthTransactionRepository) UpdateRiskScore(ctx context.Context, txnID string, score int) error {
	err := r.db.WithContext(ctx).Model(&models.CardAuthTransaction{}).
		Where("txn_id = ?", txnID).
		Update("risk_score", int32(score)).Error
	if err != nil {
		return card_errors.ErrUpdateRiskScoreFailed.WithInternal(err)
	}
	return nil
}
