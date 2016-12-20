// Package temperature shows the current CPU temperature which is read out from
// the relevant files in `/sys/class/thermal/`.
package temperature

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
	name               = "temperature"
	moduleName         = "i3gostatus.modules." + name
	defaultPeriod      = 5000
	defaultFormat      = "%d°C"
	defaultUrgentTemp  = 70
	defaultUrgentColor = "#FF0000"
)

type Config struct {
	model.BaseConfig
	UrgentTemp  int
	UrgentColor string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.BaseConfig.Period = config.GetDurationMs(configTree, c.Name+".period", defaultPeriod)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.UrgentTemp = config.GetInt(configTree, name+".urgent_temp", defaultUrgentTemp)
	c.UrgentColor = config.GetString(configTree, name+".urgent_color", defaultUrgentColor)
}

func (c *Config) Run(args *model.ModuleArgs) {
	var (
		outputBlock *model.I3BarBlock
		temperature int
		thermalFile = "/sys/class/thermal/thermal_zone0/temp"
	)

	for range time.NewTicker(c.Period).C {
		outputBlock = model.NewBlock(moduleName, c.BaseConfig, args.Index)
		data, err := ioutil.ReadFile(thermalFile)
		if err != nil {
			panic(err)
		}

		dataStr := strings.TrimSpace(string(data))

		if t, err := strconv.Atoi(dataStr); err == nil {
			temperature = int(t / 1000)
			if temperature >= c.UrgentTemp {
				outputBlock.Urgent = true
				outputBlock.Color = c.UrgentColor
			}
		}

		outputBlock.FullText = fmt.Sprintf(c.Format, temperature)
		args.OutCh <- outputBlock
	}
}
