package hedera

import (
    "net"

	protobuf "github.com/golang/protobuf/proto"
    "github.com/hashgraph/hedera-sdk-go/proto"
)

type NodeAddress struct {
	IpAddress     net.IP
	Portno        int32
	Memo          string
	RSA_PubKey    string
	NodeId        int64
	NodeAccountID AccountID
	NodeCertHash  []byte
}

func nodeAddressFromProto(node proto.NodeAddress) NodeAddress {
    return NodeAddress{
        IpAddress: node.IpAddress,
        Portno: node.Portno,
        Memo: string(node.Memo),
        RSA_PubKey: node.RSA_PubKey,
        NodeId: node.NodeId,
        NodeAccountID: accountIDFromProto(node.NodeAccountId),
        NodeCertHash: node.NodeCertHash,
    }
}

// NodeAddressBookQuery retrieves the address book
type NodeAddressBookQuery struct {
	QueryBuilder
	pb *proto.FileGetContentsQuery
}

// NewNodeAddressBookQuery creates a FileContentsQuery builder for the address book and decodes the result
func NewNodeAddressBookQuery() *NodeAddressBookQuery {
	pb := &proto.FileGetContentsQuery{
        Header: &proto.QueryHeader{},
        FileID: FileIDForAddressBook().toProto(),
    }

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_FileGetContents{FileGetContents: pb}

	return &NodeAddressBookQuery{inner, pb}
}

// Execute executes the NodeAddressBookQuery using the provided client. The value returned is a list of 
// NodeAddress from the address bookk
func (builder *NodeAddressBookQuery) Execute(client *Client) ([]NodeAddress, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return []NodeAddress{}, err
	}

    bytes := resp.GetFileGetContents().FileContents.Contents

    var nodeAddressBook proto.NodeAddressBook
    protobuf.Unmarshal(bytes, &nodeAddressBook)

    book := make([]NodeAddress, len(nodeAddressBook.NodeAddress))

    for i, node := range nodeAddressBook.NodeAddress {
        book[i] = nodeAddressFromProto(*node)
    }

	return book, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *NodeAddressBookQuery) SetMaxQueryPayment(maxPayment Hbar) *NodeAddressBookQuery {
	return &NodeAddressBookQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *NodeAddressBookQuery) SetQueryPayment(paymentAmount Hbar) *NodeAddressBookQuery {
	return &NodeAddressBookQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *NodeAddressBookQuery) SetQueryPaymentTransaction(tx Transaction) *NodeAddressBookQuery {
	return &NodeAddressBookQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
