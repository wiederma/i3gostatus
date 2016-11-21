// Package backlight outputs the current backlight brightness in percent by
// reading the relevant files in `/sys/class/backlight/`. Since I only have an
// intel GPU, there is only Intel support right now.
package backlight

import (
	"fmt"
	"io/ioutil"
	"os/exec"
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
	c.BaseConfig.Parse(moduleName, configTree)
	c.BaseConfig.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	brightnessFile := "/sys/class/backlight/intel_backlight/brightness"
	maxBrightnessFile := "/sys/class/backlight/intel_backlight/max_brightness"
	incBrightnessCmd := []string{"xbacklight", "-inc", "5"}
	decBrightnessCmd := []string{"xbacklight", "-dec", "5"}
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

		go func() {
			// TODO: Update brightness after the click event has been processes.
			for event := range args.InCh {
				switch event.Button {
				case model.MouseButtonLeft, model.MouseWheelUp:
					exec.Command(incBrightnessCmd[0], incBrightnessCmd[1:]...).Run()
				case model.MouseButtonRight, model.MouseWheelDown:
					exec.Command(decBrightnessCmd[0], decBrightnessCmd[1:]...).Run()
				}
			}
		}()

		outputBlock.FullText = fmt.Sprintf(c.Format, output)
		args.OutCh <- outputBlock
	}
}
