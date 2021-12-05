package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	
)

func CounterKey(name string) string {
	return KeyGenerate("counter", name)
}

func KeyGenerate(names ...string) string {
	return strings.Join(names, ":")
}

func NextID(name string) (ID int, err error) {

	redisDB.HSet(ctx, "counter", )


	intCmd := redisDB.Incr(ctx, CounterKey(name))
	if err := intCmd.Err(); err != nil {
		return 0, err
	} else {
		ID = int(intCmd.Val())
	}
	return
}

// ExtraID 提取ID
func ExtraID(key string) (ID int, err error) {
	idx := strings.LastIndex(key, ":")
	if idx == -1 || idx == len(key)-1 {
		err = errors.New("ID is not exist")
		return
	} else {
		ID, err = strconv.Atoi(key[idx+1:])
	}
	return
}

func BatchExtraID(keys []string) (ids []int, err error) {
	for _, key := range keys {
		if id, err := ExtraID(key); err != nil {
			break
		} else {
			ids = append(ids, id)
		}
	}
	return
}

func AnyEmptyStr(items ...string) (isEmpty bool) {
	for _, v := range items {
		if len(v) == 0 {
			return true
		}
	}
	return false
}
