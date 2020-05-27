// Copyright (C) 2017 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package tftp // import "pack.ag/tftp"

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"pack.ag/tftp/netascii"
)

const (
	defaultPort       = "69"
	defaultMode       = ModeOctet
	defaultUDPNet     = "udp"
	defaultTimeout    = time.Second
	defaultBlksize    = 512
	defaultWindowsize = 1
	defaultRetransmit = 10
)

// All connections will use these options unless overridden.
var defaultOptions = map[string]string{
	optTransferSize: "0", // Enable tsize
}

// newConn starts listening on a system assigned port and returns an initialized conn
//
// udpNet is one of "udp", "udp4", or "udp6"
// addr is the address of the target client or server
func newConn(udpNet string, mode TransferMode, addr *net.UDPAddr) (*conn, error) {
	// Start listening, an empty UDPAddr will cause the system to assign a port
	netConn, err := net.ListenUDP(udpNet, &net.UDPAddr{})
	if err != nil {
		return nil, wrapError(err, "network listen failed")
	}

	c := &conn{
		log:        newLogger(addr.String()),
		remoteAddr: addr,
		netConn:    netConn,
		blksize:    defaultBlksize,
		timeout:    defaultTimeout,
		windowsize: defaultWindowsize,
		retransmit: defaultRetransmit,
		mode:       mode,
	}
	c.rx.buf = make([]byte, 4+defaultBlksize) // +4 for headers

	return c, nil
}

func newSinglePortConn(addr *net.UDPAddr, mode TransferMode, netConn *net.UDPConn, reqChan chan []byte) *conn {
	return &conn{
		log:        newLogger(addr.String()),
		remoteAddr: addr,
		blksize:    defaultBlksize,
		timeout:    defaultTimeout,
		windowsize: defaultWindowsize,
		retransmit: defaultRetransmit,
		mode:       mode,
		buf:        make([]byte, 4+defaultBlksize), // +4 for headers
		reqChan:    reqChan,
		netConn:    netConn,
	}
}

// newConnFromHost wraps newConn and looks up the target's address from a string
//
// This function is used by Client
func newConnFromHost(udpNet string, mode TransferMode, host string) (*conn, error) {
	// Resolve server
	addr, err := net.ResolveUDPAddr(udpNet, host)
	if err != nil {
		return nil, wrapError(err, "address resolve failed")
	}

	return newConn(udpNet, mode, addr)
}

// conn handles TFTP read and write requests
type conn struct {
	log        *logger
	netConn    *net.UDPConn // Underlying network connection
	remoteAddr net.Addr     // Address of the remote server or client

	// Single Port Mode
	reqChan chan []byte
	timer   *time.Timer

	// Transfer type
	isClient bool // Whether or not we're the client, gets set by sendRequest
	isSender bool // Whether we're sending or receiving, gets set by writeSetup

	// Negotiable options
	blksize    uint16        // Size of DATA payloads
	timeout    time.Duration // How long to wait before resending packets
	windowsize uint16        // Number of DATA packets between ACKs
	mode       TransferMode  // octet or netascii
	tsize      *int64        // Size of the file being sent/received

	// Other, non-negotiable options
	retransmit int // Number of times an individual datagram will be retransmitted on error

	// Track state of transfer
	optionsParsed bool   // Whether TFTP options have been parsed yet
	window        uint16 // Packets sent since last ACK
	block         uint16 // Current block #
	catchup       bool   // Ignore incoming blocks from a window we reset
	p             []byte // bytes to be read/written (depending on send/receive)
	n             int    // byte count read/written
	tries         int    // retry counter
	err           error  // error has occurreds
	closing       bool   // connection is closing
	done          bool   // the transfer is complete

	// Buffers
	buf   []byte       // incoming data from, sized to blksize + headers
	txBuf *ringBuffer  // buffers outgoing data, retaining windowsize * blksize
	rxBuf bytes.Buffer // buffer incoming data

	// Datgrams
	tx datagram // Constructs outgoing datagrams
	rx datagram // Hold and parse current incoming datagram

	// reader/writer are rxBuf/txBuf, possibly wrapped by netascii reader/writer
	reader io.Reader
	writer io.Writer
}

// sendWriteRequest sends WRQ to server and negotiates transfer options
func (c *conn) sendWriteRequest(filename string, opts map[string]string) error {
	c.isSender = true
	// Build WRQ
	c.tx.writeWriteReq(filename, c.mode, opts)

	for state := c.sendRequest; state != nil; {
		state = state()
	}

	return c.err
}

// sendReadRequest send RRQ to server and negotiates transfer options
//
// If the server doesn't support options and responds with data, the data will be added
// to rxBuf.
func (c *conn) sendReadRequest(filename string, opts map[string]string) error {
	// Build RRQ
	c.tx.writeReadReq(filename, c.mode, opts)

	for state := c.sendRequest; state != nil; {
		state = state()
	}

	return c.err
}

func (c *conn) sendRequest() stateType {
	// Set that we're a client
	c.isClient = true

	// Send request
	if err := c.writeToNet(); err != nil {
		c.err = wrapError(err, "writing request to network")
		return nil
	}

	return c.receiveResponse
}

func (c *conn) receiveResponse() stateType {
	if c.tries >= c.retransmit {
		c.err = wrapError(c.err, "receiving request response")
		return nil
	}
	c.tries++

	addr, err := c.readFromNet()
	if err != nil {
		c.log.debug("error getting %s response from %v", c.tx.opcode(), c.remoteAddr)
		c.err = err
		return c.receiveResponse
	}

	if err := c.rx.validate(); err != nil {
		c.log.debug("error validating response from %v: %v", c.remoteAddr, err)
		c.err = wrapError(err, "validating request response")
		return nil
	}

	if c.reqChan == nil {
		// Update address
		c.remoteAddr = addr
	}
	c.log.trace("Received response from %v: %v", addr, c.rx)

	c.tries = 0

	if c.isSender {
		return c.handleWRQResponse
	}

	return c.handleRRQResponse
}

func (c *conn) handleWRQResponse() stateType {
	// Should have received OACK if server supports options, or ACK if not
	switch c.rx.opcode() {
	case opCodeOACK, opCodeACK:
		// Got OACK, parse options
		return c.writeSetup
	case opCodeERROR:
		// Received an error
		c.err = wrapError(c.remoteError(), "WRQ OACK response")
		return nil
	default:
		c.err = wrapError(&errUnexpectedDatagram{dg: c.rx.String()}, "WRQ OACK response")
		return nil
	}
}

func (c *conn) handleRRQResponse() stateType {
	// Should have received OACK if server supports options, or DATA if not
	switch c.rx.opcode() {
	case opCodeOACK:
		// Got OACK, parse options
		return c.readSetup
	case opCodeDATA:
		// Server doesn't support options,
		// write data to the buf so it's available for reading
		n, err := c.rxBuf.Write(c.rx.data())
		if err != nil {
			c.err = wrapError(err, "writing RRQ response data")
			return nil
		}
		c.block = c.rx.block()
		if uint16(n) < c.blksize {
			c.done = true
		}
		return c.readSetup
	case opCodeERROR:
		// Received an error
		c.err = wrapError(c.remoteError(), "RRQ OACK response")
		return nil
	default:
		c.err = wrapError(&errUnexpectedDatagram{dg: c.rx.String()}, "RRQ OACK response")
		return nil
	}
}

// Write implements io.Writer and wraps write().
//
// If mode is ModeNetASCII, wrap write() with netascii.EncodeWriter.
func (c *conn) Write(p []byte) (int, error) {
	// Can't write if an error has been sent/received
	if c.err != nil {
		return 0, wrapError(c.err, "checking conn err before Write")
	}

	c.p = p
	for state := c.startWrite; state != nil; {
		state = state()
	}

	return c.n, wrapError(c.err, "writing")
}

type stateType func() stateType

func (c *conn) startWrite() stateType {
	if !c.optionsParsed {
		// Options won't be parsed before first write so that API consumer
		// has opportunity to set tsize with ReadRequest.WriteSize()
		return c.writeSetup
	}
	return c.write
}

// writeSetup parses options and sets up buffers before
// first write.
func (c *conn) writeSetup() stateType {
	// Set that we're sending
	c.isSender = true

	ackOpts, err := c.parseOptions()
	if err != nil {
		return c.error(err, "parsing options")
	}

	// Set buf size
	if len(c.buf) != int(c.blksize) {
		c.buf = make([]byte, c.blksize)
	}

	// Init ringBuffer
	c.txBuf = newRingBuffer(int(c.windowsize), int(c.blksize))

	c.writer = c.txBuf
	if c.mode == ModeNetASCII {
		c.writer = netascii.NewWriter(c.writer)
	}

	// Client setup is done, ready to send data
	if c.isClient {
		return nil
	}

	// Sending DATA ACKs when there are no options
	if len(ackOpts) == 0 {
		return c.write
	}

	// Send OACK
	return c.sendOACK(ackOpts)
}

func (c *conn) sendOACK(o options) stateType {
	return func() stateType {
		c.log.trace("Sending OACK to %s\n", c.remoteAddr)
		c.tx.writeOptionAck(o)
		if err := c.writeToNet(); err != nil {
			return c.error(err, "writing OACK")
		}

		return c.getAck
	}
}

func (c *conn) error(err error, desc string) stateType {
	return func() stateType {
		c.err = wrapError(err, desc)
		return nil
	}
}

// write writes adds data to txBuf and writes data to netConn in chunks of
// blksize, until the last chunk of <blksize, which signals transfer completion.
func (c *conn) write() stateType {
	// Copy to buffer
	read, err := c.writer.Write(c.p)
	if err != nil {
		c.n, c.err = read, wrapError(err, "writing data to txBuf before write")
		return nil
	}
	c.n = read

	return c.writeData
}

// writeData writes a single DATA datagram
func (c *conn) writeData() stateType {
	if c.closing && c.done {
		return nil
	}
	if c.txBuf.Len() < int(c.blksize) && !c.closing {
		return nil
	}

	c.block++

	// Read data from txBuf
	n, err := c.txBuf.Read(c.buf)
	if err != nil && err != io.EOF {
		c.err = wrapError(err, "reading data from txBuf before writing to network")
		return nil
	}
	c.tx.writeData(c.block, c.buf[:n])

	// Send w.tx datagram
	c.log.trace("Sending block %d with %d bytes to %s\n", c.block, n, c.remoteAddr)
	err = c.writeToNet()
	if err != nil {
		c.err = wrapError(err, "writing data to network")
		return nil
	}

	// Increment the window
	c.window++

	// If this is last block, move to get ack immediately
	if uint16(n) < c.blksize {
		c.done = true
		return c.getAck
	}

	// Continue on if we haven't reached the windowsize
	if c.window < c.windowsize {
		return c.writeData
	}

	// Reset window
	c.window = 0

	return c.getAck
}

// Read implements io.Reader and wraps read()
//
// If mode is ModeNetASCII, read() is wrapped with netascii.ReadDecoder
func (c *conn) Read(p []byte) (int, error) {
	c.n = 0
	if c.err != nil {
		// Can't read if an error has been sent/received
		return 0, wrapError(c.err, "checking conn error before Read")
	}

	c.p = p
	for state := c.startRead; state != nil; {
		state = state()
	}
	return c.n, c.err
}

func (c *conn) startRead() stateType {
	if !c.optionsParsed {
		return c.readSetup
	}
	return c.read
}

// readSetup parses options and sets up buffers before
// first read.
func (c *conn) readSetup() stateType {
	c.reader = &c.rxBuf
	if c.mode == ModeNetASCII {
		c.reader = netascii.NewReader(c.reader)
	}

	ackOpts, err := c.parseOptions()
	if err != nil {
		c.err = wrapError(err, "read setup")
		return nil
	}

	// Set buf size
	if needed := int(c.blksize + 4); len(c.rx.buf) != needed {
		c.rx.buf = make([]byte, needed)
	}

	// If there we're not options negotiated, send ACK
	// Client never sends OACK
	if len(ackOpts) == 0 || c.isClient {
		c.log.trace("Sending ACK to %s\n", c.remoteAddr)
		c.tx.writeAck(c.block)
	} else {
		c.log.trace("Sending OACK to %s\n", c.remoteAddr)
		c.tx.writeOptionAck(ackOpts)
	}

	// Send ACK/OACK
	err = c.writeToNet()
	if err != nil {
		c.err = err
		return nil
	}

	if c.isClient {
		return nil
	}

	return c.read
}

// read reads data from netConn until p is full or the connection is
// complete.
func (c *conn) read() stateType {
	if c.rxBuf.Len() >= len(c.p) || c.done {
		// Read buffered data into p
		n, err := c.reader.Read(c.p)
		c.n = n
		if err != nil && err != io.EOF { // Ignore EOF from bytes.Buffer
			c.err = wrapError(err, "reading from rxBuf after read")
		}
		// If done, signal that there's nothing more to read by io.EOF
		if c.done && c.rxBuf.Len() == 0 {
			c.err = io.EOF
		}
		return nil
	}

	// Read next datagram
	return c.readData
}

// readDatagram reads a single datagram into rx
func (c *conn) readData() stateType {
	if c.tries >= c.retransmit {
		c.log.debug("Max retries exceeded")
		c.sendError(ErrCodeNotDefined, "max retries reached")
		c.err = wrapError(ErrMaxRetries, "reading data")
		return nil
	}
	c.tries++

	c.log.trace("Waiting for DATA from %s\n", c.remoteAddr)
	_, err := c.readFromNet()
	if err != nil {
		c.log.debug("error receiving block %d: %v", c.block+1, err)
		c.log.trace("Resending ACK for %d\n", c.block)
		if err := c.sendAck(c.block); err != nil {
			c.log.debug("resending ACK %v", err)
		}
		c.window = 0
		return c.readData
	}

	// validate datagram
	if err := c.rx.validate(); err != nil {
		c.err = wrapError(err, "validating read data")
		return nil
	}

	// Check for opcode
	switch op := c.rx.opcode(); op {
	case opCodeDATA:
	case opCodeERROR:
		// Received an error
		c.err = wrapError(c.remoteError(), "reading data")
		return nil
	default:
		c.err = wrapError(&errUnexpectedDatagram{dg: c.rx.String()}, "read data response")
		return nil
	}

	c.log.trace("Received block %d\n", c.rx.block())
	c.tries = 0

	return c.ackData
}

// ackData handles block sequence, windowing, and acknowledgements
func (c *conn) ackData() stateType {
	switch diff := c.rx.block() - c.block; {
	case diff == 1:
		// Next block as expected; increment window and block
		c.log.trace("ackData diff: %d, current block: %d, rx block %d", diff, c.block, c.rx.block())
		c.block++
		c.window++
		c.catchup = false
	case diff == 0:
		// Same block again, ignore
		c.log.trace("ackData diff: %d, current block: %d, rx block %d", diff, c.block, c.rx.block())
		return c.read
	case diff > c.windowsize:
		c.log.trace("ackData diff: %d, current block: %d, rx block %d", diff, c.block, c.rx.block())
		// Sender is behind, missed ACK? Wait for catchup
		return c.read
	case diff <= c.windowsize:
		c.log.trace("ackData diff: %d, current block: %d, rx block %d", diff, c.block, c.rx.block())
		// We missed blocks
		if c.catchup {
			// Ignore, we need to catchup with server
			return c.read
		}
		// ACK previous block, reset window, and return sequnce error
		c.log.debug("Missing blocks between %d and %d. Resetting to block %d", c.block, c.rx.block(), c.block)
		if err := c.sendAck(c.block); err != nil {
			c.err = wrapError(err, "sending missed block(s) ACK")
			return nil
		}
		c.window = 0
		c.catchup = true
		return c.read
	}

	// Add data to buffer
	n, err := c.rxBuf.Write(c.rx.data())
	if err != nil {
		c.err = wrapError(err, "writing to rxBuf after read")
		return nil
	}

	if n < int(c.blksize) {
		// Reveived last DATA, we're done
		c.done = true
	}

	if c.window < c.windowsize && n >= int(c.blksize) {
		// We haven't reached the window
		return c.read
	}

	// Reached the windowsize or final data, send ACK and reset window
	c.log.trace("window %d, windowsize: %d, offset: %d, blksize: %d", c.window, c.windowsize, c.rx.offset, c.blksize)
	c.window = 0
	c.log.trace("Window %d reached, sending ACK for %d\n", c.windowsize, c.block)
	if err := c.sendAck(c.block); err != nil {
		c.err = wrapError(err, "sending DATA ACK")
		return nil
	}

	return c.read
}

// Close flushes any remaining data to be transferred and closes netConn
func (c *conn) Close() error {
	c.log.debug("Closing connection to %s\n", c.remoteAddr)

	if c.reqChan == nil {
		defer func() {
			// Close network even if another error occurs
			err := c.netConn.Close()
			if err != nil {
				c.log.debug("error closing network connection:", err)
			}
			if c.err == nil {
				c.err = err
			}
		}()
	}

	// Can't write if an error has been sent/received
	if c.err != nil && c.err != io.EOF {
		return wrapError(c.err, "checking conn err before Close")
	}

	// netasciiEnc needs to be flushed if it's in use
	if flusher, ok := c.writer.(interface {
		Flush() error
	}); ok {
		c.log.trace("flushing writer")
		if err := flusher.Flush(); err != nil {
			return wrapError(err, "flushing writer")
		}
	}

	// Write any remaining data, or 0 length DATA to end transfer
	if c.txBuf != nil {
		c.closing = true
		c.Write([]byte{})
	}

	if c.err == io.EOF {
		return nil
	}

	return c.err
}

// parseOACK parses the options from a datagram and returns the successfully
// negotiated options.
func (c *conn) parseOptions() (options, error) {
	ackOpts := make(map[string]string)

	// parse and set options
	for opt, val := range c.rx.options() {
		switch opt {
		case optBlocksize:
			size, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return nil, &errParsingOption{option: opt, value: val}
			}
			c.blksize = uint16(size)
			ackOpts[opt] = val
		case optTimeout:
			seconds, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return nil, &errParsingOption{option: opt, value: val}
			}
			c.timeout = time.Second * time.Duration(seconds)
			ackOpts[opt] = val
		case optTransferSize:
			tsize, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, &errParsingOption{option: opt, value: val}
			}
			if c.isSender && c.tsize != nil {
				// We're sender, send tsize
				ackOpts[opt] = strconv.FormatInt(*c.tsize, 10)
				continue
			}
			c.tsize = &tsize
		case optWindowSize:
			size, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return nil, &errParsingOption{option: opt, value: val}
			}
			c.windowsize = uint16(size)
			ackOpts[opt] = val
		}
	}

	c.optionsParsed = true

	return ackOpts, nil
}

// sendError sends ERROR datagram to remote host
func (c *conn) sendError(code ErrorCode, msg string) {
	c.log.debug("Sending error code %s to %s: %s\n", code, c.remoteAddr, msg)

	// Check error message length
	if len(msg) > int((c.blksize - 1)) { // -1 for NULL terminator
		c.log.debug("error message is larger than blksize, truncating")
		msg = msg[:c.blksize-1]
	}

	// Send error
	c.tx.writeError(code, msg)
	if err := c.writeToNet(); err != nil {
		c.log.debug("sending ERROR: %v", err)
	}
}

// sendAck sends ACK
func (c *conn) sendAck(block uint16) error {
	c.tx.writeAck(block)

	c.log.trace("Sending ACK for %d to %s\n", block, c.remoteAddr)
	return wrapError(c.writeToNet(), "sending ACK")
}

// getAck reads ACK, validates structure and checks for ERROR
//
// If the received ACK is for a previous block, indicating the receiver missed data,
// it will rollback the transfer to the ACK'd block and reset the window.
func (c *conn) getAck() stateType {
	c.tries++
	if c.tries > c.retransmit {
		c.log.debug("Max retries exceeded")
		c.sendError(ErrCodeNotDefined, "max retries reached")
		c.err = wrapError(ErrMaxRetries, "reading ack")
		return nil
	}

	c.log.trace("Waiting for ACK from %s\n", c.remoteAddr)
	sAddr, err := c.readFromNet()
	if err != nil {
		c.log.trace("Error waiting for ACK: %v", err)
		c.err = wrapError(err, "waiting for ACK")
		return c.getAck
	}

	// Send error to requests not from requesting client. May consider
	// ignoring entirely.
	// RFC1350:
	// "If a source TID does not match, the packet should be
	// discarded as erroneously sent from somewhere else.  An error packet
	// should be sent to the source of the incorrect packet, while not
	// disturbing the transfer."
	if c.reqChan == nil && sAddr.String() != c.remoteAddr.String() {
		c.log.err("Received unexpected datagram from %v, expected %v\n", sAddr, c.remoteAddr)
		go func() {
			var err datagram
			err.writeError(ErrCodeUnknownTransferID, "Unexpected TID")
			// Don't care about an error here, just a courtesy
			_, _ = c.netConn.WriteTo(err.bytes(), sAddr)
		}()

		return c.getAck // Read another datagram
	}

	// Validate received datagram
	if err := c.rx.validate(); err != nil {
		c.err = wrapError(err, "ACK validation failed")
		return nil
	}

	// Check opcode
	switch op := c.rx.opcode(); op {
	case opCodeACK:
		c.log.trace("Got ACK for block %d\n", c.rx.block())
		// continue on
	case opCodeERROR:
		c.err = wrapError(c.remoteError(), "error receiving ACK")
		return nil
	default:
		c.err = wrapError(&errUnexpectedDatagram{c.rx.String()}, "error receiving ACK")
		return nil
	}

	// Check block #
	if rxBlock := c.rx.block(); rxBlock != c.block {
		if rxBlock > c.block {
			// Out of order ACKs can cause this scenario, ignore the ACK
			c.log.debug("Received ACK > current block, ignoring.")
			return c.getAck
		}
		c.log.debug("Expected ACK for block %d, got %d. Resetting to block %d.", c.block, rxBlock, rxBlock)
		c.txBuf.UnreadSlots(int(c.block - rxBlock))
		c.block = rxBlock
		c.window = 0

		// Reset done in case error on final send
		c.done = false
	}

	c.tries = 0

	if c.tx.opcode() == opCodeOACK { // TODO: Avoid checking tx opcode?
		return c.write
	}
	return c.writeData
}

// remoteError formats the error in rx, sets err and returns the error.
func (c *conn) remoteError() error {
	c.err = &errRemoteError{dg: c.rx.String()}
	return c.err
}

// readFromNet reads from netConn into b.
func (c *conn) readFromNet() (net.Addr, error) {
	if c.reqChan != nil {
		// Setup timer
		if c.timer == nil {
			c.timer = time.NewTimer(c.timeout)
		} else {
			c.timer.Reset(c.timeout)
		}

		// Single port mode
		select {
		case c.rx.buf = <-c.reqChan:
			c.rx.offset = len(c.rx.buf)
			return nil, nil
		case <-c.timer.C:
			return nil, errors.New("timeout reading from channel")
		}
	}

	if err := c.netConn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, wrapError(err, "setting network read deadline")
	}
	n, addr, err := c.netConn.ReadFrom(c.rx.buf)
	c.rx.offset = n
	return addr, err
}

// writeToNet writes tx to netConn.
func (c *conn) writeToNet() error {
	if err := c.netConn.SetWriteDeadline(time.Now().Add(c.timeout * time.Duration(c.retransmit))); err != nil {
		return wrapError(err, "setting network write deadline")
	}
	_, err := c.netConn.WriteTo(c.tx.bytes(), c.remoteAddr)
	return err
}

// ringBuffer wraps a bytes.Buffer, adding the ability to unread data
// up to the number of slots.
type ringBuffer struct {
	bytes.Buffer
	slots int
	size  int

	buf      []byte // buffer space
	slotsLen []int  // len of data written to each slot
	current  int    // current to be read or written to
	head     int    // head of buffer
}

// newRingBuffer initializes a new ringBuffer
func newRingBuffer(slots int, size int) *ringBuffer {
	return &ringBuffer{
		buf:      make([]byte, size*slots),
		slotsLen: make([]int, size*slots),
		slots:    slots,
		size:     size,
	}
}

// Len returns bytes.Buffer.Len() + any buffer space between current and head
func (r *ringBuffer) Len() int {
	bufInUse := (r.head - r.current) * r.size
	return r.Buffer.Len() + bufInUse
}

// Read reads data from from byte.Buffer if current and head are equal.
// If current is behind head, data will be read from buf.
func (r *ringBuffer) Read(p []byte) (int, error) {
	slot := r.current % r.slots
	offset := slot * r.size

	if r.current != r.head {
		// Copy data out of buf and increment current
		len := offset + r.slotsLen[slot]
		n := copy(p, r.buf[offset:len])
		r.current++
		return n, nil
	}

	// Read from Buffer and copy read data into current slot
	n, err := r.Buffer.Read(p)
	n = copy(r.buf[offset:offset+n], p[:n])
	r.slotsLen[slot] = n

	// Increment current and head
	r.current++
	r.head = r.current
	return n, err
}

// UnreadSlots decrements the current slot, resulting in the
// new reads going to the ringBuffer until current catches up to head
func (r *ringBuffer) UnreadSlots(n int) {
	r.current -= n
}

// readerFunc is an adapter type to convert a function
// to a io.Reader
type readerFunc func([]byte) (int, error)

// Read implements io.Reader
func (f readerFunc) Read(p []byte) (int, error) {
	return f(p)
}

// writerFunc is an adapter type to convert a function
// to a io.Writer
type writerFunc func([]byte) (int, error)

// Write implements io.Writer
func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}

func errorDefer(fn func() error, log *logger, msg string) {
	if err := fn(); err != nil {
		log.debug(msg+": %v", err)
	}
}
