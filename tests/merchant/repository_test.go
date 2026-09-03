package merchant_test

import (
	models "github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/repository"
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

type MerchantRepositoryTestSuite struct {
	suite.Suite
	ts       *tests.TestSuite
	gormDB   *gorm.DB
	pool     *pgxpool.Pool
	repo     repository.Repositories
	userRepo user_repo.Repositories
	userID   int
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
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

	// Create a user for merchant ownership
	user, err := s.userRepo.UserCommand().CreateUser(context.Background(), &requests.CreateUserRequest{
		FirstName: "Merchant",
		LastName:  "Owner",
		Email:     fmt.Sprintf("merchant.owner-%d-%d@example.com", time.Now().UnixNano(), time.Now().UnixNano()%10000),
		Password:  "password123",
	})
	s.Require().NoError(err)
	s.userID = int(user.UserID)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.ts.Teardown()
}

func (s *MerchantRepositoryTestSuite) createSeedMerchant() (*models.MerchantAllFieldsRow, error) {
	return s.repo.CreateMerchant(context.Background(), &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   fmt.Sprintf("Test Merchant %d", time.Now().UnixNano()),
	})
}

func (s *MerchantRepositoryTestSuite) TestCreateMerchant() {
	ctx := context.Background()
	req := &requests.CreateMerchantRequest{
		UserID: s.userID,
		Name:   fmt.Sprintf("Test Merchant %d", time.Now().UnixNano()),
	}

	res, err := s.repo.CreateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *MerchantRepositoryTestSuite) TestFindAllMerchants() {
	_, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	res, err := s.repo.FindAllMerchants(ctx, &requests.FindAllMerchants{
		Page:     1,
		PageSize: 10,
		Search:   "",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(res), 1)
}

func (s *MerchantRepositoryTestSuite) TestFindByApiKey() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	found, err := s.repo.FindByApiKey(ctx, merchant.ApiKey)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(merchant.ApiKey, found.ApiKey)
}

func (s *MerchantRepositoryTestSuite) TestFindByMerchantId() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	found, err := s.repo.FindByMerchantId(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(found)
	s.Equal(merchant.MerchantID, found.MerchantID)
}

func (s *MerchantRepositoryTestSuite) TestUpdateMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	merchantID := int(merchant.MerchantID)
	req := &requests.UpdateMerchantRequest{
		MerchantID: &merchantID,
		Name:       "Updated Merchant",
		UserID:     s.userID,
		Status:     "active",
	}

	res, err := s.repo.UpdateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Merchant", res.Name)
}

func (s *MerchantRepositoryTestSuite) TestTrashMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	trashed, err := s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(trashed)
}

func (s *MerchantRepositoryTestSuite) TestRestoreMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	restored, err := s.repo.RestoreMerchant(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.NotNil(restored)
}

func (s *MerchantRepositoryTestSuite) TestDeleteMerchantPermanent() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	success, err := s.repo.DeleteMerchantPermanent(ctx, int(merchant.MerchantID))
	s.NoError(err)
	s.True(success)
}

func (s *MerchantRepositoryTestSuite) TestRestoreAllMerchant() {
	merchant, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	_, err = s.repo.TrashedMerchant(ctx, int(merchant.MerchantID))
	s.Require().NoError(err)

	success, err := s.repo.RestoreAllMerchant(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *MerchantRepositoryTestSuite) TestDeleteAllMerchantPermanent() {
	_, err := s.createSeedMerchant()
	s.Require().NoError(err)
	ctx := context.Background()

	success, err := s.repo.DeleteAllMerchantPermanent(ctx)
	s.NoError(err)
	s.True(success)
}

func (s *MerchantRepositoryTestSuite) TestDuplicateApiKeyConflict() {
	ctx := context.Background()

	// Insert a merchant directly via GORM to control the api_key
	seed := models.Merchant{
		UserID: int32(s.userID),
		ApiKey: "duplicate-api-key-12345",
		Name:   "First Merchant",
	}
	err := s.gormDB.WithContext(ctx).Create(&seed).Error
	s.Require().NoError(err)

	// Second insert with same api_key must violate unique constraint
	dup := models.Merchant{
		UserID: int32(s.userID),
		ApiKey: "duplicate-api-key-12345",
		Name:   "Second Merchant",
	}
	err = s.gormDB.WithContext(ctx).Create(&dup).Error
	s.Require().Error(err, "duplicate api_key must violate the unique constraint")

	appErr := sharederrors.ErrConstraintOrFailed(err, "Merchant", "create merchant")
	s.Require().NotNil(appErr)
	s.Equal(http.StatusConflict, appErr.Code, "23505 must map to 409")
	s.Contains(appErr.Message, "already exists")
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}
