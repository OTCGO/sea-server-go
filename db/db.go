package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/hzxiao/goutil/log"
	"time"
)

type DB struct {
	uri    string
	engine *xorm.Engine
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
	engine, err := xorm.NewEngine("mysql", db.uri)
	if err != nil {
		return err
	}

	d.engine = engine
	return nil
}

var db *DB

func InitDB(uri string) (err error) {
	db = &DB{uri: uri}
	db.engine, err = xorm.NewEngine("mysql", db.uri)
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
