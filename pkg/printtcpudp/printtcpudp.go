package printtcpudp

import (
	"fmt"

	"github.com/Gazmasater/netlink/internal/netlinkdecode"
)

func PrintPacketInfo(packet netlinkdecode.PacketInfo) {
	fmt.Printf("srsIp:%s, dstIp:%s, srcPort:%s, dstPort:%s, protocol:%s\n",
		packet.SrcIP, packet.DstIP, packet.SrcPort,
		packet.DstPort, packet.Protocol)
}
