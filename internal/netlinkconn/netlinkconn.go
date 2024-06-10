package netlinkconn

import (
	"github.com/mdlayher/netlink"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

const (
	NETLINK_NETFILTER = unix.NETLINK_NETFILTER
	NFNLGRP_NFTRACE   = unix.NFNLGRP_NFTRACE
)

func ConnectToNetlink(logger *zap.Logger) (*netlink.Conn, error) {
	// Подключение к Netlink
	conn, err := netlink.Dial(NETLINK_NETFILTER, nil)
	if err != nil {
		logger.Error("Ошибка подключения", zap.Error(err))
		return nil, err
	}

	// Присоединение к группе Netlink для отслеживания трассировок пакетов
	if err := conn.JoinGroup(NFNLGRP_NFTRACE); err != nil {
		logger.Error("Ошибка подписки на группу", zap.Error(err))
		conn.Close()
		return nil, err
	}

	logger.Info("Слушаем Netlink сообщения...")
	return conn, nil
}
