package syncblk

import (
	"fmt"
	"github.com/OTCGO/sea-server-go/db"
	"github.com/OTCGO/sea-server-go/pkg/bigfloalt"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/log"
	"sort"
	"strings"
)

type SyncHistory struct {
}

func (sa *SyncHistory) Name() string {
	return HistoryTask
}

func (sa *SyncHistory) Sync(block goutil.Map) error {
	if block == nil {
		return fmt.Errorf("block is nil")
	}

	height := int(block.GetInt64("index")) + 1

	var err error
	var (
		gTxids, sTxids []string
	)
	for _, tx := range block.GetMapArray("tx") {
		for _, vin := range tx.GetMapArray("vin") {
			gTxids = append(gTxids, vin.GetString("txid"))
		}
		if tx.GetString("type") == "InvocationTransaction" {
			sTxids = append(sTxids, tx.GetString("txid"))
		}
	}
	utxoMap, err := rpcUtxoByTxids(gTxids)
	if err != nil {
		log.Error("[SyncHistory] rpc get tx async err: %v", err)
		return fmt.Errorf("rpc get tx async fail(%v)", err)
	}
	logMap, err := rpcLogByTxids(sTxids)
	if err != nil {
		log.Error("[SyncHistory] rpc get app log async err: %v", err)
		return fmt.Errorf("rpc get app log async fail(%v)", err)
	}

	var historyList []*db.History
	blockTime := block.GetInt64("time")
	for _, tx := range block.GetMapArray("tx") {
		txid := tx.GetString("txid")
		//nep5
		if tx.GetString("type") == "InvocationTransaction" {
			appLog, ok := formatAppLog(logMap.GetMap(tx.GetString("txid")))
			if !ok {
				continue
			}
			hs, err := parseNep5History(appLog, txid, int(blockTime))
			if err != nil {
				log.Error("[SyncHistory] parse nep5 history by height(%v) err: %v", height, err)
				return fmt.Errorf("parse nep5 history fail(%v)", err)
			}
			historyList = append(historyList, hs...)
		}

		//global
		utxoM, voutM := goutil.Map{}, goutil.Map{}
		for _, vin := range tx.GetMapArray("vin") {
			utxo := utxoMap.GetMapP(fmt.Sprintf("%v/%v", vin.GetString("txid"), vin.GetInt64("vout")))
			if utxo == nil {
				log.Error("[SyncHistory] lack of utxo(%v, %v)", vin.GetString("txid"), vin.GetInt64("vout"))
				return fmt.Errorf("lack of utxo")
			}
			key := utxo.GetString("asset")+"_"+utxo.GetString("address")
			value := utxo.GetString("value")
			_, ok := utxoM[key]
			if ok {
				value, _ = bigfloalt.Add(value, utxoM.GetStringP(key+"/value"))
			}
			utxo.Set("value", value)
			utxoM.Set(key, utxo)
		}

		for _, vout := range tx.GetMapArray("vout") {
			key := vout.GetString("asset")+"_"+vout.GetString("address")
			value := vout.GetString("value")
			_, ok := voutM[key]
			if ok {
				value, _ = bigfloalt.Add(value, voutM.GetStringP(key+"/value"))
			}
			vout.Set("value", value)
			voutM.Set(key, vout)
		}

		if isTransferItself(utxoM, voutM) {
			continue
		}

		keys := utxoM.Keys()
		sort.Strings(keys)
		for i, key := range keys {
			utxo := utxoM.GetMap(key)
			if voutM.Exist(key) {
				gt, _ := bigfloalt.Gt(utxo.GetString("value"), voutM.GetStringP(key+"/value"))
				if gt {
					diff, _ := bigfloalt.Sub(utxo.GetString("value"), voutM.GetStringP(key+"/value"))
					utxo.Set("value", diff)
					delete(voutM, key)
				}
			}
			historyList = append(historyList, &db.History{
				Txid:      txid,
				Operation: "out",
				Asset:     strings.TrimPrefix(utxo.GetString("asset"), "0x"),
				IndexN:    i,
				Address:   utxo.GetString("address"),
				Value:     utxo.GetString("value"),
				Timepoint: int(blockTime),
			})
		}

		keys = voutM.Keys()
		sort.Strings(keys)
		for i, key := range keys {
			vout := voutM.GetMap(key)
			historyList = append(historyList, &db.History{
				Txid:      txid,
				Operation: "in",
				Asset:     strings.TrimPrefix(vout.GetString("asset"), "0x"),
				IndexN:    i,
				Address:   vout.GetString("address"),
				Value:     vout.GetString("value"),
				Timepoint: int(blockTime),
			})
		}
	}

	for i := range historyList {
		_, err = db.InsertOrIgnoreHistory(historyList[i])
		if err != nil {
			log.Error("[SyncHistory] update history by height(%v) err: %v", height, err)
			return fmt.Errorf("update history fail(%v)", err)
		}
	}
	//update status
	err = db.MustUpdateStatus(db.Status{Name: sa.Name(), UpdateHeight: height})
	if err != nil {
		log.Error("[SyncHistory] update status by height(%v) err: %v", height, err)
		return fmt.Errorf("update status fail(%v)", err)
	}
	return nil
}

func (sa *SyncHistory) BlockHeight() (int, int, error) {
	status, err := db.GetStatus(sa.Name())
	if err != nil {
		log.Error("[SyncHistory] get status err: %v", err)
		return 0, 0, fmt.Errorf("get status fail(%v)", err)
	}

	assetStatus, err := db.GetStatus(AssetsTask)
	if err != nil {
		log.Error("[SyncHistory] get asset status err: %v", err)
		return 0, 0, fmt.Errorf("get asset status fail(%v)", err)
	}
	return status.UpdateHeight, assetStatus.UpdateHeight, nil
}

func (sa *SyncHistory) Block(height int) (goutil.Map, error) {
	b, err := db.GetBlock(height)
	if err != nil {
		return nil, err
	}
	return b.Raw, nil
}

func (sa *SyncHistory) Threads() int {
	return 1
}

func rpcUtxoByTxids(txids []string) (goutil.Map, error) {
	txids = goutil.RemoveDupString(txids)
	r := goutil.Map{}
	for _, txid := range txids {
		var result goutil.Map
		rpcErr := neo.Rpc(superNode.FastestNode.Value(), neo.MethodGetRawTransaction, []interface{}{txid, 1}, &result)
		if rpcErr != nil {
			return nil, rpcErr
		}
		r.Set(txid, result.Get("vout"))
	}
	return r, nil
}

func rpcLogByTxids(txids []string) (goutil.Map, error) {
	txids = goutil.RemoveDupString(txids)
	r := goutil.Map{}
	for _, txid := range txids {
		var result goutil.Map
		rpcErr := neo.Rpc(superNode.FastestNode.Value(), neo.MethodGetApplicationLog, []interface{}{txid}, &result)
		if rpcErr != nil {
			return nil, rpcErr
		}
		r.Set(txid, result)
	}
	return r, nil
}

func formatAppLog(appLog goutil.Map) (goutil.Map, bool) {
	if appLog == nil {
		return nil, false
	}
	var vmstate = appLog.GetString("vmstate")
	if vmstate == "HALT, BREAK" {
		return appLog, true
	} else if vmstate != "" {
		return nil, false
	}

	return formatAppLog(appLog.GetMapP("executions/0"))
}

func parseNep5History(appLog goutil.Map, txid string, blockTime int) ([]*db.History, error) {
	var historyList []*db.History
	for i, noti := range appLog.GetMapArray("notifications") {
		if !isNep5TransferNotification(noti) {
			continue
		}
		asset := strings.TrimPrefix(noti.GetString("contract"), "0x")
		decimals, err := db.GetAssetDecimals(asset)
		if err != nil {
			return nil, fmt.Errorf("get decimals by asset(%v) fail(%v)", asset, err)
		}
		var value string
		if noti.GetStringP("state/value/3/type") == "Integer" {
			value, err = bigfloalt.Format(noti.GetStringP("state/value/3/value"), 10, decimals)
		} else {
			value, err = bigfloalt.Format(noti.GetStringP("state/value/3/value"), 16, decimals)
		}
		if err != nil {
			return nil, err
		}
		fromSh := noti.GetStringP("state/value/1/value")
		if fromSh != "" {
			from, err := neo.ScriptHash2Address(neo.HexDecodeString(fromSh))
			if err == nil {
				historyList = append(historyList, &db.History{
					Txid:      txid,
					Operation: "out",
					Asset:     asset,
					IndexN:    i,
					Address:   from,
					Value:     value,
					Timepoint: blockTime,
				})
			}
		}
		to, _ := neo.ScriptHash2Address(neo.HexDecodeString(noti.GetStringP("state/value/2/value")))
		historyList = append(historyList, &db.History{
			Txid:      txid,
			Operation: "in",
			Asset:     asset,
			IndexN:    i,
			Address:   to,
			Value:     value,
			Timepoint: blockTime,
		})
	}
	return historyList, nil
}

func isNep5TransferNotification(noti goutil.Map) bool {
	if noti == nil {
		return false
	}
	if len(noti.GetMapArrayP("state/value")) == 4 &&
		noti.GetStringP("state/value/0/value") == neo.HexEncode("transfer") {
		return true
	}
	return false
}

func isTransferItself(utxoM, voutM goutil.Map) bool {
	if len(utxoM) == len(voutM) && len(voutM) == 1 {
		k1 := utxoM.Keys()[0]
		k2 := voutM.Keys()[0]
		if k1 == k2 && utxoM.GetStringP(k1+"/value") == voutM.GetStringP(k2+"/value") {
			return true
		}
	}
	return false
}