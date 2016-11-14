package temperature

import (
	"time"
	"io/ioutil"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name       = "temperature"
	moduleName = "i3gostatus.modules." + name
)

type Config struct {
	model.BaseConfig
	Output string
}

func (config *Config) ReadConfig(configTree *toml.TomlTree) {
	config.BaseConfig.ReadConfig(name, configTree)
}

func (config *Config) Run(out chan *model.I3BarBlockWrapper, index int) {
	ticker := time.NewTicker(config.Period)
	outputBlock := model.NewBlock(moduleName, config.BaseConfig, index)
	var temperatureStr string

	for range ticker.C {
		data, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
		if err != nil {
			panic(err)
		}

		temperatureStr = strings.TrimSuffix(strings.TrimSpace(string(data)), "000")

		outputBlock.FullText = temperatureStr
		out <- outputBlock
	}
}
