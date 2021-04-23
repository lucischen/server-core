package env

import (
	"log"
	"os"
	"strconv"
)

func GetStr(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Env '%v' is empty", key)
	}

	return value
}

func GetInt(key string, greaterThanZero bool) int {
	v := GetStr(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("%v must be integer not %v. err: %v", key, v, err)
	}

	if greaterThanZero && i <= 0 {
		log.Fatalf("%v must greater than 0, got: %v", key, i)
	}

	return i
}
