package i3gostatus

import (
	"fmt"
	"os"
	"time"

	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

func WriteHeader() {
	header := model.NewHeader()
	fmt.Println(utils.Json(header))
	// i3bar is a streaming JSON parser, so we need to open the endless array.
	fmt.Println("[")
}

func Run(options *runtimeOptions) {
	nModules := 0
	configTree := loadConfig(options.configPath)
	outChannel := make(chan *model.I3BarBlockWrapper)
	rateLimit := findFastestPeriod(configTree)
	rateTimer := time.NewTimer(rateLimit)

	if len(EnabledModules) == 0 {
		fmt.Println("No modules are enabled!")
		os.Exit(1)
	}

	WriteHeader()

	for i, v := range EnabledModules {
		v.ReadConfig(configTree)
		go v.Run(outChannel, i)
		nModules += 1
	}

	outputSlice := make([]*model.I3BarBlock, nModules)

	for {
		select {
		case block := <-outChannel:
			outputSlice[block.Index] = &block.I3BarBlock
		case <-rateTimer.C:
			rateTimer.Reset(rateLimit)
			fmt.Println(fmt.Sprintf("%s,", utils.Json(outputSlice)))
		}
	}
}
