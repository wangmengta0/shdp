package main

import (
	"context"
	"fmt"
	"log"
	"shdp/internal/middle/RabbitMQ"
	"shdp/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"shdp/config"
	"shdp/internal/handler"
	"shdp/internal/repository"
	"shdp/internal/router"
	"shdp/internal/service"
)

func main() {
	// 1. 初始化配置
	config.InitConfig()

	// 2. 初始化 MySQL (GORM)
	db := initMySQL()

	// 3. 初始化 Redis (go-redis)
	rdb := initRedis()

	// 4. 依赖注入：组装 MVC 各层
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo, rdb)
	userHandler := handler.NewUserHandler(userService)

	voucherRepo := repository.NewVoucherRepo(db)
	voucherService := service.NewVoucherService(voucherRepo, rdb)
	voucherHandler := handler.NewVoucherHandler(voucherService)

	RabbitMQ.InitRabbitMQ()
	RabbitMQ.StartCacheDeleteConsumer(rdb)
	RabbitMQ.StartSeckillOrderConsumer(voucherRepo)
	handlers := handler.NewGroup(userHandler, voucherHandler)

	err := utils.InitSnowflake(1)
	if err != nil {
		log.Fatalf("初始化雪花算法失败: %v", err)
	}
	// 5. 初始化 Gin 引擎
	r := gin.Default()

	// 6. 注册路由 (将 Redis 客户端传入以便中间件使用)
	router.SetUpRouter(r, rdb, handlers)

	// 7. 启动服务
	addr := fmt.Sprintf(":%d", config.Conf.Server.Port)
	log.Printf("服务启动成功，监听端口 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// initMySQL 初始化 MySQL 连接
func initMySQL() *gorm.DB {
	dsn := config.Conf.MySQL.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印 SQL 方便调试
	})
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	// 设置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(config.Conf.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Conf.MySQL.MaxOpenConns)

	log.Println("MySQL 连接成功！")

	return db
}

// initRedis 初始化 Redis 连接
func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	// 测试连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 连接成功！")
	return rdb
}
