// Package database provides database connection management.
package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/RAiWorks/RapidGo/v2/core/config"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DBConfig holds database connection configuration.
type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// NewDBConfig reads database configuration from environment variables.
func NewDBConfig() DBConfig {
	return DBConfig{
		Driver:          config.Env("DB_DRIVER", ""),
		Host:            config.Env("DB_HOST", "localhost"),
		Port:            config.Env("DB_PORT", "5432"),
		Name:            config.Env("DB_NAME", "rapidgo_dev"),
		User:            config.Env("DB_USER", ""),
		Password:        config.Env("DB_PASSWORD", ""),
		SSLMode:         config.Env("DB_SSL_MODE", "disable"),
		MaxOpenConns:    config.EnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    config.EnvInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: time.Duration(config.EnvInt("DB_CONN_MAX_LIFETIME", 5)) * time.Minute,
		ConnMaxIdleTime: time.Duration(config.EnvInt("DB_CONN_MAX_IDLE_TIME", 3)) * time.Minute,
	}
}

// DSN returns the data source name for the configured driver.
func (cfg DBConfig) DSN() string {
	switch cfg.Driver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)
	case "sqlite":
		return cfg.Name
	default:
		return ""
	}
}

// Connect establishes a database connection using environment configuration.
func Connect() (*gorm.DB, error) {
	return ConnectWithConfig(NewDBConfig())
}

// ConnectWithConfig establishes a database connection using the provided configuration.
func ConnectWithConfig(cfg DBConfig) (*gorm.DB, error) {
	dialector, err := newDialector(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newGormLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// NewReadDBConfig reads read-replica configuration from DB_READ_* environment
// variables. Each setting falls back to the corresponding DB_* value, then to
// the same defaults used by NewDBConfig.
func NewReadDBConfig() DBConfig {
	return DBConfig{
		Driver:          config.Env("DB_READ_DRIVER", config.Env("DB_DRIVER", "")),
		Host:            config.Env("DB_READ_HOST", config.Env("DB_HOST", "localhost")),
		Port:            config.Env("DB_READ_PORT", config.Env("DB_PORT", "5432")),
		Name:            config.Env("DB_READ_NAME", config.Env("DB_NAME", "rapidgo_dev")),
		User:            config.Env("DB_READ_USER", config.Env("DB_USER", "")),
		Password:        config.Env("DB_READ_PASSWORD", config.Env("DB_PASSWORD", "")),
		SSLMode:         config.Env("DB_READ_SSL_MODE", config.Env("DB_SSL_MODE", "disable")),
		MaxOpenConns:    config.EnvInt("DB_READ_MAX_OPEN_CONNS", config.EnvInt("DB_MAX_OPEN_CONNS", 25)),
		MaxIdleConns:    config.EnvInt("DB_READ_MAX_IDLE_CONNS", config.EnvInt("DB_MAX_IDLE_CONNS", 10)),
		ConnMaxLifetime: time.Duration(config.EnvInt("DB_READ_CONN_MAX_LIFETIME", config.EnvInt("DB_CONN_MAX_LIFETIME", 5))) * time.Minute,
		ConnMaxIdleTime: time.Duration(config.EnvInt("DB_READ_CONN_MAX_IDLE_TIME", config.EnvInt("DB_CONN_MAX_IDLE_TIME", 3))) * time.Minute,
	}
}

// newGormLogger returns a GORM logger configured based on environment.
// In development (APP_ENV=development) or when DB_LOG=true, it logs all
// queries with execution times. In production it stays silent.
func newGormLogger() gormlogger.Interface {
	env := strings.ToLower(config.Env("APP_ENV", "production"))
	dbLog := strings.ToLower(config.Env("DB_LOG", "false"))

	if env != "development" && dbLog != "true" {
		return gormlogger.Default.LogMode(gormlogger.Silent)
	}

	threshold := time.Duration(config.EnvInt("DB_SLOW_THRESHOLD_MS", 200)) * time.Millisecond

	return gormlogger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             threshold,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

// newDialector creates the appropriate GORM dialector for the configured driver.
func newDialector(cfg DBConfig) (gorm.Dialector, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Open(cfg.DSN()), nil
	case "mysql":
		return mysql.Open(cfg.DSN()), nil
	case "sqlite":
		return sqlite.Open(cfg.DSN()), nil
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER: %s", cfg.Driver)
	}
}
