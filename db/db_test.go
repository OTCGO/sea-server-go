package db

func init() {
	err := InitDB("root:123456@/mysql?charset=utf8&loc=Local&parseTime=true")
	if err != nil {
		panic(err)
	}

	_, err = db.engine.ImportFile("../sea_test.sql")
	if err != nil {
		panic(err)
	}
	db.engine.Close()
	db.uri = "root:123456@/sea_test?charset=utf8&loc=Local&parseTime=true"
	err = db.reconnect()
	if err != nil {
		panic(err)
	}
	db.engine.ShowSQL(true)
}

func deleteAll() {
	db.engine.Table(TableBlock).Exec("delete from block")
	db.engine.Table(TableStatus).Exec("delete from status")
	db.engine.Table(TableAssets).Exec("delete from assets")
	db.engine.Table(TableUpt).Exec("delete from upt")
	db.engine.Table(TableUtxos).Exec("delete from utxos")
	db.engine.Table(TableBalance).Exec("delete from balance")
}
