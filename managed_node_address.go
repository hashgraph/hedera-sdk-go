package hedera

import (
	"regexp"
	"strconv"
)

var hostAndPort = regexp.MustCompile(`^(\S+):(\d+)$`)

type _ManagedNodeAddress struct {
	address *string
	port    uint32
}

func _ManagedNodeAddressFromString(str string) *_ManagedNodeAddress {
	hostAndPortMatch := hostAndPort.FindStringSubmatch(str)

	if len(hostAndPortMatch) > 1 {
		temp, err := strconv.ParseUint(hostAndPortMatch[2], 10, 64)
		if err != nil {
			return &_ManagedNodeAddress{}
		}
		return &_ManagedNodeAddress{
			address: &hostAndPortMatch[1],
			port:    uint32(temp)}
	}

	panic("failed to parse node address")
}

func (address *_ManagedNodeAddress) _IsTransportSecurity() bool {
	return address.port == 50212 || address.port == 433
}

func (address *_ManagedNodeAddress) _ToInsecure() *_ManagedNodeAddress {
	var port uint32

	switch address.port {
	case 50212:
		port = 50211
	case 433:
		port = 5600
	}

	return &_ManagedNodeAddress{
		address: address.address,
		port:    port,
	}
}

func (address *_ManagedNodeAddress) _ToSecure() *_ManagedNodeAddress {
	var port uint32

	switch address.port {
	case 50211:
		return &_ManagedNodeAddress{
			address: address.address,
			port:    50212,
		}
	case 5600:
		return &_ManagedNodeAddress{
			address: address.address,
			port:    433,
		}
	}

	return &_ManagedNodeAddress{
		address: address.address,
		port:    port,
	}
}

func (address *_ManagedNodeAddress) _String() string {
	if address.address != nil {
		return *address.address + ":" + strconv.FormatInt(int64(address.port), 10)
	}

	return ""
}
