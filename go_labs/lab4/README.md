# Lab 4 - Smart Home System

A smart home control system implementing three behavioral design patterns: Command, Mediator, and Observer.

## Behavioral Design Patterns Overview

**Behavioral design patterns** deal with object collaboration and responsibility distribution:

- **Command** - Encapsulates a request as an object, enabling undo/redo and request queuing
- **Mediator** - Defines an object that coordinates interactions between a set of objects
- **Observer** - Defines a one-to-many dependency so that when one object changes state, all dependents are notified
- **Chain of Responsibility** - Passes requests along a chain of handlers
- **Iterator** - Provides a way to access elements of a collection sequentially
- **State** - Allows an object to alter its behavior when its internal state changes
- **Strategy** - Defines a family of algorithms, encapsulating each one and making them interchangeable
- **Template Method** - Defines the skeleton of an algorithm, deferring some steps to subclasses
- **Visitor** - Separates algorithms from the objects on which they operate
- **Memento** - Captures and restores an object's internal state

This lab implements **Command**, **Mediator**, and **Observer** patterns.

## Usage

```bash
cd lab4
go run .
```

### Example Session

```
════════════════════════════════════════════════
              SMART HOME CONTROL
════════════════════════════════════════════════

> light on
[COMMAND]
  [exec] Light ON

> scene movie
[SCENE]
  [scene] Movie Night
     -> Light dim -> 20%
     -> Thermostat -> 24°C

> motion
[SENSOR EVENT]
  [sensor] Motion Sensor detected: motion
  [log] EventLog: event=motion source=Motion Sensor
     -> Hub notified: motion_detected
     -> Light ON
     -> Alarm ARMED

> undo
[UNDO]
  [undo] Light OFF
```

## Testing

```bash
cd lab4
go test ./...
```

## Architecture & Design Patterns

### Project Structure

```
lab4/
├── domain/           # Core device types (pattern-agnostic)
│   └── device.go     # Light, Thermostat, Alarm types
├── command/          # Command pattern implementation
│   └── command.go    # Commands with Execute/Undo and History
├── mediator/         # Mediator pattern implementation
│   └── mediator.go   # SmartHub coordinates devices
├── observer/         # Observer pattern implementation
│   └── observer.go   # Sensor subjects and observers
├── main.go           # Composition root and CLI
└── README.md
```

### Domain Models

Core device types shared across all patterns:

```go
// device.go
type Device interface {
    Name() string
}

type Light struct {
    on         bool
    brightness int
}

type Thermostat struct {
    temp int
}

type Alarm struct {
    armed bool
}
```

### Command Pattern

- **Location**: `command/command.go:1-168`
- **Purpose**: Encapsulates device actions as objects with undo support

```go
// Command interface
type Command interface {
    Execute()
    Undo()
}

// History is the Invoker that stores and executes commands
type History struct {
    stack []Command
}

func (h *History) Execute(c Command) {
    c.Execute()
    h.stack = append(h.stack, c)
}

func (h *History) Undo() bool {
    if len(h.stack) == 0 {
        return false
    }
    last := h.stack[len(h.stack)-1]
    h.stack = h.stack[:len(h.stack)-1]
    last.Undo()
    return true
}
```

**Concrete Commands**:
- `LightOnCommand` / `LightOffCommand` - Toggle light state
- `DimCommand` - Adjust brightness with previous level stored
- `ThermostatCommand` - Set temperature with previous value stored
- `AlarmCommand` - Arm/disarm with previous state stored

**Benefits**:
- **Undo Support**: Each command stores previous state for reversal
- **Decoupling**: Invoker doesn't know about receivers
- **History**: Commands can be logged, queued, or replayed

### Mediator Pattern

- **Location**: `mediator/mediator.go:1-96`
- **Purpose**: Coordinates device interactions without direct coupling

```go
// Mediator interface
type Mediator interface {
    Notify(sender domain.Device, event string)
}

// SmartHub is the ConcreteMediator
type SmartHub struct {
    light      *domain.Light
    thermostat *domain.Thermostat
    alarm      *domain.Alarm
}

func (h *SmartHub) Notify(sender domain.Device, event string) {
    switch event {
    case "motion_detected":
        h.light.On()
        h.light.SetBrightness(100)
        h.alarm.Arm()
    case "door_opened":
        if h.alarm.IsArmed() {
            h.Notify(h.alarm, "alarm_triggered")
        }
    case "alarm_triggered":
        h.light.On()
        h.light.SetBrightness(100)
    }
}

// Scene executes coordinated actions
func (h *SmartHub) Scene(name string) {
    switch name {
    case "Movie Night":
        h.light.SetBrightness(20)
        h.thermostat.SetTemp(24)
    case "Away":
        h.light.Off()
        h.thermostat.SetTemp(18)
        h.alarm.Arm()
    case "Morning":
        h.light.On()
        h.thermostat.SetTemp(22)
        h.alarm.Disarm()
    }
}
```

**Supported Events**:
- `motion_detected` - Turns on light and arms alarm
- `door_opened` - Triggers alarm if armed
- `alarm_triggered` - Turns light to full brightness

**Predefined Scenes**:
- **Movie Night** - Dim lights (20%), temperature 24°C
- **Away** - Lights off, temperature 18°C, alarm armed
- **Morning** - Lights on, temperature 22°C, alarm disarmed

**Benefits**:
- **Loose Coupling**: Devices don't reference each other directly
- **Centralized Logic**: All coordination rules in one place
- **Extensibility**: New devices or interactions added without changing existing devices

### Observer Pattern

- **Location**: `observer/observer.go:1-88`
- **Purpose**: One-to-many dependency for event notification

```go
// Event payload
type Event struct {
    Source  string
    Type    string
    Payload string
}

// Observer interface
type Observer interface {
    OnNotify(event Event)
}

// Subject interface
type Subject interface {
    Subscribe(o Observer)
    Unsubscribe(o Observer)
    Notify(event Event)
}

// Sensor is the ConcreteSubject
type Sensor struct {
    name      string
    observers []Observer
}

func (s *Sensor) Detect(eventType string) {
    s.Notify(Event{
        Source:  s.name,
        Type:    eventType,
        Payload: s.name,
    })
}

// Logger is a ConcreteObserver
type Logger struct {
    name string
}

func (l *Logger) OnNotify(event Event) {
    fmt.Printf("[log] %s: event=%s source=%s\n", l.name, event.Type, event.Source)
}
```

**Components**:
- `Sensor` (Subject) - Motion sensor, door sensor
- `Logger` (Observer) - Logs all events
- `SensorHubAdapter` (Observer) - Bridges Observer to Mediator patterns

**Benefits**:
- **Decoupling**: Subjects don't know concrete observer types
- **Dynamic Subscription**: Observers can subscribe/unsubscribe at runtime
- **Broadcast**: Single event notifies multiple observers

### Pattern Integration

The patterns are integrated in `main.go` through a `SensorHubAdapter`:

```go
// SensorHubAdapter bridges Observer -> Mediator patterns
type SensorHubAdapter struct {
    hub   *mediator.SmartHub
    alarm *domain.Alarm
}

func (a *SensorHubAdapter) OnNotify(event observer.Event) {
    switch event.Type {
    case "motion":
        a.hub.Notify(a.hub.Light(), "motion_detected")
    case "door_open":
        a.hub.Notify(a.hub.Alarm(), "door_opened")
    }
}
```

**Data Flow**:
1. **Sensor** detects event (Observer Subject)
2. **SensorHubAdapter** receives notification (Observer)
3. **SmartHub** coordinates device response (Mediator)
4. **Devices** execute actions

## CLI Commands

| Command | Pattern | Description |
|---------|---------|-------------|
| `light <on\|off\|dim N>` | Command | Control light with undo support |
| `temp <N>` | Command | Set thermostat (10-35°C) |
| `alarm <arm\|disarm>` | Command | Control alarm system |
| `scene <movie\|away\|morning>` | Mediator | Run coordinated scene |
| `motion` | Observer | Simulate motion detection |
| `door` | Observer | Simulate door sensor |
| `undo` | Command | Undo last command |
| `status` | - | Show device states |
| `help` | - | Show available commands |
| `quit` | - | Exit application |

## Pattern Benefits

### Command Pattern
- Undo support for all device operations
- Decoupling between UI and execution
- Command history tracking

### Mediator Pattern
- Centralized device coordination
- Devices don't communicate directly
- Multi-device scenes as single calls

### Observer Pattern
- Event-driven sensor responses
- Add observers without modifying sensors
- Centralized event logging

## Conclusion

This implementation shows how behavioral patterns address interaction challenges in a smart home system. The **Command** pattern encapsulates device actions with undo support, the **Mediator** coordinates device interactions without coupling them, and the **Observer** enables event-driven sensor responses.

Each pattern package depends only on domain types, with integration in the composition root (`main.go`).
