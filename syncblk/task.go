package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/syncblk/supernode"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/container"
	"github.com/hzxiao/goutil/log"
	"sort"
	"sync"
	"time"
)

const dur = time.Second * 5

const (
	BlockTask   = "block"
	AssetsTask  = "assets"
	UtxoTask    = "utxo"
	BalanceTask = "balance"
	HistoryTask = "history"
)

var (
	tasks     = map[string]SyncTask{}
	superNode *supernode.NodeInfo
)

type SyncTask interface {
	Name() string
	Handle(block goutil.Map) (interface{}, error)
	Sync(block goutil.Map) error
	BlockHeight() (int, int, error)
	Block(height int) (goutil.Map, error)
	Threads() int
}

func Init() error {
	tasks := []SyncTask{&SyncBlock{threads: config.Conf.SyncBlockThreads},
		&SyncAssets{}, &SyncUtxo{}, &SyncBalance{}, &SyncHistory{}}
	if config.Conf.OnlySyncBlock {
		tasks = tasks[:1]
	}
	err := Register(tasks...)
	if err != nil {
		return err
	}
	superNode = supernode.SuperNode
	return nil
}

func Register(task ...SyncTask) error {
	for _, t := range task {
		_, found := tasks[t.Name()]
		if found {
			return fmt.Errorf("dup task name(%v)", t.Name())
		}
		tasks[t.Name()] = t
	}
	return nil
}

func SyncAll() {
	for name := range tasks {
		go runTask(tasks[name])
	}
}

func runTask(task SyncTask) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("[Sync] task(%v) panic by err: %v", e)
		}
	}()
	log.Info("[Sync] task(%v) start run", task.Name())
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	var err error
	for {
		var saveHeight, latestHeight int
		saveHeight, latestHeight, err = task.BlockHeight()
		if err != nil {
			log.Error("[Sync] task(%v) get block height err: %v", task.Name(), err)
		}

		if saveHeight == latestHeight {
			<-ticker.C
			continue
		}

		var wg = &sync.WaitGroup{}
		var blocks = container.NewSafeSlice(0)
		heights, count := splitHeight(saveHeight+1, latestHeight+1, task.Threads(), 10)
		for _, info := range heights {
			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for i := start; i < end; i++ {
					block, err := task.Block(i)
					if err != nil {
						log.Error("[Sync] task(%v) get block for height(%v) err: %v", task.Name(), i, err)
						return
					}
					blocks.Append(block)
				}
			}(info["start"], info["end"])
		}
		wg.Wait()

		if blocks.Len() != count {
			log.Error("[Sync] task(%v) get blocks less than count(%v) err: %v", task.Name(), count, err)
			continue
		}

		blockMaps := make([]goutil.Map, 0)
		blocks.Range(func(v interface{}) bool {
			blockMaps = append(blockMaps, goutil.MapV(v))
			return true
		})
		sort.Sort(Blocks(blockMaps))
		for _, b := range blockMaps {
			err = task.Sync(b)
			if err != nil {
				log.Error("[Sync] task(%v) do sync at height(%v) err: %v", task.Name(), b.GetInt64("index")+1, err)
				break
			}
			log.Info("[Sync] task(%v) do sync success at height(%v)", task.Name(), b.GetInt64("index")+1)
		}
	}
}

func splitHeight(start, end, threads, size int) (heights []map[string]int, count int) {
	if start >= end {
		return
	}
	num := (end - start) / size
	if num < threads && (end-start)%size > 0 {
		num += 1
	}
	curStart := start
	for i := 0; i < num && i < threads; i++ {
		curEnd := curStart + size
		if curEnd > end {
			curEnd = end
		}
		heights = append(heights, map[string]int{
			"start": curStart,
			"end":   curEnd,
		})
		count += curEnd - curStart
		curStart += size
	}
	return
}
