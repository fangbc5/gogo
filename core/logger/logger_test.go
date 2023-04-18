package logger

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger(WithLevel(TraceLevel))
	h1 := NewHelper(l).WithFields(map[string]interface{}{"key1": "val1"})
	h1.Trace("trace_msg1")
	h1.Warn("warn_msg1")

	h2 := NewHelper(l).WithFields(map[string]interface{}{"key2": "val2"})
	h2.Trace("trace_msg2")
	h2.Warn("warn_msg2")

	h3 := NewHelper(l)
	h3.Info("test_msg")
	h3.Error("test_error")
	ctx := context.TODO()
	ctx = context.WithValue(ctx, &loggerKey{}, h3)
	v := ctx.Value(&loggerKey{})
	ll := v.(*Helper)
	ll.Info("test_msg")
}
