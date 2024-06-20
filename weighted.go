package robin

import (
	"cmp"
	"slices"
	"sync"
	"sync/atomic"
)

// WItem - generic weighted item to be returned by randomizer
type WItem[T any] struct {
	item        T
	weight      int
	occurrences int
}

// WRR - weighted round robin randomizer
type WRR[T any] struct {
	cl    atomic.Int32
	cr    atomic.Int32
	items []WItem[T]
	mtx   sync.RWMutex
}

func (w *WRR[T]) Next() T {
	var res T
	w.mtx.RLock()
	defer w.mtx.RUnlock()
	for i, v := range w.items {
		if v.occurrences%int(w.cl.Load()) < v.weight {
			w.items[i].occurrences++
			res = v.item
			w.cr.Store(w.cr.Add(1))
			break
		}
	}
	if w.cr.Load() >= w.cl.Load() {
		w.cr.Store(0)
		for i := range w.items {
			w.items[i].occurrences = 0
		}
	}
	return res
}

func (w *WRR[T]) Add(item T, weight int) {
	w.cl.Store(w.cl.Add(int32(weight)))
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.items = append(w.items, WItem[T]{item: item, weight: weight})
	slices.SortFunc(w.items, func(i, j WItem[T]) int { return cmp.Compare(j.weight, i.weight) })
}

// Reset resets the WRR to its initial state.
// It clears all items, resets the cycle length to 0, and releases the lock on the internal mutex.
// This function is safe for concurrent use.
func (w *WRR[T]) Reset() {
	w.cl.Store(0)
	w.cr.Store(0)
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.items = w.items[:0]
}
