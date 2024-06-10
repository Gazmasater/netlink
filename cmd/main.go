package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Gazmasater/netlink/internal/netlinkconn"
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

	conn, err := netlinkconn.ConnectToNetlink() //TODO Соединение нужно создавать в том месте где оно планируется использоваться. Нет смысла это выносить наружу (см TODO)
	if err != nil {
		logger.Fatalf(ctx, "Ошибка установления соединения: %v", err)
	}
	netlinkprocess.ProcessNetlinkMessages(ctx, conn)

	logger.SetLevel(zap.InfoLevel)
	logger.Info(ctx, "-= BYE =-")
}
