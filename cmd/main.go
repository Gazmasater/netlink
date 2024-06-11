package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Gazmasater/netlink/internal/netlinkprocess"

	//                                       ^^^^^^^^^^^^ netlinkconn и netlinkprocess перенести в один пакет (см TODO)
	"github.com/Gazmasater/netlink/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger.SetLevel(zap.DebugLevel)
	logger.Info(ctx, "-= HELLO =-")

	netlinkprocess.ProcessNetlinkMessages(ctx)

	logger.SetLevel(zap.InfoLevel)
	logger.Info(ctx, "-= BYE =-")
}
