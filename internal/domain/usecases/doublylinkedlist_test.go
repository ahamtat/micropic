package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testValues = []struct {
	value interface{}
}{
	{2},
	{12},
	{85},
	{06},
}

func createAndFillList(t *testing.T) *doublyLinkedList {
	list := &doublyLinkedList{}
	for _, v := range testValues {
		item := list.PushHead(v)
		require.NotNil(t, item, "item should be not nil")
	}
	return list
}

func TestDoublyLinkedList_PushHead(t *testing.T) {
	list := createAndFillList(t)
	require.Equal(t, len(testValues), list.length, "lengths should be equal")
}

func TestDoublyLinkedList_PopTail(t *testing.T) {
	list := createAndFillList(t)
	v := list.PopTail()
	require.Equal(t, testValues[0], v, "values should be equal")
	require.Equal(t, len(testValues)-1, list.length, "lengths should be equal")
}

func TestDoublyLinkedList_GetHead(t *testing.T) {
	list := createAndFillList(t)
	require.NotNil(t, list.GetHead(), "list head should be not nil")
}

func TestDoublyLinkedList_GetLength(t *testing.T) {
	list := createAndFillList(t)
	require.Equal(t, len(testValues), list.GetLength(), "lengths should be equal")
}

func TestDoublyLinkedList_Clean(t *testing.T) {
	list := createAndFillList(t)
	list.Clean()
	require.Equal(t, 0, list.length, "length should be zero")
}

func TestDoublyLinkedList_MoveHead(t *testing.T) {
	list := createAndFillList(t)
	tail := list.tail
	list.MoveHead(tail)
	v := list.GetHead().Value
	require.Equal(t, testValues[0], v, "values should be equal")
}
