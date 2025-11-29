// Package mediator implements the Mediator pattern.
// Encapsulates how devices interact through a central hub.
package mediator

import (
	"fmt"

	"smarthome/domain"
)

// Mediator defines the interface for communication between devices.
type Mediator interface {
	Notify(sender domain.Device, event string)
}

// SmartHub is the ConcreteMediator coordinating all devices.
type SmartHub struct {
	light      *domain.Light
	thermostat *domain.Thermostat
	alarm      *domain.Alarm
}

// NewSmartHub creates a new hub with the given devices.
func NewSmartHub(light *domain.Light, thermostat *domain.Thermostat, alarm *domain.Alarm) *SmartHub {
	return &SmartHub{
		light:      light,
		thermostat: thermostat,
		alarm:      alarm,
	}
}

// Devices returns all hub devices for external use.
func (h *SmartHub) Devices() []domain.Device {
	return []domain.Device{h.light, h.thermostat, h.alarm}
}

// Accessors for devices
func (h *SmartHub) Light() *domain.Light           { return h.light }
func (h *SmartHub) Thermostat() *domain.Thermostat { return h.thermostat }
func (h *SmartHub) Alarm() *domain.Alarm           { return h.alarm }

// Notify handles inter-device communication through the hub.
func (h *SmartHub) Notify(sender domain.Device, event string) {
	switch event {
	case "motion_detected":
		fmt.Printf("     -> Hub notified: %s\n", event)
		h.light.On()
		h.light.SetBrightness(100)
		fmt.Println("     -> Light ON")
		h.alarm.Arm()
		fmt.Println("     -> Alarm ARMED")
	case "alarm_triggered":
		fmt.Printf("     -> Hub notified: %s\n", event)
		h.light.On()
		h.light.SetBrightness(100)
		fmt.Println("     -> Light ON (alarm response)")
	case "door_opened":
		fmt.Printf("     -> Hub notified: %s\n", event)
		if h.alarm.IsArmed() {
			fmt.Println("     -> Alarm TRIGGERED!")
			h.Notify(h.alarm, "alarm_triggered")
		}
	}
}

// Scene executes a predefined scene (coordinated actions).
func (h *SmartHub) Scene(name string) {
	fmt.Printf("\n  [scene] %s\n", name)
	switch name {
	case "Movie Night":
		h.light.SetBrightness(20)
		fmt.Println("     -> Light dim -> 20%")
		h.thermostat.SetTemp(24)
		fmt.Printf("     -> Thermostat -> 24°C\n")
	case "Away":
		h.light.Off()
		fmt.Println("     -> Light OFF")
		h.thermostat.SetTemp(18)
		fmt.Printf("     -> Thermostat -> 18°C\n")
		h.alarm.Arm()
		fmt.Println("     -> Alarm ARMED")
	case "Morning":
		h.light.On()
		h.light.SetBrightness(100)
		fmt.Println("     -> Light ON")
		h.thermostat.SetTemp(22)
		fmt.Printf("     -> Thermostat -> 22°C\n")
		h.alarm.Disarm()
		fmt.Println("     -> Alarm DISARMED")
	}
}
