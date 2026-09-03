package repository

import (
	"fmt"
	"context"
	"errors"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cardPaymentRepository struct {
	db *gorm.DB
}

func NewCardPaymentRepository(db *gorm.DB) CardPaymentRepository {
	return &cardPaymentRepository{
		db: db,
	}
}

func (r *cardPaymentRepository) PostPaymentIdempotent(ctx context.Context, req *requests.PostPaymentRequest, expectedBillingStatus, targetBillingStatus string) (*CardPaymentResult, error) {
	if req.ReferenceID != "" {
		existing, err := r.GetPaymentByReferenceID(ctx, req.ReferenceID)
		if err == nil && existing != nil {
			if validationErr := validatePaymentReplay(existing, req); validationErr != nil {
				return nil, validationErr
			}
			return &CardPaymentResult{Payment: existing, Replayed: true}, nil
		}
		if err != nil {
		// Not found is OK, other errors are not
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrGetCardPaymentFailed.WithInternal(err)
		}
		}
	}

	var billingID *int32
	if req.BillingID != nil {
		billingID = new(int32)
		*billingID = int32(*req.BillingID)
	}

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, card_errors.ErrInsertCardPaymentFailed.WithInternal(tx.Error)
	}
	defer func() { _ = tx.Rollback() }()

	if billingID != nil {
		// Lock the billing row
		var cycle models.BillingCycle
		if lockErr := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("billing_id = ?", *billingID).First(&cycle).Error; lockErr != nil {
			return nil, card_errors.ErrGetBillingCyclesFailed.WithInternal(lockErr)
		}
		if cycle.Status != expectedBillingStatus {
			_ = tx.Rollback()
			if req.ReferenceID != "" {
				existing, lookupErr := r.GetPaymentByReferenceID(ctx, req.ReferenceID)
				if lookupErr == nil {
					if validationErr := validatePaymentReplay(existing, req); validationErr != nil {
						return nil, validationErr
					}
					return &CardPaymentResult{Payment: existing, Replayed: true}, nil
				}
			}
			return nil, sharedErrors.NewBadRequestError("billing cycle is not in a payable state")
		}
		if req.Amount != int64(cycle.AmountDue) {
			return nil, sharedErrors.NewBadRequestError("payment amount must equal billing amount due")
		}
		if updateErr := tx.Model(&models.BillingCycle{}).
			Where("billing_id = ? AND status = ?", *billingID, expectedBillingStatus).
			Update("status", targetBillingStatus).Error; updateErr != nil {
			return nil, sharedErrors.NewBadRequestError("billing cycle is not in a payable state")
		}
	}

	payment := models.CardPayment{
		PaymentUuid:    uuid.New().String(),
		CardNumber:     req.CardNumber,
		BillingID:      billingID,
		Amount:         req.Amount,
		PaymentChannel: req.PaymentChannel,
		ReferenceID:    req.ReferenceID,
		Status:         "completed",
	}

	if err := tx.Create(&payment).Error; err != nil {
		if isUniqueViolation(err) && req.ReferenceID != "" {
			_ = tx.Rollback()
			existing, lookupErr := r.GetPaymentByReferenceID(ctx, req.ReferenceID)
			if lookupErr == nil {
				if validationErr := validatePaymentReplay(existing, req); validationErr != nil {
					return nil, validationErr
				}
				return &CardPaymentResult{Payment: existing, Replayed: true}, nil
			}
		}
		return nil, card_errors.ErrInsertCardPaymentFailed.WithInternal(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, card_errors.ErrInsertCardPaymentFailed.WithInternal(err)
	}

	return &CardPaymentResult{Payment: &payment}, nil
}

func isUniqueViolation(err error) bool {
	var stateErr interface{ SQLState() string }
	return errors.As(err, &stateErr) && stateErr.SQLState() == "23505"
}

func validatePaymentReplay(existing *models.CardPayment, req *requests.PostPaymentRequest) error {
	if existing == nil || existing.CardNumber != req.CardNumber || existing.Amount != req.Amount || existing.PaymentChannel != req.PaymentChannel {
		return sharedErrors.ErrConflict.WithMessage("card payment reference already used with different payload")
	}
	if req.BillingID == nil {
		if existing.BillingID != nil {
			return sharedErrors.ErrConflict.WithMessage("card payment reference already used with different billing cycle")
		}
		return nil
	}
	if existing.BillingID == nil || *existing.BillingID != int32(*req.BillingID) {
		return sharedErrors.ErrConflict.WithMessage("card payment reference already used with different billing cycle")
	}
	return nil
}

func (r *cardPaymentRepository) GetPaymentByReferenceID(ctx context.Context, referenceID string) (*models.CardPayment, error) {
	var res models.CardPayment
	err := r.db.WithContext(ctx).Where("reference_id = ?", referenceID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, card_errors.ErrGetCardPaymentFailed.WithInternal(err)
	}
	return &res, nil
}

func (r *cardPaymentRepository) PostPayment(ctx context.Context, req *requests.PostPaymentRequest) (*models.CardPayment, error) {
	expectedBillingStatus, targetBillingStatus := "", ""
	if req.BillingID != nil {
		expectedBillingStatus = "unpaid"
		targetBillingStatus = "paid"
	}
	result, err := r.PostPaymentIdempotent(ctx, req, expectedBillingStatus, targetBillingStatus)
	if err != nil {
		return nil, err
	}
	return result.Payment, nil
}

func (r *cardPaymentRepository) GetPaymentHistory(ctx context.Context, cardNumber string, page, pageSize int) ([]*models.CardPaymentListRow, error) {
	offset := (page - 1) * pageSize

	var totalCount int64
	_ = r.db.WithContext(ctx).Table("card_payments").
		Where("card_number = ?", cardNumber).
		Count(&totalCount).Error

	var res []*models.CardPaymentListRow
	err := r.db.WithContext(ctx).Table("card_payments").
		Select(fmt.Sprintf("card_payments.*, %d as total_count", totalCount)).
		Where("card_number = ?", cardNumber).
		Order("created_at DESC").
		Limit(pageSize).Offset(int(math.Max(0, float64(offset)))).
		Scan(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetCardPaymentFailed.WithInternal(err)
	}

	return res, nil
}

func (r *cardPaymentRepository) CountPayments(ctx context.Context, cardNumber string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("card_payments").
		Where("card_number = ?", cardNumber).
		Count(&count).Error
	if err != nil {
		return 0, card_errors.ErrGetCardPaymentFailed.WithInternal(err)
	}
	return int(count), nil
}
