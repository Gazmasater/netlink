package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Gazmasater/netlink/internal/netlinkprocess"
	"github.com/Gazmasater/netlink/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger.SetLevel(zap.DebugLevel)
	logger.Info(ctx, "-= HELLO =-")

	collector := &netlinkprocess.Collector{}
	if err := collector.Run(ctx); err != nil {
		logger.Error(ctx, "Error running collector", zap.Error(err))
	}
	logger.SetLevel(zap.InfoLevel)
	logger.Info(ctx, "-= BYE =-")
}
