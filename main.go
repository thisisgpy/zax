package main

import (
	"fmt"
	"time"
	"zax/model"
	"zax/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {

	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", "zax:7788uJmki*@tcp(rm-wz988oqn7627g91t3so.mysql.rds.aliyuncs.com:3306)/zax?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		fmt.Printf("连接数据库失败: %v", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	txWrapper := service.NewTxWrapper(db)

	org := model.SysOrg{
		ID:         1,
		Code:       "0001",
		Name:       "测试组织",
		NameAbbr:   "测试",
		Comment:    "测试组织",
		ParentID:   0,
		IsDeleted:  false,
		CreateTime: time.Now(),
		CreateBy:   "admin",
		UpdateTime: time.Now(),
		UpdateBy:   "admin",
	}

	e := txWrapper.RunTx(func(tx *sqlx.Tx) error {
		res, err := tx.NamedExec("INSERT INTO sys_org (id, code, name, name_abbr, comment, parent_id, is_deleted, create_time, create_by, update_time, update_by) VALUES (:id, :code, :name, :name_abbr, :comment, :parent_id, :is_deleted, :create_time, :create_by, :update_time, :update_by)", org)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		fmt.Println(rows)
		return nil
	})
	if e != nil {
		fmt.Println(e)
	}
}
