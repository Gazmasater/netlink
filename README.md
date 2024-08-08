sudo nft add table inet filter

sudo nft add chain inet filter input { type filter hook input priority 0 \; }

sudo nft add rule inet filter input tcp dport 0-65535 log prefix \"TCP: \" flags all counter

sudo nft add rule inet filter input tcp dport 0-65535 meta nftrace set 1
