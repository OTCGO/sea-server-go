package server

import (
	"github.com/OTCGO/sea-server-go/config"
	"github.com/OTCGO/sea-server-go/job/node"
	"github.com/OTCGO/sea-server-go/job/sync"
	"github.com/gin-gonic/gin"
	"github.com/hzxiao/goutil"
	"strconv"
)

//
// /:height/:task/mockSync
func mockSync(c *gin.Context) {
	height, err := strconv.Atoi(c.Param("height"))
	if err != nil {
		WriteJSON(c, nil, err)
		return
	}
	stats, err := sync.HandleOneHeight(height, c.Param("task"))
	WriteJSON(c, goutil.Map{"data": stats}, err)
}

//
// /adm/stats
func stats(c *gin.Context) {
	res := goutil.Map{
		"openSync":      config.Conf.OpenSync,
		"onlyOpenBlock": config.Conf.OnlySyncBlock,
		"sync":          sync.Stats(),
		"node":          node.SuperNode.Status(),
	}
	WriteJSON(c, res, nil)
}