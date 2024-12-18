package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	Log              = "exporter.log"
	LogMaxSize int64 = 1 // Megabytes
	LogFile, _       = os.OpenFile(Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Logging          = zerolog.New(LogFile).Level(zerolog.InfoLevel)

	DebugLog     = "exporter_debug.log"
	DebugLogFile *os.File
	Debugging    zerolog.Logger
	Debug        bool
)

func StdoutLogging(module string, logmsg string) {
	Logging.Info().Time("dt", time.Now()).Str("module", module).Msg(logmsg)
}

func StderrLogging(module string, logmsg string, err error) {
	Logging.Error().Time("dt", time.Now()).Str("module", module).Err(err).Msg(logmsg)
}

func FatalLogging(module string, logmsg string, err error) {
	Logging.Fatal().Time("dt", time.Now()).Str("module", module).Err(err).Msg(logmsg)
}

func DebugLogging(module string, logmsg string) {
	Debugging.Debug().Time("dt", time.Now()).Str("module", module).Msg(logmsg)
}
