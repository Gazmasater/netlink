package printtcpudp

import (
	"github.com/Gazmasater/netlink/internal/netlinkparser"

	"fmt"
)

func PrintPacketInfo(packet netlinkparser.PacketInfo) {
	fmt.Printf("srsIp:%s, dstIp:%s, srcPort:%s, dstPort:%s, protocol:%s",
		packet.SrcIP, packet.DstIP, packet.SrcPort,
		packet.DstPort, packet.Protocol)
}
