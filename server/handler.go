package server

import (
	"github.com/OTCGO/sea-server-go/syncblk"
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
	stats, err := syncblk.HandleOneHeight(height, c.Param("task"))
	WriteJSON(c, goutil.Map{"data": stats}, err)
}
