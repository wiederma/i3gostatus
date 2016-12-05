package syncthing

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

const (
	name       = "syncthing"
	moduleName = "i3gostatus.modules." + name
)

var logger *log.Logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)
var xdgOpen string = utils.Which("xdg-open")

type Config struct {
	model.BaseConfig
	STUrl      string
	UpString   string
	UpColor    string
	DownString string
	DownColor  string
}

func startSyncthing() {
	systemctl := utils.Which("systemctl")
	if systemctl == "" {
		logger.Println("systemctl is not available")
		return
	}

	exec.Command(systemctl, "--user", "start", "syncthing.service").CombinedOutput()
}

func stopSyncthing() {
	systemctl := utils.Which("systemctl")
	if systemctl == "" {
		logger.Println("systemctl is not available")
		return
	}

	exec.Command(systemctl, "--user", "stop", "syncthing.service").CombinedOutput()
}

func isUp(url string) bool {
	var up bool

	if csrfToken == "" {
		initHTTPSession(url)
	}

	if resp, err := stGet(url, "/rest/system/ping"); err == nil {
		// I do not feel motivated to parse JSON now...
		// This should suffice in most cases.
		if resp == `{"ping":"pong"}` {
			up = true
		} else {
			up = false
		}
	} else if _, ok := err.(noActiveSessionError); ok {
		logger.Printf("Warning: %s", err)
		logger.Println("Renewing http session...")
		initHTTPSession(url)
		up = false
	} else {
		up = false
	}

	return up
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Period = config.GetDurationMs(configTree, c.Name+".period", 10000)
	c.STUrl = config.GetString(configTree, name+".st_url", "http://localhost:8384")
	c.UpString = config.GetString(configTree, name+".up_string", "ST UP")
	c.DownString = config.GetString(configTree, name+".down_string", "ST DOWN")
	c.UpColor = config.GetString(configTree, name+".down_string", "#00FF00")
	c.DownColor = config.GetString(configTree, name+".down_string", "#FF0000")
}

func (c *Config) Run(args *model.ModuleArgs) {
	logger.Println("Started Syncthing module")
	logger.Printf("Configuration: %+v\n", c)

	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	initHTTPSession(c.STUrl)

	if xdgOpen != "" {
		go clickHandlers.NewListener(args, outputBlock, c)
	} else {
		logger.Println("xdg-open is not available.")
		logger.Println("No click handler available")
	}

	for range time.NewTicker(c.Period).C {
		if isUp(c.STUrl) {
			outputBlock.Color = c.UpColor
			outputBlock.FullText = c.UpString
		} else {
			outputBlock.Color = c.DownColor
			outputBlock.FullText = c.DownString
		}

		args.OutCh <- outputBlock
	}
}
