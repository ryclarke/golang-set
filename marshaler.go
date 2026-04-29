package mapset

import "go.mongodb.org/mongo-driver/bson/bsontype"

// assert that Marshaler still implements the Set interface and can be used interchangeably with its wrapped Set.
var _ Set[string] = new(Marshaler[string])

// Marshaler is a wrapper type that provides inline JSON and BSON marshaling/unmarshaling capabilities for any Set implementation.
// This is particularly useful for embedding a Set within a struct that needs to be serialized, without having to implement custom marshaling
// methods for the struct itself. This must be a concrete type (i.e. not an interface) so that marshaler can properly invoke the methods.
type Marshaler[T comparable] struct {
	Set[T] // Embed the Set anonymously to allow passthrough access to its methods.
}

func NewMarshaler[T comparable](s Set[T]) Marshaler[T] {
	return Marshaler[T]{Set: s}
}

// MarshalJSON implements the json.Marshaler interface.
func (m Marshaler[T]) MarshalJSON() ([]byte, error) {
	if m.Set == nil {
		return []byte("null"), nil
	}

	return m.Set.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Marshaler[T]) UnmarshalJSON(b []byte) error {
	if m.Set == nil {
		m.Set = NewSet[T]()
	}

	return m.Set.UnmarshalJSON(b)
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (m Marshaler[T]) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if m.Set == nil {
		return bsontype.Array, []byte("[]"), nil
	}

	return m.Set.MarshalBSONValue()
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (m *Marshaler[T]) UnmarshalBSONValue(bt bsontype.Type, b []byte) error {
	if m.Set == nil {
		m.Set = NewSet[T]()
	}

	return m.Set.UnmarshalBSONValue(bt, b)
}
