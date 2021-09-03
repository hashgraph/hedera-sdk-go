package hedera

import (
	"fmt"
)

type _IPv4AddressPart struct {
	left  byte
	right byte
}

func (ip *_IPv4AddressPart) SetLeft(left byte) *_IPv4AddressPart {
	ip.left = left
	return ip
}

func (ip *_IPv4AddressPart) SetRight(right byte) *_IPv4AddressPart {
	ip.right = right
	return ip
}

func (ip *_IPv4AddressPart) GetLeft() byte {
	return ip.left
}

func (ip *_IPv4AddressPart) GetRight() byte {
	return ip.right
}

func (ip *_IPv4AddressPart) String() string {
	return fmt.Sprintf("%d.%d", uint(ip.left), uint(ip.right))
}
