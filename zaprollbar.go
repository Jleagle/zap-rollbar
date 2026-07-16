package zaprollbar

import (
	"fmt"

	"github.com/rollbar/rollbar-go"
	"go.uber.org/zap/zapcore"
)

type rollbarCore struct {
	client   *rollbar.Client
	minLevel zapcore.Level
	fields   map[string]interface{}
}

// Option defines a functional option for NewCore.
type Option func(*rollbarCore)

// WithMinLevel sets the minimum log level to send to Rollbar.
// Defaults to WarnLevel.
func WithMinLevel(level zapcore.Level) Option {
	return func(c *rollbarCore) {
		c.minLevel = level
	}
}

// NewCore creates a new zapcore.Core that sends logs to Rollbar using the provided rollbar.Client.
func NewCore(client *rollbar.Client, opts ...Option) zapcore.Core {
	core := &rollbarCore{
		client:   client,
		minLevel: zapcore.WarnLevel,
		fields:   make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(core)
	}

	return core
}

func (c *rollbarCore) Enabled(level zapcore.Level) bool {
	return level >= c.minLevel
}

func (c *rollbarCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}
	for k, v := range enc.Fields {
		clone.fields[k] = v
	}
	return clone
}

func (c *rollbarCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

func (c *rollbarCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	extras := make(map[string]interface{}, len(c.fields)+len(fields))
	for k, v := range c.fields {
		extras[k] = v
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}
	for k, v := range enc.Fields {
		extras[k] = v
	}

	level := rollbar.ERR
	switch entry.Level {
	case zapcore.DebugLevel:
		level = rollbar.DEBUG
	case zapcore.InfoLevel:
		level = rollbar.INFO
	case zapcore.WarnLevel:
		level = rollbar.WARN
	case zapcore.ErrorLevel:
		level = rollbar.ERR
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		level = rollbar.CRIT
	}

	if entry.Level >= zapcore.ErrorLevel {
		c.client.ErrorWithExtras(level, fmt.Errorf(entry.Message), extras)
	} else {
		c.client.MessageWithExtras(level, entry.Message, extras)
	}

	return nil
}

func (c *rollbarCore) Sync() error {
	c.client.Wait()
	return nil
}

func (c *rollbarCore) clone() *rollbarCore {
	newFields := make(map[string]interface{}, len(c.fields))
	for k, v := range c.fields {
		newFields[k] = v
	}
	return &rollbarCore{
		client:   c.client,
		minLevel: c.minLevel,
		fields:   newFields,
	}
}
