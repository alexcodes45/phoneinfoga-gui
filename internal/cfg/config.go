package cfg

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	PreferServe bool   `mapstructure:"prefer_serve"`
	BinaryPath  string `mapstructure:"phoneinfoga_binary"`
	ProxyURL    string `mapstructure:"proxy_url"`
	TimeoutMs   int    `mapstructure:"timeout_ms"`
}

func AppConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil { return "", err }
	dir := filepath.Join(base, "phoneinfoga-desktop")
	if err := os.MkdirAll(dir, 0o755); err != nil { return "", err }
	return dir, nil
}

func Load() (*Config, error) {
	dir, err := AppConfigDir()
	if err != nil { return nil, err }
	v := viper.New()
	v.SetConfigFile(filepath.Join(dir, "config.yaml"))
	v.SetDefault("prefer_serve", true)
	v.SetDefault("timeout_ms", 60000)

	if err := v.ReadInConfig(); err != nil {
		// first run: ignore
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil { return nil, err }
	return &c, nil
}

func Save(c *Config) error {
	dir, err := AppConfigDir()
	if err != nil { return err }
	v := viper.New()
	v.Set("prefer_serve", c.PreferServe)
	v.Set("phoneinfoga_binary", c.BinaryPath)
	v.Set("proxy_url", c.ProxyURL)
	v.Set("timeout_ms", c.TimeoutMs)
	return v.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}
