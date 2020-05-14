package logging

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"os"
	"time"
)

var (
	// This timestamp format differs from RFC3339Nano by using .000 instead
	// of .999999999 which changes the timestamp from 9 variable to 3 fixed
	// decimals (.130 instead of .130987456).
	timestampFormat = log.TimestampFormat(
		func() time.Time { return time.Now().UTC() },
		"2006-01-02T15:04:05.000Z07:00",
	)
)

// AllowedLevel is a settable identifier for the minimum level a log entry
// must be have.
type AllowedLevel struct {
	s string
	o level.Option
}

func (l *AllowedLevel) String() string {
	return l.s
}

// Set updates the value of the allowed level.
func (l *AllowedLevel) Set(s string) error {
	switch s {
	case "debug":
		l.o = level.AllowDebug()
	case "info":
		l.o = level.AllowInfo()
	case "warn":
		l.o = level.AllowWarn()
	case "error":
		l.o = level.AllowError()
	default:
		return errors.Errorf("unrecognized log level %q", s)
	}
	l.s = s
	return nil
}

func (l *AllowedLevel) Type() string {
	return "string"
}

// AllowedFormat is a settable identifier for the output format that the logger can have.
type AllowedFormat struct {
	s string
}

func (f *AllowedFormat) String() string {
	return f.s
}

// Set updates the value of the allowed format.
func (f *AllowedFormat) Set(s string) error {
	switch s {
	case "logfmt", "json":
		f.s = s
	default:
		return errors.Errorf("unrecognized log format %q", s)
	}
	return nil
}

func (f *AllowedFormat) Type() string {
	return "string"
}

// Config is a struct containing configurable settings for the logger
type Config struct {
	Level  *AllowedLevel
	Format *AllowedFormat
}

func New(config *Config) log.Logger {
	var l log.Logger
	if config.Format.s == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}

	l = level.NewFilter(l, config.Level.o)
	l = log.With(l, "ts", timestampFormat, "caller", log.DefaultCaller)
	return l
}
