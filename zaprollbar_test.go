package zaprollbar

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewCore(t *testing.T) {
	core := NewCore("token", WithEnvironment("test"))
	if core == nil {
		t.Fatal("expected core to be non-nil")
	}

	if !core.Enabled(zapcore.WarnLevel) {
		t.Error("expected WarnLevel to be enabled by default")
	}

	if core.Enabled(zapcore.InfoLevel) {
		t.Error("expected InfoLevel to be disabled by default")
	}
}

func TestWith(t *testing.T) {
	core := NewCore("token")
	newCore := core.With([]zapcore.Field{
		zap.String("foo", "bar"),
	})

	if newCore == nil {
		t.Fatal("expected newCore to be non-nil")
	}

	if newCore == core {
		t.Error("expected newCore to be a clone, not the same instance")
	}
}

func TestEnabled(t *testing.T) {
	core := NewCore("token", WithMinLevel(zapcore.ErrorLevel))

	if core.Enabled(zapcore.WarnLevel) {
		t.Error("expected WarnLevel to be disabled")
	}

	if !core.Enabled(zapcore.ErrorLevel) {
		t.Error("expected ErrorLevel to be enabled")
	}

	if !core.Enabled(zapcore.FatalLevel) {
		t.Error("expected FatalLevel to be enabled")
	}
}
