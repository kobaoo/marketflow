package app

import (
	"sync"
	"time"
)


type key struct {
    Exchange string
    Symbol   string
}

type tick struct {
    Ts    time.Time
    Price float64
}

type WindowStore struct {
    mu   sync.RWMutex
    data map[key][]tick 
}

func NewWindowStore() *WindowStore {
    return &WindowStore{data: make(map[key][]tick)}
}

func (w *WindowStore) Add(exchange, symbol string, price float64, ts time.Time) {
    k := key{exchange, symbol}
    w.mu.Lock()
    defer w.mu.Unlock()

    arr := w.data[k]
    arr = append(arr, tick{Ts: ts, Price: price})

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

func (w *WindowStore) Snapshot() map[key][]tick {
    w.mu.RLock()
    defer w.mu.RUnlock()
    out := make(map[key][]tick, len(w.data))
    for k, v := range w.data {
        vv := make([]tick, len(v))
        copy(vv, v)
        out[k] = vv
    }
    return out
}

func (w *WindowStore) Latest(exchange, symbol string) (float64, bool) {
    k := key{exchange, symbol}
    w.mu.RLock()
    defer w.mu.RUnlock()
    arr, ok := w.data[k]
    if !ok || len(arr) == 0 {
        return 0, false
    }
    return arr[len(arr)-1].Price, true
}
