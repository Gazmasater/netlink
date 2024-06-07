package util

import (
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

// TODO лучше использовать общепринятую семантику для таких функций: func Decode(msg netlink.Message) error
/*
	Примерно как должен выглядеть такой метод в твоем случае:
	const NFTA_TRACE_MAX = unix.NFTA_TRACE_PAD
	const (
		TlHeaderLen = 8
		NlHeaderLen = 20
	)

	type Trace struct {
		SAddr    net.IP
		DAddr    net.IP
		SPort    uint16
		DPort    uint16
	}
	func (t *Trace) Decode(msg netlink.Message) error {
		attrs, err := netlink.AttrDecode(msg[4:], NFTA_TRACE_MAX, nil)
		if err != nil {
			return err
		}
		if attrs[unix.NFTA_TRACE_NETWORK_HEADER] != nil {
			netHeader := attrs[unix.NFTA_TRACE_NETWORK_HEADER].Data
			l := len(netHeader)
			if l < NlHeaderLen {
				return errors.Errorf("incorrect NlHeader binary length=%d", l)
			}
			t.SAddr = make(net.IP, net.IPv4len)
			t.DAddr = make(net.IP, net.IPv4len)

			copy(t.SAddr, netHeader[12:16])
			copy(t.DAddr, netHeader[16:20])
		}
		if attrs[unix.NFTA_TRACE_TRANSPORT_HEADER] != nil {
			transportHeader :=attrs[unix.NFTA_TRACE_TRANSPORT_HEADER].Data
			if l := len(transportHeader); l < TlHeaderLen {
				return errors.Errorf("incorrect TlHeader binary length=%d", l)
			}
			t.SPort = binary.BigEndian.Uint16(transportHeader[:2])
			t.DPort = binary.BigEndian.Uint16(transportHeader[2:4])
		}
		return nil
	}
*/
// ParseMessage разбирает и выводит информацию из Netlink сообщения
func ParseMessage(msg netlink.Message) {
	//TODO нет проверки что ты принимаешь из нетлинк именно тот тип сообщения!
	// Заголовок IPv4 начинается с 96-го байта
	ipHeader := msg.Data[96:]
	//					 ^^ magic number. Все такие значения переопредели через константы + здесь будет паника если длина Data < 96

	// Проверяем, что длина IP заголовка достаточна
	if len(ipHeader) < 20 {
		//			   ^^ magic number. Все такие значения переопредели через константы
		fmt.Println("IP заголовок не найден в данных")
		//   ^^^^^^ Вместо Println нужно возвращать ошибку: return errors.New("IP заголовок не найден в данных")
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
	//	^^^^^ Название функции не отражает того что она еще и печатать лог будет. Обычно тут только парсинг, а печать уже в отдельной функции (смотри комментарии в main.go)

}
