package database

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewGormClient creates a new *gorm.DB connection using the same viper DB_*
// configuration that the existing pgxpool client uses. The returned *gorm.DB
// is intended to live alongside the legacy pgxpool.Pool during the strangler
// migration and will be removed once all services have moved to GORM.
func NewGormClient(logger logger.LoggerInterface) (*gorm.DB, error) {
	dbDriver := viper.GetString("DB_DRIVER")
	if dbDriver != "postgres" && dbDriver != "pgx" && dbDriver != "" {
		logger.Error("GORM postgres driver only supports PostgreSQL", zap.String("DB_DRIVER", dbDriver))
		return nil, fmt.Errorf("gorm postgres driver only supports PostgreSQL, got: %s", dbDriver)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Jakarta",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USERNAME"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
	)

	maxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}

	maxIdleConns := viper.GetInt("DB_MIN_IDLE_CONNS")
	if maxIdleConns <= 0 {
		maxIdleConns = 50
	}

	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
	}

	connMaxIdleTime := viper.GetDuration("DB_CONN_MAX_IDLE_TIME")
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 30 * time.Minute
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Error("Failed to open GORM connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB from GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to get underlying sql.DB from GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database via GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database via GORM: %w", err)
	}

	logger.Debug("GORM database connection established successfully",
		zap.Int("MaxOpenConns", maxOpenConns),
		zap.Int("MaxIdleConns", maxIdleConns),
		zap.Duration("ConnMaxLifetime", connMaxLifetime),
		zap.Duration("ConnMaxIdleTime", connMaxIdleTime),
	)

	return gormDB, nil
}
