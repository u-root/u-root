package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaTaPrioUnspec           = iota
	tcaTaPrioPrioMap          /* struct tc_mqprio_qopt */
	tcaTaPrioSchedEntryList   /* nested of entry */
	tcaTaPrioSchedBaseTime    /* s64 */
	tcaTaPrioSchedSingleEntry /* single entry */
	tcaTaPrioSchedClockID     /* s32 */
	tcaTaPrioPad
	tcaTaPrioAdminSched              /* The admin sched, only used in dump */
	tcaTaPrioSchedCycleTime          /* s64 */
	tcaTaPrioSchedCycleTimeExtension /* s64 */
	tcaTaPrioFlags                   /* u32 */
	tcaTaPrioTxTimeDelay             /* u32 */
	tcaTaPrioTcEntry                 /* nest */
)

// TaPrio contains TaPrio attributes
type TaPrio struct {
	PrioMap                 *MqPrioQopt
	SchedBaseTime           *int64
	SchedClockID            *int32
	SchedCycleTime          *int64
	SchedCycleTimeExtension *int64
	Flags                   *uint32
	TxTimeDelay             *uint32
}

// unmarshalTaPrio parses the TaPrio-encoded data and stores the result in the value pointed to by info.
func unmarshalTaPrio(data []byte, info *TaPrio) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaTaPrioPrioMap:
			opt := &MqPrioQopt{}
			err := unmarshalStruct(ad.Bytes(), opt)
			multiError = concatError(multiError, err)
			info.PrioMap = opt
		case tcaTaPrioSchedBaseTime:
			info.SchedBaseTime = int64Ptr(ad.Int64())
		case tcaTaPrioSchedClockID:
			info.SchedClockID = int32Ptr(ad.Int32())
		case tcaTaPrioSchedCycleTime:
			info.SchedCycleTime = int64Ptr(ad.Int64())
		case tcaTaPrioSchedCycleTimeExtension:
			info.SchedCycleTimeExtension = int64Ptr(ad.Int64())
		case tcaTaPrioFlags:
			info.Flags = uint32Ptr(ad.Uint32())
		case tcaTaPrioTxTimeDelay:
			info.TxTimeDelay = uint32Ptr(ad.Uint32())

		case tcaTaPrioPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalTaPrio()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}

// marshalTaPrio returns the binary encoding of TaPrio
func marshalTaPrio(info *TaPrio) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("TaPrio: %w", ErrNoArg)
	}
	options := []tcOption{}

	// TODO: improve logic and check combinations
	var multiError error

	if info.PrioMap != nil {
		data, err := marshalStruct(info.PrioMap)
		multiError = concatError(multiError, err)
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTaPrioPrioMap, Data: data})
	}
	if info.SchedBaseTime != nil {
		options = append(options, tcOption{Interpretation: vtInt64, Type: tcaTaPrioSchedBaseTime, Data: int64Value(info.SchedBaseTime)})
	}
	if info.SchedClockID != nil {
		options = append(options, tcOption{Interpretation: vtInt32, Type: tcaTaPrioSchedClockID, Data: int32Value(info.SchedClockID)})
	}
	if info.SchedCycleTime != nil {
		options = append(options, tcOption{Interpretation: vtInt64, Type: tcaTaPrioSchedCycleTime, Data: int64Value(info.SchedCycleTime)})
	}
	if info.SchedCycleTimeExtension != nil {
		options = append(options, tcOption{Interpretation: vtInt64, Type: tcaTaPrioSchedCycleTimeExtension, Data: int64Value(info.SchedCycleTimeExtension)})
	}
	if info.Flags != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTaPrioFlags, Data: uint32Value(info.Flags)})
	}
	if info.TxTimeDelay != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTaPrioTxTimeDelay, Data: uint32Value(info.TxTimeDelay)})
	}

	if multiError != nil {
		return []byte{}, multiError
	}
	return marshalAttributes(options)
}
