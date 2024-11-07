package data

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"service-checker/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db: db,
	}, cleanup, nil
}
func NewDB(c *conf.Data) (*gorm.DB, func(), error) {
	db, cleanup, err := newDB(c.Database)
	if err != nil {
		return nil, nil, err
	}

	// Automatic migration
	if err = db.AutoMigrate(); err != nil {
		return nil, nil, err
	}

	return db, cleanup, err
}

func newDB(database *conf.Data_Database) (*gorm.DB, func(), error) {
	db, err := gorm.Open(mysql.Open(database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	// Setting the log Level
	db.Logger = db.Logger.LogMode(logger.LogLevel(database.LogLevel))

	// Set database configuration
	var sqlDB *sql.DB
	if sqlDB, err = db.DB(); err != nil {
		return nil, nil, err
	}

	if database.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(int(database.MaxIdleConns))
	} else {
		sqlDB.SetMaxIdleConns(10)
	}

	if database.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(int(database.MaxOpenConns))
	} else {
		sqlDB.SetMaxOpenConns(100)
	}

	if database.ConnMaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(database.ConnMaxLifetime.AsDuration())
	} else {
		sqlDB.SetConnMaxLifetime(time.Second * 300)
	}

	cleanup := func() {
		_ = sqlDB.Close()
	}

	return db, cleanup, err
}
