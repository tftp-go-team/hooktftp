package logger

import (
	"fmt"
	"log/syslog"
)

var (
	instance *syslog.Writer
	err error
)

func Initialize(tag string) error {
	instance, err = syslog.New(syslog.LOG_INFO, tag)
	return err
}

func Alert(format string, a ...interface{}) (err error) {
	return instance.Alert(fmt.Sprintf(format, a...))
}

func Crit(format string, a ...interface{}) (err error) {
	return instance.Crit(fmt.Sprintf(format, a...))
}

func Debug(format string, a ...interface{}) (err error) {
	return instance.Debug(fmt.Sprintf(format, a...))
}

func Emerg(format string, a ...interface{}) (err error) {
	return instance.Emerg(fmt.Sprintf(format, a...))
}

func Err(format string, a ...interface{}) (err error) {
	return instance.Err(fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) (err error) {
	return instance.Info(fmt.Sprintf(format, a...))
}

func Notice(format string, a ...interface{}) (err error) {
	return instance.Notice(fmt.Sprintf(format, a...))
}

func Warning(format string, a ...interface{}) (err error) {
	return instance.Warning(fmt.Sprintf(format, a...))
}

