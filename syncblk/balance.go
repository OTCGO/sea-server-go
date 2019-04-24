package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"strings"
)

type SyncBalance struct {
}

func (sa *SyncBalance) Name() string {
	return BalanceTask
}

func (sa *SyncBalance) Sync(block goutil.Map) (err error) {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index")) + 1
	for _, info := range block.GetMapArray("info") {
		address, asset := info.GetString("address"), info.GetString("asset")
		balance, err := rpcGetBalance(asset, address)
		if err != nil {
			log.Error("[SyncBalance] rpc get balance by asset(%v), address(%v) err: %v", asset, address, err)
			return fmt.Errorf("rpc get balance fail(%v)", err)
		}
		b := &db.Balance{
			Asset:             asset,
			Address:           address,
			Value:             balance,
			LastUpdatedHeight: height,
		}
		_, err = db.InsertOrUpdateBalance(b)
		if err != nil {
			log.Error("[SyncBalance] update balance by height(%v) err: %v", height, err)
			return fmt.Errorf("update balance fail(%v)", err)
		}
	}
	err = db.MustUpdateStatus(db.Status{Name: sa.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncBalance] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}
	return nil
}

func (sa *SyncBalance) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(sa.Name())
	if err != nil {
		log.Error("[SyncBalance] get status err: %v", err)
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	blockStatus, err := db.GetStatus(UtxoTask)
	if err != nil {
		log.Error("[SyncBalance] get utxo status err: %v", err)
		return 0, 0, fmt.Errorf("get utxo status fail(%v)", err)
	}
	return status.UpdateHeight, blockStatus.UpdateHeight, nil
}

func (sa *SyncBalance) Block(height int) (goutil.Map, error) {
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

	balance, err = neo.FormatBigFloat(value, base, asset.Decimals)
	if err != nil {
		return
	}
	return
}

func rpcGetGlobalBalance(asset, address string) (balance string, err error) {
	var result goutil.Map
	err = neo.Rpc(neo.MethodGetAccountState, []string{address}, &result)
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
