package logging

import (
	"log/slog"
	"time"
)

const (
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelDebug = slog.LevelDebug
)

type (
	Logger         = slog.Logger
	Attr           = slog.Attr
	Level          = slog.Level
	Handler        = slog.Handler
	Value          = slog.Value
	HandlerOptions = slog.HandlerOptions
	LogValuer      = slog.LogValuer
)

var (
	NewTextHandler = slog.NewTextHandler
	NewJSONHandler = slog.NewJSONHandler
	New            = slog.New
	SetDefault     = slog.SetDefault

	StringAttr   = slog.String
	BoolAttr     = slog.Bool
	Float64Attr  = slog.Float64
	AnyAttr      = slog.Any
	DurationAttr = slog.Duration
	IntAttr      = slog.Int
	Int64Attr    = slog.Int64
	Uint64Attr   = slog.Uint64

	GroupValue = slog.GroupValue
	Group      = slog.Group
)

func Float32Attr(key string, val float32) Attr {
	return slog.Float64(key, float64(val))
}

func UInt32Attr(key string, val uint32) Attr {
	return slog.Int(key, int(val))
}

func Int32Attr(key string, val int32) Attr {
	return slog.Int(key, int(val))
}

func TimeAttr(key string, time time.Time) Attr {
	return slog.String(key, time.String())
}

func ErrAttr(err error) Attr {
	return slog.String("error", err.Error())
}
