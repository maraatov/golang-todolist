package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	todo "todolists"
	"todolists/internal/db"
	"todolists/internal/handler"
	"todolists/internal/repository"
	"todolists/internal/service"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка инициализации конфига: %s", err.Error())
	}

	dbPool, err := db.Connect("postgres://postgres:Aldiyar2004@localhost:5432/practices?sslmode=disable")
	if err != nil {
		logrus.Fatalln(err)
		return
	}

	repos := repository.NewRepository(dbPool)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Ошибка старта сервера: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
