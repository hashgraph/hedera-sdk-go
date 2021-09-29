package hedera

import (
	"regexp"
	"strconv"
)

type _ManagedNodeAddress struct {
	name    *string
	address *string
	port    uint32
}

func _ManagedNodeAddressFromString(str string) *_ManagedNodeAddress {
	hostAndPort := regexp.MustCompile(`(^.*):(\d+$)`)
	inProcess := regexp.MustCompile(`in-process:(.*)`)

	hostAndPortMatch := hostAndPort.FindStringSubmatch(str)
	inProcessMatch := inProcess.FindStringSubmatch(str)

	//for _, s := range hostAndPortMatch {
	//	println(s)
	//}
	//for _, s := range inProcessMatch {
	//	println(s)
	//}

	if hostAndPortMatch != nil {
		if len(hostAndPortMatch) > 1 {
			temp, err := strconv.ParseUint(hostAndPortMatch[2], 10, 64)
			if err != nil {
				return &_ManagedNodeAddress{}
			}
			return &_ManagedNodeAddress{
				name:    nil,
				address: &hostAndPortMatch[1],
				port:    uint32(temp),
			}
		} else {
			panic("failed to parse node address")
		}
	} else if inProcess != nil {
		return &_ManagedNodeAddress{
			name:    &inProcessMatch[1],
			address: nil,
			port:    0,
		}
	} else {
		panic("failed to parse node address")
	}
}

func (address *_ManagedNodeAddress) _IsInProcess() bool {
	return address.name != nil
}

func (address *_ManagedNodeAddress) _IsTransportSecurity() bool {
	return address.port == 50212 || address.port == 433
}

func (address *_ManagedNodeAddress) _ToInsecure() *_ManagedNodeAddress {
	var port uint32

	switch address.port {
	case 50212:
		port = 50211
		break
	case 433:
		port = 5600
		break
	}

	return &_ManagedNodeAddress{
		name:    address.name,
		address: address.address,
		port:    port,
	}
}

func (address *_ManagedNodeAddress) _ToSecure() *_ManagedNodeAddress {
	var port uint32

	switch address.port {
	case 50211:
		port = 50212
		break
	case 5600:
		port = 433
		break
	}

	return &_ManagedNodeAddress{
		name:    address.name,
		address: address.address,
		port:    port,
	}
}

func (address *_ManagedNodeAddress) _String() string {
	if address.name != nil {
		return "in-process:" + *address.name
	}

	if address.address != nil {
		return *address.address + ":" + strconv.FormatInt(int64(address.port), 10)
	}

	return ""
}
