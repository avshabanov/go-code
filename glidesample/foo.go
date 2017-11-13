package glidesample

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Foo is a sample exported function
func Foo() string {
	return fmt.Sprintf("Hello from Glidesample library version %s", Version)
}

// Foox is a sample exported function
func Foox(a, b int) int {
	return 100 + a + b
}

// FooService is a sample service
type FooService struct {
	Logger *zap.Logger
}

// RunFooService runs foo service
func RunFooService(s *FooService) string {
	logger := s.Logger
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	logger.Info("Something happened!",
		// Structured context as strongly typed Field values.
		zap.String("url", "http://example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	return "DONE"
}
