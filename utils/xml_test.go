package utils

import "testing"

type Person struct {
	Name   string `xml:"name"`
	Age    int16  `xml:"age"`
	Gender uint8  `xml:"gender"`
	Teachers P
}

type P struct {
	Name   string
	Age    int16
	Gender uint8
	Person []*Person
}

func TestMarshal(t *testing.T) {
	d := Person{
		Name:   "A",
		Age:    21,
		Gender: 1,
		Teachers: P{
			"B",
			25,
			1,
			[]*Person{&Person{
				"C",
				22,
				0,
				P{}}}}}

	b, err := Marshal(d)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}
