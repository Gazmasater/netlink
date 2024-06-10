package main

import (
	"fmt"
	"log"

	"github.com/Gazmasater/netlink/logger"
	"github.com/Gazmasater/netlink/netlinkconnect"
	"github.com/Gazmasater/netlink/netlinkparser"
	"github.com/Gazmasater/netlink/printttcpudp"

	"go.uber.org/zap"
)

func main() {

	logger, err := logger.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Logger initialized successfully")
	defer logger.Sync()

	conn, err := netlinkconnect.ConnectToNetlink(logger)
	if err != nil {
		logger.Error("Ошибка подключения к Netlink", zap.Error(err))
		return
	}

	// Бесконечный цикл для приема и обработки сообщений
	for {
		// Получение сообщений от Netlink
		msgs, err := conn.Receive()
		if err != nil {
			logger.Error("Ошибка получения сообщения", zap.Error(err))
			continue
		}

		// Обработка каждого полученного сообщения
		for _, msg := range msgs {
			// Проверка, что тип сообщения соответствует требуемому и длина данных достаточна
			if len(msg.Data) >= 96 {
				// Вызов функции ParseMessage для обработки сообщения
				packet, err := netlinkparser.Decode(msg)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// Выводим информацию о пакете
				printttcpudp.PrintPacketInfo(packet)
			}
		}

	}
}
