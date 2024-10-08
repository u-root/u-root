package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"

// wrapper for go callback function
extern int violationNotify(void* p);
*/
import "C"
import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"
)

type policyCondition string

const (
	DbePolicy     = policyCondition("Double-bit ECC error")
	PCIePolicy    = policyCondition("PCI error")
	MaxRtPgPolicy = policyCondition("Max Retired Pages Limit")
	ThermalPolicy = policyCondition("Thermal Limit")
	PowerPolicy   = policyCondition("Power Limit")
	NvlinkPolicy  = policyCondition("Nvlink Error")
	XidPolicy     = policyCondition("XID Error")
)

type PolicyViolation struct {
	Condition policyCondition
	Timestamp time.Time
	Data      interface{}
}

type policyIndex int

const (
	dbePolicyIndex policyIndex = iota
	pciePolicyIndex
	maxRtPgPolicyIndex
	thermalPolicyIndex
	powerPolicyIndex
	nvlinkPolicyIndex
	xidPolicyIndex
)

type policyConditionParam struct {
	typ   uint32
	value uint32
}

type DbePolicyCondition struct {
	Location  string
	NumErrors uint
}

type PciPolicyCondition struct {
	ReplayCounter uint
}

type RetiredPagesPolicyCondition struct {
	SbePages uint
	DbePages uint
}

type ThermalPolicyCondition struct {
	ThermalViolation uint
}

type PowerPolicyCondition struct {
	PowerViolation uint
}

type NvlinkPolicyCondition struct {
	FieldId uint16
	Counter uint
}

type XidPolicyCondition struct {
	ErrNum uint
}

var (
	policyChanOnce sync.Once
	policyMapOnce  sync.Once

	// callbacks maps PolicyViolation channels with policy
	// captures C callback() value for each violation condition
	callbacks map[string]chan PolicyViolation

	// paramMap maps C.dcgmPolicy_t.parms index and limits
	// to be used in setPolicy() for setting user selected policies
	paramMap map[policyIndex]policyConditionParam
)

func makePolicyChannels() {
	policyChanOnce.Do(func() {
		callbacks = make(map[string]chan PolicyViolation)
		callbacks["dbe"] = make(chan PolicyViolation, 1)
		callbacks["pcie"] = make(chan PolicyViolation, 1)
		callbacks["maxrtpg"] = make(chan PolicyViolation, 1)
		callbacks["thermal"] = make(chan PolicyViolation, 1)
		callbacks["power"] = make(chan PolicyViolation, 1)
		callbacks["nvlink"] = make(chan PolicyViolation, 1)
		callbacks["xid"] = make(chan PolicyViolation, 1)
	})
}

func makePolicyParmsMap() {
	const (
		policyFieldTypeBool    = 0
		policyFieldTypeLong    = 1
		policyBoolValue        = 1
		policyMaxRtPgThreshold = 10
		policyThermalThreshold = 100
		policyPowerThreshold   = 250
	)

	policyMapOnce.Do(func() {
		paramMap = make(map[policyIndex]policyConditionParam)
		paramMap[dbePolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeBool,
			value: policyBoolValue,
		}

		paramMap[pciePolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeBool,
			value: policyBoolValue,
		}

		paramMap[maxRtPgPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeLong,
			value: policyMaxRtPgThreshold,
		}

		paramMap[thermalPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeLong,
			value: policyThermalThreshold,
		}

		paramMap[powerPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeLong,
			value: policyPowerThreshold,
		}

		paramMap[nvlinkPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeBool,
			value: policyBoolValue,
		}

		paramMap[xidPolicyIndex] = policyConditionParam{
			typ:   policyFieldTypeBool,
			value: policyBoolValue,
		}
	})
}

// ViolationRegistration is a go callback function for dcgmPolicyRegister() wrapped in C.violationNotify()
//
//export ViolationRegistration
func ViolationRegistration(data unsafe.Pointer) int {
	var con policyCondition
	var timestamp time.Time
	var val interface{}

	response := *(*C.dcgmPolicyCallbackResponse_t)(unsafe.Pointer(data))

	switch response.condition {
	case C.DCGM_POLICY_COND_DBE:
		dbe := (*C.dcgmPolicyConditionDbe_t)(unsafe.Pointer(&response.val))
		con = DbePolicy
		timestamp = createTimeStamp(dbe.timestamp)
		val = DbePolicyCondition{
			Location:  dbeLocation(int(dbe.location)),
			NumErrors: *uintPtr(dbe.numerrors),
		}
	case C.DCGM_POLICY_COND_PCI:
		pci := (*C.dcgmPolicyConditionPci_t)(unsafe.Pointer(&response.val))
		con = PCIePolicy
		timestamp = createTimeStamp(pci.timestamp)
		val = PciPolicyCondition{
			ReplayCounter: *uintPtr(pci.counter),
		}
	case C.DCGM_POLICY_COND_MAX_PAGES_RETIRED:
		mpr := (*C.dcgmPolicyConditionMpr_t)(unsafe.Pointer(&response.val))
		con = MaxRtPgPolicy
		timestamp = createTimeStamp(mpr.timestamp)
		val = RetiredPagesPolicyCondition{
			SbePages: *uintPtr(mpr.sbepages),
			DbePages: *uintPtr(mpr.dbepages),
		}
	case C.DCGM_POLICY_COND_THERMAL:
		thermal := (*C.dcgmPolicyConditionThermal_t)(unsafe.Pointer(&response.val))
		con = ThermalPolicy
		timestamp = createTimeStamp(thermal.timestamp)
		val = ThermalPolicyCondition{
			ThermalViolation: *uintPtr(thermal.thermalViolation),
		}
	case C.DCGM_POLICY_COND_POWER:
		pwr := (*C.dcgmPolicyConditionPower_t)(unsafe.Pointer(&response.val))
		con = PowerPolicy
		timestamp = createTimeStamp(pwr.timestamp)
		val = PowerPolicyCondition{
			PowerViolation: *uintPtr(pwr.powerViolation),
		}
	case C.DCGM_POLICY_COND_NVLINK:
		nvlink := (*C.dcgmPolicyConditionNvlink_t)(unsafe.Pointer(&response.val))
		con = NvlinkPolicy
		timestamp = createTimeStamp(nvlink.timestamp)
		val = NvlinkPolicyCondition{
			FieldId: uint16(nvlink.fieldId),
			Counter: *uintPtr(nvlink.counter),
		}
	case C.DCGM_POLICY_COND_XID:
		xid := (*C.dcgmPolicyConditionXID_t)(unsafe.Pointer(&response.val))
		con = XidPolicy
		timestamp = createTimeStamp(xid.timestamp)
		val = XidPolicyCondition{
			ErrNum: *uintPtr(xid.errnum),
		}
	}

	err := PolicyViolation{
		Condition: con,
		Timestamp: timestamp,
		Data:      val,
	}

	switch con {
	case DbePolicy:
		callbacks["dbe"] <- err
	case PCIePolicy:
		callbacks["pcie"] <- err
	case MaxRtPgPolicy:
		callbacks["maxrtpg"] <- err
	case ThermalPolicy:
		callbacks["thermal"] <- err
	case PowerPolicy:
		callbacks["power"] <- err
	case NvlinkPolicy:
		callbacks["nvlink"] <- err
	case XidPolicy:
		callbacks["xid"] <- err
	}
	return 0
}

func setPolicy(groupId GroupHandle, condition C.dcgmPolicyCondition_t, paramList []policyIndex) (err error) {
	var policy C.dcgmPolicy_t
	policy.version = makeVersion1(unsafe.Sizeof(policy))
	policy.mode = C.dcgmPolicyMode_t(C.DCGM_OPERATION_MODE_AUTO)
	policy.action = C.DCGM_POLICY_ACTION_NONE
	policy.isolation = C.DCGM_POLICY_ISOLATION_NONE
	policy.validation = C.DCGM_POLICY_VALID_NONE
	policy.condition = condition

	// iterate on paramMap for given policy conditions
	for _, key := range paramList {
		conditionParam, exists := paramMap[policyIndex(key)]
		if !exists {
			return fmt.Errorf("Error: Invalid Policy condition, %v does not exist", key)
		}
		// set policy condition parameters
		// set condition type (bool or longlong)
		policy.parms[key].tag = conditionParam.typ

		// set condition val (violation threshold)
		// policy.parms.val is a C union type
		// cgo docs: Go doesn't have support for C's union type
		// C union types are represented as a Go byte array
		binary.LittleEndian.PutUint32(policy.parms[key].val[:], conditionParam.value)
	}

	var statusHandle C.dcgmStatus_t

	result := C.dcgmPolicySet(handle.handle, groupId.handle, &policy, statusHandle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("Error setting policies: %s", err)
	}

	log.Println("Policy successfully set.")

	return
}

func registerPolicy(ctx context.Context, groupId GroupHandle, typ ...policyCondition) (<-chan PolicyViolation, error) {
	// init policy globals for internal API
	makePolicyChannels()
	makePolicyParmsMap()

	// make a list of policy conditions for setting their parameters
	var paramKeys []policyIndex
	// get all conditions to be set in setPolicy()
	var condition C.dcgmPolicyCondition_t = 0
	for _, t := range typ {
		switch t {
		case DbePolicy:
			paramKeys = append(paramKeys, dbePolicyIndex)
			condition |= C.DCGM_POLICY_COND_DBE
		case PCIePolicy:
			paramKeys = append(paramKeys, pciePolicyIndex)
			condition |= C.DCGM_POLICY_COND_PCI
		case MaxRtPgPolicy:
			paramKeys = append(paramKeys, maxRtPgPolicyIndex)
			condition |= C.DCGM_POLICY_COND_MAX_PAGES_RETIRED
		case ThermalPolicy:
			paramKeys = append(paramKeys, thermalPolicyIndex)
			condition |= C.DCGM_POLICY_COND_THERMAL
		case PowerPolicy:
			paramKeys = append(paramKeys, powerPolicyIndex)
			condition |= C.DCGM_POLICY_COND_POWER
		case NvlinkPolicy:
			paramKeys = append(paramKeys, nvlinkPolicyIndex)
			condition |= C.DCGM_POLICY_COND_NVLINK
		case XidPolicy:
			paramKeys = append(paramKeys, xidPolicyIndex)
			condition |= C.DCGM_POLICY_COND_XID
		}
	}

	var err error
	if err = setPolicy(groupId, condition, paramKeys); err != nil {
		return nil, err
	}

	var finishCallback unsafe.Pointer
	result := C.dcgmPolicyRegister(handle.handle, groupId.handle, C.dcgmPolicyCondition_t(condition), C.fpRecvUpdates(C.violationNotify), C.fpRecvUpdates(finishCallback))

	if err = errorString(result); err != nil {
		return nil, &DcgmError{msg: C.GoString(C.errorString(result)), Code: result}
	}
	log.Println("Listening for violations...")

	violation := make(chan PolicyViolation, len(typ))
	go func() {
		defer func() {
			log.Println("unregister policy violation...")
			close(violation)
			unregisterPolicy(groupId, condition)
		}()
		for {
			select {
			case dbe := <-callbacks["dbe"]:
				violation <- dbe
			case pcie := <-callbacks["pcie"]:
				violation <- pcie
			case maxrtpg := <-callbacks["maxrtpg"]:
				violation <- maxrtpg
			case thermal := <-callbacks["thermal"]:
				violation <- thermal
			case power := <-callbacks["power"]:
				violation <- power
			case nvlink := <-callbacks["nvlink"]:
				violation <- nvlink
			case xid := <-callbacks["xid"]:
				violation <- xid
			case <-ctx.Done():
				return
			}
		}
	}()

	return violation, err
}

func unregisterPolicy(groupId GroupHandle, condition C.dcgmPolicyCondition_t) {
	result := C.dcgmPolicyUnregister(handle.handle, groupId.handle, condition)

	if err := errorString(result); err != nil {
		log.Println(fmt.Errorf("error unregistering policy: %s", err))
	}
}

func createTimeStamp(t C.longlong) time.Time {
	tm := int64(t) / 1000000
	ts := time.Unix(tm, 0)
	return ts
}

func dbeLocation(location int) string {
	switch location {
	case 0:
		return "L1"
	case 1:
		return "L2"
	case 2:
		return "Device"
	case 3:
		return "Register"
	case 4:
		return "Texture"
	}
	return "N/A"
}
