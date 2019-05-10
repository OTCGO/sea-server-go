package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"sync/atomic"
)

type SyncBlock struct {
	threads int
	height int32
	status int32
}

func (sb *SyncBlock) Name() string {
	return BlockTask
}

func (sb *SyncBlock) Sync(block goutil.Map) error {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index")) + 1
	res, err := sb.Handle(block)
	if err != nil {
		log.Error("[SyncBlock] handle at height(%v) err: %v", height, err)
		return fmt.Errorf("handle fail(%v)", err)
	}

	b, _ := res.(*db.Block)
	_, err = db.InsertBlock(b)
	if err != nil {
		log.Error("[SyncBlock] insert block by height(%v) err: %v", height, err)
		return fmt.Errorf("insert block fail(%v)", err)
	}

	err = db.MustUpdateStatus(db.Status{Name: sb.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncBlock] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}

	atomic.StoreInt32(&sb.height, int32(height))
	return nil
}

func (sb *SyncBlock) Handle(block goutil.Map) (interface{}, error) {
	var sysFee int64
	for _, item := range block.GetMapArray("tx") {
		sysFee += item.GetInt64("sys_fee")
	}
	var totalSysFee = sysFee
	height := int(block.GetInt64("index")) + 1
	if height > 1 {
		prevBlock, err := db.GetBlock(height - 1)
		if err != nil {
			return nil, fmt.Errorf("get prev block fail(%v)", err)
		}
		totalSysFee += int64(prevBlock.TotalSysFee)
	}

	b := &db.Block{
		Height:      height,
		SysFee:      int(sysFee),
		TotalSysFee: int(totalSysFee),
		Raw:         block,
	}
	return b, nil
}

func (sb *SyncBlock) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(sb.Name())
	if err != nil {
		log.Error("[SyncBlock] get status err: %v", err)
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	var latestHeight int
	err = neo.Rpc(superNode.FastestNode.Value(), neo.MethodGetBlockCount, []interface{}{}, &latestHeight)
	if err != nil {
		log.Error("[SyncBlock] rpc get block count err: %v", err)
		return 0, 0, fmt.Errorf("rpc get block count fail(%v)", err)
	}
	return status.UpdateHeight, latestHeight, nil
}

func (sb *SyncBlock) Block(height int) (goutil.Map, error) {
	var block goutil.Map
	err := neo.Rpc(superNode.FastestNode.Value(), neo.MethodGetBlock, []int{height - 1, 1}, &block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (sb *SyncBlock) Threads() int {
	return sb.threads
}

func (sb *SyncBlock) SetStatus(status int32)  {
	atomic.StoreInt32(&sb.status, status)
}

func (sb *SyncBlock) Stats() goutil.Map {
	return goutil.Map{
		"height": atomic.LoadInt32(&sb.height),
		"status": atomic.LoadInt32(&sb.status),
	}
}

type Blocks []goutil.Map

func (blocks Blocks) Len() int {
	return len(blocks)
}

func (blocks Blocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}

func (blocks Blocks) Less(i, j int) bool {
	return blocks[i].GetInt64("index") < blocks[j].GetInt64("index")
}