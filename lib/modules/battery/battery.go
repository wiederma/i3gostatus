// Package static adds a static string to i3bar. Its main purpose is
// demonstrating the module API of `i3gostatus` and it acts as a template for
// new modules.
package battery

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name               = "battery"
	moduleName         = "i3gostatus.modules." + name
	defaultPeriod      = 5000
	powerSupplyBaseDir = "/sys/class/power_supply/"
	defaultFormat      = `BAT: {{.EnergyNowPerc | printf "%.0f" }}%`
	defaultMinWidth    = `BAT: 90%`
)

var logger = log.New(os.Stderr, "i3gostatus ", log.LstdFlags)

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
	c.MinWidth = config.GetString(configTree, name+".min_width", defaultMinWidth)
}

type batteryStats struct {
	Status           string
	EnergyNow        int
	EnergyFull       int
	EnergyFullDesign int
	VoltageNow       int

	// Calculated fields
	EnergyNowPerc float64
	Degradation   float64
}

func numberOfBatteries() int {
	fd, err := os.Open(powerSupplyBaseDir)
	if err != nil {
		logger.Panicln(err)
	}
	defer fd.Close()

	dirs, err := fd.Readdirnames(0)
	if err != nil {
		logger.Panicln(err)
	}

	n := 0
	for _, dir := range dirs {
		if strings.HasPrefix(dir, "BAT") {
			path := filepath.Join(powerSupplyBaseDir, dir)

			if fd, err := os.Open(path); err == nil {
				stat, err := fd.Stat()
				if err != nil {
					logger.Panicln(err)
				}

				if stat.IsDir() {
					n++
				}
			} else {
				logger.Panicln(err)
			}
		}
	}

	return n
}

func getBatteryStats() *batteryStats {
	stats := &batteryStats{}

	for i := 0; i < numberOfBatteries(); i++ {
		basepath := filepath.Join(powerSupplyBaseDir, fmt.Sprintf("BAT%d", i))
		energyFullStr, _ := ioutil.ReadFile(filepath.Join(basepath, "energy_full"))
		energyNowStr, _ := ioutil.ReadFile(filepath.Join(basepath, "energy_now"))
		energyFull, _ := strconv.Atoi(strings.TrimSpace(string(energyFullStr)))
		energyNow, _ := strconv.Atoi(strings.TrimSpace(string(energyNowStr)))

		stats.EnergyNow += energyNow
		stats.EnergyFull += energyFull
	}

	stats.EnergyNowPerc = (float64(stats.EnergyNow) / float64(stats.EnergyFull)) * 100

	return stats
}

func (c *Config) Run(args *model.ModuleArgs) {
	var outStr string
	t := template.Must(template.New("battery").Parse(c.Format))

	for range time.NewTicker(c.Period).C {
		outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
		batStats := getBatteryStats()
		buf := bytes.NewBufferString(outStr)

		if err := t.Execute(buf, batStats); err != nil {
			logger.Panicln(err)
		}

		outputBlock.FullText = buf.String()
		args.OutCh <- outputBlock
	}
}
