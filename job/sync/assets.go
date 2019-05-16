package sync

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"strconv"
	"strings"
	"sync/atomic"
)

type SyncAssets struct {
	height int32
	status int32
}

func (sa *SyncAssets) Name() string {
	return AssetsTask
}

func (sa *SyncAssets) Sync(block goutil.Map) error {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index"))
	res, err := sa.Handle(block)
	if err != nil {
		log.Error("[SyncAssets] handle at height(%v) err: %v", height, err)
		return err
	}
	if res != nil {
		assets, ok := res.([]*db.Assets)
		if !ok {
			return fmt.Errorf("assert err")
		}
		for _, asset := range assets {
			_, err = db.InsertAssets(asset)
			if err != nil {
				log.Error("[SyncAssets] insert asset(%v) at height(%v) err: %v", asset.Asset, height, err)
				return fmt.Errorf("insert asset(%v) fail(%v)", asset.Asset, err)
			}
		}
	}
	err = db.MustUpdateStatus(db.Status{Name: sa.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncAssets] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update s  tatus fail(%v)", err)
	}

	atomic.StoreInt32(&sa.height, int32(height))
	return nil
}

func (sa *SyncAssets) Handle(block goutil.Map) (interface{}, error) {
	var assets []*db.Assets
	for _, tx := range block.GetMapArray("tx") {
		if tx.GetString("type") == "RegisterTransaction" {
			asset := parseGlobalAsset(tx)
			assets = append(assets, asset)
		}
		if neo.IsRegisterNep5AssetTx(tx) {
			asset, err := parseNep5Asset(tx)
			if err != nil {
				return nil, fmt.Errorf("parse nep5 asset by txid(%v) fail(%v)", tx.GetString("txid"), err)
			}
			if asset == nil {
				continue
			}
			assets = append(assets, asset)
		}
	}
	return assets, nil
}

func (sa *SyncAssets) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(sa.Name())
	if err != nil {
		log.Error("[SyncAssets] get status err: %v", err)
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	blockStatus, err := db.GetStatus(BlockTask)
	if err != nil {
		log.Error("[SyncAssets] get block status err: %v", err)
		return 0, 0, fmt.Errorf("get block status fail(%v)", err)
	}
	return status.UpdateHeight, blockStatus.UpdateHeight, nil
}

func (sa *SyncAssets) Block(height int) (goutil.Map, error) {
	b, err := db.GetBlock(height)
	if err != nil {
		return nil, err
	}
	return b.Raw, nil
}

func (sa *SyncAssets) Threads() int {
	return 1
}

func (sa *SyncAssets) SetStatus(status int32)  {
	atomic.StoreInt32(&sa.status, status)
}

func (sa *SyncAssets) Stats() goutil.Map {
	return goutil.Map{
		"height": atomic.LoadInt32(&sa.height),
		"status": atomic.LoadInt32(&sa.status),
	}
}

func parseGlobalAsset(tx goutil.Map) *db.Assets {
	return &db.Assets{
		Asset:    strings.TrimPrefix(tx.GetString("txid"), "0x"),
		Type:     tx.GetStringP("asset/type"),
		Name:     tx.GetStringP("asset/name/0/name"),
		Decimals: int(tx.GetInt64P("asset/precision")),
	}
}

func parseNep5Asset(tx goutil.Map) (*db.Assets, error) {
	contract, err := neo.ParseContract(tx.GetString("script"))
	if err != nil {
		return nil, fmt.Errorf("parse contract fail(%v)", err)
	}
	funcs := []string{"decimals", "totalSupply", "name", "symbol"}
	var result = map[string]string{}
	for _, f := range funcs {
		v, success, err := rpcInvoke([]string{contract.Contract, f})
		if err != nil {
			return nil, fmt.Errorf("rpc invoke func(%v) fail(%v)", f, err)
		}
		if !success {
			return nil, nil
		}
		result[f] = v.GetString("value")
	}
	var asset = &db.Assets{
		Type:         "NEP5",
		Asset:        contract.Contract,
		ContractName: contract.ContractName,
		Version:      contract.Version,
	}
	for f, v := range result {
		switch f {
		case "decimals":
			asset.Decimals, _ = strconv.Atoi(v)
		case "totalSupply":
		case "name":
			asset.Name = neo.HexDecode(v)
		case "symbol":
			asset.Symbol = neo.HexDecode(v)
		}
	}

	return asset, nil
}

func rpcInvoke(params interface{}) (goutil.Map, bool, error) {
	var invokeFail = func(r goutil.Map) bool {
		if r == nil {
			return false
		}
		return strings.HasPrefix(r.GetString("state"), "FAULT")
	}

	var result = goutil.Map{}
	err := neo.Rpc(superNode.FastestNode.Value(), neo.MethodInvokeFunction, params, &result)
	if err != nil {
		return nil, false, err
	}
	if invokeFail(result) {
		return nil, false, nil
	}

	return result.GetMapP("stack/0"), true, nil
}
