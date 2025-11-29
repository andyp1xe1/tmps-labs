// Package domain defines the core device types.
package domain

// Device is the common interface for all smart home devices.
type Device interface {
	Name() string
}

// Light represents a smart light with on/off and brightness control.
type Light struct {
	on         bool
	brightness int
}

// NewLight creates a light with default state (off, 100% brightness).
func NewLight() *Light {
	return &Light{on: false, brightness: 100}
}

func (l *Light) Name() string    { return "Light" }
func (l *Light) IsOn() bool      { return l.on }
func (l *Light) Brightness() int { return l.brightness }
func (l *Light) On()             { l.on = true }
func (l *Light) Off()            { l.on = false }
func (l *Light) SetBrightness(b int) {
	l.brightness = b
	l.on = true
}

// Thermostat represents a smart thermostat with temperature control.
type Thermostat struct {
	temp int
}

// NewThermostat creates a thermostat with default temperature (20°C).
func NewThermostat() *Thermostat {
	return &Thermostat{temp: 20}
}

func (t *Thermostat) Name() string     { return "Thermostat" }
func (t *Thermostat) Temp() int        { return t.temp }
func (t *Thermostat) SetTemp(temp int) { t.temp = temp }

// Alarm represents a smart alarm system.
type Alarm struct {
	armed bool
}

// NewAlarm creates an alarm with default state (disarmed).
func NewAlarm() *Alarm {
	return &Alarm{armed: false}
}

func (a *Alarm) Name() string  { return "Alarm" }
func (a *Alarm) IsArmed() bool { return a.armed }
func (a *Alarm) Arm()          { a.armed = true }
func (a *Alarm) Disarm()       { a.armed = false }
