package config

import (
	"encoding/json"
	"fmt"
	"github.com/jorgejr568/cloudflare-cli/internal/utils"
	"os"
)

type LocalConfig struct {
	CloudflareAPIKeyEntry string `json:"cloudflare_api_key"`
}

func (c LocalConfig) CloudflareAPIKey() string {
	return c.CloudflareAPIKeyEntry
}

func acquireLocalConfigPath() string {
	return fmt.Sprintf("%s/.config/cloudflare-cli", os.Getenv("HOME"))
}

func createPathIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	return nil
}

func SaveLocalConfig(config LocalConfig) error {
	path := acquireLocalConfigPath()
	if err := createPathIfNotExists(path); err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/config.json", path))
	if err != nil {
		return err
	}
	defer utils.LogErrorIfError(file.Close())

	fileData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if _, err := file.Write(fileData); err != nil {
		return err
	}

	return nil
}

func LoadLocalConfig() (LocalConfig, error) {
	path := acquireLocalConfigPath()
	file, err := os.Open(fmt.Sprintf("%s/config.json", path))
	if err != nil {
		if os.IsNotExist(err) {
			return LocalConfig{}, nil
		}
		return LocalConfig{}, err
	}
	defer utils.LogErrorIfError(file.Close())

	var config LocalConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return LocalConfig{}, err
	}

	return config, nil
}
