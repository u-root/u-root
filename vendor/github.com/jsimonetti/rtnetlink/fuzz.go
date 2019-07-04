// +build gofuzz

package rtnetlink

func Fuzz(data []byte) int {
	return fuzzLinkMessage(data)
	//return fuzzAddressMessage(data)
	//return fuzzRouteMessage(data)
	//return fuzzNeighMessage(data)
}

func fuzzLinkMessage(data []byte) int {
	m := &LinkMessage{}
	if err := (m).UnmarshalBinary(data); err != nil {
		return 0
	}

	if _, err := m.MarshalBinary(); err != nil {
		panic(err)
	}

	return 1
}

func fuzzAddressMessage(data []byte) int {
	m := &AddressMessage{}
	if err := (m).UnmarshalBinary(data); err != nil {
		return 0
	}

	if _, err := m.MarshalBinary(); err != nil {
		panic(err)
	}

	return 1
}

func fuzzRouteMessage(data []byte) int {
	m := &RouteMessage{}
	if err := (m).UnmarshalBinary(data); err != nil {
		return 0
	}

	if _, err := m.MarshalBinary(); err != nil {
		panic(err)
	}

	return 1
}

func fuzzNeighMessage(data []byte) int {
	m := &NeighMessage{}
	if err := (m).UnmarshalBinary(data); err != nil {
		return 0
	}

	if _, err := m.MarshalBinary(); err != nil {
		panic(err)
	}

	return 1
}
