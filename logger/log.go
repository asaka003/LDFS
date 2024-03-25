package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

//自定义logger配置初始化
func InitLog() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()

	var core zapcore.Core
	// mode := viper.GetString("app.mode")
	mode := "dev"
	if mode == "dev" {
		//开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			//输出到日志文件
			zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel),
			//输出到终端
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.DebugLevel),
		)
	} else {
		//这里设置记录的日志信息为Debug等级以上，日志输出到文件中
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	}

	//zap.AddCaller()用于添加调用信息
	Logger = zap.New(core, zap.AddCaller())
}

//使用日志切割
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: "./test.log",
		// MaxSize:    viper.GetInt("log.max_size"),    //M
		// MaxBackups: viper.GetInt("log.max_backups"), //最大备份数量
		// MaxAge:     viper.GetInt("log.max_age"),     //最大保存天数
		MaxSize:    10,    //M
		MaxBackups: 5,     //最大备份数量
		MaxAge:     2,     //最大保存天数
		Compress:   false, //是否进行压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

//以哪种编码形式写入日志，这里测试选用JSON
func getEncoder() zapcore.Encoder {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoder)
}
