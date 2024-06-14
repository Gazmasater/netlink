package printtcpudp

import (
	"fmt"

	"github.com/Gazmasater/netlink/internal/netlinkdecode"
)

type ConsolePacketPrinter struct{}

func (c ConsolePacketPrinter) PrintHeader(header string) {
	fmt.Println(header)
}

func (c ConsolePacketPrinter) PrintPacketInfo(packet netlinkdecode.PacketInfo) {
	fmt.Printf("srcIp:%s, dstIp:%s, srcPort:%s, dstPort:%s, protocol:%s\n",
		packet.SrcIP, packet.DstIP, packet.SrcPort,
		packet.DstPort, packet.Protocol)
}
