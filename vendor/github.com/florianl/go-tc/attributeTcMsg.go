package tc

import (
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

func extractTcmsgAttributes(action int, data []byte, info *Attribute) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var options []byte
	var xStats []byte
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaKind:
			info.Kind = ad.String()
		case tcaOptions:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is known at this moment,
			// so we save it for later
			options = ad.Bytes()
		case tcaChain:
			info.Chain = uint32Ptr(ad.Uint32())
		case tcaXstats:
			// the evaluation of this field depends on tcaKind.
			// there is no guarantee, that kind is know at this moment,
			// so we save it for later
			xStats = ad.Bytes()
		case tcaStats:
			tcstats := &Stats{}
			err := unmarshalStruct(ad.Bytes(), tcstats)
			multiError = concatError(multiError, err)
			info.Stats = tcstats
		case tcaStats2:
			tcstats2 := &Stats2{}
			err := unmarshalStruct(ad.Bytes(), tcstats2)
			multiError = concatError(multiError, err)
			info.Stats2 = tcstats2
		case tcaHwOffload:
			info.HwOffload = uint8Ptr(ad.Uint8())
		case tcaEgressBlock:
			info.EgressBlock = uint32Ptr(ad.Uint32())
		case tcaIngressBlock:
			info.IngressBlock = uint32Ptr(ad.Uint32())
		case tcaStab:
			stab := &Stab{}
			err := unmarshalStab(ad.Bytes(), stab)
			multiError = concatError(multiError, err)
			info.Stab = stab
		case tcaPad:
			// padding does not contain data, we just skip it
		case tcaExtWarnMsg:
			info.ExtWarnMsg = ad.String()
		default:
			return fmt.Errorf("extractTcmsgAttributes()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}

	if err = concatError(multiError, ad.Err()); err != nil {
		return err
	}

	if len(options) > 0 {
		if (action&actionMask == actionQdisc) && hasQOpt(info.Kind) {
			err = extractQOpt(options, info, info.Kind)
		} else {
			err = extractTCAOptions(options, info, info.Kind)
		}
		multiError = concatError(multiError, err)
	}

	if len(xStats) > 0 {
		tcxstats := &XStats{}
		err := extractXStats(xStats, tcxstats, info.Kind)
		multiError = concatError(multiError, err)
		info.XStats = tcxstats
	}
	return multiError
}

func hasQOpt(kind string) bool {
	classful := map[string]bool{
		"hfsc": true,
		"qfq":  true,
		"htb":  true,
	}
	if _, ok := classful[kind]; ok {
		return true
	}
	return false
}

func extractQOpt(data []byte, tc *Attribute, kind string) error {
	var multiError error
	switch kind {
	case "hfsc":
		info := &HfscQOpt{}
		err := unmarshalHfscQOpt(data, info)
		multiError = concatError(multiError, err)
		tc.HfscQOpt = info
	case "qfq":
		info := &Qfq{}
		err := unmarshalQfq(data, info)
		multiError = concatError(multiError, err)
		tc.Qfq = info
	case "htb":
		info := &Htb{}
		err := unmarshalHtb(data, info)
		multiError = concatError(multiError, err)
		tc.Htb = info
	default:
		return fmt.Errorf("no QOpts for %s", kind)
	}
	return multiError
}

func extractTCAOptions(data []byte, tc *Attribute, kind string) error {
	var multiError error
	switch kind {
	case "choke":
		info := &Choke{}
		err := unmarshalChoke(data, info)
		multiError = concatError(multiError, err)
		tc.Choke = info
	case "fq_codel":
		info := &FqCodel{}
		err := unmarshalFqCodel(data, info)
		multiError = concatError(multiError, err)
		tc.FqCodel = info
	case "codel":
		info := &Codel{}
		err := unmarshalCodel(data, info)
		multiError = concatError(multiError, err)
		tc.Codel = info
	case "fq":
		info := &Fq{}
		err := unmarshalFq(data, info)
		multiError = concatError(multiError, err)
		tc.Fq = info
	case "pie":
		info := &Pie{}
		err := unmarshalPie(data, info)
		multiError = concatError(multiError, err)
		tc.Pie = info
	case "hhf":
		info := &Hhf{}
		err := unmarshalHhf(data, info)
		multiError = concatError(multiError, err)
		tc.Hhf = info
	case "htb":
		info := &Htb{}
		err := unmarshalHtb(data, info)
		multiError = concatError(multiError, err)
		tc.Htb = info
	case "hfsc":
		info := &Hfsc{}
		err := unmarshalHfsc(data, info)
		multiError = concatError(multiError, err)
		tc.Hfsc = info
	case "dsmark":
		info := &Dsmark{}
		err := unmarshalDsmark(data, info)
		multiError = concatError(multiError, err)
		tc.Dsmark = info
	case "drr":
		info := &Drr{}
		err := unmarshalDrr(data, info)
		multiError = concatError(multiError, err)
		tc.Drr = info
	case "cbq":
		info := &Cbq{}
		err := unmarshalCbq(data, info)
		multiError = concatError(multiError, err)
		tc.Cbq = info
	case "atm":
		info := &Atm{}
		err := unmarshalAtm(data, info)
		multiError = concatError(multiError, err)
		tc.Atm = info
	case "pfifo_fast":
		fallthrough
	case "prio":
		info := &Prio{}
		err := unmarshalPrio(data, info)
		multiError = concatError(multiError, err)
		tc.Prio = info
	case "tbf":
		info := &Tbf{}
		err := unmarshalTbf(data, info)
		multiError = concatError(multiError, err)
		tc.Tbf = info
	case "sfb":
		info := &Sfb{}
		err := unmarshalSfb(data, info)
		multiError = concatError(multiError, err)
		tc.Sfb = info
	case "sfq":
		info := &Sfq{}
		err := unmarshalSfq(data, info)
		multiError = concatError(multiError, err)
		tc.Sfq = info
	case "red":
		info := &Red{}
		err := unmarshalRed(data, info)
		multiError = concatError(multiError, err)
		tc.Red = info
	case "pfifo":
		limit := &FifoOpt{}
		err := unmarshalStruct(data, limit)
		multiError = concatError(multiError, err)
		tc.Pfifo = limit
	case "mqprio":
		info := &MqPrio{}
		err := unmarshalMqPrio(data, info)
		multiError = concatError(multiError, err)
		tc.MqPrio = info
	case "bfifo":
		limit := &FifoOpt{}
		err := unmarshalStruct(data, limit)
		multiError = concatError(multiError, err)
		tc.Bfifo = limit
	case "clsact":
		return extractClsact(data)
	case "ingress":
		return extractIngress(data)
	case "qfq":
		info := &Qfq{}
		err := unmarshalQfq(data, info)
		multiError = concatError(multiError, err)
		tc.Qfq = info
	case "basic":
		info := &Basic{}
		err := unmarshalBasic(data, info)
		multiError = concatError(multiError, err)
		tc.Basic = info
	case "bpf":
		info := &Bpf{}
		err := unmarshalBpf(data, info)
		multiError = concatError(multiError, err)
		tc.BPF = info
	case "cgroup":
		info := &Cgroup{}
		err := unmarshalCgroup(data, info)
		multiError = concatError(multiError, err)
		tc.Cgroup = info
	case "u32":
		info := &U32{}
		err := unmarshalU32(data, info)
		multiError = concatError(multiError, err)
		tc.U32 = info
	case "flower":
		info := &Flower{}
		err := unmarshalFlower(data, info)
		multiError = concatError(multiError, err)
		tc.Flower = info
	case "rsvp":
		info := &Rsvp{}
		err := unmarshalRsvp(data, info)
		multiError = concatError(multiError, err)
		tc.Rsvp = info
	case "route4":
		info := &Route4{}
		err := unmarshalRoute4(data, info)
		multiError = concatError(multiError, err)
		tc.Route4 = info
	case "fw":
		info := &Fw{}
		err := unmarshalFw(data, info)
		multiError = concatError(multiError, err)
		tc.Fw = info
	case "flow":
		info := &Flow{}
		err := unmarshalFlow(data, info)
		multiError = concatError(multiError, err)
		tc.Flow = info
	case "matchall":
		info := &Matchall{}
		err := unmarshalMatchall(data, info)
		multiError = concatError(multiError, err)
		tc.Matchall = info
	case "netem":
		info := &Netem{}
		err := unmarshalNetem(data, info)
		multiError = concatError(multiError, err)
		tc.Netem = info
	case "cake":
		info := &Cake{}
		err := unmarshalCake(data, info)
		multiError = concatError(multiError, err)
		tc.Cake = info
	case "plug":
		info := &Plug{}
		err := unmarshalPlug(data, info)
		multiError = concatError(multiError, err)
		tc.Plug = info
	case "tcindex":
		info := &TcIndex{}
		err := unmarshalTcIndex(data, info)
		multiError = concatError(multiError, err)
		tc.TcIndex = info
	case "cbs":
		info := &Cbs{}
		err := unmarshalCbs(data, info)
		multiError = concatError(multiError, err)
		tc.Cbs = info
	case "taprio":
		info := &TaPrio{}
		err := unmarshalTaPrio(data, info)
		multiError = concatError(multiError, err)
		tc.TaPrio = info
	default:
		return fmt.Errorf("extractTCAOptions(): unsupported kind %s: %w", kind, ErrUnknownKind)
	}

	return multiError
}

func extractXStats(data []byte, tc *XStats, kind string) error {
	var multiError error
	switch kind {
	case "sfb":
		info := &SfbXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Sfb = info
	case "sfq":
		info := &SfqXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Sfq = info
	case "red":
		info := &RedXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Red = info
	case "choke":
		info := &ChokeXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Choke = info
	case "htb":
		info := &HtbXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Htb = info
	case "cbq":
		info := &CbqXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Cbq = info
	case "codel":
		info := &CodelXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Codel = info
	case "hhf":
		info := &HhfXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Hhf = info
	case "pie":
		info := &PieXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Pie = info
	case "fq_codel":
		info := &FqCodelXStats{}
		err := unmarshalFqCodelXStats(data, info)
		multiError = concatError(multiError, err)
		tc.FqCodel = info
	case "fq":
		info := &FqQdStats{}
		// Pad out data to size of our FqQdStats struct to handle
		// unmarshalling data from older kernel versions with smaller structs
		qd := make([]byte, binary.Size(info))
		copy(qd, data)
		err := unmarshalStruct(qd, info)
		multiError = concatError(multiError, err)
		tc.Fq = info
	case "hfsc":
		info := &HfscXStats{}
		err := unmarshalStruct(data, info)
		multiError = concatError(multiError, err)
		tc.Hfsc = info
	default:
		return fmt.Errorf("extractXStats(): unsupported kind: %s", kind)
	}
	return multiError
}

func extractClsact(data []byte) error {
	// Clsact is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("clsact is parameterless: %w", ErrInvalidArg)
	}
	return nil
}

func extractIngress(data []byte) error {
	// Ingress is parameterless - so we expect to options
	if len(data) != 0 {
		return fmt.Errorf("extractIngress()\t%v", data)
	}
	return nil
}

const (
	tcaUnspec = iota
	tcaKind
	tcaOptions
	tcaStats
	tcaXstats
	tcaRate
	tcaFcnt
	tcaStats2
	tcaStab
	tcaPad
	tcaDumpInvisible
	tcaChain
	tcaHwOffload
	tcaIngressBlock
	tcaEgressBlock
	tcaDumpFlags
	tcaExtWarnMsg
)
