package app

import (
	"marketflow/internal/domain"
	"sync"
	"time"
)

type WindowStore struct {
	mu   sync.RWMutex
	data map[domain.Key][]domain.PriceTick
}

func NewWindowStore() domain.WindowStore {
	return &WindowStore{data: make(map[domain.Key][]domain.PriceTick)}
}

func (w *WindowStore) Add(msg domain.PriceTick) {
	k := domain.Key{Exchange: msg.Exchange, Symbol: msg.Symbol}
	w.mu.Lock()
	defer w.mu.Unlock()

	arr := w.data[k]
	arr = append(arr, msg)

	cutoff := msg.Ts.Add(-time.Minute)
	i := 0
	for i < len(arr) && arr[i].Ts.Before(cutoff) {
		i++
	}
	if i > 0 {
		arr = arr[i:]
	}
	w.data[k] = arr
}

func (w *WindowStore) Snapshot() map[domain.Key][]domain.PriceTick {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make(map[domain.Key][]domain.PriceTick, len(w.data))
	for k, v := range w.data {
		vv := make([]domain.PriceTick, len(v))
		copy(vv, v)
		out[k] = vv
	}
	return out
}

func (w *WindowStore) GetLatestPriceBySymbol(symbol string) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var (
		found  bool
		latest time.Time
		price  float64
	)

	for k, arr := range w.data {
		if k.Symbol != symbol || len(arr) == 0 {
			continue
		}
		last := arr[len(arr)-1]
		if !found || last.Ts.After(latest) {
			found = true
			latest = last.Ts
			price = last.Price
		}
	}

	if !found {
		return -1, domain.ErrNotFound
	}

	return price, nil
}

func (w *WindowStore) GetLatestPriceBySymbolAndExchange(symbol, exchange string) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	arr, ok := w.data[domain.Key{Exchange: exchange, Symbol: symbol}]
	if !ok || len(arr) == 0 {
		return 0, domain.ErrNotFound
	}
	return arr[len(arr)-1].Price, nil
}

func (w *WindowStore) GetHighestPriceBySymbolAndPeriod(symbol string, period time.Duration) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cutoff := time.Now().Add(-period)

	var (
		found  bool
		maxVal float64
	)

	for k, arr := range w.data {
		if k.Symbol != symbol || len(arr) == 0 {
			continue
		}

		// Тики в arr идут по времени по возрастанию (мы append'им и чистим старые в начале).
		// Идём с конца, пока находим точки внутри окна.
		for i := len(arr) - 1; i >= 0; i-- {
			t := arr[i]
			if t.Ts.Before(cutoff) {
				// дальше только ещё более старые тики -> выходим из этого массива
				break
			}
			if !found || t.Price > maxVal {
				found = true
				maxVal = t.Price
			}
		}
	}

	if !found {
		return 0, domain.ErrNotFound
	}
	return maxVal, nil
}

// -------- HIGHEST (symbol+period+exchange) --------

func (w *WindowStore) GetHighestPriceBySymbolAndPeriodAndExchange(symbol, exchange string, period time.Duration) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cutoff := time.Now().Add(-period)

	arr, ok := w.data[domain.Key{Exchange: exchange, Symbol: symbol}]
	if !ok || len(arr) == 0 {
		return 0, domain.ErrNotFound
	}

	found := false
	var maxVal float64
	for i := len(arr) - 1; i >= 0; i-- {
		t := arr[i]
		if t.Ts.Before(cutoff) {
			break
		}
		if !found || t.Price > maxVal {
			found = true
			maxVal = t.Price
		}
	}
	if !found {
		return 0, domain.ErrNotFound
	}
	return maxVal, nil
}

// ===================== LOWEST ======================

func (w *WindowStore) GetLowestPriceBySymbol(symbol string) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	found := false
	var minVal float64

	for k, arr := range w.data {
		if k.Symbol != symbol || len(arr) == 0 {
			continue
		}
		// без периода — по всем тикам, что в окне стора (обычно ~последняя минута)
		for _, t := range arr {
			if !found || t.Price < minVal {
				found = true
				minVal = t.Price
			}
		}
	}
	if !found {
		return 0, domain.ErrNotFound
	}
	return minVal, nil
}

func (w *WindowStore) GetLowestPriceBySymbolAndExchange(symbol, exchange string) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	arr, ok := w.data[domain.Key{Exchange: exchange, Symbol: symbol}]
	if !ok || len(arr) == 0 {
		return 0, domain.ErrNotFound
	}

	found := false
	var minVal float64
	for _, t := range arr {
		if !found || t.Price < minVal {
			found = true
			minVal = t.Price
		}
	}
	if !found {
		return 0, domain.ErrNotFound
	}
	return minVal, nil
}

func (w *WindowStore) GetLowestPriceBySymbolAndPeriod(symbol string, period time.Duration) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cutoff := time.Now().Add(-period)

	found := false
	var minVal float64
	for k, arr := range w.data {
		if k.Symbol != symbol || len(arr) == 0 {
			continue
		}
		for i := len(arr) - 1; i >= 0; i-- {
			t := arr[i]
			if t.Ts.Before(cutoff) {
				break
			}
			if !found || t.Price < minVal {
				found = true
				minVal = t.Price
			}
		}
	}
	if !found {
		return 0, domain.ErrNotFound
	}
	return minVal, nil
}

func (w *WindowStore) GetLowestPriceBySymbolAndPeriodAndExchange(symbol, exchange string, period time.Duration) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cutoff := time.Now().Add(-period)

	arr, ok := w.data[domain.Key{Exchange: exchange, Symbol: symbol}]
	if !ok || len(arr) == 0 {
		return 0, domain.ErrNotFound
	}

	found := false
	var minVal float64
	for i := len(arr) - 1; i >= 0; i-- {
		t := arr[i]
		if t.Ts.Before(cutoff) {
			break
		}
		if !found || t.Price < minVal {
			found = true
			minVal = t.Price
		}
	}
	if !found {
		return 0, domain.ErrNotFound
	}
	return minVal, nil
}

// ===================== AVERAGE =====================

func (w *WindowStore) GetAvgPriceBySymbol(symbol string) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var sum float64
	var cnt int64

	for k, arr := range w.data {
		if k.Symbol != symbol || len(arr) == 0 {
			continue
		}
		for _, t := range arr {
			sum += t.Price
			cnt++
		}
	}
	if cnt == 0 {
		return 0, domain.ErrNotFound
	}
	return sum / float64(cnt), nil
}

func (w *WindowStore) GetAvgPriceBySymbolAndExchange(symbol, exchange string) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	arr, ok := w.data[domain.Key{Exchange: exchange, Symbol: symbol}]
	if !ok || len(arr) == 0 {
		return 0, domain.ErrNotFound
	}

	var sum float64
	for _, t := range arr {
		sum += t.Price
	}
	return sum / float64(len(arr)), nil
}

func (w *WindowStore) GetAvgPriceBySymbolAndPeriod(symbol string, period time.Duration) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cutoff := time.Now().Add(-period)

	var sum float64
	var cnt int64
	for k, arr := range w.data {
		if k.Symbol != symbol || len(arr) == 0 {
			continue
		}
		for i := len(arr) - 1; i >= 0; i-- {
			t := arr[i]
			if t.Ts.Before(cutoff) {
				break
			}
			sum += t.Price
			cnt++
		}
	}
	if cnt == 0 {
		return 0, domain.ErrNotFound
	}
	return sum / float64(cnt), nil
}

func (w *WindowStore) GetAvgPriceBySymbolAndPeriodAndExchange(symbol, exchange string, period time.Duration) (float64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	cutoff := time.Now().Add(-period)

	arr, ok := w.data[domain.Key{Exchange: exchange, Symbol: symbol}]
	if !ok || len(arr) == 0 {
		return 0, domain.ErrNotFound
	}

	var sum float64
	var cnt int64
	for i := len(arr) - 1; i >= 0; i-- {
		t := arr[i]
		if t.Ts.Before(cutoff) {
			break
		}
		sum += t.Price
		cnt++
	}
	if cnt == 0 {
		return 0, domain.ErrNotFound
	}
	return sum / float64(cnt), nil
}
