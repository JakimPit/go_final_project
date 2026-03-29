package server

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"go_final_project/pkg/api"
)

func Run(port string) {
	api.Init()

	// Swagger UI доступен по адресу /swagger/
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Сервер запущен на порту %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
