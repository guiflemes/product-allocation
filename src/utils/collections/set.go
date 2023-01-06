package collections

type Set[T comparable] map[T]bool

func NewSet[T comparable](size int) Set[T] {
	return make(Set[T], size)
}

func (s Set[T]) Add(values ...T) {
	for _, value := range values {
		s[value] = true
	}
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	return s[value]
}

func (a Set[T]) Union(b Set[T]) Set[T] {
	small, large := smallLarge(a, b)

	for value := range small {
		large.Add(value)
	}
	return large
}

func (a Set[T]) Difference(b Set[T]) Set[T] {
	resultSet := NewSet[T](0)
	for value := range a {
		if !b.Contains(value) {
			resultSet.Add(value)
		}
	}
	return resultSet
}

func (a Set[T]) Intersection(b Set[T]) Set[T] {
	small, large := smallLarge(a, b)

	resultSet := NewSet[T](0)
	for value := range small {
		if large.Contains(value) {
			resultSet.Add(value)
		}
	}
	return resultSet
}

func smallLarge[T comparable](a, b Set[T]) (Set[T], Set[T]) {
	small, large := b, a
	if len(b) > len(a) {
		small, large = a, b
	}

	return small, large
}

func (a Set[T]) Equals(b Set[T]) bool {
	return len(a.Difference(b)) == 0 && len(b.Difference(a)) == 0
}

func (a Set[T]) ToSlice() []T {
	var slice []T

	for v := range a {
		slice = append(slice, v)
	}

	return slice
}
