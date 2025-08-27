package initialization

import (
	"database/sql"
	"fmt"

	"github.com/SomeHowMicroservice/shm-be/user/config"
	"github.com/SomeHowMicroservice/shm-be/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var allModels = []interface{}{
	&model.User{},
	&model.Role{},
	&model.Profile{},
	&model.Measurement{},
	&model.Address{},
}

type DB struct {
	Gorm *gorm.DB
	sql *sql.DB
}

func InitDB(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s channel_binding=%s",
		cfg.Database.DBHost,
		cfg.Database.DBName,
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBSSLMode,
		cfg.Database.DBChannelBinding,
	)
	gDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("kết nối PostgreSQL thất bại: %w", err)
	}

	if err := runAutoMigrations(gDB); err != nil {
		return nil, fmt.Errorf("chuyển dịch DB thất bại: %w", err)
	}

	sqlDB, err := gDB.DB()
	if err != nil {
		return nil, fmt.Errorf("không lấy được sql.DB: %w", err)
	}

	return &DB{
		gDB,
		sqlDB,
	}, nil
}

func (d *DB) Close() {
	_ = d.sql.Close()
}

func runAutoMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(allModels...); err != nil {
		return err
	}

	return nil
}
