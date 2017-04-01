package jsonfield

import "testing"

func TestJsonTest(t *testing.T) {
	jj := new(A)
	jj.A, jj.B = 111, 222

	if byts, err := Marshal(jj); err != nil || string(byts) != `{"A":111,"B":222}` {
		t.Error(string(byts), err)
	}

	if byts, err := Marshal(jj, "A"); err != nil || string(byts) != `{"A":111}` {
		t.Error(string(byts), err)
	}

	if byts, err := Marshal(jj, "B"); err != nil || string(byts) != `{"B":222}` {
		t.Error(string(byts), err)
	}
}

type A struct {
	A int
	B int
}
