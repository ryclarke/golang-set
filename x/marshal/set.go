// Package marshal implements a wrapper around the mapset Set interface
// providing inline JSON and BSON marshaling/unmarshaling capabilities.
// This is particularly useful for embedding a Set within a struct that
// needs to be serialized, as otherwise the outer struct would need to
// implement custom marshaling logic to handle each Set field.
//
// The zero value *Set is safe for use without initialization and will
// marshal as an empty set (i.e. "[]" in JSON). When unmarshaling, if
// the Set is initialized it will keep the same underlying mapset.Set
// instance and simply clear and repopulate it with the new values. If
// the Set is not initialized, it will be initialized using the default
// thread-safe mapset implementation with the unmarshaled values.
package marshal

import (
	"encoding/json"
	"errors"
	"fmt"

	mapset "github.com/deckarep/golang-set/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var errWrongBSONType = errors.New("cannot unmarshal BSON type into Set")

var (
	// assert compatibility with mapset.Set interface (marshal.Set can be used anywhere a mapset.Set is expected)
	_ mapset.Set[string] = new(Set[string])

	// assert that marshal.Set implements the necessary interfaces for JSON marshaling/unmarshaling.
	_ json.Marshaler   = new(Set[string])
	_ json.Unmarshaler = new(Set[string])

	// assert that marshal.Set implements the necessary interfaces for BSON marshaling/unmarshaling.
	_ bson.ValueMarshaler   = new(Set[string])
	_ bson.ValueUnmarshaler = new(Set[string])
)

// NewSet wraps the given mapset.Set in the marshalable Set wrapper type.
func NewSet[T comparable](set mapset.Set[T]) *Set[T] {
	return &Set[T]{set}
}

// Set is a wrapper type that provides inline JSON and BSON marshaling/unmarshaling capabilities for any mapset.Set implementation.
// This is particularly useful for embedding a Set within a struct that needs to be serialized, without having to implement custom marshaling
// methods for the outer struct. This must be a concrete type (i.e. not an interface) so that the marshaler can properly invoke the methods.
type Set[T comparable] struct {
	mapset.Set[T] // Embed the mapset.Set anonymously to allow passthrough access to its methods.
}

// MarshalJSON implements the json.Marshaler interface.
func (m Set[T]) MarshalJSON() ([]byte, error) {
	if m.Set == nil {
		// Marshal a nil inner set as if it were merely empty
		// to guarantee safe usage of Set's zero value.
		return json.Marshal([]T{})
	}

	return json.Marshal(m.ToSlice())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Set[T]) UnmarshalJSON(b []byte) error {
	var items []T
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	if m.Set == nil {
		m.Set = mapset.NewSet(items...)
	} else {
		m.Clear()
		m.Append(items...)
	}

	return nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (m Set[T]) MarshalBSONValue() (bt byte, b []byte, err error) {
	var btype bson.Type

	if m.Set == nil {
		// Marshal a nil inner set as if it were merely empty
		// to guarantee safe usage of Set's zero value.
		btype, b, err = bson.MarshalValue([]T{})

		return byte(btype), b, err
	}

	btype, b, err = bson.MarshalValue(m.ToSlice())

	return byte(btype), b, err
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (m *Set[T]) UnmarshalBSONValue(bt byte, b []byte) error {
	if bt != byte(bson.TypeArray) {
		// Only unmarshaling from BSON arrays is supported.
		return fmt.Errorf("%w: expected BSON array, got type %v", errWrongBSONType, bt)
	}

	var items []T
	if err := bson.UnmarshalValue(bson.Type(bt), b, &items); err != nil {
		return err
	}

	if m.Set == nil {
		m.Set = mapset.NewSet(items...)
	} else {
		m.Clear()
		m.Append(items...)
	}

	return nil
}
