package ds

// Set is a collection of unique elements.
// It is implemented using a map with the element type as the key and an empty struct as the value.
// It is not concurrency-safe, so consider wrapping it with a sync.RWMutex if you need to use it concurrently.
type Set[E comparable] map[E]struct{}

func NewSet[E comparable](items ...E) Set[E] {
	set := make(Set[E])
	for _, v := range items {
		set.Add(v)
	}
	return set
}

func (s Set[E]) Add(value E) {
	s[value] = struct{}{}
}

func (s Set[E]) Remove(value E) {
	delete(s, value)
}

func (s Set[E]) Contains(value E) bool {
	_, contains := s[value]
	return contains
}
