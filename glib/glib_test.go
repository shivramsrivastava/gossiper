package glib

import (
	"fmt"
	"testing"
)

//A dummy fucntion that will create a Glib and override the ML to Fake MemberList for to test easily
func CreateGlibforTest() *Glib {

	g := NewGlib("test", "north", false)

	g.list = &Fake_memberlist{}

	return g
}

func Test_Init(T *testing.T) {

	g := CreateGlibforTest()

	err := g.Init()

	if err != nil {
		T.Fail()
	}
}
