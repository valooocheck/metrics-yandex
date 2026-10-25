package main

import (
	"log"
	"net/http"

	"github.com/valooocheck/metrics-yandex/internal/handler"
	"github.com/valooocheck/metrics-yandex/internal/repository"
)

func main() {
	storage := repository.NewMemStorage()
	h := handler.NewMetricsHandler(storage)

	// Хендлер передаётся напрямую, без http.ServeMux: мультиплексор
	// нормализует путь и отвечает редиректом (307/301), а редиректы
	// по заданию не поддерживаются.
	log.Println("server listening on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", h); err != nil {
		log.Fatal(err)
	}
}
