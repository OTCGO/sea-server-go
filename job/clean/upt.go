package clean

import (
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/job/sync"
	"sync/atomic"
)

type UptCleaner struct {
	cleaned int32
}

func (uc *UptCleaner) Name() string {
	return "upt"
}

func (uc *UptCleaner) BlockHeight() (int, int, error) {
	height, err := getTaskHeightFromDB(sync.BalanceTask)
	if err != nil {
		return 0, 0, err
	}
	return int(uc.cleaned), height[sync.BalanceTask], nil
}

func (uc *UptCleaner) Clean(start, end int) error {
	err := db.DeleteUpt(start, end)
	if err != nil {
		return err
	}

	cleaned := -1
	if end > 0 {
		cleaned = end - 1
	}
	atomic.StoreInt32(&uc.cleaned, int32(cleaned))
	return nil
}
