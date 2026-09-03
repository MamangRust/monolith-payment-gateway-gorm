package card_test

import (
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	sharederrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"

	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type CardRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	gormDB   *gorm.DB
	pool     *pgxpool.Pool
	repo     *repository.Repositories
	userRepo user_repo.Repositories
	userID   int
}

func (s *CardRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.pool = pool

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.gormDB = gormDB
	s.repo = repository.NewRepositories(gormDB)
	s.userRepo = user_repo.NewRepositories(gormDB)

	// Create a user for card ownership
	user, err := s.userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Card",
		LastName:  "Owner",
		Email:     fmt.Sprintf("card.owner-%d-%d@example.com", time.Now().UnixNano(), time.Now().UnixNano()%10000),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *CardRepositoryTestSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.ts.Teardown()
}

func (s *CardRepositoryTestSuite) createSeedCard() (*models.CardCreateRow, error) {
	return s.repo.CardCommand.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	})
}

func (s *CardRepositoryTestSuite) TestCreateCard() {
	ctx := context.Background()
	req := &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	}

	res, err := s.repo.CardCommand.CreateCard(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *CardRepositoryTestSuite) TestFindAllCards() {
	_, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.CardQuery.FindAllCards(ctx, &requests.FindAllCards{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *CardRepositoryTestSuite) TestFindById() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	found, err := s.repo.CardQuery.FindById(ctx, int(card.CardID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(card.CardID, found.CardID)
}

func (s *CardRepositoryTestSuite) TestFindByActive() {
	_, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.CardQuery.FindByActive(ctx, &requests.FindAllCards{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *CardRepositoryTestSuite) TestFindByTrashed() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.CardCommand.TrashedCard(ctx, int(card.CardID))
	s.Require().NoError(err)

	res, err := s.repo.CardQuery.FindByTrashed(ctx, &requests.FindAllCards{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *CardRepositoryTestSuite) TestUpdateCard() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	req := &requests.UpdateCardRequest{
		CardID:       int(card.CardID),
		UserID:       s.userID,
		CardType:     "credit",
		ExpireDate:   time.Now().AddDate(6, 0, 0),
		CVV:          "456",
		CardProvider: "MasterCard",
	}

	res, err := s.repo.CardCommand.UpdateCard(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *CardRepositoryTestSuite) TestTrashCard() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	trashed, err := s.repo.CardCommand.TrashedCard(ctx, int(card.CardID))
	s.NoError(err)
	s.NotNil(trashed)
}

func (s *CardRepositoryTestSuite) TestRestoreCard() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.CardCommand.TrashedCard(ctx, int(card.CardID))
	s.Require().NoError(err)

	restored, err := s.repo.CardCommand.RestoreCard(ctx, int(card.CardID))
	s.NoError(err)
	s.NotNil(restored)
}

func (s *CardRepositoryTestSuite) TestDeleteCardPermanent() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.CardCommand.TrashedCard(ctx, int(card.CardID))
	s.Require().NoError(err)

	success, err := s.repo.CardCommand.DeleteCardPermanent(ctx, int(card.CardID))
	s.NoError(err)
	s.True(success)
}

func (s *CardRepositoryTestSuite) TestRestoreAllCard() {
	card, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.CardCommand.TrashedCard(ctx, int(card.CardID))
	s.Require().NoError(err)

	success, err := s.repo.CardCommand.RestoreAllCard(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *CardRepositoryTestSuite) TestDeleteAllCardPermanent() {
	_, err := s.createSeedCard()
	s.Require().NoError(err)
	ctx := context.Background()

	success, err := s.repo.CardCommand.DeleteAllCardPermanent(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *CardRepositoryTestSuite) TestDuplicateCardNumberConflict() {
	ctx := context.Background()

	// Insert a card directly via GORM to control the card_number
	seed := models.Card{
		UserID:       int32(s.userID),
		CardNumber:   "9999999999999999",
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		Cvv:          "123",
		CardProvider: "Visa",
	}
	err := s.gormDB.WithContext(ctx).Create(&seed).Error
	s.Require().NoError(err)

	// Second insert with same card_number must violate unique constraint
	dup := models.Card{
		UserID:       int32(s.userID),
		CardNumber:   "9999999999999999",
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		Cvv:          "123",
		CardProvider: "Visa",
	}
	err = s.gormDB.WithContext(ctx).Create(&dup).Error
	s.Require().Error(err, "duplicate card_number must violate the unique constraint")

	appErr := sharederrors.ErrConstraintOrFailed(err, "Card", "create card")
	s.Require().NotNil(appErr)
	s.Equal(http.StatusConflict, appErr.Code, "23505 must map to 409")
	s.Contains(appErr.Message, "already exists")
}

func (s *CardRepositoryTestSuite) TestConcurrentCardUpdatesNoCorruption() {
	ctx := context.Background()
	card, err := s.createSeedCard()
	s.Require().NoError(err)

	base := requests.UpdateCardRequest{
		CardID:       int(card.CardID),
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	}
	g1 := base
	g1.CardType = "credit"
	g2 := base
	g2.CVV = "999"

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, e := s.repo.CardCommand.UpdateCard(ctx, &g1)
		errs <- e
	}()
	go func() {
		defer wg.Done()
		_, e := s.repo.CardCommand.UpdateCard(ctx, &g2)
		errs <- e
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		s.NoError(err)
	}

	found, err := s.repo.CardQuery.FindById(ctx, int(card.CardID))
	s.Require().NoError(err)
	valid := (found.CardType == "credit" && found.Cvv == "123") ||
		(found.CardType == "debit" && found.Cvv == "999")
	s.True(valid, "torn card state after concurrent updates: %s/%s", found.CardType, found.Cvv)
}

func TestCardRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CardRepositoryTestSuite))
}
