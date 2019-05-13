package db

import "github.com/hzxiao/goutil"

type Status struct {
	ID           int
	Name         string
	UpdateHeight int
}

type Block struct {
	Height      int        `json:"height,omitempty"`
	SysFee      int        `json:"sysFee"`
	TotalSysFee int        `json:"totalSysFee"`
	Raw         goutil.Map `json:"-"`
}

type Assets struct {
	ID           int    `json:"id,omitempty"`
	Asset        string `json:"asset,omitempty"`
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	Version      string `json:"version,omitempty"`
	Decimals     int    `json:"decimals"`
	ContractName string `json:"contractName,omitempty"`
}

type Utxos struct {
	ID          int    `json:"id,omitempty"`
	Txid        string `json:"txid,omitempty"`
	IndexN      int    `json:"indexN,omitempty"`
	Address     string `json:"address,omitempty"`
	Value       string `json:"value"`
	Asset       string `json:"asset,omitempty"`
	Height      int    `json:"height,omitempty"`
	SpentTxid   string `json:"spentTxid,omitempty"`
	SpentHeight int    `json:"spentHeight,omitempty"`
	ClaimTxid   string `json:"claimTxid,omitempty"`
	ClaimHeight int    `json:"claimHeight,omitempty"`
	Status      int    `json:"status"`
}

type Upt struct {
	ID           int
	Address      string
	Asset        string
	UpdateHeight int
}

type Balance struct {
	ID                int    `json:"id,omitempty"`
	Address           string `json:"address,omitempty"`
	Asset             string `json:"asset,omitempty"`
	Value             string `json:"value"`
	LastUpdatedHeight int    `json:"lastUpdatedHeight,omitempty"`
}

type History struct {
	ID        int    `json:"id,omitempty"`
	Txid      string `json:"txid,omitempty"`
	Operation string `json:"operation,omitempty"`
	IndexN    int    `json:"indexN,omitempty"`
	Address   string `json:"address,omitempty"`
	Value     string `json:"value"`
	Asset     string `json:"asset,omitempty"`
	Timepoint int    `json:"timepoint,omitempty"`
}
