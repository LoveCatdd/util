package log

import "github.com/LoveCatdd/util/pkg/lib/core/config"

// ZapConfig
// @Description: zap日志配置结构体
type ZapConfig struct {
	Zap struct {
		Prefix     string         `mapstructure:"prefix"`
		TimeFormat string         `mapstructure:"timeFormat"`
		Level      string         `mapstructure:"level"`
		Caller     bool           `mapstructure:"caller"`
		StackTrace bool           `mapstructure:"stackTrace"`
		Writer     string         `mapstructure:"writer"` //日志输出到哪里 file | console | both
		Encode     string         `mapstructure:"encode"`
		LogFile    *LogFileConfig `mapstructure:"logFile"`
		Enable     bool           `mapstructure:"logFile"` // 开关zap日志
	} `mapstructure:"zap"`
}

// LogFileConfig
// @Description: 日志文件配置结构体
type LogFileConfig struct {
	MaxSize  int      `mapstructure:"maxSize"`
	BackUps  int      `mapstructure:"backups"`
	Compress bool     `mapstructure:"compress"`
	Output   []string `mapstructure:"output"`
	Errput   []string `mapstructure:"errput"`
}

func (*ZapConfig) FileType() string {
	return config.VIPER_YAML
}
