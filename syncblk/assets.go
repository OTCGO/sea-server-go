package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"strings"
)

type SyncAssets struct {
}

func (sa *SyncAssets) Name() string {
	return AssetsTask
}

func (sa *SyncAssets) Sync(block goutil.Map) error {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index")) + 1

	var err error
	for _, tx := range block.GetMapArray("tx") {
		if tx.GetString("type") == "RegisterTransaction" {
			asset := parseGlobalAsset(tx)
			_, err = db.InsertAssets(asset)
			if err != nil {
				log.Error("[SyncAssets] insert global asset(%v) at height(%v) err: %v", asset.Asset, height, err)
				return fmt.Errorf("insert global asset(%v) fail(%v)", asset.Asset, err)
			}
		}
		if neo.IsRegisterNep5AssetTx(tx) {
			asset, err := parseNep5Asset(tx)
			if err != nil {
				log.Error("[SyncAssets] parse nep5 asset by txid(%v) at height(%v) err: %v", tx.GetString("txid"), height, err)
				return fmt.Errorf("parse nep5 asset by txid(%v) fail(%v)", tx.GetString("txid"), err)
			}
			_, err = db.InsertAssets(asset)
			if err != nil {
				log.Error("[SyncAssets] insert global asset(%v) at height(%v) err: %v", asset.Asset, height, err)
				return fmt.Errorf("insert global asset(%v) fail(%v)", asset.Asset, err)
			}
		}
	}
	err = db.MustUpdateStatus(db.Status{Name: sa.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncAssets] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}
	return nil
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

func parseGlobalAsset(tx goutil.Map) *db.Assets {
	return &db.Assets{
		Asset:    strings.TrimLeft(tx.GetString("txid"), "0x"),
		Type:     tx.GetStringP("asset/type"),
		Name:     tx.GetStringP("asset/name/0/name"),
		Decimals: int(tx.GetInt64P("asset/precision")),
	}
}

//TODO parse nep5 asset
func parseNep5Asset(tx goutil.Map) (*db.Assets, error) {
	return nil, nil
}