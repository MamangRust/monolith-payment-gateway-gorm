package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"gorm.io/gorm"
)

type cardQueryRepository struct {
	db *gorm.DB
}

func NewCardQueryRepository(db *gorm.DB) CardQueryRepository {
	return &cardQueryRepository{
		db: db,
	}
}

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (r *cardQueryRepository) FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListRow, error) {
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("cards").Where("deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(card_number ILIKE ? OR card_type ILIKE ? OR card_provider ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, card_errors.ErrFindAllCardsFailed.WithInternal(err)
	}

	var cards []*models.CardListRow
	query := r.db.WithContext(ctx).Table("cards").
		Select(fmt.Sprintf("cards.*, %d as total_count", totalCount)).
		Where("cards.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(cards.card_number ILIKE ? OR cards.card_type ILIKE ? OR cards.card_provider ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if err := query.Order("cards.card_id DESC").Limit(pageSize).Offset(offset).Scan(&cards).Error; err != nil {
		return nil, card_errors.ErrFindAllCardsFailed.WithInternal(err)
	}

	return cards, nil
}

func (r *cardQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, error) {
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("cards").Where("deleted_at IS NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(card_number ILIKE ? OR card_type ILIKE ? OR card_provider ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, card_errors.ErrFindActiveCardsFailed.WithInternal(err)
	}

	var cards []*models.CardListWithDeletedRow
	query := r.db.WithContext(ctx).Table("cards").
		Select(fmt.Sprintf("cards.*, %d as total_count", totalCount)).
		Where("cards.deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("(cards.card_number ILIKE ? OR cards.card_type ILIKE ? OR cards.card_provider ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if err := query.Order("cards.card_id DESC").Limit(pageSize).Offset(offset).Scan(&cards).Error; err != nil {
		return nil, card_errors.ErrFindActiveCardsFailed.WithInternal(err)
	}

	return cards, nil
}

func (r *cardQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*models.CardListWithDeletedRow, error) {
	page, pageSize := normalizePagination(req.Page, req.PageSize)
	offset := (page - 1) * pageSize

	var totalCount int64
	countQuery := r.db.WithContext(ctx).Table("cards").Where("deleted_at IS NOT NULL")
	if req.Search != "" {
		countQuery = countQuery.Where("(card_number ILIKE ? OR card_type ILIKE ? OR card_provider ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, card_errors.ErrFindTrashedCardsFailed.WithInternal(err)
	}

	var cards []*models.CardListWithDeletedRow
	query := r.db.WithContext(ctx).Table("cards").
		Select(fmt.Sprintf("cards.*, %d as total_count", totalCount)).
		Where("cards.deleted_at IS NOT NULL")
	if req.Search != "" {
		query = query.Where("(cards.card_number ILIKE ? OR cards.card_type ILIKE ? OR cards.card_provider ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if err := query.Order("cards.card_id DESC").Limit(pageSize).Offset(offset).Scan(&cards).Error; err != nil {
		return nil, card_errors.ErrFindTrashedCardsFailed.WithInternal(err)
	}

	return cards, nil
}

func (r *cardQueryRepository) FindById(ctx context.Context, card_id int) (*models.CardAllFieldsRow, error) {
	if card_id <= 0 {
		return nil, sharedErrors.NewBadRequestError("card ID must be greater than zero")
	}

	var card models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Table("cards").Where("card_id = ? AND deleted_at IS NULL", card_id).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrCardNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &card, nil
}

func (r *cardQueryRepository) FindCardByUserId(ctx context.Context, user_id int) (*models.CardAllFieldsRow, error) {
	var card models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Table("cards").Where("user_id = ? AND deleted_at IS NULL", user_id).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrCardNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &card, nil
}

func (r *cardQueryRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	var card models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Table("cards").Where("card_number = ? AND deleted_at IS NULL", card_number).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrCardNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &card, nil
}

func (r *cardQueryRepository) FindUserCardByCardNumber(ctx context.Context, card_number string) (*models.CardByEmailRow, error) {
	var card models.CardByEmailRow
	err := r.db.WithContext(ctx).Table("cards").
		Select("cards.*, users.email").
		Joins("JOIN users ON users.user_id = cards.user_id").
		Where("cards.card_number = ? AND cards.deleted_at IS NULL", card_number).
		First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrCardNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &card, nil
}

func (r *cardQueryRepository) FindByIdUpdatedAt(ctx context.Context, card_id int, since time.Time) (*models.CardAllFieldsRow, error) {
	var card models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Table("cards").Where("card_id = ? AND deleted_at IS NULL AND updated_at >= ?", card_id, since).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrCardNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &card, nil
}

func (r *cardQueryRepository) FindCardByCardNumberForUpdate(ctx context.Context, card_number string) (*models.CardAllFieldsRow, error) {
	var card models.CardAllFieldsRow
	err := r.db.WithContext(ctx).Table("cards").Clauses().Where("card_number = ? AND deleted_at IS NULL", card_number).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, card_errors.ErrCardNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &card, nil
}
