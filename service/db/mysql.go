package db

import (
	"context"
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// SQLCommon ...
type (
	DB = gorm.DB
	// Model ...
	Model = gorm.Model
	// Association ...
	Association = gorm.Association
)

var (
	// errSlowCommand ...
	errSlowCommand = errors.New("mysql slow command")

	// ErrRecordNotFound returns a "record not found error". Occurs only when attempting to query the database with a struct; querying with a slice won't return this error
	ErrRecordNotFound = gorm.ErrRecordNotFound
	// ErrInvalidTransaction occurs when you are trying to `Commit` or `Rollback`
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	// ErrMissingWhereClause missing where clause
	ErrMissingWhereClause = gorm.ErrMissingWhereClause
	// ErrUnsupportedRelation unsupported relations
	ErrUnsupportedRelation = gorm.ErrUnsupportedRelation
	// ErrInvalidData unsupported data
	ErrInvalidData = gorm.ErrInvalidData
)

// WithContext ...
func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	db.WithContext(ctx)
	return db
}

// Open ...
func Open(options *Config) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(options.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set the connection options
	if options.Debug {
		db = db.Debug()
	}
	sqlDB.SetMaxIdleConns(options.MaxIdleConns)
	sqlDB.SetMaxOpenConns(options.MaxOpenConns)

	if options.ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(options.ConnMaxLifetime)
	}

	return
}
