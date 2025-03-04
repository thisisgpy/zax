package main

import (
	"zax/config"
	"zax/handler"
	"zax/repository"
	"zax/service"
	"zax/util"

	_ "github.com/go-sql-driver/mysql"
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

	// 初始化雪花算法
	idGen, _ := util.NewSnowflake(1)

	// 初始化事务助手
	txHelper := util.NewTxHelper(db)

	// 数据访问层初始化
	orgRepo := repository.NewOrgRepository(db)

	// 业务层初始化
	orgService := service.NewOrgService(logger, idGen, txHelper, orgRepo)

	// 控制器初始化
	orgHandler := handler.NewOrgHandler(orgService)

	// 初始化gin
	r := config.GinInit(logger)

	// 路由注册
	handler.RegisterOrgHandlers(r, orgHandler)

	logger.Info("Server started on port 8899")

	if err := r.Run(":8899"); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
