package syncblk

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

	height := int(block.GetInt64("index")) + 1
	res, _ := su.Handle(block)

	m, _ := res.(map[string]interface{})
	utxos, _ := m["utxos"].([]*db.Utxos)
	for _, utxo := range utxos {
		err = db.UpdateUtxoClaim(utxo)
		if err != nil {
			log.Error("[SyncUtxo] update utxo by txid(%v) height(%v) err: %v", utxo.Txid, height, err)
			return fmt.Errorf("update utxo fail(%v)", err)
		}
	}

	addressInfo, _ := m["addressInfo"].(map[string]struct{})
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
	height := int(block.GetInt64("index")) + 1

	var utxos = make([]*db.Utxos, 0)
	var addressInfo = map[string]struct{}{}
	for _, tx := range block.GetMapArray("tx") {
		txid := tx.GetString("txid")
		//vout
		for _, vout := range tx.GetMapArray("vout") {
			uxto := &db.Utxos{
				Txid:    txid,
				IndexN:  int(vout.GetInt64("n")),
				Asset:   strings.TrimPrefix(vout.GetString("asset"), "0x"),
				Value:   vout.GetString("value"),
				Address: vout.GetString("address"),
				Height:  height,
				Status:  1,
			}
			utxos = append(utxos, uxto)
			addressInfo[fmt.Sprintf("%v-%v", uxto.Address, uxto.Asset)] = struct{}{}
		}
		//vin
		for _, vin := range tx.GetMapArray("vin") {
			uxto := &db.Utxos{
				Txid:        vin.GetString("txid"),
				SpentHeight: height,
				SpentTxid:   txid,
				IndexN:      int(vin.GetInt64("vout")),
			}
			utxos = append(utxos, uxto)
			addressInfo[fmt.Sprintf("%v-%v", uxto.Address, uxto.Asset)] = struct{}{}
		}

		//claim
		for _, claim := range tx.GetMapArray("claim") {
			uxto := &db.Utxos{
				Txid:        claim.GetString("txid"),
				ClaimHeight: height,
				ClaimTxid:   txid,
				IndexN:      int(claim.GetInt64("vout")),
			}
			utxos = append(utxos, uxto)
		}
	}

	return map[string]interface{}{
		"utxos":       utxos,
		"addressInfo": addressInfo,
	}, nil
}

func (su *SyncUtxo) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(su.Name())
	if err != nil {
		log.Error("[SyncUtxo] get status err: %v", err)
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	blockStatus, err := db.GetStatus(BlockTask)
	if err != nil {
		log.Error("[SyncUtxo] get block status err: %v", err)
		return 0, 0, fmt.Errorf("get block status fail(%v)", err)
	}
	return status.UpdateHeight, blockStatus.UpdateHeight, nil
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

func (su *SyncUtxo) SetStatus(status int32)  {
	atomic.StoreInt32(&su.status, status)
}

func (su *SyncUtxo) Stats() goutil.Map {
	return goutil.Map{
		"height": atomic.LoadInt32(&su.height),
		"status": atomic.LoadInt32(&su.status),
	}
}
