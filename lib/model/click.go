package model

type ClickHandlers struct {
	HandleLeftClick   func(*ModuleArgs, *I3BarBlock, interface{})
	HandleMiddleClick func(*ModuleArgs, *I3BarBlock, interface{})
	HandleRightClick  func(*ModuleArgs, *I3BarBlock, interface{})
	HandleWheelUp     func(*ModuleArgs, *I3BarBlock, interface{})
	HandleWheelDown   func(*ModuleArgs, *I3BarBlock, interface{})
}

func (handlers *ClickHandlers) NewListener(args *ModuleArgs, block *I3BarBlock, data interface{}) {
	for event := range args.InCh {
		switch event.Button {
		case MouseButtonLeft:
			if handlers.HandleLeftClick != nil {
				handlers.HandleLeftClick(args, block, data)
			}
		case MouseButtonMiddle:
			if handlers.HandleMiddleClick != nil {
				handlers.HandleMiddleClick(args, block, data)
			}
		case MouseButtonRight:
			if handlers.HandleRightClick != nil {
				handlers.HandleRightClick(args, block, data)
			}
		case MouseWheelUp:
			if handlers.HandleWheelUp != nil {
				handlers.HandleWheelUp(args, block, data)
			}
		case MouseWheelDown:
			if handlers.HandleWheelDown != nil {
				handlers.HandleWheelDown(args, block, data)
			}
		default:
			continue
		}
	}
}
