package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	_logger *zap.Logger
	Config  = new(ZapConfig)
)

func zapEncoder() zapcore.Encoder {

	// 新建一个配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "Time",
		LevelKey:      "Level",
		NameKey:       "Logger",
		CallerKey:     "Caller",
		MessageKey:    "Message",
		StacktraceKey: "StackTrace",
		LineEnding:    zapcore.DefaultLineEnding,
		FunctionKey:   zapcore.OmitKey,
	}
	// 自定义时间格式
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(Config.Zap.Prefix + t.Format(Config.Zap.TimeFormat) +
			"\ttraceId: " + traceId.String())
	}

	// 日志级别小写
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 秒级时间间隔
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// 简短的调用者输出
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 完整的序列化logger名称
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	// 最终的日志编码 json或者console
	switch Config.Zap.Encode {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func zapWriteSyncer() zapcore.WriteSyncer {
	syncers := make([]zapcore.WriteSyncer, 0, 2)

	// 如果开启了日志控制台输出，就加入控制台书写器
	if Config.Zap.Writer == WRITEBOTH || Config.Zap.Writer == WRITECONSOLE {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	// 如果开启了日志文件存储，就根据文件路径切片加入书写器
	if Config.Zap.Writer == WRITEBOTH || Config.Zap.Writer == WRITEFILE {
		// 添加日志输出器
		for _, path := range Config.Zap.LogFile.Output {
			logger := &lumberjack.Logger{
				Filename:   path,                        //文件路径
				MaxSize:    Config.Zap.LogFile.MaxSize,  //分割文件的大小
				MaxBackups: Config.Zap.LogFile.BackUps,  //备份次数
				Compress:   Config.Zap.LogFile.Compress, // 是否压缩
				LocalTime:  true,                        //使用本地时间
			}
			syncers = append(syncers, zapcore.Lock(zapcore.AddSync(logger)))
		}
	}
	return zap.CombineWriteSyncers(syncers...)
}

func InitZap() {

	// 构建编码器
	encoder := zapEncoder()

	// 构建日志级别
	levelEnabler := zapLevelEnabler(Config)

	// 最后获得Core和 Options
	subCore, options := tee(1, encoder, levelEnabler)

	// 注入logger
	_logger = zap.New(subCore, options...)
}

// 将所有合并
func tee(skip int, encoder zapcore.Encoder, levelEnabler zapcore.LevelEnabler) (core zapcore.Core, options []zap.Option) {
	sink := zapWriteSyncer()
	return zapcore.NewCore(encoder, sink, levelEnabler), buildOptions(skip)
}

// 构建Option
func buildOptions(skip int) (options []zap.Option) {

	if Config.Zap.Caller {
		options = append(options, zap.AddCaller())
	}

	if Config.Zap.StackTrace {
		options = append(options, zap.AddStacktrace(zapcore.WarnLevel))
	}
	options = append(options, zap.AddCallerSkip(skip))
	return
}

type Field = zap.Field

func Debug(a ...any) {

	_logger.Log(zapcore.DebugLevel, fmt.Sprintln(a...))
}

func Info(a ...any) {
	_logger.Log(zapcore.InfoLevel, fmt.Sprintln(a...))
}

func Warn(a ...any) {
	_logger.Log(zapcore.WarnLevel, fmt.Sprintln(a...))
}

func Error(a ...any) {
	_logger.Log(zapcore.ErrorLevel, fmt.Sprintln(a...))
}

func Panic(a ...any) {
	_logger.Log(zapcore.PanicLevel, fmt.Sprintln(a...))
}

func Fatal(a ...any) {
	_logger.Log(zapcore.FatalLevel, fmt.Sprint(a...))
}

func Debugf(format string, a ...any) {

	_logger.Log(zapcore.DebugLevel, fmt.Sprintf(format, a...))
}

func Infof(format string, a ...any) {
	_logger.Log(zapcore.InfoLevel, fmt.Sprintf(format, a...))
}

func Warnf(format string, a ...any) {
	_logger.Log(zapcore.WarnLevel, fmt.Sprintf(format, a...))
}

func Errorf(format string, a ...any) {
	_logger.Log(zapcore.ErrorLevel, fmt.Sprintf(format, a...))
}

func Panicf(format string, a ...any) {
	_logger.Log(zapcore.PanicLevel, fmt.Sprintf(format, a...))
}

func Fatalf(format string, a ...any) {
	_logger.Log(zapcore.FatalLevel, fmt.Sprintf(format, a...))
}

func Sync() error {
	return _logger.Sync()
}

func OutZapLog(skip int) *zap.Logger {
	// 构建编码器
	encoder := zapEncoder()

	// 构建日志级别
	levelEnabler := zapLevelEnabler(Config)

	// 最后获得Core和 Options
	subCore, options := tee(skip, encoder, levelEnabler)

	// 注入logger
	return zap.New(subCore, options...)
}
