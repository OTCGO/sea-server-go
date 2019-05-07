package db

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestInsertOrUpdateUpt(t *testing.T) {
	deleteAll()

	upt := &Upt{Address: "1", Asset: "2", UpdateHeight: 1}
	ok, err := InsertOrUpdateUpt(upt)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = InsertOrUpdateUpt(upt)
	assert.NoError(t, err)
	assert.False(t, ok)

	upt.UpdateHeight = 2
	ok, err = InsertOrUpdateUpt(upt)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestInsertOrUpdateBalance(t *testing.T) {
	deleteAll()

	b := &Balance{
		Address:           "1",
		Asset:             "2",
		Value:             "3",
		LastUpdatedHeight: 2,
	}
	ok, err := InsertOrUpdateBalance(b)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestInsertOrIgnoreHistory(t *testing.T) {
	deleteAll()

	h := &History{
		Txid:      "123",
		Operation: "in",
		IndexN:    0,
		Address:   "1234",
	}
	ok, err := InsertOrIgnoreHistory(h)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = InsertOrIgnoreHistory(h)
	assert.NoError(t, err)
	assert.False(t, ok)

	h.Operation = "out"
	ok, err = InsertOrIgnoreHistory(h)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestUpdateUtxoVinAndRet(t *testing.T) {
	deleteAll()

	utxo := &Utxos{
		Txid:    "1",
		IndexN:  0,
		Asset:   "1",
		Address: "1",
	}
	_, err := InsertUtxo(utxo)
	assert.NoError(t, err)

	vin := &Utxos{
		SpentTxid:   "2",
		SpentHeight: 2,
		Txid:        utxo.Txid,
		IndexN:      0,
	}
	err = UpdateUtxoVinAndRet(vin)
	assert.NoError(t, err)

	assert.Equal(t, utxo.Address, vin.Address)
}
