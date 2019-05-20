package sync

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"strings"
	"sync/atomic"
)

type SyncUtxo struct {
	height int32
	status int32
}

func (su *SyncUtxo) Name() string {
	return UtxoTask
}

func (su *SyncUtxo) Sync(block goutil.Map) (err error) {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index"))
	res, _ := su.Handle(block)

	m, _ := res.(map[string]interface{})
	var addressInfo = map[string]struct{}{}
	vouts, _ := m["vouts"].([]*db.Utxos)
	for _, utxo := range vouts {
		_, err = db.InsertUtxo(utxo)
		if err != nil {
			log.Error("[SyncUtxo] insert utxo by txid(%v) height(%v) err: %v", utxo.Txid, height, err)
			return fmt.Errorf("insert utxo fail(%v)", err)
		}
		addressInfo[fmt.Sprintf("%v-%v", utxo.Address, utxo.Asset)] = struct{}{}
	}

	vins, _ := m["vins"].([]*db.Utxos)
	for _, utxo := range vins {
		err = db.UpdateUtxoVinAndRet(utxo)
		if err != nil {
			log.Error("[SyncUtxo] update utxo by txid(%v) height(%v) err: %v", utxo.Txid, height, err)
			return fmt.Errorf("update utxo fail(%v)", err)
		}
		addressInfo[fmt.Sprintf("%v-%v", utxo.Address, utxo.Asset)] = struct{}{}
	}

	claims, _ := m["claims"].([]*db.Utxos)
	for _, utxo := range claims {
		err = db.UpdateUtxoClaim(utxo)
		if err != nil {
			log.Error("[SyncUtxo] update utxo claim by txid(%v) height(%v) err: %v", utxo.Txid, height, err)
			return fmt.Errorf("update utxo claim fail(%v)", err)
		}
	}

	for k := range addressInfo {
		infos := strings.Split(k, "-")
		upt := &db.Upt{Address: infos[0], Asset: infos[1], UpdateHeight: height}
		_, err = db.InsertOrUpdateUpt(upt)
		if err != nil {
			log.Error("[SyncUtxo] update upt by height(%v) err: %v", height, err)
			return fmt.Errorf("update upt fail(%v)", err)
		}
	}

	err = db.MustUpdateStatus(db.Status{Name: su.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncUtxo] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}

	atomic.StoreInt32(&su.height, int32(height))
	return nil
}

func (su *SyncUtxo) Handle(block goutil.Map) (interface{}, error) {
	height := int(block.GetInt64("index"))

	var (
		vouts  = make([]*db.Utxos, 0)
		vins   = make([]*db.Utxos, 0)
		claims = make([]*db.Utxos, 0)
	)
	for _, tx := range block.GetMapArray("tx") {
		txid := tx.GetString("txid")
		//vout
		for _, vout := range tx.GetMapArray("vout") {
			uxto := &db.Utxos{
				Txid:        txid,
				IndexN:      int(vout.GetInt64("n")),
				Asset:       strings.TrimPrefix(vout.GetString("asset"), "0x"),
				Value:       vout.GetString("value"),
				Address:     vout.GetString("address"),
				Height:      height,
				Status:      1,
			}
			vouts = append(vouts, uxto)
		}
		//vin
		for _, vin := range tx.GetMapArray("vin") {
			uxto := &db.Utxos{
				Txid:        vin.GetString("txid"),
				SpentHeight: height,
				SpentTxid:   txid,
				IndexN:      int(vin.GetInt64("vout")),
			}
			vins = append(vins, uxto)
		}

		//claim
		for _, claim := range tx.GetMapArray("claim") {
			uxto := &db.Utxos{
				Txid:        claim.GetString("txid"),
				ClaimHeight: height,
				ClaimTxid:   txid,
				IndexN:      int(claim.GetInt64("vout")),
			}
			claims = append(claims, uxto)
		}
	}

	return map[string]interface{}{
		"vouts":  vouts,
		"vins":   vins,
		"claims": claims,
	}, nil
}

func (su *SyncUtxo) BlockHeight() (int, int, error) {
	height, err := getTaskHeightFromDB(su.Name(), BlockTask)
	if err != nil {
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	return height[su.Name()], height[BlockTask], nil
}

func (su *SyncUtxo) Block(height int) (goutil.Map, error) {
	b, err := db.GetBlock(height)
	if err != nil {
		return nil, err
	}
	return b.Raw, nil
}

func (su *SyncUtxo) Threads() int {
	return 1
}

func (su *SyncUtxo) SetStatus(status int32) {
	atomic.StoreInt32(&su.status, status)
}

func (su *SyncUtxo) Stats() goutil.Map {
	return goutil.Map{
		"height": atomic.LoadInt32(&su.height),
		"status": atomic.LoadInt32(&su.status),
	}
}
