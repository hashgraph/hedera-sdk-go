package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type FeeData struct {
	NodeData    *FeeComponents
	NetworkData *FeeComponents
	ServiceData *FeeComponents
}

func _FeeDataFromProtobuf(feeData *proto.FeeData) (FeeData, error) {
	if feeData == nil {
		return FeeData{}, errParameterNull
	}

	nodeData, err := _FeeComponentsFromProtobuf(feeData.Nodedata)
	if err != nil {
		return FeeData{}, err
	}

	networkData, err := _FeeComponentsFromProtobuf(feeData.Networkdata)
	if err != nil {
		return FeeData{}, err
	}

	serviceData, err := _FeeComponentsFromProtobuf(feeData.Servicedata)
	if err != nil {
		return FeeData{}, err
	}

	return FeeData{
		NodeData:    &nodeData,
		NetworkData: &networkData,
		ServiceData: &serviceData,
	}, nil
}

func (feeData FeeData) _ToProtobuf() *proto.FeeData {
	var nodeData *proto.FeeComponents
	if feeData.NodeData != nil {
		nodeData = feeData.NodeData._ToProtobuf()
	}

	var networkData *proto.FeeComponents
	if feeData.NetworkData != nil {
		networkData = feeData.NetworkData._ToProtobuf()
	}

	var serviceData *proto.FeeComponents
	if feeData.ServiceData != nil {
		serviceData = feeData.ServiceData._ToProtobuf()
	}

	return &proto.FeeData{
		Nodedata:    nodeData,
		Networkdata: networkData,
		Servicedata: serviceData,
	}
}

func (feeData FeeData) ToBytes() []byte {
	data, err := protobuf.Marshal(feeData._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FeeDataFromBytes(data []byte) (FeeData, error) {
	if data == nil {
		return FeeData{}, errByteArrayNull
	}
	pb := proto.FeeData{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeData{}, err
	}

	info, err := _FeeDataFromProtobuf(&pb)
	if err != nil {
		return FeeData{}, err
	}

	return info, nil
}

func (feeData FeeData) String() string {
	return fmt.Sprintf("\nNodedata: %s\nNetworkdata: %s\nServicedata: %s\n", feeData.NodeData.String(), feeData.NetworkData.String(), feeData.ServiceData.String())
}
