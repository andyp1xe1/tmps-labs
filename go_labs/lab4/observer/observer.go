// Package observer implements the Observer pattern.
// Defines a one-to-many dependency where subjects notify observers of events.
package observer

import "fmt"

// Event represents a notification payload.
type Event struct {
	Source  string
	Type    string
	Payload string
}

// Observer is the observer interface - objects that want to be notified.
type Observer interface {
	OnNotify(event Event)
}

// Subject is the subject interface - objects that notify observers.
type Subject interface {
	Subscribe(o Observer)
	Unsubscribe(o Observer)
	Notify(event Event)
}

// Sensor is the ConcreteSubject that detects events and notifies observers.
type Sensor struct {
	name      string
	observers []Observer
}

// NewSensor creates a new sensor with a name.
func NewSensor(name string) *Sensor {
	return &Sensor{name: name}
}

func (s *Sensor) Name() string { return s.name }

// Subscribe adds an observer to the notification list.
func (s *Sensor) Subscribe(o Observer) {
	s.observers = append(s.observers, o)
}

// Unsubscribe removes an observer from the notification list.
func (s *Sensor) Unsubscribe(o Observer) {
	for i, obs := range s.observers {
		if obs == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return
		}
	}
}

// Notify sends an event to all observers.
func (s *Sensor) Notify(event Event) {
	for _, o := range s.observers {
		o.OnNotify(event)
	}
}

// Detect simulates the sensor detecting something and notifying observers.
func (s *Sensor) Detect(eventType string) {
	fmt.Printf("  [sensor] %s detected: %s\n", s.name, eventType)
	s.Notify(Event{
		Source:  s.name,
		Type:    eventType,
		Payload: s.name,
	})
}

// Logger is a ConcreteObserver that logs all events.
type Logger struct {
	name string
}

// NewLogger creates a new logger observer.
func NewLogger(name string) *Logger {
	return &Logger{name: name}
}

// OnNotify implements the Observer interface.
func (l *Logger) OnNotify(event Event) {
	fmt.Printf("  [log] %s: event=%s source=%s\n", l.name, event.Type, event.Source)
}
