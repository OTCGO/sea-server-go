package clean

import (
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/job/sync"
	"math"
	"sync/atomic"
)

type BlockCleaner struct {
	cleaned int32
}

func (bc *BlockCleaner) Name() string {
	return "block"
}

func (bc *BlockCleaner) BlockHeight() (int, int, error) {
	height, err := getTaskHeightFromDB(sync.HistoryTask, sync.BalanceTask)
	if err != nil {
		return 0, 0, err
	}
	min := int(math.Max(float64(height[sync.BalanceTask]), float64(height[sync.HistoryTask])))
	return int(bc.cleaned), min, nil
}

func (bc *BlockCleaner) Clean(start, end int) error {
	err := db.CleanBlockRawData(start, end)
	if err != nil {
		return err
	}

	cleaned := -1
	if end > 0 {
		cleaned = end - 1
	}
	atomic.StoreInt32(&bc.cleaned, int32(cleaned))
	return nil
}
