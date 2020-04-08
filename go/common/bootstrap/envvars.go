package bootstrap

import (
	"log"
	"os"
	"strconv"
)

func GetIntEnvVar(key string) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("missing required env var %v", key)
	}

	var err error
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("cannot parse %v, error: %v", key, err)
	}

	log.Printf("%v set to %v", key, value)

	return result
}

func GetEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("missing required env var %v", key)
	}

	log.Printf("%v set to %v", key, value)

	return value
}

func GetOptionalEnvVar(key string, def string) string {
	strValue, exists := os.LookupEnv(key)
	result := def
	if exists {
		result = strValue
	}

	log.Printf("%v set to %v", key, result)

	return result
}

func GetOptionalBoolEnvVar(key string, def bool) bool {
	strValue, exists := os.LookupEnv(key)
	result := def
	if exists {
		var err error
		result, err = strconv.ParseBool(strValue)
		if err != nil {
			log.Panicf("cannot parse %v, error: %v", key, err)
		}
	}

	log.Printf("%v set to %v", key, result)

	return result
}

func GetBoolEnvVar(key string) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("missing required env var %v", key)
	}

	var err error
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Panicf("cannot parse %v, error: %v", key, err)
	}

	log.Printf("%v set to %v", key, value)

	return result
}

func GetOptionalIntEnvVar(key string, def int) int {
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

func GetOptionalFloatEnvVar(key string, def float64) float64 {
	strValue, exists := os.LookupEnv(key)
	result := def
	if exists {
		var err error
		result, err = strconv.ParseFloat(strValue, 64)
		if err != nil {
			log.Panicf("cannot parse %v, error: %v", key, err)
		}
	}

	log.Printf("%v set to %v", key, result)

	return result
}
