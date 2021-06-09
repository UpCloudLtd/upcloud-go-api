package logger

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	DebugEnvVar = "UPCLOUD_SDK_DEBUG"
)

var (
	Logger = logrus.New()
)

func enabledDebug() bool {
	env, ok := os.LookupEnv(DebugEnvVar)
	if !ok {
		return false
	}

	v, err := strconv.ParseBool(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: environment variable %s has invalid bool value", DebugEnvVar)
	}

	return v
}

func SetupLogger() {
	if enabledDebug() {
		Logger.SetLevel(logrus.DebugLevel)
	}

	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote:           true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
	// More config goes here.
}

func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}