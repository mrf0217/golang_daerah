package config

import (
    "database/sql"
    "github.com/jmoiron/sqlx"
)

func InitGolangDBX() *sqlx.DB {
    return sqlx.NewDb(InitGolangDB(), "postgres")
}

func InitTrafficDBX() *sqlx.DB {
    return sqlx.NewDb(InitTrafficDB(), "postgres")
}

func InitMySQLDBX() *sqlx.DB {
    return sqlx.NewDb(InitMySQLDB(), "mysql")
}

func InitPassengerPlaneDBX() *sqlx.DB {
    return sqlx.NewDb(InitPassengerPlaneDB(), "mysql")
}

func InitTerminalDBX() *sqlx.DB {
    return sqlx.NewDb(InitTerminalDB(), "mysql")
}