package tc

import (
	"fmt"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Qdisc represents the queueing discipline part of traffic control
type Qdisc struct {
	Tc
}

// Qdisc allows to read and alter queues
func (tc *Tc) Qdisc() *Qdisc {
	return &Qdisc{*tc}
}

// Add creates a new queueing discipline
func (qd *Qdisc) Add(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_NEWQDISC, info)
	if err != nil {
		return err
	}
	return qd.action(unix.RTM_NEWQDISC, netlink.Create|netlink.Excl, &info.Msg, options)
}

// Replace add/remove a queueing discipline. If the node does not exist yet it is created
func (qd *Qdisc) Replace(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_NEWQDISC, info)
	if err != nil {
		return err
	}
	return qd.action(unix.RTM_NEWQDISC, netlink.Create|netlink.Replace, &info.Msg, options)
}

// Link performs a replace on an existing queueing discipline
func (qd *Qdisc) Link(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_NEWQDISC, info)
	if err != nil {
		return err
	}
	return qd.action(unix.RTM_NEWQDISC, netlink.Replace, &info.Msg, options)
}

// Delete removes a queueing discipline
func (qd *Qdisc) Delete(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_DELQDISC, info)
	if err != nil {
		return err
	}
	return qd.action(unix.RTM_DELQDISC, netlink.HeaderFlags(0), &info.Msg, options)
}

// Change modifies a queueing discipline 'in place'
func (qd *Qdisc) Change(info *Object) error {
	if info == nil {
		return ErrNoArg
	}
	options, err := validateQdiscObject(unix.RTM_NEWQDISC, info)
	if err != nil {
		return err
	}
	return qd.action(unix.RTM_NEWQDISC, netlink.HeaderFlags(0), &info.Msg, options)
}

// Get fetches all queueing disciplines
func (qd *Qdisc) Get() ([]Object, error) {
	return qd.get(unix.RTM_GETQDISC, &Msg{})
}

func validateQdiscObject(action int, info *Object) ([]tcOption, error) {
	options := []tcOption{}
	if info.Ifindex == 0 {
		return options, ErrInvalidDev
	}

	// TODO: improve logic and check combinations
	var data []byte
	var err error
	switch info.Kind {
	case "cbs":
		data, err = marshalCbs(info.Cbs)
	case "cake":
		data, err = marshalCake(info.Cake)
	case "choke":
		data, err = marshalChoke(info.Choke)
	case "pfifo":
		data, err = marshalStruct(info.Pfifo)
	case "bfifo":
		data, err = marshalStruct(info.Bfifo)
	case "tbf":
		data, err = marshalTbf(info.Tbf)
	case "sfb":
		data, err = marshalSfb(info.Sfb)
	case "sfq":
		data, err = marshalSfq(info.Sfq)
	case "red":
		data, err = marshalRed(info.Red)
	case "qfq":
		// qfq is parameterless
		// parameters are used in its corresponding class
	case "pie":
		data, err = marshalPie(info.Pie)
	case "mqprio":
		data, err = marshalMqPrio(info.MqPrio)
	case "hhf":
		data, err = marshalHhf(info.Hhf)
	case "hfsc":
		data, err = marshalHfscQOpt(info.HfscQOpt)
	case "fq":
		data, err = marshalFq(info.Fq)
	case "dsmark":
		data, err = marshalDsmark(info.Dsmark)
	case "drr":
		data, err = marshalDrr(info.Drr)
	case "codel":
		data, err = marshalCodel(info.Codel)
	case "cbq":
		data, err = marshalCbq(info.Cbq)
	case "atm":
		data, err = marshalAtm(info.Atm)
	case "fq_codel":
		data, err = marshalFqCodel(info.FqCodel)
	case "htb":
		data, err = marshalHtb(info.Htb)
	case "netem":
		data, err = marshalNetem(info.Netem)
	case "prio":
		data, err = marshalPrio(info.Prio)
	case "plug":
		data, err = marshalPlug(info.Plug)
	case "taprio":
		data, err = marshalTaPrio(info.TaPrio)
	case "clsact":
		// clsact is parameterless
	case "ingress":
		// ingress is parameterless
	default:
		return options, fmt.Errorf("%s: %w", info.Kind, ErrNotImplemented)
	}
	if err != nil {
		return options, err
	}
	if len(data) < 1 && action == unix.RTM_NEWQDISC {
		if info.Kind != "clsact" && info.Kind != "ingress" && info.Kind != "qfq" {
			return options, ErrNoArg
		}
	} else {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: data})
	}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaKind, Data: info.Kind})

	if (info.Stats != nil || info.XStats != nil || info.Stats2 != nil) && action != unix.RTM_DELQDISC {
		return options, ErrNotImplemented
	}

	if info.EgressBlock != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaEgressBlock, Data: uint32Value(info.EgressBlock)})
	}
	if info.IngressBlock != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIngressBlock, Data: uint32Value(info.IngressBlock)})
	}
	if info.HwOffload != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaHwOffload, Data: uint8Value(info.HwOffload)})
	}
	if info.Chain != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaChain, Data: uint32Value(info.Chain)})
	}
	if info.Stab != nil {
		data, err := marshalStab(info.Stab)
		if err != nil {
			return options, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaStab, Data: data})

	}
	return options, nil
}
