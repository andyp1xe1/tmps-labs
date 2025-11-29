// Package command implements the Command pattern.
// Encapsulates requests as objects with Execute/Undo operations.
package command

import (
	"fmt"

	"smarthome/domain"
)

// Command is the command interface declaring Execute and Undo operations.
type Command interface {
	Execute()
	Undo()
}

// History is the Invoker that stores and executes commands.
type History struct {
	stack []Command
}

// Execute runs a command and adds it to history.
func (h *History) Execute(c Command) {
	c.Execute()
	h.stack = append(h.stack, c)
}

// Undo reverses the last command.
func (h *History) Undo() bool {
	if len(h.stack) == 0 {
		return false
	}
	last := h.stack[len(h.stack)-1]
	h.stack = h.stack[:len(h.stack)-1]
	last.Undo()
	return true
}

// Len returns the number of commands in history.
func (h *History) Len() int {
	return len(h.stack)
}

// LightOnCommand turns the light on.
type LightOnCommand struct {
	light *domain.Light
	wasOn bool
}

func NewLightOnCommand(l *domain.Light) *LightOnCommand {
	return &LightOnCommand{light: l}
}

func (c *LightOnCommand) Execute() {
	c.wasOn = c.light.IsOn()
	c.light.On()
	fmt.Println("  [exec] Light ON")
}

func (c *LightOnCommand) Undo() {
	if !c.wasOn {
		c.light.Off()
	}
	fmt.Println("  [undo] Light OFF")
}

// LightOffCommand turns the light off.
type LightOffCommand struct {
	light *domain.Light
	wasOn bool
}

func NewLightOffCommand(l *domain.Light) *LightOffCommand {
	return &LightOffCommand{light: l}
}

func (c *LightOffCommand) Execute() {
	c.wasOn = c.light.IsOn()
	c.light.Off()
	fmt.Println("  [exec] Light OFF")
}

func (c *LightOffCommand) Undo() {
	if c.wasOn {
		c.light.On()
	}
	fmt.Println("  [undo] Light ON")
}

// DimCommand sets light brightness.
type DimCommand struct {
	light     *domain.Light
	level     int
	prevLevel int
}

func NewDimCommand(l *domain.Light, level int) *DimCommand {
	return &DimCommand{light: l, level: level}
}

func (c *DimCommand) Execute() {
	c.prevLevel = c.light.Brightness()
	c.light.SetBrightness(c.level)
	fmt.Printf("  [exec] Light dim -> %d%%\n", c.level)
}

func (c *DimCommand) Undo() {
	c.light.SetBrightness(c.prevLevel)
	fmt.Printf("  [undo] Light dim -> %d%%\n", c.prevLevel)
}

// ThermostatCommand sets the thermostat temperature.
type ThermostatCommand struct {
	thermostat *domain.Thermostat
	temp       int
	prevTemp   int
}

func NewThermostatCommand(t *domain.Thermostat, temp int) *ThermostatCommand {
	return &ThermostatCommand{thermostat: t, temp: temp}
}

func (c *ThermostatCommand) Execute() {
	c.prevTemp = c.thermostat.Temp()
	c.thermostat.SetTemp(c.temp)
	fmt.Printf("  [exec] Thermostat -> %d°C\n", c.temp)
}

func (c *ThermostatCommand) Undo() {
	c.thermostat.SetTemp(c.prevTemp)
	fmt.Printf("  [undo] Thermostat -> %d°C\n", c.prevTemp)
}

// AlarmCommand arms or disarms the alarm.
type AlarmCommand struct {
	alarm    *domain.Alarm
	arm      bool
	wasArmed bool
}

func NewAlarmCommand(a *domain.Alarm, arm bool) *AlarmCommand {
	return &AlarmCommand{alarm: a, arm: arm}
}

func (c *AlarmCommand) Execute() {
	c.wasArmed = c.alarm.IsArmed()
	if c.arm {
		c.alarm.Arm()
		fmt.Println("  [exec] Alarm ARMED")
	} else {
		c.alarm.Disarm()
		fmt.Println("  [exec] Alarm DISARMED")
	}
}

func (c *AlarmCommand) Undo() {
	if c.wasArmed {
		c.alarm.Arm()
		fmt.Println("  [undo] Alarm ARMED")
	} else {
		c.alarm.Disarm()
		fmt.Println("  [undo] Alarm DISARMED")
	}
}
