package netlinkplus

import "golang.org/x/sys/unix"

func (p proto) String() string {
	switch p {
	case unix.IPPROTO_TCP:
		return "tcp"
	case unix.IPPROTO_UDP:
		return "udp"
	case unix.IPPROTO_UDPLITE:
		return "udplite"
	case unix.IPPROTO_ESP:
		return "esp"
	case unix.IPPROTO_AH:
		return "ah"
	case unix.IPPROTO_ICMP:
		return "icmp"
	case unix.IPPROTO_ICMPV6:
		return "icmpv6"
	case unix.IPPROTO_COMP:
		return "comp"
	case unix.IPPROTO_DCCP:
		return "dccp"
	case unix.IPPROTO_SCTP:
		return "sctp"
	}
	return "unknown"
}
