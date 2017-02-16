// https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=839940
// https://github.com/falconindy/ponymix

package pulseaudio

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
	name       = "pulseaudio"
	moduleName = "i3gostatus.modules." + name
	// TODO: Use template/text
	defaultFormat = "♪: %s"
)

var logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)
var ponymix string

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
}

func getCurrentVolume() string {
	out, err := exec.Command(ponymix, "get-volume").Output()
	if err != nil {
		logger.Printf("Failed querying volume: %s\n", err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func isMuted() bool {
	// ponymix is-muted returns 1 or 0. Therefore if there is no error
	// the output sink is muted.
	if err := exec.Command(ponymix, "is-muted").Run(); err == nil {
		return true
	}
	return false
}

func toggleMute() error {
	if err := exec.Command(ponymix, "toggle").Run(); err != nil {
		logger.Printf("Failed toggle device: %s\n", err)
		return err
	}

	return nil
}

func increaseVolume(val int) error {
	valStr := strconv.Itoa(val)
	if err := exec.Command(ponymix, "increase", valStr).Run(); err != nil {
		logger.Printf("increaseVolume() failed: %s\n", err)
		return err
	}
	return nil
}

func decreaseVolume(val int) error {
	valStr := strconv.Itoa(val)
	if err := exec.Command(ponymix, "decrease", valStr).Run(); err != nil {
		logger.Printf("decreaseVolume() failed: %s\n", err)
		return err
	}
	return nil
}

func (c *Config) Run(args *model.ModuleArgs) {
	ponymix = utils.Which("ponymix")
	if ponymix == "" {
		logger.Println("ponymix is not available; terminating module.")
		logger.Println("https://github.com/falconindy/ponymix")
		return
	}

	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	go clickHandlers.NewListener(args, outputBlock, c)

	for range time.NewTicker(c.Period).C {
		if isMuted() {
			outputBlock.Color = "#FF0000"
		} else {
			// Use standard color; unsetting it let the JSON
			// marshaller omit the relevant field.
			outputBlock.Color = ""
		}

		outputBlock.FullText = fmt.Sprintf(c.Format, getCurrentVolume())
		args.OutCh <- outputBlock
	}
}
