package zap_factory

import (
	"agent/global/variable"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func CreateZapFactory(entry func(zapcore.Entry) error) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()

	timePrecision := variable.Config.Logs.TimePrecision
	var recordTimeFormat string
	switch timePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006-01-02 15:04:05"

	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at" // The time key field for generating JSON format log is ts by default. After modification, it is convenient to import the log to elk server

	var encoder zapcore.Encoder
	switch variable.Config.Logs.TextFormat {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // Normal mode
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // JSON format
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // Normal mode
	}

	//Writer
	fileName := variable.BasePath + variable.Config.Logs.LogName
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,                        //Location of log files
		MaxSize:    variable.Config.Logs.MaxSize,    //The maximum size (in megabytes) of the log file before cutting
		MaxBackups: variable.Config.Logs.MaxBackups, //Maximum number of old files to keep
		MaxAge:     variable.Config.Logs.MaxAge,     //Maximum number of days to keep old files
		Compress:   variable.Config.Logs.Compress,   //Compress / archive old files
	}
	writer := zapcore.AddSync(lumberJackLogger)
	// Start initializing zap log core parameters，
	//Parameter 1: encoder
	//Parameter 2: writer
	//Parameter 3: parameter level and debug level support logging of all subsequent functions. If it is at the high level of fatal, the log can be written only at the level of > = fatal
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry))
}
