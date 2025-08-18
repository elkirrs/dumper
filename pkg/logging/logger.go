package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type LoggerOptions struct {
	Level      Level
	AddSource  bool
	IsJSON     bool
	SetDefault bool
	Enabled    bool
	File       string
	AlsoStdout bool
	FileHandle *os.File
}

type Logs struct {
	*Logger
	file *os.File
}

var defaultWriter = os.Stdout

type LoggerOption func(*LoggerOptions)

func NewLogger(opts ...LoggerOption) *Logs {
	config := &LoggerOptions{
		Level:      LevelInfo,
		AddSource:  true,
		IsJSON:     true,
		SetDefault: true,
		Enabled:    true,
	}

	for _, opt := range opts {
		opt(config)
	}

	if !config.Enabled {
		logger := New(slog.NewTextHandler(io.Discard, nil))
		return &Logs{Logger: logger}
	}

	var writers []io.Writer

	// Если указан файл
	if config.File != "" {
		f, err := os.OpenFile(config.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			panic("cannot open log file: " + err.Error())
		}
		config.FileHandle = f
		writers = append(writers, f)
	}

	// Если включено дублирование в консоль
	if config.AlsoStdout && config.File != "" {
		writers = append(writers, os.Stdout)
	}

	// Если не указано ни одно место, по умолчанию os.Stdout
	if len(writers) == 0 {
		writers = append(writers, defaultWriter)
	}

	multiWriter := io.MultiWriter(writers...)

	options := &HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	var h Handler = NewTextHandler(multiWriter, options)
	if config.IsJSON {
		h = NewJSONHandler(multiWriter, options)
	}

	logger := New(h)
	if config.SetDefault {
		SetDefault(logger)
	}

	return &Logs{Logger: logger, file: config.FileHandle}
}

func WithLevel(level string) LoggerOption {
	return func(o *LoggerOptions) {
		var l Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			l = LevelInfo
		}

		o.Level = l
	}
}

func WithAddSource(addSource bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AddSource = addSource
	}
}

func WithIsJSON(isJSON bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.IsJSON = isJSON
	}
}

func WithSetDefault(setDefault bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.SetDefault = setDefault
	}
}

func WithFile(file string) LoggerOption {
	return func(o *LoggerOptions) {
		o.File = file
	}
}

func WithAlsoStdout(also bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.AlsoStdout = also
	}
}

func WithEnabled(enabled bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.Enabled = enabled
	}
}

func (l *Logs) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}

func Default() *Logger {
	return slog.Default()
}
