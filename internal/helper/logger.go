package helper

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/rs/zerolog/log"
)

const (
	LoggerLevelTrace = "LoggerLevelTrace"
	LoggerLevelDebug = "LoggerLevelDebug"
	LoggerLevelInfo  = "LoggerLevelInfo"
	LoggerLevelWarn  = "LoggerLeveWarn"
	LoggerLevelError = "LoggerLevelError"
	LoggerLevelFatal = "LoggerLevelFatal"
	LoggerLevelPanic = "LoggerLevelPanic"
)

func Logger(level, message string, err error) {
	if err == nil && (level == "" || message == "") {
		log.Error().Stack().Err(errors.New("all params log is required")).Send()
	}

	pc, _, line, _ := runtime.Caller(1)
	path := runtime.FuncForPC(pc).Name()

	switch level {
	case LoggerLevelDebug:
		log.Debug().Str("message", message).Send()
	case LoggerLevelInfo:
		log.Info().Str("message", message).Send()
	case LoggerLevelWarn:
		log.Warn().Str("message", message).Send()
	case LoggerLevelError:
		log.Error().Str("path", path).Str("line", fmt.Sprint(line)).Err(err).Send()
	case LoggerLevelFatal:
		log.Fatal().Str("path", path).Str("line", fmt.Sprint(line)).Err(err).Send()
	case LoggerLevelPanic:
		log.Panic().Str("path", path).Str("line", fmt.Sprint(line)).Err(err).Send()
	default:
		log.Error().Stack().Err(errors.New("logger level invalid")).Send()
	}
}
