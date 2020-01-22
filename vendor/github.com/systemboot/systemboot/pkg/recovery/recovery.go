package recovery

// Recoverer interface offers recovering
// from critical errors in different ways.
// Currently permissiverecoverer with log
// output and securerecovery with shutdown
// capabilites are supported.
type Recoverer interface {
	Recover(message string) error
}
