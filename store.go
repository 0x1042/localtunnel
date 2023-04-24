package main

import (
	"sync"
)

type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Delete(key K) { m.m.Delete(key) }

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	if v, has := m.m.Load(key); has {
		return v.(V), has
	}
	return value, ok
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if v, ok := m.m.LoadAndDelete(key); ok {
		return v.(V), loaded
	}
	return value, loaded
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	act, ok := m.m.LoadOrStore(key, value)
	return act.(V), ok
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}
