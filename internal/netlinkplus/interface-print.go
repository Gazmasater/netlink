package netlinkplus

import (
	"context"

	"github.com/Gazmasater/netlink/pkg/logger"
)

type PacketPrinter interface {
	PrintHeader(header string)
	PrintPacketInfo(packet PacketInfo)
}

type consolePacketPrinter struct{}

func NewPrinter() PacketPrinter {
	return &consolePacketPrinter{}
}

func (c consolePacketPrinter) PrintHeader(header string) {
	logger.Info(context.Background(), header)
}

func (c consolePacketPrinter) PrintPacketInfo(packet PacketInfo) {
	logger.Infof(context.Background(), "srcIp:%s, dstIp:%s, srcPort:%d, dstPort:%d, protocol:%s",
		packet.SrcIP, packet.DstIP, packet.SrcPort,
		packet.DstPort, packet.Protocol)
}
