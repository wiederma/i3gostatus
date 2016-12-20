package cpu

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name            = "cpu"
	moduleName      = "i3gostatus.modules." + name
	defaultFormat   = `CPU: {{.BusyPerc | printf "%.0f"}}%`
	defaultMinWidth = `CPU: 90%` // It is not so likely to get 100%, so let's save that char.
)

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.MinWidth = config.GetString(configTree, name+".min_width", defaultMinWidth)
}

type cpuStat struct {
	busyTime float64
	idleTime float64
	BusyPerc float64
	IdlePerc float64
}

func getCPUTimes() (float64, float64) {
	statFileFd, err := os.Open("/proc/stat")
	if err != nil {
		panic(err)
	}
	defer statFileFd.Close()
	r := bufio.NewReader(statFileFd)
	busyTime := 0
	idleTime := 0

	// TODO: Only read first line
	if data, _, err := r.ReadLine(); err == nil {
		cols := strings.Split(string(data), " ")

		if strings.Compare(cols[0], "cpu") == 0 {
			for i, val := range cols {
				if i == 0 || i == 1 {
					continue
				}

				t, err := strconv.Atoi(val)
				if err != nil {
					panic(err)
				}

				// See proc(5); file "/proc/stat"
				switch i {
				case 0, 1:
					continue
				case 2, 3, 4, 7, 8, 9, 10, 11:
					busyTime += t
				case 5, 6:
					idleTime += t
				}
			}
		}
	} else {
		panic(err)
	}

	return float64(busyTime), float64(idleTime)
}

func getCPUStats() *cpuStat {
	busyTime, idleTime := getCPUTimes()
	return &cpuStat{
		busyTime: busyTime,
		idleTime: idleTime,
	}
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	t := template.Must(template.New("cpu").Parse(c.Format))
	stats := &cpuStat{}
	statsPrev := &cpuStat{}
	var outStr string

	for range time.NewTicker(c.Period).C {
		buf := bytes.NewBufferString(outStr)

		stats = getCPUStats()
		// We have to compare the previous value with the current value, so we
		// have to do the calculation in the main loop.
		stats.BusyPerc = (stats.busyTime - statsPrev.busyTime) / (stats.busyTime + stats.idleTime - statsPrev.busyTime - statsPrev.idleTime) * 100
		stats.IdlePerc = 1 - stats.BusyPerc
		statsPrev = stats

		if err := t.Execute(buf, stats); err == nil {
			outputBlock.FullText = buf.String()
		} else {
			outputBlock.FullText = fmt.Sprint(err)
		}

		args.OutCh <- outputBlock
	}
}
