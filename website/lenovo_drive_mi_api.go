package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/dingsongjie/go-project-template/configs"
	"github.com/dingsongjie/go-project-template/pkg/db"
	"github.com/dingsongjie/go-project-template/pkg/log"

	routers "github.com/dingsongjie/go-project-template/website/routes"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.GetLogger()
	godotenv.Load(".env")
	configs.ConfigDb(os.Getenv("USER_INFO_DB_CONNECTION_STRING"), os.Getenv("NEW_PATH_DB_CONNECTION_STRING"))
	configs.ConfigGin(os.Getenv("GIN_MODE"))
	logger := log.GetLogger()
	defer logger.Sync()

	gin.SetMode(configs.GinMode)
	r := gin.Default()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// 配置数据库

	db.AddDbClients()

	routers.AddRouter(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", zap.Error(err))
	}
	logger.Info("Server exiting")
}
