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
	"github.com/rumpelsepp/i3gostatus/lib/utils"
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
	c.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
}

func getBrightness() float64 {
	brightnessFile := "/sys/class/backlight/intel_backlight/brightness"
	maxBrightnessFile := "/sys/class/backlight/intel_backlight/max_brightness"
	var res float64

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
		res = float64(val)
	} else {
		panic(err)
	}

	if val, err := strconv.Atoi(maxBrightnessStr); err == nil {
		res = (res / float64(val)) * 100
	} else {
		panic(err)
	}

	return res
}

func clickHandler(args *model.ModuleArgs) {
	// TODO: Update brightness after the click event has been processes.
	cmd, err := utils.Which("xbacklight")
	if err != nil {
		switch err.(type) {
		case utils.CommandNotAvailError:
			// TODO: Log a warning here (once the logging system is there...)
			return
		default:
			panic(err)
		}
	}

	incBrightnessCmd := []string{cmd, "-inc", "5"}
	decBrightnessCmd := []string{cmd, "-dec", "5"}

	for event := range args.InCh {
		switch event.Button {
		case model.MouseButtonLeft, model.MouseWheelUp:
			exec.Command(incBrightnessCmd[0], incBrightnessCmd[1:]...).Run()
		case model.MouseButtonRight, model.MouseWheelDown:
			exec.Command(decBrightnessCmd[0], decBrightnessCmd[1:]...).Run()
		default:
			continue
		}
	}
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	var output float64

	go clickHandler(args)

	for range time.NewTicker(c.Period).C {
		output = getBrightness()
		outputBlock.FullText = fmt.Sprintf(c.Format, output)
		args.OutCh <- outputBlock
	}
}
