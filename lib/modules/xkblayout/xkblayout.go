package xkblayout

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

const (
	name          = "xkblayout"
	moduleName    = "i3gostatus.modules." + name
	defaultFormat = "Xkbd: %s"
	defaultPeriod = 10000
)

var logger *log.Logger
var setxkbmap string

func init() {
	logger = log.New(os.Stderr, "[xkblayout] ", log.LstdFlags)
}

type Config struct {
	model.BaseConfig
	Layouts []string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.Period = config.GetDurationMs(configTree, c.Name+".period", defaultPeriod)
	c.Layouts = config.GetStringList(configTree, name+".layouts", []string{})
}

func queryCurrentLayout() string {
	res := ""

	// We can use this, as the main entry point of this module
	// checks if the variable is actually set.
	out, err := exec.Command(setxkbmap, "-query").Output()
	if err != nil {
		logger.Fatalf("Error occured: %s", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "layout") {
			res = strings.TrimSpace(strings.Split(line, ":")[1])
		}

		if strings.HasPrefix(line, "variant") {
			res += " " + strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return res
}

func setLayout(spec string) {
	s := strings.SplitAfter(spec, " ")

	switch len(s) {
	case 1:
		exec.Command(setxkbmap, "-layout", s[0]).CombinedOutput()
	case 2:
		exec.Command(setxkbmap, "-layout", s[0], "-variant", s[1]).CombinedOutput()
	default:
		logger.Println("Layoutstring is broken")
		return
	}
}

func (c *Config) Run(args *model.ModuleArgs) {
	setxkbmap = utils.Which("setxkbmap")
	if setxkbmap == "" {
		logger.Println("setxkbmap is not available; terminating module.")
		return
	}

	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)

	if len(c.Layouts) > 0 {
		go clickHandlers.NewListener(args, outputBlock, c)
	} else {
		logger.Println("No layouts sepcified.")
		logger.Println("No click handler available.")
	}

	for range time.NewTicker(c.Period).C {
		outputBlock.FullText = fmt.Sprintf(c.Format, queryCurrentLayout())
		args.OutCh <- outputBlock
	}
}
