package clean

import (
	"github.com/OTCGO/sea-server-go/db"
	"github.com/hzxiao/goutil/log"
	"time"
)

type Cleaner interface {
	Name() string
	BlockHeight() (int, int, error)
	Clean(start, end int) error
}

var cleanTasks []Cleaner

func init() {
	cleanTasks = append(cleanTasks, &BlockCleaner{}, &UptCleaner{})
}

func StartClean() {
	for _, cleaner := range cleanTasks {
		go runCleanTask(cleaner)
	}
}

func runCleanTask(cleaner Cleaner) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		cleaned, latest, err := cleaner.BlockHeight()
		if err != nil {
			log.Error("[Clean] cleaner(%v) get block height err: %v", cleaner.Name(), err)
			continue
		}
		if cleaned+1 >= latest {
			continue
		}

		err = cleaner.Clean(cleaned+1, latest)
		if err != nil {
			log.Error("[Clean] cleaner(%v) do clean from %v to %v err: %v", cleaner.Name(), cleaned+1, latest, err)
		}
	}
}

func getTaskHeightFromDB(tasks ...string) (map[string]int, error) {
	ss, err := db.GetStatusByNames(tasks...)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int)
	for _, task := range tasks {
		res[task] = -1
	}

	for _, s := range ss {
		res[s.Name] = s.UpdateHeight
	}

	return res, nil
}
