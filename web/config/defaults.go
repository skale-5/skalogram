package config

import (
	"fmt"
	"sort"
)

func defaults() map[string]string {
	return map[string]string{
		"PG_HOST":     "127.0.0.1",
		"PG_PORT":     "5432",
		"PG_USER":     "postgres",
		"PG_PASSWORD": "postgres",
		"PG_DBNAME":   "skalogram",

		"REDIS_HOST": "127.0.0.1",
		"REDIS_PORT": "6379",
		"CACHE_TTL":  "60s",

		"STORAGE_TYPE":          "gs",
		"STORAGE_BUCKET":        "skalogram-posts-dev",
		"STORAGE_BUCKET_REGION": "eu-west3",

		"LISTEN_ADDR": "0.0.0.0",
		"LISTEN_PORT": "8080",
	}
}

func PrintDefaults() {
	fmt.Println("Default configurations:")
	keys := make([]string, 0, len(defaults()))
	for k := range defaults() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("\t%s=\"%s\"\n", k, defaults()[k])
	}
}
