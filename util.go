package main

import (
	"errors"
	"strconv"
	"strings"
)

func CounterKey(name string) string {
	return KeyGenerate("counter", name)
}

func KeyGenerate(names ...string) string {
	return strings.Join(names, ":")
}

func NextID(name string) (ID int, err error) {
	intCmd := redisDB.Incr(ctx, CounterKey(name))
	if err := intCmd.Err(); err != nil {
		return 0, err
	} else {
		ID = int(intCmd.Val())
	}
	return
}

func ExtraID(key string) (ID int, err error) {
	idx := strings.LastIndex(":", key)
	if idx == -1 || idx == len(key)-1 {
		err = errors.New("ID is not exist")
		return
	} else {
		ID, err = strconv.Atoi(key[idx+1:])
	}
	return
}
