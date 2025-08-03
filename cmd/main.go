package main

import (
	"DragDrop-Files/initialization"
	"DragDrop-Files/internal/delivery/http"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.Info("start init server")
	if err := initialization.LoadConfiguration(); err != nil {
		logrus.Fatal(err.Error())
	}
	if err := initialization.ErrorInitialization(); err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("end init server")
}

// @title           DragDropFiles
// @version         1.0.0
// @description     Сервис по сохранению файлов

// @host      		localhost:1941
// @BasePath  		/api/
func main() {
	var serverInstance server.Server
	routes, businessDatabase := initialization.InitLayers()
	go run(serverInstance, routes, &initialization.ConfigService.Server)
	stop()
	serverInstance.Stop(context.Background(), businessDatabase)
}

func run(server server.Server, routes *http.Http, config *entity.ServerConfig) {
	ginEngine := routes.InitHTTPHttps(config)
	certificates := initialization.ConfigService.Certificates

	if err := server.Run(config.Port, ginEngine, certificates); err != nil {
		if err.Error() != "http: Server closed" {
			logrus.Fatalf("error occurred while running http server: %s", nil, err.Error())
		}
	}
}

func stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
}
