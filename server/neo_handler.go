package server

import (
	"fmt"
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

// height get neo height
// /height
func neoHeight(c *gin.Context) {
	res, err := neoService.Height()
	WriteJSON(c, res, err)
}

// block get neo block by height
// /block/:height
func neoBlock(c *gin.Context) {
	h, err := strconv.Atoi("height")
	if err != nil {
		WriteJSON(c, nil, err)
		return
	}
	if h <= 0 {
		WriteJSON(c, nil, fmt.Errorf("invalid height"))
		return
	}

	res, err := neoService.Block(h)
	WriteJSON(c, res, err)
}