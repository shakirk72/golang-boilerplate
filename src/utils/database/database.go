package database

import (
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/models"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	cfg := config.App.Database
	logLevel := gormlogger.Info
	if config.App.App.IsProd {
		logLevel = gormlogger.Error
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get sql.DB from gorm", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	logger.Info("database connected")
}

func Migrate() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		logger.Fatal("auto-migration failed", zap.Error(err))
	}
	logger.Info("database migration complete")
}

// Ping checks that the database connection is still alive.
func Ping() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
