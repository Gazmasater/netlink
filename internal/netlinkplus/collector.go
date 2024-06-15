package netlinkplus

//       ^^^^^^^^^^^^ слишком длинное и сложное название для пакета - переименовать файл и пакет в collector

import (
	"context"
	"sync"

	"github.com/Gazmasater/netlink/pkg/logger"
	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type (
	Collector interface {
		Run(context.Context) error
		Close() error
	}
	collectorImpl struct {
		onceRun   sync.Once
		onceClose sync.Once
		stop      chan struct{}
		stopped   chan struct{}
	}
)

func NewCollector() Collector {
	return &collectorImpl{
		stop: make(chan struct{}),
	}
}

func (c *collectorImpl) Run(ctx context.Context) (err error) {
	var doRun bool

	c.onceRun.Do(func() {
		doRun = true
		c.stopped = make(chan struct{})
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
		//								^^^^^^^^^^^^^^^^^^ ошибки перевести на английский иначе будет мешанина в языках ведь err будет уже содержать ошибку на инглише
	}

	defer conn.Close()

	// Присоединение к группе Netlink для отслеживания трассировок пакетов
	if err = conn.JoinGroup(unix.NFNLGRP_NFTRACE); err != nil {
		return errors.WithMessage(err, "Ошибка подписки на группу")
		//								^^^^^^^^^^^^^^^^^^ ошибки перевести на английский иначе будет мешанина в языках ведь err будет уже содержать ошибку на инглише
	}

	incoming := make(chan any, 1)

	go func() {
		defer close(incoming) // Закрываем канал при завершении горутины
		var e error
		var v any
		for e == nil {
			if v, e = conn.Receive(); e != nil {
				v = e
			}
			select {
			case <-ctx.Done():
				return
			case <-c.stop:
				return
			case incoming <- v:
			}
		}
	}()

loop:
	for err == nil {
		select {
		case v, ok := <-incoming:
			if !ok {
				break loop
			}
			switch t := v.(type) {
			case error:
				err = t
			case []netlink.Message:
				for _, msg := range t {
					var pktInfo PacketInfo
					flag, err := pktInfo.Decode(msg.Data)
					if err != nil {
						break
					}
					if !flag {
						printer := NewPrinter()
						printer.PrintHeader("Packet Information")
						printer.PrintPacketInfo(pktInfo)
					}

				}
			}
		case <-ctx.Done():
			log.Info("will exit cause ctx canceled")
			err = ctx.Err()
		case <-c.stop:
			log.Info("will exit cause it has closed")
			break loop
		}
	}
	//TODO это все надо проверить что корректно работает!
	return err
}

func (c *collectorImpl) Close() error {
	c.onceClose.Do(func() {
		close(c.stop)
		c.onceRun.Do(func() {}) // Сбрасываем onceRun
		if c.stopped != nil {
			<-c.stopped
		}
	})
	return nil
}
