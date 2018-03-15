package wifi

type StubWifiWorker struct{}

func (StubWifiWorker) ScanInterfaces() ([]string, error) {
	return nil, nil
}

func (StubWifiWorker) ScanWifi() ([]WifiOption, error) {
	return nil, nil
}

func (StubWifiWorker) ScanCurrentWifi() (string, error) {
	return "", nil
}

func (StubWifiWorker) Connect(a ...string) error {
	return nil
}
