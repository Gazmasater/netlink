package printtcpudp

import "github.com/Gazmasater/netlink/internal/netlinkdecode"

type PacketPrinter interface {
	PrintHeader(header string)
	PrintPacketInfo(packet netlinkdecode.PacketInfo)
}
