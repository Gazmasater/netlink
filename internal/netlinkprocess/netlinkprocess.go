package netlinkprocess

import (
	"context"

	"github.com/Gazmasater/netlink/internal/netlinkparser"
	"github.com/Gazmasater/netlink/pkg/logger"
	"github.com/Gazmasater/netlink/pkg/printtcpudp"

	"github.com/mdlayher/netlink"
)

// TODO Обычно, когда речь идет о процессах, их обычно описывают как минимум 2мя методами:
// Run(ctx) error и Close() error (См ниже)
func ProcessNetlinkMessages(ctx context.Context, conn *netlink.Conn) {
	//                                           ^^^^^^^^^^^^^^^^^ netlink.Conn нужно создавать и закрывать внутри метода, нет смысла тянуть его извне
	log := logger.FromContext(ctx).Named("collector")
	for {
		select {
		case <-ctx.Done():
			log.Info("Завершение работы по сигналу")
			return
		default: //TODO вот смешивать case и default в данном случае не оч хорошо! Лучше реализовать прием сообщений через канал,
			// т.е. как писал в прошлый раз здесь должно быть:
			/*
				case msgs, ok := <-RcvNetlinkMsgs:
				if !ok {
					log.Fatal("netlink channel has already closed")
				}
				for _, msg := range msgs {
					msg.Decode()
					log.Infof("msg: %s", msg)
				}
				и никаких default
			*/
			//Вот этот метод  conn.Receive() он вызывает блокирующий системный вызов пока не появятся данные на прием,
			// соответственно в этом месте весь твой цикл будет стопиться навсегда пока не придут данные если делать
			// как ты сделал через default, соответственно не отработается и ctx.Done().
			// Поэтому как я написал выше нужно обернуть твой метод приема данных через канал
			msgs, err := conn.Receive()
			if err != nil {
				log.Errorf("Ошибка получения сообщения: %v", err) //TODO Вообще наверное весь текст ошибок лучше писать на английском, потому что иначе будет смесь языков, потому что внутренние ошибки будут выводиться на английском
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
					//           ^^^^^^^^^ Вот эту функцию нет смысла выносить в pkg.
					// Лучше внутри нашего пакета создать интерфейс Printer,
					// в котором в зависимости от реализации он сможет либо печатать это в консоль,
					// либо в файл, либо еще куда-то. Так и тестить будет проще
				}
			}
		}
	}
}

//TODO Вот как примерно нужно реализовать:
/*
type Collector struct {
	onceRun   sync.Once
	onceClose sync.Once
	stop      chan struct{}
	stopped   chan struct{}
}

func (c *Collector) Run(ctx context.Context) error {
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
		close(r.stopped)
	}()
	for {
		select {
		case <-ctx.Done():
			log.Info("will exit cause ctx canceled")
			return ctx.Err()
		case <-c.stop:
			log.Info("will exit cause it has closed")
			return nil
		case msgs, ok := <-RcvNetlinkMsgs:
			if !ok {
				log.Info("netlink channel has already closed")
				return errors.New("netlink channel has already closed")
			}
			for _, msg := range msgs {
				var pktInfo PacketInfo
				pktInfo.Decode(msg)
				netlink.Printer(pktInfo)
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

*/
