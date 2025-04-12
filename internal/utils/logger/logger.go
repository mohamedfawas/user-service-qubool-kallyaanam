package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger with additional methods
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(isDevelopment bool) (*Logger, error) {
	var config zap.Config
	if isDevelopment {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: zapLogger,
	}, nil
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

// UserProfileEvent logs events related to user profiles
func (l *Logger) UserProfileEvent(ctx context.Context, event string, userID string, profileID string, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("event", event),
		zap.String("module", "user_profile"),
	}

	if userID != "" {
		baseFields = append(baseFields, zap.String("user_id", userID))
	}

	if profileID != "" {
		baseFields = append(baseFields, zap.String("profile_id", profileID))
	}

	// Add request ID if available in context
	if requestID, ok := ctx.Value("request_id").(string); ok {
		baseFields = append(baseFields, zap.String("request_id", requestID))
	}

	l.Info("User profile event", append(baseFields, fields...)...)
}
