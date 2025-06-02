package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first  *ListItem
	last   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.length++

	item := &ListItem{
		Value: v,
		Next:  l.first,
	}
	if l.first != nil {
		l.first.Prev = item
	} else {
		l.last = item
	}
	l.first = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.length++
	item := &ListItem{
		Value: v,
		Prev:  l.last,
	}

	if l.last != nil {
		l.last.Next = item
	} else {
		l.first = item
	}
	l.last = item

	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.first = i.Next
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
