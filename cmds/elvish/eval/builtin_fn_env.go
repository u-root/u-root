package eval

import (
	"errors"
	"os"
)

var errNonExistentEnvVar = errors.New("non-existent environment variable")

func init() {
	addBuiltinFns(map[string]interface{}{
		"has-env":   hasEnv,
		"get-env":   getEnv,
		"set-env":   os.Setenv,
		"unset-env": os.Unsetenv,
	})
}

func hasEnv(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func getEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", errNonExistentEnvVar
	}
	return value, nil
}
