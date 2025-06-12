package log

// This is just a "log" wrapper

import (
    "fmt"
    "log"
    "os"
)

type Logger interface {
    Info(messages ...string)
    Infof(tmp string, messages ...any)
    Debug(messages ...string)
    Debugf(tmp string, messages ...any)
    Error(messages ...string)
    Errorf(tmp string, messages ...any)
    Fatal(messages ...string)
}

type LogLevel string

const INFO_LEVEL LogLevel = "INFO"
const DEBUG_LEVEL LogLevel = "DEBUG"
const ERROR_LEVEL LogLevel = "ERROR"
const FATAL_LEVEL LogLevel = "FATAL"

type LoggerImpl struct {
    name  string
    level LogLevel
}

var loggerBag map[string]Logger

func GetOrCreateLogger(name string, defaultLevel string) Logger {
    if loggerBag == nil {
        loggerBag = make(map[string]Logger)
    }

    if loggerBag[name] == nil {
        loggerBag[name] = NewLogger(name, defaultLevel)
    }

    return loggerBag[name]
}

func NewLogger(name string, level string) Logger {
    if level == string(DEBUG_LEVEL) {
        return &LoggerImpl{name: name, level: DEBUG_LEVEL}
    } else if level == string(INFO_LEVEL) {
        return &LoggerImpl{name: name, level: INFO_LEVEL}
    } else {
        // other levels cannot be chosen. Default to INFO
        logger := &LoggerImpl{name: name, level: INFO_LEVEL}
        logger.Info("Invalid log level provided. Defaulting to INFO")
        return logger
    }
}

func (l *LoggerImpl) Infof(tmp string, messages ...any) {
    l._logf(INFO_LEVEL, tmp, messages...)
}

func (l *LoggerImpl) Info(messages ...string) {
    // Always log
    l._log(INFO_LEVEL, messages...)
}

func (l *LoggerImpl) Debugf(tmp string, messages ...any) {
    // Only log if debug level is set
    if l.level == DEBUG_LEVEL {
        l._logf(DEBUG_LEVEL, tmp, messages...)
    }
}

func (l *LoggerImpl) Debug(messages ...string) {
    // Only log if debug level is set
    if l.level == DEBUG_LEVEL {
        l._log(DEBUG_LEVEL, messages...)
    }
}

func (l *LoggerImpl) Errorf(tmp string, messages ...any) {
    l._logf(ERROR_LEVEL, tmp, messages...)
}

func (l *LoggerImpl) Error(messages ...string) {
    l._log(ERROR_LEVEL, messages...)
}

func (l *LoggerImpl) Fatal(messages ...string) {
    // Log and terminate
    l._log(INFO_LEVEL, messages...)
    os.Exit(1)
}

func (l *LoggerImpl) _log(logLevel LogLevel, messages ...string) {
    prefix := fmt.Sprintf("[%v] [%v]:", string(logLevel), l.name)

    args := make([]any, len(messages) + 1)
    args[0] = prefix
    for i, msg := range messages {
        args[i + 1] = msg
    }

    log.Println(args...)
}

func (l *LoggerImpl) _logf(logLevel LogLevel, tmp string, messages ...any) {
    finalMessage := fmt.Sprintf(tmp, messages...)
    l._log(logLevel, finalMessage)
}
