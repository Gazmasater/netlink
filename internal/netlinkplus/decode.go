// TODO перенести этот файл в пакет netlink (см TODO)
package netlinkplus

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/Gazmasater/netlink/pkg/logger"
	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

const (
	NFTNL_TRACE_NETWORK_HEADER = 1 << iota
	NFTNL_TRACE_TRANSPORT_HEADER

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
		Flag     uint16
		logger   logger.TypeOfLogger
	}
)

func (p *PacketInfo) SetLogger(log logger.TypeOfLogger) {
	p.logger = log
}

func (pkt *PacketInfo) IsReady() bool {
	requiredFlags := uint16(NFTNL_TRACE_NETWORK_HEADER | NFTNL_TRACE_TRANSPORT_HEADER)
	return pkt.Flag&requiredFlags == requiredFlags
}

func (pkt *PacketInfo) LogPacketInfo() {
	var protocolName string
	switch pkt.Protocol {
	case 6:
		protocolName = "TCP"
	case 17:
		protocolName = "UDP"
	default:
		protocolName = "Unknown"
	}
	pkt.logger.Infof("Packet Information: SrcIP=%s, DstIP=%s, SrcPort=%d, DstPort=%d, Protocol=%s",
		pkt.SrcIP, pkt.DstIP, pkt.SrcPort, pkt.DstPort, protocolName)
}

func (pkt *PacketInfo) LogPacketFile() {
	var protocolName string
	switch pkt.Protocol {
	case 6:
		protocolName = "TCP"
	case 17:
		protocolName = "UDP"
	default:
		protocolName = "Unknown"
	}

	// Open the file for writing (append mode), create it if it doesn't exist
	file, err := os.OpenFile("packet_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the packet information to the file
	logMessage := fmt.Sprintf("Packet Information: SrcIP=%s, DstIP=%s, SrcPort=%d, DstPort=%d, Protocol=%s\n",
		pkt.SrcIP, pkt.DstIP, pkt.SrcPort, pkt.DstPort, protocolName)

	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

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
			pkt.Flag |= NFTNL_TRACE_NETWORK_HEADER
		case unix.NFTA_TRACE_TRANSPORT_HEADER:
			b := ad.Bytes()
			if l := len(b); l < TlHeaderLen {
				return errors.Errorf("incorrect TlHeader binary length=%d", l)
			}
			pkt.SrcPort = binary.BigEndian.Uint16(b[:2])
			pkt.DstPort = binary.BigEndian.Uint16(b[2:4])
			pkt.Flag |= NFTNL_TRACE_TRANSPORT_HEADER
		}
	}
	if ad.Err() != nil {
		return errors.WithMessage(err, "failed to unmarshal attribute")
	}
	return nil
}
