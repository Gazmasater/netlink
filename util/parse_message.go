package util

import (
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

// ParseMessage разбирает и выводит информацию из Netlink сообщения
func ParseMessage(msg netlink.Message) {
	// Заголовок IPv4 начинается с 96-го байта
	ipHeader := msg.Data[96:]

	// Проверяем, что длина IP заголовка достаточна
	if len(ipHeader) < 20 {
		fmt.Println("IP заголовок не найден в данных")
		return
	}

	// Извлечение IP адресов и протокола
	srcIP := ipHeader[12:16]
	dstIP := ipHeader[16:20]
	protocol := ipHeader[9]

	// Извлечение портов
	srcPort := binary.BigEndian.Uint16(ipHeader[24:26])
	dstPort := binary.BigEndian.Uint16(ipHeader[26:28])

	fmt.Printf("Источник IP: %d.%d.%d.%d, Пункт назначения IP: %d.%d.%d.%d, Протокол: %d\n",
		srcIP[0], srcIP[1], srcIP[2], srcIP[3], dstIP[0], dstIP[1], dstIP[2], dstIP[3], protocol)
	fmt.Printf("Источник порт (UDP): %d, Пункт назначения порт (UDP): %d\n", srcPort, dstPort)

}
