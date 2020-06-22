package iscsinl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/vishvananda/netlink/nl"
	"golang.org/x/sys/unix"
)

// STOP_CONN_RECOVER - when stopping connection clean up I/O on that connection
const STOP_CONN_RECOVER = 0x3

// IscsiEvent iscsi_if.h:enum iscsi_uevent_e
type IscsiEvent uint32

// IscsiEvents that we use
const (
	UEVENT_BASE                  IscsiEvent = 10
	KEVENT_BASE                             = 100
	ISCSI_UEVENT_CREATE_SESSION             = UEVENT_BASE + 1
	ISCSI_UEVENT_DESTROY_SESSION            = UEVENT_BASE + 2
	ISCSI_UEVENT_CREATE_CONN                = UEVENT_BASE + 3
	ISCSI_UEVENT_DESTROY_CONN               = UEVENT_BASE + 4
	ISCSI_UEVENT_BIND_CONN                  = UEVENT_BASE + 5
	ISCSI_UEVENT_SET_PARAM                  = UEVENT_BASE + 6
	ISCSI_UEVENT_START_CONN                 = UEVENT_BASE + 7
	ISCSI_UEVENT_STOP_CONN                  = UEVENT_BASE + 8
	ISCSI_UEVENT_SEND_PDU                   = UEVENT_BASE + 9

	ISCSI_KEVENT_RECV_PDU       = KEVENT_BASE + 1
	ISCSI_KEVENT_CONN_ERROR     = KEVENT_BASE + 2
	ISCSI_KEVENT_IF_ERROR       = KEVENT_BASE + 3
	ISCSI_KEVENT_CREATE_SESSION = KEVENT_BASE + 6
)

// IscsiParam iscsi_if.h:enum iscsi_param
type IscsiParam uint32

// IscsiParams up until INITIATOR_NAME
const (
	ISCSI_PARAM_MAX_RECV_DLENGTH IscsiParam = iota
	ISCSI_PARAM_MAX_XMIT_DLENGTH
	ISCSI_PARAM_HDRDGST_EN
	ISCSI_PARAM_DATADGST_EN
	ISCSI_PARAM_INITIAL_R2T_EN
	ISCSI_PARAM_MAX_R2T
	ISCSI_PARAM_IMM_DATA_EN
	ISCSI_PARAM_FIRST_BURST
	ISCSI_PARAM_MAX_BURST
	ISCSI_PARAM_PDU_INORDER_EN
	ISCSI_PARAM_DATASEQ_INORDER_EN
	ISCSI_PARAM_ERL
	ISCSI_PARAM_IFMARKER_EN
	ISCSI_PARAM_OFMARKER_EN
	ISCSI_PARAM_EXP_STATSN
	ISCSI_PARAM_TARGET_NAME
	ISCSI_PARAM_TPGT
	ISCSI_PARAM_PERSISTENT_ADDRESS
	ISCSI_PARAM_PERSISTENT_PORT
	ISCSI_PARAM_SESS_RECOVERY_TMO
	ISCSI_PARAM_CONN_PORT
	ISCSI_PARAM_CONN_ADDRESS
	ISCSI_PARAM_USERNAME
	ISCSI_PARAM_USERNAME_IN
	ISCSI_PARAM_PASSWORD
	ISCSI_PARAM_PASSWORD_IN
	ISCSI_PARAM_FAST_ABORT
	ISCSI_PARAM_ABORT_TMO
	ISCSI_PARAM_LU_RESET_TMO
	ISCSI_PARAM_HOST_RESET_TMO
	ISCSI_PARAM_PING_TMO
	ISCSI_PARAM_RECV_TMO
	ISCSI_PARAM_IFACE_NAME
	ISCSI_PARAM_ISID
	ISCSI_PARAM_INITIATOR_NAME
)

var paramToString = map[IscsiParam]string{
	ISCSI_PARAM_TARGET_NAME:        "Target Name",
	ISCSI_PARAM_INITIATOR_NAME:     "Initiator Name",
	ISCSI_PARAM_MAX_RECV_DLENGTH:   "Max Recv DLength",
	ISCSI_PARAM_MAX_XMIT_DLENGTH:   "Max Xmit DLenght",
	ISCSI_PARAM_FIRST_BURST:        "First Burst",
	ISCSI_PARAM_MAX_BURST:          "Max Burst",
	ISCSI_PARAM_PDU_INORDER_EN:     "PDU Inorder EN",
	ISCSI_PARAM_DATASEQ_INORDER_EN: "Data Seq In Order EN",
	ISCSI_PARAM_INITIAL_R2T_EN:     "Inital R2T EN",
	ISCSI_PARAM_IMM_DATA_EN:        "Immediate Data EN",
	ISCSI_PARAM_EXP_STATSN:         "Exp Statsn",
	ISCSI_PARAM_HDRDGST_EN:         "HDR Digest EN",
	ISCSI_PARAM_DATADGST_EN:        "Data Digest EN",
	ISCSI_PARAM_PING_TMO:           "Ping TMO",
	ISCSI_PARAM_RECV_TMO:           "Recv TMO",
}

func (p IscsiParam) String() string {
	val, ok := paramToString[p]
	if !ok {
		return fmt.Sprintf("IscsiParam(%d)", int(p))
	}

	return val
}

// IscsiErr iscsi_if.h:enum iscsi_err
type IscsiErr uint32

// IscsiErr iscsi_if.h:enum iscsi_err
const (
	ISCSI_OK                      IscsiErr = 0
	ISCSI_ERR_BASE                         = 1000
	ISCSI_ERR_DATASN                       = ISCSI_ERR_BASE + 1
	ISCSI_ERR_DATA_OFFSET                  = ISCSI_ERR_BASE + 2
	ISCSI_ERR_MAX_CMDSN                    = ISCSI_ERR_BASE + 3
	ISCSI_ERR_EXP_CMDSN                    = ISCSI_ERR_BASE + 4
	ISCSI_ERR_BAD_OPCODE                   = ISCSI_ERR_BASE + 5
	ISCSI_ERR_DATALEN                      = ISCSI_ERR_BASE + 6
	ISCSI_ERR_AHSLEN                       = ISCSI_ERR_BASE + 7
	ISCSI_ERR_PROTO                        = ISCSI_ERR_BASE + 8
	ISCSI_ERR_LUN                          = ISCSI_ERR_BASE + 9
	ISCSI_ERR_BAD_ITT                      = ISCSI_ERR_BASE + 10
	ISCSI_ERR_CONN_FAILED                  = ISCSI_ERR_BASE + 11
	ISCSI_ERR_R2TSN                        = ISCSI_ERR_BASE + 12
	ISCSI_ERR_SESSION_FAILED               = ISCSI_ERR_BASE + 13
	ISCSI_ERR_HDR_DGST                     = ISCSI_ERR_BASE + 14
	ISCSI_ERR_DATA_DGST                    = ISCSI_ERR_BASE + 15
	ISCSI_ERR_PARAM_NOT_FOUND              = ISCSI_ERR_BASE + 16
	ISCSI_ERR_NO_SCSI_CMD                  = ISCSI_ERR_BASE + 17
	ISCSI_ERR_INVALID_HOST                 = ISCSI_ERR_BASE + 18
	ISCSI_ERR_XMIT_FAILED                  = ISCSI_ERR_BASE + 19
	ISCSI_ERR_TCP_CONN_CLOSE               = ISCSI_ERR_BASE + 20
	ISCSI_ERR_SCSI_EH_SESSION_RST          = ISCSI_ERR_BASE + 21
	ISCSI_ERR_NOP_TIMEDOUT                 = ISCSI_ERR_BASE + 22
)

func (e IscsiErr) String() string {
	switch e {
	case ISCSI_OK:
		return "ISCSI_OK"
	case ISCSI_ERR_BASE:
		return "ISCSI_ERR_BASE"
	case ISCSI_ERR_DATASN:
		return "ISCSI_ERR_DATASN"
	case ISCSI_ERR_DATA_OFFSET:
		return "ISCSI_ERR_DATA_OFFSET"
	case ISCSI_ERR_MAX_CMDSN:
		return "ISCSI_ERR_MAX_CMDSN"
	case ISCSI_ERR_EXP_CMDSN:
		return "ISCSI_ERR_EXP_CMDSN"
	case ISCSI_ERR_BAD_OPCODE:
		return "ISCSI_ERR_BAD_OPCODE"
	case ISCSI_ERR_DATALEN:
		return "ISCSI_ERR_DATALEN"
	case ISCSI_ERR_AHSLEN:
		return "ISCSI_ERR_AHSLEN"
	case ISCSI_ERR_PROTO:
		return "ISCSI_ERR_PROTO"
	case ISCSI_ERR_LUN:
		return "ISCSI_ERR_LUN"
	case ISCSI_ERR_BAD_ITT:
		return "ISCSI_ERR_BAD_ITT"
	case ISCSI_ERR_CONN_FAILED:
		return "ISCSI_ERR_CONN_FAILED"
	case ISCSI_ERR_R2TSN:
		return "ISCSI_ERR_R2TSN"
	case ISCSI_ERR_SESSION_FAILED:
		return "ISCSI_ERR_SESSION_FAILED"
	case ISCSI_ERR_HDR_DGST:
		return "ISCSI_ERR_HDR_DGST"
	case ISCSI_ERR_DATA_DGST:
		return "ISCSI_ERR_DATA_DGST"
	case ISCSI_ERR_PARAM_NOT_FOUND:
		return "ISCSI_ERR_PARAM_NOT_FOUND"
	case ISCSI_ERR_NO_SCSI_CMD:
		return "ISCSI_ERR_NO_SCSI_CMD"
	case ISCSI_ERR_INVALID_HOST:
		return "ISCSI_ERR_INVALID_HOST"
	case ISCSI_ERR_XMIT_FAILED:
		return "ISCSI_ERR_XMIT_FAILED"
	case ISCSI_ERR_TCP_CONN_CLOSE:
		return "ISCSI_ERR_TCP_CONN_CLOSE"
	case ISCSI_ERR_SCSI_EH_SESSION_RST:
		return "ISCSI_ERR_SCSI_EH_SESSION_RST"
	case ISCSI_ERR_NOP_TIMEDOUT:
		return "ISCSI_ERR_NOP_TIMEDOUT"
	default:
		return strconv.Itoa(int(e))
	}
}

// All Iscsu[U/K]event structs must be the same size,
// with UserArg always 24 bytes, ReturnArg always 16 bytes

// iSCSIUEvent is a generic uevent, for reading header content
type iSCSIUEvent struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	UserArg         [24]byte
	ReturnArg       [16]byte
}

// iSCSIUEventCreateSession corresponds with ISCSI_UEVENT_CREATE_SESSION
type iSCSIUEventCreateSession struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	CSession        struct {
		InitialCmdSN uint32
		CmdsMax      uint16
		QueueDepth   uint16
		_            [16]byte
	}
	CSessionRet struct {
		Sid    uint32
		HostNo uint32
		_      [8]byte
	}
}

type iSCSIUEventDestroySession struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	DSession        struct {
		Sid uint32
		_   [20]byte
	}
	DSessionRet struct {
		Ret uint32
		_   [12]byte
	}
}

type iSCSIUEventCreateConnection struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	CConn           struct {
		Sid uint32
		Cid uint32
		_   [16]byte
	}
	CConnRet struct {
		Sid uint32
		Cid uint32
		_   [8]byte
	}
}
type iSCSIUEventDestroyConnection struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	DConn           struct {
		Sid uint32
		Cid uint32
		_   [16]byte
	}
	DConnRet struct {
		Ret uint32
		_   [12]byte
	}
}

type iSCSIUEventBindConnection struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	BConn           struct {
		Sid          uint32
		Cid          uint32
		TransportEph uint64
		IsLeading    uint32
		_            [4]byte
	}
	BConnRet struct {
		Ret uint32
		_   [12]byte
	}
}

type iSCSIUEventSetParam struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	SetParam        struct {
		Sid   uint32
		Cid   uint32
		Param IscsiParam
		Len   uint32
		_     [8]byte
	}
	SetParamRet struct {
		Ret uint32
		_   [12]byte
	}
}

type iSCSIUEventStartConnection struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	StartConn       struct {
		Sid uint32
		Cid uint32
		_   [16]byte
	}
	StartConnRet struct {
		Ret uint32
		_   [12]byte
	}
}

type iSCSIUEventStopConnection struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	StopConn        struct {
		Sid        uint32
		Cid        uint32
		ConnHandle uint64
		Flag       uint32
		_          [4]byte
	}
	StopConnRet struct {
		Ret uint32
		_   [12]byte
	}
}

type iSCSIUEventSendPDU struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	SendPDU         struct {
		Sid      uint32
		Cid      uint32
		HdrSize  uint32
		DataSize uint32
		_        [8]byte
	}
	SendPDURet struct {
		Ret int32
		_   [12]byte
	}
}

type iSCSIKEventConnError struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	_               [24]byte
	ConnErr         struct {
		Sid   uint32
		Cid   uint32
		Error IscsiErr
		_     [4]byte
	}
}

type iSCSIKEventRecvEvent struct {
	Type            IscsiEvent
	IfError         uint32
	TransportHandle uint64
	_               [24]byte
	RecvReq         struct {
		Sid        uint32
		Cid        uint32
		RecvHandle uint64
	}
}

// IscsiIpcConn is a single netlink connection
type IscsiIpcConn struct {
	Conn            *nl.NetlinkSocket
	TransportHandle uint64
	nextSeqNr       uint32
}

// ConnectNetlink connects to the iscsi netlink socket, and if successful returns
// an IscsiIpcConn ready to accept commands.
func ConnectNetlink() (*IscsiIpcConn, error) {
	conn, err := nl.Subscribe(unix.NETLINK_ISCSI, 1)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile("/sys/class/iscsi_transport/tcp/handle")
	if err != nil {
		return nil, err
	}
	handle, err := strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
	if err != nil {
		return nil, err
	}

	conn.SetReceiveTimeout(&unix.Timeval{Sec: 30})
	conn.SetSendTimeout(&unix.Timeval{Sec: 30})

	return &IscsiIpcConn{Conn: conn, TransportHandle: handle}, nil
}

// WaitFor reads ipcs until event of type Type is received, discarding others.
// (This presumes we're the only ones using netlink on this host, and we only
// have one outstanding request we're waiting for...)
func (c *IscsiIpcConn) WaitFor(Type IscsiEvent) (*syscall.NetlinkMessage, error) {
	for {
		msgs, _, err := c.Conn.Receive()

		if err != nil {
			return nil, err
		}

		for _, msg := range msgs {
			reader := bytes.NewReader(msg.Data)
			var uevent iSCSIUEvent
			err = binary.Read(reader, binary.LittleEndian, &uevent)
			if err != nil {
				return nil, err
			}
			if uevent.TransportHandle != c.TransportHandle {
				return nil, fmt.Errorf("wrong transport handle: %v", uevent.TransportHandle)
			}
			if uevent.Type == Type {
				return &msg, nil
			} else if uevent.Type == ISCSI_KEVENT_CONN_ERROR {
				reader.Seek(0, 0)
				var connErr iSCSIKEventConnError
				binary.Read(reader, binary.LittleEndian, &connErr)
				return nil, fmt.Errorf("connection error: %+v", connErr)
			} else if uevent.Type == ISCSI_KEVENT_IF_ERROR {
				return nil, fmt.Errorf("interface error: %v (invalid netlink message?)", uevent.IfError)
			}

			log.Printf("Dumping unexpected event of type %d", uevent.Type)
		}
	}
}

// FillNetlink aids attaching binary data to netlink request
func FillNetlink(request *nl.NetlinkRequest, data ...interface{}) error {
	var buf bytes.Buffer

	for _, item := range data {
		switch v := item.(type) {
		case []byte:
			_, err := buf.Write(v)
			if err != nil {
				return err
			}
		default:
			err := binary.Write(&buf, binary.LittleEndian, item)
			if err != nil {
				return err
			}
		}
	}
	request.AddRawData(buf.Bytes())

	return nil
}

// DoNetlink send netlink and listen to response
// ueventP *must* be a pointer to iSCSIUEvent
func (c *IscsiIpcConn) DoNetlink(ueventP unsafe.Pointer, data ...interface{}) error {
	uevent := (*iSCSIUEvent)(ueventP)
	request := nl.NetlinkRequest{
		NlMsghdr: unix.NlMsghdr{
			Len:   uint32(unix.SizeofNlMsghdr),
			Type:  uint16(uevent.Type),
			Flags: uint16(1),
			Seq:   atomic.AddUint32(&c.nextSeqNr, 1),
		},
	}

	if err := FillNetlink(&request, append([]interface{}{*uevent}, data...)...); err != nil {
		return err
	}
	if err := c.Conn.Send(&request); err != nil {
		return err
	}

	response, err := c.WaitFor(uevent.Type)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(response.Data)
	return binary.Read(reader, binary.LittleEndian, uevent)
}

// CreateSession creates a new kernel iSCSI session, returning the new
// iscsi_session id and scsi_host id
func (c *IscsiIpcConn) CreateSession(cmdsMax uint16, queueDepth uint16) (sid uint32, hostID uint32, err error) {
	cSession := iSCSIUEventCreateSession{
		Type:            ISCSI_UEVENT_CREATE_SESSION,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	cSession.CSession.CmdsMax = cmdsMax
	cSession.CSession.QueueDepth = queueDepth

	if err := c.DoNetlink(unsafe.Pointer(&cSession)); err != nil {
		return 0, 0, err
	}
	log.Println("Created new session ", cSession)

	return cSession.CSessionRet.Sid, cSession.CSessionRet.HostNo, nil
}

// DestroySession attempts to destroy the identified session. May fail if
// connections are still active
func (c *IscsiIpcConn) DestroySession(sid uint32) error {
	dSession := iSCSIUEventDestroySession{
		Type:            ISCSI_UEVENT_DESTROY_SESSION,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	dSession.DSession.Sid = sid

	if err := c.DoNetlink(unsafe.Pointer(&dSession)); err != nil {
		return err
	}

	if dSession.DSessionRet.Ret != 0 {
		return fmt.Errorf("error destroying session: %d", dSession.DSessionRet.Ret)
	}

	return nil
}

// CreateConnection creates a new iSCSI connection for an existing session,
// returning the connection id on success
func (c *IscsiIpcConn) CreateConnection(sid uint32) (cid uint32, err error) {

	cConn := iSCSIUEventCreateConnection{
		Type:            ISCSI_UEVENT_CREATE_CONN,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	cConn.CConn.Sid = sid

	if err := c.DoNetlink(unsafe.Pointer(&cConn)); err != nil {
		return 0, err
	}
	log.Println("Created new connection", cConn)

	return cConn.CConnRet.Cid, nil
}

// DestroyConnection attempts to destroy the identified connection. May fail if
// connection is still running
func (c *IscsiIpcConn) DestroyConnection(sid uint32, cid uint32) error {
	dConn := iSCSIUEventDestroyConnection{
		Type:            ISCSI_UEVENT_DESTROY_CONN,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	dConn.DConn.Sid = sid
	dConn.DConn.Cid = cid

	if err := c.DoNetlink(unsafe.Pointer(&dConn)); err != nil {
		return err
	}

	if dConn.DConnRet.Ret != 0 {
		return fmt.Errorf("error destroying connection: %d", dConn.DConnRet.Ret)
	}

	return nil
}

// BindConnection binds a TCP socket to the given kernel connection. fd must be
// the current program's fd for the TCP socket.
func (c *IscsiIpcConn) BindConnection(sid uint32, cid uint32, fd int) error {
	bConn := iSCSIUEventBindConnection{
		Type:            ISCSI_UEVENT_BIND_CONN,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	bConn.BConn.Sid = sid
	bConn.BConn.Cid = cid
	bConn.BConn.TransportEph = uint64(fd)
	bConn.BConn.IsLeading = 1

	err := c.DoNetlink(unsafe.Pointer(&bConn))
	if err != nil {
		return err
	}
	log.Println("Binded new connection ", bConn)

	if bConn.BConnRet.Ret != 0 {
		return fmt.Errorf("error binding connection: %d", bConn.BConnRet.Ret)
	}

	return nil
}

// SetParam sets a single parameter for the given connection. value will be
// null terminated and must not contain null bytes
func (c *IscsiIpcConn) SetParam(sid uint32, cid uint32, param IscsiParam, value string) error {
	setParam := iSCSIUEventSetParam{
		Type:            ISCSI_UEVENT_SET_PARAM,
		TransportHandle: c.TransportHandle,
	}
	setParam.SetParam.Sid = sid
	setParam.SetParam.Cid = cid
	setParam.SetParam.Param = param
	setParam.SetParam.Len = uint32(len(value) + 1) // Null terminatorrequest

	err := c.DoNetlink(unsafe.Pointer(&setParam), []byte(value+"\x00"))
	if err != nil {
		return err
	}

	if setParam.SetParamRet.Ret != 0 {
		return fmt.Errorf("error setting param %d: %d", param, setParam.SetParamRet.Ret)
	}
	return nil
}

// StartConnection starts the given connection. The connection should be bound
// and logged in.
func (c *IscsiIpcConn) StartConnection(sid uint32, cid uint32) error {
	sConn := iSCSIUEventStartConnection{
		Type:            ISCSI_UEVENT_START_CONN,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	sConn.StartConn.Sid = sid
	sConn.StartConn.Cid = cid

	err := c.DoNetlink(unsafe.Pointer(&sConn))
	if err != nil {
		return err
	}

	if sConn.StartConnRet.Ret != 0 {
		return fmt.Errorf("error starting connection: %d", sConn.StartConnRet.Ret)
	}
	return nil
}

// StopConnection attempts to stop the identified connection.
func (c *IscsiIpcConn) StopConnection(sid uint32, cid uint32) error {
	sConn := iSCSIUEventStopConnection{
		Type:            ISCSI_UEVENT_STOP_CONN,
		IfError:         0,
		TransportHandle: c.TransportHandle,
	}
	sConn.StopConn.Sid = sid
	sConn.StopConn.Cid = cid
	sConn.StopConn.Flag = STOP_CONN_RECOVER

	if err := c.DoNetlink(unsafe.Pointer(&sConn)); err != nil {
		return err
	}

	if sConn.StopConnRet.Ret != 0 {
		return fmt.Errorf("error stopping connection: %d", sConn.StopConnRet.Ret)
	}
	return nil
}

// PduLike interface for sending PDUs
type PduLike interface {
	// Length of the PDU header (iscsi_hdr/etc.)
	HeaderLen() uint32
	// Length of the PDU data
	DataLen() uint32
	// Header + Data
	Serialize() []byte
}

// RecvPDU waits for a PDU for the given connection, and returns the raw
// PDU with header on success. RecvPDU assumes a single iSCSI connection,
// and will error if a different connection receives a PDU
func (c *IscsiIpcConn) RecvPDU(sid uint32, cid uint32) ([]byte, error) {
	response, err := c.WaitFor(ISCSI_KEVENT_RECV_PDU)
	if err != nil {
		return nil, err
	}

	var recvReq iSCSIKEventRecvEvent
	reader := bytes.NewReader(response.Data)
	binary.Read(reader, binary.LittleEndian, &recvReq)

	if recvReq.RecvReq.Sid != sid || recvReq.RecvReq.Cid != cid {
		return nil, errors.New("unexpected PDU for different session... is another initiator running?")
	}

	log.Printf("Response length: %v, data length: %v", response.Header.Len, len(response.Data))

	return ioutil.ReadAll(reader)
}

// SendPDU sends the given PDU on the given connection
func (c *IscsiIpcConn) SendPDU(sid uint32, cid uint32, pdu PduLike) error {
	sendPdu := iSCSIUEventSendPDU{
		Type:            ISCSI_UEVENT_SEND_PDU,
		TransportHandle: c.TransportHandle,
	}
	sendPdu.SendPDU.Sid = sid
	sendPdu.SendPDU.Cid = cid
	sendPdu.SendPDU.HdrSize = pdu.HeaderLen()
	sendPdu.SendPDU.DataSize = pdu.DataLen()

	if err := c.DoNetlink(unsafe.Pointer(&sendPdu), pdu.Serialize()); err != nil {
		return err
	}

	if sendPdu.SendPDURet.Ret != 0 {
		return fmt.Errorf("SendPDU had unexpected error code %d", sendPdu.SendPDURet.Ret)
	}
	return nil
}
