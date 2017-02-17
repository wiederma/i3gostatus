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
	defaultFormat      = `BAT{{.Index}}: {{.EnergyNowPerc | printf "%.0f" }}%`
	defaultMinWidth    = `BAT0: 90%`
	defaultIndex       = "sum"
)

var logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)

type Config struct {
	model.BaseConfig
	Index string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
	c.MinWidth = config.GetString(configTree, name+".min_width", defaultMinWidth)
	c.Index = config.GetString(configTree, name+".index", defaultIndex)
}

type batteryStats struct {
	Status           string
	Index            int
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

func getBatteryStats() []*batteryStats {
	nbats := numberOfBatteries()
	stats := make([]*batteryStats, nbats+1)
	statsSum := &batteryStats{Index: nbats}

	for i := 0; i < nbats; i++ {
		basepath := filepath.Join(powerSupplyBaseDir, fmt.Sprintf("BAT%d", i))
		energyFullStr, _ := ioutil.ReadFile(filepath.Join(basepath, "energy_full"))
		energyNowStr, _ := ioutil.ReadFile(filepath.Join(basepath, "energy_now"))
		energyFull, _ := strconv.Atoi(strings.TrimSpace(string(energyFullStr)))
		energyNow, _ := strconv.Atoi(strings.TrimSpace(string(energyNowStr)))

		statsSum.EnergyFull += energyFull
		statsSum.EnergyNow += energyNow

		stats[i] = &batteryStats{
			Index:         i,
			EnergyFull:    energyFull,
			EnergyNow:     energyNow,
			EnergyNowPerc: (float64(energyNow) / float64(energyFull)) * 100,
		}
	}

	statsSum.EnergyNowPerc = (float64(statsSum.EnergyNow) / float64(statsSum.EnergyFull)) * 100
	stats[nbats] = statsSum

	return stats
}

func (c *Config) curIndex() int {
	var showIndex int

	if c.Index == "sum" {
		showIndex = numberOfBatteries()
	} else {
		showIndex, _ = strconv.Atoi(c.Index)
	}

	return showIndex
}

func (c *Config) renderTemplate() string {
	var outStr string
	t := template.Must(template.New("battery").Parse(c.Format))
	buf := bytes.NewBufferString(outStr)
	batStats := getBatteryStats()

	if err := t.Execute(buf, batStats[c.curIndex()]); err != nil {
		logger.Panicln(err)
	}

	return buf.String()
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)

	go clickHandlers.NewListener(args, outputBlock, c)

	for range time.NewTicker(c.Period).C {
		outputBlock.FullText = c.renderTemplate()
		args.OutCh <- outputBlock
	}
}
