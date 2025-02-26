package main

import (
	"fmt"
	"time"
	"zax/config"
	"zax/model"
	"zax/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	// 初始化日志
	zapLogger := config.InitLogger()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	// 初始化数据库
	db, err := config.InitDB()
	if err != nil {
		logger.Errorf("连接数据库失败: %v", err)
		return
	}

	dbUtil := service.NewDBUtil(db)

	org := model.SysOrg{
		ID:         6,
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

	e := dbUtil.RunTx(func(tx *sqlx.Tx) error {
		res, err := tx.NamedExec("INSERT INTO sys_org (id, code, name, name_abbr, comment, parent_id, is_deleted, create_time, create_by, update_time, update_by) VALUES (:id, :code, :name, :name_abbr, :comment, :parent_id, :is_deleted, :create_time, :create_by, :update_time, :update_by)", org)
		if err != nil {
			return err
		}
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		logger.Infof("插入组织成功.影响行数 %d", rows)
		return nil
	})
	if e != nil {
		fmt.Println(e)
	}
}
