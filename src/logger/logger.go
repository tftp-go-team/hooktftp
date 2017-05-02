package logger

import (
	"fmt"
	"log/syslog"
	"os"
)

var (
	instance *syslog.Writer
	err      error
)

func Initialize(tag string) error {
	instance, err = syslog.New(syslog.LOG_INFO, tag)
	return err
}

func Close() error {
	return instance.Close()
}

func Alert(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Alert(str)
}

func Crit(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Crit(str)
}

func Debug(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Debug(str)
}

func Emerg(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Emerg(str)
}

func Err(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Err(str)
}

func Info(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Info(str)
}

func Notice(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Notice(str)
}

func Warning(format string, a ...interface{}) (err error) {
	str := fmt.Sprintf(format, a...)
	if instance == nil {
		_, ferr := fmt.Fprintln(os.Stderr, str)
		return ferr
	}
	return instance.Warning(str)
}
