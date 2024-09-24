package tc

import (
	"fmt"
)

// SfqQopt contains SFQ attributes
type SfqQopt struct {
	Quantum       uint32 /* Bytes per round allocated to flow */
	PerturbPeriod int32  /* Period of hash perturbation */
	Limit         uint32 /* Maximal packets in queue */
	Divisor       uint32 /* Hash divisor  */
	Flows         uint32 /* Maximal number of flows  */
}

// Sfq contains attributes of the SFQ discipline
// https://man7.org/linux/man-pages/man8/sfq.8.html
type Sfq struct {
	V0 SfqQopt

	Depth    uint32 /* max number of packets per flow */
	Headdrop uint32

	/* SFQRED parameters */
	Limit    uint32 /* HARD maximal flow queue length (bytes) */
	QthMin   uint32 /* Min average length threshold (bytes) */
	QthMax   uint32 /* Max average length threshold (bytes) */
	Wlog     uint8  /* log(W)		*/
	Plog     uint8  /* log(P_max/(qth_max-qth_min))	*/
	ScellLog uint8  /* cell size for idle damping */
	Flags    uint8
	MaxP     uint32 /* probability, high resolution */
}

// unmarshalSfq parses the Sfq-encoded data and stores the result in the value pointed to by info.
func unmarshalSfq(data []byte, info *Sfq) error {
	return unmarshalStruct(data, info)
}

// marshalSfq returns the binary encoding of Sfq
func marshalSfq(info *Sfq) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("Sfq: %w", ErrNoArg)
	}

	// TODO: improve logic and check combinations
	return marshalStruct(info)
}
