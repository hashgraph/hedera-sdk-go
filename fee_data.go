package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type FeeData struct {
	NodeData    *FeeComponents
	NetworkData *FeeComponents
	ServiceData *FeeComponents
}

func newFeeData() FeeData {
	return FeeData{
		NodeData:    nil,
		NetworkData: nil,
		ServiceData: nil,
	}
}

func feeDataFromProtobuf(feeData *services.FeeData) (FeeData, error) {
	if feeData == nil {
		return FeeData{}, errParameterNull
	}

	nodeData, err := feeComponentsFromProtobuf(feeData.Nodedata)
	if err != nil {
		return FeeData{}, err
	}

	networkData, err := feeComponentsFromProtobuf(feeData.Networkdata)
	if err != nil {
		return FeeData{}, err
	}

	serviceData, err := feeComponentsFromProtobuf(feeData.Servicedata)
	if err != nil {
		return FeeData{}, err
	}

	return FeeData{
		NodeData:    &nodeData,
		NetworkData: &networkData,
		ServiceData: &serviceData,
	}, nil
}

func (feeData FeeData) toProtobuf() *services.FeeData {
	var nodeData *services.FeeComponents
	if feeData.NodeData != nil {
		nodeData = feeData.NodeData.toProtobuf()
	}

	var networkData *services.FeeComponents
	if feeData.NetworkData != nil {
		networkData = feeData.NetworkData.toProtobuf()
	}

	var serviceData *services.FeeComponents
	if feeData.ServiceData != nil {
		serviceData = feeData.ServiceData.toProtobuf()
	}

	return &services.FeeData{
		Nodedata:    nodeData,
		Networkdata: networkData,
		Servicedata: serviceData,
	}
}

func (feeData FeeData) ToBytes() []byte {
	data, err := protobuf.Marshal(feeData.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FeeDataFromBytes(data []byte) (FeeData, error) {
	if data == nil {
		return FeeData{}, errByteArrayNull
	}
	pb := services.FeeData{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeData{}, err
	}

	info, err := feeDataFromProtobuf(&pb)
	if err != nil {
		return FeeData{}, err
	}

	return info, nil
}

func (feeData FeeData) String() string {
	return fmt.Sprintf("\nNodedata: %s\nNetworkdata: %s\nServicedata: %s\n", feeData.NodeData.String(), feeData.NetworkData.String(), feeData.ServiceData.String())
}
