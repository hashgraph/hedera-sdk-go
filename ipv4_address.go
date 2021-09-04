package hedera

type _IPv4Address struct {
	network _IPv4AddressPart
	host    _IPv4AddressPart
}

func _Ipv4AddressFromProtobuf(byte []byte) _IPv4Address {
	return _IPv4Address{
		network: _IPv4AddressPart{
			left:  byte[0],
			right: byte[1],
		},
		host: _IPv4AddressPart{
			left:  byte[2],
			right: byte[3],
		},
	}
}

func (ip *_IPv4Address) _ToProtobuf() []byte {
	return []byte{ip.network.left, ip.network.right, ip.host.left, ip.host.right}
}

func (ip *_IPv4Address) String() string {
	return ip.network.String() + "." + ip.host.String()
}
