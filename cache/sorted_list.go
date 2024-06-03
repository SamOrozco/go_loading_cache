package cache

// SortedList is a list that is sorted as you add items to the list it will be ordered by the given comparing function
type SortedList[K any] interface {
	Add(k K)
	GetAll() []K
	GetFirstN(num int) []K
}

type node[K any] struct {
	value K
	next  *node[K]
}

type LinkedSortedList[K any] struct {
	// CompareFunc the function to use to compare the items in the list
	// if the response is less than 0 the left item is less than the right
	// if the response is 0 the items are equal
	// if the response is greater than 0 the left item is greater than the right
	CompareFunc func(left K, right K) int

	// root node
	root *node[K]
}

func (s LinkedSortedList[K]) Add(k K) {
	newNode := &node[K]{value: k}

	// no root then set the root to the new node
	if s.root == nil {
		s.root = newNode
		return
	}

	// if our value is less than zero we insert our value at the beginning and set the root to the new node
	if s.CompareFunc(k, s.root.value) < 0 {
		newNode.next = s.root
		s.root = newNode
		return
	}

	// insert where we are less than a value
	current := s.root
	for current.next != nil && s.CompareFunc(k, current.next.value) > 0 {
		current = current.next
	}
	newNode.next = current.next
	current.next = newNode
}

func (s LinkedSortedList[K]) GetAll() []K {
	response := make([]K, 0)
	current := s.root
	for current != nil {
		response = append(response, current.value)
		current = current.next
	}
	return response
}

func (s LinkedSortedList[K]) GetFirstN(idx int) []K {
	response := make([]K, 0)
	current := s.root
	for current != nil && idx > 0 {
		response = append(response, current.value)
		current = current.next
		idx--
	}
	return response
}
