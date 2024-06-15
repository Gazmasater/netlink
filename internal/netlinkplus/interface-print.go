package netlinkplus

type PacketPrinter interface {
	PrintHeader(header string)
	PrintPacketInfo(packet PacketInfo)
}
