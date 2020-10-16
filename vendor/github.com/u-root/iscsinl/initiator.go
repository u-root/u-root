// Package iscsinl acts as an initiator for bootstrapping an iscsi connection
// Partial implementation of RFC3720 login and NETLINK_ISCSI, just enough to
// get a connection going.
package iscsinl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// Login constants
const (
	ISCSI_OP_LOGIN     = 0x03
	ISCSI_OP_LOGIN_RSP = 0x23
	ISCSI_OP_IMMEDIATE = 0x40

	ISCSI_VERSION = 0x00

	ISCSI_FLAG_LOGIN_TRANSIT  = 0x80
	ISCSI_FLAG_LOGIN_CONTINUE = 0x40
)

// IscsiLoginStage corresponds to iSCSI login stage
type IscsiLoginStage uint8

// Login stages
const (
	ISCSI_SECURITY_NEGOTIATION_STAGE IscsiLoginStage = 0
	ISCSI_OP_PARMS_NEGOTIATION_STAGE                 = 1
	ISCSI_FULL_FEATURE_PHASE                         = 3
)

func hton24(buf *[3]byte, num int) {
	buf[0] = uint8(((num) >> 16) & 0xFF)
	buf[1] = uint8(((num) >> 8) & 0xFF)
	buf[2] = uint8((num) & 0xFF)
}

func ntoh24(buf [3]byte) uint {
	return (uint(buf[0]) << 16) | (uint(buf[1]) << 8) | uint(buf[2])
}

func hton48(buf *[6]byte, num int) {
	buf[0] = uint8(((num) >> 40) & 0xFF)
	buf[1] = uint8(((num) >> 32) & 0xFF)
	buf[2] = uint8(((num) >> 24) & 0xFF)
	buf[3] = uint8(((num) >> 16) & 0xFF)
	buf[4] = uint8(((num) >> 8) & 0xFF)
	buf[5] = uint8((num) & 0xFF)
}

// LoginHdr is the header for ISCSI_OP_LOGIN
// See: RFC3720 10.12.
type LoginHdr struct {
	Opcode     uint8
	Flags      uint8
	MaxVersion uint8
	MinVersion uint8
	HLength    uint8
	DLength    [3]uint8
	Isid       [6]uint8
	Tsih       uint16
	Itt        uint32
	Cid        uint16
	Rsvd3      uint16
	CmdSN      uint32
	ExpStatSN  uint32
	Rsvd5      [16]uint8
}

// LoginRspHdr is the header for ISCSI_OP_LOGIN_RSP
// See: RFC3720 10.13.
type LoginRspHdr struct {
	Opcode        uint8
	Flags         uint8
	MaxVersion    uint8
	ActiveVersion uint8
	HLength       uint8
	DLength       [3]uint8
	Isid          [6]uint8
	Tsih          uint16
	Itt           uint32
	Rsvd3         uint32
	StatSN        uint32
	ExpCmdSN      uint32
	MaxCmdSN      uint32
	StatusClass   uint8
	StatusDetail  uint8
	Rsvd5         [10]uint8
}

// IscsiLoginPdu is an iSCSI Login Request PDU
type IscsiLoginPdu struct {
	Header       LoginHdr
	TextSegments bytes.Buffer
}

// HeaderLen gives the length of the PDU header
func (l *IscsiLoginPdu) HeaderLen() uint32 {
	return uint32(binary.Size(l.Header))
}

// DataLen gives the length of all data segements for this PDU
func (l *IscsiLoginPdu) DataLen() uint32 {
	return uint32(l.TextSegments.Len())
}

// Serialize to network order bytes
func (l *IscsiLoginPdu) Serialize() []byte {
	var buf bytes.Buffer

	hton24(&l.Header.DLength, int(l.DataLen()))
	binary.Write(&buf, binary.LittleEndian, l.Header)
	buf.Write(l.TextSegments.Bytes())
	return buf.Bytes()
}

// AddParam the key=value string to the login payload and adds null terminator
func (l *IscsiLoginPdu) AddParam(keyvalue string) {
	l.TextSegments.WriteString(keyvalue)
	l.TextSegments.WriteByte(0)
}

// ReReadPartitionTable opens the given file and reads partition table from it
func ReReadPartitionTable(devname string) error {
	f, err := os.OpenFile(devname, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	_, err = unix.IoctlGetInt(int(f.Fd()), unix.BLKRRPART)
	return err
}

// IscsiOptions configures iSCSI session.
type IscsiOptions struct {
	InitiatorName string
	Address       string
	Volume        string

	// See RFC7143 Section 13 for these.
	// Max data per single incoming iSCSI packet
	MaxRecvDLength int
	// Max data per single outgoing iSCSI packet
	MaxXmitDLength int
	// Max unsolicited data per iSCSI command sequence
	FirstBurstLength int
	// Max data per iSCSI command sequence
	MaxBurstLength int
	// CRC32C or None
	HeaderDigest string
	// CRC32C or None
	DataDigest string
	// Seconds to wait for heartbeat response before declaring the connection dead
	PingTimeout int32
	// Seconds to wait on an idle connection before sending a heartbeat
	RecvTimeout int32
	// Max iSCSI commands outstanding
	CmdsMax uint16
	// Max IOs outstanding
	QueueDepth uint16
	// Require initial Ready To Transfer (R2T) (false enables unsolicited data)
	InitialR2T bool
	// Enable iSCSI Immediate Data
	ImmediateData bool
	// Require iSCSI data sequence to be sent order by offset
	DataPDUInOrder bool
	// Require packets in an iSCSI data sequence to be sent in order by sequence number
	DataSequenceInOrder bool

	// Scheduler to configure for the blockdev
	Scheduler string

	// ScanTimeout is the total time to wait for block devices to appear in
	// the file system after initiating a scan.
	ScanTimeout time.Duration
}

// IscsiTargetSession represents an iSCSI session and a single connection to a target
type IscsiTargetSession struct {
	opts   IscsiOptions
	cid    uint32
	hostID uint32
	sid    uint32

	// Update this on login response
	tsih      uint16
	expCmdSN  uint32
	maxCmdSN  uint32
	expStatSN uint32
	currStage IscsiLoginStage

	// Seconds to wait for heartbeat response before declaring the connection dead
	pingTimeout int32
	// Seconds to wait on an idle connection before sending a heartbeat
	recvTimeout int32

	blockDevName []string

	conn    *net.TCPConn
	netlink *IscsiIpcConn
}

const (
	oneMegabyte = 1048576
	oneMinute   = 60
)

var defaultOpts = IscsiOptions{
	MaxRecvDLength:      oneMegabyte,
	MaxXmitDLength:      oneMegabyte,
	FirstBurstLength:    oneMegabyte,
	MaxBurstLength:      oneMegabyte,
	HeaderDigest:        "CRC32C",
	DataDigest:          "CRC32C",
	PingTimeout:         oneMinute,
	RecvTimeout:         oneMinute,
	CmdsMax:             128,
	QueueDepth:          16,
	InitialR2T:          false,
	ImmediateData:       true,
	DataPDUInOrder:      true,
	DataSequenceInOrder: true,
	Scheduler:           "noop",
	ScanTimeout:         3 * time.Second,
}

// Option is a functional API for setting optional configuration.
type Option func(i *IscsiOptions)

// WithTarget adds the target address and volume to the config.
func WithTarget(addr, volume string) Option {
	return func(i *IscsiOptions) {
		i.Address = addr
		i.Volume = volume
	}
}

// WithScanTimeout sets the timeout to wait for devices to appear after sending the scan event.
func WithScanTimeout(dur time.Duration) Option {
	return func(i *IscsiOptions) {
		i.ScanTimeout = dur
	}
}

// WithInitiator adds the initiator name to the config.
func WithInitiator(initiatorName string) Option {
	return func(i *IscsiOptions) {
		i.InitiatorName = initiatorName
	}
}

// WithCmdsMax sets the maximum number of outstanding iSCSI commands.
func WithCmdsMax(n uint16) Option {
	return func(i *IscsiOptions) {
		i.CmdsMax = n
	}
}

// WithQueueDepth sets the maximum number of outstanding IOs.
func WithQueueDepth(n uint16) Option {
	return func(i *IscsiOptions) {
		i.QueueDepth = n
	}
}

// WithScheduler sets the block device scheduler.
func WithScheduler(sched string) Option {
	return func(i *IscsiOptions) {
		i.Scheduler = sched
	}
}

// WithDigests sets both the header and data digest. Acceptable values: None or CRC32C.
func WithDigests(digest string) Option {
	return func(i *IscsiOptions) {
		i.HeaderDigest = digest
		i.DataDigest = digest
	}
}

// NewSession constructs an IscsiTargetSession
func NewSession(netlink *IscsiIpcConn, opts ...Option) *IscsiTargetSession {
	i := &IscsiTargetSession{
		opts:    defaultOpts,
		netlink: netlink,
	}
	// Apply optional arguments from user.
	for _, opt := range opts {
		opt(&i.opts)
	}
	return i
}

// Connect creates a kernel iSCSI session and connection, connects to the
// target, and binds the connection to the kernel session.
func (s *IscsiTargetSession) Connect() error {
	var err error
	s.sid, s.hostID, err = s.netlink.CreateSession(s.opts.CmdsMax, s.opts.QueueDepth)
	if err != nil {
		return err
	}

	s.cid, err = s.netlink.CreateConnection(s.sid)
	if err != nil {
		return err
	}

	resolvedAddr, err := net.ResolveTCPAddr("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.conn, err = net.DialTCP("tcp", nil, resolvedAddr)
	if err != nil {
		return err
	}

	file, err := s.conn.File()
	if err != nil {
		return err
	}
	defer file.Close()
	fd := file.Fd()

	return s.netlink.BindConnection(s.sid, s.cid, int(fd))
}

// Start starts the kernel iSCSI session. Call this after successfully
// logging in and setting all desired parameters.
func (s *IscsiTargetSession) Start() error {
	return s.netlink.StartConnection(s.sid, s.cid)
}

// TearDown stops and destroys the connection & session
// in case of partially created session, stopping connections/destroying
// connections won't work, so try it all
func (s *IscsiTargetSession) TearDown() error {
	sConnErr := s.netlink.StopConnection(s.sid, s.cid)

	dConnErr := s.netlink.DestroyConnection(s.sid, s.cid)

	if err := s.netlink.DestroySession(s.sid); err != nil {
		return fmt.Errorf("failure to destroy session DestroySession:%v DestroyConnection:%v StopConnection:%v", err, dConnErr, sConnErr)
	}
	return nil
}

func netlinkBoolStr(pred bool) string {
	if pred {
		return "1"
	}
	return "0"
}

func iscsiParseBool(inval string) (bool, error) {
	if inval == "Yes" {
		return true, nil
	} else if inval == "No" {
		return false, nil
	}
	return false, fmt.Errorf("invalid bool: %s", inval)
}

func iscsiBoolStr(pred bool) string {
	if pred {
		return "Yes"
	}
	return "No"
}

// SetParams sets some desired parameters for the kernel session
func (s *IscsiTargetSession) SetParams() error {
	params := []struct {
		p IscsiParam
		v string
	}{
		{ISCSI_PARAM_TARGET_NAME, s.opts.Volume},
		{ISCSI_PARAM_INITIATOR_NAME, s.opts.InitiatorName},
		{ISCSI_PARAM_MAX_RECV_DLENGTH, fmt.Sprintf("%d", s.opts.MaxRecvDLength)},
		{ISCSI_PARAM_MAX_XMIT_DLENGTH, fmt.Sprintf("%d", s.opts.MaxXmitDLength)},
		{ISCSI_PARAM_FIRST_BURST, fmt.Sprintf("%d", s.opts.FirstBurstLength)},
		{ISCSI_PARAM_MAX_BURST, fmt.Sprintf("%d", s.opts.MaxBurstLength)},
		{ISCSI_PARAM_PDU_INORDER_EN, netlinkBoolStr(s.opts.DataPDUInOrder)},
		{ISCSI_PARAM_DATASEQ_INORDER_EN, netlinkBoolStr(s.opts.DataSequenceInOrder)},
		{ISCSI_PARAM_INITIAL_R2T_EN, netlinkBoolStr(s.opts.InitialR2T)},
		{ISCSI_PARAM_IMM_DATA_EN, netlinkBoolStr(s.opts.ImmediateData)},
		{ISCSI_PARAM_EXP_STATSN, fmt.Sprintf("%d", s.expStatSN)},
		{ISCSI_PARAM_HDRDGST_EN, netlinkBoolStr(s.opts.HeaderDigest == "CRC32C")},
		{ISCSI_PARAM_DATADGST_EN, netlinkBoolStr(s.opts.DataDigest == "CRC32C")},
		{ISCSI_PARAM_PING_TMO, fmt.Sprintf("%d", s.opts.PingTimeout)},
		{ISCSI_PARAM_RECV_TMO, fmt.Sprintf("%d", s.opts.RecvTimeout)},
	}

	for _, pp := range params {
		log.Printf("Setting param %s to %v", pp.p.String(), pp.v)
		if err := s.netlink.SetParam(s.sid, s.cid, pp.p, pp.v); err != nil {
			return err
		}
	}
	return nil
}

// writeFile is ioutil.WriteFile but disallows creating new file
func writeFile(filename string, contents string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	wlen, err := file.WriteString(contents)
	if err != nil && wlen < len(contents) {
		err = io.ErrShortWrite
	}
	// If Close() fails this likely indicates a write failure.
	if errClose := file.Close(); err == nil {
		err = errClose
	}
	return err
}

// ReScan triggers a scsi host scan so the kernel creates a block device for the
// newly attached session, then waits for the block device to be created
func (s *IscsiTargetSession) ReScan() error {
	// The three wildcards stand for channel, SCSI target ID, and LUN.
	if err := writeFile(fmt.Sprintf("/sys/class/scsi_host/host%d/scan", s.hostID), "- - -"); err != nil {
		return err
	}

	var matches []string
	start := time.Now()
	// The kernel may add devices it finds through scanning at any time. If
	// a scan yields multiple, kernel will not add them atomically. We wait
	// until at least one device has appeared, and no new devices have
	// appeared for 100ms. We also time out based on the user defined
	// ScanTimeout.
	for elapsed := time.Now().Sub(start); elapsed <= s.opts.ScanTimeout; {
		log.Printf("Waiting for device...")
		time.Sleep(100 * time.Millisecond)
		newMatches, err := filepath.Glob(fmt.Sprintf(
			"/sys/class/iscsi_session/session%d/device/target*/*/block/*/uevent", s.sid))
		if err != nil {
			return err
		}
		if len(newMatches) > 0 {
			if len(matches) == len(newMatches) {
				break
			}
			matches = newMatches
		}
	}

	found := false
	for _, match := range matches {
		contents, err := ioutil.ReadFile(match)
		if err != nil {
			log.Printf("error reading file for %v err=%v skipping error\n", match, err)
			continue
		}

		for _, kv := range strings.Split(string(contents), "\n") {
			splitkv := strings.Split(kv, "=")
			if splitkv[0] == "DEVNAME" {
				s.blockDevName = append(s.blockDevName, splitkv[1])
				found = true
			}
		}
	}

	if !found {
		return errors.New("could not find any device DEVNAMEs")
	}
	return nil

}

// ConfigureBlockDevs will set blockdev params for this iSCSI session, and returns blockdev name
func (s *IscsiTargetSession) ConfigureBlockDevs() ([]string, error) {
	if err := s.ReScan(); err != nil {
		return nil, err
	}

	for i := range s.blockDevName {
		for {
			log.Printf("Waiting for sysfs...")
			time.Sleep(30 * time.Millisecond)
			_, err := os.Stat(fmt.Sprintf("/sys/block/%v/queue/nr_requests", s.blockDevName[i]))
			if !os.IsNotExist(err) {
				break
			}
		}
		params := []struct {
			filen string
			val   string
		}{
			{fmt.Sprintf("/sys/block/%v/queue/nr_requests", s.blockDevName[i]), fmt.Sprintf("%d", s.opts.QueueDepth)},
			{fmt.Sprintf("/sys/block/%v/queue/scheduler", s.blockDevName[i]), s.opts.Scheduler},
			{fmt.Sprintf("/sys/block/%v/queue/rotational", s.blockDevName[i]), "0"},
		}

		for _, pp := range params {
			if err := writeFile(pp.filen, pp.val); err != nil {
				return nil, err
			}
		}
	}
	return s.blockDevName, nil
}

// processOperationalParam assigns params returned from the target. Errors if
// we cannot continue with negotiation.
func (s *IscsiTargetSession) processOperationalParam(keyvalue string) error {
	split := strings.Split(keyvalue, "=")
	if len(split) != 2 {
		return fmt.Errorf("invalid format for operational param \"%v\"", keyvalue)
	}
	key, value := split[0], split[1]

	if value == "Reject" {
		return fmt.Errorf("target rejected parameter %q", key)
	}

	switch key {
	case "HeaderDigest":
		s.opts.HeaderDigest = value
	case "DataDigest":
		s.opts.DataDigest = value
	case "InitialR2T":
		val, err := iscsiParseBool(value)
		if err != nil {
			return err
		}
		s.opts.InitialR2T = val || s.opts.InitialR2T
	case "ImmediateData":
		val, err := iscsiParseBool(value)
		if err != nil {
			return err
		}
		s.opts.ImmediateData = val && s.opts.ImmediateData
	case "MaxRecvDataSegmentLength":
		length, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		s.opts.MaxXmitDLength = int(length)
	case "MaxBurstLength":
		length, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		s.opts.MaxBurstLength = int(length)
	case "FirstBurstLength":
		length, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		s.opts.FirstBurstLength = int(length)
	case "DataPDUInOrder":
		val, err := iscsiParseBool(value)
		if err != nil {
			return err
		}
		s.opts.DataPDUInOrder = val || s.opts.DataPDUInOrder
	case "DataSequenceInOrder":
		val, err := iscsiParseBool(value)
		if err != nil {
			return err
		}
		s.opts.DataSequenceInOrder = val || s.opts.DataSequenceInOrder
	default:
		log.Printf("Ignoring unknown param \"%v\"", keyvalue)
	}
	return nil
}

// processOperationalParams processes all parameters in a login response
func (s *IscsiTargetSession) processOperationalParams(data []byte) error {
	params := strings.Split(string(data), "\x00")
	// Annoyingly, strings.Split will always have an empty string at the end
	// An empty string in the middle of params suggests we have an otherwise
	// malformed request, since we shouldn't expect double nul bytes
	params = params[0 : len(params)-1]
	for _, param := range params {
		if err := s.processOperationalParam(param); err != nil {
			return err
		}
	}
	return nil
}

func (s *IscsiTargetSession) processLoginResponse(response []byte) error {
	var loginRespPdu LoginRspHdr
	reader := bytes.NewReader(response)
	if err := binary.Read(reader, binary.LittleEndian, &loginRespPdu); err != nil {
		return err
	}
	if loginRespPdu.Opcode != ISCSI_OP_LOGIN_RSP {
		return fmt.Errorf("unexpected response pdu opcode %d", loginRespPdu.Opcode)
	}

	if loginRespPdu.StatusClass != 0 {
		return fmt.Errorf("error in login response %d %d", loginRespPdu.StatusClass, loginRespPdu.StatusDetail)
	}

	s.maxCmdSN = loginRespPdu.MaxCmdSN
	s.expCmdSN = loginRespPdu.ExpCmdSN
	s.tsih = loginRespPdu.Tsih
	s.expStatSN = loginRespPdu.StatSN + 1
	if (loginRespPdu.Flags & ISCSI_FLAG_LOGIN_TRANSIT) != 0 {
		s.currStage = IscsiLoginStage(loginRespPdu.Flags & 0x03)
	}

	// dLength generally != the length of the rest of the netlink buffer
	dLength := int(ntoh24(loginRespPdu.DLength))
	if dLength == 0 {
		return nil
	}
	theRest := make([]byte, dLength)
	read, err := reader.Read(theRest)
	if err != nil {
		return err
	}
	if read != dLength {
		return errors.New("unexpected EOF reading PDU data")
	}
	return s.processOperationalParams(theRest)
}

// Login - RFC iSCSI login
// https://www.ietf.org/rfc/rfc3720.txt
// For now "negotiates" no auth security.
func (s *IscsiTargetSession) Login() error {
	log.Println("Starting login...")

	for s.currStage != ISCSI_OP_PARMS_NEGOTIATION_STAGE {
		loginReq := IscsiLoginPdu{
			Header: LoginHdr{
				Opcode:     ISCSI_OP_LOGIN | ISCSI_OP_IMMEDIATE,
				MaxVersion: ISCSI_VERSION,
				MinVersion: ISCSI_VERSION,
				ExpStatSN:  s.expStatSN,
				Tsih:       s.tsih,
				Flags:      uint8((s.currStage << 2) | ISCSI_OP_PARMS_NEGOTIATION_STAGE | ISCSI_FLAG_LOGIN_TRANSIT),
			},
		}
		hton48(&loginReq.Header.Isid, int(s.sid))
		loginReq.AddParam("AuthMethod=None")
		// RFC 3720 page 36 last line, https://tools.ietf.org/html/rfc3720#page-36
		// The session type is defined during login with the key=value parameter
		// in the login command.
		loginReq.AddParam("SessionType=Normal")
		loginReq.AddParam(fmt.Sprintf("InitiatorName=%s", s.opts.InitiatorName))
		loginReq.AddParam(fmt.Sprintf("TargetName=%s", s.opts.Volume))

		if err := s.netlink.SendPDU(s.sid, s.cid, &loginReq); err != nil {
			return fmt.Errorf("sendPDU: %v", err)
		}

		response, err := s.netlink.RecvPDU(s.sid, s.cid)
		if err != nil {
			return fmt.Errorf("recvpdu: %v", err)
		}
		if err = s.processLoginResponse(response); err != nil {
			return err
		}
	}

	for s.currStage != ISCSI_FULL_FEATURE_PHASE {
		loginReq := IscsiLoginPdu{
			Header: LoginHdr{
				Opcode:     ISCSI_OP_LOGIN | ISCSI_OP_IMMEDIATE,
				MaxVersion: ISCSI_VERSION,
				MinVersion: ISCSI_VERSION,
				ExpStatSN:  s.expStatSN,
				Tsih:       s.tsih,
				Flags:      uint8((s.currStage << 2) | ISCSI_FULL_FEATURE_PHASE | ISCSI_FLAG_LOGIN_TRANSIT),
			},
		}
		hton48(&loginReq.Header.Isid, int(s.sid))
		loginReq.AddParam(fmt.Sprintf("InitiatorName=%s", s.opts.InitiatorName))
		loginReq.AddParam(fmt.Sprintf("TargetName=%s", s.opts.Volume))
		loginReq.AddParam("SessionType=Normal")
		loginReq.AddParam(fmt.Sprintf("MaxRecvDataSegmentLength=%d", s.opts.MaxRecvDLength))
		loginReq.AddParam(fmt.Sprintf("FirstBurstLength=%d", s.opts.FirstBurstLength))
		loginReq.AddParam(fmt.Sprintf("MaxBurstLength=%d", s.opts.MaxBurstLength))
		loginReq.AddParam(fmt.Sprintf("HeaderDigest=%v", s.opts.HeaderDigest))
		loginReq.AddParam(fmt.Sprintf("DataDigest=%v", s.opts.DataDigest))
		loginReq.AddParam(fmt.Sprintf("InitialR2T=%v", iscsiBoolStr(s.opts.InitialR2T)))
		loginReq.AddParam(fmt.Sprintf("ImmediateData=%v", iscsiBoolStr(s.opts.ImmediateData)))
		loginReq.AddParam(fmt.Sprintf("DataPDUInOrder=%v", iscsiBoolStr(s.opts.DataPDUInOrder)))
		loginReq.AddParam(fmt.Sprintf("DataSequenceInOrder=%v", iscsiBoolStr(s.opts.DataSequenceInOrder)))

		if err := s.netlink.SendPDU(s.sid, s.cid, &loginReq); err != nil {
			return fmt.Errorf("sendpdu2: %v", err)
		}

		response, err := s.netlink.RecvPDU(s.sid, s.cid)
		if err != nil {
			return fmt.Errorf("recvpdu2: %v", err)
		}
		if err = s.processLoginResponse(response); err != nil {
			return err
		}
	}
	return nil

}

// MountIscsi connects to the given iscsi target and mounts it, returning the
// device name on success
func MountIscsi(opts ...Option) ([]string, error) {
	netlink, err := ConnectNetlink()
	if err != nil {
		return nil, err
	}

	session := NewSession(netlink, opts...)
	if err = session.Connect(); err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	if err := session.Login(); err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	if err := session.SetParams(); err != nil {
		return nil, fmt.Errorf("params: %v", err)
	}

	if err := session.Start(); err != nil {
		return nil, fmt.Errorf("start: %v", err)
	}

	devnames, err := session.ConfigureBlockDevs()
	if err != nil {
		return nil, err
	}

	for i := range devnames {
		if err := ReReadPartitionTable("/dev/" + devnames[i]); err != nil {
			return nil, err
		}
	}

	return devnames, nil
}

// TearDownIscsi tears down the specified session
func TearDownIscsi(sid uint32, cid uint32) error {
	netlink, err := ConnectNetlink()
	if err != nil {
		return err
	}
	session := IscsiTargetSession{sid: sid, cid: cid, netlink: netlink}

	return session.TearDown()
}
