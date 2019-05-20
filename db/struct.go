package db

import "github.com/hzxiao/goutil"

type Status struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	UpdateHeight int    `json:"update_height,omitempty"`
}

type Block struct {
	Height      int        `json:"height,omitempty"`
	SysFee      int        `json:"sys_fee"`
	TotalSysFee int        `json:"total_sys_fee"`
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
	ContractName string `json:"contract_name,omitempty"`
}

type Utxos struct {
	ID          int    `json:"id,omitempty"`
	Txid        string `json:"txid,omitempty"`
	IndexN      int    `json:"index_n,omitempty"`
	Address     string `json:"address,omitempty"`
	Value       string `json:"value"`
	Asset       string `json:"asset,omitempty"`
	Height      int    `json:"height,omitempty"`
	SpentTxid   string `json:"spent_txid,omitempty"`
	SpentHeight int    `json:"spent_height,omitempty"`
	ClaimTxid   string `json:"claim_txid,omitempty"`
	ClaimHeight int    `json:"claim_height,omitempty"`
	Status      int    `json:"status"`
}

type Upt struct {
	ID           int    `json:"id,omitempty"`
	Address      string `json:"address,omitempty"`
	Asset        string `json:"asset,omitempty"`
	UpdateHeight int    `json:"update_height"`
}

type Balance struct {
	ID                int    `json:"id,omitempty"`
	Address           string `json:"address,omitempty"`
	Asset             string `json:"asset,omitempty"`
	Value             string `json:"value"`
	LastUpdatedHeight int    `json:"last_updated_height,omitempty"`
}

type History struct {
	ID        int    `json:"id,omitempty"`
	Txid      string `json:"txid,omitempty"`
	Operation string `json:"operation,omitempty"`
	IndexN    int    `json:"index_n,omitempty"`
	Address   string `json:"address,omitempty"`
	Value     string `json:"value"`
	Asset     string `json:"asset,omitempty"`
	Timepoint int    `json:"timepoint,omitempty"`
}
