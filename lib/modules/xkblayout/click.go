package xkblayout

import (
	"fmt"

	"github.com/rumpelsepp/i3gostatus/lib/model"
)

var nextIndex int = 0
var prevIndex int = 0

var clickHandlers = &model.ClickHandlers{
	HandleRightClick: onRightClick,
	HandleLeftClick:  onLeftClick,
}

func determineIndexes(layouts []string) {
	for i, l := range layouts {
		if l == queryCurrentLayout() {
			switch i {
			case 0:
				nextIndex = i + 1
				prevIndex = len(layouts) - 1
			// curIndex == last element
			case len(layouts) - 1:
				nextIndex = 0
				prevIndex = i - 1
			default:
				nextIndex = i + 1
				prevIndex = i - 1
			}

			return
		}
	}
}

func onRightClick(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)
	determineIndexes(config.Layouts)

	setLayout(config.Layouts[prevIndex])
	block.FullText = fmt.Sprintf(config.Format, queryCurrentLayout())
	args.ClickEventCh <- block
}

func onLeftClick(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)
	determineIndexes(config.Layouts)

	setLayout(config.Layouts[nextIndex])
	block.FullText = fmt.Sprintf(config.Format, queryCurrentLayout())
	args.ClickEventCh <- block
}
