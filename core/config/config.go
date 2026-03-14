package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

// Load reads the .env file and sets environment variables.
// If no .env file is found, it logs a message and continues
// (system environment variables are still available).
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}
}

// LoadConfig reads environment variables into a struct T using struct tags.
//
// Supported tags:
//   - `env:"KEY"` — env var name to read (required for the field to be populated)
//   - `default:"value"` — fallback if env var is empty/unset
//   - `validate:"required"` — passed to go-playground/validator
//
// Supported field types: string, int, int64, bool, float64, time.Duration, []string (comma-separated).
//
// Example:
//
//	type AppConfig struct {
//	    Port    int           `env:"PORT" default:"8080" validate:"min=1,max=65535"`
//	    Debug   bool          `env:"APP_DEBUG" default:"false"`
//	    Name    string        `env:"APP_NAME" default:"app"`
//	    Timeout time.Duration `env:"TIMEOUT" default:"30s"`
//	}
//	cfg, err := config.LoadConfig[AppConfig]()
func LoadConfig[T any]() (T, error) {
	var cfg T
	v := reflect.ValueOf(&cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := v.Field(i)

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		raw := os.Getenv(envKey)
		if raw == "" {
			raw = field.Tag.Get("default")
		}
		if raw == "" {
			continue
		}

		if err := setField(fv, raw); err != nil {
			var zero T
			return zero, fmt.Errorf("config: field %s: %w", field.Name, err)
		}
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		var zero T
		return zero, fmt.Errorf("config: validation failed: %w", err)
	}

	return cfg, nil
}

// MustLoadConfig is like LoadConfig but panics on error.
// Use in application startup where a missing/invalid config is fatal.
func MustLoadConfig[T any]() T {
	cfg, err := LoadConfig[T]()
	if err != nil {
		panic(fmt.Sprintf("config: %v", err))
	}
	return cfg
}

// setField sets a reflect.Value from a raw string, supporting common Go types.
func setField(fv reflect.Value, raw string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(raw)
	case reflect.Int, reflect.Int64:
		if fv.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(raw)
			if err != nil {
				return err
			}
			fv.Set(reflect.ValueOf(d))
		} else {
			n, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return err
			}
			fv.SetInt(n)
		}
	case reflect.Bool:
		fv.SetBool(raw == "true" || raw == "1")
	case reflect.Float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	case reflect.Slice:
		if fv.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(raw, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			fv.Set(reflect.ValueOf(parts))
		} else {
			return fmt.Errorf("unsupported slice type %s", fv.Type().Elem().Kind())
		}
	default:
		return fmt.Errorf("unsupported type %s", fv.Kind())
	}
	return nil
}
