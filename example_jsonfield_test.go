// Package jsonfield_test

package jsonfield_test

import (
	"github.com/mnhkahn/gogogo/logger"
	"github.com/mnhkahn/jsonfield"
)

type A struct {
	A int
	B int
}

func Example() {
	jj := new(A)
	jj.A, jj.B = 111, 222

	buf, _ := jsonfield.Marshal(jj)
	logger.Info(string(buf)) // prints {"A":111,"B":222}

	buf, _ = jsonfield.Marshal(jj, "A")
	logger.Info(string(buf)) // prints {"A":111}
}

func ExampleMarshal() {
	jj := new(A)
	jj.A, jj.B = 111, 222
	byts, err := jsonfield.Marshal(jj, "A") // {"A":111}
	_, _ = byts, err
}

