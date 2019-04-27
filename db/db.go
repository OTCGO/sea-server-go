package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/hzxiao/goutil/log"
	"time"
)

type DB struct {
	uri    string
	engine *xorm.Engine
}

func (d *DB) init() error {
	engine, err := xorm.NewEngine("mysql", db.uri)
	if err != nil {
		return err
	}

	engine.SetMapper(core.GonicMapper{})
	db.engine = engine
	return nil
}

func (d *DB) loopPing() {
	var err error
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		<-ticker.C
		err = d.engine.Ping()
		if err == nil {
			continue
		}
		log.Error("[DB] ping database fail(%v) and try to reconnect...", err)
		err = d.reconnect()
		if err != nil {
			log.Error("[DB] reconnect to database fail: %v", err)
		} else {
			log.Info("[DB] reconnect to database success.")
		}
	}
}

func (d *DB) reconnect() error {
	return d.init()
}

func (d *DB) Insert(obj interface{}) (int64, error) {
	return d.engine.Insert(obj)
}

func (d *DB) InsertOrIgnore(obj, ignore interface{}) (bool, error) {
	exists, err := db.engine.Exist(ignore)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	effected, err := db.engine.Insert(obj)
	if err != nil {
		return false, err
	}
	return effected == 1, nil
}

var db *DB

func InitDB(uri string) (err error) {
	db = &DB{uri: uri}
	err = db.init()
	if err != nil {
		return
	}
	if err = db.engine.Ping(); err != nil {
		return
	}
	log.Info("[DB] connect to database success.")
	go db.loopPing()
	return
}

const (
	TableStatus  = "status"
	TableBlock   = "block"
	TableAssets  = "assets"
	TableUtxos   = "utxos"
	TableUpt     = "upt"
	TableBalance = "balance"
	TableHistory = "history"
)

func InitStatus(names ...string) error {
	for _, name := range names {
		exists, err := db.engine.Exist(&Status{Name: name})
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		err = InsertStatus(&Status{Name: name, UpdateHeight: 0})
		if err != nil {
			return fmt.Errorf("insert status by name(%v) err: %v", err)
		}
	}
	return nil
}
