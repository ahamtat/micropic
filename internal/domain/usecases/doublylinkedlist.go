package usecases

import (
	"sync"
)

type dllItem struct {
	Value interface{}
	Next  *dllItem
	Prev  *dllItem
}

type doublyLinkedList struct {
	mx     sync.Mutex
	head   *dllItem
	tail   *dllItem
	length int
}

func (l *doublyLinkedList) GetHead() *dllItem {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.head
}

func (l *doublyLinkedList) GetLength() int {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.length
}

// Push value to the head of doubly linked list
func (l *doublyLinkedList) PushHead(v interface{}) *dllItem {
	l.mx.Lock()
	defer l.mx.Unlock()

	i := &dllItem{
		Value: v,
		Next:  l.head,
	}
	if l.head != nil {
		l.head.Prev = i
	} else {
		// First item in list
		l.tail = i
	}
	l.head = i
	l.length++

	return i
}

// Move item to the head of doubly linked list
func (l *doublyLinkedList) MoveHead(i *dllItem) {
	l.mx.Lock()
	defer l.mx.Unlock()

	// Extract item from list
	if i.Prev != nil {
		i.Prev.Next = i.Next
		i.Prev = nil
	}

	// Move item to list's head
	l.head.Prev = i
	i.Next = l.head
	l.head = i
}

// Remove tail item and return its value
func (l *doublyLinkedList) PopTail() interface{} {
	l.mx.Lock()
	defer l.mx.Unlock()

	params := l.tail.Value
	if l.tail.Prev != nil {
		l.tail.Prev.Next = nil
	}
	l.length--

	return params
}

func (l *doublyLinkedList) Clean() {
	for l.GetLength() != 0 {
		_ = l.PopTail()
	}
}
