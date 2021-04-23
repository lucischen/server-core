package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/francoispqt/onelog"
)

type LogField interface {
	String(string, interface{}) LogField
	WithCaller() LogField
	Submit()
}

// Mock
type Mock struct{}

func (m *Mock) String(k string, v interface{}) LogField {
	return m
}

func (m *Mock) WithCaller() LogField {
	return m
}

func (m *Mock) Submit() {}

// Error
type Error struct {
	logger      *onelog.Logger
	callerLevel int
	title       string
	m           map[string]string
	withCaller  bool
}

func (e *Error) String(k string, v interface{}) LogField {
	e.m[k] = fmt.Sprintf("%+v", v)

	return e
}

func (e *Error) WithCaller() LogField {
	e.withCaller = true
	return e
}

func (e *Error) Submit() {
	e.logger.ErrorWithFields(e.title, func(oe onelog.Entry) {
		for k, v := range e.m {
			oe.String(k, v)
		}

		if e.withCaller {
			oe.String("CALLER", generateCallerList(e.callerLevel))
		}
	})
}

// Info
type Info struct {
	logger      *onelog.Logger
	callerLevel int
	title       string
	m           map[string]string
	withCaller  bool
}

func (i *Info) String(k string, v interface{}) LogField {
	i.m[k] = fmt.Sprintf("%+v", v)

	return i
}

func (i *Info) WithCaller() LogField {
	i.withCaller = true
	return i
}

func (i *Info) Submit() {
	i.logger.InfoWithFields(i.title, func(oe onelog.Entry) {
		for k, v := range i.m {
			oe.String(k, v)
		}
		if i.withCaller {
			oe.String("CALLER", generateCallerList(i.callerLevel))
		}
	})
}

// Debug
type Debug struct {
	logger      *onelog.Logger
	callerLevel int
	title       string
	m           map[string]string
	withCaller  bool
}

func (d *Debug) String(k string, v interface{}) LogField {
	d.m[k] = fmt.Sprintf("%+v", v)

	return d
}

func (d *Debug) WithCaller() LogField {
	d.withCaller = true
	return d
}

func (d *Debug) Submit() {
	d.logger.DebugWithFields(d.title, func(oe onelog.Entry) {
		for k, v := range d.m {
			oe.String(k, v)
		}
		if d.withCaller {
			oe.String("CALLER", generateCallerList(d.callerLevel))
		}
	})
}

func generateCallerList(callerLevel int) string {
	var callers strings.Builder

	for i := startCallerLevel; ; i++ {
		_, file, line, ok := runtime.Caller(i)

		if !ok || i == getCallerLevels(callerLevel) {
			break
		}

		var caller strings.Builder
		caller.WriteString(file)
		caller.WriteString(":")
		caller.WriteString(strconv.Itoa(line))

		callers.WriteString(caller.String())
		callers.WriteString(" ")
	}

	return callers.String()
}
