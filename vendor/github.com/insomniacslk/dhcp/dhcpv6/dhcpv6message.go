package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/iana"
	"github.com/u-root/u-root/pkg/rand"
	"github.com/u-root/u-root/pkg/uio"
)

const MessageHeaderSize = 4

// Message represents a DHCPv6 Message as defined by RFC 3315 Section 6.
type Message struct {
	MessageType   MessageType
	TransactionID TransactionID
	Options       Options
}

var randomRead = rand.Read

// GenerateTransactionID generates a random 3-byte transaction ID.
func GenerateTransactionID() (TransactionID, error) {
	var tid TransactionID
	n, err := randomRead(tid[:])
	if err != nil {
		return tid, err
	}
	if n != len(tid) {
		return tid, fmt.Errorf("invalid random sequence: shorter than 3 bytes")
	}
	return tid, nil
}

// GetTime returns a time integer suitable for DUID-LLT, i.e. the current time counted
// in seconds since January 1st, 2000, midnight UTC, modulo 2^32
func GetTime() uint32 {
	now := time.Since(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	return uint32((now.Nanoseconds() / 1000000000) % 0xffffffff)
}

// NewSolicitWithCID creates a new SOLICIT message with CID.
func NewSolicitWithCID(duid Duid, modifiers ...Modifier) (*Message, error) {
	m, err := NewMessage()
	if err != nil {
		return nil, err
	}
	m.MessageType = MessageTypeSolicit
	m.AddOption(&OptClientId{Cid: duid})
	oro := new(OptRequestedOption)
	oro.SetRequestedOptions([]OptionCode{
		OptionDNSRecursiveNameServer,
		OptionDomainSearchList,
	})
	m.AddOption(oro)
	m.AddOption(&OptElapsedTime{})
	// FIXME use real values for IA_NA
	iaNa := &OptIANA{}
	iaNa.IaId = [4]byte{0xfa, 0xce, 0xb0, 0x0c}
	iaNa.T1 = 0xe10
	iaNa.T2 = 0x1518
	m.AddOption(iaNa)
	// Apply modifiers
	for _, mod := range modifiers {
		mod(m)
	}
	return m, nil
}

// NewSolicit creates a new SOLICIT message with DUID-LLT, using the
// given network interface's hardware address and current time
func NewSolicit(ifaceHWAddr net.HardwareAddr, modifiers ...Modifier) (*Message, error) {
	duid := Duid{
		Type:          DUID_LLT,
		HwType:        iana.HWTypeEthernet,
		Time:          GetTime(),
		LinkLayerAddr: ifaceHWAddr,
	}
	return NewSolicitWithCID(duid, modifiers...)
}

// NewAdvertiseFromSolicit creates a new ADVERTISE packet based on an SOLICIT packet.
func NewAdvertiseFromSolicit(sol *Message, modifiers ...Modifier) (*Message, error) {
	if sol == nil {
		return nil, errors.New("SOLICIT cannot be nil")
	}
	if sol.Type() != MessageTypeSolicit {
		return nil, errors.New("The passed SOLICIT must have SOLICIT type set")
	}
	// build ADVERTISE from SOLICIT
	adv := &Message{
		MessageType:   MessageTypeAdvertise,
		TransactionID: sol.TransactionID,
	}
	// add Client ID
	cid := sol.GetOneOption(OptionClientID)
	if cid == nil {
		return nil, errors.New("Client ID cannot be nil in SOLICIT when building ADVERTISE")
	}
	adv.AddOption(cid)

	// apply modifiers
	for _, mod := range modifiers {
		mod(adv)
	}
	return adv, nil
}

// NewRequestFromAdvertise creates a new REQUEST packet based on an ADVERTISE
// packet options.
func NewRequestFromAdvertise(adv *Message, modifiers ...Modifier) (*Message, error) {
	if adv == nil {
		return nil, errors.New("ADVERTISE cannot be nil")
	}
	if adv.MessageType != MessageTypeAdvertise {
		return nil, fmt.Errorf("The passed ADVERTISE must have ADVERTISE type set")
	}
	// build REQUEST from ADVERTISE
	req := &Message{
		MessageType:   MessageTypeRequest,
		TransactionID: adv.TransactionID,
	}
	// add Client ID
	cid := adv.GetOneOption(OptionClientID)
	if cid == nil {
		return nil, fmt.Errorf("Client ID cannot be nil in ADVERTISE when building REQUEST")
	}
	req.AddOption(cid)
	// add Server ID
	sid := adv.GetOneOption(OptionServerID)
	if sid == nil {
		return nil, fmt.Errorf("Server ID cannot be nil in ADVERTISE when building REQUEST")
	}
	req.AddOption(sid)
	// add Elapsed Time
	req.AddOption(&OptElapsedTime{})
	// add IA_NA
	iaNa := adv.GetOneOption(OptionIANA)
	if iaNa == nil {
		return nil, fmt.Errorf("IA_NA cannot be nil in ADVERTISE when building REQUEST")
	}
	req.AddOption(iaNa)
	// add OptRequestedOption
	oro := OptRequestedOption{}
	oro.SetRequestedOptions([]OptionCode{
		OptionDNSRecursiveNameServer,
		OptionDomainSearchList,
	})
	req.AddOption(&oro)
	// add OPTION_VENDOR_CLASS, only if present in the original request
	// TODO implement OptionVendorClass
	vClass := adv.GetOneOption(OptionVendorClass)
	if vClass != nil {
		req.AddOption(vClass)
	}

	// apply modifiers
	for _, mod := range modifiers {
		mod(req)
	}
	return req, nil
}

// NewReplyFromMessage creates a new REPLY packet based on a
// Message. The function is to be used when generating a reply to
// REQUEST, CONFIRM, RENEW, REBIND, RELEASE and INFORMATION-REQUEST packets.
func NewReplyFromMessage(msg *Message, modifiers ...Modifier) (*Message, error) {
	if msg == nil {
		return nil, errors.New("Message cannot be nil")
	}
	switch msg.Type() {
	case MessageTypeRequest, MessageTypeConfirm, MessageTypeRenew,
		MessageTypeRebind, MessageTypeRelease, MessageTypeInformationRequest:
	default:
		return nil, errors.New("Cannot create REPLY from the passed message type set")
	}

	// build REPLY from MESSAGE
	rep := &Message{
		MessageType:   MessageTypeReply,
		TransactionID: msg.TransactionID,
	}
	// add Client ID
	cid := msg.GetOneOption(OptionClientID)
	if cid == nil {
		return nil, errors.New("Client ID cannot be nil when building REPLY")
	}
	rep.AddOption(cid)

	// apply modifiers
	for _, mod := range modifiers {
		mod(rep)
	}
	return rep, nil
}

// Type returns this message's message type.
func (m Message) Type() MessageType {
	return m.MessageType
}

// GetInnerMessage returns the message itself.
func (m *Message) GetInnerMessage() (*Message, error) {
	return m, nil
}

// AddOption adds an option to this message.
func (m *Message) AddOption(option Option) {
	m.Options.Add(option)
}

// UpdateOption updates the existing options with the passed option, adding it
// at the end if not present already
func (m *Message) UpdateOption(option Option) {
	m.Options.Update(option)
}

// IsNetboot returns true if the machine is trying to netboot. It checks if
// "boot file" is one of the requested options, which is useful for
// SOLICIT/REQUEST packet types, it also checks if the "boot file" option is
// included in the packet, which is useful for ADVERTISE/REPLY packet.
func (m *Message) IsNetboot() bool {
	if m.IsOptionRequested(OptionBootfileURL) {
		return true
	}
	if optbf := m.GetOneOption(OptionBootfileURL); optbf != nil {
		return true
	}
	return false
}

// IsOptionRequested takes an OptionCode and returns true if that option is
// within the requested options of the DHCPv6 message.
func (m *Message) IsOptionRequested(requested OptionCode) bool {
	for _, optoro := range m.GetOption(OptionORO) {
		for _, o := range optoro.(*OptRequestedOption).RequestedOptions() {
			if o == requested {
				return true
			}
		}
	}
	return false
}

// String returns a short human-readable string for this message.
func (m *Message) String() string {
	return fmt.Sprintf("Message(messageType=%s transactionID=%s, %d options)",
		m.MessageType, m.TransactionID, len(m.Options))
}

// Summary prints all options associated with this message.
func (m *Message) Summary() string {
	ret := fmt.Sprintf(
		"Message\n"+
			"  messageType=%s\n"+
			"  transactionid=%s\n",
		m.MessageType,
		m.TransactionID,
	)
	ret += "  options=["
	if len(m.Options) > 0 {
		ret += "\n"
	}
	for _, opt := range m.Options {
		ret += fmt.Sprintf("    %v\n", opt.String())
	}
	ret += "  ]\n"
	return ret
}

// ToBytes returns the serialized version of this message as defined by RFC
// 3315, Section 5.
func (m *Message) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write8(uint8(m.MessageType))
	buf.WriteBytes(m.TransactionID[:])
	buf.WriteBytes(m.Options.ToBytes())
	return buf.Data()
}

// GetOption returns the options associated with the code.
func (m *Message) GetOption(code OptionCode) []Option {
	return m.Options.Get(code)
}

// GetOneOption returns the first associated option with the code from this
// message.
func (m *Message) GetOneOption(code OptionCode) Option {
	return m.Options.GetOne(code)
}

// IsRelay returns whether this is a relay message or not.
func (m *Message) IsRelay() bool {
	return false
}
