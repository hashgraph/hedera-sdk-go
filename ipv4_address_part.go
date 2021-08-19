package hedera

import (
	"fmt"
)

type ipv4AddressPart struct {
	left  byte
	right byte
}

func (ip *ipv4AddressPart) SetLeft(left byte) *ipv4AddressPart {
	ip.left = left
	return ip
}

func (ip *ipv4AddressPart) SetRight(right byte) *ipv4AddressPart {
	ip.right = right
	return ip
}

func (ip *ipv4AddressPart) GetLeft() byte {
	return ip.left
}

func (ip *ipv4AddressPart) GetRight() byte {
	return ip.right
}

func (ip *ipv4AddressPart) String() string {
	return fmt.Sprintf("%d.%d", uint(ip.left), uint(ip.right))
}
