package netlinkplus

import (
	"fmt"
)

type consolePacketPrinter struct{}

func NewPrinter() PacketPrinter {
	return &consolePacketPrinter{}
}

func (c consolePacketPrinter) PrintHeader(header string) {
	fmt.Println(header)
	// ^^^^^^^ тут должен быть логгер а не Println
}

func (c consolePacketPrinter) PrintPacketInfo(packet PacketInfo) {
	fmt.Printf("srcIp:%s, dstIp:%s, srcPort:%d, dstPort:%d, protocol:%s\n",
		//   ^^^^^^^ тут должен быть логгер а не Printf
		packet.SrcIP, packet.DstIP, packet.SrcPort,
		packet.DstPort, packet.Protocol)
}
