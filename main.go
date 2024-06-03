package main

import (
	"fmt"
	"os"

	"github.com/mdlayher/netlink"
	"test.com/data"
	"test.com/util"
)

func main() {

	// Подключение к Netlink
	conn, err := netlink.Dial(data.NETLINK_NETFILTER, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Слушаем Netlink сообщения...")

	// Присоединение к группе Netlink для отслеживания трассировок пакетов
	if err := conn.JoinGroup(data.NFNLGRP_NFTRACE); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подписки на группу: %v\n", err)
		os.Exit(1)
	}

	// Бесконечный цикл для приема и обработки сообщений
	for {
		// Получение сообщений от Netlink
		msgs, err := conn.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка получения сообщения: %v\n", err)
			continue
		}

		// Обработка каждого полученного сообщения
		for _, msg := range msgs {
			// Проверка, что тип сообщения соответствует требуемому и длина данных достаточна
			if len(msg.Data) >= 96 {
				// Вызов функции ParseMessage для обработки сообщения
				util.ParseMessage(msg)
			}
		}
	}
}
