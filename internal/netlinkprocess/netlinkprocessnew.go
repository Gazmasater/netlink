package netlinkprocess

import (
	"context"
	"sync"

	"github.com/Gazmasater/netlink/internal/netlinkdecode"
	"github.com/Gazmasater/netlink/pkg/printtcpudp"
	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"

	"github.com/Gazmasater/netlink/pkg/logger"
)

type Collector struct {
	onceRun        sync.Once
	onceClose      sync.Once
	stop           chan struct{}
	stopped        chan struct{}
	RcvNetlinkMsgs chan []netlink.Message
}

func (c *Collector) Run(ctx context.Context) error {
	var doRun bool

	c.onceRun.Do(func() {
		doRun = true
		c.stopped = make(chan struct{})
		c.RcvNetlinkMsgs = make(chan []netlink.Message)
	})
	if !doRun {
		return errors.New("it has been run or closed yet")
	}
	log := logger.FromContext(ctx).Named("collector")
	log.Info("start")
	defer func() {
		log.Info("stop")
		close(c.stopped)
	}()

	conn, err := netlink.Dial(unix.NETLINK_NETFILTER, nil)
	if err != nil {
		return errors.WithMessage(err, "Ошибка подключения")
	}

	// Присоединение к группе Netlink для отслеживания трассировок пакетов
	if err := conn.JoinGroup(unix.NFNLGRP_NFTRACE); err != nil {
		conn.Close()
		return errors.WithMessage(err, "Ошибка подписки на группу")
	}

	go func() {
		for {
			msgs, err := conn.Receive()
			if err != nil {
				log.Errorf("Ошибка приема сообщений от netlink: %v", err)
				return
			}

			for _, msg := range msgs {
				if len(msg.Data) == 160 {
					c.RcvNetlinkMsgs <- msgs

				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Info("will exit cause ctx canceled")
			return ctx.Err()
		case <-c.stop:
			log.Info("will exit cause it has closed")
			return nil
		case msgs, ok := <-c.RcvNetlinkMsgs:
			if !ok {
				log.Info("netlink channel has already closed")
				return errors.New("netlink channel has already closed")
			}
			for _, msg := range msgs {
				var pktInfo netlinkdecode.PacketInfo

				pktInfo, _ = netlinkdecode.Decode(msg)
				printtcpudp.PrintPacketInfo(pktInfo)
			}
		}
	}
}
func (c *Collector) Close() error {
	c.onceClose.Do(func() {
		close(c.stop)
		c.onceRun.Do(func() {})
		if c.stopped != nil {
			<-c.stopped
		}
	})
	return nil
}
