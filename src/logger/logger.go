package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(logName string, silence, caller bool) *zap.Logger {
	var logger *zap.Logger

	fileName := fmt.Sprintf("log/%v", logName)

	//获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	//encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	//encoder := zapcore.NewConsoleEncoder(encoderConfig)
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	//文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName, //日志文件存放目录
		MaxSize:    16,       //文件大小限制,单位MB
		MaxBackups: 15,       //最大保留日志文件数量
		MaxAge:     90,       //日志文件保留天数
		Compress:   true,     //是否压缩处理
	})

	// 是否静默
	var syncer zapcore.WriteSyncer
	switch silence {
	case true:
		syncer = zapcore.NewMultiWriteSyncer(fileWriteSyncer)
	case false:
		syncer = zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout))
	}

	//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	fileCore := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)

	// 是否添加定位
	switch caller {
	case true:
		logger = zap.New(fileCore, zap.AddCaller()) //AddCaller()为显示文件名和行号
	case false:
		logger = zap.New(fileCore)
	}

	return logger
}

func Data(values ...interface{}) []zap.Field {
	log := zap.NewExample()
	if len(values) == 0 || len(values)%2 != 0 {
		log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", values))
		return nil
	}
	// Zap.Field is used when keyvals pairs appear
	var data []zap.Field
	for i := 0; i < len(values); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(values[i]), values[i+1]))
	}
	return data
}
