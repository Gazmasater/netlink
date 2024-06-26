package nftrace

import (
	"context"
	"sync"

	"github.com/Gazmasater/pkg/logger"
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
		return errors.WithMessage(err, "failed to create netlink connection")
	}

	defer conn.Close()

	if err = conn.JoinGroup(unix.NFNLGRP_NFTRACE); err != nil {
		return errors.WithMessage(err, "failed to join to the netlink NFTRACE group")
	}

	incoming := make(chan any, 1)

	go func() {
		defer close(incoming)
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
					var trace Trace
					err := trace.Decode(msg.Data)
					if err != nil {
						break
					}

					if trace.IsReady() {
						log.Info(trace.String())
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
