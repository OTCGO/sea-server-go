package db

import "github.com/hzxiao/goutil"

type Status struct {
	ID           int
	Name         string
	UpdateHeight int
}

type Block struct {
	Height      int
	SysFee      int
	TotalSysFee int
	Raw         goutil.Map
}

type Assets struct {
	ID           int    `json:"id,omitempty"`
	Asset        string `json:"asset,omitempty"`
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	Version      string `json:"version,omitempty"`
	Decimals     int    `json:"decimals,omitempty"`
	ContractName string `json:"contractName,omitempty"`
}

type Utxos struct {
	ID          int
	Txid        string
	IndexN      int
	Address     string
	Value       string
	Asset       string
	Height      int
	SpentTxid   string
	SpentHeight int
	ClaimTxid   string
	ClaimHeight int
	Status      int
}

type Upt struct {
	ID           int
	Address      string
	Asset        string
	UpdateHeight int
}

type Balance struct {
	ID                int
	Address           string
	Asset             string
	Value             string
	LastUpdatedHeight int
}

type History struct {
	ID        int
	Txid      string
	Operation string
	IndexN    int
	Address   string
	Value     string
	Asset     string
	Timepoint int
}
