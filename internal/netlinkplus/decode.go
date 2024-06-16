// TODO перенести этот файл в пакет netlink (см TODO)
package netlinkplus

import (
	"encoding/binary"
	"net"

	"github.com/Gazmasater/netlink/internal/data"

	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

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
