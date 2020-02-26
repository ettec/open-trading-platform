package bootstrap

import (
	"log"
	"os"
	"strconv"
)

func GetBootstrapEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("missing required env var %v", key)
	}

	log.Printf("%v set to %v", key, value)

	return value
}

func GetOptionalBootstrapIntEnvVar(key string, def int) int {
	strValue, exists := os.LookupEnv(key)
	result := def
	if exists {
		var err error
		result, err = strconv.Atoi(strValue)
		if err != nil {
			log.Panicf("cannot parse %v, error: %v", key, err)
		}
	}

	log.Printf("%v set to %v", key, result)

	return result
}