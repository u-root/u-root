package tc

import "fmt"

// Prio contains attributes of the prio discipline
type Prio struct {
	Bands   uint32
	PrioMap [16]uint8
}

// unmarshalPrio parses the Prio-encoded data and stores the result in the value pointed to by info.
func unmarshalPrio(data []byte, info *Prio) error {
	return unmarshalStruct(data, info)
}

// marshalPrio returns the binary encoding of MqPrio
func marshalPrio(info *Prio) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("Prio: %w", ErrNoArg)
	}
	return marshalStruct(info)
}
