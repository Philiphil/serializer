package serializer

import (
	"fmt"
	"testing"
)

func TestSerializer_Deserialize(t *testing.T) {
	s := NewSerializer(JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		panic(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
	if o != testDeserializedResult {
		panic("!")
	}

}

func TestSerializer_Deserialize2(t *testing.T) {
	test2 := Recursive{
		Hidden{1, 2},
		Hidden{3, 4},
	}
	expected2 := Recursive{
		Hidden{1, 0},
		Hidden{0, 0},
	}

	s := NewSerializer(JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		panic(err)
	}
	o := Recursive{}
	err = s.Deserialize(serialized, &o)
	if o != expected2 {
		fmt.Println(o)
		fmt.Println(expected2)
		panic("!")
	}

}

func TestSerializer_Deserialize3(t *testing.T) {
	test2 := []Recursive{
		{
			Hidden{1, 2},
			Hidden{3, 4},
		},
	}
	expected2 := []Recursive{
		{
			Hidden{1, 0},
			Hidden{0, 0},
		},
	}
	s := NewSerializer(JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		panic(err)
	}
	o := []Recursive{}
	err = s.Deserialize(serialized, &o)
	if o[0] != expected2[0] {
		fmt.Println(o)
		fmt.Println(expected2)
		panic("!")
	}

}

func TestSerializer_MergeObjects(t *testing.T) {
	target := Test{
		11, 11, 11, 11, 11, 11, 11,
	}
	result := Test{
		9, 11, 7, 6, 11, 11, 3,
	}
	s := NewSerializer(JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		panic(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
	err = s.MergeObjects(&target, &o)
	if err != nil {
		panic(err)
	}
	if target != result {
		fmt.Println(target)
		fmt.Println(result)
		panic("!")
	}

}

func TestSerializer_DeserializeAndMerge(t *testing.T) {
	target := Test{
		11, 11, 11, 11, 11, 11, 11,
	}
	result := Test{
		9, 11, 7, 6, 11, 11, 3,
	}
	s := NewSerializer(JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		panic(err)
	}
	err = s.DeserializeAndMerge(serialized, &target)
	if err != nil {
		panic(err)
	}
	if target != result {
		fmt.Println(target)
		fmt.Println(result)
		panic("!")
	}

}

type MyStruct struct {
	Name  string `json:"name" group:"group1"`
	Age   int    `json:"age" group:"group2"`
	Email string `json:"email"`
}
