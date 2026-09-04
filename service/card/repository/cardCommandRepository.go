package repository

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardCommandRepository struct {
	db *gorm.DB
}

func NewCardCommandRepository(db *gorm.DB) CardCommandRepository {
	return &cardCommandRepository{db: db}
}

func (r *cardCommandRepository) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*models.CardCreateRow, error) {
	number, err := randomvcc.RandomCardNumber()
	if err != nil {
		return nil, card_errors.ErrCreateCardFailed.WithInternal(err)
	}
	var result models.CardCreateRow
	err = r.db.WithContext(ctx).Raw(`INSERT INTO cards (user_id, card_number, card_type, expire_date, cvv, card_provider, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, 'active', current_timestamp, current_timestamp) RETURNING card_id, user_id, card_number, card_type, expire_date, cvv, card_provider, status, created_at, updated_at`,
		request.UserID, number, request.CardType, request.ExpireDate, request.CVV, request.CardProvider).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrCreateCardFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*models.CardUpdateRow, error) {
	var result models.CardUpdateRow
	err := r.db.WithContext(ctx).Raw(`UPDATE cards SET card_type = ?, expire_date = ?, cvv = ?, card_provider = ?, updated_at = current_timestamp WHERE card_id = ? AND deleted_at IS NULL RETURNING card_id, card_number, card_type`,
		request.CardType, request.ExpireDate, request.CVV, request.CardProvider, request.CardID).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrUpdateCardFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) TrashedCard(ctx context.Context, cardId int) (*models.CardTrashRow, error) {
	var result models.CardTrashRow
	err := r.db.WithContext(ctx).Raw(`UPDATE cards SET deleted_at = current_timestamp WHERE card_id = ? AND deleted_at IS NULL RETURNING card_id, user_id, card_number, card_type, expire_date, cvv, card_provider, status, created_at, updated_at, deleted_at`, cardId).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrTrashCardFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) RestoreCard(ctx context.Context, cardId int) (*models.CardRestoreRow, error) {
	var result models.CardRestoreRow
	err := r.db.WithContext(ctx).Raw(`UPDATE cards SET deleted_at = NULL WHERE card_id = ? AND deleted_at IS NOT NULL RETURNING card_id, user_id, card_number, card_type, expire_date, cvv, card_provider, status, created_at, updated_at, deleted_at`, cardId).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrRestoreCardFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) DeleteCardPermanent(ctx context.Context, cardId int) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM cards WHERE card_id = ? AND deleted_at IS NOT NULL`, cardId).Error
	if err != nil {
		return false, card_errors.ErrDeleteCardPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *cardCommandRepository) RestoreAllCard(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`UPDATE cards SET deleted_at = NULL WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, card_errors.ErrRestoreAllCardsFailed.WithInternal(err)
	}
	return true, nil
}

func (r *cardCommandRepository) DeleteAllCardPermanent(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Exec(`DELETE FROM cards WHERE deleted_at IS NOT NULL`).Error
	if err != nil {
		return false, card_errors.ErrDeleteAllCardsPermanentFailed.WithInternal(err)
	}
	return true, nil
}

func (r *cardCommandRepository) ToggleCardStatus(ctx context.Context, request *requests.ToggleCardStatusRequest) (*models.CardAllFieldsRow, error) {
	var result models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE cards SET status = CASE WHEN status = 'active' THEN 'inactive' ELSE 'active' END, updated_at = current_timestamp WHERE card_id = ? AND deleted_at IS NULL RETURNING card_id, user_id, card_number, card_type, expire_date, cvv, card_provider, status, outstanding_balance, credit_limit, reward_points, created_at, updated_at`,
		request.CardID).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrUpdateCardFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) UpdateCreditLimit(ctx context.Context, request *requests.UpdateCreditLimitRequest) (*models.CardAllFieldsRow, error) {
	var result models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE cards SET credit_limit = ?, updated_at = current_timestamp WHERE card_id = ? AND deleted_at IS NULL RETURNING card_id, user_id, card_number, card_type, expire_date, cvv, card_provider, status, outstanding_balance, credit_limit, reward_points, created_at, updated_at`,
		int32(request.CreditLimit), request.CardID).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrUpdateCreditLimitFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) RedeemPoints(ctx context.Context, request *requests.RedeemPointsRequest) (*models.CardAllFieldsRow, error) {
	var result models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Raw(`UPDATE cards SET reward_points = reward_points - ?, updated_at = current_timestamp WHERE card_id = ? AND deleted_at IS NULL AND reward_points >= ? RETURNING card_id, user_id, card_number, card_type, expire_date, cvv, card_provider, status, outstanding_balance, credit_limit, reward_points, created_at, updated_at`,
		int32(request.Points), request.CardID, int32(request.Points)).Scan(&result).Error
	if err != nil {
		return nil, card_errors.ErrRedeemRewardPointsFailed.WithInternal(err)
	}
	return &result, nil
}

func (r *cardCommandRepository) ProcessBillingCycles(ctx context.Context, billingCycleDay int) ([]*models.BillingCycle, error) {
	var results []*models.BillingCycle

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert billing cycles for credit cards with outstanding balance
		err := tx.Raw(`
			INSERT INTO billing_cycles (card_number, cycle_start, cycle_end, amount_due, due_date, status)
			SELECT
				c.card_number,
				date_trunc('month', CURRENT_DATE),
				date_trunc('month', CURRENT_DATE) + interval '1 month' - interval '1 day',
				c.outstanding_balance,
				CURRENT_DATE,
				'unpaid'
			FROM cards c
			WHERE c.card_type = 'credit'
			  AND c.status = 'active'
			  AND c.deleted_at IS NULL
			  AND c.outstanding_balance > 0
			  AND ? BETWEEN 1 AND 31
			  AND EXTRACT(DAY FROM CURRENT_DATE) = ?
			ON CONFLICT (card_number, cycle_start, cycle_end) DO NOTHING
			RETURNING billing_id, card_number, cycle_start, cycle_end, amount_due, due_date, status, created_at, updated_at
		`, billingCycleDay, billingCycleDay).Scan(&results).Error
		if err != nil {
			return err
		}

		// Queue outbox events for newly inserted cycles
		for _, c := range results {
			err = tx.Exec(`
				INSERT INTO outbox_events (event_key, topic, message_key, payload)
				VALUES (?, 'card.statement.generated', ?, jsonb_build_object('billing_id', ?::int, 'card_number', ?::text, 'billing_cycle_day', ?::int, 'cycle_start', ?::timestamp, 'cycle_end', ?::timestamp, 'amount_due', ?::int))
				ON CONFLICT (event_key) DO NOTHING
			`,
				fmt.Sprintf("card-statement:%d", c.BillingID),
				fmt.Sprintf("%d", c.BillingID),
				c.BillingID, c.CardNumber, billingCycleDay, c.CycleStart, c.CycleEnd, c.AmountDue,
			).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, card_errors.ErrGetBillingCyclesFailed.WithInternal(err)
	}
	return results, nil
}
