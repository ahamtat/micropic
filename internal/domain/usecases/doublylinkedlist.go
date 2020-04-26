package usecases

import (
	"sync"
)

type DllItem struct {
	Value interface{}
	Next  *DllItem
	Prev  *DllItem
}

type DoublyLinkedList struct {
	mx     sync.Mutex
	head   *DllItem
	Tail   *DllItem
	Length int
}

func (l *DoublyLinkedList) GetHead() *DllItem {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.head
}

func (l *DoublyLinkedList) GetLength() int {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.Length
}

// Push value to the head of doubly linked list.
func (l *DoublyLinkedList) PushHead(v interface{}) *DllItem {
	l.mx.Lock()
	defer l.mx.Unlock()

	i := &DllItem{
		Value: v,
		Next:  l.head,
	}
	if l.head != nil {
		l.head.Prev = i
	} else {
		// First item in list
		l.Tail = i
	}
	l.head = i
	l.Length++

	return i
}

// Move item to the head of doubly linked list.
func (l *DoublyLinkedList) MoveHead(i *DllItem) {
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

// Remove tail item and return its value.
func (l *DoublyLinkedList) PopTail() interface{} {
	l.mx.Lock()
	defer l.mx.Unlock()

	params := l.Tail.Value
	if l.Tail.Prev != nil {
		l.Tail.Prev.Next = nil
	}
	l.Length--

	return params
}

func (l *DoublyLinkedList) Clean() {
	for l.GetLength() != 0 {
		_ = l.PopTail()
	}
}
