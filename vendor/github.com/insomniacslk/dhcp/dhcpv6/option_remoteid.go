package dhcpv6

import (
	"fmt"

	"github.com/u-root/u-root/pkg/uio"
)

// OptRemoteId implemens the Remote ID option.
//
// https://www.ietf.org/rfc/rfc4649.txt
type OptRemoteId struct {
	enterpriseNumber uint32
	remoteId         []byte
}

func (op *OptRemoteId) Code() OptionCode {
	return OptionRemoteID
}

// ToBytes serializes this option to a byte stream.
func (op *OptRemoteId) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write32(op.enterpriseNumber)
	buf.WriteBytes(op.remoteId)
	return buf.Data()
}

func (op *OptRemoteId) EnterpriseNumber() uint32 {
	return op.enterpriseNumber
}

func (op *OptRemoteId) SetEnterpriseNumber(enterpriseNumber uint32) {
	op.enterpriseNumber = enterpriseNumber
}

func (op *OptRemoteId) RemoteID() []byte {
	return op.remoteId
}

func (op *OptRemoteId) SetRemoteID(remoteId []byte) {
	op.remoteId = append([]byte(nil), remoteId...)
}

func (op *OptRemoteId) String() string {
	return fmt.Sprintf("OptRemoteId{enterprisenum=%v, remoteid=%v}",
		op.enterpriseNumber, op.remoteId,
	)
}

// ParseOptRemoteId builds an OptRemoteId structure from a sequence of bytes.
// The input data does not include option code and length bytes.
func ParseOptRemoteId(data []byte) (*OptRemoteId, error) {
	var opt OptRemoteId
	buf := uio.NewBigEndianBuffer(data)
	opt.enterpriseNumber = buf.Read32()
	opt.remoteId = buf.ReadAll()
	return &opt, buf.FinError()
}
