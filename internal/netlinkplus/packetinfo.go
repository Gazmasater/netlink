package netlinkplus

import (
	"github.com/Gazmasater/netlink/internal/data"
	"github.com/Gazmasater/netlink/pkg/logger"
)

type (
	proto      uint8
	PacketInfo struct {
		SrcIP    string
		DstIP    string
		SrcPort  uint16
		DstPort  uint16
		Protocol proto
		Flag     uint16
		logger   logger.TypeOfLogger
	}
)

func (p *PacketInfo) SetLogger(log logger.TypeOfLogger) {
	p.logger = log
}

func (pkt *PacketInfo) IsReady() bool {
	requiredFlags := uint16(data.NFTNL_TRACE_NETWORK_HEADER | data.NFTNL_TRACE_TRANSPORT_HEADER)
	return pkt.Flag&requiredFlags == requiredFlags
}
