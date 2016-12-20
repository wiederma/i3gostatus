// Package backlight outputs the current backlight brightness in percent by
// spawning `xbacklight` processes and parsing its output.
package backlight

import (
	"fmt"
	"log"
	"os"
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
	defaultFormat = "%.0f%%"
)

type Config struct {
	model.BaseConfig
}

var logger *log.Logger
var xbacklight string

func init() {
	logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)
	xbacklight = utils.Which("xbacklight")
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
}

func getBrightness() float64 {
	brightness, _ := exec.Command(xbacklight, "-get").Output()
	res, _ := strconv.ParseFloat(strings.TrimSpace(string(brightness)), 64)

	return res
}

func setBrightness(val int) {
	err := exec.Command(xbacklight, "set", strconv.Itoa(val)).Run()
	if err != nil {
		logger.Printf("Command failed: %s\n", err)
	}
}

func incBrightness(val int) {
	err := exec.Command(xbacklight, "-inc", strconv.Itoa(val)).Run()
	if err != nil {
		logger.Printf("Command failed: %s\n", err)
	}
}

func decBrightness(val int) {
	err := exec.Command(xbacklight, "-dec", strconv.Itoa(val)).Run()
	if err != nil {
		logger.Printf("Command failed: %s\n", err)
	}
}

func (c *Config) Run(args *model.ModuleArgs) {
	if xbacklight == "" {
		logger.Println("xbacklight is not available. Terminating backlight module.")
		return
	}

	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	go clickHandlers.NewListener(args, outputBlock, c)

	for range time.NewTicker(c.Period).C {
		outputBlock.FullText = fmt.Sprintf(c.Format, getBrightness())
		args.OutCh <- outputBlock
	}
}
