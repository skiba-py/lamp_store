package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level, format string) (*zap.Logger, error) {
	var config zap.Config

	// Уровень логирования
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	// Настраиваем конфигурацию в зависимости от формата
	if format == "json" {
		config = zap.Config{
			Level:       zap.NewAtomicLevelAt(zapLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		config = zap.Config{
			Level:       zap.NewAtomicLevelAt(zapLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	return config.Build()
}
