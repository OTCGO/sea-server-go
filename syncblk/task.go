package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/syncblk/supernode"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
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
	tasks = map[string]SyncTask{}
	superNode  *supernode.NodeInfo
)

type SyncTask interface {
	Name() string
	Sync(block goutil.Map) error
	BlockHeight() (int, int, error)
	Block(height int) (goutil.Map, error)
}

func Init() error {
	err := Register(&SyncBlock{}, &SyncAssets{}, &SyncUtxo{}, &SyncBalance{}, &SyncHistory{})
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

		var h = saveHeight + 1
		for h <= latestHeight {
			block, err := task.Block(h)
			if err != nil {
				log.Error("[Sync] task(%v) get block for height(%v) err: %v", task.Name(), h, err)
				continue
			}
			err = task.Sync(block)
			if err != nil {
				log.Error("[Sync] task(%v) do sync err: %v", task.Name(), err)
				continue
			}
			log.Info("[Sync] task(%v) do sync success at height(%v)", h)
			h++
		}
	}
}
