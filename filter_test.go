package serializer

import (
	"testing"
)

type Test struct {
	Test0 int `group:"test"`
	Test1 int `group:"testo"`
	Test2 int `group:"test"`
	Test3 int `group:"testo,test"`
	Test4 int
	test5 int
	Test6 int `group:"test"`
}

var test = Test{
	9, -8, 7, 6, -5, -4, 3,
}
var testDeserializedResult = Test{
	9, 0, 7, 6, 0, 0, 3,
}

func TestExecute(t *testing.T) {
	s := NewSerializer(JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		panic(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
}