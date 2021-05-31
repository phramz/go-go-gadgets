package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var _ Factory = &defaultFactory{}

// Fields alias for multiple extra fields
type Fields = map[string]interface{}

// Factory abstraction for log instance configuration
type Factory interface {
	WithVerbose() Factory
	WithFormatText() Factory
	WithFormatJSON() Factory
	WithCallerReporting() Factory
	WithFields(fields Fields) Factory
	WithField(key string, value interface{}) Factory
	Get() Logger
}

// NewFactory returns a new instanced of logger.Factory
func NewFactory() Factory {
	return &defaultFactory{
		verbose:         false,
		formatter:       &logrus.TextFormatter{},
		callerReporting: false,
		fields:          make(map[string]interface{}),
	}
}

type defaultFactory struct {
	verbose         bool
	formatter       logrus.Formatter
	callerReporting bool
	fields          map[string]interface{}
}

func (d *defaultFactory) Get() Logger {
	instance := logrus.New()

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	instance.SetOutput(os.Stdout)
	instance.SetLevel(logrus.InfoLevel)
	instance.SetFormatter(d.formatter)

	if d.verbose {
		// Only logrus the warning severity or above.
		instance.SetLevel(logrus.DebugLevel)
	}

	if d.callerReporting {
		// Send caller info
		instance.SetReportCaller(true)
	}

	return instance.WithFields(d.fields)
}

func (d *defaultFactory) WithFields(fields Fields) Factory {
	f := d.clone()

	for k, v := range fields {
		f = f.WithField(k, v).(*defaultFactory)
	}

	return f
}

func (d *defaultFactory) WithField(key string, value interface{}) Factory {
	f := d.clone()

	d.fields[key] = value

	return f
}

func (d *defaultFactory) WithCallerReporting() Factory {
	f := d.clone()
	f.callerReporting = true

	return f
}

func (d *defaultFactory) WithVerbose() Factory {
	f := d.clone()
	f.verbose = true

	return f
}

func (d *defaultFactory) WithFormatText() Factory {
	f := d.clone()
	f.formatter = &logrus.TextFormatter{}

	return f
}

func (d *defaultFactory) WithFormatJSON() Factory {
	f := d.clone()
	f.formatter = &logrus.JSONFormatter{}

	return f
}

func (d *defaultFactory) clone() *defaultFactory {
	return &defaultFactory{
		verbose:         d.verbose,
		formatter:       d.formatter,
		callerReporting: d.callerReporting,
		fields:          d.fields,
	}
}
