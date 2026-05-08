/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package logger

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

var logLevel = new(slog.LevelVar)

func init() {
	logLevel.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{
		Level: logLevel,
		//AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func SetDebugLevel() {
	logLevel.Set(slog.LevelDebug)
}

func SetInfoLevel() {
	logLevel.Set(slog.LevelInfo)
}

func SetWarningLevel() {
	logLevel.Set(slog.LevelWarn)
}

func SetErrorLevel() {
	logLevel.Set(slog.LevelError)
}

func addSource() slog.Attr {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	return slog.Group("source",
		slog.String("file", file),
		slog.Int("line", line),
	)
}

func Debug(s string, a ...any) {
	slog.LogAttrs(nil, slog.LevelDebug, fmt.Sprintf(s, a...), addSource())
}

func Info(s string, a ...any) {
	slog.LogAttrs(nil, slog.LevelInfo, fmt.Sprintf(s, a...), addSource())
}

func Warn(s string, a ...any) {
	slog.LogAttrs(nil, slog.LevelWarn, fmt.Sprintf(s, a...), addSource())
}

func Error(s string, a ...any) {
	slog.LogAttrs(nil, slog.LevelError, fmt.Sprintf(s, a...), addSource())
}
