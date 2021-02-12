package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mx       *sync.Mutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mx.Lock()
	defer l.mx.Unlock()

	item, ok := l.items[key]
	if !ok {
		newItem := cacheItem{
			key:   key,
			value: value,
		}
		itemList := l.queue.PushFront(newItem)
		l.items[key] = itemList
		if l.queue.Len() > l.capacity {
			lastItem := l.queue.Back()
			l.queue.Remove(lastItem)
			delete(l.items, lastItem.Value.(cacheItem).key)
		}
		return false
	}
	item.Value = cacheItem{
		key:   key,
		value: value,
	}
	l.queue.PushFront(item)
	return true
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mx.Lock()
	defer l.mx.Unlock()

	item, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.queue.MoveToFront(item)
	return item.Value.(cacheItem).value, true
}

func (l *lruCache) Clear() {
	l.capacity = 0
	l.queue = &list{}
	l.items = make(map[Key]*listItem)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	if capacity > 0 {
		return &lruCache{
			capacity: capacity,
			queue:    &list{},
			items:    make(map[Key]*listItem),
			mx:       &sync.Mutex{},
		}
	}
	return nil
}
