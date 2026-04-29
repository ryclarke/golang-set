package marshal

import (
	"bytes"
	"encoding/json"
	"testing"

	mapset "github.com/deckarep/golang-set/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type testStruct struct {
	Str string `bson:"str"           json:"str"`
	Int int    `bson:"int,omitempty" json:"int,omitempty"`
}

type wrapper struct {
	Field     string           `bson:"field"   json:"field"`
	StringSet *Set[string]     `bson:"strings" json:"strings"`
	StructSet *Set[testStruct] `bson:"structs" json:"structs"`
}

func TestSet_MarshalUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input *Set[string] // initial input value
		want  []byte       // marshaled JSON output
		want2 *Set[string] // final value after unmarshal round-trip
	}{
		{
			name:  "marshal nil set (container zero value)",
			input: new(Set[string]),
			want:  []byte("[]"),
			want2: NewSet(mapset.NewSet[string]()),
		},
		{
			name:  "marshal initialized empty set",
			input: NewSet(mapset.NewSet[string]()),
			want:  []byte("[]"),
			want2: NewSet(mapset.NewSet[string]()),
		},
		{
			name:  "marshal non-empty set",
			input: NewSet(mapset.NewSet("value")),
			want:  []byte("[\"value\"]"),
			want2: NewSet(mapset.NewSet("value")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("json.Marshal got %s, want %s", string(got), string(tt.want))
			}

			var got2 Set[string]
			if err := json.Unmarshal(got, &got2); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if !tt.want2.Equal(got2.Set) {
				t.Errorf("json.Unmarshal got %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestSet_MarshalUnmarshalJSON_Wrapper(t *testing.T) {
	tests := []struct {
		name  string
		input wrapper // initial input value
		want  []byte  // marshaled JSON output
		want2 wrapper // final value after unmarshal round-trip
	}{
		{
			name:  "marshal struct zero value",
			input: wrapper{},
			want:  []byte(`{"field":"","strings":null,"structs":null}`),
			want2: wrapper{},
		},
		{
			name: "marshal initialized empty set",
			input: wrapper{
				StringSet: NewSet(mapset.NewSet[string]()),
				StructSet: NewSet(mapset.NewSet[testStruct]()),
			},
			want: []byte(`{"field":"","strings":[],"structs":[]}`),
			want2: wrapper{
				StringSet: NewSet(mapset.NewSet[string]()),
				StructSet: NewSet(mapset.NewSet[testStruct]()),
			},
		},
		{
			name: "marshal non-empty set",
			input: wrapper{
				Field:     "value",
				StringSet: NewSet(mapset.NewSet("setValue")),
				StructSet: NewSet(mapset.NewSet(testStruct{Str: "structValue", Int: 42})),
			},
			want: []byte(`{"field":"value","strings":["setValue"],"structs":[{"str":"structValue","int":42}]}`),
			want2: wrapper{
				Field:     "value",
				StringSet: NewSet(mapset.NewSet("setValue")),
				StructSet: NewSet(mapset.NewSet(testStruct{Str: "structValue", Int: 42})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("json.Marshal got %s, want %s", string(got), string(tt.want))
			}

			got2 := wrapper{}
			if err := json.Unmarshal(got, &got2); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if (tt.want2.StringSet != nil && !tt.want2.StringSet.Equal(got2.StringSet.Set)) ||
				(tt.want2.StructSet != nil && !tt.want2.StructSet.Equal(got2.StructSet.Set)) ||
				tt.want2.Field != got2.Field {
				t.Errorf("json.Unmarshal got %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestSet_MarshalUnmarshalBSON(t *testing.T) {
	tests := []struct {
		name  string
		input *Set[string] // initial input value
		want  []byte       // marshaled BSON output
		want2 *Set[string] // final value after unmarshal round-trip
	}{
		{
			name:  "marshal nil set (container zero value)",
			input: new(Set[string]),
			want: []byte{
				5, 0, 0, 0, 0,
			}, // BSON array with zero elements
			want2: NewSet(mapset.NewSet[string]()),
		},
		{
			name:  "marshal initialized empty set",
			input: NewSet(mapset.NewSet[string]()),
			want: []byte{
				5, 0, 0, 0, 0,
			}, // BSON array with zero elements
			want2: NewSet(mapset.NewSet[string]()),
		},
		{
			name:  "marshal non-empty set",
			input: NewSet(mapset.NewSet("value")),
			want: []byte{
				18, 0, 0, 0, 2, 48, 0, 6, 0, 0, 0, 'v', 'a', 'l', 'u', 'e', 0, 0,
			}, // BSON array with one string element "value"
			want2: NewSet(mapset.NewSet("value")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt, got, err := bson.MarshalValue(tt.input)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("bson.MarshalValue got %v, want %v", got, tt.want)
			}

			if bt != bson.TypeArray {
				t.Errorf("bson.MarshalValue got type %v, want %v", bt, bson.TypeArray)
			}

			var got2 Set[string]
			if err := bson.UnmarshalValue(bt, got, &got2); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if !tt.want2.Equal(got2.Set) {
				t.Errorf("bson.UnmarshalValue got %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestSet_MarshalUnmarshalBSON_Wrapper(t *testing.T) {
	tests := []struct {
		name  string
		input wrapper // initial input value
		want  []byte  // marshaled BSON output
		want2 wrapper // final value after unmarshal round-trip
	}{
		{
			name:  "marshal struct zero value",
			input: wrapper{},
			want: []byte{
				35, 0, 0, 0,
				2, 'f', 'i', 'e', 'l', 'd', 0, 1, 0, 0, 0, 0,
				10, 's', 't', 'r', 'i', 'n', 'g', 's', 0,
				10, 's', 't', 'r', 'u', 'c', 't', 's', 0,
				0,
			}, // BSON representation of {"field":"","strings":null,"structs":null}
			want2: wrapper{},
		},
		{
			name: "marshal initialized empty set",
			input: wrapper{
				StringSet: NewSet(mapset.NewSet[string]()),
				StructSet: NewSet(mapset.NewSet[testStruct]()),
			},
			want: []byte{
				45, 0, 0, 0,
				2, 'f', 'i', 'e', 'l', 'd', 0, 1, 0, 0, 0, 0,
				4, 's', 't', 'r', 'i', 'n', 'g', 's', 0, 5, 0, 0, 0, 0,
				4, 's', 't', 'r', 'u', 'c', 't', 's', 0, 5, 0, 0, 0, 0, 0,
			}, // BSON representation of {"field":"","strings":[],"structs":[]}
			want2: wrapper{
				StringSet: NewSet(mapset.NewSet[string]()),
				StructSet: NewSet(mapset.NewSet[testStruct]()),
			},
		},
		{
			name: "marshal non-empty set",
			input: wrapper{
				Field:     "value",
				StringSet: NewSet(mapset.NewSet("setValue")),
				StructSet: NewSet(mapset.NewSet(testStruct{Str: "structValue", Int: 42})),
			},
			want: []byte{
				104, 0, 0, 0,
				2, 'f', 'i', 'e', 'l', 'd', 0, 6, 0, 0, 0, 'v', 'a', 'l', 'u', 'e', 0,
				4, 's', 't', 'r', 'i', 'n', 'g', 's', 0, 21, 0, 0, 0,
				2, '0', 0, 9, 0, 0, 0, 's', 'e', 't', 'V', 'a', 'l', 'u', 'e', 0, 0,
				4, 's', 't', 'r', 'u', 'c', 't', 's', 0, 43, 0, 0, 0,
				3, '0', 0, 35, 0, 0, 0,
				2, 's', 't', 'r', 0, 12, 0, 0, 0, 's', 't', 'r', 'u', 'c', 't', 'V', 'a', 'l', 'u', 'e', 0,
				16, 'i', 'n', 't', 0, 42, 0, 0, 0, 0, 0, 0,
			}, // BSON representation of {"field":"value","strings":["setValue"],"structs":[{"str":"structValue","int":42}]}
			want2: wrapper{
				Field:     "value",
				StringSet: NewSet(mapset.NewSet("setValue")),
				StructSet: NewSet(mapset.NewSet(testStruct{Str: "structValue", Int: 42})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt, got, err := bson.MarshalValue(tt.input)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("bson.MarshalValue got %v, want %v", got, tt.want)
			}

			got2 := wrapper{}
			if err := bson.UnmarshalValue(bt, got, &got2); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			if (tt.want2.StringSet != nil && !tt.want2.StringSet.Equal(got2.StringSet.Set)) ||
				(tt.want2.StructSet != nil && !tt.want2.StructSet.Equal(got2.StructSet.Set)) ||
				tt.want2.Field != got2.Field {
				t.Errorf("bson.UnmarshalValue got %v, want %v", got2, tt.want2)
			}
		})
	}
}
