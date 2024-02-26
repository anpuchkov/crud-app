package main

import (
	"Api/configs"
	"Api/psql/database"
	"Api/psql/handler"
	"Api/server"
	"Api/service"
	"context"
	"log"
	"net/http"
)

func main() {
	cfg, err := configs.ConfigInit()
	if err != nil {
		log.Println("failed to initialize config: ", err)
		return
	}
	ctx := context.Background()

	db, err := database.InitPostgresConnection(ctx, *cfg.DBConfig)
	if err != nil {
		log.Println("unable to connect to the database: ", err)
		return
	}

	defer db.Close()

	repo := handler.NewService(db)
	serv := service.NewService(repo)
	handle := server.NewHandler(serv)

	ser := &http.Server{
		Addr:    ":8080",
		Handler: handle.InitHandler(),
	}

	log.Println("server is running on port 8080")

	if err := ser.ListenAndServe(); err != nil {
		log.Fatalf("ser error: %v", err)
	}
}
