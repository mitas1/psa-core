package logging

import (
	"os"

	logging "github.com/op/go-logging"
)

const LOGGER_NAME = "psa-backend"

var (
	logFormat = logging.MustStringFormatter(
		`%{color}%{time:2006/02/01-15:04:05.000} %{shortpkg}: %{shortfile} %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	fileLogFormat = logging.MustStringFormatter(
		`%{level:-8s} %{time:15:04:05.99} %{shortpkg:5s}.%{shortfunc:-12s}> %{message}`,
	)
)

func GetLogger() *logging.Logger {
	return logging.MustGetLogger(LOGGER_NAME)
}

func SetupLogger(filename *string) *logging.Logger {
	logger := GetLogger()

	backendStderr := logging.AddModuleLevel(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logFormat,
	))

	logging.SetBackend(backendStderr)

	if filename != nil {
		logFile, err := os.Create(*filename)
		if err != nil {
			os.Exit(1)
		}

		backendFile := logging.AddModuleLevel(logging.NewBackendFormatter(
			logging.NewLogBackend(logFile, "", 0),
			fileLogFormat,
		))

		logging.SetBackend(backendStderr, backendFile)
	}

	return logger
}
