package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	API      APIConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env         string
	Port        string
	Name        string
	Version     string
	IsProd      bool
	CORSOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type APIConfig struct {
	BaseURL string
	Key     string
	Timeout int
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

var App *Config

// envDefaults holds per-environment default values that can still be
// overridden by the corresponding .env keys at runtime.
type envDefaults struct {
	APITimeout int
	JWTExpiry  int
}

func getEnvDefaults(env string) envDefaults {
	switch env {
	case EnvProd:
		return prodDefaults
	case EnvQA:
		return qaDefaults
	default:
		return devDefaults
	}
}

// Load reads .env, detects APP_ENV (default: dev), and populates App config.
func Load(envFile string) {
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("warning: could not load %s, falling back to OS env: %v", envFile, err)
	}

	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "" {
		env = EnvDev
	}
	if env != EnvDev && env != EnvQA && env != EnvProd {
		log.Printf("warning: unknown APP_ENV=%q, falling back to dev", env)
		env = EnvDev
	}

	defaults := getEnvDefaults(env)
	timeout, _ := strconv.Atoi(getEnv("EXTERNAL_API_TIMEOUT", strconv.Itoa(defaults.APITimeout)))
	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", strconv.Itoa(defaults.JWTExpiry)))

	corsRaw := getEnv("CORS_ORIGINS", "*")
	var corsOrigins []string
	for _, o := range strings.Split(corsRaw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			corsOrigins = append(corsOrigins, o)
		}
	}

	App = &Config{
		App: AppConfig{
			Env:         env,
			Port:        getEnv("APP_PORT", Common.DefaultPort),
			Name:        Common.AppName,
			Version:     Common.AppVersion,
			IsProd:      env == EnvProd,
			CORSOrigins: corsOrigins,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", Common.DefaultDBPort),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "gomvc_dev"),
			SSLMode:  getEnv("DB_SSLMODE", Common.DefaultSSLMode),
			Timezone: getEnv("DB_TIMEZONE", Common.DefaultTZ),
		},
		API: APIConfig{
			BaseURL: getEnv("EXTERNAL_API_BASE_URL", ""),
			Key:     getEnv("EXTERNAL_API_KEY", ""),
			Timeout: timeout,
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", ""),
			ExpiryHours: expiryHours,
		},
	}

	log.Printf("config loaded [env=%s]", env)
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.Name, d.Port, d.SSLMode, d.Timezone,
	)
}

// Validate checks that all required config values are present.
func (c *Config) Validate() error {
	var missing []string
	if c.JWT.Secret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if c.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.Database.Password == "" && c.App.IsProd {
		missing = append(missing, "DB_PASSWORD (required in prod)")
	}
	if len(missing) > 0 {
		return errors.New("missing required config: " + strings.Join(missing, ", "))
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
