package wifi

type StubWifiWorker struct {
	ScanInterfacesOut  []string
	ScanWifiOut        []WifiOption
	ScanCurrentWifiOut string
}

func (stub StubWifiWorker) ScanInterfaces() ([]string, error) {
	return stub.ScanInterfacesOut, nil
}

func (stub StubWifiWorker) ScanWifi() ([]WifiOption, error) {
	return stub.ScanWifiOut, nil
}

func (stub StubWifiWorker) ScanCurrentWifi() (string, error) {
	return stub.ScanCurrentWifiOut, nil
}

func (StubWifiWorker) Connect(a ...string) error {
	return nil
}
