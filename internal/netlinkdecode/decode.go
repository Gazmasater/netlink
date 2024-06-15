// TODO перенести этот файл в пакет netlink (см TODO)
package netlinkdecode

import (
	"encoding/binary"
	"net"

	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	// Transport layer header length
	TlHeaderLen = 8
	// Network layer header length
	NlHeaderLen = 20
	// Offset attribute data in the nft netlink group message
	NlNftAttrOffset = 4
)

type (
	proto      uint8
	PacketInfo struct {
		SrcIP    string
		DstIP    string
		SrcPort  uint16
		DstPort  uint16
		Protocol proto
	}
)

func (p proto) String() string {
	switch p {
	case unix.IPPROTO_TCP:
		return "tcp"

	case unix.IPPROTO_UDP:
		return "udp"

	case unix.IPPROTO_UDPLITE:
		return "udplite"

	case unix.IPPROTO_ESP:
		return "esp"

	case unix.IPPROTO_AH:
		return "ah"

	case unix.IPPROTO_ICMP:
		return "icmp"

	case unix.IPPROTO_ICMPV6:
		return "icmpv6"

	case unix.IPPROTO_COMP:
		return "comp"

	case unix.IPPROTO_DCCP:
		return "dccp"

	case unix.IPPROTO_SCTP:
		return "sctp"

	}

	return "unknown"
}

func (pkt *PacketInfo) Decode(b []byte) error {
	//TODO это все надо проверить что корректно работает!
	ad, err := netlink.NewAttributeDecoder(b[NlNftAttrOffset:])
	if err != nil {
		return errors.WithMessage(err, "failed to create new nl attribute decoder")
	}
	ad.ByteOrder = binary.BigEndian
	for ad.Next() {
		switch ad.Type() {
		case unix.NFTA_TRACE_NETWORK_HEADER:
			b := ad.Bytes()
			l := len(b)
			if l < NlHeaderLen {
				return errors.Errorf("incorrect NlHeader binary length=%d", l)
			}
			srcIP := make(net.IP, net.IPv4len)
			dstIP := make(net.IP, net.IPv4len)

			copy(srcIP, b[12:16])
			copy(dstIP, b[16:20])

			pkt.SrcIP = srcIP.String()
			pkt.DstIP = dstIP.String()

			pkt.Protocol = proto(b[9])

		case unix.NFTA_TRACE_TRANSPORT_HEADER:
			b := ad.Bytes()
			if l := len(b); l < TlHeaderLen {
				return errors.Errorf("incorrect TlHeader binary length=%d", l)
			}

			pkt.SrcPort = binary.BigEndian.Uint16(b[:2])
			pkt.DstPort = binary.BigEndian.Uint16(b[2:4])
		}

	}
	if ad.Err() != nil {
		return errors.WithMessage(err, "failed to unmarshal attribute")
	}

	return nil

	/*
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
	*/
}
