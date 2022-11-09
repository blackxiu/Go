package balancing

import (
	"errors"
	"strconv"
	"sync"

	"github.com/shenzhendev/hmyrpc/rpc"
)

type WeightMap struct {
	len   int
	lock  sync.RWMutex
	Nodes []*nodeWithWeight
}

type nodeWithWeight struct {
	rpc             rpc.RPCClient
	InitWeight      int64
	currentWeight   int64
	effectiveWeight int64
}

func NewWeightBalance() *WeightMap {
	return &WeightMap{
		len:   0,
		Nodes: nil,
	}
}

func (m *WeightMap) Add(endpoint rpc.RPCClient, params ...string) error {
	if endpoint.Url() == "" {
		return errors.New("endpoint url is empty")
	}
	if len(params) == 0 {
		return errors.New("params is empty")
	}
	// weight
	weightInt, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return err
	}
	// lock
	m.lock.Lock()
	defer m.lock.Unlock()

	// add new rpc client
	m.Nodes = append(m.Nodes, &nodeWithWeight{
		rpc:             endpoint,
		InitWeight:      weightInt,
		currentWeight:   0,
		effectiveWeight: weightInt,
	})
	// add length
	m.len += 1

	return nil
}

func (m *WeightMap) Get(key string) (choose rpc.RPCClient) {
	return m.next()
}

func (m *WeightMap) next() (choose rpc.RPCClient) {
	// init total weight
	totalWeight := int64(0)
	chooseIndex := 0
	// lock
	m.lock.Lock()
	defer m.lock.Unlock()
	// for
	for i := 0; i < m.len; i++ {
		totalWeight += m.Nodes[i].effectiveWeight
		m.Nodes[i].currentWeight += m.Nodes[i].effectiveWeight
		if choose == nil || m.Nodes[i].currentWeight > m.currentWeight(choose) {
			choose = m.Nodes[i].rpc
			chooseIndex = i
		}
	}
	// if choose is nil, return default
	if choose == nil {
		choose = m.Nodes[0].rpc
		chooseIndex = 0
	}
	//
	m.Nodes[chooseIndex].currentWeight -= totalWeight
	return
}

func (m *WeightMap) currentWeight(endpoint rpc.RPCClient) int64 {
	for _, ep := range m.Nodes {
		if ep.rpc == endpoint {
			return ep.currentWeight
		}
	}
	// if not match, return default
	return m.Nodes[0].currentWeight
}

func (m *WeightMap) Mark(endpoint rpc.RPCClient, up bool) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, ep := range m.Nodes {
		if ep.rpc == endpoint {
			if up && ep.effectiveWeight < ep.InitWeight {
				ep.effectiveWeight++
			}
			if !up && ep.effectiveWeight >= 1 {
				ep.effectiveWeight--
			}
		}
	}
	return nil
}
