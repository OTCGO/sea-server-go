package db

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestInsertStatus(t *testing.T) {
	deleteAll()

	var err error
	s1 := &Status{
		Name:         "utxo",
		UpdateHeight: 0,
	}
	err = InsertStatus(s1)
	assert.NoError(t, err)

	s2 := &Status{
		Name:         "assets",
		UpdateHeight: 0,
	}
	err = InsertStatus(s2)
	assert.NoError(t, err)

	//test err
	err = InsertStatus(s2)
	assert.Error(t, err)

	err = InsertStatus(nil)
	assert.Error(t, err)

	err = InsertStatus(&Status{})
	assert.Error(t, err)
}

func TestMustUpdateStatus(t *testing.T) {
	deleteAll()

	var err error
	s := &Status{
		Name:         "utxo",
		UpdateHeight: 0,
	}
	err = InsertStatus(s)
	assert.NoError(t, err)

	s1 := Status{Name: s.Name, UpdateHeight: 1}
	err = MustUpdateStatus(s1)
	assert.NoError(t, err)

	exist, err := db.engine.Exist(&Status{ID: s.ID, UpdateHeight: s1.UpdateHeight})
	assert.NoError(t, err)
	assert.True(t, exist)

	s1.UpdateHeight = 0
	err = MustUpdateStatus(s1)
	assert.NoError(t, err)

	exist, err = db.engine.Exist(&Status{ID: s.ID, UpdateHeight: s1.UpdateHeight})
	assert.NoError(t, err)
	assert.True(t, exist)

	s2 := Status{Name: "x", UpdateHeight: 2}
	err = MustUpdateStatus(s2)
	assert.NoError(t, err)
}

func TestFindAllStatus(t *testing.T) {
	deleteAll()

	var err error
	s1 := &Status{
		Name:         "utxo",
		UpdateHeight: 0,
	}
	err = InsertStatus(s1)
	assert.NoError(t, err)

	s2 := &Status{
		Name:         "assets",
		UpdateHeight: 0,
	}
	err = InsertStatus(s2)
	assert.NoError(t, err)

	ss, err := FindAllStatus()
	assert.NoError(t, err)
	assert.Len(t, ss, 2)
}

func TestGetStatus(t *testing.T) {
	deleteAll()

	var err error
	s1 := &Status{
		Name:         "utxo",
		UpdateHeight: 1,
	}
	err = InsertStatus(s1)
	assert.NoError(t, err)

	s, err := GetStatus(s1.Name)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, s1.UpdateHeight, s.UpdateHeight)

	_, err = GetStatus("x")
	assert.Error(t, err)
}

func TestInsertBlock(t *testing.T) {
	deleteAll()

	var err error
	b1 := &Block{
		Height: 1,
		Raw: goutil.Map{
			"index": 0,
		},
	}
	ok, err := InsertBlock(b1)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = InsertBlock(b1)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestGetBlock(t *testing.T) {
	deleteAll()

	var err error
	b1 := &Block{
		Height: 1,
		SysFee: 2,
		Raw: goutil.Map{
			"index": "0",
		},
	}
	ok, err := InsertBlock(b1)
	assert.NoError(t, err)
	assert.True(t, ok)

	b, err := GetBlock(b1.Height)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b1.SysFee, b.SysFee)
	assert.Equal(t, goutil.Map{"index": "0"}, b.Raw)

	_, err = GetBlock(0)
	assert.Error(t, err)
}

func TestCleanBlockRawData(t *testing.T) {
	deleteAll()

	var err error
	b1 := &Block{
		Height: 1,
		SysFee: 2,
		Raw: goutil.Map{
			"index": "0",
		},
	}
	ok, err := InsertBlock(b1)
	assert.NoError(t, err)
	assert.True(t, ok)

	err = CleanBlockRawData(0, 2)
	assert.NoError(t, err)

	b, err := GetBlock(b1.Height)
	assert.NoError(t, err)
	assert.Nil(t, b.Raw)
}

func TestInsertAssets(t *testing.T) {
	deleteAll()

	var err error
	a1 := &Assets{
		Asset:    "0xcc",
		Decimals: 8,
	}
	ok, err := InsertAssets(a1)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = InsertAssets(a1)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestGetStatusByNames(t *testing.T) {
	deleteAll()

	var err error
	s1 := &Status{
		Name:         "utxo",
		UpdateHeight: 0,
	}
	err = InsertStatus(s1)
	assert.NoError(t, err)

	s2 := &Status{
		Name:         "assets",
		UpdateHeight: 0,
	}
	err = InsertStatus(s2)
	assert.NoError(t, err)

	ss, err := GetStatusByNames("utxo", "assets")
	assert.NoError(t, err)
	assert.Len(t, ss, 2)

	ss, err = GetStatusByNames("utxo", "assets", "x")
	assert.NoError(t, err)
	assert.Len(t, ss, 2)

	ss, err = GetStatusByNames("x")
	assert.NoError(t, err)
	assert.Len(t, ss, 0)

}