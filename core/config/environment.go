package config

// AppEnv returns the current application environment (development, production, testing).
func AppEnv() string {
	return Env("APP_ENV", "development")
}

// IsProduction returns true if APP_ENV is "production".
func IsProduction() bool {
	return AppEnv() == "production"
}

// IsDevelopment returns true if APP_ENV is "development".
func IsDevelopment() bool {
	return AppEnv() == "development"
}

// IsTesting returns true if APP_ENV is "testing".
func IsTesting() bool {
	return AppEnv() == "testing"
}

// IsLocal returns true if APP_ENV is "local" or "development".
func IsLocal() bool {
	return AppEnv() == "local" || AppEnv() == "development"
}

// IsDebug returns true if APP_DEBUG is "true" or "1".
func IsDebug() bool {
	return EnvBool("APP_DEBUG", false)
}
