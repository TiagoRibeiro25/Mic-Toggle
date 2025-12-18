package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Hotkey           string `json:"hotkey"`
	PlayBeep         bool   `json:"playBeep"`
	ShowNotification bool   `json:"showNotification"`
}

func defaultConfig() *Config {
	return &Config{
		Hotkey: "Ctrl+Shift+M",
	}
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "mic-toggle")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, "config.json"), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := defaultConfig()
			Save(cfg)
			return cfg, nil
		}
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
