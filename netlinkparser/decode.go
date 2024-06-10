package netlinkparser

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

const NFTA_TRACE_MAX = unix.NFTA_TRACE_PAD
const (
	TlHeaderLen = 8
	NlHeaderLen = 20
)

// ParseMessage разбирает и выводит информацию из Netlink сообщения
type PacketInfo struct {
	SrcIP    string
	DstIP    string
	SrcPort  string
	DstPort  string
	Protocol string
}

func Decode(msg netlink.Message) (PacketInfo, error) {

	// Заголовок IPv4 начинается с 96-го байта
	ipHeader := msg.Data[96:]

	// Проверяем, что длина IP заголовка достаточна
	if len(ipHeader) < NlHeaderLen {
		return PacketInfo{}, fmt.Errorf("IP заголовок не найден в данных")
	}

	// Извлечение IP адресов и протокола
	srcIP := ipHeader[12:16]
	dstIP := ipHeader[16:20]
	protocol := ipHeader[9]
	var protocol_string string
	switch protocol {
	case 6:
		protocol_string = "TCP"
	case 17:
		protocol_string = "UDP"
	}

	src_IP := net.IP(srcIP).String()
	dst_IP := net.IP(dstIP).String()

	// Извлечение портов
	srcPort := binary.BigEndian.Uint16(ipHeader[24:26])
	dstPort := binary.BigEndian.Uint16(ipHeader[26:28])

	srcPort_str := strconv.FormatUint(uint64(srcPort), 10)
	dstPort_str := strconv.FormatUint(uint64(dstPort), 10)

	return PacketInfo{
		SrcIP:    src_IP,
		DstIP:    dst_IP,
		SrcPort:  srcPort_str,
		DstPort:  dstPort_str,
		Protocol: protocol_string,
	}, nil
}

func PrintPacketInfo(packet PacketInfo) {
	fmt.Printf("Источник IP: %s\n", packet.SrcIP)
	fmt.Printf("Пункт назначения IP: %s\n", packet.DstIP)
	fmt.Printf("Источник порт: %s\n", packet.SrcPort)
	fmt.Printf("Пункт назначения порт: %s\n", packet.DstPort)
	fmt.Printf("Протокол: %s\n", packet.Protocol)
	fmt.Println()
}
