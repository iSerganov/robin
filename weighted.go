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
	cycleWeight atomic.Int32
	cr          atomic.Int32
	items       []WItem[T]
	mtx         sync.RWMutex
}

// Next returns the next item from the weighted round robin randomizer.
// It iterates over the items, checks if the current item's occurrences modulo cycle length is less than its weight,
// if true, it increments the occurrences of the current item, sets the result to the current item, increments the cycle counter,
// and breaks the loop. If the cycle counter is greater than or equal to the cycle length, it resets the cycle counter and
// occurrences of all items to 0.
//
// The function is thread-safe and uses a read lock on the internal mutex to ensure that the items slice is not modified
// concurrently.
//
// The function returns the next item from the weighted round robin randomizer.
func (w *WRR[T]) Next() T {
	var res T
	w.mtx.RLock()
	for i, v := range w.items {
		if v.occurrences%int(w.cycleWeight.Load()) < v.weight {
			w.items[i].occurrences++
			res = v.item
			_ = w.cr.Add(1)
			break
		}
	}
	w.mtx.RUnlock()
	if w.cr.Load() >= w.cycleWeight.Load() {
		w.cr.Store(0)
		w.mtx.RLock()
		for i := range w.items {
			w.items[i].occurrences = 0
		}
		w.mtx.RUnlock()
	}
	return res
}

// Add adds a new item with its weight to the weighted round robin randomizer.
// It increments the cycle length by the given weight, acquires a write lock on the internal mutex,
// appends the new item to the items slice, sorts the items slice in descending order of weight,
// and releases the lock.
//
// Parameters:
// - item: The item to be added to the randomizer.
// - weight: The weight of the item. The higher the weight, the more likely the item will be selected.
//
// Return:
// This function does not return any value.
func (w *WRR[T]) Add(item T, weight int) {
	_ = w.cycleWeight.Add(int32(weight))
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.items = append(w.items, WItem[T]{item: item, weight: weight})
	slices.SortFunc(w.items, func(i, j WItem[T]) int { return cmp.Compare(j.weight, i.weight) })
}

// Reset resets the WRR to its initial state.
// It clears all items, resets the cycle length to 0, and releases the lock on the internal mutex.
// This function is safe for concurrent use.
func (w *WRR[T]) Reset() {
	w.cycleWeight.Store(0)
	w.cr.Store(0)
	w.mtx.Lock()
	defer w.mtx.Unlock()
	w.items = w.items[:0]
}
