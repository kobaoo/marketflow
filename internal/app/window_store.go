package app

import (
	"marketflow/internal/domain"
	"sync"
	"time"
)

type WindowStore struct {
    mu sync.RWMutex
    data map[domain.Key][]domain.Tick
}

func NewWindowStore() domain.WindowStore {
    return &WindowStore{data: make(map[domain.Key][]domain.Tick)}
}

func (w *WindowStore) Add(exchange, symbol string, price float64, ts time.Time) {
    k := domain.Key{Exchange:exchange, Symbol: symbol}
    w.mu.Lock()
    defer w.mu.Unlock()

    arr := w.data[k]
    arr = append(arr, domain.Tick{Ts: ts, Price: price})

    cutoff := ts.Add(-time.Minute)
    i := 0
    for i < len(arr) && arr[i].Ts.Before(cutoff) {
        i++
    }
    if i > 0 {
        arr = arr[i:]
    }
    w.data[k] = arr
}

func (w *WindowStore) Snapshot() map[domain.Key][]domain.Tick {
    w.mu.RLock()
    defer w.mu.RUnlock()
    out := make(map[domain.Key][]domain.Tick, len(w.data))
    for k, v := range w.data {
        vv := make([]domain.Tick, len(v))
        copy(vv, v)
        out[k] = vv
    }
    return out
}

func (w *WindowStore) Latest(exchange, symbol string) (float64, bool) {
    k := domain.Key{Exchange: exchange, Symbol: symbol}
    w.mu.RLock()
    defer w.mu.RUnlock()
    arr, ok := w.data[k]
    if !ok || len(arr) == 0 {
        return 0, false
    }
    return arr[len(arr)-1].Price, true
}
