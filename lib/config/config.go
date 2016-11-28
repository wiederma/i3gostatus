package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "config ", log.LstdFlags)
}

func Path() string {
	var configPath string

	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		configPath = xdgHome + "/i3gostatus/config.toml"
	} else {
		configPath = os.Getenv("HOME") + "/.config/i3gostatus/config.toml"
	}

	logger.Printf("Using config file: %s", configPath)

	return configPath
}

func Load(path string) *toml.TomlTree {
	var configTree *toml.TomlTree
	var err error

	configTree, err = toml.LoadFile(path)
	if err != nil {
		logger.Println(err)
		logger.Println("Using default config...")

		defaultConfig := `modules = ["datetime"]`
		configTree, err = toml.Load(defaultConfig)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return configTree
}

func GetString(configTree *toml.TomlTree, path string, def string) string {
	return configTree.GetDefault(path, def).(string)
}

func GetInt(configTree *toml.TomlTree, path string, def int) int {
	return configTree.GetDefault(path, def).(int)
}

func GetFloat64(configTree *toml.TomlTree, path string, def float64) float64 {
	return configTree.GetDefault(path, def).(float64)
}
func GetDuration(configTree *toml.TomlTree, path string, def int) time.Duration {
	return time.Duration(configTree.GetDefault(path, int64(def)).(int64))
}

func GetDurationMs(configTree *toml.TomlTree, path string, def int) time.Duration {
	return time.Duration(configTree.GetDefault(path, int64(def)).(int64)) * time.Millisecond
}
