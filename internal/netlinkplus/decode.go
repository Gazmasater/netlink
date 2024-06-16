// TODO перенести этот файл в пакет netlink (см TODO)
package netlinkplus

import (
	"encoding/binary"
	"net"

	"github.com/Gazmasater/netlink/internal/data"

	"github.com/Gazmasater/netlink/pkg/logger"
	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
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
	requiredFlags := uint16(data.NFTNL_TRACE_NETWORK_HEADER | data.NFTNL_TRACE_TRANSPORT_HEADER)
	return pkt.Flag&requiredFlags == requiredFlags
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
	ad, err := netlink.NewAttributeDecoder(b[data.NlNftAttrOffset:])
	if err != nil {
		return errors.WithMessage(err, "failed to create new nl attribute decoder")
	}
	ad.ByteOrder = binary.BigEndian
	for ad.Next() {
		switch ad.Type() {
		case unix.NFTA_TRACE_NETWORK_HEADER:
			b := ad.Bytes()
			l := len(b)
			if l < data.NlHeaderLen {
				return errors.Errorf("incorrect NlHeader binary length=%d", l)
			}
			srcIP := make(net.IP, net.IPv4len)
			dstIP := make(net.IP, net.IPv4len)
			copy(srcIP, b[12:16])
			copy(dstIP, b[16:20])
			pkt.SrcIP = srcIP.String()
			pkt.DstIP = dstIP.String()
			pkt.Protocol = proto(b[9])
			pkt.Flag |= data.NFTNL_TRACE_NETWORK_HEADER
		case unix.NFTA_TRACE_TRANSPORT_HEADER:
			b := ad.Bytes()
			if l := len(b); l < data.TlHeaderLen {
				return errors.Errorf("incorrect TlHeader binary length=%d", l)
			}
			pkt.SrcPort = binary.BigEndian.Uint16(b[:2])
			pkt.DstPort = binary.BigEndian.Uint16(b[2:4])
			pkt.Flag |= data.NFTNL_TRACE_TRANSPORT_HEADER
		}
	}
	if ad.Err() != nil {
		return errors.WithMessage(err, "failed to unmarshal attribute")
	}
	return nil
}
