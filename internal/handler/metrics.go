package handler

import (
	"net/http"
	"strconv"
	"strings"

	models "github.com/valooocheck/metrics-yandex/internal/model"
	"github.com/valooocheck/metrics-yandex/internal/repository"
)

// MetricsHandler обрабатывает HTTP-запросы к серверу метрик.
type MetricsHandler struct {
	storage repository.Storage
}

// NewMetricsHandler создаёт хендлер с переданным хранилищем.
func NewMetricsHandler(storage repository.Storage) *MetricsHandler {
	return &MetricsHandler{storage: storage}
}

// ServeHTTP разбирает путь /update/<тип>/<имя>/<значение> вручную,
// чтобы избежать автоматических редиректов http.ServeMux.
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	// parts: [update, <тип>, <имя>, <значение>]
	if len(parts) < 2 || parts[0] != "update" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(parts) < 4 {
		// Путь короче ожидаемого — считаем, что имя метрики отсутствует.
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metricType {
	case models.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.storage.UpdateCounter(metricName, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.storage.UpdateGauge(metricName, value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
