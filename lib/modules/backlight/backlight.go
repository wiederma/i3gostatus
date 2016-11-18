// Package backlight outputs the current backlight brightness in percent by
// reading the relevant files in `/sys/class/backlight/`. Since I only have an
// intel GPU, there is only Intel support right now.
package backlight

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name          = "backlight"
	moduleName    = "i3gostatus.modules." + name
	defaultPeriod = 3000
	defaultFormat = "%.0f %%"
)

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.BaseConfig.Period = config.GetDurationMs(configTree, c.Name+".period", defaultPeriod)
	c.Format = config.GetString(configTree, c.Name+".format", defaultFormat)
}

func (c *Config) Run(out chan *model.I3BarBlockWrapper, index int) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, index)
	brightnessFile := "/sys/class/backlight/intel_backlight/brightness"
	maxBrightnessFile := "/sys/class/backlight/intel_backlight/max_brightness"
	var output float64

	for range time.NewTicker(c.Period).C {
		brightness, err := ioutil.ReadFile(brightnessFile)
		if err != nil {
			panic(err)
		}

		maxBrightness, err := ioutil.ReadFile(maxBrightnessFile)
		if err != nil {
			panic(err)
		}

		brightnessStr := strings.TrimSpace(string(brightness))
		maxBrightnessStr := strings.TrimSpace(string(maxBrightness))

		if val, err := strconv.Atoi(brightnessStr); err == nil {
			output = float64(val)
		} else {
			panic(err)
		}

		if val, err := strconv.Atoi(maxBrightnessStr); err == nil {
			output = (output / float64(val)) * 100
		} else {
			panic(err)
		}

		outputBlock.FullText = fmt.Sprintf(c.Format, output)
		out <- outputBlock
	}
}
