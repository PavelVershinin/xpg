package cache

import (
	"strings"
	"sync"
	"time"
)

type cacheRow struct {
	validDate time.Time
	value     interface{}
}

var (
	mu   sync.Mutex
	list map[string]cacheRow
)

func init() {
	list = make(map[string]cacheRow)
}

// Set Сохранит в кеш значение value с ключем key на период periodValidity
func Set(key string, value interface{}, periodValidity time.Duration) {
	if periodValidity <= 0 {
		periodValidity = 3600 * time.Second
	}
	mu.Lock()
	list[key] = cacheRow{
		validDate: time.Now().Add(periodValidity),
		value:     value,
	}
	mu.Unlock()
}
// Get Получит из кеша значение по ключу key
func Get(key string) (res interface{}, ok bool) {
	mu.Lock()
	if data, exists := list[key]; exists {
		if ok = data.validDate.After(time.Now()); ok {
			res = data.value
		} else {
			delete(list, key)
		}
	}
	mu.Unlock()
	return res, ok
}

// Delete Удалит из кеша значение по ключу key
// Если keyFragment == true удалит все значения с частично совпавшими ключами
func Delete(key string, keyFragment bool) {
	go func() {
		mu.Lock()
		if key == "" && keyFragment {
			list = make(map[string]cacheRow)
		} else if !keyFragment {
			delete(list, key)
		} else {
			for k, _ := range list {
				if strings.Contains(k, key) {
					delete(list, k)
				}
			}
		}
		mu.Unlock()
	}()
}
