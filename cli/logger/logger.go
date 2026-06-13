package logger

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

func Init() {
	err := godotenv.Load("../.env")
	logPath := os.Getenv("LOG_PATH")

	if err != nil {
		fmt.Println("Unable to load env", err)
	}
	deployeType := os.Getenv("DEPLOYED_TYPE")
	if logPath == "" {
		logPath = "./log/app.log"
	}
	if deployeType == "VM" {
		// Create the lumberjack logger for file rotation
		logFile := &lumberjack.Logger{
			Filename:   logPath, // Log file location
			MaxSize:    10,      // Max size in MB before rotating
			MaxBackups: 3,       // Max number of backup files to keep
			MaxAge:     28,      // Max age in days to keep old log files
			Compress:   true,    // Compress old log files
		}
		// Create a Zap encoder (JSON format)
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

		// Create a Zap core using lumberjack for log rotation
		writeSyncer := zapcore.AddSync(logFile)                          // Use lumberjack's writer syncer
		core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel) // Info level logging

		// Create the final Zap logger with rotation
		Log = zap.New(core)
		fmt.Println("Logging to file with rotation (VM)")

	} else {
		// Default logging to stdout for containerized deployment
		Log = zap.Must(zap.NewProduction())
		//fmt.Println("Logging to terminal (Container)")
	}
}
