package util

import "os"

func GetEnvWithDefault(name, fallback string) string {
	if envVar, ok := os.LookupEnv(name); ok {
		return envVar
	}
	return fallback
}
