//+build plan9 windows

package nltest

func isSyscallError(_ error) bool {
	return false
}
