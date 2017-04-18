package battery

// FIXME: This code is a mess.

import (
	"time"

	"github.com/godbus/dbus"
)

const busName = "org.freedesktop.UPower"

const Unkown = 0

// Possible power sources
const (
	LinePower = 1
	Battery   = 2
	UPS       = 3
	Monitor   = 4
	Mouse     = 5
	Keyboard  = 6
	PDA       = 7
	Phone     = 8
)

// Possible states
const (
	Charging         = 1
	Discharging      = 2
	Empty            = 3
	FullyCharged     = 4
	PendingCharge    = 5
	PendingDischarge = 6
)

type Properties struct {
	NativePath       string
	Vendor           string
	Model            string
	Serial           string
	UpdateTime       uint64
	Type             uint32
	PowerSupply      bool
	HasHistory       bool
	HasStatistics    bool
	Online           bool
	Energy           float64
	EnergyEmpty      float64
	EnergyFull       float64
	EnergyFullDesign float64
	EnergyRate       float64
	Voltage          float64
	TimeToEmpty      time.Duration
	TimeToFull       time.Duration
	Percentage       float64
	IsPresent        bool
	State            uint32
	IsRechargeable   bool
	Capacity         float64
	Technology       uint32
}

func enumerateDevices() []dbus.ObjectPath {
	method := "org.freedesktop.UPower.EnumerateDevices"
	var objects []dbus.ObjectPath

	conn, err := dbus.SystemBus()
	if err != nil {
		logger.Panicln(err)
	}

	obj := conn.Object(busName, "/org/freedesktop/UPower")
	err = obj.Call(method, 0).Store(&objects)
	if err != nil {
		logger.Panicln(err)
	}

	return objects
}

func getAllProperties(dev dbus.ObjectPath) Properties {
	variants := map[string]dbus.Variant{}
	props := Properties{}

	conn, err := dbus.SystemBus()
	if err != nil {
		logger.Panicln(err)
	}

	obj := conn.Object(busName, dev)
	if err := obj.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.freedesktop.UPower.Device").Store(&variants); err != nil {
		logger.Panicln(err)
	}

	props.NativePath = variants["NativePath"].Value().(string)
	props.Vendor = variants["Vendor"].Value().(string)
	props.Model = variants["Model"].Value().(string)
	props.Serial = variants["Serial"].Value().(string)
	props.UpdateTime = variants["UpdateTime"].Value().(uint64)
	props.Type = variants["Type"].Value().(uint32)
	props.PowerSupply = variants["PowerSupply"].Value().(bool)
	props.HasHistory = variants["HasHistory"].Value().(bool)
	props.HasStatistics = variants["HasStatistics"].Value().(bool)
	props.Online = variants["Online"].Value().(bool)
	props.Energy = variants["Energy"].Value().(float64)
	props.EnergyEmpty = variants["EnergyEmpty"].Value().(float64)
	props.EnergyFull = variants["EnergyFull"].Value().(float64)
	props.EnergyFullDesign = variants["EnergyFullDesign"].Value().(float64)
	props.EnergyRate = variants["EnergyRate"].Value().(float64)
	props.Voltage = variants["Voltage"].Value().(float64)
	props.TimeToEmpty = time.Duration(time.Duration(variants["TimeToEmpty"].Value().(int64)) * time.Second)
	props.TimeToFull = time.Duration(time.Duration(variants["TimeToFull"].Value().(int64)) * time.Second)
	props.Percentage = variants["Percentage"].Value().(float64)
	props.IsPresent = variants["IsPresent"].Value().(bool)
	props.State = variants["State"].Value().(uint32)
	props.IsRechargeable = variants["IsRechargeable"].Value().(bool)
	props.Capacity = variants["Capacity"].Value().(float64)
	props.Technology = variants["Technology"].Value().(uint32)

	return props
}

func SignalChanged() chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)

	conn, err := dbus.SystemBus()
	if err != nil {
		logger.Panicln(err)
	}

	call := conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		"sender='org.freedesktop.UPower',path='/org/freedesktop/UPower/devices/DisplayDevice',type='signal'",
	)
	if call.Err != nil {
		logger.Panicln(call.Err)
	}

	conn.Signal(ch)

	return ch
}
