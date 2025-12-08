package database

import (
	"golang_daerah/config"

	"github.com/jmoiron/sqlx"
)

func InitAllDatabases() map[string]*sqlx.DB {
	dbs := make(map[string]*sqlx.DB)

	dbs["terminal"] = config.InitTerminalDBX()
	dbs["passenger"] = config.InitPassengerPlaneDBX()
	dbs["auth"] = config.InitAuthDBX()
	dbs["traffic"] = config.InitTrafficDBX()
	dbs["golang"] = config.InitGolangDBX()
	dbs["mysql"] = config.InitMySQLDBX()
	dbs["passanger"] = config.InitMySQLDBX_passanger()

	return dbs
}

func CloseAllDatabases(dbs map[string]*sqlx.DB) {
	for _, db := range dbs {
		db.Close()
	}
}
