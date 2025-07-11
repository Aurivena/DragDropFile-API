package main

import (
	"DragDrop-Files/initialization"
	"DragDrop-Files/internal/action"
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/persistence"
	"DragDrop-Files/internal/route"
	"DragDrop-Files/models"
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

	businessDatabase := persistence.NewBusinessDatabase(initialization.ConfigService)
	minioClient := persistence.NewMinioStorage(initialization.ConfigService.Minio)
	sources := persistence.Sources{
		BusinessDB: businessDatabase,
	}
	persistences := persistence.NewPersistence(&sources)
	domains := domain.NewDomain(persistences, initialization.ConfigService, minioClient)
	actions := action.NewAction(domains)

	routes := route.NewRoute(actions)
	certificates := initialization.ConfigService.Certificates

	go run(serverInstance, routes, &initialization.ConfigService.Server, certificates)
	stop()
	serverInstance.Stop(context.Background(), businessDatabase)
}

func run(server server.Server, routes *route.Route, config *models.ServerConfig, certificates models.CertificatesConfig) {
	ginEgine := routes.InitHTTPRoutes(config)

	if err := server.Run(config.Port, ginEgine, certificates); err != nil {
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
