package typedsyncmap

import "sync"

type typedSyncMap[K comparable, V any] struct {
	sync.Map
}

type TypeMapInterface[K comparable, V any] interface {
	Load(key K) (value V, ok bool)
	Store(key K, value V)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	LoadAndDelete(key K) (value V, loaded bool)
	Delete(key K)
	Swap(key K, value V) (previous V, loaded bool)
	CompareAndSwap(key K, old V, new V) (swapped bool)
	CompareAndDelete(key K, old V) (deleted bool)
	Range(func(key K, value V) (shouldContinue bool))
	Clear()
}

func NewTypedSyncMap[K comparable, V any]() TypeMapInterface[K, V] {
	return &typedSyncMap[K, V]{}
}

func (t *typedSyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := t.Map.Load(key)
	if !ok {
		return value, false
	}
	return v.(V), true
}

func (t *typedSyncMap[K, V]) Store(key K, value V) {
	t.Map.Store(key, value)
}

func (t *typedSyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, l := t.Map.LoadOrStore(key, value)
	return a.(V), l
}

func (t *typedSyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, l := t.Map.LoadAndDelete(key)
	if !l {
		return value, false
	}
	return v.(V), true
}

func (t *typedSyncMap[K, V]) Delete(key K) {
	t.Map.Delete(key)
}

func (t *typedSyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	p, l := t.Map.Swap(key, value)
	if !l {
		return previous, false
	}
	return p.(V), true
}

func (t *typedSyncMap[K, V]) CompareAndSwap(key K, old V, new V) (swapped bool) {
	return t.Map.CompareAndSwap(key, old, new)
}

func (t *typedSyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return t.Map.CompareAndDelete(key, old)
}

func (t *typedSyncMap[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	t.Map.Range(
		func(k, v any) bool {
			// given the Store method ensures V is always stored, this should be safe
			return f(k.(K), v.(V))
		},
	)
}

// Clear removes all keys and values from the map.
func (t *typedSyncMap[K, V]) Clear() {
	t.Map.Range(
		func(key, value interface{}) bool {
			t.Map.Delete(key)
			return true
		},
	)
}
