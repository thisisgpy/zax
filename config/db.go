package config

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "zax:7788uJmki*@tcp(rm-wz988oqn7627g91t3so.mysql.rds.aliyuncs.com:3306)/zax?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}
