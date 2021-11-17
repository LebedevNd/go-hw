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
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	front := l.front
	if front == nil {
		front = l.back
	}

	l.front = &ListItem{
		v,
		front,
		nil,
	}

	if l.len == 0 {
		l.back = l.front
	}

	if front != nil {
		front.Prev = l.front
	}

	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	back := l.back
	if back == nil {
		back = l.front
	}
	l.back = &ListItem{
		v,
		nil,
		back,
	}

	if l.len == 0 {
		l.front = l.back
	}

	if back != nil {
		back.Next = l.back
	}

	l.len++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
