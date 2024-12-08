package config

import (
	"os"
	"strings"
)

func GetValue(in string) string {
	if strings.HasPrefix("env", in) {
		envVar := in[4:]
		return os.Getenv(envVar)
	}

	return in
}

func GetValueWithDefault(in string, def string) string {
	result := GetValue(in)
	if result == "" {
		return def
	}

	return result
}
