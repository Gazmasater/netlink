package netlinkplus

import (
	"fmt"
	"os"
)

func (pkt *PacketInfo) LogPacketInfo() {
	var protocolName string
	switch pkt.Protocol {
	case 6:
		protocolName = "TCP"
	case 17:
		protocolName = "UDP"
	default:
		protocolName = "Unknown"
	}
	pkt.logger.Infof("Packet Information: SrcIP=%s, DstIP=%s, SrcPort=%d, DstPort=%d, Protocol=%s",
		pkt.SrcIP, pkt.DstIP, pkt.SrcPort, pkt.DstPort, protocolName)
}

func (pkt *PacketInfo) LogPacketFile() {
	var protocolName string
	switch pkt.Protocol {
	case 6:
		protocolName = "TCP"
	case 17:
		protocolName = "UDP"
	default:
		protocolName = "Unknown"
	}

	// Open the file for writing (append mode), create it if it doesn't exist
	file, err := os.OpenFile("packet_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		pkt.logger.Fatalf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the packet information to the file
	logMessage := fmt.Sprintf("Packet Information: SrcIP=%s, DstIP=%s, SrcPort=%d, DstPort=%d, Protocol=%s\n",
		pkt.SrcIP, pkt.DstIP, pkt.SrcPort, pkt.DstPort, protocolName)

	if _, err := file.WriteString(logMessage); err != nil {
		pkt.logger.Fatalf("Error writing to file: %v\n", err)
	}
}
