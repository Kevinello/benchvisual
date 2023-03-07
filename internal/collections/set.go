package collections

// Set A collection of unique comparable items. Uses a map with only true values
// to accomplish set functionality.
//
//	@author kevineluo
//	@update 2023-03-01 03:52:15
type Set[T comparable] map[T]struct{}

// NewSet Create a new empty set with the specified initial size.
//
//	@param size int
//	@return Set
//	@author kevineluo
//	@update 2023-03-01 03:52:21
func NewSet[T comparable](size int) Set[T] {
	return make(Set[T], size)
}

// Add Add a new key to the set
//
//	@param s Set[T]
//	@return Add
//	@author kevineluo
//	@update 2023-03-01 03:52:26
func (s Set[T]) Add(key T) {
	s[key] = struct{}{}
}

// Remove Remove a key from the set. If the key is not in the set then noop
//
//	@param s Set[T]
//	@return Remove
//	@author kevineluo
//	@update 2023-03-01 03:52:28
func (s Set[T]) Remove(key T) {
	delete(s, key)
}

// Contains Check if Set s contains key
//
//	@param s Set[T]
//	@return Contains
//	@author kevineluo
//	@update 2023-03-01 03:52:29
func (s Set[T]) Contains(key T) bool {
	_, exist := s[key]
	return exist
}

// Difference A difference B
// NOTE: A-B != B-A
//
//	@param a Set[T]
//	@return Difference
//	@author kevineluo
//	@update 2023-03-01 03:52:42
func (s Set[T]) Difference(b Set[T]) Set[T] {
	resultSet := NewSet[T](0)
	for key := range s {
		if !b.Contains(key) {
			resultSet.Add(key)
		}
	}
	return resultSet
}

// Union A union B
//
//	@param a Set[T]
//	@return Union
//	@author kevineluo
//	@update 2023-03-01 03:54:43
func (s Set[T]) Union(b Set[T]) Set[T] {
	small, large := smallLarge(s, b)

	for key := range small {
		large.Add(key)
	}
	return large
}

// Intersection A intersect B
//
//	@param s Set[T]
//	@return Intersection
//	@author kevineluo
//	@update 2023-03-01 03:54:59
func (s Set[T]) Intersection(b Set[T]) Set[T] {
	small, large := smallLarge(s, b)

	resultSet := NewSet[T](0)
	for key := range small {
		if large.Contains(key) {
			resultSet.Add(key)
		}
	}
	return resultSet
}

// returns the small and large sets according to their len
func smallLarge[T comparable](a, b Set[T]) (Set[T], Set[T]) {
	small, large := b, a
	if len(b) > len(a) {
		small, large = a, b
	}

	return small, large
}

// ToSlice Turn a Set into a slice
//
//	@param s Set[T]
//	@return ToSlice
//	@author kevineluo
//	@update 2023-03-01 04:43:18
func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for key := range s {
		slice = append(slice, key)
	}

	return slice
}

// Equals A == B (all elements of A are in B and vice versa)
//
//	@param a Set[T]
//	@return Equals
//	@author kevineluo
//	@update 2023-03-01 04:43:01
func (s Set[T]) Equals(b Set[T]) bool {
	return len(s.Difference(b)) == 0 && len(b.Difference(s)) == 0
}

// -------------------------------------------------
// SLICE HELPERS

// SliceToSet Create a Set from a slice.
//
//	@param s []T
//	@return Set
//	@author kevineluo
//	@update 2023-03-01 04:42:55
func SliceToSet[T comparable](s []T) Set[T] {
	set := NewSet[T](len(s))
	for _, item := range s {
		set.Add(item)
	}
	return set
}

// MapSliceToSet Map a slice to a set using a function f
//
//	@param s []S
//	@param f func(s S) T
//	@return Set
//	@author kevineluo
//	@update 2023-03-01 04:42:04
func MapSliceToSet[S any, T comparable](s []S, f func(s S) T) Set[T] {
	set := NewSet[T](len(s))
	for _, item := range s {
		set.Add(f(item))
	}
	return set
}

// SliceUnion Union two slices, The provided slices do not need to be unique. Order not guaranteed.
//
//	@param a []T
//	@param b []T
//	@return []T
//	@author kevineluo
//	@update 2023-03-01 04:42:00
func SliceUnion[T comparable](a, b []T) []T {
	aSet, bSet := SliceToSet(a), SliceToSet(b)
	union := aSet.Union(bSet)
	return union.ToSlice()
}

// SliceIntersection Intersection of two slices, The provided slices do not need to be unique. Order not guaranteed.
//
//	@param a []T
//	@param b []T
//	@return []T
//	@author kevineluo
//	@update 2023-03-01 04:41:58
func SliceIntersection[T comparable](a, b []T) []T {
	aSet, bSet := SliceToSet(a), SliceToSet(b)
	intersection := aSet.Intersection(bSet)
	return intersection.ToSlice()
}

// SliceDifference Difference of two slices(A-B). Slices do not need to be unique. Order not guaranteed.
//
//	@param a []T
//	@param b []T
//	@return []T
//	@author kevineluo
//	@update 2023-03-01 04:41:53
func SliceDifference[T comparable](a, b []T) []T {
	aSet, bSet := SliceToSet(a), SliceToSet(b)
	difference := aSet.Difference(bSet)
	return difference.ToSlice()
}

// SliceEqual judge if two slice is equal(unique element)
//
//	@param a []T
//	@param b []T
//	@return bool
//	@author kevineluo
//	@update 2023-03-01 04:46:22
func SliceEqual[T comparable](a, b []T) bool {
	aSet, bSet := SliceToSet(a), SliceToSet(b)
	return aSet.Equals(bSet)
}
