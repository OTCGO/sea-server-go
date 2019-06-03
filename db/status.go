package db

import (
	"fmt"
	"github.com/hzxiao/goutil"
)

const upsertStatusSql = `INSERT INTO status(name,update_height) VALUES (?, ?) ON DUPLICATE KEY UPDATE update_height = ?;`

func InsertStatus(status *Status) error {
	if status == nil {
		return fmt.Errorf("status is nil")
	}
	if status.Name == "" {

		return fmt.Errorf("status.name is empty")
	}

	_, err := db.InsertMap(status, goutil.Struct2Map(status))
	if err != nil {
		return err
	}
	return nil
}

func MustUpdateStatus(status Status) error {
	_, err := db.engine.Exec(upsertStatusSql, status.Name, status.UpdateHeight, status.UpdateHeight)
	return err
}

func FindAllStatus() ([]*Status, error) {
	var s []*Status
	err := db.engine.Find(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GetStatus(name string) (*Status, error) {
	var status Status
	ok, err := db.engine.Where("name = ?", name).Get(&status)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &status, nil
}

func GetStatusByNames(names ...string) ([]*Status, error) {
	if len(names) == 0 {
		return nil, nil
	}

	var status []*Status
	err := db.engine.In("name", names).Find(&status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func InsertBlock(block *Block) (bool, error) {
	if block == nil {
		return false, fmt.Errorf("block is nil")
	}

	return db.InsertOrIgnore(block, &Block{Height: block.Height}, false)
}

func GetBlock(height int) (*Block, error) {
	var block Block
	ok, err := db.engine.Where("height = ?", height).Get(&block)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &block, nil
}

func CleanBlockRawData(start, end int) error {
	_, err := db.engine.Table(TableBlock).Exec("UPDATE block SET raw = NULL WHERE height >= ? AND height < ?;", start, end)
	return err
}

func InsertAssets(assets *Assets) (bool, error) {
	if assets == nil {
		return false, fmt.Errorf("assets is nil")
	}

	return db.InsertOrIgnore(assets, &Assets{Asset: assets.Asset}, false)
}

func GetAsset(asset string) (*Assets, error) {
	var a Assets
	ok, err := db.engine.Where("asset = ?", asset).Get(&a)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &a, nil
}

func GetAssetDecimals(asset string) (int, error) {
	var a Assets
	ok, err := db.engine.Where("asset = ?", asset).Get(&a)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return a.Decimals, nil
}