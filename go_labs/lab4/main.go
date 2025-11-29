// Smart Home System - Behavioral Design Patterns Demo
//
// This project demonstrates three Gang of Four behavioral patterns:
//   - Command:  Encapsulates device actions with undo support
//   - Mediator: SmartHub coordinates devices without direct coupling
//   - Observer: Sensors notify subscribers of detected events
//
// Architecture:
//   - domain/   : Shared device types (pattern-agnostic)
//   - command/  : Command pattern (depends on domain)
//   - mediator/ : Mediator pattern (depends on domain)
//   - observer/ : Observer pattern (self-contained)
//   - main.go   : Composition root - wires patterns together
//
// No cross-pattern dependencies - patterns only share domain types.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"smarthome/command"
	"smarthome/domain"
	"smarthome/mediator"
	"smarthome/observer"
)

// SensorHubAdapter bridges Observer -> Mediator patterns.
// Lives in main.go (composition root) to avoid cross-pattern imports.
type SensorHubAdapter struct {
	hub   *mediator.SmartHub
	alarm *domain.Alarm
}

// OnNotify implements observer.Observer interface.
func (a *SensorHubAdapter) OnNotify(event observer.Event) {
	switch event.Type {
	case "motion":
		a.hub.Notify(a.hub.Light(), "motion_detected")
	case "door_open":
		fmt.Println("     -> Door opened, checking alarm status")
		a.hub.Notify(a.hub.Alarm(), "door_opened")
	}
}

// SmartHome integrates all three patterns.
type SmartHome struct {
	// Domain devices (single source of truth)
	light      *domain.Light
	thermostat *domain.Thermostat
	alarm      *domain.Alarm

	// Patterns
	hub     *mediator.SmartHub // Mediator
	history *command.History   // Command invoker

	// Sensors (Observer subjects)
	motionSensor *observer.Sensor
	doorSensor   *observer.Sensor
	logger       *observer.Logger
}

func NewSmartHome() *SmartHome {
	// Create domain devices (shared state)
	light := domain.NewLight()
	thermostat := domain.NewThermostat()
	alarm := domain.NewAlarm()

	home := &SmartHome{
		// Domain
		light:      light,
		thermostat: thermostat,
		alarm:      alarm,

		// Patterns - all use the same domain devices
		hub:     mediator.NewSmartHub(light, thermostat, alarm),
		history: &command.History{},

		// Sensors
		motionSensor: observer.NewSensor("Motion Sensor"),
		doorSensor:   observer.NewSensor("Door Sensor"),
		logger:       observer.NewLogger("EventLog"),
	}

	// Wire up Observer -> Mediator via adapter
	hubAdapter := &SensorHubAdapter{hub: home.hub, alarm: home.alarm}

	// Subscribe adapter and logger to sensors
	home.motionSensor.Subscribe(hubAdapter)
	home.motionSensor.Subscribe(home.logger)
	home.doorSensor.Subscribe(hubAdapter)
	home.doorSensor.Subscribe(home.logger)

	return home
}

func main() {
	home := NewSmartHome()
	reader := bufio.NewReader(os.Stdin)

	printHeader()
	printHelp()

	for {
		fmt.Print("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToLower(parts[0])
		args := parts[1:]

		switch cmd {
		case "help", "h":
			printHelp()

		case "status", "s":
			home.printStatus()

		case "light":
			home.handleLight(args)

		case "temp":
			home.handleTemp(args)

		case "alarm":
			home.handleAlarm(args)

		case "scene":
			home.handleScene(args)

		case "motion":
			fmt.Println("\n[SENSOR EVENT]")
			home.motionSensor.Detect("motion")

		case "door":
			fmt.Println("\n[SENSOR EVENT]")
			home.doorSensor.Detect("door_open")

		case "undo", "u":
			fmt.Println("\n[UNDO]")
			if !home.history.Undo() {
				fmt.Println("  Nothing to undo")
			}

		case "quit", "q", "exit":
			fmt.Println("\nGoodbye!")
			return

		default:
			fmt.Printf("  Unknown command: %s (type 'help')\n", cmd)
		}
	}
}

func (h *SmartHome) handleLight(args []string) {
	if len(args) == 0 {
		fmt.Println("  Usage: light <on|off|dim N>")
		return
	}

	fmt.Println("\n[COMMAND]")
	switch args[0] {
	case "on":
		h.history.Execute(command.NewLightOnCommand(h.light))
	case "off":
		h.history.Execute(command.NewLightOffCommand(h.light))
	case "dim":
		if len(args) < 2 {
			fmt.Println("  Usage: light dim <0-100>")
			return
		}
		level, err := strconv.Atoi(args[1])
		if err != nil || level < 0 || level > 100 {
			fmt.Println("  Invalid brightness level")
			return
		}
		h.history.Execute(command.NewDimCommand(h.light, level))
	default:
		fmt.Println("  Usage: light <on|off|dim N>")
	}
}

func (h *SmartHome) handleTemp(args []string) {
	if len(args) == 0 {
		fmt.Println("  Usage: temp <temperature>")
		return
	}

	temp, err := strconv.Atoi(args[0])
	if err != nil || temp < 10 || temp > 35 {
		fmt.Println("  Invalid temperature (10-35)")
		return
	}

	fmt.Println("\n[COMMAND]")
	h.history.Execute(command.NewThermostatCommand(h.thermostat, temp))
}

func (h *SmartHome) handleAlarm(args []string) {
	if len(args) == 0 {
		fmt.Println("  Usage: alarm <arm|disarm>")
		return
	}

	fmt.Println("\n[COMMAND]")
	switch args[0] {
	case "arm":
		h.history.Execute(command.NewAlarmCommand(h.alarm, true))
	case "disarm":
		h.history.Execute(command.NewAlarmCommand(h.alarm, false))
	default:
		fmt.Println("  Usage: alarm <arm|disarm>")
	}
}

func (h *SmartHome) handleScene(args []string) {
	if len(args) == 0 {
		fmt.Println("  Scenes: movie, away, morning")
		return
	}

	fmt.Println("\n[SCENE]")
	switch strings.ToLower(args[0]) {
	case "movie":
		h.hub.Scene("Movie Night")
	case "away":
		h.hub.Scene("Away")
	case "morning":
		h.hub.Scene("Morning")
	default:
		fmt.Println("  Unknown scene. Available: movie, away, morning")
	}
}

func (h *SmartHome) printStatus() {
	fmt.Println("\n┌─────────────────────────────────────┐")
	fmt.Println("│           DEVICE STATUS             │")
	fmt.Println("├─────────────────────────────────────┤")

	if h.light.IsOn() {
		fmt.Printf("│  Light:      ON (%d%% brightness)    │\n", h.light.Brightness())
	} else {
		fmt.Println("│  Light:      OFF                    │")
	}

	fmt.Printf("│  Thermostat: %d°C                    │\n", h.thermostat.Temp())

	if h.alarm.IsArmed() {
		fmt.Println("│  Alarm:      ARMED                  │")
	} else {
		fmt.Println("│  Alarm:      DISARMED               │")
	}

	fmt.Printf("│  History:    %d commands             │\n", h.history.Len())
	fmt.Println("└─────────────────────────────────────┘")
}

func printHeader() {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════")
	fmt.Println("              SMART HOME CONTROL")
	fmt.Println("════════════════════════════════════════════════")
}

func printHelp() {
	fmt.Println()
	fmt.Println("  COMMANDS")
	fmt.Println("  ────────────────────────────────────")
	fmt.Println("  light <on|off|dim N>  Control light")
	fmt.Println("  temp <N>              Set temperature")
	fmt.Println("  alarm <arm|disarm>    Control alarm")
	fmt.Println("  scene <name>          Run scene (movie/away/morning)")
	fmt.Println("  motion                Simulate motion sensor")
	fmt.Println("  door                  Simulate door sensor")
	fmt.Println("  undo                  Undo last command")
	fmt.Println("  status                Show device status")
	fmt.Println("  help                  Show this help")
	fmt.Println("  quit                  Exit")
}
