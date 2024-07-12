package main

import (
	"chatService/pkg"
	"chatService/pkg/handler"
	"chatService/pkg/repository"
	"chatService/pkg/service"
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chanOs := make(chan os.Signal)
	signal.Notify(chanOs, syscall.SIGINT, syscall.SIGTERM)
	initConfig()
	db := initDb()
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)
	server := new(pkg.Server)
	go func() {
		if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to run server: %s", err)
		}
	}()
	<-chanOs
	err := server.ShutDown(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	cancel()
}

func initConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func initDb() *sqlx.DB {
	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err)
	}
	return db
}
