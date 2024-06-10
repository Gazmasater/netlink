package netlinkconn

import (
	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	NETLINK_NETFILTER = unix.NETLINK_NETFILTER
	NFNLGRP_NFTRACE   = unix.NFNLGRP_NFTRACE
	//                       ^^^^^^^^^^^^ какой смысл в переопределении этих констант?
)

// TODO Нет смысла выносить создание подключения в отдельный пакет. Соединение нужно создавать там где оно и планируется использоваться, т.е. в твоем случае в netlinkprocess.go (см TODO)
func ConnectToNetlink() (*netlink.Conn, error) {
	// Подключение к Netlink
	conn, err := netlink.Dial(NETLINK_NETFILTER, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "Ошибка подключения")
	}

	// Присоединение к группе Netlink для отслеживания трассировок пакетов
	if err := conn.JoinGroup(NFNLGRP_NFTRACE); err != nil {
		conn.Close()
		return nil, errors.WithMessage(err, "Ошибка подписки на группу")
	}

	return conn, nil
}
