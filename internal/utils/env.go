package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvString(varName string, defaultValue string) string {
	val := os.Getenv(varName)
	if val == "" {
		return defaultValue
	}

	return val
}

func GetEnvInt(varName string, defaultValue int) int {
	valStr := GetEnvString(varName, "")
	if valStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valStr)
	if err != nil {
		log.Print(err)
		return defaultValue
	}

	return value
}

func GetEnvBool(varName string, defaultValue bool) bool {
	valStr := GetEnvString(varName, "")
	if valStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valStr)
	if err != nil {
		log.Print(err)
		return defaultValue
	}

	return value
}
