package nftrace

const (
	NFTNL_TRACE_NETWORK_HEADER = 1 << iota
	NFTNL_TRACE_TRANSPORT_HEADER
)

const (
	// Transport layer header length
	TlHeaderLen = 8
	// Network layer header length
	NlHeaderLen = 20
	// Offset attribute data in the nft netlink group message
	NlNftAttrOffset = 4
)
