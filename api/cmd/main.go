package main

import (
	"api/config"
	"api/infrastructure/mysql/db"
	"api/server"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := config.GetConfig()
	db.NewMainDB(conf.DB)

	server.Run(ctx)
}
