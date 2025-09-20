package config

import (
	"os"
	"strings"
)

// Config centralizes the application's configurations.
type Config struct {
	ProviderType string // "file"
	FilePath     string // path to packs file (when ProviderType="file")
	EnvVar       string
	HTTPAddr     string
}

// Load reads the environment variables and builds the Config.
// Can be expanded to include flags or YAML files in the future.
func Load() Config {
	return Config{
		ProviderType: getEnv("PACK_PROVIDER", "file"),
		FilePath:     getEnv("PACK_SIZES_FILE", "./packs.csv"),
		EnvVar:       getEnv("PACK_SIZES_ENV", "PACK_SIZES"),
		HTTPAddr:     getEnv("HTTP_ADDR", ":8080"),
	}
}

// -------- helpers --------

func getEnv(key, def string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return def
}
