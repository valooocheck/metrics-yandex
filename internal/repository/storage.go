package repository

import (
	"errors"

	models "github.com/valooocheck/metrics-yandex/internal/model"
)

// ErrNotFound возвращается, когда запрошенная метрика отсутствует в хранилище.
var ErrNotFound = errors.New("metric not found")

// Storage описывает интерфейс взаимодействия с хранилищем метрик.
type Storage interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, value int64) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAll() ([]models.Metrics, error)
}
