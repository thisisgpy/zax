package main

import (
	"encoding/json"
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
	// 	ID:         3,
	// 	Code:       "0003",
	// 	Name:       "测试组织3",
	// 	NameAbbr:   "测试3",
	// 	Comment:    "这是一个测试组织3",
	// 	ParentID:   1,
	// 	CreateTime: time.Now(),
	// 	CreateBy:   "admin",
	// 	UpdateTime: time.Now(),
	// 	UpdateBy:   "admin",
	// }

	orgTrees, err := orgService.FindOrgTrees()
	if err != nil {
		logger.Errorf("查询组织树失败: %v", err)
		return
	}
	jsonData, _ := json.Marshal(orgTrees)
	fmt.Println(string(jsonData))
	wg.Wait()
}
