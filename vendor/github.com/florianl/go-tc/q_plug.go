package tc

import "fmt"

// PlugAction defines actions for plug.
type PlugAction int32

// Various Plug actions.
const (
	PlugBuffer PlugAction = iota
	PlugReleaseOne
	PlugReleaseIndefinite
	PlugLimit
)

// Plug contains attributes of the plug discipline
type Plug struct {
	Action PlugAction
	Limit  uint32
}

func marshalPlug(info *Plug) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("Plug: %w", ErrNoArg)
	}
	return marshalStruct(info)
}

func unmarshalPlug(data []byte, info *Plug) error {
	// So far the kernel does not implement this functionality.
	return ErrNotImplemented
}
