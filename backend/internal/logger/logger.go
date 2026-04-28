// Package logger provides structured JSON logging with component tagging.
// Log lines are machine-parseable so monitoring tools can filter by component,
// level, or error type without regex scraping.
package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type entry struct {
	Time      string `json:"time"`
	Level     Level  `json:"level"`
	Component string `json:"component"`
	Message   string `json:"msg"`
	Error     string `json:"error,omitempty"`
}

var out = log.New(os.Stdout, "", 0)

func write(level Level, component, msg, errStr string) {
	e := entry{
		Time:      time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Component: component,
		Message:   msg,
		Error:     errStr,
	}
	b, _ := json.Marshal(e)
	out.Println(string(b))
}

func Info(component, msg string) {
	write(LevelInfo, component, msg, "")
}

func Warn(component, msg string, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	write(LevelWarn, component, msg, errStr)
}

func Error(component, msg string, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	write(LevelError, component, msg, errStr)
}
