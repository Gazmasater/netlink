package printtcpudp

import (
	"github.com/Gazmasater/netlink/internal/netlinkparser"

	"fmt"
)

func PrintPacketInfo(packet netlinkparser.PacketInfo) {
	fmt.Printf("Источник IP: %s\n", packet.SrcIP)
	fmt.Printf("Пункт назначения IP: %s\n", packet.DstIP)
	fmt.Printf("Источник порт: %s\n", packet.SrcPort)
	fmt.Printf("Пункт назначения порт: %s\n", packet.DstPort)
	fmt.Printf("Протокол: %s\n", packet.Protocol)
	fmt.Println()
}
