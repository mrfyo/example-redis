package main

import "strings"

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
