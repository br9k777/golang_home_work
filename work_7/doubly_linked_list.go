package doubleLinkedList

// var (
// 	log *zap.Logger
// )

// func init() {
// 	var err error
// 	if log, err = zap.NewProduction(); err != nil {
// 		fmt.Fprint(os.Stderr, err)
// 		os.Exit(1)
// 	}
// 	zap.ReplaceGlobals(log)
// }

//List Контайнер элементов
type List struct {
	lenght      int   // длинна списка
	first, last *Item // ссылки на первый и последний элементы списка
}

//Len длинна списка
func (l *List) Len() int {
	return l.lenght
}

//First первый Item
func (l *List) First() *Item {
	return l.first
}

//Last последний Item
func (l *List) Last() *Item {
	return l.last
}

//PushFront добавить значение в начало
func (l *List) PushFront(v interface{}) {
	item := &Item{
		value: v,
		next:  l.First(),
		list:  l,
	}
	if item.next == nil {
		l.last = item
	} else {
		item.next.prev = item
	}
	l.first = item
	l.lenght++
}

//PushBack добавить значение в конец
func (l *List) PushBack(v interface{}) {
	item := &Item{
		value: v,
		prev:  l.Last(),
		list:  l,
	}
	if item.prev == nil {
		l.first = item
	} else {
		item.prev.next = item
	}
	l.last = item
	l.lenght++
}

//Item элемент списка
type Item struct {
	value      interface{} //самы элемента
	next, prev *Item       // ссылки на следующий и предидущий элементы
	list       *List       // ссылка на контейнер в котором находится элемент
}

//Value возвращает значение
func (i *Item) Value() interface{} {
	if i == nil {
		return nil
	}
	return i.value
}

//Next возвращает следующий элемент списка
func (i *Item) Next() *Item {
	if i == nil {
		return nil
	}
	return i.next
}

//Prev возвращает предыдущий элемент списка
func (i *Item) Prev() *Item {
	if i == nil {
		return nil
	}
	return i.prev
}

//Remove удалить элемент из списка
func (i *Item) Remove() {
	if i == nil {
		return
	}
	if i.prev == nil {
		i.list.first = i.next
	} else {
		i.prev.next = i.next
	}
	if i.next == nil {
		i.list.last = i.prev
	} else {
		i.next.prev = i.prev
	}
	i.list.lenght--
}
