package repository

import (
	"sort"
	"sync"

	models "github.com/valooocheck/metrics-yandex/internal/model"
)

// MemStorage — in-memory реализация Storage.
type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage создаёт пустое хранилище метрик.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateGauge замещает предыдущее значение gauge.
func (s *MemStorage) UpdateGauge(name string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
	return nil
}

// UpdateCounter прибавляет значение к предыдущему значению counter.
func (s *MemStorage) UpdateCounter(name string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[name] += value
	return nil
}

// GetGauge возвращает значение gauge по имени.
func (s *MemStorage) GetGauge(name string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.gauges[name]
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

// GetCounter возвращает значение counter по имени.
func (s *MemStorage) GetCounter(name string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.counters[name]
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

// GetAll возвращает снимок всех метрик, отсортированный по типу и имени.
func (s *MemStorage) GetAll() ([]models.Metrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Metrics, 0, len(s.gauges)+len(s.counters))

	for name, value := range s.gauges {
		v := value
		result = append(result, models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &v,
		})
	}

	for name, value := range s.counters {
		d := value
		result = append(result, models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &d,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].MType != result[j].MType {
			return result[i].MType < result[j].MType
		}
		return result[i].ID < result[j].ID
	})

	return result, nil
}
