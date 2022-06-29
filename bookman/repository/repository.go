package repository

import (
	"bookman/config"
	"bookman/entity"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initClient(cfg *config.RDB) error {
	if db != nil {
		return nil
	}

	const dsnTemplate = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s"
	dsn := fmt.Sprintf(dsnTemplate, cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.TimeZone)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func paginate(value any, pagination *entity.Pagination) func(db *gorm.DB) *gorm.DB {
	var total int64
	db.Model(value).Count(&total)
	pagination.Total = int(total)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.Offset).Limit(pagination.Limit)
	}
}
