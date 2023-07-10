package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	docs "lenovo-drive-mi-api/docs"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	UserInfoConnectionString string
	NewPathConnectionString  string
	logger                   *zap.Logger
)

func main() {
	logger, _ = zap.NewProductionConfig().Build()
	defer logger.Sync()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Failed to load environment variables: %v", zap.Error(err))
	}
	r := gin.Default()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	docs.SwaggerInfo.BasePath = "/"
	// 连接数据库
	UserInfoConnectionString = os.Getenv("USER_INFO_DB_CONNECTION_STRING")
	NewPathConnectionString = os.Getenv("NEW_PATH_DB_CONNECTION_STRING")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "health",
		})
	})
	r.POST("/getUserInfomation", getAllInformation)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
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
