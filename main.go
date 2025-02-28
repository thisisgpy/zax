package main

import (
	"fmt"
	"sync"
	"zax/config"
	"zax/repository"
	"zax/service"
	"zax/util"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

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

	orgRepo := repository.NewOrgRepository(db)

	txHelper := util.NewTxHelper(db)

	orgService := service.NewOrgService(logger, txHelper, orgRepo)

	// sysOrg := &model.SysOrg{
	// 	ID:         100,
	// 	Code:       "0004",
	// 	Name:       "总公司4",
	// 	NameAbbr:   "总部4",
	// 	Comment:    "这是一个测试组织4",
	// 	ParentID:   0,
	// 	IsDeleted:  false,
	// 	CreateTime: time.Now(),
	// 	CreateBy:   "admin",
	// 	UpdateTime: time.Now(),
	// 	UpdateBy:   "admin",
	// }

	code, _ := orgService.GenerateOrgCode(100)
	fmt.Println(code)
	wg.Wait()
}
