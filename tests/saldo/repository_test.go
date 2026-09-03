package saldo_test

import (
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	tests "github.com/MamangRust/monolith-payment-gateway-test"
	user_repo "github.com/MamangRust/monolith-payment-gateway-user/repository"

	"github.com/stretchr/testify/suite"
)

type SaldoRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	repo     repository.Repositories
	queries  *gorm.DB
	cardRepo *card_repo.Repositories
	userRepo user_repo.Repositories
	userID   int
}

func (s *SaldoRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, gormErr := gorm.Open(postgres.Open(s.ts.DBURL), &gorm.Config{})
	if gormErr != nil {
		s.Require().NoError(gormErr)
	}
	s.queries = gormDB
	s.repo = repository.NewRepositories(gormDB)
	s.cardRepo = card_repo.NewRepositories(gormDB)
	s.userRepo = user_repo.NewRepositories(gormDB)

	// Create user
	user, err := s.userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Saldo",
		LastName:  "Owner",
		Email:     fmt.Sprintf("saldo.owner-%d@example.com", time.Now().UnixNano()),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *SaldoRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *SaldoRepositoryTestSuite) createSeedSaldo() (*models.CreateSaldoRow, error) {
	card, err := s.cardRepo.CardCommand.CreateCard(context.Background(), &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	})
	if err != nil {
		return nil, err
	}

	return s.repo.CreateSaldo(context.Background(), &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 100000,
	})
}

func (s *SaldoRepositoryTestSuite) TestCreateSaldo() {
	ctx := context.Background()

	card, _ := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "123",
		CardProvider: "Visa",
	})

	req := &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 100000,
	}

	res, err := s.repo.CreateSaldo(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *SaldoRepositoryTestSuite) TestCreateSaldoIsIdempotentAgainstDefaultEvent() {
	ctx := context.Background()
	card, err := s.cardRepo.CardCommand.CreateCard(ctx, &requests.CreateCardRequest{
		UserID:       s.userID,
		CardType:     "debit",
		ExpireDate:   time.Now().AddDate(5, 0, 0),
		CVV:          "789",
		CardProvider: "Visa",
	})
	s.Require().NoError(err)

	_, err = s.repo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 1000000,
	})
	s.Require().NoError(err)

	// The card-created Kafka event carries the default zero balance. It must
	// not overwrite a balance that was already initialized through the API.
	err = s.repo.CreateSaldoIfNotExists(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 0,
	})
	s.Require().NoError(err)

	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- s.repo.CreateSaldoIfNotExists(ctx, &requests.CreateSaldoRequest{
				CardNumber:   card.CardNumber,
				TotalBalance: 0,
			})
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		s.Require().NoError(err)
	}

	found, err := s.repo.FindByCardNumber(ctx, card.CardNumber)
	s.Require().NoError(err)
	s.Equal(int32(1000000), found.TotalBalance)

	// A repeated API command remains atomic and updates the existing active row
	// rather than creating a second saldo.
	_, err = s.repo.CreateSaldo(ctx, &requests.CreateSaldoRequest{
		CardNumber:   card.CardNumber,
		TotalBalance: 1250000,
	})
	s.Require().NoError(err)

	found, err = s.repo.FindByCardNumber(ctx, card.CardNumber)
	s.Require().NoError(err)
	s.Equal(int32(1250000), found.TotalBalance)
}

func (s *SaldoRepositoryTestSuite) TestFindAllSaldos() {
	_, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindAllSaldos(ctx, &requests.FindAllSaldos{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *SaldoRepositoryTestSuite) TestFindById() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	found, err := s.repo.FindById(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(saldo.SaldoID, found.SaldoID)
}

func (s *SaldoRepositoryTestSuite) TestFindByActive() {
	_, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindByActive(ctx, &requests.FindAllSaldos{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *SaldoRepositoryTestSuite) TestFindByTrashed() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	res, err := s.repo.FindByTrashed(ctx, &requests.FindAllSaldos{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *SaldoRepositoryTestSuite) TestUpdateSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	id := int(saldo.SaldoID)
	req := &requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   saldo.CardNumber,
		TotalBalance: 200000,
	}

	res, err := s.repo.UpdateSaldo(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *SaldoRepositoryTestSuite) TestTrashSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	trashed, err := s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.NotNil(trashed)
}

func (s *SaldoRepositoryTestSuite) TestRestoreSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	restored, err := s.repo.RestoreSaldo(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.NotNil(restored)
}

func (s *SaldoRepositoryTestSuite) TestDeleteSaldoPermanent() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	success, err := s.repo.DeleteSaldoPermanent(ctx, int(saldo.SaldoID))
	s.NoError(err)
	s.True(success)
}

func (s *SaldoRepositoryTestSuite) TestRestoreAllSaldo() {
	saldo, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedSaldo(ctx, int(saldo.SaldoID))
	s.Require().NoError(err)

	success, err := s.repo.RestoreAllSaldo(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *SaldoRepositoryTestSuite) TestDeleteAllSaldoPermanent() {
	_, err := s.createSeedSaldo()
	s.Require().NoError(err)
	ctx := context.Background()

	success, err := s.repo.DeleteAllSaldoPermanent(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *SaldoRepositoryTestSuite) TestConcurrentSaldoDeltasApplyAll() {
	ctx := context.Background()
	seed, err := s.createSeedSaldo()
	s.Require().NoError(err)

	// Two concurrent guarded delta updates must both apply — the atomic
	// UPDATE ... total_balance = total_balance + $2 serializes at row level,
	// so no delta is lost (mirror of Test7_TopupConcurrentUpdatesApplyAllDeltas).
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		e := s.queries.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = CURRENT_TIMESTAMP WHERE card_number = ? AND deleted_at IS NULL`, 500, seed.CardNumber).Error
		errs <- e
	}()
	go func() {
		defer wg.Done()
		e := s.queries.Exec(`UPDATE saldos SET total_balance = total_balance + ?, updated_at = CURRENT_TIMESTAMP WHERE card_number = ? AND deleted_at IS NULL`, 300, seed.CardNumber).Error
		errs <- e
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		s.NoError(err)
	}

	found, err := s.repo.FindByCardNumber(ctx, seed.CardNumber)
	s.Require().NoError(err)
	s.Equal(int32(100800), found.TotalBalance, "both concurrent deltas must be applied (no lost update)")
}

func TestSaldoRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SaldoRepositoryTestSuite))
}
