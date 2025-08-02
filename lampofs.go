package lampofs

import (
	"io"
	"time"
)

type LampEventType = string

const (
	READ    LampEventType = "READ"
	WRITE   LampEventType = "WRITE"
	PUT     LampEventType = "PUT"
	DELETE  LampEventType = "DELETE"
	PREPEND LampEventType = "UPDATE-PREPEND"
	APPEND  LampEventType = "UPDATE-APPEND"
)

type LampEvent struct {
	Type      LampEventType
	Path      string
	Timestamp int64
	Data      interface{}
}

type Driver interface {
	Read(path string) (io.ReadCloser, error)
	Write(path string, data []byte) error
	Put(path string, data []byte) error
	Delete(path string) error
	Update(path string, data []byte, prepend bool) error
}

type Lampo struct {
	driver Driver
	events []func(event LampEvent)
}

type LampoOption func(*Lampo)

func NewLampo(driver Driver, opts ...LampoOption) *Lampo {
	lampo := &Lampo{
		driver: driver,
		events: make([]func(LampEvent), 0),
	}

	for _, opt := range opts {
		opt(lampo)
	}

	return lampo
}

func (l *Lampo) On(handler func(event LampEvent)) {
	l.events = append(l.events, handler)
}

func (l *Lampo) Read(path string) (io.ReadCloser, error) {
	reader, err := l.driver.Read(path)
	if err != nil {
		return nil, err
	}

	l.fireEvent(LampEvent{
		Type:      READ,
		Path:      path,
		Timestamp: time.Now().Unix(),
	})

	return reader, nil
}

func (l *Lampo) Write(path string, data []byte) error {
	err := l.driver.Write(path, data)
	if err != nil {
		return err
	}

	l.fireEvent(LampEvent{
		Type:      WRITE,
		Path:      path,
		Timestamp: time.Now().Unix(),
		Data:      len(data),
	})

	return nil
}

func (l *Lampo) Put(path string, data []byte) error {
	err := l.driver.Put(path, data)
	if err != nil {
		return err
	}

	l.fireEvent(LampEvent{
		Type:      PUT,
		Path:      path,
		Timestamp: time.Now().Unix(),
		Data:      len(data),
	})

	return nil
}

func (l *Lampo) Delete(path string) error {
	err := l.driver.Delete(path)
	if err != nil {
		return err
	}

	l.fireEvent(LampEvent{
		Type:      DELETE,
		Path:      path,
		Timestamp: time.Now().Unix(),
	})

	return nil
}

func (l *Lampo) Update(path string, data []byte, prepend bool) error {
	err := l.driver.Update(path, data, prepend)
	if err != nil {
		return err
	}

	action := APPEND
	if prepend {
		action = PREPEND
	}

	l.fireEvent(LampEvent{
		Type:      action,
		Path:      path,
		Timestamp: time.Now().Unix(),
		Data:      len(data),
	})

	return nil
}

func (l *Lampo) fireEvent(event LampEvent) {
	for _, handler := range l.events {
		handler(event)
	}
}
