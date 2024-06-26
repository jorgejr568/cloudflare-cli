package config

import "os"

type envConfig struct{}

func (c envConfig) CloudflareAPIKey() string {
	return os.Getenv("CLOUDFLARE_API_KEY")
}
