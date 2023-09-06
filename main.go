package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"webb-auth/conf"
	"webb-auth/db"
	"webb-auth/redis"
	"webb-auth/router"
	"webb-auth/wechat"

	"github.com/fhy/utils-golang/utils"
	"github.com/fhy/utils-golang/wggo"

	"github.com/sirupsen/logrus"
)

func init() {
	env := flag.String("env", "local", "Configuration file to use")
	flag.Parse()
	logrus.Info("init auth's config")
	conf.Init(*env)
	utils.InitLogger(&conf.Conf.Log)
	db.Init(conf.Conf.DbType, conf.Conf.Dbconfig)
	redis.Init(&conf.Conf.Redis)
	wechat.Init(&conf.Conf.Redis, &conf.Conf.WeChat)
}

func main() {
	wggo.Wg.Wait()
	logrus.Info("Authentication server Start")

	// sync.Once is used to make sure we can close the channel at different execution stages(SIGTERM or when the config is loaded).
	type closeOnce struct {
		C     chan struct{}
		once  sync.Once
		Close func()
	}
	// Wait until the server is ready to handle reloading.
	reloadReady := &closeOnce{
		C: make(chan struct{}),
	}
	reloadReady.Close = func() {
		reloadReady.once.Do(func() {
			close(reloadReady.C)
		})
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Conf.Server.IP, conf.Conf.Server.Port),
		Handler: router.GetGinEngine(),
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	<-quit
	logrus.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	wggo.Wg.Wait()
	logrus.Println("Server exiting")
}
