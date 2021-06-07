package util

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	out io.Writer = ioutil.Discard
	// If out is set by SetOutputFile, outFile is set and keeps the same value
	// as out. Otherwise, outFile is nil.
	outFile *os.File
	loggers []*log.Logger
)

// GetLogger gets a logger with a prefix.
func GetLogger(prefix string) *log.Logger {
	logger := log.New(out, prefix, log.LstdFlags)
	loggers = append(loggers, logger)
	return logger
}

// SetOutput redirects the output of all loggers obtained with GetLogger to the
// new io.Writer. If the old output was a file opened by SetOutputFile, it is
// closed.
func SetOutput(newout io.Writer) {
	if outFile != nil {
		outFile.Close()
		outFile = nil
	}
	out = newout
	outFile = nil
	for _, logger := range loggers {
		logger.SetOutput(out)
	}
}

// SetOutputFile redirects the output of all loggers obtained with GetLogger to
// the named file. If the old output was a file opened by SetOutputFile, it is
// closed. The new file is truncated. SetOutFile("") is equivalent to
// SetOutput(ioutil.Discard).
func SetOutputFile(fname string) error {
	if fname == "" {
		SetOutput(ioutil.Discard)
		return nil
	}
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	SetOutput(file)
	outFile = file
	return nil
}
