package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardRewardRepository struct {
	db *gorm.DB
}

func NewCardRewardRepository(db *gorm.DB) CardRewardRepository {
	return &cardRewardRepository{db: db}
}

// earnPoints calculates points based on amount and MCC category.
func earnPoints(amount int64, mcc string) int32 {
	base := int32(amount / 10000)
	switch mcc {
	case "5812", "4511", "5541":
		base *= 2
	}
	if base < 1 {
		base = 1
	}
	return base
}

func (r *cardRewardRepository) EarnRewards(ctx context.Context, req *requests.EarnRewardsRequest) (*models.CardReward, error) {
	points := earnPoints(req.Amount, req.Mcc)

	reward := models.CardReward{
		CardNumber:   req.CardNumber,
		TxnID:        req.TxnID,
		Amount:       req.Amount,
		Mcc:          req.Mcc,
		PointsEarned: points,
		ExpiresAt:    time.Now().AddDate(1, 0, 0),
	}
	if err := r.db.WithContext(ctx).Create(&reward).Error; err != nil {
		return nil, card_errors.ErrEarnRewardsFailed.WithInternal(err)
	}

	return &reward, nil
}

func (r *cardRewardRepository) GetBalance(ctx context.Context, cardNumber string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(
		"SELECT COALESCE(SUM(points_earned), 0) FROM card_rewards WHERE card_number = ? AND redeemed = false AND expires_at > NOW()",
		cardNumber,
	).Scan(&total).Error
	if err != nil {
		return 0, card_errors.ErrGetRewardBalanceFailed.WithInternal(err)
	}

	return total, nil
}

func (r *cardRewardRepository) GetHistory(ctx context.Context, cardNumber string) ([]*models.CardReward, error) {
	var res []*models.CardReward
	err := r.db.WithContext(ctx).
		Where("card_number = ?", cardNumber).
		Order("created_at DESC").
		Find(&res).Error
	if err != nil {
		return nil, card_errors.ErrGetRewardHistoryFailed.WithInternal(err)
	}
	return res, nil
}

func (r *cardRewardRepository) RedeemRewards(ctx context.Context, cardNumber string, points int64) (int64, error) {
	// Find redeemable reward IDs ordered by creation date
	var rows []struct {
		RewardID     int32
		PointsEarned int32
	}
	err := r.db.WithContext(ctx).Table("card_rewards").
		Select("reward_id, points_earned").
		Where("card_number = ? AND redeemed = false AND expires_at > NOW()", cardNumber).
		Order("created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return 0, card_errors.ErrRedeemRewardsFailed.WithInternal(err)
	}

	var ids []int32
	var total int64
	for _, row := range rows {
		if total >= points {
			break
		}
		ids = append(ids, row.RewardID)
		total += int64(row.PointsEarned)
	}

	if total < points {
		return 0, sharedErrors.ErrBadRequest.WithMessage("Insufficient reward points")
	}

	// Mark rewards as redeemed
	result := r.db.WithContext(ctx).Model(&models.CardReward{}).
		Where("reward_id IN ?", ids).
		Update("redeemed", true)
	if result.Error != nil {
		return 0, card_errors.ErrRedeemRewardsFailed.WithInternal(result.Error)
	}

	return total, nil
}


