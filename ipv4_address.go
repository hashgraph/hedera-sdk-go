package hedera

type ipv4Address struct {
	network ipv4AddressPart
	host    ipv4AddressPart
}

func ipv4AddressFromProtobuf(byte []byte) ipv4Address {
	return ipv4Address{
		network: ipv4AddressPart{
			left:  byte[0],
			right: byte[1],
		},
		host: ipv4AddressPart{
			left:  byte[2],
			right: byte[3],
		},
	}
}

func (ip *ipv4Address) toProtobuf() []byte {
	return []byte{ip.network.left, ip.network.right, ip.host.left, ip.host.right}
}

func (ip *ipv4Address) String() string {
	return ip.network.String() + "." + ip.host.String()
}
