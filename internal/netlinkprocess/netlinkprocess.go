package netlinkprocess

import (
	"context"

	"github.com/Gazmasater/netlink/internal/netlinkparser"
	"github.com/Gazmasater/netlink/pkg/printtcpudp"

	"github.com/mdlayher/netlink"
	"go.uber.org/zap"
)

func ProcessNetlinkMessages(ctx context.Context, conn *netlink.Conn, logger *zap.Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Завершение работы по сигналу")
			return
		default:
			msgs, err := conn.Receive()
			if err != nil {
				logger.Error("Ошибка получения сообщения", zap.Error(err))
				continue
			}

			for _, msg := range msgs {
				if len(msg.Data) >= 96 {
					packet, err := netlinkparser.Decode(msg)
					if err != nil {
						logger.Sugar().Fatal(err)
						return
					}

					// Выводим информацию о пакете
					printtcpudp.PrintPacketInfo(packet)
				}
			}
		}
	}
}
