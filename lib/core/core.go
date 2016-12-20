package i3gostatus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/registry"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[i3gostatus] ", log.LstdFlags)
}

func writeHeader(options *runtimeOptions) {
	header := model.NewHeader(options.clickEvents)
	fmt.Println(utils.Json(header))
	// i3bar is a streaming JSON parser, so we need to open the endless array.
	fmt.Println("[")
}

func readStdin(outChannels map[string]chan *model.I3ClickEvent) {
	scanner := bufio.NewScanner(os.Stdin)
	var inputStr string

	for scanner.Scan() {
		// Trim the endless JSON array stuff. It causes parse errors,
		// since we do line by line JSON parsing here.
		inputStr = strings.Trim(scanner.Text(), "[, ")
		clickEvent := &model.I3ClickEvent{}

		if err := json.Unmarshal([]byte(inputStr), clickEvent); err == nil {
			outChannels[clickEvent.Instance] <- clickEvent
		}
	}
}

func Run(options *runtimeOptions) {
	configTree := config.Load(options.configPath)
	enabledModules := registry.Initialize(configTree)
	rateLimit := utils.FindFastestModule(configTree)
	rateTimer := time.NewTimer(rateLimit)
	outChannel := make(chan *model.I3BarBlock)
	clickEventChannel := make(chan *model.I3BarBlock)
	outSlice := make([]*model.I3BarBlock, len(enabledModules))
	// The relevant inChannel is only used when click_events is enabled.
	// If click_events is disabled, it is never written to  the channel.
	inChannels := make(map[string]chan *model.I3ClickEvent)

	logger.Printf("Runtime options set: %+v", options)

	if len(enabledModules) == 0 {
		fmt.Fprintln(os.Stderr, "No modules are enabled!")
		os.Exit(1)
	}

	writeHeader(options)

	for i, v := range enabledModules {
		v.ParseConfig(configTree)
		id := reflect.ValueOf(v).Elem().FieldByName("Instance").String()
		inChannel := make(chan *model.I3ClickEvent)
		go v.Run(&model.ModuleArgs{inChannel, outChannel, clickEventChannel, i})
		// Add it it to the channel map. The click_event handler must be able
		// to somehow find the correct channel.
		inChannels[id] = inChannel
	}

	if options.clickEvents {
		go readStdin(inChannels)
	}

	for {
		select {
		case block := <-outChannel:
			outSlice[block.Index] = block
		case block := <-clickEventChannel:
			outSlice[block.Index] = block
			fmt.Println(fmt.Sprintf("%s,", utils.Json(outSlice)))
		case <-rateTimer.C:
			rateTimer.Reset(rateLimit)
			fmt.Println(fmt.Sprintf("%s,", utils.Json(outSlice)))
		}
	}
}
