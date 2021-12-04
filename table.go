package main

import (
	"log"
	"strconv"
)

type Table interface {
	SetID(ID int)
	TableName() string
	ToMap() map[string]interface{}
}

func CreateItem(table Table) (err error) {
	ID, err := NextID(table.TableName())
	if err != nil {
		return
	}
	table.SetID(ID)
	key := KeyGenerate(table.TableName(), strconv.Itoa(ID))
	intCmd := redisDB.HSet(ctx, key, table.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

func UpdateItem(ID int, table Table) (err error) {
	key := KeyGenerate(table.TableName(), strconv.Itoa(ID))
	intCmd := redisDB.HSet(ctx, key, table.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

func RemoveItem(ID int, table Table) (err error) {
	key := KeyGenerate(table.TableName(), strconv.Itoa(ID))
	intCmd := redisDB.Del(ctx, key)
	if err = intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

