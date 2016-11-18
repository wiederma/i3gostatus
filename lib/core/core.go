package i3gostatus

import (
	"fmt"
	"os"
	"time"

	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/registry"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

func writeHeader() {
	header := model.NewHeader()
	fmt.Println(utils.Json(header))
	// i3bar is a streaming JSON parser, so we need to open the endless array.
	fmt.Println("[")
}

func Run(options *runtimeOptions) {
	configTree := config.Load(options.configPath)
	enabledModules := registry.Initialize(configTree)
	rateLimit := utils.FindFastestModule(configTree)
	rateTimer := time.NewTimer(rateLimit)
	outChannel := make(chan *model.I3BarBlockWrapper)
	outSlice := make([]*model.I3BarBlock, len(enabledModules))

	if len(enabledModules) == 0 {
		fmt.Fprintln(os.Stderr, "No modules are enabled!")
		os.Exit(1)
	}

	writeHeader()

	for i, v := range enabledModules {
		v.ParseConfig(configTree)
		go v.Run(outChannel, i)
	}

	for {
		select {
		case block := <-outChannel:
			outSlice[block.Index] = &block.I3BarBlock
		case <-rateTimer.C:
			rateTimer.Reset(rateLimit)
			fmt.Println(fmt.Sprintf("%s,", utils.Json(outSlice)))
		}
	}
}
