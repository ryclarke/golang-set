/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2013 - 2022 Ralph Caraveo (deckarep@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Package mapset implements a simple and  set collection.
// Items stored within it are unordered and unique. It supports
// typical set operations: membership testing, intersection, union,
// difference, symmetric difference and cloning.
//
// Package mapset provides two implementations of the Set
// interface. The default implementation is safe for concurrent
// access, but a non-thread-safe implementation is also provided for
// programs that can benefit from the slight speed improvement and
// that can enforce mutual exclusion through other means.
package mapset

// NewSet creates and returns a new set with the given elements.
// Operations on the resulting set are thread-safe.
func NewSet[T comparable](vs ...T) Set[T] {
	s := newThreadSafeSetWithSize[T](len(vs))
	s.uss.append(vs...)
	return s
}

// NewSetWithSize creates and returns a reference to an empty set with a specified
// capacity. Operations on the resulting set are thread-safe.
func NewSetWithSize[T comparable](cardinality int) Set[T] {
	return newThreadSafeSetWithSize[T](cardinality)
}

// NewThreadUnsafeSet creates and returns a new set with the given elements.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeSet[T comparable](vs ...T) Set[T] {
	s := newThreadUnsafeSetWithSize[T](len(vs))
	s.append(vs...)
	return s
}

// NewThreadUnsafeSetWithSize creates and returns a reference to an empty set with
// a specified capacity. Operations on the resulting set are not thread-safe.
func NewThreadUnsafeSetWithSize[T comparable](cardinality int) Set[T] {
	return newThreadUnsafeSetWithSize[T](cardinality)
}

// NewSetFromMapKeys creates and returns a new set with the given keys of the map.
// Operations on the resulting set are thread-safe.
func NewSetFromMapKeys[T comparable, V any](val map[T]V) Set[T] {
	s := newThreadSafeSetWithSize[T](len(val))

	for k := range val {
		s.uss.add(k)
	}

	return s
}

// NewThreadUnsafeSetFromMapKeys creates and returns a new set with the given keys of the map.
// Operations on the resulting set are not thread-safe.
func NewThreadUnsafeSetFromMapKeys[T comparable, V any](val map[T]V) Set[T] {
	s := newThreadUnsafeSetWithSize[T](len(val))

	for k := range val {
		s.add(k)
	}

	return s
}

// Elements returns an iterator that yields the elements of the set. Starting
// with Go 1.23, users can use a for loop to iterate over it.
func Elements[T comparable](s Set[T]) func(func(element T) bool) {
	return func(yield func(element T) bool) {
		s.Each(func(t T) bool {
			return !yield(t)
		})
	}
}
