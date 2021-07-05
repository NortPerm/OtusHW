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
	List // Remove me after realization.
	// Place your code here.
}

func (l list) Len() int {
	return 0
}

func (l list) Front() *ListItem {
	return nil
}

func (l list) Back() *ListItem {
	return nil
}

func NewList() List {
	return new(list)
}
