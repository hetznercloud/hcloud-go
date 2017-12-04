package conversion

import (
	"fmt"
	"testing"
)

func TestDefaultConversion(t *testing.T) {
	type test1 struct {
		Field1    string
		Field2    int
		FieldElse string
	}
	type test2 struct {
		Field1 string
		Field2 int
	}

	t1 := &test1{
		Field1:    "field1",
		Field2:    222,
		FieldElse: "field_else",
	}
	t2 := &test2{}

	c := NewConverter()
	err := c.DefaultConvert(t1, t2)
	fmt.Println(err, t2)
}
