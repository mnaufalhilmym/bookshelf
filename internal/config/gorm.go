package config

import (
	"fmt"
	"time"

	"github.com/mnaufalhilmym/gotracing"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(
	database string,
	poolIdle int,
	poolMax int,
	poolLifetime int,
) *gorm.DB {
	db, err := gorm.Open(
		sqlite.Open(database),
		&gorm.Config{
			TranslateError: true,
			Logger: logger.New(&gormTracingWriter{}, logger.Config{
				SlowThreshold:        time.Second,
				ParameterizedQueries: true,
				LogLevel:             logger.Info,
			}),
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	connection, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get *sql.DB: %w", err))
	}

	connection.SetMaxIdleConns(poolIdle)
	connection.SetMaxOpenConns(poolMax)
	connection.SetConnMaxLifetime(time.Duration(poolLifetime * int(time.Second)))

	if err := connection.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping the database: %w", err))
	}

	return db
}

type gormTracingWriter struct{}

func (*gormTracingWriter) Printf(format string, args ...any) {
	gotracing.Tracef(format, args...)
}
