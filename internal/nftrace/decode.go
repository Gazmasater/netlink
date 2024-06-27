package nftrace

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type (
	proto   uint8
	pktInfo struct {
		SrcIP    string
		DstIP    string
		SrcPort  uint16
		DstPort  uint16
		Protocol proto
		// tcp flags: syn, ack, etc...
		Flags uint8 // 3 bits
	}
	Trace struct {
		Data pktInfo
		Flag uint16
	}
)

func (t *Trace) IsReady() bool {
	return t.Flag&NFTNL_TRACE_NETWORK_HEADER != 0 && t.Flag&NFTNL_TRACE_TRANSPORT_HEADER != 0
}

func (t *Trace) String() string {
	return fmt.Sprintf("Packet Information: SrcIP=%s, DstIP=%s, SrcPort=%d, DstPort=%d, Protocol=%s",
		t.Data.SrcIP, t.Data.DstIP, t.Data.SrcPort, t.Data.DstPort, t.Data.Protocol)
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

func (t *Trace) Decode(b []byte) error {

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
			t.Data.SrcIP = srcIP.String()
			t.Data.DstIP = dstIP.String()
			t.Data.Protocol = proto(b[9])

			t.Flag = NFTNL_TRACE_NETWORK_HEADER
		case unix.NFTA_TRACE_TRANSPORT_HEADER:
			b := ad.Bytes()
			if l := len(b); l < TlHeaderLen {
				return errors.Errorf("incorrect TlHeader binary length=%d", l)
			}
			t.Data.SrcPort = binary.BigEndian.Uint16(b[:2])
			t.Data.DstPort = binary.BigEndian.Uint16(b[2:4])
			if t.Data.Protocol == proto(unix.IPPROTO_TCP) {
				t.Data.Flags = (b[13] >> 1)
			}

			t.Flag |= NFTNL_TRACE_TRANSPORT_HEADER
		}
	}
	if ad.Err() != nil {
		return errors.WithMessage(err, "failed to unmarshal attribute")
	}
	return nil
}
