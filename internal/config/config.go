package config

import (
	"github.com/google/uuid"
	"os"
)

// config package is used to load the configuration from the environment variables
// the logic behind having separate config package is usability
// we can use the same package in other services as well to load the configuration

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Staging     Environment = "staging"
	Testing     Environment = "testing"
)

type ObjectStoreConfig struct {
	StoreType string
	Bucket    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

type DbConfig struct {
	Type             string `json:"db_type"`
	ConnectionString string `json:"connection_string"`
}

type SupabaseConfig struct {
	ProjectRef string `json:"project_ref"`
	ApiKey     string `json:"api_key"`
	JwtSecret  string `json:"jwt"`
}

type Config struct {
	Environment       string `json:"environment"`
	DbConfig          DbConfig
	ObjectStoreConfig ObjectStoreConfig
	SupabaseConfig    SupabaseConfig
	AdminUserID       uuid.UUID
}

var AppConfig *Config

// IsDevelopment returns true if the environment is development
func IsDevelopment() bool {
	return AppConfig.Environment == string(Development)
}

// IsProduction returns true if the environment is production
func IsProduction() bool {
	return AppConfig.Environment == string(Production)
}

// IsStaging returns true if the environment is staging
func IsStaging() bool {
	return AppConfig.Environment == string(Staging)
}

// IsTesting returns true if the environment is testing
func IsTesting() bool {
	return AppConfig.Environment == string(Testing)
}

// LoadConfig loads the configuration from the environment variables
func LoadConfig() *Config {
	Env := os.Getenv("ENVIRONMENT")
	if Env == "" {
		Env = string(Development)
	}

	// load db config
	DbType := os.Getenv("DB_TYPE")
	if DbType == "" {
		panic("DB_TYPE is not set")
	}

	DbConnString := os.Getenv("DB_CONNECTION_STRING")
	if DbType != "sqlite" && DbConnString == "" {
		panic("DB_CONNECTION_STRING is not set")
	}

	// load object supabase config
	SupabaseProjectRef := os.Getenv("SUPABASE_PROJECT_REF")
	if SupabaseProjectRef == "" {
		panic("SUPABASE_PROJECT_REF is not set")
	}

	SupabaseApiKey := os.Getenv("SUPABASE_API_KEY")
	if SupabaseApiKey == "" {
		panic("SUPABASE_API_KEY is not set")
	}

	SupabaseJwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if SupabaseJwtSecret == "" {
		panic("SUPABASE_JWT_SECRET is not set")
	}

	AdminUserID, err := uuid.Parse(os.Getenv("ADMIN_USER_ID"))
	if err != nil {
		panic(err)
	}

	AppConfig = &Config{
		Environment: Env,
		DbConfig: DbConfig{
			Type:             DbType,
			ConnectionString: DbConnString,
		},
		SupabaseConfig: SupabaseConfig{
			ProjectRef: SupabaseProjectRef,
			ApiKey:     SupabaseApiKey,
			JwtSecret:  SupabaseJwtSecret,
		},
		AdminUserID: AdminUserID,
	}

	return AppConfig
}
