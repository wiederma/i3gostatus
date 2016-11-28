package config

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

func Path() string {
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		return xdgHome + "/i3gostatus/config.toml"
	}
	return os.Getenv("HOME") + "/.config/i3gostatus/config.toml"
}

func Load(path string) *toml.TomlTree {
	var configTree *toml.TomlTree
	var err error

	configTree, err = toml.LoadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Using default config...")

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
