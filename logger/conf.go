package logger

import (
	"os"
	"strconv"
)

const (
	EnvLogFileName    = "LOG_FILE_NAME"
	EnvLogFileMaxSize = "LOG_FILE_MAX_SIZE"
	EnvLogFileMaxAge  = "LOG_FILE_MAX_AGE"
	EnvLogCompress    = "LOG_COMPRESS"
	EnableAtomicLevel = "ENABLE_ATOMIC_LEVEL"
)

type Config struct {
	Filename          string
	MaxSize           int
	MaxAge            int
	Compress          bool
	EnableAtomicLevel bool
}

var (
	conf = Config{
		Filename:          os.ExpandEnv("$HOME/irishub-sync/sync_server.log"),
		MaxSize:           20,
		MaxAge:            7,
		Compress:          true,
		EnableAtomicLevel: true,
	}
)

func init() {
	fileName, found := os.LookupEnv(EnvLogFileName)
	if found {
		conf.Filename = fileName
	}

	maxSize, found := os.LookupEnv(EnvLogFileMaxSize)
	if found {
		if size, err := strconv.Atoi(maxSize); err == nil {
			conf.MaxSize = size
		}
	}

	maxAge, found := os.LookupEnv(EnvLogFileMaxAge)
	if found {
		if age, err := strconv.Atoi(maxAge); err == nil {
			conf.MaxAge = age
		}
	}

	compress, found := os.LookupEnv(EnvLogCompress)
	if found {
		if compre, err := strconv.ParseBool(compress); err == nil {
			conf.Compress = compre
		}
	}

	enableAtomicLevel, found := os.LookupEnv(EnableAtomicLevel)
	if found {
		if atomicLevel, err := strconv.ParseBool(enableAtomicLevel); err == nil {
			conf.EnableAtomicLevel = atomicLevel
		}
	}
}
