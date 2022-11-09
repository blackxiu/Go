package balancing

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"

	"github.com/shenzhendev/hmyrpc/rpc"
)

type ConsistentHashMap struct {
	hash    func(data []byte) uint32
	keys    uint32Slice
	replica int
	hashMap map[uint32]rpc.RPCClient
	len     int
	lock    sync.RWMutex
}

func NewConsistentHashBalance(replica int, sign func(data []byte) uint32) *ConsistentHashMap {
	chmap := &ConsistentHashMap{
		hash:    sign,
		replica: replica,
		hashMap: make(map[uint32]rpc.RPCClient),
	}
	if sign == nil {
		chmap.hash = crc32.ChecksumIEEE // in redis ...
	}
	return chmap
}

func (m *ConsistentHashMap) Add(endpoint rpc.RPCClient, params ...string) error {
	if endpoint.Url() == "" {
		return errors.New("endpoint url is empty")
	}
	// lock
	m.lock.Lock()
	defer m.lock.Unlock()
	// for
	for i := 0; i < m.replica; i++ {
		hash := m.hash([]byte(strconv.Itoa(i) + endpoint.Url()))
		m.keys = append(m.keys, hash)
		m.hashMap[hash] = endpoint
	}
	sort.Sort(m.keys)
	return nil
}

func (m *ConsistentHashMap) Mark(endpoint rpc.RPCClient, up bool) error {
	panic("not implemented")
}

func (m *ConsistentHashMap) Get(key string) rpc.RPCClient {
	// if empty
	if m.IsEmpty() {
		return nil
	}
	//
	h := m.hash([]byte(key))
	index := sort.Search(len(m.keys), func(i int) bool {
		//fmt.Println(m.keys[i], h, m.hashMap[m.keys[i]].Url, i, m.keys[i] > h == true)
		return m.keys[i] > h
	})
	if index == len(m.keys) {
		index = 0 // default
	}
	// lock
	m.lock.RLock()
	defer m.lock.RUnlock()
	//
	fmt.Println(m.hashMap[m.keys[index]])
	return m.hashMap[m.keys[index]]
}

func (m *ConsistentHashMap) IsEmpty() bool {
	return len(m.keys) == 0
}

type uint32Slice []uint32

func (s uint32Slice) Len() int {
	return len(s)
}

func (s uint32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s uint32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
