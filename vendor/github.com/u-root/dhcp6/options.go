package dhcp6

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"sort"
)

var (
	// errInvalidOptions is returned when invalid options data is encountered
	// during parsing.  The data could report an incorrect length or have
	// trailing bytes which are not part of the option.
	errInvalidOptions = errors.New("invalid options data")

	// errInvalidOptionRequest is returned when a valid duration cannot be parsed
	// from OptionOptionRequest, because an odd number of bytes are present.
	errInvalidOptionRequest = errors.New("invalid option value for OptionRequestOption")
)

// Options is a map of OptionCode keys with a slice of byte slice values.
// Its methods can be used to easily check for and parse additional
// information from a client request.  If raw data is needed, the map
// can be accessed directly.
type Options map[OptionCode][][]byte

// Add adds a new OptionCode key and BinaryMarshaler struct's bytes to the
// Options map.
func (o Options) Add(key OptionCode, value encoding.BinaryMarshaler) error {
	// Special case: since OptionRapidCommit actually has zero length, it is
	// possible for an option key to appear with no value.
	if value == nil {
		o.addRaw(key, nil)
		return nil
	}

	b, err := value.MarshalBinary()
	if err != nil {
		return err
	}

	o.addRaw(key, b)
	return nil
}

// addRaw adds a new OptionCode key and raw value byte slice to the
// Options map.
func (o Options) addRaw(key OptionCode, value []byte) {
	o[key] = append(o[key], value)
}

// Get attempts to retrieve the first value specified by an OptionCode
// key.  If a value is found, get returns the value and boolean true.
// If it is not found, Get returns nil and boolean false.
func (o Options) Get(key OptionCode) ([]byte, bool) {
	// Empty map has no key/value pairs
	if len(o) == 0 {
		return nil, false
	}

	// Check for value by key
	v, ok := o[key]
	if !ok {
		return nil, false
	}

	// Some options can actually have zero length (OptionRapidCommit),
	// so just return an empty byte slice if this is the case
	if len(v) == 0 {
		return []byte{}, true
	}

	return v[0], true
}

// ClientID returns the Client Identifier Option value, as described in RFC
// 3315, Section 22.2.
//
// The DUID returned allows unique identification of a client to a server.
//
// The boolean return value indicates if OptionClientID was present in the
// Options map.  The error return value indicates if a known, valid DUID type
// could be parsed from the option.
func (o Options) ClientID() (DUID, bool, error) {
	v, ok := o.Get(OptionClientID)
	if !ok {
		return nil, false, nil
	}

	d, err := parseDUID(v)
	return d, true, err
}

// ServerID returns the Server Identifier Option value, as described in RFC
// 3315, Section 22.3.
//
// The DUID returned allows unique identification of a server to a client.
//
// The boolean return value indicates if OptionServerID was present in the
// Options map.  The error return value indicates if a known, valid DUID type
// could be parsed from the option.
func (o Options) ServerID() (DUID, bool, error) {
	v, ok := o.Get(OptionServerID)
	if !ok {
		return nil, false, nil
	}

	d, err := parseDUID(v)
	return d, true, err
}

// IANA returns the Identity Association for Non-temporary Addresses Option
// value, as described in RFC 3315, Section 22.4.
//
// Multiple IANA values may be present in a single DHCP request.
//
// The boolean return value indicates if OptionIANA was present in the Options
// map.  The error return value indicates if one or more valid IANAs could not
// be parsed from the option.
func (o Options) IANA() ([]*IANA, bool, error) {
	// Client may send multiple IANA option requests, so we must
	// access the map directly
	vv, ok := o[OptionIANA]
	if !ok {
		return nil, false, nil
	}

	// Parse each IA_NA value
	iana := make([]*IANA, len(vv), len(vv))
	for i := range vv {
		ia := new(IANA)
		if err := ia.UnmarshalBinary(vv[i]); err != nil {
			return nil, true, err
		}

		iana[i] = ia
	}

	return iana, true, nil
}

// IATA returns the Identity Association for Temporary Addresses Option
// value, as described in RFC 3315, Section 22.5.
//
// Multiple IATA values may be present in a single DHCP request.
//
// The boolean return value indicates if OptionIATA was present in the Options
// map.  The error return value indicates if one or more valid IATAs could not
// be parsed from the option.
func (o Options) IATA() ([]*IATA, bool, error) {
	// Client may send multiple IATA option requests, so we must
	// access the map directly
	vv, ok := o[OptionIATA]
	if !ok {
		return nil, false, nil
	}

	// Parse each IA_NA value
	iata := make([]*IATA, len(vv), len(vv))
	for i := range vv {
		ia := new(IATA)
		if err := ia.UnmarshalBinary(vv[i]); err != nil {
			return nil, true, err
		}

		iata[i] = ia
	}

	return iata, true, nil
}

// IAAddr returns the Identity Association Address Option value, as described
// in RFC 3315, Section 22.6.
//
// The IAAddr option must always appear encapsulated in the Options map of a
// IANA or IATA option.  Multiple IAAddr values may be
// present in a single DHCP request.
//
// The boolean return value indicates if OptionIAAddr was present in the Options
// map.  The error return value indicates if one or more valid IAAddrs could not
// be parsed from the option.
func (o Options) IAAddr() ([]*IAAddr, bool, error) {
	// Client may send multiple IAAddr option requests, so we must
	// access the map directly
	vv, ok := o[OptionIAAddr]
	if !ok {
		return nil, false, nil
	}

	// Parse each IAAddr value
	iaaddr := make([]*IAAddr, len(vv), len(vv))
	for i := range vv {
		iaa := new(IAAddr)
		if err := iaa.UnmarshalBinary(vv[i]); err != nil {
			return nil, true, err
		}

		iaaddr[i] = iaa
	}

	return iaaddr, true, nil
}

// OptionRequest returns the Option Request Option value, as described in RFC
// 3315, Section 22.7.
//
// The slice of OptionCode values indicates the options a DHCP client is
// interested in receiving from a server.
//
// The boolean return value indicates if OptionORO was present in the Options
// map.  The error return value indicates if a valid OptionCode slice could be
// parsed from the option.
func (o Options) OptionRequest() (OptionRequestOption, bool, error) {
	v, ok := o.Get(OptionORO)
	if !ok {
		return nil, false, nil
	}

	var oro OptionRequestOption
	err := oro.UnmarshalBinary(v)
	return oro, true, err
}

// Preference returns the Preference Option value, as described in RFC 3315,
// Section 22.8.
//
// The integer preference value is sent by a server to a client to affect the
// selection of a server by the client.
//
// The boolean return value indicates if OptionPreference was present in the
// Options map.  The error return value indicates if a valid integer value
// could not be parsed from the option.
func (o Options) Preference() (Preference, bool, error) {
	v, ok := o.Get(OptionPreference)
	if !ok {
		return 0, false, nil
	}

	p := new(Preference)
	err := p.UnmarshalBinary(v)
	return *p, true, err
}

// ElapsedTime returns the Elapsed Time Option value, as described in RFC 3315,
// Section 22.9.
//
// The time.Duration returned reports the time elapsed during a DHCP
// transaction, as reported by a client.
//
// The boolean return value indicates if OptionElapsedTime was present in the
// Options map.  The error return value indicates if a valid duration could be
// parsed from the option.
func (o Options) ElapsedTime() (ElapsedTime, bool, error) {
	v, ok := o.Get(OptionElapsedTime)
	if !ok {
		return 0, false, nil
	}

	t := new(ElapsedTime)
	err := t.UnmarshalBinary(v)
	return *t, true, err
}

// RelayMessageOption returns the Relay Message Option value, as described in RFC 3315,
// Section 22.10.
//
// The RelayMessage option carries a DHCP message in a Relay-forward or
// Relay-reply message.
//
// The boolean return value indicates if OptionRelayMsg was present in the
// Options map.  The error return value indicates if a valid OptionRelayMsg could be
// parsed from the option.
func (o Options) RelayMessageOption() (RelayMessageOption, bool, error) {
	v, ok := o.Get(OptionRelayMsg)
	if !ok {
		return nil, false, nil
	}

	r := new(RelayMessageOption)
	err := r.UnmarshalBinary(v)
	return *r, true, err
}

// Authentication returns the Authentication Option value, as described in RFC 3315,
// Section 22.11.
//
// The Authentication option carries authentication information to
// authenticate the identity and contents of DHCP messages.
//
// The boolean return value indicates if Authentication was present in the
// Options map.  The error return value indicates if a valid authentication could be
// parsed from the option.
func (o Options) Authentication() (*Authentication, bool, error) {
	v, ok := o.Get(OptionAuth)
	if !ok {
		return nil, false, nil
	}

	a := new(Authentication)
	err := a.UnmarshalBinary(v)
	return a, true, err
}

// Unicast returns the IP from a Unicast Option value, described in RFC 3315,
// Section 22.12.
//
// The IP return value indicates a server's IPv6 address, which a client may
// use to contact the server via unicast.
//
// The boolean return value indicates if OptionUnicast was present in the
// Options map.  The error return value indicates if a valid IPv6 address
// could not be parsed from the option.
func (o Options) Unicast() (IP, bool, error) {
	v, ok := o.Get(OptionUnicast)
	if !ok {
		return nil, false, nil
	}

	var ip IP
	err := ip.UnmarshalBinary(v)
	return ip, true, err
}

// StatusCode returns the Status Code Option value, described in RFC 3315,
// Section 22.13.
//
// The StatusCode return value may be used to determine a code and an
// explanation for the status.
//
// The boolean return value indicates if OptionStatusCode was present in the
// Options map.  The error return value indicates if a valid StatusCode could
// not be parsed from the option.
func (o Options) StatusCode() (*StatusCode, bool, error) {
	v, ok := o.Get(OptionStatusCode)
	if !ok {
		return nil, false, nil
	}

	s := new(StatusCode)
	err := s.UnmarshalBinary(v)
	return s, true, err
}

// RapidCommit returns the Rapid Commit Option value, described in RFC 3315,
// Section 22.14.
//
// The boolean return value indicates if OptionRapidCommit was present in the
// Options map, and thus, if Rapid Commit should be used.
//
// The error return value indicates if a valid Rapid Commit Option could not
// be parsed.
func (o Options) RapidCommit() (bool, error) {
	v, ok := o.Get(OptionRapidCommit)
	if !ok {
		return false, nil
	}

	// Data must be completely empty; presence of the Rapid Commit option
	// indicates it is requested.
	if len(v) != 0 {
		return false, io.ErrUnexpectedEOF
	}

	return true, nil
}

// UserClass returns the User Class Option value, described in RFC 3315,
// Section 22.15.
//
// The Data structure returned contains any raw class data present in
// the option.
//
// The boolean return value indicates if OptionUserClass was present in the
// Options map.  The error return value indicates if any errors were present
// in the class data.
func (o Options) UserClass() (Data, bool, error) {
	v, ok := o.Get(OptionUserClass)
	if !ok {
		return nil, false, nil
	}

	var d Data
	err := d.UnmarshalBinary(v)
	return d, true, err
}

// VendorClass returns the Vendor Class Option value, described in RFC 3315,
// Section 22.16.
//
// The VendorClass structure returned contains VendorClass in
// the option.
//
// The boolean return value indicates if OptionVendorClass was present in the
// Options map.  The error return value indicates if any errors were present
// in the VendorClass data.
func (o Options) VendorClass() (*VendorClass, bool, error) {
	v, ok := o.Get(OptionVendorClass)
	if !ok {
		return nil, false, nil
	}

	vc := new(VendorClass)
	err := vc.UnmarshalBinary(v)
	return vc, true, err
}

// VendorOpts returns the Vendor-specific Information Option value, described in RFC 3315,
// Section 22.17.
//
// The VendorOpts structure returned contains Vendor-specific Information data present in
// the option.
//
// The boolean return value indicates if VendorOpts was present in the
// Options map.  The error return value indicates if any errors were present
// in the class data.
func (o Options) VendorOpts() (*VendorOpts, bool, error) {
	v, ok := o.Get(OptionVendorOpts)
	if !ok {
		return nil, false, nil
	}

	vo := new(VendorOpts)
	err := vo.UnmarshalBinary(v)
	return vo, true, err
}

// InterfaceID returns the Interface-Id Option value, described in RFC 3315,
// Section 22.18.
//
// The InterfaceID structure returned contains any raw class data present in
// the option.
//
// The boolean return value indicates if InterfaceID was present in the
// Options map.  The error return value indicates if any errors were present
// in the interface-id data.
func (o Options) InterfaceID() (InterfaceID, bool, error) {
	v, ok := o.Get(OptionInterfaceID)
	if !ok {
		return nil, false, nil
	}

	var i InterfaceID
	err := i.UnmarshalBinary(v)
	return i, true, err
}

// IAPD returns the Identity Association for Prefix Delegation Option value,
// described in RFC 3633, Section 9.
//
// Multiple IAPD values may be present in a a single DHCP request.
//
// The boolean return value indicates if OptionIAPD was present in the Options
// map.  The error return value indicates if one or more valid IAPDs could not
// be parsed from the option.
func (o Options) IAPD() ([]*IAPD, bool, error) {
	// Client may send multiple IAPD option requests, so we must
	// access the map directly
	vv, ok := o[OptionIAPD]
	if !ok {
		return nil, false, nil
	}

	// Parse each IA_PD value
	iapd := make([]*IAPD, len(vv))
	for i := range vv {
		ia := new(IAPD)
		if err := ia.UnmarshalBinary(vv[i]); err != nil {
			return nil, true, err
		}

		iapd[i] = ia
	}

	return iapd, true, nil
}

// IAPrefix returns the Identity Association Prefix Option value, as described
// in RFC 3633, Section 10.
//
// Multiple IAPrefix values may be present in a a single DHCP request.
//
// The boolean return value indicates if OptionIAPrefix was present in the
// Options map.  The error return value indicates if one or more valid
// IAPrefixes could not be parsed from the option.
func (o Options) IAPrefix() ([]*IAPrefix, bool, error) {
	// Client may send multiple IAPrefix option requests, so we must
	// access the map directly
	vv, ok := o[OptionIAPrefix]
	if !ok {
		return nil, false, nil
	}

	// Parse each IAPrefix value
	iaprefix := make([]*IAPrefix, len(vv))
	for i := range vv {
		ia := new(IAPrefix)
		if err := ia.UnmarshalBinary(vv[i]); err != nil {
			return nil, true, err
		}

		iaprefix[i] = ia
	}

	return iaprefix, true, nil
}

// RemoteIdentifier returns the Remote Identifier, described in RFC 4649.
//
// This option may be added by DHCPv6 relay agents that terminate
// switched or permanent circuits and have mechanisms to identify the
// remote host end of the circuit.
//
// The boolean return value indicates if OptionRemoteIdentifier was present in the
// Options map.  The error return value indicates if any errors were present
// in the class data.
func (o Options) RemoteIdentifier() (*RemoteIdentifier, bool, error) {
	v, ok := o.Get(OptionRemoteIdentifier)
	if !ok {
		return nil, false, nil
	}

	r := new(RemoteIdentifier)
	err := r.UnmarshalBinary(v)
	return r, true, err
}

// BootFileURL returns the Boot File URL Option value, described in RFC 5970,
// Section 3.1.
//
// The URL return value contains a URL which may be used by clients to obtain
// a boot file for PXE.
//
// The boolean return value indicates if OptionBootFileURL was present in the
// Options map.  The error return value indicates if a valid boot file URL
// could not be parsed from the option.
func (o Options) BootFileURL() (*URL, bool, error) {
	v, ok := o.Get(OptionBootFileURL)
	if !ok {
		return nil, false, nil
	}

	u := new(URL)
	err := u.UnmarshalBinary(v)
	return u, true, err
}

// BootFileParam returns the Boot File Parameters Option value, described in
// RFC 5970, Section 3.2.
//
// The Data structure returned contains any parameters needed for a boot
// file, such as a root filesystem label or a path to a configuration file for
// further chainloading.
//
// The boolean return value indicates if OptionBootFileParam was present in
// the Options map.  The error return value indicates if valid boot file
// parameters could not be parsed from the option.
func (o Options) BootFileParam() (Data, bool, error) {
	v, ok := o.Get(OptionBootFileParam)
	if !ok {
		return nil, false, nil
	}

	var d Data
	err := d.UnmarshalBinary(v)
	return d, true, err
}

// ClientArchType returns the Client System Architecture Type Option value,
// described in RFC 5970, Section 3.3.
//
// The ArchTypes slice returned contains a list of one or more ArchType values.
// The first ArchType listed is the client's most preferable value.
//
//
// The boolean return value indicates if OptionClientArchType was present in
// the Options map.  The error return value indicates if a valid list of
// ArchType values could not be parsed from the option.
func (o Options) ClientArchType() (ArchTypes, bool, error) {
	v, ok := o.Get(OptionClientArchType)
	if !ok {
		return nil, false, nil
	}

	var a ArchTypes
	err := a.UnmarshalBinary(v)
	return a, true, err
}

// NII returns the Client Network Interface Identifier Option value, described
// in RFC 5970, Section 3.4.
//
// The NII value returned indicates a client's level of Universal Network
// Device Interface (UNDI) support.
//
// The boolean return value indicates if OptionNII was present in
// the Options map.  The error return value indicates if a valid list of
// ArchType values could not be parsed from the option.
func (o Options) NII() (*NII, bool, error) {
	v, ok := o.Get(OptionNII)
	if !ok {
		return nil, false, nil
	}

	n := new(NII)
	err := n.UnmarshalBinary(v)
	return n, true, err
}

// byOptionCode implements sort.Interface for optslice.
type byOptionCode optslice

func (b byOptionCode) Len() int               { return len(b) }
func (b byOptionCode) Less(i int, j int) bool { return b[i].Code < b[j].Code }
func (b byOptionCode) Swap(i int, j int)      { b[i], b[j] = b[j], b[i] }

// enumerate returns an ordered slice of option data from the Options map,
// for use with sending responses to clients.
func (o Options) enumerate() optslice {
	// Send all values for a given key
	var options optslice
	for k, v := range o {
		for _, vv := range v {
			options = append(options, option{
				Code: k,
				Data: vv,
			})
		}
	}

	sort.Sort(byOptionCode(options))
	return options
}

// parseOptions returns a slice of option code and values from an input byte
// slice.  It is used with various different types to enable parsing of both
// top-level options, and options embedded within other options.  If options
// data is malformed, it returns errInvalidOptions.
func parseOptions(b []byte) (Options, error) {
	var length int
	options := make(Options)

	buf := bytes.NewBuffer(b)

	for buf.Len() > 3 {
		// 2 bytes: option code
		o := option{}
		o.Code = OptionCode(binary.BigEndian.Uint16(buf.Next(2)))

		// 2 bytes: option length
		length = int(binary.BigEndian.Uint16(buf.Next(2)))

		// If length indicated is zero, skip to next iteration
		if length == 0 {
			continue
		}

		// N bytes: option data
		o.Data = buf.Next(length)
		// Set slice's max for option's data
		o.Data = o.Data[:len(o.Data):len(o.Data)]

		// If option data has less bytes than indicated by length,
		// return an error
		if len(o.Data) < length {
			return nil, errInvalidOptions
		}

		options.addRaw(o.Code, o.Data)
	}

	// Report error for any trailing bytes
	if buf.Len() != 0 {
		return nil, errInvalidOptions
	}

	return options, nil
}

// option represents an individual DHCP Option, as defined in RFC 3315,
// Section 22.  An Option carries both an OptionCode and its raw Data.  The
// format of option data varies depending on the option code.
type option struct {
	Code OptionCode
	Data []byte
}

// optslice is a slice of option values, and is used to help marshal option
// values into binary form.
type optslice []option

// count returns the number of bytes that this slice of options will occupy
// when marshaled to binary form.
func (o optslice) count() int {
	var c int
	for _, oo := range o {
		// 2 bytes: option code
		// 2 bytes: option length
		// N bytes: option data
		c += 2 + 2 + len(oo.Data)
	}

	return c
}

// write writes the option slice into the provided buffer.  The caller must
// ensure that a large enough buffer is provided to write to avoid panics.
func (o optslice) write(p []byte) {
	var i int
	for _, oo := range o {
		// 2 bytes: option code
		binary.BigEndian.PutUint16(p[i:i+2], uint16(oo.Code))
		i += 2

		// 2 bytes: option length
		binary.BigEndian.PutUint16(p[i:i+2], uint16(len(oo.Data)))
		i += 2

		// N bytes: option data
		copy(p[i:i+len(oo.Data)], oo.Data)
		i += len(oo.Data)
	}
}
