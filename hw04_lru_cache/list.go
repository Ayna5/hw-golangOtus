package hw04_lru_cache //nolint:golint,stylecheck,revive

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый элемент списка
	Back() *listItem                   // последний элемент списка
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало

}

type listItem struct {
	Value interface{} // значение
	Next  *listItem   // следующий элемент
	Prev  *listItem   // предыдущий элемент
}

type list struct {
	first, last *listItem
	len         int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.first
}

func (l *list) Back() *listItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *listItem {
	item := &listItem{
		Value: v,
	}
	if l.len == 0 {
		l.last = item
		l.first = item
		l.len++
		return item
	}
	item.Prev = nil
	item.Next = l.first
	l.first.Prev = item
	l.first = item
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *listItem {
	item := &listItem{
		Value: v,
	}
	if l.len == 0 {
		l.first = item
		l.last = item
		l.len++
		return item
	}
	item.Prev = l.last
	item.Next = nil
	l.last.Next = item
	l.last = item
	l.len++
	return item
}

func (l *list) Remove(i *listItem) {
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if i.Prev == nil {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{}
}
