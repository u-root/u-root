package dcgm

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	dcgmInitCounter int
	mux             sync.Mutex
)

// Init starts DCGM, based on the user selected mode
// DCGM can be started in 3 differengt modes:
// 1. Embedded: Start hostengine within this process
// 2. Standalone: Connect to an already running nv-hostengine at the specified address
// Connection address can be passed as command line args: -connect "IP:PORT/Socket" -socket "isSocket"
// 3. StartHostengine: Open an Unix socket to start and connect to the nv-hostengine and terminate before exiting
func Init(m mode, args ...string) (cleanup func(), err error) {
	mux.Lock()
	if dcgmInitCounter < 0 {
		count := fmt.Sprintf("%d", dcgmInitCounter)
		err = fmt.Errorf("Shutdown() is called %s times, before Init()", count[1:])
	}
	if dcgmInitCounter == 0 {
		err = initDcgm(m, args...)

		if err != nil {
			mux.Unlock()

			return nil, err
		}
	}
	dcgmInitCounter += 1
	mux.Unlock()

	return func() {
		if err := Shutdown(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to shutdown DCGM with error: `%v`", err)
		}
	}, err
}

// Shutdown stops DCGM and destroy all connections
func Shutdown() (err error) {
	mux.Lock()
	if dcgmInitCounter <= 0 {
		err = fmt.Errorf("Init() needs to be called before Shutdown()")
	}
	if dcgmInitCounter == 1 {
		err = shutdown()
	}
	dcgmInitCounter -= 1
	mux.Unlock()

	return
}

// GetAllDeviceCount counts all GPUs on the system
func GetAllDeviceCount() (uint, error) {
	return getAllDeviceCount()
}

func GetEntityGroupEntities(entityGroup Field_Entity_Group) ([]uint, error) {
	return getEntityGroupEntities(entityGroup)
}

// GetSupportedDevices returns only DCGM supported GPUs
func GetSupportedDevices() ([]uint, error) {
	return getSupportedDevices()
}

// GetDeviceInfo describes the given device
func GetDeviceInfo(gpuId uint) (Device, error) {
	return getDeviceInfo(gpuId)
}

// GetDeviceStatus monitors GPU status including its power, memory and GPU utilization
func GetDeviceStatus(gpuId uint) (DeviceStatus, error) {
	return latestValuesForDevice(gpuId)
}

// GetDeviceTopology returns device topology corresponding to the gpuId
func GetDeviceTopology(gpuId uint) ([]P2PLink, error) {
	return getDeviceTopology(gpuId)
}

// WatchPidFields lets DCGM start recording stats for GPU process
// It needs to be called before calling GetProcessInfo
func WatchPidFields() (GroupHandle, error) {
	return watchPidFields(time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples)
}

// GetProcessInfo provides detailed per GPU stats for this process
func GetProcessInfo(group GroupHandle, pid uint) ([]ProcessInfo, error) {
	return getProcessInfo(group, pid)
}

// HealthCheckByGpuId monitors GPU health for any errors/failures/warnings
func HealthCheckByGpuId(gpuId uint) (DeviceHealth, error) {
	return healthCheckByGpuId(gpuId)
}

// ListenForPolicyViolations sets GPU usage and error policies and notifies in case of any violations
func ListenForPolicyViolations(ctx context.Context, typ ...policyCondition) (<-chan PolicyViolation, error) {
	groupId := GroupAllGPUs()
	return registerPolicy(ctx, groupId, typ...)
}

// Introspect returns DCGM hostengine memory and CPU usage
func Introspect() (DcgmStatus, error) {
	return introspect()
}

// Get all of the profiling metric groups for a given GPU group.
func GetSupportedMetricGroups(gpuId uint) ([]MetricGroup, error) {
	return getSupportedMetricGroups(gpuId)
}

func GetNvLinkLinkStatus() ([]NvLinkStatus, error) {
	return getNvLinkLinkStatus()
}
