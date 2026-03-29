// @title           Планировщик задач API
// @version         1.0
// @description     REST API для управления задачами с поддержкой повторений
// @host            localhost:7540
// @BasePath        /

package main

import (
	"log"
	"os"

	_ "go_final_project/docs"
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	dbFile := "scheduler.db"
	if envDB := os.Getenv("TODO_DBFILE"); envDB != "" {
		dbFile = envDB
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	server.Run(port)
}
