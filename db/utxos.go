package db

import (
	"fmt"
	"github.com/hzxiao/goutil"
)

const (
	upsertUptSql    = `INSERT INTO upt(address,asset,update_height) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE update_height = ?`
	upserBalanceSql = `INSERT INTO balance(address,asset,value,last_updated_height) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE value = ?,last_updated_height = ?`
)

func InsertUtxo(t *Utxos) (bool, error) {
	if t == nil {
		return false, fmt.Errorf("utxo is nil")
	}

	return db.InsertOrIgnore(t, &Utxos{Txid: t.Txid, IndexN: t.IndexN})
}

func UpdateUtxoVinAndRet(t *Utxos) error {
	if t == nil {
		return fmt.Errorf("utxo is nil")
	}

	update := goutil.Map{
		"spent_txid":   t.Txid,
		"spent_height": t.SpentHeight,
		"status":       0,
	}
	_, err := db.engine.Table(TableUtxos).Where("txid = ? AND index_n = ?", t.Txid, t.IndexN).Update(update)
	if err != nil {
		return err
	}

	var u Utxos
	_, err = db.engine.Table(TableUtxos).Where("txid = ? AND index_n = ?", t.Txid, t.IndexN).Get(&u)
	*t = u

	return err
}

func UpdateUtxoClaim(t *Utxos) error {
	if t == nil {
		return fmt.Errorf("utxo is nil")
	}

	update := goutil.Map{
		"claim_txid":   t.Txid,
		"claim_height": t.SpentHeight,
	}
	_, err := db.engine.Table(TableUtxos).Where("txid = ? AND index_n = ?", t.Txid, t.IndexN).Update(update)
	return err
}

func InsertOrUpdateUpt(upt *Upt) (bool, error) {
	if upt == nil {
		return false, fmt.Errorf("upt is nil")
	}

	result, err := db.engine.Exec(upsertUptSql, upt.Address, upt.Asset, upt.UpdateHeight, upt.UpdateHeight)
	if err != nil {
		return false, err
	}
	effected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return effected > 0, nil
}

func ListUptByHeight(height int) ([]*Upt, error) {
	var upts []*Upt
	err := db.engine.Where("update_height = ?", height).Find(&upts)
	if err != nil {
		return nil, err
	}
	return upts, nil
}

func InsertOrUpdateBalance(b *Balance) (bool, error) {
	if b == nil {
		return false, fmt.Errorf("balance is nil")
	}

	result, err := db.engine.Exec(upserBalanceSql, b.Address, b.Asset, b.Value, b.LastUpdatedHeight,
		b.Value, b.LastUpdatedHeight)
	if err != nil {
		return false, err
	}
	effected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return effected > 0, nil
}

func InsertOrIgnoreHistory(h *History) (bool, error) {
	if h == nil {
		return false, fmt.Errorf("history is nil")
	}

	return db.InsertOrIgnore(h, &History{
		Txid:      h.Txid,
		Operation: h.Operation,
		IndexN:    h.IndexN,
	})
}
