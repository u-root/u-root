package tc

import (
	"fmt"

	"github.com/florianl/go-tc/internal/unix"
	"github.com/mdlayher/netlink"
)

// Actions represents the actions part of rtnetlink
type Actions struct {
	Tc
}

// tcamsg is Actions specific
type tcaMsg struct {
	Family uint8
	Pad1   uint8
	Pad2   uint16
}

// Actions allows to read and alter actions
func (tc *Tc) Actions() *Actions {
	return &Actions{*tc}
}

// Add creates a new actions
func (a *Actions) Add(info []*Action) error {
	if len(info) == 0 {
		return ErrNoArg
	}
	options, err := validateActionsObject(unix.RTM_NEWACTION, info)
	if err != nil {
		return err
	}
	return a.action(unix.RTM_NEWACTION, netlink.Create|netlink.Excl, tcaMsg{
		Family: unix.AF_UNSPEC,
	}, options)
}

// Replace add/remove an actions. If the node does not exist yet it is created
func (a *Actions) Replace(info []*Action) error {
	if len(info) == 0 {
		return ErrNoArg
	}
	options, err := validateActionsObject(unix.RTM_NEWACTION, info)
	if err != nil {
		return err
	}
	return a.action(unix.RTM_NEWACTION, netlink.Create, tcaMsg{
		Family: unix.AF_UNSPEC,
	}, options)
}

// Delete removes an actions
func (a *Actions) Delete(info []*Action) error {
	if len(info) == 0 {
		return ErrNoArg
	}
	options, err := validateActionsObject(unix.RTM_DELACTION, info)
	if err != nil {
		return err
	}
	return a.action(unix.RTM_DELACTION, netlink.HeaderFlags(0), tcaMsg{
		Family: unix.AF_UNSPEC,
	}, options)
}

// Get fetches all actions
func (a *Actions) Get(actions []*Action) ([]*Action, error) {
	var results []*Action
	var data []byte
	tcminfo, err := marshalStruct(tcaMsg{
		Family: unix.AF_UNSPEC,
	})
	if err != nil {
		return results, err
	}

	data = append(data, tcminfo...)
	options, err := validateActionsObject(unix.RTM_GETACTION, actions)
	if err != nil {
		return results, err
	}

	attrs, err := marshalAttributes(options)
	if err != nil {
		return results, err
	}
	data = append(data, attrs...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(unix.RTM_GETACTION),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := a.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		// The first 4 bytes contain tcaMsg - which is skipped here.
		if err := unmarshalRoot(msg.Data[4:], &results); err != nil {
			return results, err
		}
	}
	return results, nil
}

func validateActionsObject(cmd int, info []*Action) ([]tcOption, error) {
	options := []tcOption{}

	data, err := marshalActions(cmd, info)
	if err != nil {
		return options, err
	}
	options = append(options, tcOption{Interpretation: vtBytes, Type: 1 /*TCA_ROOT_TAB*/, Data: data})

	return options, nil
}

const (
	tcaRootUnspec = iota
	tcaRootTab
	tcaRootFlags
	tcaRootCount
	tcaRootTimeDelta
	tcaRootExtWarnMsg
)

func unmarshalRoot(data []byte, actions *[]*Action) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var multiError error
	for ad.Next() {
		switch ad.Type() {
		case tcaRootTab:
			err := unmarshalActions(ad.Bytes(), actions)
			multiError = concatError(multiError, err)
		case tcaRootFlags:
			_ = ad.Uint64()
		case tcaRootCount:
			_ = ad.Uint32()
		case tcaRootTimeDelta:
			_ = ad.Uint32()
		case tcaRootExtWarnMsg:
			_ = ad.String()
		default:
			return fmt.Errorf("unmarshalRoot()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return concatError(multiError, ad.Err())
}
