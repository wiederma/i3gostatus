package i3gostatus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/registry"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

func writeHeader(options *runtimeOptions) {
	header := model.NewHeader(options.clickEvents)
	fmt.Println(utils.Json(header))
	// i3bar is a streaming JSON parser, so we need to open the endless array.
	fmt.Println("[")
}

func readStdin(out chan *model.I3ClickEvent) {
	scanner := bufio.NewScanner(os.Stdin)
	var inputStr string

	for scanner.Scan() {
		// Trim the endless JSON array stuff. It causes parse errors,
		// since we do line by line JSON parsing here.
		inputStr = strings.Trim(scanner.Text(), "[, ")
		clickEvent := &model.I3ClickEvent{}

		if err := json.Unmarshal([]byte(inputStr), clickEvent); err == nil {
			out <- clickEvent
		}
	}
}

func Run(options *runtimeOptions) {
	configTree := config.Load(options.configPath)
	enabledModules := registry.Initialize(configTree)
	rateLimit := utils.FindFastestModule(configTree)
	rateTimer := time.NewTimer(rateLimit)
	outChannel := make(chan *model.I3BarBlockWrapper)
	outSlice := make([]*model.I3BarBlock, len(enabledModules))
	// inChannel is only used when click_events is enabled.
	// If click_events is disabled, it is never written to
	// the channel.
	inChannel := make(chan *model.I3ClickEvent)

	if len(enabledModules) == 0 {
		fmt.Fprintln(os.Stderr, "No modules are enabled!")
		os.Exit(1)
	}

	writeHeader(options)

	if options.clickEvents {
		go readStdin(inChannel)
	}

	for i, v := range enabledModules {
		v.ParseConfig(configTree)
		go v.Run(outChannel, inChannel, i)
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
