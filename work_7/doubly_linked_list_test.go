package doubleLinkedList

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestSetLogger(t *testing.T) {
	var err error
	var log *zap.Logger
	if log, err = zap.NewDevelopment(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)
}

func (l *List) PrintList(w io.Writer) {
	fmt.Fprintf(w, "Список длинной %d первый элемент %d последний элемент %d\n", l.Len(), l.First().Value(), l.Last().Value())
	for item := l.First(); item != nil; item = item.next {
		fmt.Fprintf(w, "%d ", item.Value())
	}
	fmt.Fprintf(w, "\n")
}

func TestInsert(t *testing.T) {
	// t.Skip()
	writer := os.Stdout
	list := new(List)
	fmt.Fprintf(writer, "Делаем вствку элеметов с 31 -60 PushBack\n")
	for i := 31; i < 61; i++ {
		list.PushBack(i)
	}
	list.PrintList(writer)
	fmt.Fprintf(writer, "Делаем вствку элеметов с 1 -30 PushFront\n")
	for i := 1; i < 31; i++ {
		list.PushFront(i)
	}
	list.PrintList(writer)
}

func TestRemovae(t *testing.T) {
	// t.Skip()
	writer := os.Stdout
	list := new(List)
	fmt.Fprintf(writer, "Делаем новый список с 1 - 80 PushBack\n")
	for i := 1; i < 81; i++ {
		list.PushBack(i)
	}
	list.PrintList(writer)
	fmt.Fprintf(writer, "Удаляем случайные 8 элементов \n")
	elem := list.First()
	rand.Seed(int64(time.Now().Second()))
	var r int
	for i := 0; i < 8; i++ {
		r = rand.Intn(10)
		for j := 0; j < r; j++ {
			elem = elem.next
		}
		deleteElem := elem
		elem = elem.next
		deleteElem.Remove()
		fmt.Fprintf(writer, "Удален элемент %d\n", deleteElem.Value())
	}
	list.PrintList(writer)
	r = rand.Intn(10)
	fmt.Fprintf(writer, "Удаляем %d элементов в начале\n", r)
	for j := 0; j < r; j++ {
		list.First().Remove()
	}
	r = rand.Intn(10)
	fmt.Fprintf(writer, "Удаляем %d элементов в конце\n", r)
	for j := 0; j < r; j++ {
		list.Last().Remove()
	}
	list.PrintList(writer)
}
