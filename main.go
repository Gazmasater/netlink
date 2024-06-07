package main

import (
	"fmt"
	"os"

	"github.com/mdlayher/netlink"
	"test.com/data"
	"test.com/util"
	//		  ^^^^ В го не принято называть пакеты по типу util, tools, data и т.д. - это все общие имена не отражающие сути
)

// TODO вся реализация должна быть скрыта во внутренних слоях. В main только инициализация
func main() {

	// Подключение к Netlink
	conn, err := netlink.Dial(data.NETLINK_NETFILTER, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
		//	^^^^^^^ - нужно использовать пакет для логирования типа log, logrus или "go.uber.org/zap"
		os.Exit(1)
		// ^^^^^^ Логирование и завершение можно совместить в одну функцию log.Fatal которая есть в любом пакете для логирования
	}
	defer conn.Close()

	fmt.Println("Слушаем Netlink сообщения...") //TODO кажется этот лог надо перенести уже после подключения к группе

	// Присоединение к группе Netlink для отслеживания трассировок пакетов
	if err := conn.JoinGroup(data.NFNLGRP_NFTRACE); err != nil {
		//					  ^^^^^^^^^^^^^^^^^^^ эти константы уже определены в пакете "golang.org/x/sys/unix"
		fmt.Fprintf(os.Stderr, "Ошибка подписки на группу: %v\n", err)
		os.Exit(1)
	}
	//TODO не предусмотрен выход из бесконечного цикла. Обычно это делается через context.Context
	/*
			ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		 	defer cancel()
			stoped := make(chan struct{})
			go func() {
				defer close(stoped)
		  		for {
		   			select {
		   			case <-ctx.Done():
		    			return
					case msgs, ok := <-RcvNetlinkMsgs:
						if !ok {
							log.Fatal("netlink channel has already closed")
						}
						for _, msg := range msgs {
							msg.Decode()
							log.Infof("msg: %s", msg)
						}
		   			}
		 	 	}
		 	}()
		 	<-stoped
	*/
	// Бесконечный цикл для приема и обработки сообщений
	for {
		// Получение сообщений от Netlink
		msgs, err := conn.Receive()
		//				^^^^ это блокирующий вызов (нужно это учитывать в реализации)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка получения сообщения: %v\n", err)
			continue
			//^^^^^ почему? Если получили ошибку то какой смысл продолжать? будем только накапливать лог с ошибками!
		}

		// Обработка каждого полученного сообщения
		for _, msg := range msgs {
			// Проверка, что тип сообщения соответствует требуемому и длина данных достаточна
			if len(msg.Data) >= 96 {
				//				^^ magic number
				// Вызов функции ParseMessage для обработки сообщения
				util.ParseMessage(msg)
				//    ^^^^^^^^^^ если посмотришь библиотеки связанные с нетлинком там этот метод называется Decode. Лучше придерживаться уже общепринятых наименований
			}
		}
	}
}
