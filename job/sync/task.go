package sync

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/job/node"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/container"
	"github.com/hzxiao/goutil/log"
	"github.com/hzxiao/goutil/slice"
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

// task status
const (
	taskStoped int32 = iota
	taskRunning
)

var (
	tasks     = map[string]SyncTask{}
	superNode *node.NodeInfo
)

type SyncTask interface {
	Name() string
	Handle(block goutil.Map) (interface{}, error)
	Sync(block goutil.Map) error
	BlockHeight() (int, int, error)
	Block(height int) (goutil.Map, error)
	Threads() int
	SetStatus(status int32)
	Stats() goutil.Map
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
	superNode = node.SuperNode
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
		task.SetStatus(taskStoped)
	}()
	task.SetStatus(taskRunning)
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
			log.Error("[Sync] task(%v) get blocks less than count(%v)", task.Name(), count)
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

//HandleOneHeight handle task by given height and return result
func HandleOneHeight(height int, name string) ([]goutil.Map, error) {
	if !slice.ContainsString([]string{"all", AssetsTask, BalanceTask, BlockTask, UtxoTask, HistoryTask}, name) {
		return nil, fmt.Errorf("unknown task(%v)", name)
	}
	tasks := []SyncTask{&SyncBlock{}, &SyncAssets{}, &SyncUtxo{}, &SyncBalance{}, &SyncHistory{}}
	if name != "all" {
		for _, task := range tasks {
			if task.Name() == name {
				tasks = []SyncTask{task}
				break
			}
		}
	}
	var stats []goutil.Map
	for _, task := range tasks {
		stat := goutil.Map{"task": task.Name()}
		block, err := task.Block(height)
		if err != nil {
			stat.Set("err", err)
			stats = append(stats, stat)
			continue
		}
		res, err := task.Handle(block)
		if err != nil {
			stat.Set("err", err)
		}
		stat.Set("data", res)
		stats = append(stats, stat)
	}
	return stats, nil
}
