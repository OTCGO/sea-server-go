package sync

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/pkg/bigfloalt"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"math"
	"strings"
	"sync/atomic"
)

type SyncBalance struct {
	height int32
	status int32
}

func (sb *SyncBalance) Name() string {
	return BalanceTask
}

func (sb *SyncBalance) Sync(block goutil.Map) (err error) {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index")) + 1
	res, err := sb.Handle(block)
	if err != nil {
		log.Error("[SyncBalance] handle at height(%v) err: %v", height, err)
		return err
	}
	if res != nil {
		balances, ok := res.([]*db.Balance)
		if !ok {
			return fmt.Errorf("assert err")
		}
		for _, b := range balances {
			_, err = db.InsertOrUpdateBalance(b)
			if err != nil {
				log.Error("[SyncBalance] update balance by height(%v) err: %v", height, err)
				return fmt.Errorf("update balance fail(%v)", err)
			}
		}
	}
	err = db.MustUpdateStatus(db.Status{Name: sb.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncBalance] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}

	atomic.StoreInt32(&sb.height, int32(height))
	return nil
}

func (sb *SyncBalance) Handle(block goutil.Map) (interface{}, error) {
	height := int(block.GetInt64("index")) + 1
	var balances []*db.Balance
	for _, info := range block.GetMapArray("info") {
		address, asset := info.GetString("address"), info.GetString("asset")
		balance, err := rpcGetBalance(asset, address)
		if err != nil {
			return nil, fmt.Errorf("rpc get balance fail(%v)", err)
		}
		b := &db.Balance{
			Asset:             asset,
			Address:           address,
			Value:             balance,
			LastUpdatedHeight: height,
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (sb *SyncBalance) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(sb.Name())
	if err != nil {
		log.Error("[SyncBalance] get status err: %v", err)
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	assetStatus, err := db.GetStatus(AssetsTask)
	if err != nil {
		log.Error("[SyncBalance] get asset status err: %v", err)
		return 0, 0, fmt.Errorf("get asset status fail(%v)", err)
	}

	utxoStatus, err := db.GetStatus(UtxoTask)
	if err != nil {
		log.Error("[SyncBalance] get utxo status err: %v", err)
		return 0, 0, fmt.Errorf("get utxo status fail(%v)", err)
	}

	min := math.Min(float64(assetStatus.UpdateHeight), float64(utxoStatus.UpdateHeight))
	return status.UpdateHeight, int(min), nil
}

func (sb *SyncBalance) Block(height int) (goutil.Map, error) {
	upts, err := db.ListUptByHeight(height)
	if err != nil {
		return nil, err
	}
	block := goutil.Map{"index": height - 1}
	var info []goutil.Map
	for _, upt := range upts {
		info = append(info, goutil.Map{
			"address": upt.Address,
			"asset":   upt.Asset,
		})
	}
	block.Set("info", info)
	return block, nil
}

func (sb *SyncBalance) Threads() int {
	return 1
}

func (sb *SyncBalance) SetStatus(status int32)  {
	atomic.StoreInt32(&sb.status, status)
}

func (sb *SyncBalance) Stats() goutil.Map {
	return goutil.Map{
		"height": atomic.LoadInt32(&sb.height),
		"status": atomic.LoadInt32(&sb.status),
	}
}

func rpcGetBalance(asset, address string) (string, error) {
	if len(asset) == 40 {
		return rpcGetNep5Balance(asset, address)
	} else if len(asset) == 64 {
		return rpcGetGlobalBalance(asset, address)
	}
	return "", fmt.Errorf("invalid asset(%v)", asset)
}

func rpcGetNep5Balance(contract, address string) (balance string, err error) {
	hash, err := neo.Address2ScriptHash(address)
	if err != nil {
		return
	}
	params := []interface{}{
		contract,
		"balanceOf",
		[]goutil.Map{
			{
				"type":  "Hash160",
				"value": neo.ReverseBigLitterEndian(string(hash)),
			},
		},
	}
	v, success, err := rpcInvoke(params)
	if err != nil {
		return
	}
	if !success {
		err = fmt.Errorf("rpc invoke func fail")
		return
	}

	asset, err := db.GetAsset(contract)
	if err != nil {
		return "", fmt.Errorf("get asset(%v) from db err: %v", contract, err)
	}
	value := v.GetString("value")
	var base int
	switch v.GetString("type") {
	case "ByteArray":
		value = neo.ReverseBigLitterEndian(value)
		base = 16
	case "Integer":
		base = 10
	default:
		return "", fmt.Errorf("unkown value type(%v)", v.GetString("type"))
	}

	balance, err = bigfloalt.Format(value, base, asset.Decimals)
	if err != nil {
		return
	}
	return
}

func rpcGetGlobalBalance(asset, address string) (balance string, err error) {
	var result goutil.Map
	err = neo.Rpc(superNode.FastestNode.Value(), neo.MethodGetAccountState, []string{address}, &result)
	if err != nil {
		return
	}

	for _, b := range result.GetMapArray("balances") {
		if asset == strings.TrimPrefix(b.GetString("asset"), "0x") {
			balance = b.GetString("value")
			return
		}
	}
	return "0", nil
}
