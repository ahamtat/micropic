package logger

import (
	"log"
	"sync"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger
var once sync.Once

// Init initializes a thread-safe singleton logger
// This would be called from a main method when the application starts up
// This function would ideally, take zap configuration, but is left out
// in favor of simplicity using the example logger.
func Init(logLevel, filePath string) {
	// once ensures the singleton is initialized only once
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stderr"}
		if len(filePath) > 0 {
			config.OutputPaths = append(config.OutputPaths, filePath)
		}
		config.ErrorOutputPaths = []string{"stderr"}
		if len(filePath) > 0 {
			config.ErrorOutputPaths = append(config.ErrorOutputPaths, filePath)
		}
		config.EncoderConfig = zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			//CallerKey:    "caller",
			//EncodeCaller: zapcore.ShortCallerEncoder,
		}

		var level zapcore.Level
		err := level.UnmarshalText([]byte(logLevel))
		if err != nil {
			log.Fatalf("can't marshal level string: %v", logLevel)
		}
		config.Level = zap.NewAtomicLevelAt(level)

		logger, err := config.Build()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
		defer func() {
			if err := logger.Sync(); err != nil {
				return
				// Do not process sync err as told in:
				// https://github.com/uber-go/zap/issues/328
				//log.Fatalf("can't sync logger: %v", err)
			}
		}()
		sugar = logger.Sugar()
	})
}

// Get ordinary logger
func GetLogger() *zap.Logger {
	if sugar == nil {
		return nil
	}
	return sugar.Desugar()
}

// Debug logs a debug message with the given fields
func Debug(message string, fields ...interface{}) {
	sugar.Debugw(message, fields...)
}

// Info logs a debug message with the given fields
func Info(message string, fields ...interface{}) {
	sugar.Infow(message, fields...)
}

// Warn logs a debug message with the given fields
func Warn(message string, fields ...interface{}) {
	sugar.Warnw(message, fields...)
}

// Error logs a debug message with the given fields
func Error(message string, fields ...interface{}) {
	sugar.Errorw(message, fields...)
}

// Fatal logs a message than calls os.Exit(1)
func Fatal(message string, fields ...interface{}) {
	sugar.Fatalw(message, fields...)
}
