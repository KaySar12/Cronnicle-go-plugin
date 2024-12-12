package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type MyHandler struct{}

func (h MyHandler) Enabled(context context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelDebug:
		return false
	case slog.LevelInfo:
		fallthrough
	case slog.LevelWarn:
		fallthrough
	case slog.LevelError:
		return true
	default:
		panic("unreachable")
	}
}

func (h MyHandler) Handle(context context.Context, record slog.Record) error {
	message := record.Message

	//appends each attribute to the message
	//An attribute is of the form `<key>=<value>` and specified as in `slog.Error(<message>, <key>, <value>, ...)`.
	record.Attrs(func(attr slog.Attr) bool {
		message += fmt.Sprintf(" %v", attr)
		return true
	})

	timestamp := record.Time.Format(time.RFC3339)

	switch record.Level {
	case slog.LevelDebug:
		fallthrough
	case slog.LevelInfo:
		fallthrough
	case slog.LevelWarn:
		fmt.Fprintf(os.Stderr, "[%v] %v %v\n", record.Level, timestamp, message)
	case slog.LevelError:
		fmt.Fprintf(os.Stderr, "!!!ERROR!!! %v %v\n", timestamp, message)
	default:
		panic("unreachable")
	}

	return nil
}

// for advanced users
func (h MyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("unimplemented")
}

// for advanced users
func (h MyHandler) WithGroup(name string) slog.Handler {
	panic("unimplemented")
}
