package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Gazmasater/netlink/internal/netlinkconn"
	"github.com/Gazmasater/netlink/internal/netlinkprocess"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := netlinkconn.ConnectToNetlink(logger)
	if err != nil {
		logger.Error("Ошибка установления соединения", zap.Error(err))
		return
	}
	netlinkprocess.ProcessNetlinkMessages(ctx, conn, logger)

}
