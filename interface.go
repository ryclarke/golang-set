package mapset

import (
	"encoding/json"
	"fmt"
)

// Set is the primary interface provided by the mapset package. It represents an unordered
// set of data and a large number of operations that can be applied to that set.
type Set[T comparable] interface {
	Comparator[T]
	Mutator[T]
	Builder[T]
}

// Comparator defines the methods for inspecting and comparing sets.
type Comparator[T comparable] interface {
	fmt.Stringer   // String returns a string representation of the set.
	json.Marshaler // MarshalJSON creates a JSON array from the set.

	// Cardinality returns the number of elements in the set.
	Cardinality() int

	// IsEmpty determines if there are elements in the set.
	IsEmpty() bool

	// ToSlice returns the members of the set as a slice.
	ToSlice() []T

	// Contains returns whether the given item is in the set.
	//
	// Contains(All|Any) may cause the argument to escape to the heap.
	// See: https://github.com/deckarep/golang-set/issues/118
	Contains(val T) bool

	// ContainsAll returns whether the given items are all in the set.
	ContainsAll(val ...T) bool

	// ContainsAny returns whether at least one of the given items are in the set.
	ContainsAny(val ...T) bool

	// Equal determines if two sets are equal to each other. If they have the same
	// cardinality and contain the same elements, they are considered equal. The
	// order in which the elements were added is irrelevant.
	//
	// Note that the argument to Equal must be of the same type as the receiver
	// of the method. Otherwise, Equal will panic.
	Equal(other Set[T]) bool

	// Intersects returns whether at least one element in the given set is in the set.
	Intersects(other Set[T]) bool

	// IsSubset determines if every element in this set is in the other set.
	//
	// Note that the argument to IsSubset must be of the same type as the receiver
	// of the method. Otherwise, IsSubset will panic.
	IsSubset(other Set[T]) bool

	// IsProperSubset determines if every element in this set is in
	// the other set but the two sets are not equal.
	//
	// Note that the argument to IsProperSubset must be of the same type
	// as the receiver of the method. Otherwise, IsProperSubset will panic.
	IsProperSubset(other Set[T]) bool

	// IsSuperset determines if every element in the other set is in this set.
	//
	// Note that the argument to IsSuperset must be of the same type as the receiver
	// of the method. Otherwise, IsSuperset will panic.
	IsSuperset(other Set[T]) bool

	// IsProperSuperset determines if every element in the other set
	// is in this set but the two sets are not equal.
	//
	// Note that the argument to IsProperSuperset must be of the same type
	// as the receiver of the method. Otherwise, IsProperSuperset will panic.
	IsProperSuperset(other Set[T]) bool

	// Each iterates over elements and executes the passed func against each element.
	// If passed func returns true, stop iteration at the time.
	Each(func(T) bool)

	// Iter returns a channel of elements that you can
	// range over.
	Iter() <-chan T

	// Iterator returns an Iterator object that you can
	// use to range over the set.
	Iterator() *Iterator[T]
}

// Mutator defines the methods for managing the contents of sets.
type Mutator[T comparable] interface {
	json.Unmarshaler // UnmarshalJSON recreates a set from a JSON array.

	// Add adds an element to the set. Returns whether the item was added.
	Add(val T) bool

	// Append multiple elements to the set. Returns the number of elements added.
	Append(val ...T) int

	// AppendFrom elements from another set into this set. (shorthand of s.Append(other.ToSlice()...))
	// Similar to Union but it modifies the calling set in place instead of returning a new set.
	// Returns the number of elements added.
	AppendFrom(other Set[T]) int

	// Remove removes a single element from the set.
	Remove(i T)

	// RemoveAll removes multiple elements from the set.
	RemoveAll(i ...T)

	// RemoveFrom removes elements from another set from this set. (shorthand of s.RemoveAll(other.ToSlice()...))
	// Similar to Difference but it modifies the calling set in place instead of returning a new set.
	RemoveFrom(other Set[T])

	// Pop removes and returns an arbitrary item from the set.
	Pop() (T, bool)

	// PopN removes and returns up to n arbitrary items from the set.
	// It returns a slice of the removed items and the actual number of items removed.
	// If the set is empty or n is less than or equal to 0s, it returns an empty slice and 0.
	// If n is greater than the set's size, all items are removed and returned.
	PopN(n int) ([]T, int)

	// Clear removes all elements from the set, leaving the empty set.
	Clear()
}

// Builder defines the methods for creating new sets from an existing set.
type Builder[T comparable] interface {
	// Clone returns a clone of the set using the same
	// implementation, duplicating all keys.
	Clone() Set[T]

	// Difference returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not also
	// elements of other.
	//
	// Note that the argument to Difference
	// must be of the same type as the receiver
	// of the method. Otherwise, Difference will
	// panic.
	Difference(other Set[T]) Set[T]

	// SymmetricDifference returns a new set with all elements which are
	// in either this set or the other set but not in both.
	//
	// Note that the argument to SymmetricDifference
	// must be of the same type as the receiver
	// of the method. Otherwise, SymmetricDifference
	// will panic.
	SymmetricDifference(other Set[T]) Set[T]

	// Union returns a new set with all elements in both sets.
	//
	// Note that the argument to Union must be of the
	// same type as the receiver of the method.
	// Otherwise, Union will panic.
	Union(other Set[T]) Set[T]

	// Intersect returns a new set containing only the elements
	// that exist only in both sets.
	//
	// Note that the argument to Intersect
	// must be of the same type as the receiver
	// of the method. Otherwise, Intersect will
	// panic.
	Intersect(other Set[T]) Set[T]

	// Filter iterates over elements and executes the passed func against each element.
	// If passed func returns true, the element will be added to the returned set.
	Filter(func(T) bool) Set[T]
}
