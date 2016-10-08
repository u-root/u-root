// Copyright (C) 2016 Kale Blankenship. All rights reserved.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details

package trivialt

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/vcabbage/trivialt/netascii"
)

const (
	defaultPort       = 69
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
func newConn(udpNet string, mode transferMode, addr *net.UDPAddr) (*conn, error) {
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

func newSinglePortConn(addr *net.UDPAddr, mode transferMode, netConn *net.UDPConn, reqChan chan []byte) *conn {
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
func newConnFromHost(udpNet string, mode transferMode, host string) (*conn, error) {
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
	mode       transferMode  // octet or netascii
	tsize      *int64        // Size of the file being sent/received

	// Other, non-negotiable options
	retransmit int // Number of times an individual datagram will be retransmitted on error

	// Track state of transfer
	optionsParsed bool   // Whether TFTP options have been parsed yet
	window        uint16 // Packets sent since last ACK
	block         uint16 // Current block #
	catchup       bool   // Ignore incoming blocks from a window we reset

	// Buffers
	buf   []byte       // incoming data from, sized to blksize + headers
	txBuf *ringBuffer  // buffers outgoing data, retaining windowsize * blksize
	rxBuf bytes.Buffer // buffer incoming data

	// Datgrams
	tx datagram // Constructs outgoing datagrams
	rx datagram // Hold and parse current incoming datagram

	// netascii encoder/decoder wrap Read/Write methods when transfer
	// is in netascii mode
	netasciiEnc *netascii.Writer
	netasciiRdr *netascii.Reader

	// Indicates the transfer is complete
	done bool

	// Indicates an ERROR has been sent or received
	err error
}

// sendWriteRequest sends WRQ to server and negotiates transfer options
func (c *conn) sendWriteRequest(filename string, opts map[string]string) error {
	// Build WRQ
	c.tx.writeWriteReq(filename, c.mode, opts)

	// Send request
	if err := c.sendRequest(); err != nil {
		return wrapError(err, "sending WRQ")
	}

	// Should have received OACK if server supports options, or DATA if not
	switch op := c.rx.opcode(); op {
	case opCodeOACK, opCodeACK:
		// Got OACK, parse options
		if err := c.writeSetup(); err != nil {
			return wrapError(err, "parsing OACK to WRQ")
		}
	case opCodeERROR:
		// Received an error
		return wrapError(c.remoteError(), "WRQ OACK response")
	default:
		return wrapError(&errUnexpectedDatagram{dg: c.rx.String()}, "WRQ OACK response")
	}

	return nil
}

// sendReadRequest send RRQ to server and negotiates transfer options
//
// If the server doesn't support options and responds with data, the data will be added
// to rxBuf.
func (c *conn) sendReadRequest(filename string, opts map[string]string) error {
	// Build RRQ
	c.tx.writeReadReq(filename, c.mode, opts)

	// Send request
	if err := c.sendRequest(); err != nil {
		return wrapError(err, "sending RRQ")
	}

	// Should have received OACK if server supports options, or DATA if not
	switch op := c.rx.opcode(); op {
	case opCodeOACK:
		// Got OACK, parse options
		if err := c.readSetup(); err != nil {
			return wrapError(err, "got OACK, read setup")
		}
	case opCodeDATA:
		// Server doesn't support options,
		// write data to the buf so it's available for reading

		var writer io.Writer = &c.rxBuf
		if c.mode == ModeNetASCII {
			na := netascii.NewWriter(writer)
			defer errorDefer(na.Flush, c.log, "error flushing netascii encoder")
			writer = na
		}

		n, err := writer.Write(c.rx.data())
		if err != nil {
			return wrapError(err, "writing RRQ response data")
		}

		// Set optionsParsed
		c.optionsParsed = true

		// If less than blksize, we're done
		if n < int(c.blksize) {
			c.done = true
		}

		// Ack data
		if err := c.sendAck(c.rx.block()); err != nil {
			return wrapError(err, "sending ACK for RRQ DATA")
		}
		c.block = c.rx.block()
	case opCodeERROR:
		// Received an error
		return wrapError(c.remoteError(), "RRQ OACK response")
	default:
		return wrapError(&errUnexpectedDatagram{dg: c.rx.String()}, "RRQ OACK response")
	}

	return nil
}

func (c *conn) sendRequest() error {
	// Set that we're a client
	c.isClient = true

	// Send request
	if err := c.writeToNet(); err != nil {
		return wrapError(err, "writing request to network")
	}

	// Receive response
	for retries := 0; ; {
		addr, err := c.readFromNet()
		if err == nil {
			if c.reqChan == nil {
				// Update address
				c.remoteAddr = addr
			}
			break
		}

		if retries < c.retransmit {
			c.log.debug("error getting %s response from %v", c.tx.opcode(), c.remoteAddr)
			retries++
			continue
		}

		return wrapError(err, "receiving request response")
	}

	// Extract and validate response
	return wrapError(c.rx.validate(), "validating request response")
}

// Write implements io.Writer and wraps write().
//
// If mode is ModeNetASCII, wrap write() with netascii.EncodeWriter.
func (c *conn) Write(p []byte) (int, error) {
	// Can't write if an error has been sent/received
	if c.err != nil {
		return 0, wrapError(c.err, "checking conn err before Write")
	}

	if c.mode == ModeNetASCII {
		if c.netasciiEnc == nil {
			c.netasciiEnc = netascii.NewWriter(writerFunc(c.write))
		}
		n, err := c.netasciiEnc.Write(p)
		return n, wrapError(err, "writing through netascii encoder")
	}

	n, err := c.write(p)
	return n, wrapError(err, "writing through standard writer")
}

// writeSetup parses options and sets up buffers before
// first write.
func (c *conn) writeSetup() error {
	// Set that we're sending
	c.isSender = true

	ackOpts, err := c.parseOptions()
	if err != nil {
		return wrapError(err, "write setup")
	}

	// Set buf size
	if len(c.buf) != int(c.blksize) {
		c.buf = make([]byte, c.blksize)
	}

	// Init ringBuffer
	c.txBuf = newRingBuffer(int(c.windowsize), int(c.blksize))

	// Sending DATA ACKs when there are no options
	// Client doesn't send OACK
	if len(ackOpts) == 0 || c.isClient {
		return nil
	}

	// Send OACK
	c.log.trace("Sending OACK to %s\n", c.remoteAddr)
	c.tx.writeOptionAck(ackOpts)
	if err := c.writeToNet(); err != nil {
		return wrapError(err, "sending OACK in response to options")
	}

	return wrapError(c.getAck(), "write setup")
}

// write writes adds data to txBuf and writes data to netConn in chunks of
// blksize, until the last chunk of <blksize, which signals transfer completion.
func (c *conn) write(p []byte) (int, error) {
	// Options won't be parsed before first write so that API consumer
	// has opportunity to set tsize with ReadRequest.WriteSize()
	if !c.optionsParsed {
		if err := c.writeSetup(); err != nil {
			return 0, wrapError(err, "parsing options before write")
		}
	}
	// Copy to buffer
	read, err := c.txBuf.Write(p)
	if err != nil {
		return read, wrapError(err, "writing data to txBuf before write")
	}

	for c.txBuf.Len() >= int(c.blksize) {
		if err := c.writeData(); err != nil {
			return read, wrapError(err, "writing data")
		}
		// Increment the window
		c.window++

		// Continue on if we haven't reached the windowsize
		if c.window < c.windowsize {
			continue
		}

		// Reset window
		c.window = 0

		// Get ACK
		retries := 0
		for {
			err := c.getAck()
			if err == nil {
				break
			}

			// Return if the transfer has erred
			if c.err != nil {
				return read, wrapError(c.err, "writing data")
			}

			// Retry until maxRetransmit
			retries++
			c.log.debug("Error receiving ACK (retry %d): %v\n", retries, err)
			if retries > c.retransmit {
				c.log.debug("Max retries exceeded")
				c.sendError(ErrCodeNotDefined, "max retries reached")
				return read, wrapError(err, "receiving ACK after writing data")
			}
		}
	}

	return read, nil
}

// writeData writes a single DATA datagram
func (c *conn) writeData() error {
	c.block++

	// Read data from txBuf
	n, err := c.txBuf.Read(c.buf)
	if err != nil && err != io.EOF {
		return wrapError(err, "reading data from txBuf before writing to network")
	}
	c.tx.writeData(c.block, c.buf[:n])

	// Send w.tx datagram
	c.log.trace("Sending block %d with %d bytes to %s\n", c.block, n, c.remoteAddr)
	return wrapError(c.writeToNet(), "writing data to network")
}

// Read implements io.Reader and wraps read()
//
// If mode is ModeNetASCII, read() is wrapped with netascii.ReadDecoder
func (c *conn) Read(b []byte) (int, error) {
	// Can't read if an error has been sent/received
	if c.err != nil {
		return 0, wrapError(c.err, "checking conn error before Read")
	}

	if c.mode == ModeNetASCII {
		if c.netasciiRdr == nil {
			c.netasciiRdr = netascii.NewReader(readerFunc(c.read))
		}
		n, err := c.netasciiRdr.Read(b)
		if err != io.EOF {
			err = wrapError(err, "Read from netascii decoder")
		}
		return n, err
	}

	n, err := c.read(b)
	if err != io.EOF {
		err = wrapError(err, "Read from standard reader")
	}
	return n, err
}

// readSetup parses options and sets up buffers before
// first read.
func (c *conn) readSetup() error {
	ackOpts, err := c.parseOptions()
	if err != nil {
		return wrapError(err, "read setup")
	}

	// Set buf size
	needed := int(c.blksize + 4)
	if len(c.rx.buf) != needed {
		c.rx.buf = make([]byte, needed)
	}

	// If there we're not options negotiated, send ACK
	// Client never sends OACK
	if len(ackOpts) == 0 || c.isClient {
		c.log.trace("Sending ACK to %s\n", c.remoteAddr)
		c.tx.writeAck(0)
	} else {
		c.log.trace("Sending OACK to %s\n", c.remoteAddr)
		c.tx.writeOptionAck(ackOpts)
	}

	// Send ACK/OACK
	return wrapError(c.writeToNet(), "sending ACK/OACK in response to options")
}

// read reads data from netConn until p is full or the connection is
// complete.
func (c *conn) read(p []byte) (int, error) {
	if !c.optionsParsed {
		if err := c.readSetup(); err != nil {
			return 0, wrapError(err, "parsing options before initial read")
		}
	}

	for l := len(p); c.rxBuf.Len() < l && !c.done; {
		// Read next datagram
		if err := c.readData(); err != nil {
			return 0, wrapError(err, "reading data")
		}

		if err := c.ackData(); err != nil {
			if err == errBlockSequence {
				continue
			}
			return 0, wrapError(err, "reading data")
		}

		// Add data to buffer
		n, err := c.rxBuf.Write(c.rx.data())
		if err != nil {
			return 0, wrapError(err, "writing to rxBuf after read")
		}

		if n < int(c.blksize) {
			// Reveived last DATA, we're done
			c.done = true
		}
	}

	// Read buffered data into p
	read, err := c.rxBuf.Read(p)
	if err != nil && err != io.EOF { // Ignore EOF from bytes.Buffer
		return read, wrapError(err, "reading from rxBuf after read")
	}
	// If done, signal that there's nothing more to read by io.EOF
	if c.done && c.rxBuf.Len() == 0 {
		return read, io.EOF
	}

	return read, nil
}

// readDatagram reads a single datagram into rx
func (c *conn) readData() error {
	for retries := 0; ; {
		c.log.trace("Waiting for DATA from %s\n", c.remoteAddr)
		_, err := c.readFromNet()
		if err == nil {
			break
		}

		c.log.debug("error receiving block %d: %v", c.block+1, err)
		if retries == c.retransmit {
			c.log.debug("Max retransmit reached, ending transfer")
			return wrapError(err, "reading data")
		}

		c.log.trace("Resending ACK for %d\n", c.block)
		if err := c.sendAck(c.block); err != nil {
			c.log.debug("resending ACK %v", err)
		}
		c.window = 0
		retries++
	}

	// validate datagram
	if err := c.rx.validate(); err != nil {
		return wrapError(err, "validating read data")
	}

	// Check for opcode
	switch op := c.rx.opcode(); op {
	case opCodeDATA:
	case opCodeERROR:
		// Received an error
		return wrapError(c.remoteError(), "reading data")
	default:
		return wrapError(&errUnexpectedDatagram{dg: c.rx.String()}, "read data response")
	}

	c.log.trace("Received block %d\n", c.rx.block())

	return nil
}

// ackData handles block sequence, windowing, and acknowledgements
func (c *conn) ackData() error {
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
		return errBlockSequence
	case diff > c.windowsize:
		c.log.trace("ackData diff: %d, current block: %d, rx block %d", diff, c.block, c.rx.block())
		// Sender is behind, missed ACK? Wait for catchup
		return errBlockSequence
	case diff <= c.windowsize:
		c.log.trace("ackData diff: %d, current block: %d, rx block %d", diff, c.block, c.rx.block())
		// We missed blocks
		if c.catchup {
			// Ignore, we need to catchup with server
			return errBlockSequence
		}
		// ACK previous block, reset window, and return sequnce error
		c.log.debug("Missing blocks between %d and %d. Resetting to block %d", c.block, c.rx.block(), c.block)
		if err := c.sendAck(c.block); err != nil {
			return wrapError(err, "sending missed block(s) ACK")
		}
		c.window = 0
		c.catchup = true
		return errBlockSequence
	}

	// If we've reached the windowsize, send ACK and reset window
	if c.window >= c.windowsize || c.rx.offset < int(c.blksize) {
		c.log.trace("window %d, windowsize: %d, offset: %d, blksize: %d", c.window, c.windowsize, c.rx.offset, c.blksize)
		c.window = 0
		c.log.trace("Window %d reached, sending ACK for %d\n", c.windowsize, c.block)
		if err := c.sendAck(c.block); err == nil {
			return wrapError(err, "sending DATA ACK")
		}
	}

	return nil
}

// Close flushes any remaining data to be transferred and closes netConn
func (c *conn) Close() (err error) {
	c.log.debug("Closing connection to %s\n", c.remoteAddr)

	if c.reqChan == nil {
		defer func() {
			// Close network even if another error occurs
			//
			cErr := c.netConn.Close()
			if cErr != nil {
				c.log.debug("error closing network connection:", cErr)
			}
			if err == nil {
				err = cErr
			}
		}()
	}

	// Can't write if an error has been sent/received
	if c.err != nil {
		return wrapError(c.err, "checking conn err before Close")
	}

	// netasciiEnc needs to be flushed if it's in use
	if c.netasciiEnc != nil {
		c.log.trace("Flushing netascii encoder")
		if err := c.netasciiEnc.Flush(); err != nil {
			return wrapError(err, "flushing netascii encoder before Close")
		}
	}

	// Write any remaining data, or 0 length DATA to end transfer
	if c.txBuf != nil {
		for retries := 0; ; {
			if err := c.writeData(); err != nil {
				return wrapError(err, "writing final data before Close")
			}

			err := c.getAck()
			if err == nil {
				// Recheck data, window could have been missed
				if c.txBuf.Len() >= int(c.blksize) {
					c.log.trace("%d", c.txBuf.Len())
					c.write([]byte{})
					continue
				}
				if c.txBuf.Len() > 0 {
					continue
				}
				break
			}

			// Return if the transfer has erred
			if c.err != nil {
				return wrapError(c.err, "checking conn error sending ACK for final data before Close")
			}

			// Retry until maxRetransmit
			retries++
			c.log.debug("Error receiving ACK (retry %d): %v\n", retries, err)
			if retries > c.retransmit {
				c.log.debug("Max retries exceeded")
				c.sendError(ErrCodeNotDefined, "max retries reached")
				return wrapError(err, "writing final data ACK before Close")
			}
		}
	}

	return nil
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

	// Set error
	c.err = &errRemoteError{dg: c.tx.String()}
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
func (c *conn) getAck() error {
	for {
		c.log.trace("Waiting for ACK from %s\n", c.remoteAddr)
		sAddr, err := c.readFromNet()
		if err != nil {
			return wrapError(err, "network read failed")
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

			continue // Read another datagram
		}
		break
	}

	// Validate received datagram
	if err := c.rx.validate(); err != nil {
		return wrapError(err, "ACK validation failed")
	}

	// Check opcode
	switch op := c.rx.opcode(); op {
	case opCodeACK:
		c.log.trace("Got ACK for block %d\n", c.rx.block())
		// continue on
	case opCodeERROR:
		return wrapError(c.remoteError(), "error receiving ACK")
	default:
		return wrapError(&errUnexpectedDatagram{c.rx.String()}, "error receiving ACK")
	}

	// Check block #
	if rxBlock := c.rx.block(); rxBlock != c.block {
		if rxBlock > c.block {
			// Out of order ACKs can cause this scenario, ignore the ACK
			return nil
		}
		c.log.debug("Expected ACK for block %d, got %d. Resetting to block %d.", c.block, rxBlock, rxBlock)
		c.txBuf.UnreadSlots(int(c.block - rxBlock))
		c.block = rxBlock
		c.window = 0
	}

	return nil
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
