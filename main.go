package main

import (
	"flag"
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/job/node"
	"github.com/OTCGO/sea-server-go/job/sync"
	"github.com/OTCGO/sea-server-go/server"
	"github.com/hzxiao/goutil/log"
	"github.com/hzxiao/goutil/version"
	"os"
	"os/signal"
	"syscall"
)

var (
	conf string
	ver  bool
)

func init() {
	flag.BoolVar(&ver, "v", false, "show version info.")
	flag.StringVar(&conf, "conf", "seago.toml", "default config file")
}

func main() {
	flag.Parse()

	var err error
	if ver {
		err = version.Print()
		if err != nil {
			panic(err)
		}
		return
	}
	err = config.Init(conf)
	if err != nil {
		panic(err)
	}

	//log init
	err = log.SetLogger(config.Conf.Log.OutputWay, config.Conf.Log.Output)
	if err != nil {
		panic(err)
	}

	err = initTaskAndServer()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func initTaskAndServer() (err error) {
	err = db.InitDB(config.Conf.MySQLUrl)
	if err != nil {
		return
	}

	if config.Conf.OpenSync {
		node.Init()
		err = sync.Init()
		if err != nil {
			return
		}
		sync.SyncAll()
	}

	server.Init()
	go server.Run()

	return
}
