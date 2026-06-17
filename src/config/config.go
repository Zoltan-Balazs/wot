package config

import "os"

type Config struct {
	CloudantURL string
	DBName      string
	DesignDoc   string
	IndexName   string
}

var instance *Config

func Get() Config {
	if instance != nil {
		return *instance
	}
	instance = &Config{
		CloudantURL: withDefault("CLOUDANT_URL", "https://mikerhodes.cloudant.com"),
		DBName:      withDefault("CLOUDANT_DB", "airportdb"),
		DesignDoc:   withDefault("CLOUDANT_DESIGN_DOC", "view1"),
		IndexName:   withDefault("CLOUDANT_INDEX", "geo"),
	}
	return *instance
}

func withDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
