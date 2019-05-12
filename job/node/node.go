package node

import (
	"expvar"
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/container"
	"github.com/hzxiao/goutil/log"
	"github.com/hzxiao/goutil/pool"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	SuperNode *NodeInfo
	grPool    *pool.Pool
)

func Init() {
	defaultNode := config.Conf.NeoUrl
	SuperNode = &NodeInfo{
		FastestNode:     expvar.NewString("fastestNode"),
		SupportLogNode:  expvar.NewString("supportLogNode"),
		fastestNodes:    container.NewSafeSlice(0),
		supportRpcNodes: container.NewSafeSlice(0),
		supportLogNodes: container.NewSafeSlice(0),
	}
	SuperNode.FastestNode.Set(defaultNode)
	SuperNode.SupportLogNode.Set(defaultNode)
	grPool = pool.NewPool(20, 20)

	go ScanNodes()
	go ScanHeight()
	go SuperNode.updateLoop()
}

type NodeInfo struct {
	Height         int32
	FastestNode    *expvar.String
	SupportLogNode *expvar.String

	fastestNodes    *container.SafeSlice
	supportRpcNodes *container.SafeSlice
	supportLogNodes *container.SafeSlice
}

func (ni *NodeInfo) updateLoop() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		nodes := getAllString(ni.fastestNodes)
		if len(nodes) > 0 {
			ni.FastestNode.Set(nodes[rand.Intn(len(nodes))])
		}

		nodes = getAllString(ni.supportLogNodes)
		if len(nodes) > 0 {
			ni.SupportLogNode.Set(nodes[rand.Intn(len(nodes))])
		}

		<-ticker.C
	}
}

func (ni *NodeInfo) Status() goutil.Map {
	return goutil.Map{
		"height":         atomic.LoadInt32(&ni.Height),
		"fastestNode":    ni.FastestNode.Value(),
		"supportLogNode": ni.SupportLogNode.Value(),
		"fast":           getAllString(ni.fastestNodes),
		"rpc":            getAllString(ni.supportRpcNodes),
		"log":            getAllString(ni.supportLogNodes),
	}
}

func getAllString(s *container.SafeSlice) []string {
	var values []string

	s.Range(func(value interface{}) bool {
		values = append(values, value.(string))
		return true
	})
	return values
}

func asyncRpcHeight(wg *sync.WaitGroup, url string, result *sync.Map) {
	grPool.JobQueue <- func() {
		defer wg.Done()
		var height int
		err := neo.RpcTimeout(url, neo.MethodGetBlockCount, []int{}, time.Second*10, &height)
		if err != nil {
			log.Error("[ScanNodes] rpc get block count by url(%v) err: %v", url, err)
			return
		}
		result.Store(url, height)
	}
}

func ScanNodes() {
	var asyncRpcGetAppLog = func(wg *sync.WaitGroup, url string, result *container.SafeSlice) {
		grPool.JobQueue <- func() {
			wg.Done()
			var r goutil.Map
			err := RpcGetAppLog(url, "", time.Second*30, &r)
			if err != nil {
				log.Error("[ScanNodes] rpc get app log by url(%v) err: %v", url, err)
				return
			}
			result.Append(url)
		}
	}
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		seeds := config.Conf.Seeds
		urls := goutil.RemoveDupString(append(seeds, getAllString(SuperNode.supportLogNodes)...))
		var wg sync.WaitGroup
		var result = sync.Map{}
		for _, url := range urls {
			wg.Add(1)
			asyncRpcHeight(&wg, url, &result)
		}
		wg.Wait()

		var maxHeight int
		var supportRpcNodes = container.NewSafeSlice(0)
		result.Range(func(key, value interface{}) bool {
			url, height := key.(string), value.(int)
			if height > maxHeight {
				maxHeight = height
			}
			supportRpcNodes.Append(url)
			return true
		})
		SuperNode.supportRpcNodes = supportRpcNodes

		var fastestNodes []interface{}
		result.Range(func(key, value interface{}) bool {
			url, height := key.(string), value.(int)
			if height == maxHeight {
				fastestNodes = append(fastestNodes, url)
			}
			return true
		})
		SuperNode.fastestNodes.Cover(fastestNodes)

		var supportLogNodes = container.NewSafeSlice(0)
		for _, url := range getAllString(SuperNode.supportRpcNodes) {
			wg.Add(1)
			asyncRpcGetAppLog(&wg, url, supportLogNodes)
		}
		wg.Wait()
		SuperNode.supportLogNodes = supportLogNodes

		<-ticker.C
	}
}

func ScanHeight() {
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()
	for {
		var wg sync.WaitGroup
		var result = sync.Map{}
		for _, url := range getAllString(SuperNode.fastestNodes) {
			wg.Add(1)
			asyncRpcHeight(&wg, url, &result)
		}
		wg.Wait()

		var maxHeight int
		result.Range(func(key, value interface{}) bool {
			height := value.(int)
			if height > maxHeight {
				maxHeight = height
			}
			return true
		})
		if atomic.LoadInt32(&SuperNode.Height) < int32(maxHeight) {
			atomic.StoreInt32(&SuperNode.Height, int32(maxHeight))
		}

		<-ticker.C
	}
}

func RpcGetAppLog(url, txid string, timeout time.Duration, result interface{}) error {
	if txid == "" {
		if config.Conf.Net == "testnet" {
			txid = "0x1bae5666ef5d645bb7d6edbe53a179763fda44a1b4ec6a49c2051883e03d0ba1"
		} else {
			txid = "0xd4e01144f6088028bc5af0e7e5e5dc9a0d133d54154275a966abd346d2319ff0"
		}
	}
	return neo.RpcTimeout(url, neo.MethodGetApplicationLog, []string{txid}, timeout, &result)
}
