package service

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/job/sync"
	"github.com/hzxiao/goutil"
	"math"
)

func Height() (goutil.Map, error) {
	height := math.MaxInt32
	//first, query from sync task
	if config.Conf.OpenSync {
		for _, stat := range sync.Stats() {
			task := stat.GetString("task")
			if task == sync.UtxoTask || task == sync.HistoryTask {
				height = int(math.Min(float64(height), stat.GetFloat64("height")))
			}
		}
	}
	if height < math.MaxInt32 {
		return goutil.Map{"height": height}, nil
	}
	//query from db
	utxoStatus, err := db.GetStatus(sync.UtxoTask)
	if err != nil {
		return nil, fmt.Errorf("get utxo status error: %v", err)
	}
	historyStatus, err := db.GetStatus(sync.HistoryTask)
	if err != nil {
		return nil, fmt.Errorf("get history status error: %v", err)
	}

	return goutil.Map{
		"height": math.Min(float64(utxoStatus.UpdateHeight), float64(historyStatus.UpdateHeight)),
	}, nil
}

func Block(height int) (goutil.Map, error) {
	block, err := db.GetBlock(height)
	if err != nil {
		return nil, fmt.Errorf("get block(%v) error: %v", height, err)
	}

	return goutil.Struct2Map(block), nil
}
