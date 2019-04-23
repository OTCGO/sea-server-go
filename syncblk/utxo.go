package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"strings"
)

type SyncUtxo struct {
}

func (sa *SyncUtxo) Name() string {
	return UtxoTask
}

func (sa *SyncUtxo) Sync(block goutil.Map) (err error) {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index")) + 1

	for _, tx := range block.GetMapArray("tx") {
		txid := tx.GetString("txid")
		//vin
		for _, vin := range tx.GetMapArray("vin") {
			uxto := &db.Utxos{
				Txid:        vin.GetString("txid"),
				SpentHeight: height,
				SpentTxid:   txid,
				IndexN:      int(vin.GetInt64("vout")),
			}
			err = db.UpdateUtxoVinAndRet(uxto)
			if err != nil {
				log.Error("[SyncUtxo] update vin by txid(%v) height(%v) err: %v", txid, height, err)
				return fmt.Errorf("update vin fail(%v)", err)
			}
		}
		//vout
		for _, vout := range tx.GetMapArray("vout") {
			uxto := &db.Utxos{
				Txid:    txid,
				IndexN:  int(vout.GetInt64("n")),
				Asset:   strings.TrimLeft(vout.GetString("asset"), "0x"),
				Value:   vout.GetString("value"),
				Address: vout.GetString("address"),
				Height:  height,
			}
			_, err = db.InsertUtxo(uxto)
			if err != nil {
				log.Error("[SyncUtxo] insert vout by txid(%v) height(%v) err: %v", txid, height, err)
				return fmt.Errorf("insert vout fail(%v)", err)
			}
		}
		//claim
		for _, claim := range tx.GetMapArray("claim") {
			uxto := &db.Utxos{
				Txid:        claim.GetString("txid"),
				ClaimHeight: height,
				ClaimTxid:   txid,
				IndexN:      int(claim.GetInt64("vout")),
			}
			err = db.UpdateUtxoClaim(uxto)
			if err != nil {
				log.Error("[SyncUtxo] update claim by txid(%v) height(%v) err: %v", txid, height, err)
				return fmt.Errorf("update claim fail(%v)", err)
			}
		}
	}
	err = db.MustUpdateStatus(db.Status{Name: sa.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncUtxo] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}
	return nil
}

func (sa *SyncUtxo) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(sa.Name())
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

func (sa *SyncUtxo) Block(height int) (goutil.Map, error) {
	b, err := db.GetBlock(height)
	if err != nil {
		return nil, err
	}
	return b.Raw, nil
}
