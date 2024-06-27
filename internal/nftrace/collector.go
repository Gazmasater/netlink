package nftrace

import (
	"context"
	"log"
	"sync"

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

	log.Println("start")
	defer func() {
		log.Println("stop")
		close(c.stopped)
	}()

	// Создание Netlink сокета
	fd, err := unix.Socket(
		unix.AF_NETLINK,
		unix.SOCK_RAW,
		unix.NETLINK_GENERIC,
	)
	if err != nil {
		return errors.WithMessage(err, "failed to create netlink socket")
	}
	defer unix.Close(fd)

	// Привязка сокета к Netlink адресу
	addr := unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: 0,
		Pid:    0,
	}
	if err := unix.Bind(fd, &addr); err != nil {
		return errors.WithMessage(err, "failed to bind socket to address")
	}

	// Установка группы для прослушивания
	if err := unix.SetsockoptInt(fd, unix.SOL_NETLINK, unix.NETLINK_ADD_MEMBERSHIP, unix.NFNLGRP_NFTRACE); err != nil {
		return errors.WithMessage(err, "failed to join NFTRACE group")
	}

	// Канал для приема данных
	incoming := make(chan interface{}, 1)

	go func() {
		defer close(incoming)
		buf := make([]byte, unix.Getpagesize())
		for {
			n, _, err := unix.Recvfrom(fd, buf, 0)
			if err != nil {
				if err == unix.EINTR {
					continue
				}
				incoming <- err
				return
			}
			incoming <- buf[:n]
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
			case []byte:
				// Здесь можно обработать полученные данные
				log.Printf("Received data: %v\n", t)
				// Пример обработки данных, необходима дополнительная логика
			}
		case <-ctx.Done():
			log.Println("will exit cause ctx canceled")
			err = ctx.Err()
		case <-c.stop:
			log.Println("will exit cause it has closed")
			break loop
		}
	}
	return err
}

func (c *collectorImpl) Close() error {
	c.onceClose.Do(func() {
		close(c.stop)
		c.onceRun.Do(func() {})
		if c.stopped != nil {
			<-c.stopped
		}
	})
	return nil
}
