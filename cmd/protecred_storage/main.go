package main

import (
	"flag"
	"log"

	"protected-storage-server/internal/app"
	"protected-storage-server/internal/common/utils"
	"protected-storage-server/internal/config"
)

func main() {
	utils.LoadEnvironments(".env")

	utils.HandleFlag()
	flag.Parse()

	serverAddress := utils.GetServerAddress()
	dbAddress := utils.GetDBAddress()
	key := utils.GetKey()

	cfg := config.New(serverAddress, key, dbAddress)
	grpcApp, err := app.NewGRPC(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = grpcApp.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
