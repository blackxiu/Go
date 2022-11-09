package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"mylearn/day07_FileAndjsonAndErr/example/example1/balance"
)

type HashBalance struct {
}

//注册一个接口实例
func init() {
	balance.RegisterBalancer("hash", &HashBalance{})
}

func (p *HashBalance) DoBalance(insts []*balance.Instance, key ...string) (inst *balance.Instance, err error) {
	var defKey = fmt.Sprintf("%d", rand.Int()) //默认key  随机的
	if len(key) > 0 {
		defKey = key[0]
	}

	lens := len(insts)
	if lens == 0 {
		err = fmt.Errorf("No backend instance")
		return
	}

	//哈希计算
	crcTable := crc32.MakeTable(crc32.IEEE)
	hashVal := crc32.Checksum([]byte(defKey), crcTable)
	index := int(hashVal) % lens
	inst = insts[index]
	return
}
