package netlinkdecode

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/mdlayher/netlink"
)

const (
	NlHeaderLen = 20
)

type PacketInfo struct {
	SrcIP    string
	DstIP    string
	SrcPort  string
	DstPort  string
	Protocol string
}

func Decode(msg netlink.Message) (PacketInfo, error) {
	ipHeader := msg.Data[0:]

	// Проверяем, что длина IP заголовка достаточна
	if len(ipHeader) < NlHeaderLen {
		return PacketInfo{}, fmt.Errorf("IP заголовок не найден в данных")
	}

	// Извлечение IP адресов и протокола
	srcIP := ipHeader[108:112]
	dstIP := ipHeader[112:116]
	protocol := ipHeader[105]
	var protocolString string
	switch protocol {
	case 6:
		protocolString = "TCP"
	case 17:
		protocolString = "UDP"
	}

	srcIPStr := net.IP(srcIP).String()
	dstIPStr := net.IP(dstIP).String()

	// Извлечение портов
	srcPort := binary.BigEndian.Uint16(ipHeader[120:122])
	dstPort := binary.BigEndian.Uint16(ipHeader[122:124])

	srcPortStr := strconv.FormatUint(uint64(srcPort), 10)
	dstPortStr := strconv.FormatUint(uint64(dstPort), 10)

	if protocolString == "TCP" {

		flagsByte1 := ipHeader[133]
		bitString1 := fmt.Sprintf("%08b", flagsByte1)
		fmt.Printf("Все биты 133-го байта: %s\n", bitString1)
	}

	// Формирование структуры PacketInfo
	packetInfo := PacketInfo{
		SrcIP:    srcIPStr,
		DstIP:    dstIPStr,
		SrcPort:  srcPortStr,
		DstPort:  dstPortStr,
		Protocol: protocolString,
	}

	return packetInfo, nil
}
