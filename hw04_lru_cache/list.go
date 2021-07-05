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

func NewListItem(v interface{}) *ListItem {
	return &ListItem{Value: v, Next: nil, Prev: nil}
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	node := NewListItem(v)
	if l.len == 0 {
		l.back = node
	} else {
		l.front.Prev = node
		node.Next = l.front
	}
	l.front = node
	l.len++
	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	node := NewListItem(v)
	if l.len == 0 {
		l.front = node
	} else {
		l.back.Next = node
		node.Prev = l.back
	}
	l.back = node
	l.len++
	return node
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
	i = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return &list{}
}

func intoSlice(l List) []int {
	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int)) // beware, can be panicked
	}
	return elems
}
